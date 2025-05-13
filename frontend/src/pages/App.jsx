// src/pages/App.jsx
"use client";

import React, { useEffect, useState, useMemo } from "react";
import ElementPicker from "../components/ElementPicker";
import ControlsPanel from "../components/ControlsPanel";
import RecipeResults from "../components/RecipeResults";
import TreeVisualizer from "../components/TreeVisualizer";
import { useSearch } from "../hooks/useSearch";
import testTubeIcon from "../assets/test-tube.png";
import "../../public/App.css";
import "../assets/background.css";
import mapper from "../../../database/mapper2.json";

function App() {
  const [currentView, setCurrentView] = useState("landing");
  const [selectedElements, setSelectedElements] = useState([]);
  const [liveTreeData, setLiveTreeData] = useState([]);  const {
    searchParams,
    setSearchParams,
    searchResults,
    isLoading,
    executionTime,
    nodesVisited,
    progress,
    totalRecipes,
    startSearch,
    resetSearch,
  } = useSearch();

// Mengubah hasil pencarian dari backend (searchResults) agar setiap node dan child-nya
// memiliki properti icon, sehingga TreeVisualizer bisa menampilkan gambar/icon elemen.
// useMemo digunakan agar proses ini hanya dijalankan ulang jika searchResults berubah.
const processedResults = useMemo(() => {
  return (searchResults || []).map((result) => ({
    ...result,
    icon: mapper[result.name] || "", // Tambahkan icon untuk node utama
    children: (result.children || []).map((child) => ({
      ...child,
      icon: mapper[child.name] || "", // Tambahkan icon untuk setiap child
    })),
  }));
}, [searchResults]);

  // Efek untuk melakukan animasi "live update" pada visualisasi pohon.
// Setiap 100ms, satu node dari processedResults akan ditambahkan ke liveTreeData,
// sehingga pohon akan muncul secara bertahap (mirip animasi progres pencarian).
useEffect(() => {
  if (!processedResults.length || isLoading) return;
  setLiveTreeData(processedResults);
}, [processedResults, isLoading]);

  // Fungsi untuk mengubah tampilan ke halaman pencarian saat tombol "Start Exploring" ditekan.
  const handleStartExploring = () => {
    setCurrentView("search");
  };

  const handleStartSearch = () => {
    if (!selectedElements.length) return;
    startSearch(selectedElements);
    setCurrentView("results");
  };

  const handleBackToSearch = () => {
    resetSearch();
    setCurrentView("search");
  };

  const renderLandingPage = () => (
    <div className="landing-container">
      <div className="landing-content">
        <h1 className="discover-text">Discover the Secrets of</h1>
        <div className="title-container">
          <h1 className="title">Little Alchemy 2</h1>
          <img src={testTubeIcon || "/placeholder.svg"} alt="Test tube" className="test-tube-icon" />
        </div>
        <h2 className="subtitle">Recipe Finder</h2>
        <div className="basic-elements">
          <div className="element-card">
            <div className="element-icon water-icon">💧</div>
            <h3>Water</h3>
            <p className="element-description">The fluid of life</p>
          </div>
          <div className="element-card">
            <div className="element-icon fire-icon">🔥</div>
            <h3>Fire</h3>
            <p className="element-description">The energy of transformation</p>
          </div>
          <div className="element-card">
            <div className="element-icon earth-icon">🌋</div>
            <h3>Earth</h3>
            <p className="element-description">The foundation of creation</p>
          </div>
          <div className="element-card">
            <div className="element-icon air-icon">💨</div>
            <h3>Air</h3>
            <p className="element-description">The breath of existence</p>
          </div>
        </div>
        <p className="description">
          Uncover the secrets of Little Alchemy 2 with our advanced recipe finder.
          <br />
          Explore paths from basics to complex creations using BFS, DFS, and Bidirectional search.
        </p>
        <button className="start-button clickable-button" onClick={handleStartExploring}>
          <img src={testTubeIcon || "icon.svg"} alt="" className="button-icon" /> Start Exploring
        </button>
      </div>
    </div>
  );

  const renderSearchPage = () => (
    <div className="search-container">
      <h1 className="search-title">Search Elements</h1>
      <p className="search-subtitle">Select an element below to find its recipes</p>
      <div className="search-grid">
        <ElementPicker
          algorithm={searchParams.algorithm}
          onElementSelect={setSelectedElements}
        />
        <div className="controls-container">
          <ControlsPanel searchParams={searchParams} setSearchParams={setSearchParams} />
          <button
            className="search-button clickable-button"
            onClick={handleStartSearch}
            disabled={!selectedElements.length}
          >
            <span className="search-icon">🔍</span> Start Search
          </button>
        </div>
      </div>
    </div>
  );

  const renderResultsPage = () => {
    const selectedElementData = selectedElements[0]
      ? { name: selectedElements[0].name, icon: mapper[selectedElements[0].name] }
      : null;

    return (
      <div className="results-container">
        <div className="results-header">
          <button className="back-button clickable-button" onClick={handleBackToSearch}>
            ← Back to Search
          </button>
          <button className="print-button clickable-button">
            <span className="print-icon">🖨️</span> Print Tree
          </button>
        </div>
        <div className="results-content">          <RecipeResults
            selectedElement={selectedElements}
            algorithm={searchParams.algorithm}
            recipeType={searchParams.recipeType}
            progress={progress}
            executionTime={executionTime}
            nodesVisited={nodesVisited}
            totalRecipes={totalRecipes}
          />
          <div className="visualization-container">
            <h2 className="visualization-title">Recipe Visualization</h2>
            <TreeVisualizer
              results={liveTreeData}
              selectedElement={selectedElementData}
              isLoading={isLoading}
            />
          </div>
        </div>
      </div>
    );
  };

  return (
    <div className="app bg-stars">
      {currentView === "landing" && renderLandingPage()}
      {currentView === "search" && renderSearchPage()}
      {currentView === "results" && renderResultsPage()}
    </div>
  );
}

export default App;