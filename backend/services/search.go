// File ini adalah service utama untuk pencarian resep di backend Little Alchemy 2.
// Berisi inisialisasi database, pembacaan mapper, dan implementasi algoritma BFS, DFS, dan Bidirectional.

package services

import (
  "database/sql"      // Untuk koneksi database SQLite
  "encoding/json"     // Untuk decode file JSON
  "log"               // Untuk logging error
  "os"                // Untuk akses file
  "time"              // Untuk pengukuran waktu eksekusi
  _ "github.com/mattn/go-sqlite3" // Driver SQLite
)

var db *sql.DB                    // Variabel global untuk koneksi database
var mapper map[string]string      // Variabel global untuk mapping nama elemen ke gambar/icon. jadi gambarnya dalam string dan nama elemen dalam string juga

func init() {
  var err error
  db, err = sql.Open("sqlite3", "../database/alchemy.db") // Buka database SQLite
  
  if err != nil {
    log.Fatalf("Gagal membuka database: %v", err)
  } else {
    log.Printf("Database ditemukan")
  }

  file, err := os.Open("../database/mapper2.json") // Buka file mapper2.json
  if err != nil {
    log.Fatalf("Gagal membuka mapper.json: %v", err)
  }
  defer file.Close() // Tutup file setelah selesai
  if err := json.NewDecoder(file).Decode(&mapper); err != nil {
    log.Fatalf("Gagal mendekode mapper.json: %v", err) // Jika gagal decode, log error
  }
}

type Node struct {
  Name     string      // Nama elemen
  Children []*Node     // Anak-anak node (hasil kombinasi)
}

// Mengambil anak-anak dari elemen tertentu dari database
func getChildren(name string) ([]*Node, error) { // cari anak dari elemen
  rows, err := db.Query("SELECT DISTINCT element FROM elements WHERE item1 = ? OR item2 = ?", name, name) // Ambil elemen yang memiliki item1 atau item2 sama dengan nama elemen
// Dengan tabel yang sama seperti di atas, jika name = "Water", maka query akan mencari:


  if err != nil {
    return nil, err // Jika gagal query, kembalikan error
  }
  defer rows.Close() // Tutup rows setelah selesai

  var children []*Node // Slice untuk menyimpan anak-anak node
  for rows.Next() { // Iterasi setiap hasil query
    var childName string
    if err := rows.Scan(&childName); err != nil {  // kalo gagal scan, kembalikan error
      return nil, err 
    }
    children = append(children, &Node{Name: childName}) // kalo berhasil, tambahkan ke slice children
  }
  return children, nil // Kembalikan slice children pas udah selesai
}

//KODE INI DIUBAH (ATAS)
// RecipeStep represents a single step in a recipe creation
type RecipeStep struct {
  Result string  // Element yang dihasilkan
  Item1  string  // Bahan pertama
  Item2  string  // Bahan kedua
}

