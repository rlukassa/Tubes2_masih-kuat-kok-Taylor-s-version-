"use client"

import { useState, useEffect } from "react"

export function useSearch(initialItems = []) {
  const [searchParams, setSearchParams] = useState({
    algorithm: "BFS",
    recipeType: "Best",
    maxRecipes: 5,
  })

  // Search results and metrics for recipe search
  const [searchResults, setSearchResults] = useState([])
  const [isLoading, setIsLoading] = useState(false)
  const [executionTime, setExecutionTime] = useState(0)
  const [nodesVisited, setNodesVisited] = useState(0)
  const [progress, setProgress] = useState(0)

  // element filtering functionality
  const [searchTerm, setSearchTerm] = useState("")
  const [items, setItems] = useState(initialItems)
  const [filteredItems, setFilteredItems] = useState(initialItems)

  // Update items when initialItems changes
  useEffect(() => {
    setItems(initialItems)
    setFilteredItems(initialItems) // Show all items by default
  }, [initialItems])

  // Filter items whenever searchTerm change
  useEffect(() => {
    if (!searchTerm.trim()) {
      // When search is empty, show all items
      setFilteredItems(items)
      return
    }

    const searchTermLower = searchTerm.toLowerCase()
    
    console.log("Search term:", searchTermLower)
    console.log("Available items:", items)
    
    const filtered = items.filter(item => {
      // Handle null 
      if (!item) return false
      
      // Check if item has a name property
      if (typeof item === 'object') {
        const itemName = item.name ? String(item.name) : '';
        
        const nameMatch = itemName.toLowerCase().includes(searchTermLower)
        
        console.log(`Item "${itemName}", matches: ${nameMatch}`)
        
        const itemDesc = item.description ? String(item.description) : '';
        
        // Check if description contains search term 
        const descriptionMatch = itemDesc.toLowerCase().includes(searchTermLower)
        
        return nameMatch || descriptionMatch
      } else if (typeof item === 'string') {
        return item.toLowerCase().includes(searchTermLower)
      }
      
      return false
    })

    console.log("Filtered items:", filtered)
    setFilteredItems(filtered)
  }, [searchTerm, items])

  const startSearch = async (element) => {
    setIsLoading(true)
    setSearchResults([])
    setExecutionTime(0)
    setNodesVisited(0)
    setProgress(0)

    try {
      // Simulate progress updates
      const progressInterval = setInterval(() => {
        setProgress((prev) => {
          if (prev >= 95) {
            clearInterval(progressInterval)
            return 95
          }
          return prev + 5
        })
      }, 200)

      // Simulate API call to backend
      const startTime = performance.now()

      const response = await fetch(`/api/search`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          elementId: element.id,
          algorithm: searchParams.algorithm,
          recipeType: searchParams.recipeType,
          maxRecipes: searchParams.maxRecipes,
        }),
      })

      const data = await response.json()
      const endTime = performance.now()

      clearInterval(progressInterval)
      setProgress(100)

      // Set results
      setSearchResults(data.results || [])
      setExecutionTime(Math.round(endTime - startTime))
      setNodesVisited(data.nodesVisited || 0)

      // Simulate delay 
      setTimeout(() => {
        setIsLoading(false)
      }, 500)
    } catch (error) {
      console.error("Search error:", error)
      setIsLoading(false)
      setProgress(0)
    }
  }

  const resetSearch = () => {
    setSearchResults([])
    setExecutionTime(0)
    setNodesVisited(0)
    setProgress(0)
  }

  return {
    searchParams,
    setSearchParams,
    searchResults,
    isLoading,
    executionTime,
    nodesVisited,
    progress,
    startSearch,
    resetSearch,
    
    searchTerm,
    setSearchTerm,
    items,
    setItems,
    filteredItems,
    hasResults: filteredItems.length > 0,
    isSearching: searchTerm.trim() !== ""
  }
}