// RecipeResults.jsx with improved responsiveness
import React from "react";

export default function RecipeResults({
  selectedElement,
  algorithm,
  recipeType,
  progress,
  executionTime,
  nodesVisited,
  totalRecipes,
}) {
  return (
    <div className="recipe-details">
      <h2>Recipe Details</h2>

      <div className="detail-section">
        <h3>Target Element(s)</h3>
        <div className="detail-value target-element">
          {Array.isArray(selectedElement) ? (
            selectedElement.map((el, index) => (
              <span
                key={index}
                className="element-item-detail"
                style={{
                  display: "inline-flex",
                  alignItems: "center",
                  margin: "0 10px",
                }}
              >
                <img
                  src={el.icon}
                  alt={el.name}
                  style={{
                    width: "30px",
                    height: "30px",
                    marginRight: "10px",
                  }}
                />
                {el.name}
                {index === 0 && selectedElement.length > 1 ? " → " : ""}
              </span>
            ))
          ) : selectedElement ? (
            <span
              className="element-item-detail"
              style={{
                display: "inline-flex",
                alignItems: "center",
              }}
            >
              <img
                src={selectedElement.icon}
                alt={selectedElement.name}
                style={{
                  width: "60px",
                  height: "60px",
                  marginRight: "10px",
                }}
                className="target-element-icon"
              />
              {selectedElement.name}
            </span>
          ) : (
            "None"
          )}
        </div>
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

      {/* Total Recipes info */}
      <div className="detail-section">
        <h3>Total Recipes</h3>
        <div className="detail-value">{totalRecipes || 0}</div>
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
  );
}