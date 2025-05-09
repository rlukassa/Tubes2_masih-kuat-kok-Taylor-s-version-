package services

import (
    "database/sql"
    "encoding/json"
    "log"
    "os"
    "time"

    _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var mapper map[string]string

func init() {
    var err error
    // Path database yang benar
    db, err = sql.Open("sqlite3", "c:/Users/USER/OneDrive/Desktop/tubes2stima/Tubes2Stima/database/alchemy.db")
    if err != nil {
        log.Fatalf("Gagal membuka database: %v", err)
    }
    if err = db.Ping(); err != nil {
        log.Fatalf("Gagal terhubung ke database: %v", err)
    }
    log.Println("Berhasil terhubung ke database")

    // Load mapper.json
    file, err := os.Open("C:/Users/USER/OneDrive/Desktop/tubes2stima/Tubes2Stima/database/mapper2.json")
    if err != nil {
        log.Fatalf("Gagal membuka mapper.json: %v", err)
    }
    defer file.Close()
    if err := json.NewDecoder(file).Decode(&mapper); err != nil {
        log.Fatalf("Gagal mendekode mapper.json: %v", err)
    }
    log.Println("Berhasil memuat mapper.json")
}

// FetchElements mengambil semua elemen unik dari database
func FetchElements() ([]map[string]interface{}, error) {
    // Periksa apakah tabel 'elements' ada
    var tableExists bool
    err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='elements'").Scan(&tableExists)
    if err != nil {
        log.Printf("Gagal memeriksa keberadaan tabel: %v", err)
        return nil, err
    }
    if !tableExists {
        log.Println("Tabel 'elements' tidak ada di database")
        return []map[string]interface{}{}, nil
    }

    rows, err := db.Query("SELECT DISTINCT element FROM elements")
    if err != nil {
        log.Printf("Gagal menjalankan query elemen: %v", err)
        return nil, err
    }
    defer rows.Close()

    var elements []map[string]interface{}
    for rows.Next() {
        var element string
        if err := rows.Scan(&element); err != nil {
            log.Printf("Gagal memindai baris: %v", err)
            return nil, err
        }
        elements = append(elements, map[string]interface{}{
            "name":  element,
            "image": mapper[element],
        })
    }
    log.Printf("Berhasil mengambil %d elemen", len(elements))
    return elements, nil
}

// Node merepresentasikan node elemen dalam pohon pencarian
type Node struct {
    Name     string
    Children []*Node
}

// getChildren mengambil elemen yang menggunakan name sebagai item1 atau item2
func getChildren(name string) ([]*Node, error) {
    rows, err := db.Query("SELECT DISTINCT element FROM elements WHERE item1 = ? OR item2 = ?", name, name)
    if err != nil {
        log.Printf("Gagal mengambil anak: %v", err)
        return nil, err
    }
    defer rows.Close()

    var children []*Node
    for rows.Next() {
        var childName string
        if err := rows.Scan(&childName); err != nil {
            log.Printf("Gagal memindai anak: %v", err)
            return nil, err
        }
        children = append(children, &Node{Name: childName})
    }
    return children, nil
}

// getParents mengambil elemen yang merupakan induk dari name tertentu
func getParents(name string) ([]*Node, error) {
    rows, err := db.Query("SELECT DISTINCT item1 FROM elements WHERE element = ? AND item1 IS NOT NULL UNION SELECT DISTINCT item2 FROM elements WHERE element = ? AND item2 IS NOT NULL", name, name)
    if err != nil {
        log.Printf("Gagal mengambil induk: %v", err)
        return nil, err
    }
    defer rows.Close()

    var parents []*Node
    for rows.Next() {
        var parentName string
        if err := rows.Scan(&parentName); err != nil {
            log.Printf("Gagal memindai induk: %v", err)
            return nil, err
        }
        parents = append(parents, &Node{Name: parentName})
    }
    return parents, nil
}

// BFS melakukan pencarian lebar pertama untuk resep
func BFS(elementName string, recipeType string, maxRecipes int) ([]interface{}, int, float64) {
    start := time.Now()
    visited := make(map[string]bool)
    queue := []*Node{{Name: elementName}}
    results := []interface{}{}
    nodesVisited := 0

    for len(queue) > 0 && (recipeType != "Limit" || len(results) < maxRecipes) {
        node := queue[0]
        queue = queue[1:]
        if visited[node.Name] {
            continue
        }
        visited[node.Name] = true
        nodesVisited++

        children, err := getChildren(node.Name)
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

        queue = append(queue, children...)
        time.Sleep(100 * time.Millisecond) // Simulasi penundaan
    }

    return results, nodesVisited, float64(time.Since(start).Milliseconds())
}

// DFS melakukan pencarian kedalaman pertama untuk resep
func DFS(elementName string, recipeType string, maxRecipes int) ([]interface{}, int, float64) {
    start := time.Now()
    visited := make(map[string]bool)
    stack := []*Node{{Name: elementName}}
    results := []interface{}{}
    nodesVisited := 0

    for len(stack) > 0 && (recipeType != "Limit" || len(results) < maxRecipes) {
        node := stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        if visited[node.Name] {
            continue
        }
        visited[node.Name] = true
        nodesVisited++

        children, err := getChildren(node.Name)
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

        stack = append(stack, children...)
        time.Sleep(100 * time.Millisecond) // Simulasi penundaan
    }

    return results, nodesVisited, float64(time.Since(start).Milliseconds())
}

// Bidirectional melakukan pencarian dua arah untuk resep
func Bidirectional(elementName string, targetName string, recipeType string, maxRecipes int) ([]interface{}, int, float64) {
    start := time.Now()
    forwardVisited := make(map[string]bool)
    backwardVisited := make(map[string]bool)
    forwardQueue := []*Node{{Name: elementName}}
    backwardQueue := []*Node{{Name: targetName}}
    results := []interface{}{}
    nodesVisited := 0

    for (len(forwardQueue) > 0 || len(backwardQueue) > 0) && (recipeType != "Limit" || len(results) < maxRecipes) {
        // Langkah maju
        if len(forwardQueue) > 0 {
            node := forwardQueue[0]
            forwardQueue = forwardQueue[1:]
            if forwardVisited[node.Name] {
                continue
            }
            forwardVisited[node.Name] = true
            nodesVisited++

            children, err := getChildren(node.Name)
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
            forwardQueue = append(forwardQueue, children...)
        }

        // Langkah mundur
        if len(backwardQueue) > 0 {
            node := backwardQueue[0]
            backwardQueue = backwardQueue[1:]
            if backwardVisited[node.Name] {
                continue
            }
            backwardVisited[node.Name] = true
            nodesVisited++

            parents, err := getParents(node.Name)
            if err != nil {
                log.Printf("Gagal mengambil induk untuk elemen %s: %v", node.Name, err)
                continue
            }
            backwardQueue = append(backwardQueue, parents...)
        }

        time.Sleep(100 * time.Millisecond) // Simulasi penundaan
    }

    return results, nodesVisited, float64(time.Since(start).Milliseconds())
}