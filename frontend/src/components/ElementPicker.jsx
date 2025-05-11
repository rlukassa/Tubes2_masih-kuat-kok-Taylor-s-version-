// src/components/ElementPicker.jsx
"use client";

import { useState, useEffect } from "react";
import mapper from "../../../database/mapper2.json";

export default function ElementPicker({ algorithm, onElementSelect }) {
  const [searchTerm, setSearchTerm] = useState("");
  const [filteredElements, setFilteredElements] = useState([]);
  const [selectedElements, setSelectedElements] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const elementsPerPage = 8;
  const [totalPages, setTotalPages] = useState(1);

  useEffect(() => {
    const allElements = Object.keys(mapper).map((key) => ({
      name: key,
      icon: mapper[key],
    }));
    setFilteredElements(allElements);
    setTotalPages(Math.ceil(allElements.length / elementsPerPage));
  }, []);

  useEffect(() => {
    const allElements = Object.keys(mapper).map((key) => ({
      name: key,
      icon: mapper[key],
    }));
    const filtered = allElements.filter((element) =>
      element.name.toLowerCase().includes(searchTerm.toLowerCase())
    );
    setFilteredElements(filtered);
    setTotalPages(Math.ceil(filtered.length / elementsPerPage));
    setCurrentPage(1);
  }, [searchTerm]);

  const handleElementClick = (element) => {
    let newSelectedElements;
    newSelectedElements = [element];
    setSelectedElements(newSelectedElements);
    onElementSelect(newSelectedElements);
  };

  const handleSearchChange = (e) => {
    setSearchTerm(e.target.value);
  };

  const handlePageChange = (page) => {
    setCurrentPage(page);
  };

  const getElementsForCurrentPage = () => {
    const startIndex = (currentPage - 1) * elementsPerPage;
    return filteredElements.slice(startIndex, startIndex + elementsPerPage);
  };

  const renderPaginationButtons = () => {
    const buttons = [];
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
      </button>
    );
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
        </button>
      );
    }
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
        </button>
      );
    }
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
        </span>
      );
    }
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
        </button>
      );
    }
    return buttons;
  };

  return (
    <div className="element-picker-container" style={{ padding: "20px" }}>
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
              ...(selectedElements[0]?.name === element.name
                ? { border: "2px solid blue" }
                : selectedElements[1]?.name === element.name
                ? { border: "2px solid red" }
                : { border: "1px solid #ccc" }),
            }}
          >
            <div className="element-icon">
              <img src={element.icon || "/placeholder.svg"} alt={element.name} style={{ width: "70px", height: "70px" }} />
            </div>
            <div className="element-name" style={{ marginTop: "10px", fontSize: "14px" }}>
              {element.name}
            </div>
          </div>
        ))}
      </div>
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
  );
}