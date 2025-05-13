// src/hooks/useSearch.js
// Custom React hook untuk mengelola proses pencarian resep di aplikasi Little Alchemy 2.
// Menyediakan state dan fungsi untuk mengatur parameter pencarian, hasil, status loading, serta eksekusi pencarian ke backend.

import { useState, useEffect } from "react";

export function useSearch() {
  // State untuk parameter pencarian (algoritma, tipe resep, jumlah maksimal resep)
  const [searchParams, setSearchParams] = useState({
    algorithm: "BFS",      // Algoritma default: BFS
    recipeType: "One",     // Tipe resep default: One
    maxRecipes: 5,         // Jumlah maksimal resep default: 5
  });

  // State untuk hasil pencarian dan status proses
  const [searchResults, setSearchResults] = useState([]); // Menyimpan hasil pencarian dari backend
  const [isLoading, setIsLoading] = useState(false);      // Status loading saat pencarian berlangsung
  const [executionTime, setExecutionTime] = useState(0);  // Lama waktu eksekusi pencarian (ms)
  const [nodesVisited, setNodesVisited] = useState(0);    // Jumlah node yang dikunjungi selama pencarian
  const [progress, setProgress] = useState(0);            // Progress pencarian (untuk animasi/progress bar)
  const [totalRecipes, setTotalRecipes] = useState(0);    // Total resep yang ditemukan
  // Fungsi untuk memulai pencarian resep
  const startSearch = async (elements) => {
    setIsLoading(true);           // Set status loading
    setSearchResults([]);         // Reset hasil sebelumnya
    setExecutionTime(0);          // Reset waktu eksekusi    setNodesVisited(0);           // Reset jumlah node
    setProgress(0);               // Reset progress
    setTotalRecipes(0);           // Reset total resep

    try {
      // Siapkan request body untuk dikirim ke backend
      const requestBody = {
        elementName: elements[0].name,           // Nama elemen awal
        algorithm: searchParams.algorithm,       // Algoritma pencarian
        recipeType: searchParams.recipeType,     // Tipe resep
        maxRecipes: searchParams.maxRecipes,     // Jumlah maksimal resep
      };
      // Jika algoritma Bidirectional, tambahkan targetName
      if (searchParams.algorithm === "Bidirectional" && elements.length > 1) {
        requestBody.targetName = elements[1].name;
      }

      // Kirim request ke backend (pastikan port dan endpoint sudah benar)
      const response = await fetch("/api/search", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify(requestBody),
});

      if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
      const data = await response.json();      // Simpan hasil pencarian ke state
      setSearchResults(data.results);
      setNodesVisited(data.nodesVisited);
      setExecutionTime(data.executionTime);
      setTotalRecipes(data.results ? data.results.length : 0); // Set total recipes
      setProgress(100); // Progress selesai
    } catch (error) {
      console.error("Search error:", error);
      alert("Gagal melakukan pencarian. Pastikan backend berjalan.");
    } finally {
      setIsLoading(false); // Selesai loading
    }
  };

  // Fungsi untuk mereset hasil pencarian dan progress
  const resetSearch = () => {
    setSearchResults([]);
    setExecutionTime(0);
    setNodesVisited(0);
    setProgress(0);
    setTotalRecipes(0);
  };

  // Return semua state dan fungsi yang dibutuhkan komponen lain
  return {
    searchParams,      // Parameter pencarian
    setSearchParams,   // Setter parameter pencarian
    searchResults,     // Hasil pencarian
    isLoading,         // Status loading
    executionTime,     // Lama eksekusi
    nodesVisited,      // Jumlah node dikunjungi
    progress,          // Progress pencarian
    totalRecipes,      // Total resep yang ditemukan
    startSearch,       // Fungsi untuk memulai pencarian
    resetSearch,       // Fungsi untuk mereset hasil
  };
}