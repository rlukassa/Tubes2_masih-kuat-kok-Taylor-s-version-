// src/components/ControlsPanel.jsx
"use client";

import { useState } from "react";

export default function ControlsPanel({ searchParams, setSearchParams }) {
  const [maxRecipes, setMaxRecipes] = useState(5);

  const handleAlgorithmChange = (algorithm) => {
    setSearchParams((prev) => ({ ...prev, algorithm }));
  };

  const handleRecipeTypeChange = (recipeType) => {
    setSearchParams((prev) => ({ ...prev, recipeType }));
  };

  const handleMaxRecipesChange = (e) => {
    const value = Number.parseInt(e.target.value);
    setMaxRecipes(value);
    setSearchParams((prev) => ({ ...prev, maxRecipes: value }));
  };

  return (
    <div className="controls-panel">
      <div className="algorithm-options">
        <h3>Algorithm Options</h3>
        <div className="option">
          <input
            type="radio"
            id="bfs"
            name="algorithm"
            checked={searchParams.algorithm === "BFS"}
            onChange={() => handleAlgorithmChange("BFS")}
          />
          <label htmlFor="bfs">
            <strong>BFS</strong>
            <div className="option-description">Breadth First Search</div>
          </label>
        </div>
        <div className="option">
          <input
            type="radio"
            id="dfs"
            name="algorithm"
            checked={searchParams.algorithm === "DFS"}
            onChange={() => handleAlgorithmChange("DFS")}
          />
          <label htmlFor="dfs">
            <strong>DFS</strong>
            <div className="option-description">Depth First Search</div>
          </label>
        </div>
        <div className="option">
          <input
            type="radio"
            id="bidirectional"
            name="algorithm"
            checked={searchParams.algorithm === "Bidirectional"}
            onChange={() => handleAlgorithmChange("Bidirectional")}
          />
          <label htmlFor="bidirectional">
            <strong>Bidirectional</strong>
            <div className="option-description">Search from both ends</div>
          </label>
        </div>
      </div>
      <div className="recipe-options">
        <h3>Recipe Options</h3>
        <div className="option">
          <input
            type="radio"
            id="one-recipe"
            name="recipeType"
            checked={searchParams.recipeType === "One"}
            onChange={() => handleRecipeTypeChange("One")}
          />
          <label htmlFor="one-recipe">
            <strong>One Recipe</strong>
            <div className="option-description">Find a recipe path for the element</div>
          </label>
        </div>
        <div className="option">
          <input
            type="radio"
            id="all-recipes"
            name="recipeType"
            checked={searchParams.recipeType === "All"}
            onChange={() => handleRecipeTypeChange("All")}
          />
          <label htmlFor="all-recipes">
            <strong>All Recipes</strong>
            <div className="option-description">Find all possible recipes for the element</div>
          </label>
        </div>
        <div className="option">
          <input
            type="radio"
            id="limit-recipes"
            name="recipeType"
            checked={searchParams.recipeType === "Limit"}
            onChange={() => handleRecipeTypeChange("Limit")}
          />
          <label htmlFor="limit-recipes">
            <strong>Limit Recipes</strong>
            <div className="option-description">Limit to a specified number for different recipes</div>
          </label>
        </div>
        {searchParams.recipeType === "Limit" && (
          <div className="max-recipes">
            <label htmlFor="max-recipes">Maximum Recipes</label>
            <input
              type="number"
              id="max-recipes"
              min="1"
              max="30"
              value={maxRecipes}
              onChange={handleMaxRecipesChange}
            />
          </div>
        )}
      </div>
    </div>
  );
}