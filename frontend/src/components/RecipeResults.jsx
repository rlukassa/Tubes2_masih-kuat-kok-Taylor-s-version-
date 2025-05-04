export default function RecipeResults({
    selectedElement,
    algorithm,
    recipeType,
    progress,
    executionTime,
    nodesVisited,
  }) {
    return (
      <div className="recipe-details">
        <h2>Recipe Details</h2>
  
        <div className="detail-section">
          <h3>Target Element</h3>
          <div className="detail-value">{selectedElement?.name || "None"}</div>
        </div>
  
        <div className="detail-section">
          <h3>Algorithm</h3>
          <div className="detail-value">{algorithm}</div>
        </div>
  
        <div className="detail-section">
          <h3>Recipe Type</h3>
          <div className="detail-value">{recipeType === "Best" ? "Shortest Path" : recipeType}</div>
        </div>
  
        <div className="detail-section">
          <h3>Progress</h3>
          <div className="progress-bar">
            <div className="progress-fill" style={{ width: `${progress}%` }}></div>
          </div>
        </div>
  
        <div className="metrics">
          <div className="metric">
            <div className="metric-icon">⏱️</div>
            <div className="metric-details">
              <h3>Execution Time</h3>
              <div className="metric-value">{executionTime}ms</div>
            </div>
          </div>
  
          <div className="metric">
            <div className="metric-icon">🔍</div>
            <div className="metric-details">
              <h3>Nodes Visited</h3>
              <div className="metric-value">{nodesVisited}</div>
            </div>
          </div>
        </div>
      </div>
    )
  }
  