// Mock data for elements
const MOCK_ELEMENTS = [
    { id: 1, name: "Water", category: "Basic" },
    { id: 2, name: "Fire", category: "Basic" },
    { id: 3, name: "Earth", category: "Basic" },
    { id: 4, name: "Air", category: "Basic" },
    { id: 5, name: "Heat", category: "Derived" },
    { id: 6, name: "Ice", category: "Derived" },
    { id: 7, name: "Snow", category: "Derived" },
    { id: 8, name: "Coal", category: "Derived" },
    { id: 9, name: "Mud", category: "Derived" },
    { id: 10, name: "Stone", category: "Derived" },
    { id: 11, name: "Clay", category: "Derived" },
    { id: 12, name: "Brick", category: "Derived" },
    { id: 13, name: "Sand", category: "Derived" },
    { id: 14, name: "Glass", category: "Derived" },
    { id: 15, name: "Metal", category: "Derived" },
    { id: 16, name: "Steam", category: "Derived" },
  ]
  
  /**
   * Fetch all elements from the API or use mock data
   * @returns {Promise<Array>} Array of elements
   */
  export const fetchElements = async () => {
    try {
      // Try to fetch from API first
      const response = await fetch("/api/elements")
      if (response.ok) {
        return await response.json()
      }
    } catch (error) {
      console.warn("Failed to fetch elements from API, using mock data", error)
    }
  
    // Return mock data if API fails
    return MOCK_ELEMENTS
  }
  
  /**
   * Format time in milliseconds to a readable format
   * @param {number} ms - Time in milliseconds
   * @returns {string} Formatted time
   */
  export const formatTime = (ms) => {
    if (ms < 1000) {
      return `${ms}ms`
    }
  
    const seconds = Math.floor(ms / 1000)
    const remainingMs = ms % 1000
  
    if (seconds < 60) {
      return `${seconds}.${remainingMs.toString().padStart(3, "0")}s`
    }
  
    const minutes = Math.floor(seconds / 60)
    const remainingSeconds = seconds % 60
  
    return `${minutes}m ${remainingSeconds}s`
  }
  
  /**
   * Format large numbers with commas
   * @param {number} num - Number to format
   * @returns {string} Formatted number
   */
  export const formatNumber = (num) => {
    return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",")
  }
  