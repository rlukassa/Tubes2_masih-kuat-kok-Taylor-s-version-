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

// Semua element di mana item1 = "Water" atau item2 = "Water"
// Dari tabel:

// Baris 1: Steam (item1 = Water)
// Baris 2: Mud (item1 = Water)
// Hasil query:
// ["Steam", "Mud"]
// Artinya, "Water" bisa digunakan untuk membuat "Steam" dan "Mud".
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

// Mengambil parent dari elemen tertentu dari database
func getParents(name string) ([]*Node, error) { 
  rows, err := db.Query("SELECT DISTINCT item1 FROM elements WHERE element = ? AND item1 IS NOT NULL UNION SELECT DISTINCT item2 FROM elements WHERE element = ? AND item2 IS NOT NULL", name, name)
//   Query ini digunakan untuk mencari parent dari sebuah elemen.
// elements adalah tabel yang berisi kolom element, item1, dan item2.
// Query ini mencari semua nilai unik (DISTINCT) dari item1 dan item2 yang menghasilkan element tertentu (name).
// item1 IS NOT NULL dan item2 IS NOT NULL memastikan hanya mengambil parent yang valid (tidak kosong).
// UNION digunakan agar hasil dari kedua pencarian digabungkan tanpa duplikat

 // misalkan :
// element	item1	item2
// Steam	Water	Fire
// Mud	Water	Earth
// Dust	Earth	Air
// Lava	Earth	Fire

// Jika name = "Mud", maka query akan mencari:
// Semua item1 dan item2 di mana element = "Mud"
// → Hasil: item1 = "Water", item2 = "Earth"
// Hasil query:
// ["Water", "Earth"]
// Artinya, parent dari "Mud" adalah "Water" dan "Earth".
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var parents []*Node
  for rows.Next() {
    var parentName string
    if err := rows.Scan(&parentName); err != nil {
      return nil, err
    }
    parents = append(parents, &Node{Name: parentName})
  }
  return parents, nil
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
//KODE INI DIUBAH (BAWAH)

//KODE INI SAYA HAPUS (ATAS) ~ini kode sebelumnya kas :v
// Algoritma BFS untuk pencarian resep
// func BFS(elementName string, recipeType string, maxRecipes int) ([]interface{}, int, float64) {
//   start := time.Now() // Catat waktu mulai eksekusi
//   visited := make(map[string]bool) // Map untuk menandai node yang sudah dikunjungi (agar tidak loop)
//   queue := []*Node{{Name: elementName}} // Queue BFS, mulai dari elemen awal
//   results := []interface{}{} // Hasil pencarian (akan diisi node-node hasil BFS)
//   nodesVisited := 0 // Counter node yang dikunjungi

//   for len(queue) > 0 && (recipeType != "Limit" || len(results) < maxRecipes) {
//     node := queue[0] // Ambil node pertama dari queue (FIFO)
//     queue = queue[1:] // Hapus node dari queue
//     if visited[node.Name] {
//       continue // Jika sudah dikunjungi, skip
//     }
//     visited[node.Name] = true // Tandai sudah dikunjungi
//     nodesVisited++ // Tambah counter

//     children, err := getChildren(node.Name) // Ambil anak-anak node dari database
//     if err != nil {
//       continue // Jika gagal, skip node ini
//     }
//     node.Children = children // Set anak-anak node

//     // Siapkan data anak-anak untuk response (bentuk array of map)
//     childData := make([]map[string]interface{}, len(children))
//     for i, child := range children {
//       childData[i] = map[string]interface{}{
//         "name":  child.Name,           // Nama anak
//         "image": mapper[child.Name],   // URL gambar anak
//         "children": []interface{}{},   // Anak dari anak (kosong, hanya 1 level di response)
//       }
//     }

//     // Tambahkan node ke hasil pencarian
//     results = append(results, map[string]interface{}{
//       "name":     node.Name,           // Nama node
//       "image":    mapper[node.Name],   // URL gambar node
//       "children": childData,           // Anak-anak node
//     })

//     queue = append(queue, children...) // Tambahkan anak-anak ke queue BFS
//     time.Sleep(100 * time.Millisecond) // Simulasi delay (misal untuk animasi/live update)
//   }

//   // Contoh:
//   // Jika elementName = "Water", dan di DB:
//   // element   | item1  | item2
//   // --------------------------
//   // Steam     | Water  | Fire
//   // Mud       | Water  | Earth
//   // Maka hasil BFS level 1: ["Steam", "Mud"]
//   // results = [
//   //   {name: "Water", children: [{name: "Steam"}, {name: "Mud"}]}
//   // ]

//   return results, nodesVisited, float64(time.Since(start).Milliseconds()) // Kembalikan hasil, jumlah node, dan waktu eksekusi
// }
// KODE INI SAYA HAPUS (BAWAH)

// Algoritma DFS untuk pencarian resep
func DFS(elementName string, recipeType string, maxRecipes int) ([]interface{}, int, float64) {
    start := time.Now() // Catat waktu mulai eksekusi
    visited := make(map[string]bool) // Map untuk menandai node yang sudah dikunjungi
    stack := []*Node{{Name: elementName}} // Stack DFS, mulai dari elemen awal
    results := []interface{}{} // Hasil pencarian
    nodesVisited := 0 // Counter node yang dikunjungi

    for len(stack) > 0 && (recipeType != "Limit" || len(results) < maxRecipes) {
        node := stack[len(stack)-1] // Ambil node terakhir dari stack (LIFO)
        stack = stack[:len(stack)-1] // Hapus node dari stack
        if visited[node.Name] {
            continue // Jika sudah dikunjungi, skip
        }
        visited[node.Name] = true // Tandai sudah dikunjungi
        nodesVisited++ // Tambah counter

        children, err := getChildren(node.Name) // Ambil anak-anak node dari database
        if err != nil {
            log.Printf("Gagal mengambil anak untuk elemen %s: %v", node.Name, err)
            continue
        }
        node.Children = children // Set anak-anak node

        // Tambahkan node ke hasil pencarian (di DFS, children bisa berupa node struct)
        results = append(results, map[string]interface{}{
            "name":     node.Name,
            "image":    mapper[node.Name],
            "children": children,
        })

        stack = append(stack, children...) // Tambahkan anak-anak ke stack (DFS)
        time.Sleep(100 * time.Millisecond) // Simulasi delay
    }

    // Contoh:
    // Jika elementName = "Water", hasil DFS bisa ["Water" -> "Steam" -> ...] tergantung urutan stack

    return results, nodesVisited, float64(time.Since(start).Milliseconds())
}

// Algoritma Bidirectional untuk pencarian dua arah
func Bidirectional(elementName string, targetName string, recipeType string, maxRecipes int) ([]interface{}, int, float64) {
    start := time.Now() // Catat waktu mulai eksekusi
    forwardVisited := make(map[string]bool) // Map node yang sudah dikunjungi dari depan
    backwardVisited := make(map[string]bool) // Map node yang sudah dikunjungi dari belakang
    forwardQueue := []*Node{{Name: elementName}} // Queue dari elemen awal
    backwardQueue := []*Node{{Name: targetName}} // Queue dari target
    results := []interface{}{} // Hasil pencarian
    nodesVisited := 0 // Counter node yang dikunjungi

    for (len(forwardQueue) > 0 || len(backwardQueue) > 0) && (recipeType != "Limit" || len(results) < maxRecipes) {
        // Langkah maju (dari elemen awal)
        if len(forwardQueue) > 0 {
            node := forwardQueue[0] // Ambil node pertama dari queue depan
            forwardQueue = forwardQueue[1:] // Hapus node dari queue
            if forwardVisited[node.Name] {
                continue // Jika sudah dikunjungi, skip
            }
            forwardVisited[node.Name] = true // Tandai sudah dikunjungi
            nodesVisited++ // Tambah counter

            children, err := getChildren(node.Name) // Ambil anak-anak node
            if err != nil {
                log.Printf("Gagal mengambil anak untuk elemen %s: %v", node.Name, err)
                continue
            }
            node.Children = children
            results = append(results, map[string]interface{}{
                "name":     node.Name,
                "image":    mapper[node.Name],
                "children": children,
            })
            forwardQueue = append(forwardQueue, children...) // Tambahkan anak-anak ke queue depan
        }

        // Langkah mundur (dari target)
        if len(backwardQueue) > 0 {
            node := backwardQueue[0] // Ambil node pertama dari queue belakang
            backwardQueue = backwardQueue[1:] // Hapus node dari queue
            if backwardVisited[node.Name] {
                continue // Jika sudah dikunjungi, skip
            }
            backwardVisited[node.Name] = true // Tandai sudah dikunjungi
            nodesVisited++ // Tambah counter

            parents, err := getParents(node.Name) // Ambil parent dari node (kebalikan getChildren)
            if err != nil {
                log.Printf("Gagal mengambil induk untuk elemen %s: %v", node.Name, err)
                continue
            }
            backwardQueue = append(backwardQueue, parents...) // Tambahkan parent ke queue belakang
        }

        time.Sleep(100 * time.Millisecond) // Simulasi delay
    }

    // Contoh:
    // Jika elementName = "Water", targetName = "Mud"
    // Algoritma akan mencari dari dua arah: dari "Water" ke depan, dari "Mud" ke belakang (parent)
    // Jika ditemukan node yang sama di kedua sisi, berarti ada jalur penghubung.

    return results, nodesVisited, float64(time.Since(start).Milliseconds())
}