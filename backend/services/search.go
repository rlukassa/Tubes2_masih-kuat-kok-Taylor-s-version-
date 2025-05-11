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

// BFS untuk pencarian resep
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

    if visited[current.Name] {
      continue
    }

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
      continue
    }

    // Dapatkan parent elements
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

      if !visited[item1] {
        newPath := make([]RecipeStep, len(current.Path))
        copy(newPath, current.Path)
        newPath = append(newPath, newStep)
        queue = append(queue, QueueItem{Name: item1, Path: newPath})
      }

      if !visited[item2] {
        newPath := make([]RecipeStep, len(current.Path))
        copy(newPath, current.Path)
        newPath = append(newPath, newStep)
        queue = append(queue, QueueItem{Name: item2, Path: newPath})
      }
    }
    rows.Close()
  }

  // Bangun hasil dari semua path yang ditemukan
  results := buildResultsFromPaths(elementName, allPaths)
  if len(results) == 0 {
    return getDefaultResult(elementName), nodesVisited, float64(time.Since(start).Milliseconds())
  }

  return results, nodesVisited, float64(time.Since(start).Milliseconds())
}

// DFS untuk pencarian resep
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

    if visited[current.Name] {
      continue
    }

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
      continue
    }

    // Dapatkan parent elements
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

      if !visited[item2] {
        newPath := make([]RecipeStep, len(current.Path))
        copy(newPath, current.Path)
        newPath = append(newPath, newStep)
        stack = append(stack, StackItem{Name: item2, Path: newPath})
      }

      if !visited[item1] {
        newPath := make([]RecipeStep, len(current.Path))
        copy(newPath, current.Path)
        newPath = append(newPath, newStep)
        stack = append(stack, StackItem{Name: item1, Path: newPath})
      }
    }
    rows.Close()
  }

  // Bangun hasil dari semua path yang ditemukan
  results := buildResultsFromPaths(elementName, allPaths)
  if len(results) == 0 {
    return getDefaultResult(elementName), nodesVisited, float64(time.Since(start).Milliseconds())
  }

  return results, nodesVisited, float64(time.Since(start).Milliseconds())
}

// Bidirectional untuk pencarian resep
func Bidirectional(elementName string, recipeType string, maxRecipes int) ([]interface{}, int, float64) {
  start := time.Now()
  basicElements := []string{"Earth", "Fire", "Water", "Air"}

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

  forwardPaths := make(map[string][]RecipeStep)
  backwardPaths := make(map[string][]RecipeStep)
  forwardVisited := make(map[string]bool)
  backwardVisited := make(map[string]bool)

  queue := []PathNode{{Name: elementName, Path: []RecipeStep{}, Direction: "forward"}}
  for _, basic := range basicElements {
    queue = append(queue, PathNode{Name: basic, Path: []RecipeStep{}, Direction: "backward"})
  }

  var allPaths [][]RecipeStep
  nodesVisited := 0
  intersections := make(map[string]bool)

  for len(queue) > 0 && (recipeType != "Limit" || len(allPaths) < maxRecipes) {
    current := queue[0]
    queue = queue[1:]

    if (current.Direction == "forward" && forwardVisited[current.Name]) ||
       (current.Direction == "backward" && backwardVisited[current.Name]) {
      continue
    }

    if current.Direction == "forward" {
      forwardVisited[current.Name] = true
      forwardPaths[current.Name] = current.Path
    } else {
      backwardVisited[current.Name] = true
      backwardPaths[current.Name] = current.Path
    }

    nodesVisited++

    if (current.Direction == "forward" && backwardVisited[current.Name]) ||
       (current.Direction == "backward" && forwardVisited[current.Name]) {
      intersections[current.Name] = true

      completePath := []RecipeStep{}
      if current.Direction == "forward" {
        reversedBackPath := reverseRecipeSteps(backwardPaths[current.Name])
        completePath = append(completePath, reversedBackPath...)
        completePath = append(completePath, current.Path...)
      } else {
        completePath = append(completePath, forwardPaths[current.Name]...)
        reversedCurrentPath := reverseRecipeSteps(current.Path)
        completePath = append(completePath, reversedCurrentPath...)
      }
      allPaths = append(allPaths, completePath)
    }

    if current.Direction == "forward" {
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

        if !forwardVisited[item1] {
          newPath := make([]RecipeStep, len(current.Path))
          copy(newPath, current.Path)
          newPath = append(newPath, newStep)
          queue = append(queue, PathNode{Name: item1, Path: newPath, Direction: "forward"})
        }

        if !forwardVisited[item2] {
          newPath := make([]RecipeStep, len(current.Path))
          copy(newPath, current.Path)
          newPath = append(newPath, newStep)
          queue = append(queue, PathNode{Name: item2, Path: newPath, Direction: "forward"})
        }
      }
      rows.Close()
    } else {
      rows, err := db.Query("SELECT element FROM elements WHERE (item1 = ? OR item2 = ?) AND item1 IS NOT NULL AND item2 IS NOT NULL", current.Name, current.Name)
      if err != nil {
        continue
      }

      for rows.Next() {
        var childName string
        if err := rows.Scan(&childName); err != nil {
          continue
        }

        var item1, item2 string
        err := db.QueryRow("SELECT item1, item2 FROM elements WHERE element = ? AND (item1 = ? OR item2 = ?)",
                          childName, current.Name, current.Name).Scan(&item1, &item2)
        if err != nil {
          continue
        }

        newStep := RecipeStep{
          Result: childName,
          Item1:  item1,
          Item2:  item2,
        }

        if !backwardVisited[childName] {
          newPath := make([]RecipeStep, len(current.Path))
          copy(newPath, current.Path)
          newPath = append(newPath, newStep)
          queue = append(queue, PathNode{Name: childName, Path: newPath, Direction: "backward"})
        }
      }
      rows.Close()
    }
  }

  // Bangun hasil dari semua path yang ditemukan
  results := buildResultsFromPaths(elementName, allPaths)
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

// Helper function untuk membangun hasil dari semua path
func buildResultsFromPaths(elementName string, paths [][]RecipeStep) []interface{} {
  results := []interface{}{}
  for _, path := range paths {
    if len(path) == 0 {
      continue
    }
    recipeTree := buildRecipeTree(path)
    recipeSteps := formatRecipeSteps(path)
    results = append(results, map[string]interface{}{
      "name":     elementName,
      "image":    mapper[elementName],
      "children": recipeTree,
      "recipe":   recipeSteps,
    })
  }
  return results
}

// Helper function untuk membangun pohon resep
func buildRecipeTree(path []RecipeStep) []map[string]interface{} {
  if len(path) == 0 {
    return []map[string]interface{}{}
  }

  lastStep := path[len(path)-1]
  var item1Steps, item2Steps []RecipeStep
  for _, step := range path {
    if step.Result == lastStep.Item1 {
      item1Steps = append(item1Steps, step)
    }
    if step.Result == lastStep.Item2 {
      item2Steps = append(item2Steps, step)
    }
  }

  var item1Path, item2Path []RecipeStep
  for _, step := range path[:len(path)-1] {
    if step.Result == lastStep.Item1 || containsElement(item1Path, step.Result) {
      item1Path = append(item1Path, step)
    }
    if step.Result == lastStep.Item2 || containsElement(item2Path, step.Result) {
      item2Path = append(item2Path, step)
    }
  }

  result := []map[string]interface{}{}
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