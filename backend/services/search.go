// File ini adalah service utama untuk pencarian resep di backend Little Alchemy 2.
// Berisi inisialisasi database, pembacaan mapper, dan implementasi algoritma BFS, DFS, dan Bidirectional.

package services

import (
	"database/sql"  // Untuk koneksi database SQLite
	"encoding/json" // Untuk decode file JSON
	"log"           // Untuk logging error
	// Untuk kalkulasi progress
	"os"   // Untuk akses file
	"time" // Untuk pengukuran waktu eksekusi

	_ "github.com/mattn/go-sqlite3" // Driver SQLite
	// Import utils untuk progress tracking
)

var db *sql.DB                    // Variabel global untuk koneksi database
var mapper map[string]string      // Variabel global untuk mapping nama elemen ke gambar/icon

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
  defer file.Close()
  if err := json.NewDecoder(file).Decode(&mapper); err != nil {
    log.Fatalf("Gagal mendekode mapper.json: %v", err)
  }
}

type Node struct {
  Name     string      // Nama elemen
  Children []*Node     // Anak-anak node (hasil kombinasi)
}

type RecipeStep struct {
  Result string // Element yang dihasilkan
  Item1  string // Bahan pertama
  Item2  string // Bahan kedua
}

// Mengambil anak-anak dari elemen tertentu dari database
func getChildren(name string) ([]*Node, error) {
  rows, err := db.Query("SELECT DISTINCT element FROM elements WHERE item1 = ? OR item2 = ?", name, name)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var children []*Node
  for rows.Next() {
    var childName string
    if err := rows.Scan(&childName); err != nil {
      return nil, err
    }
    children = append(children, &Node{Name: childName})
  }
  return children, nil
}

// BFS untuk pencarian resep (Murni BFS tanpa heuristik)
func BFS(elementName string, recipeType string, maxRecipes int) ([]interface{}, int, float64) {
  start := time.Now()
  visited := make(map[string]bool)
  nodesVisited := 0

  // Cek apakah elemen ada di database
  var exists bool
  err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM elements WHERE element = ?)", elementName).Scan(&exists)
  if err != nil || !exists {
    return getDefaultResult(elementName), nodesVisited, float64(time.Since(start).Milliseconds())
  }

  // Queue untuk BFS dengan tracking recipe path
  type QueueItem struct {
    Name string
    Path []RecipeStep
  }
  queue := []QueueItem{{Name: elementName, Path: []RecipeStep{}}}
  
  // Simpan semua path resep yang ditemukan
  var allPaths [][]RecipeStep
  basicElements := getBasicElements()

  // Process queue sampai kosong atau max recipes tercapai
  for len(queue) > 0 && (recipeType != "Limit" || len(allPaths) < maxRecipes) {
    current := queue[0]
    queue = queue[1:]

    // Tandai node sebagai dikunjungi dan naikkan counter
    if !visited[current.Name] {
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
  
      // Jika elemen dasar dan bukan target, simpan path
      if isBasic && current.Name != elementName && len(current.Path) > 0 {
        allPaths = append(allPaths, current.Path)
      }
  
      // Dapatkan parent elements - selalu eksplorasi semua path
      rows, err := db.Query("SELECT item1, item2 FROM elements WHERE element = ? AND item1 IS NOT NULL AND item2 IS NOT NULL", current.Name)
      if err != nil {
        continue
      }
  
      for rows.Next() {
        var item1, item2 string
        if err := rows.Scan(&item1, &item2); err != nil {
          continue
        }
  
        newStep := RecipeStep{
          Result: current.Name,
          Item1:  item1,
          Item2:  item2,
        }
  
        // Tambahkan item1 dan item2 ke queue tanpa memperhatikan kunjungan
        // Hanya cek pada saat diambil dari queue untuk menghindari pengulangan
        newPath1 := make([]RecipeStep, len(current.Path))
        copy(newPath1, current.Path)
        newPath1 = append(newPath1, newStep)
        queue = append(queue, QueueItem{Name: item1, Path: newPath1})
  
        newPath2 := make([]RecipeStep, len(current.Path))
        copy(newPath2, current.Path)
        newPath2 = append(newPath2, newStep)
        queue = append(queue, QueueItem{Name: item2, Path: newPath2})
      }
      rows.Close()
    }
  }

  // Bangun hasil dari semua path yang ditemukan
  results := buildResultsFromPaths(elementName, allPaths)
  log.Printf("Recipe results for %s: %+v", elementName, results)
  if len(results) == 0 {
    return getDefaultResult(elementName), nodesVisited, float64(time.Since(start).Milliseconds())
  }

  return results, nodesVisited, float64(time.Since(start).Milliseconds())
}

