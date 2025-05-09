package main

import (
    "github.com/gin-gonic/gin"
    "main/controllers" // Adjust path
)

func main() {
    gin.SetMode(gin.ReleaseMode) // Set Gin to Release Mode

    r := gin.Default()
    r.Use(corsMiddleware())

    // Define routes
    r.GET("/elements", controllers.GetElements)
    r.POST("/search", controllers.SearchRecipes)

    r.Run(":5000") // Runs on localhost:5000
}

func corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    }
}