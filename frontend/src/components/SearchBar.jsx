"use client"

import { useState, useEffect, useRef } from "react"
import { fetchElements } from "../utils/helpers"
import { useSearch } from "../hooks/useSearch"

export default function SearchBar({ onElementSelect }) {
  const [isDropdownOpen, setIsDropdownOpen] = useState(false)
  const dropdownRef = useRef(null)
  
  // Custom hook for search functionality
  const { 
    searchTerm, 
    setSearchTerm, 
    items: elements, 
    setItems: setElements, 
    filteredItems: filteredElements 
  } = useSearch([])

  useEffect(() => {
    const loadElements = async () => {
      try {
        const allElements = await fetchElements()
        console.log("Loaded elements:", allElements)
        setElements(allElements)
      } catch (error) {
        console.error("Error loading elements:", error)
      }
    }

    loadElements()
  }, [setElements])

  // dropdown visibility based on search
  useEffect(() => {
    setIsDropdownOpen(searchTerm.trim() !== "")
  }, [searchTerm])

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setIsDropdownOpen(false)
      }
    }

    document.addEventListener("mousedown", handleClickOutside)
    return () => {
      document.removeEventListener("mousedown", handleClickOutside)
    }
  }, [])

  const handleSearchChange = (e) => {
    const value = e.target.value
    console.log("Search input changed to:", value)
    setSearchTerm(value)
  }

  const handleElementSelect = (element) => {
    console.log("Element selected:", element)
    onElementSelect(element)
    setSearchTerm(element.name)
    setIsDropdownOpen(false)
  }

  return (
    <div className="search-bar-container" ref={dropdownRef}>
      <div className="search-input-wrapper">
        <input
          type="text"
          className="search-input"
          placeholder="Search elements..."
          value={searchTerm}
          onChange={handleSearchChange}
          onFocus={() => {
            console.log("Input focused, current search term:", searchTerm)
            setIsDropdownOpen(searchTerm.trim() !== "")
          }}
        />
        {searchTerm && (
          <button 
            className="clear-button" 
            onClick={() => {
              console.log("Clear button clicked")
              setSearchTerm("")
            }}
          >
            ×
          </button>
        )}
      </div>

      {isDropdownOpen && (
        <div className="search-dropdown">
          {filteredElements.length > 0 ? (
            filteredElements.map((element) => (
              <div 
                key={element.id} 
                className="dropdown-item" 
                onClick={() => handleElementSelect(element)}
              >
                {element.name}
                {element.description && (
                  <span className="element-description">{element.description}</span>
                )}
              </div>
            ))
          ) : (
            <div className="search-dropdown no-results">
              No elements found matching "{searchTerm}"
            </div>
          )}
        </div>
      )}
    </div>
  )
}