// DFS untuk pencarian resep (Murni DFS tanpa heuristik)
func DFS(elementName string, recipeType string, maxRecipes int) ([]interface{}, int, float64) {
  start := time.Now()
  visited := make(map[string]bool)
  nodesVisited := 0

  // Cek apakah elemen ada di database
  var exists bool
  err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM elements WHERE element = ?)", elementName).Scan(&exists)
  if err != nil || !exists {
    return getDefaultResult(elementName), nodesVisited, float64(time.Since(start).Milliseconds())
  }

  // Stack untuk DFS dengan tracking recipe path
  type StackItem struct {
    Name string
    Path []RecipeStep
  }
  stack := []StackItem{{Name: elementName, Path: []RecipeStep{}}}

  // Simpan semua path resep yang ditemukan
  var allPaths [][]RecipeStep
  basicElements := getBasicElements()

  // Process stack sampai kosong atau max recipes tercapai
  for len(stack) > 0 && (recipeType != "Limit" || len(allPaths) < maxRecipes) {
    current := stack[len(stack)-1]
    stack = stack[:len(stack)-1]

    // Proses node hanya jika belum dikunjungi
    if !visited[current.Name] {
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
  
      // Jika elemen dasar dan bukan target, simpan path
      if isBasic && current.Name != elementName && len(current.Path) > 0 {
        allPaths = append(allPaths, current.Path)
      }
  
      // Dapatkan parent elements - eksplorasi semua path
      rows, err := db.Query("SELECT item1, item2 FROM elements WHERE element = ? AND item1 IS NOT NULL AND item2 IS NOT NULL", current.Name)
      if err != nil {
        continue
      }
  
      // Menyimpan semua parent untuk diproses
      var parents []struct {
        Item1, Item2 string
      }
      
      for rows.Next() {
        var item1, item2 string
        if err := rows.Scan(&item1, &item2); err != nil {
          continue
        }
        parents = append(parents, struct {
          Item1, Item2 string
        }{Item1: item1, Item2: item2})
      }
      rows.Close()
      
      // Tambahkan parent ke stack dalam urutan terbalik untuk DFS murni
      for i := len(parents) - 1; i >= 0; i-- {
        newStep := RecipeStep{
          Result: current.Name,
          Item1:  parents[i].Item1,
          Item2:  parents[i].Item2,
        }
        
        // Tambahkan item1 dulu, maka item2 akan diproses terlebih dahulu (LIFO)
        newPath1 := make([]RecipeStep, len(current.Path))
        copy(newPath1, current.Path)
        newPath1 = append(newPath1, newStep)
        stack = append(stack, StackItem{Name: parents[i].Item1, Path: newPath1})
        
        // Kemudian tambahkan item2
        newPath2 := make([]RecipeStep, len(current.Path))
        copy(newPath2, current.Path)
        newPath2 = append(newPath2, newStep)
        stack = append(stack, StackItem{Name: parents[i].Item2, Path: newPath2})
      }
    }
  }

  // Bangun hasil dari semua path yang ditemukan
  results := buildResultsFromPaths(elementName, allPaths)
  log.Printf("Recipe results for %s: %+v", elementName, results)
  if len(results) == 0 {
    return getDefaultResult(elementName), nodesVisited, float64(time.Since(start).Milliseconds())
  }

  return results, nodesVisited, float64(time.Since(start).Milliseconds())
}

