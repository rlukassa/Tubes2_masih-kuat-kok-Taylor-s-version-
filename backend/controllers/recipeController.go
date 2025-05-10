// File ini adalah controller untuk endpoint pencarian resep pada backend Little Alchemy 2.
// Fungsinya menerima request pencarian dari frontend, memanggil service pencarian (BFS, DFS, Bidirectional),
// dan mengembalikan hasil pencarian dalam format JSON ke frontend.

package controllers

import (
  "net/http"              // Untuk kebutuhan HTTP response 
  "github.com/gin-gonic/gin" // Framework web Gin
  "main/services"         // Import service pencarian resep
)

func SearchRecipe(c *gin.Context) {
  var requestBody struct {
    ElementName string `json:"elementName"` // Nama elemen yang dicari
    Algorithm   string `json:"algorithm"`   // Algoritma pencarian (BFS, DFS, Bidirectional)
    RecipeType  string `json:"recipeType"`  // Tipe resep (misal: One Recipe)
    MaxRecipes  int    `json:"maxRecipes"`  // Maksimal jumlah resep -- buat RecipeType = "Limit .. "
    TargetName  string `json:"targetName"`  // Target untuk buat Algoritma Bidirectional 
  }

  if err := c.ShouldBindJSON(&requestBody); err != nil { // Bind dan validasi request body dari frontend
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"}) // Jika gagal, kirim error 400
    return
  }

  var results []interface{}    // Untuk menampung hasil pencarian -- menyimpan array (tree) resep ketika ditemukan
  var nodesVisited int         // Untuk menghitung node yang dikunjungi 
  var executionTime float64    // Untuk mencatat waktu eksekusi

  switch requestBody.Algorithm { // Pilih algoritma pencarian sesuai permintaan frontend
  case "BFS":
    results, nodesVisited, executionTime = services.BFS(requestBody.ElementName, requestBody.RecipeType, requestBody.MaxRecipes) // Panggil BFS
  case "DFS":
    results, nodesVisited, executionTime = services.DFS(requestBody.ElementName, requestBody.RecipeType, requestBody.MaxRecipes) // Panggil DFS
  case "Bidirectional":
    if requestBody.TargetName == "" { // Validasi jika Bidirectional harus ada target
      c.JSON(http.StatusBadRequest, gin.H{"error": "TargetName required for Bidirectional"})
      return
    }
    results, nodesVisited, executionTime = services.Bidirectional(requestBody.ElementName, requestBody.TargetName, requestBody.RecipeType, requestBody.MaxRecipes) // Panggil Bidirectional
  default:
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid algorithm"}) // Jika algoritma tidak valid, kirim error 400
    return
  }

  c.JSON(http.StatusOK, gin.H{ // Kirim hasil pencarian ke frontend dalam format JSON
    "results":       results,        // Hasil pencarian (array pohon resep)
    "nodesVisited":  nodesVisited,   // Jumlah node yang dikunjungi
    "executionTime": executionTime,  // Lama waktu eksekusi (ms)
  })
}