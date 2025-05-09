"use client"

import { useState, useEffect } from "react"
import mapper from "../../../database/mapper2.json"

export default function ElementPicker({ algorithm, onElementSelect }) {
  const [searchTerm, setSearchTerm] = useState("") // State untuk input pencarian
  const [filteredElements, setFilteredElements] = useState([]) // Elemen hasil pencarian
  const [selectedElements, setSelectedElements] = useState([]) // Elemen yang dipilih
  const [currentPage, setCurrentPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const elementsPerPage = 8

  // Inisialisasi elemen awal
  useEffect(() => {
    const allElements = Object.keys(mapper).map((key) => ({
      name: key,
      icon: mapper[key],
    }))
    setFilteredElements(allElements)
    setTotalPages(Math.ceil(allElements.length / elementsPerPage))
  }, [])

  // Fungsi untuk menangani perubahan input pencarian
  const handleSearchChange = (e) => {
    const value = e.target.value.toLowerCase()
    setSearchTerm(value)

    const results = Object.keys(mapper)
      .filter((key) => key.toLowerCase().includes(value))
      .map((key) => ({ name: key, icon: mapper[key] }))

    setFilteredElements(results)
    setCurrentPage(1) // Reset ke halaman pertama
    setTotalPages(Math.ceil(results.length / elementsPerPage))
  }

  const handleElementClick = (element) => {
    if (algorithm === "BFS" || algorithm === "DFS") {
      // Hanya satu elemen yang dapat dipilih
      setSelectedElements([element])
      onElementSelect([element])
    } else if (algorithm === "Bidirectional") {
      // Dua elemen harus dipilih
      if (selectedElements.length === 0) {
        setSelectedElements([element]) // Pilih elemen pertama
      } else if (selectedElements.length === 1) {
        setSelectedElements([...selectedElements, element]) // Pilih elemen kedua
      } else {
        setSelectedElements([element]) // Reset ke elemen pertama
      }
      onElementSelect(selectedElements)
    }
  }

  const getElementsForCurrentPage = () => {
    const startIndex = (currentPage - 1) * elementsPerPage
    return filteredElements.slice(startIndex, startIndex + elementsPerPage)
  }

  const handlePageChange = (page) => {
    setCurrentPage(page)
  }

  const getElementStyle = (element) => {
    if (selectedElements[0]?.name === element.name) {
      return { border: "2px solid blue" } // Elemen pertama
    } else if (selectedElements[1]?.name === element.name) {
      return { border: "2px solid red" } // Elemen kedua
    }
    return { border: "1px solid #ccc" } // Elemen lainnya
  }

  // Generate pagination buttons
  const renderPaginationButtons = () => {
    const buttons = []

    // Always show first page
    buttons.push(
      <button
        key={1}
        className={`pagination-button ${currentPage === 1 ? "active" : ""}`}
        onClick={() => handlePageChange(1)}
        style={{
          padding: "10px 15px",
          margin: "0 5px",
          border: "1px solid #ccc",
          borderRadius: "5px",
          backgroundColor: currentPage === 1 ? "#007BFF" : "#fff",
          color: currentPage === 1 ? "#fff" : "#000",
          cursor: "pointer",
        }}
      >
        1
      </button>,
    )

    // Show page 2 if total pages > 1
    if (totalPages > 1) {
      buttons.push(
        <button
          key={2}
          className={`pagination-button ${currentPage === 2 ? "active" : ""}`}
          onClick={() => handlePageChange(2)}
          style={{
            padding: "10px 15px",
            margin: "0 5px",
            border: "1px solid #ccc",
            borderRadius: "5px",
            backgroundColor: currentPage === 2 ? "#007BFF" : "#fff",
            color: currentPage === 2 ? "#fff" : "#000",
            cursor: "pointer",
          }}
        >
          2
        </button>,
      )
    }

    // Show page 3 if total pages > 2
    if (totalPages > 2) {
      buttons.push(
        <button
          key={3}
          className={`pagination-button ${currentPage === 3 ? "active" : ""}`}
          onClick={() => handlePageChange(3)}
          style={{
            padding: "10px 15px",
            margin: "0 5px",
            border: "1px solid #ccc",
            borderRadius: "5px",
            backgroundColor: currentPage === 3 ? "#007BFF" : "#fff",
            color: currentPage === 3 ? "#fff" : "#000",
            cursor: "pointer",
          }}
        >
          3
        </button>,
      )
    }

    // Show ellipsis if total pages > 4
    if (totalPages > 4) {
      buttons.push(
        <span
          key="ellipsis"
          style={{
            padding: "10px 15px",
            margin: "0 5px",
          }}
        >
          ...
        </span>,
      )
    }

    // Show last page if total pages > 3
    if (totalPages > 3) {
      buttons.push(
        <button
          key={totalPages}
          className={`pagination-button ${currentPage === totalPages ? "active" : ""}`}
          onClick={() => handlePageChange(totalPages)}
          style={{
            padding: "10px 15px",
            margin: "0 5px",
            border: "1px solid #ccc",
            borderRadius: "5px",
            backgroundColor: currentPage === totalPages ? "#007BFF" : "#fff",
            color: currentPage === totalPages ? "#fff" : "#000",
            cursor: "pointer",
          }}
        >
          {totalPages}
        </button>,
      )
    }

    return buttons
  }

  return (
    <div className="element-picker-container" style={{ padding: "20px" }}>
      {/* Search Bar */}
      <div className="search-bar-container" style={{ marginBottom: "20px" }}>
        <input
          type="text"
          className="search-input"
          placeholder="Search items..."
          value={searchTerm}
          onChange={handleSearchChange}
          style={{
            width: "100%",
            padding: "10px",
            fontSize: "16px",
            borderRadius: "5px",
            border: "1px solid #ccc",
          }}
        />
      </div>

      {/* Elements Grid */}
      <div className="elements-grid" style={{ display: "grid", gridTemplateColumns: "repeat(4, 1fr)", gap: "20px" }}>
        {getElementsForCurrentPage().map((element) => (
          <div
            key={element.name}
            className="element-item"
            onClick={() => handleElementClick(element)}
            style={{
              cursor: "pointer",
              textAlign: "center",
              padding: "10px",
              borderRadius: "5px",
              ...getElementStyle(element),
            }}
          >
            <div className="element-icon">
              <img
                src={element.icon || "/placeholder.svg"}
                alt={element.name}
                style={{ width: "70px", height: "70px" }}
              />
            </div>
            <div className="element-name" style={{ marginTop: "10px", fontSize: "14px" }}>
              {element.name}
            </div>
          </div>
        ))}
      </div>

      {/* Pagination */}
      <div className="pagination" style={{ marginTop: "20px", textAlign: "center" }}>
        <button
          className="pagination-button"
          onClick={() => handlePageChange(currentPage - 1)}
          disabled={currentPage === 1}
          style={{
            padding: "10px 15px",
            margin: "0 5px",
            border: "1px solid #ccc",
            borderRadius: "5px",
            backgroundColor: currentPage === 1 ? "#f0f0f0" : "#007BFF",
            color: currentPage === 1 ? "#ccc" : "#fff",
            cursor: currentPage === 1 ? "not-allowed" : "pointer",
          }}
        >
          &lt;
        </button>

        {renderPaginationButtons()}

        <button
          className="pagination-button"
          onClick={() => handlePageChange(currentPage + 1)}
          disabled={currentPage === totalPages}
          style={{
            padding: "10px 15px",
            margin: "0 5px",
            border: "1px solid #ccc",
            borderRadius: "5px",
            backgroundColor: currentPage === totalPages ? "#f0f0f0" : "#007BFF",
            color: currentPage === totalPages ? "#ccc" : "#fff",
            cursor: currentPage === totalPages ? "not-allowed" : "pointer",
          }}
        >
          &gt;
        </button>
      </div>
    </div>
  )
}