// Improved BFS untuk pencarian resep yang lebih lengkap
func BFS(elementName string, recipeType string, maxRecipes int) ([]interface{}, int, float64) {
  start := time.Now() // Catat waktu mulai eksekusi
  visited := make(map[string]bool) // Map untuk menandai node yang sudah dikunjungi
  nodesVisited := 0 // Counter node yang dikunjungi
  
  // Cek apakah elemen ada di database
  var exists bool
  err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM elements WHERE element = ?)", elementName).Scan(&exists)
  if err != nil || !exists {
    return []interface{}{}, nodesVisited, float64(time.Since(start).Milliseconds())
  }
  
  // Queue untuk BFS dengan tracking recipe path
  type QueueItem struct {
    Name string         // Nama elemen
    Path []RecipeStep   // Path resep dari elemen ini
  }
  queue := []QueueItem{{Name: elementName, Path: []RecipeStep{}}}
  
  // Hasil pencarian
  results := []interface{}{}
  
  // Cari elemen dasar (basic elements)
  basicElements := []string{}
  rows, err := db.Query("SELECT DISTINCT element FROM elements WHERE element NOT IN (SELECT DISTINCT element FROM elements WHERE item1 IS NOT NULL AND item2 IS NOT NULL)")
  if err == nil {
    defer rows.Close()
    for rows.Next() {
      var element string
      if err := rows.Scan(&element); err == nil {
        basicElements = append(basicElements, element)
      }
    }
  }
  
  // Process queue sampai kosong atau max recipes tercapai
  for len(queue) > 0 && (recipeType != "Limit" || len(results) < maxRecipes) {
    // Ambil item pertama dari queue (FIFO)
    current := queue[0]
    queue = queue[1:]
    
    // Skip jika sudah dikunjungi
    if visited[current.Name] {
      continue
    }
    
    // Tandai sudah dikunjungi dan tambah counter
    visited[current.Name] = true
    nodesVisited++
    
    // Cek apakah ini elemen dasar
    isBasic := false
    for _, basic := range basicElements {
      if current.Name == basic {
        isBasic = true
        break
      }
    }
    
    // Jika ini elemen dasar dan bukan target awal, tambahkan ke hasil dan lanjutkan
    if isBasic && current.Name != elementName {
      // Format node untuk response
      childrenData := []interface{}{}
      
      results = append(results, map[string]interface{}{
        "name":     current.Name,
        "image":    mapper[current.Name],
        "children": childrenData,
        "recipe":   "Basic element",
      })
      
      continue
    }
    
    // Jika ini target elemen dan kita memiliki recipe path, tambahkan ke hasil
    if current.Name == elementName && len(current.Path) > 0 {
      // Format recipe untuk response
      recipeTree := buildRecipeTree(current.Path)
      recipeSteps := formatRecipeSteps(current.Path)
      
      results = append(results, map[string]interface{}{
        "name":     current.Name,
        "image":    mapper[current.Name],
        "children": recipeTree,
        "recipe":   recipeSteps,
      })
      
      // Jika sudah mencapai max recipes, break
      if recipeType == "Limit" && len(results) >= maxRecipes {
        break
      }
    }
    
    // Dapatkan parent elements yang bisa dikombinasikan untuk membuat elemen ini
    rows, err := db.Query("SELECT item1, item2 FROM elements WHERE element = ? AND item1 IS NOT NULL AND item2 IS NOT NULL", current.Name)
    if err != nil {
      continue
    }
    
    // Process setiap kombinasi parent
    hasParents := false
    for rows.Next() {
      var item1, item2 string
      if err := rows.Scan(&item1, &item2); err != nil {
        continue
      }
      
      hasParents = true
      
      // Buat recipe step baru
      newStep := RecipeStep{
        Result: current.Name,
        Item1:  item1,
        Item2:  item2,
      }
      
      // Tambahkan parent1 ke queue dengan path yang diupdate
      if !visited[item1] {
        newPath1 := make([]RecipeStep, len(current.Path))
        copy(newPath1, current.Path)
        newPath1 = append(newPath1, newStep)
        queue = append(queue, QueueItem{Name: item1, Path: newPath1})
      }
      
      // Tambahkan parent2 ke queue dengan path yang diupdate
      if !visited[item2] {
        newPath2 := make([]RecipeStep, len(current.Path))
        copy(newPath2, current.Path)
        newPath2 = append(newPath2, newStep)
        queue = append(queue, QueueItem{Name: item2, Path: newPath2})
      }
    }
    rows.Close()
    
    // Jika node ini tidak punya parents, tampilkan sebagai "no recipe"
    if !hasParents && current.Name != elementName {
      childData := []interface{}{}
      
      results = append(results, map[string]interface{}{
        "name":     current.Name,
        "image":    mapper[current.Name],
        "children": childData,
        "recipe":   "No parents found",
      })
    }
    
    // Ambil children untuk node ini (untuk tree display)
    children, err := getChildren(current.Name)
    if err == nil && len(children) > 0 {
      // Format children untuk response
      childData := make([]map[string]interface{}, len(children))
      for i, child := range children {
        childData[i] = map[string]interface{}{
          "name":  child.Name,
          "image": mapper[child.Name],
          "children": []interface{}{},
        }
      }
      
      // Tambahkan node dengan children ke hasil
      if current.Name == elementName && len(results) == 0 {
        results = append(results, map[string]interface{}{
          "name":     current.Name,
          "image":    mapper[current.Name],
          "children": childData,
        })
      }
    }
    
    time.Sleep(20 * time.Millisecond) // Delay kecil untuk visualisasi
  }
  
  // Jika tidak ada results dan ini adalah elemen target, tambahkan info elemen
  if len(results) == 0 {
    children, _ := getChildren(elementName)
    childData := make([]map[string]interface{}, 0)
    if children != nil {
      childData = make([]map[string]interface{}, len(children))
      for i, child := range children {
        childData[i] = map[string]interface{}{
          "name":  child.Name,
          "image": mapper[child.Name],
          "children": []interface{}{},
        }
      }
    }
    
    results = append(results, map[string]interface{}{
      "name":     elementName,
      "image":    mapper[elementName],
      "children": childData,
      "recipe":   "This is a basic element or no recipe found",
    })
  }
  
  return results, nodesVisited, float64(time.Since(start).Milliseconds())
}

