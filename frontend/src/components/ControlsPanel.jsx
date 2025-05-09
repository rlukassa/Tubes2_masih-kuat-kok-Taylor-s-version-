"use client"

import { useState } from "react"

// ====================================== CATATAN =======================================
// ini buat page 2 yang ada di sebelah kanan, ini adalah panel kontrol yang berfungsi untuk mengatur algoritma pencarian dan jenis resep yang akan ditampilkan.

export default function ControlsPanel({ searchParams, setSearchParams }) {
  const [maxRecipes, setMaxRecipes] = useState(5)   // ini default value, bisa diubah sesuai kebutuhan, nilai ini akan digunakan untuk mengatur jumlah maksimum resep yang ditampilkan ketika recipeType adalah "Limit"

  const handleAlgorithmChange = (algorithm) => {
    setSearchParams((prev) => ({ ...prev, algorithm }))
  } // fungsi ini digunakan untuk mengubah algoritma pencarian yang digunakan dalam aplikasi. Ketika pengguna memilih algoritma baru, fungsi ini akan memperbarui parameter pencarian dengan algoritma yang dipilih.

  const handleRecipeTypeChange = (recipeType) => {
    setSearchParams((prev) => ({ ...prev, recipeType }))
  } // fungsi ini digunakan untuk mengubah jenis resep yang akan dicari. Ketika pengguna memilih jenis resep baru, fungsi ini akan memperbarui parameter pencarian dengan jenis resep yang dipilih.

  const handleMaxRecipesChange = (e) => {
    const value = Number.parseInt(e.target.value)
    setMaxRecipes(value)
    setSearchParams((prev) => ({ ...prev, maxRecipes: value }))
  } // fungsi ini digunakan untuk mengubah jumlah maksimum resep yang akan ditampilkan ketika jenis resep yang dipilih adalah "Limit". Ketika pengguna mengubah nilai input, fungsi ini akan memperbarui nilai maksimum resep dan juga memperbarui parameter pencarian.

  return (
    <div className="controls-panel"> 
      <div className="algorithm-options"> 
        <h3>Algorithm Options</h3>
{/* Pilihan algoritma yang dipilih  */}
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
  )
}