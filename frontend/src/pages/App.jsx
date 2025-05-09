"use client";

import { useState, useMemo } from "react";
import ElementPicker from "../components/ElementPicker";
import ControlsPanel from "../components/ControlsPanel";
import RecipeResults from "../components/RecipeResults";
import TreeVisualizer from "../components/TreeVisualizer";
import { useSearch } from "../hooks/useSearch";
import testTubeIcon from "../assets/test-tube.png";
import "../../public/App.css";
import "../assets/background.css";
import mapper from "../../../database/mapper2.json"; // Pastikan mapper diimpor

function App() {
  const [currentView, setCurrentView] = useState("landing"); // 'landing', 'search', 'results'
  const [selectedElement, setSelectedElement] = useState(null);
  const {
    searchParams,
    setSearchParams,
    searchResults,
    isLoading,
    executionTime,
    nodesVisited,
    progress,
    startSearch,
    resetSearch,
  } = useSearch();

  // Proses searchResults untuk menambahkan ikon dari mapper2.json
  const processedResults = useMemo(() => {
    return (searchResults || []).map((result) => ({
      ...result,
      icon: mapper[result.name] || "", // Ambil URL ikon dari mapper2.json
      children: (result.children || []).map((child) => ({
        ...child,
        icon: mapper[child.name] || "",
      })),
    }));
  }, [searchResults]);

  const handleStartExploring = () => {
    setCurrentView("search");
  };

  const handleStartSearch = () => {
    if (!selectedElement) return;
    startSearch(selectedElement);
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
          onElementSelect={setSelectedElement}
        />

        <div className="controls-container">
          <ControlsPanel searchParams={searchParams} setSearchParams={setSearchParams} />

          <button
            className="search-button clickable-button"
            onClick={handleStartSearch}
            disabled={!selectedElement}
          >
            <span className="search-icon">🔍</span> Start Search
          </button>
        </div>
      </div>
    </div>
  );

  const renderResultsPage = () => {
    // Proses elemen yang dipilih untuk mendapatkan nama dan ikon dari mapper2.json
    const selectedElementData = selectedElement
      ? { name: selectedElement.name || selectedElement, icon: mapper[selectedElement.name || selectedElement] }
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
  
        <div className="results-content">
          <RecipeResults
            selectedElement={selectedElement}
            algorithm={searchParams.algorithm}
            recipeType={searchParams.recipeType}
            progress={progress}
            executionTime={executionTime}
            nodesVisited={nodesVisited}
          />
  
          <div className="visualization-container">
            <h2 className="visualization-title">Recipe Visualization</h2>
            <TreeVisualizer
              results={processedResults} // Gunakan processedResults
              selectedElement={selectedElementData} // Teruskan elemen yang dipilih
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