// Helper function untuk membangun pohon resep
func buildRecipeTree(path []RecipeStep) []map[string]interface{} {
  if len(path) == 0 {
    return []map[string]interface{}{}
  }
  
  // Ambil step terakhir karena kita bekerja mundur dari elemen target
  lastStep := path[len(path)-1]
  
  // Cari step dimana item1 dan item2 adalah hasil
  var item1Steps, item2Steps []RecipeStep
  for _, step := range path {
    if step.Result == lastStep.Item1 {
      item1Steps = append(item1Steps, step)
    }
    if step.Result == lastStep.Item2 {
      item2Steps = append(item2Steps, step)
    }
  }
  
  // Buat subpath untuk item1 dan item2
  var item1Path, item2Path []RecipeStep
  for _, step := range path[:len(path)-1] {
    if step.Result == lastStep.Item1 || containsElement(item1Path, step.Result) {
      item1Path = append(item1Path, step)
    }
    if step.Result == lastStep.Item2 || containsElement(item2Path, step.Result) {
      item2Path = append(item2Path, step)
    }
  }
  
  // Buat children nodes
  result := []map[string]interface{}{}
  
  // Tambahkan node untuk item1
  item1Node := map[string]interface{}{
    "name":  lastStep.Item1,
    "image": mapper[lastStep.Item1],
  }
  if len(item1Steps) > 0 {
    item1Node["children"] = buildRecipeTree(item1Path)
  } else {
    item1Node["children"] = []interface{}{}
  }
  result = append(result, item1Node)
  
  // Tambahkan node untuk item2
  item2Node := map[string]interface{}{
    "name":  lastStep.Item2,
    "image": mapper[lastStep.Item2],
  }
  if len(item2Steps) > 0 {
    item2Node["children"] = buildRecipeTree(item2Path)
  } else {
    item2Node["children"] = []interface{}{}
  }
  result = append(result, item2Node)
  
  return result
}

// Helper function untuk mengecek apakah sebuah slice mengandung elemen
func containsElement(steps []RecipeStep, element string) bool {
  for _, step := range steps {
    if step.Result == element {
      return true
    }
  }
  return false
}

// Helper function untuk format recipe steps
func formatRecipeSteps(path []RecipeStep) []string {
  if len(path) == 0 {
    return []string{}
  }
  
  // Buat recipe steps dalam urutan terbalik (dari basic ke target)
  steps := make([]string, len(path))
  for i, step := range path {
    steps[len(path)-i-1] = step.Result + " = " + step.Item1 + " + " + step.Item2
  }
  
  return steps
}