// Bidirectional untuk pencarian resep (Murni Bidirectional tanpa heuristik)
func Bidirectional(elementName string, recipeType string, maxRecipes int) ([]interface{}, int, float64) {
  start := time.Now()
  basicElements := getBasicElements()

  // Cek apakah elemen ada di database
  var exists bool
  err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM elements WHERE element = ?)", elementName).Scan(&exists)
  if err != nil || !exists {
    return getDefaultResult(elementName), 0, float64(time.Since(start).Milliseconds())
  }

  type PathNode struct {
    Name      string
    Path      []RecipeStep
    Direction string
  }

  // Maps untuk menyimpan path dan status kunjungan
  forwardPaths := make(map[string][]RecipeStep)
  backwardPaths := make(map[string][]RecipeStep)
  forwardVisited := make(map[string]bool)
  backwardVisited := make(map[string]bool)

  // Queue untuk algoritma BFS dua arah
  forwardQueue := []PathNode{{Name: elementName, Path: []RecipeStep{}, Direction: "forward"}}
  backwardQueue := []PathNode{}
  
  // Inisialisasi queue mundur dengan semua elemen dasar
  for _, basic := range basicElements {
    backwardQueue = append(backwardQueue, PathNode{Name: basic, Path: []RecipeStep{}, Direction: "backward"})
  }

  var allPaths [][]RecipeStep
  nodesVisited := 0
  
  // Proses hingga salah satu queue kosong atau max recipes tercapai
  for len(forwardQueue) > 0 && len(backwardQueue) > 0 && (recipeType != "Limit" || len(allPaths) < maxRecipes) {
    // Proses satu langkah dari forward search
    if len(forwardQueue) > 0 {
      current := forwardQueue[0]
      forwardQueue = forwardQueue[1:]
      
      if !forwardVisited[current.Name] {
        forwardVisited[current.Name] = true
        forwardPaths[current.Name] = current.Path
        nodesVisited++
        
        // Cek interseksi dengan backward search
        if backwardVisited[current.Name] {
          completePath := []RecipeStep{}
          reversedBackPath := reverseRecipeSteps(backwardPaths[current.Name])
          completePath = append(completePath, current.Path...)
          completePath = append(completePath, reversedBackPath...)
          allPaths = append(allPaths, completePath)
        }
        
        // Eksplorasi parent elements (forward search)
        rows, err := db.Query("SELECT item1, item2 FROM elements WHERE element = ? AND item1 IS NOT NULL AND item2 IS NOT NULL", current.Name)
        if err == nil {
          for rows.Next() {
            var item1, item2 string
            if err := rows.Scan(&item1, &item2); err != nil {
              continue
            }
            
            newStep := RecipeStep{
              Result: current.Name,
              Item1:  item1,
              Item2:  item2,
            }
            
            // Tambahkan kedua parent ke queue tanpa cek visited
            newPath1 := make([]RecipeStep, len(current.Path))
            copy(newPath1, current.Path)
            newPath1 = append(newPath1, newStep)
            forwardQueue = append(forwardQueue, PathNode{Name: item1, Path: newPath1, Direction: "forward"})
            
            newPath2 := make([]RecipeStep, len(current.Path))
            copy(newPath2, current.Path)
            newPath2 = append(newPath2, newStep)
            forwardQueue = append(forwardQueue, PathNode{Name: item2, Path: newPath2, Direction: "forward"})
          }
          rows.Close()
        }
      }
    }
    
    // Proses satu langkah dari backward search
    if len(backwardQueue) > 0 {
      current := backwardQueue[0]
      backwardQueue = backwardQueue[1:]
      
      if !backwardVisited[current.Name] {
        backwardVisited[current.Name] = true
        backwardPaths[current.Name] = current.Path
        nodesVisited++
        
        // Cek interseksi dengan forward search
        if forwardVisited[current.Name] {
          completePath := []RecipeStep{}
          completePath = append(completePath, forwardPaths[current.Name]...)
          reversedCurrentPath := reverseRecipeSteps(current.Path)
          completePath = append(completePath, reversedCurrentPath...)
          allPaths = append(allPaths, completePath)
        }
        
        // Eksplorasi kombinasi yang dapat menghasilkan element ini (backward search)
        rows, err := db.Query("SELECT element, item1, item2 FROM elements WHERE (item1 = ? OR item2 = ?) AND item1 IS NOT NULL AND item2 IS NOT NULL", current.Name, current.Name)
        if err == nil {
          for rows.Next() {
            var resultElement, item1, item2 string
            if err := rows.Scan(&resultElement, &item1, &item2); err != nil {
              continue
            }
            
            newStep := RecipeStep{
              Result: resultElement,
              Item1:  item1,
              Item2:  item2,
            }
            
            // Tambahkan hasil kombinasi ke queue mundur
            if !backwardVisited[resultElement] {
              newPath := make([]RecipeStep, len(current.Path))
              copy(newPath, current.Path)
              newPath = append(newPath, newStep)
              backwardQueue = append(backwardQueue, PathNode{Name: resultElement, Path: newPath, Direction: "backward"})
            }
          }
          rows.Close()
        }
      }
    }
  }

  // Bangun hasil dari semua path yang ditemukan
  results := buildResultsFromPaths(elementName, allPaths)
  log.Printf("Recipe results for %s: %+v", elementName, results)
  if len(results) == 0 {
    return getDefaultResult(elementName), nodesVisited, float64(time.Since(start).Milliseconds())
  }

  return results, nodesVisited, float64(time.Since(start).Milliseconds())
}

// Helper function untuk mendapatkan elemen dasar
func getBasicElements() []string {
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
  return basicElements
}

// Fix for buildResultsFromPaths and buildRecipeTree functions

