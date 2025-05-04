"use client"

import { useState, useEffect } from "react"
import { fetchElements } from "../utils/helpers"

export default function ElementPicker({ onElementSelect, selectedElement }) {
  const [elements, setElements] = useState([])
  const [currentPage, setCurrentPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const elementsPerPage = 8

  useEffect(() => {
    const loadElements = async () => {
      const allElements = await fetchElements()
      setElements(allElements)
      setTotalPages(Math.ceil(allElements.length / elementsPerPage))
    }

    loadElements()
  }, [])

  const getElementsForCurrentPage = () => {
    const startIndex = (currentPage - 1) * elementsPerPage
    return elements.slice(startIndex, startIndex + elementsPerPage)
  }

  const handlePageChange = (page) => {
    setCurrentPage(page)
  }

  const getElementIcon = (element) => {
    switch (element.name.toLowerCase()) {
      case "water":
        return "💧"
      case "fire":
        return "🔥"
      case "earth":
        return "🌋"
      case "air":
        return "💨"
      case "heat":
        return "♨️"
      case "ice":
        return "🧊"
      case "snow":
        return "❄️"
      case "coal":
        return "🪨"
      default:
        return "🧪"
    }
  }

  return (
    <div className="element-picker">
      <div className="elements-grid">
        {getElementsForCurrentPage().map((element) => (
          <div
            key={element.id}
            className={`element-item ${selectedElement?.id === element.id ? "selected" : ""}`}
            onClick={() => onElementSelect(element)}
          >
            <div className="element-icon">{getElementIcon(element)}</div>
            <div className="element-name">{element.name}</div>
          </div>
        ))}
      </div>

      <div className="pagination">
        <button
          className="pagination-button"
          onClick={() => handlePageChange(currentPage - 1)}
          disabled={currentPage === 1}
        >
          &lt;
        </button>

        {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
          const pageNumber =
            currentPage <= 3 ? i + 1 : currentPage >= totalPages - 2 ? totalPages - 4 + i : currentPage - 2 + i

          if (pageNumber <= totalPages) {
            return (
              <button
                key={pageNumber}
                className={`pagination-button ${currentPage === pageNumber ? "active" : ""}`}
                onClick={() => handlePageChange(pageNumber)}
              >
                {pageNumber}
              </button>
            )
          }
          return null
        })}

        <button
          className="pagination-button"
          onClick={() => handlePageChange(currentPage + 1)}
          disabled={currentPage === totalPages}
        >
          &gt;
        </button>
      </div>
    </div>
  )
}