func DFS(elementName string, recipeType string, maxRecipes int) ([]interface{}, int, float64) {
  start := time.Now() // Catat waktu mulai eksekusi
  visited := make(map[string]bool) // Map untuk menandai node yang sudah dikunjungi
  nodesVisited := 0 // Counter node yang dikunjungi
  
  // Cek apakah elemen ada di database
  var exists bool
  err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM elements WHERE element = ?)", elementName).Scan(&exists)
  if err != nil || !exists {
    return []interface{}{}, nodesVisited, float64(time.Since(start).Milliseconds())
  }
  
  // Stack untuk DFS dengan tracking recipe path
  type StackItem struct {
    Name string         // Nama elemen
    Path []RecipeStep   // Path resep dari elemen ini
  }
  stack := []StackItem{{Name: elementName, Path: []RecipeStep{}}}
  
  // Hasil pencarian
  results := []interface{}{}
  
  // Cari elemen dasar (basic elements)
  basicElements := []string{}
  rows, err := db.Query("SELECT DISTINCT element FROM elements WHERE element NOT IN (SELECT DISTINCT element FROM elements WHERE item1 IS NOT NULL AND item2 IS NOT NULL)")
  if err == nil {
    defer rows.Close()
    for rows.Next() {
      var element string
      if err := rows.Scan(&element); err == nil {
        basicElements = append(basicElements, element)
      }
    }
  }
  
  // Process stack sampai kosong atau max recipes tercapai
  for len(stack) > 0 && (recipeType != "Limit" || len(results) < maxRecipes) {
    // Ambil item terakhir dari stack (LIFO - Last In First Out)
    current := stack[len(stack)-1]
    stack = stack[:len(stack)-1]
    
    // Skip jika sudah dikunjungi
    if visited[current.Name] {
      continue
    }
    
    // Tandai sudah dikunjungi dan tambah counter
    visited[current.Name] = true
    nodesVisited++
    
    // Cek apakah ini elemen dasar
    isBasic := false
    for _, basic := range basicElements {
      if current.Name == basic {
        isBasic = true
        break
      }
    }
    
    // Jika ini elemen dasar dan bukan target awal, tambahkan ke hasil dan lanjutkan
    if isBasic && current.Name != elementName {
      // Format node untuk response
      childrenData := []interface{}{}
      
      results = append(results, map[string]interface{}{
        "name":     current.Name,
        "image":    mapper[current.Name],
        "children": childrenData,
        "recipe":   "Basic element",
      })
      
      continue
    }
    
    // Jika ini target elemen dan kita memiliki recipe path, tambahkan ke hasil
    if current.Name == elementName && len(current.Path) > 0 {
      // Format recipe untuk response
      recipeTree := buildRecipeTree(current.Path)
      recipeSteps := formatRecipeSteps(current.Path)
      
      results = append(results, map[string]interface{}{
        "name":     current.Name,
        "image":    mapper[current.Name],
        "children": recipeTree,
        "recipe":   recipeSteps,
      })
      
      // Jika sudah mencapai max recipes, break
      if recipeType == "Limit" && len(results) >= maxRecipes {
        break
      }
    }
    
    // Dapatkan parent elements yang bisa dikombinasikan untuk membuat elemen ini
    rows, err := db.Query("SELECT item1, item2 FROM elements WHERE element = ? AND item1 IS NOT NULL AND item2 IS NOT NULL", current.Name)
    if err != nil {
      continue
    }
    
    // Process setiap kombinasi parent
    hasParents := false
    for rows.Next() {
      var item1, item2 string
      if err := rows.Scan(&item1, &item2); err != nil {
        continue
      }
      
      hasParents = true
      
      // Buat recipe step baru
      newStep := RecipeStep{
        Result: current.Name,
        Item1:  item1,
        Item2:  item2,
      }
      
      // Tambahkan parent2 ke stack terlebih dahulu (DFS)
      if !visited[item2] {
        newPath2 := make([]RecipeStep, len(current.Path))
        copy(newPath2, current.Path)
        newPath2 = append(newPath2, newStep)
        stack = append(stack, StackItem{Name: item2, Path: newPath2})
      }
      
      // Tambahkan parent1 ke stack (akan dikunjungi lebih dulu)
      if !visited[item1] {
        newPath1 := make([]RecipeStep, len(current.Path))
        copy(newPath1, current.Path)
        newPath1 = append(newPath1, newStep)
        stack = append(stack, StackItem{Name: item1, Path: newPath1})
      }
    }
    rows.Close()
    
    // Jika node ini tidak punya parents, tampilkan sebagai "no recipe"
    if !hasParents && current.Name != elementName {
      childData := []interface{}{}
      
      results = append(results, map[string]interface{}{
        "name":     current.Name,
        "image":    mapper[current.Name],
        "children": childData,
        "recipe":   "No parents found",
      })
    }
    
    // Ambil children untuk node ini (untuk tree display)
    children, err := getChildren(current.Name)
    if err == nil && len(children) > 0 {
      // Format children untuk response
      childData := make([]map[string]interface{}, len(children))
      for i, child := range children {
        childData[i] = map[string]interface{}{
          "name":  child.Name,
          "image": mapper[child.Name],
          "children": []interface{}{},
        }
      }
      
      // Tambahkan node dengan children ke hasil
      if current.Name == elementName && len(results) == 0 {
        results = append(results, map[string]interface{}{
          "name":     current.Name,
          "image":    mapper[current.Name],
          "children": childData,
        })
      }
    }
    
    time.Sleep(20 * time.Millisecond) // Delay kecil untuk visualisasi
  }
  
  // Jika tidak ada results dan ini adalah elemen target, tambahkan info elemen
  if len(results) == 0 {
    children, _ := getChildren(elementName)
    childData := make([]map[string]interface{}, 0)
    if children != nil {
      childData = make([]map[string]interface{}, len(children))
      for i, child := range children {
        childData[i] = map[string]interface{}{
          "name":  child.Name,
          "image": mapper[child.Name],
          "children": []interface{}{},
        }
      }
    }
    
    results = append(results, map[string]interface{}{
      "name":     elementName,
      "image":    mapper[elementName],
      "children": childData,
      "recipe":   "This is a basic element or no recipe found",
    })
  }
  
  return results, nodesVisited, float64(time.Since(start).Milliseconds())
}

