package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "main/services" // Sesuaikan dengan path package services Anda
    "log"
)

func GetElements(c *gin.Context) {
    elements, err := services.FetchElements()
    if err != nil {
        log.Printf("Error fetching elements: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if elements == nil {
        log.Println("No elements found in database")
        c.JSON(http.StatusOK, []interface{}{})
        return
    }
    c.JSON(http.StatusOK, elements)
}

func SearchRecipes(c *gin.Context) {
    var request struct {
        ElementName string `json:"elementName"`
        TargetName  string `json:"targetName"` // Untuk Bidirectional
        Algorithm   string `json:"algorithm"`
        RecipeType  string `json:"recipeType"`
        MaxRecipes  int    `json:"maxRecipes"`
    }
    if err := c.BindJSON(&request); err != nil {
        log.Printf("Invalid request: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var results []interface{}
    var nodesVisited int
    var executionTime float64

    switch request.Algorithm {
    case "BFS":
        results, nodesVisited, executionTime = services.BFS(request.ElementName, request.RecipeType, request.MaxRecipes)
    case "DFS":
        results, nodesVisited, executionTime = services.DFS(request.ElementName, request.RecipeType, request.MaxRecipes)
    case "Bidirectional":
        if request.TargetName == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "TargetName required for Bidirectional search"})
            return
        }
        results, nodesVisited, executionTime = services.Bidirectional(request.ElementName, request.TargetName, request.RecipeType, request.MaxRecipes)
    default:
        log.Printf("Invalid algorithm: %s", request.Algorithm)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid algorithm"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "results":      results,
        "nodesVisited": nodesVisited,
        "executionTime": executionTime,
    })
}