func buildResultsFromPaths(elementName string, paths [][]RecipeStep) []interface{} {
  // Jika tidak ada path yang ditemukan, kembalikan hasil kosong
  if len(paths) == 0 {
    return []interface{}{}
  }
  
  // Menggabungkan semua langkah resep unik dari semua path
  allUniqueSteps := make(map[string]RecipeStep)
  for _, path := range paths {
    for _, step := range path {
      // Gunakan kombinasi Result+Item1+Item2 sebagai kunci unik
      key := step.Result + "|" + step.Item1 + "|" + step.Item2
      allUniqueSteps[key] = step
    }
  }
  
  // Konversi map ke slice
  var combinedSteps []RecipeStep
  for _, step := range allUniqueSteps {
    combinedSteps = append(combinedSteps, step)
  }
  
  // Format steps untuk tampilan
  var recipeSteps []string
  if len(paths) > 0 && len(paths[0]) > 0 {
    recipeSteps = formatRecipeSteps(paths[0])
  }
  
  // Bangun tree terpadu dengan semua langkah yang dikumpulkan
  unifiedTree := buildUnifiedRecipeTree(elementName, combinedSteps)
  
  return []interface{}{
    map[string]interface{}{
      "name":     elementName,
      "image":    mapper[elementName],
      "children": unifiedTree,
      "recipe":   recipeSteps,
    },
  }
}

// Add this new function to build a unified tree
func buildUnifiedRecipeTree(targetElement string, allSteps []RecipeStep) []map[string]interface{} {
  // Buat map untuk melacak elemen yang sudah diproses untuk menghindari loop tak terbatas
  processedElements := make(map[string]bool)
  
  // Panggil fungsi rekursif dengan tracking elemen yang sudah diproses
  return buildRecipeTreeRecursive(targetElement, allSteps, processedElements)
}

// Fungsi rekursif baru yang membangun tree dengan kedalaman penuh
func buildRecipeTreeRecursive(elementName string, allSteps []RecipeStep, processedElements map[string]bool) []map[string]interface{} {
  // Cek apakah elemen ini sudah diproses sebelumnya (untuk menghindari loop tak terbatas)
  if processedElements[elementName] {
    return []map[string]interface{}{}
  }
  
  // Tandai elemen ini sebagai sudah diproses
  processedElements[elementName] = true
  
  // Cari recipe step untuk elemen ini
  var recipeForElement *RecipeStep
  for _, step := range allSteps {
    if step.Result == elementName {
      recipeForElement = &step
      break
    }
  }
  
  // Jika tidak ada recipe untuk elemen ini, kembalikan array kosong (elemen dasar)
  if recipeForElement == nil {
    processedElements[elementName] = false // Reset flag
    return []map[string]interface{}{}
  }
  
  // Buat node untuk kedua komponen dari elemen ini
  item1Node := map[string]interface{}{
    "name":  recipeForElement.Item1,
    "image": mapper[recipeForElement.Item1],
  }
  
  item2Node := map[string]interface{}{
    "name":  recipeForElement.Item2,
    "image": mapper[recipeForElement.Item2],
  }
  
  // Bangun children untuk item1 secara rekursif
  // Buat salinan dari processedElements untuk menghindari konflik
  item1Processed := make(map[string]bool)
  for k, v := range processedElements {
    item1Processed[k] = v
  }
  item1Node["children"] = buildRecipeTreeRecursive(recipeForElement.Item1, allSteps, item1Processed)
  
  // Bangun children untuk item2 secara rekursif
  // Buat salinan dari processedElements untuk menghindari konflik
  item2Processed := make(map[string]bool)
  for k, v := range processedElements {
    item2Processed[k] = v
  }
  item2Node["children"] = buildRecipeTreeRecursive(recipeForElement.Item2, allSteps, item2Processed)
  
  // Reset flag untuk elemen ini setelah selesai memproses
  processedElements[elementName] = false
  
  // Kembalikan array yang berisi kedua node sebagai children
  return []map[string]interface{}{item1Node, item2Node}
}

// Helper function untuk format recipe steps
func formatRecipeSteps(path []RecipeStep) []string {
  if len(path) == 0 {
    return []string{}
  }

  steps := make([]string, len(path))
  for i, step := range path {
    steps[len(path)-i-1] = step.Result + " = " + step.Item1 + " + " + step.Item2
  }
  return steps
}

// Helper function untuk membalik urutan recipe steps
func reverseRecipeSteps(steps []RecipeStep) []RecipeStep {
  if len(steps) == 0 {
    return steps
  }

  reversed := make([]RecipeStep, len(steps))
  for i, step := range steps {
    reversed[len(steps)-i-1] = RecipeStep{
      Result: step.Item1,
      Item1:  step.Result,
      Item2:  step.Item2,
    }
  }
  return reversed
}

// Helper function untuk hasil default jika tidak ada resep
func getDefaultResult(elementName string) []interface{} {
  children, _ := getChildren(elementName)
  childData := make([]map[string]interface{}, 0)
  if children != nil {
    childData = make([]map[string]interface{}, len(children))
    for i, child := range children {
      childData[i] = map[string]interface{}{
        "name":     child.Name,
        "image":    mapper[child.Name],
        "children": []interface{}{},
      }
    }
  }

  return []interface{}{
    map[string]interface{}{
      "name":     elementName,
      "image":    mapper[elementName],
      "children": childData,
      "recipe":   "This is a basic element or no recipe found",
    },
  }
}