// Algoritma Bidirectional untuk pencarian dari elemen target ke 4 elemen dasar
func Bidirectional(elementName string, recipeType string, maxRecipes int) ([]interface{}, int, float64) {
  start := time.Now() // Catat waktu mulai eksekusi
  
  // 4 elemen dasar
  basicElements := []string{"Earth", "Fire", "Water", "Air"}
  
  // Cek apakah elemen target ada di database
  var exists bool
  err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM elements WHERE element = ?)", elementName).Scan(&exists)
  if err != nil || !exists {
    return []interface{}{}, 0, float64(time.Since(start).Milliseconds())
  }
  
  // Struktur untuk menyimpan node dan path-nya
  type PathNode struct {
    Name       string
    Path       []RecipeStep
    Direction  string // "forward" atau "backward"
  }
  
  // Maps untuk tracking jalur dari kedua arah
  forwardPaths := make(map[string][]RecipeStep)  // target -> basic
  backwardPaths := make(map[string][]RecipeStep) // basic -> target
  
  // Maps untuk node yang sudah dikunjungi dari kedua arah
  forwardVisited := make(map[string]bool)
  backwardVisited := make(map[string]bool)
  
  // Queue untuk kedua arah
  queue := []PathNode{
    {Name: elementName, Path: []RecipeStep{}, Direction: "forward"},
  }
  
  // Tambahkan 4 elemen dasar ke queue (arah mundur)
  for _, basic := range basicElements {
    queue = append(queue, PathNode{
      Name:      basic,
      Path:      []RecipeStep{},
      Direction: "backward",
    })
  }
  
  // Hasil pencarian
  results := []interface{}{}
  
  // Counter node yang dikunjungi
  nodesVisited := 0
  
  // Intersection points antara jalur forward dan backward
  intersections := make(map[string]bool)
  
  // Process queue sampai kosong atau max recipes tercapai
  for len(queue) > 0 && (recipeType != "Limit" || len(results) < maxRecipes) {
    // Ambil node pertama dari queue
    current := queue[0]
    queue = queue[1:]
    
    // Skip jika sudah dikunjungi dari arah yang sama
    if (current.Direction == "forward" && forwardVisited[current.Name]) ||
       (current.Direction == "backward" && backwardVisited[current.Name]) {
      continue
    }
    
    // Tandai sudah dikunjungi dari arah tertentu
    if current.Direction == "forward" {
      forwardVisited[current.Name] = true
      forwardPaths[current.Name] = current.Path
    } else {
      backwardVisited[current.Name] = true
      backwardPaths[current.Name] = current.Path
    }
    
    // Tambah counter node yang dikunjungi
    nodesVisited++
    
    // Cek apakah node ini merupakan intersection point
    if (current.Direction == "forward" && backwardVisited[current.Name]) ||
       (current.Direction == "backward" && forwardVisited[current.Name]) {
      intersections[current.Name] = true
      
      // Jika ini adalah intersection dan belum mencapai max recipes, buat recipe lengkap
      if recipeType != "Limit" || len(results) < maxRecipes {
        // Gabungkan jalur dari kedua arah
        completePath := []RecipeStep{}
        
        if current.Direction == "forward" {
          // Gabungkan jalur backward (perlu dibalik) dan forward
          reversedBackPath := reverseRecipeSteps(backwardPaths[current.Name])
          completePath = append(completePath, reversedBackPath...)
          completePath = append(completePath, current.Path...)
        } else {
          // Gabungkan jalur forward dan backward (perlu dibalik)
          completePath = append(completePath, forwardPaths[current.Name]...)
          reversedCurrentPath := reverseRecipeSteps(current.Path)
          completePath = append(completePath, reversedCurrentPath...)
        }
        
        // Format recipe untuk response
        recipeTree := buildRecipeTree(completePath)
        recipeSteps := formatRecipeSteps(completePath)
        
        results = append(results, map[string]interface{}{
          "name":     elementName,
          "image":    mapper[elementName],
          "children": recipeTree,
          "recipe":   recipeSteps,
          "path":     current.Name, // Tambahkan info jalur yang digunakan
        })
      }
    }
    
    // Jika ini adalah arah forward (dari target ke elemen dasar)
    if current.Direction == "forward" {
      // Dapatkan parent elements dari elemen ini
      rows, err := db.Query("SELECT item1, item2 FROM elements WHERE element = ? AND item1 IS NOT NULL AND item2 IS NOT NULL", current.Name)
      if err != nil {
        continue
      }
      
      // Process setiap kombinasi parent
      for rows.Next() {
        var item1, item2 string
        if err := rows.Scan(&item1, &item2); err != nil {
          continue
        }
        
        // Buat recipe step baru
        newStep := RecipeStep{
          Result: current.Name,
          Item1:  item1,
          Item2:  item2,
        }
        
        // Tambahkan parent1 ke queue dengan path yang diupdate
        if !forwardVisited[item1] {
          newPath1 := make([]RecipeStep, len(current.Path))
          copy(newPath1, current.Path)
          newPath1 = append(newPath1, newStep)
          queue = append(queue, PathNode{Name: item1, Path: newPath1, Direction: "forward"})
        }
        
        // Tambahkan parent2 ke queue dengan path yang diupdate
        if !forwardVisited[item2] {
          newPath2 := make([]RecipeStep, len(current.Path))
          copy(newPath2, current.Path)
          newPath2 = append(newPath2, newStep)
          queue = append(queue, PathNode{Name: item2, Path: newPath2, Direction: "forward"})
        }
      }
      rows.Close()
    } else { // Jika ini adalah arah backward (dari elemen dasar ke target)
      // Dapatkan children elements dari elemen ini
      rows, err := db.Query("SELECT element FROM elements WHERE (item1 = ? OR item2 = ?) AND item1 IS NOT NULL AND item2 IS NOT NULL", current.Name, current.Name)
      if err != nil {
        continue
      }
      
      // Process setiap child
      for rows.Next() {
        var childName string
        if err := rows.Scan(&childName); err != nil {
          continue
        }
        
        // Untuk setiap child, ambil recipe lengkap (item1 + item2)
        var item1, item2 string
        err := db.QueryRow("SELECT item1, item2 FROM elements WHERE element = ? AND (item1 = ? OR item2 = ?)", 
                          childName, current.Name, current.Name).Scan(&item1, &item2)
        if err != nil {
          continue
        }
        
        // Buat recipe step baru
        newStep := RecipeStep{
          Result: childName,
          Item1:  item1,
          Item2:  item2,
        }
        
        // Tambahkan child ke queue dengan path yang diupdate
        if !backwardVisited[childName] {
          newPath := make([]RecipeStep, len(current.Path))
          copy(newPath, current.Path)
          newPath = append(newPath, newStep)
          queue = append(queue, PathNode{Name: childName, Path: newPath, Direction: "backward"})
        }
      }
      rows.Close()
    }
    
    time.Sleep(20 * time.Millisecond) // Delay kecil untuk visualisasi
  }
  
  // Jika tidak ada results dan ini adalah elemen target, tambahkan info elemen
  if len(results) == 0 {
    children, _ := getChildren(elementName)
    childData := make([]map[string]interface{}, 0)
    if children != nil {
      childData = make([]map[string]interface{}, len(children))
      for i, child := range children {
        childData[i] = map[string]interface{}{
          "name":  child.Name,
          "image": mapper[child.Name],
          "children": []interface{}{},
        }
      }
    }
    
    results = append(results, map[string]interface{}{
      "name":     elementName,
      "image":    mapper[elementName],
      "children": childData,
      "recipe":   "No path found between target and basic elements",
    })
  }
  
  return results, nodesVisited, float64(time.Since(start).Milliseconds())
}

// Helper function untuk membalik urutan recipe steps
func reverseRecipeSteps(steps []RecipeStep) []RecipeStep {
  if len(steps) == 0 {
    return steps
  }
  
  reversed := make([]RecipeStep, len(steps))
  for i, step := range steps {
    // Perlu membalik item1 dan item2 karena arah juga dibalik
    reversed[len(steps)-i-1] = RecipeStep{
      Result: step.Item1, // Ketika dibalik, result sekarang adalah item1
      Item1:  step.Result, // Dan result sebelumnya menjadi item1
      Item2:  step.Item2,  // item2 tetap sama
    }
  }
  
  return reversed
}