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

// Algoritma BFS untuk pencarian resep
func BFS(elementName string, recipeType string, maxRecipes int) ([]interface{}, int, float64) {
  start := time.Now() // Catat waktu mulai eksekusi
  visited := make(map[string]bool) // Map untuk menandai node yang sudah dikunjungi (agar tidak loop)
  queue := []*Node{{Name: elementName}} // Queue BFS, mulai dari elemen awal
  results := []interface{}{} // Hasil pencarian (akan diisi node-node hasil BFS)
  nodesVisited := 0 // Counter node yang dikunjungi

  for len(queue) > 0 && (recipeType != "Limit" || len(results) < maxRecipes) {
    node := queue[0] // Ambil node pertama dari queue (FIFO)
    queue = queue[1:] // Hapus node dari queue
    if visited[node.Name] {
      continue // Jika sudah dikunjungi, skip
    }
    visited[node.Name] = true // Tandai sudah dikunjungi
    nodesVisited++ // Tambah counter

    children, err := getChildren(node.Name) // Ambil anak-anak node dari database
    if err != nil {
      continue // Jika gagal, skip node ini
    }
    node.Children = children // Set anak-anak node

    // Siapkan data anak-anak untuk response (bentuk array of map)
    childData := make([]map[string]interface{}, len(children))
    for i, child := range children {
      childData[i] = map[string]interface{}{
        "name":  child.Name,           // Nama anak
        "image": mapper[child.Name],   // URL gambar anak
        "children": []interface{}{},   // Anak dari anak (kosong, hanya 1 level di response)
      }
    }

    // Tambahkan node ke hasil pencarian
    results = append(results, map[string]interface{}{
      "name":     node.Name,           // Nama node
      "image":    mapper[node.Name],   // URL gambar node
      "children": childData,           // Anak-anak node
    })

    queue = append(queue, children...) // Tambahkan anak-anak ke queue BFS
    time.Sleep(100 * time.Millisecond) // Simulasi delay (misal untuk animasi/live update)
  }

  // Contoh:
  // Jika elementName = "Water", dan di DB:
  // element   | item1  | item2
  // --------------------------
  // Steam     | Water  | Fire
  // Mud       | Water  | Earth
  // Maka hasil BFS level 1: ["Steam", "Mud"]
  // results = [
  //   {name: "Water", children: [{name: "Steam"}, {name: "Mud"}]}
  // ]

  return results, nodesVisited, float64(time.Since(start).Milliseconds()) // Kembalikan hasil, jumlah node, dan waktu eksekusi
}

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