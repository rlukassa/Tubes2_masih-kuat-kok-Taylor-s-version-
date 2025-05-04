"use client"

import { useState } from "react"
import ElementPicker from "../components/ElementPicker"
import ControlsPanel from "../components/ControlsPanel"
import SearchBar from "../components/SearchBar"
import RecipeResults from "../components/RecipeResults"
import TreeVisualizer from "../components/TreeVisualizer"
import { useSearch } from "../hooks/useSearch"
import testTubeIcon from "../assets/test-tube.png"
import "../../public/App.css"

export default function App() {
  const [currentView, setCurrentView] = useState("landing") // 'landing', 'search', 'results'
  const [selectedElement, setSelectedElement] = useState(null)
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
  } = useSearch()

  const handleStartExploring = () => {
    setCurrentView("search")
  }

  const handleStartSearch = () => {
    if (!selectedElement) return
    startSearch(selectedElement)
    setCurrentView("results")
  }

  const handleBackToSearch = () => {
    resetSearch()
    setCurrentView("search")
  }

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

        <button className="start-button" onClick={handleStartExploring}>
          <img src={testTubeIcon || "/placeholder.svg"} alt="" className="button-icon" /> Start Exploring
        </button>
      </div>
    </div>
  )

  const renderSearchPage = () => (
    <div className="search-container">
      <h1 className="search-title">Search Elements</h1>
      <p className="search-subtitle">Select an element below to find its recipes</p>

      <div className="search-content">
        <SearchBar onElementSelect={setSelectedElement} />

        <div className="search-grid">
          <ElementPicker onElementSelect={setSelectedElement} selectedElement={selectedElement} />

          <div className="controls-container">
            <ControlsPanel searchParams={searchParams} setSearchParams={setSearchParams} />

            <button className="search-button" onClick={handleStartSearch} disabled={!selectedElement}>
              <span className="search-icon">🔍</span> Start Search
            </button>
          </div>
        </div>
      </div>
    </div>
  )

  const renderResultsPage = () => (
    <div className="results-container">
      <div className="results-header">
        <button className="back-button" onClick={handleBackToSearch}>
          ← Back to Search
        </button>
        <button className="print-button">
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
          <TreeVisualizer results={searchResults} selectedElement={selectedElement} isLoading={isLoading} />
        </div>
      </div>
    </div>
  )

  return (
    <div className="app">
      {currentView === "landing" && renderLandingPage()}
      {currentView === "search" && renderSearchPage()}
      {currentView === "results" && renderResultsPage()}
    </div>
  )
}
