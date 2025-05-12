// src/components/TreeVisualizer.jsx
// Komponen visualisasi pohon resep menggunakan D3.js.
// Menampilkan hasil pencarian dalam bentuk tree langsung tanpa animasi.

"use client";

import React from "react";
import { useEffect, useRef, useState } from "react";
import * as d3 from "d3";

// Helper function to ensure all nodes have a children property
function ensureChildrenExist(node) {
  if (!node) return { name: "Unknown", children: [] };
  
  const newNode = { ...node };
  if (!newNode.children) {
    newNode.children = [];
  } else if (Array.isArray(newNode.children)) {
    newNode.children = newNode.children.map(ensureChildrenExist);
  }
  return newNode;
}

export default function TreeVisualizer({ 
  results, 
  selectedElement, 
  isLoading, 
  nodesVisited,  // Jumlah node yang telah dikunjungi
  searchAlgorithm // BFS, DFS, atau Bidirectional
}) {
  const svgRef = useRef(null); // Referensi ke elemen SVG untuk D3
  const containerRef = useRef(null); // Referensi ke container untuk mengambil ukuran
  const [zoom, setZoom] = useState(1); // State untuk level zoom
  
  // Render ulang tree setiap kali hasil, loading, atau zoom berubah
  useEffect(() => {
    renderTree();
  }, [results, isLoading, zoom, selectedElement]);
  
  // Fungsi untuk menghasilkan warna berdasarkan kedalaman node
  const getNodeColor = (depth) => {
    const colors = ["#4CAF50", "#2196F3", "#FFC107", "#E91E63", "#9C27B0", "#FF5722"];
    return colors[depth % colors.length];
  };
  
  // Fungsi utama untuk menggambar tree dengan D3
  const renderTree = () => {
    if (!svgRef.current || !containerRef.current || !results || results.length === 0) return;
    
    d3.select(svgRef.current).selectAll("*").remove(); // Bersihkan SVG sebelum render ulang
      const width = containerRef.current.clientWidth; // Lebar container
    const height = containerRef.current.clientHeight; // Tinggi container
    
    // Ensure each node in results has a children property
    const processedResults = results.map(ensureChildrenExist);
    
    // Data root tree langsung dari hasil
    const rootData = ensureChildrenExist({
      name: selectedElement?.name || "No Element Selected",
      image: selectedElement?.icon || "",
      children: processedResults
    });
    
    // Layout pohon dengan ukuran tertentu
    const treeLayout = d3.tree().size([height - 100, width - 200]);
    const root = d3.hierarchy(rootData);
    
    // Tambahkan depth sebagai properti untuk pewarnaan
    root.descendants().forEach((d, i) => {
      d.id = i;
      d.depth = d.depth;
    });
    
    treeLayout(root);
    
    // Inisialisasi SVG dan group utama
    const svg = d3.select(svgRef.current)
      .attr("width", width)
      .attr("height", height);
    
    // Tambahkan transformasi dengan zoom
    const g = svg.append("g")
      .attr("transform", `translate(100, 50) scale(${zoom})`);
    
    // Definisikan perilaku zoom
    const zoomer = d3.zoom()
      .scaleExtent([0.5, 3])
      .on("zoom", (event) => {
        g.attr("transform", event.transform);
      });
    
    // Terapkan zoom ke SVG
    svg.call(zoomer);
    
    // Gambar garis antar node (link)
    const links = g.selectAll(".link")
      .data(root.links())
      .enter()
      .append("path")
      .attr("class", "link")
      .attr("d", d3.linkHorizontal()
        .x((d) => d.y)
        .y((d) => d.x))
      .attr("stroke", (d) => getNodeColor(d.source.depth))
      .attr("fill", "none")
      .attr("stroke-width", 2);
    
    // Gambar node (elemen) beserta icon dan nama
    const nodes = g.selectAll(".node")
      .data(root.descendants())
      .enter()
      .append("g")
      .attr("class", "node")
      .attr("transform", (d) => `translate(${d.y}, ${d.x})`);
    
    // Tentukan ukuran node berdasarkan kedalaman
    const nodeSize = (depth) => {
      return Math.max(50 - (depth * 5), 30); // Ukuran menurun dengan kedalaman, minimal 30px
    };
    
    // Tambahkan lingkaran sebagai latar belakang node
    g.selectAll(".node")
      .append("circle")
      .attr("r", (d) => nodeSize(d.depth) / 2)
      .attr("fill", (d) => getNodeColor(d.depth))
      .attr("stroke", "#333")
      .attr("stroke-width", 1);
    
    // Gambar icon elemen (atau placeholder jika tidak ada)
    g.selectAll(".node")
      .append("image") 
      .attr("xlink:href", (d) => d.data.image || "/placeholder.svg")
      .attr("width", (d) => nodeSize(d.depth))
      .attr("height", (d) => nodeSize(d.depth))
      .attr("x", (d) => -nodeSize(d.depth) / 2)
      .attr("y", (d) => -nodeSize(d.depth) / 2);
    
    // Tampilkan nama elemen di bawah icon
    g.selectAll(".node")
      .append("text")
      .attr("dy", (d) => nodeSize(d.depth) / 2 + 15)
      .attr("text-anchor", "middle")
      .text((d) => d.data.name)
      .style("font-size", (d) => Math.max(13 - d.depth, 9) + "px")
      .style("font-weight", (d) => d.depth === 0 ? "bold" : "normal")
      .style("fill", "#333");
    
    // Tambahkan informasi resep jika ada
    g.selectAll(".node")
      .filter((d) => d.data.recipe && Array.isArray(d.data.recipe))
      .append("text")
      .attr("dy", (d) => nodeSize(d.depth) / 2 + 30)
      .attr("text-anchor", "middle")
      .text((d) => d.data.recipe[0] || "")
      .style("font-size", "10px")
      .style("font-style", "italic")
      .style("fill", "#666");
  };
  
  // Tampilkan informasi pencarian yang sudah selesai
  const renderSearchInfo = () => {
    if (!results || results.length === 0) return null;
    
    return (
      <div className="search-info" style={{ 
        position: "absolute", 
        top: "10px", 
        left: "10px", 
        background: "rgba(255,255,255,0.8)",
        padding: "10px",
        borderRadius: "5px",
        boxShadow: "0 2px 5px rgba(0,0,0,0.2)",
        zIndex: 10
      }}>
        <h3>Search Results</h3>
        <p><strong>Algorithm:</strong> {searchAlgorithm}</p>
        <p><strong>Nodes visited:</strong> {nodesVisited}</p>
        <p><strong>Recipes found:</strong> {results.length}</p>
      </div>
    );
  };
  
  // Tampilkan spinner selama loading
  const spinnerStyle = `
    @keyframes spin {
      0% { transform: rotate(0deg); }
      100% { transform: rotate(360deg); }
    }
  `;
  
  // Tampilkan pesan jika tidak ada hasil
  if (!isLoading && (!results || results.length === 0)) {
    return (
      <div className="tree-empty" style={{ 
        display: "flex", 
        justifyContent: "center", 
        alignItems: "center", 
        height: "600px",
        flexDirection: "column",
        background: "#f9f9f9",
        borderRadius: "8px"
      }}>
        <svg width="100" height="100" viewBox="0 0 24 24" fill="none" stroke="#999" strokeWidth="1.5">
          <path d="M9.879 9.879A3 3 0 1 0 12 15c-.883 0-1.68-.377-2.227-.984m0 0a2.94 2.94 0 0 1-.237-.368m0 0c.25-.65.437-1.363.548-2.121.142-.972.228-1.977.228-3.001 0-2.304-.549-4.195-1.228-5.527m0 0a2.98 2.98 0 0 0-.227-.368M6.879 6.879l-5 5m10.242 10.242l-5-5" />
        </svg>
        <p style={{ margin: "20px 0", fontSize: "16px", color: "#666" }}>
          No recipe found for the selected element.
        </p>
        <p style={{ fontSize: "14px", color: "#999" }}>
          Try selecting another element or using a different search algorithm.
        </p>
      </div>
    );
  }
  
  // Tampilkan loading state jika sedang loading
  if (isLoading) {
    return (
      <div className="tree-loading" style={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        height: "600px",
        flexDirection: "column",
        borderRadius: "8px"
      }}>
        <style>{spinnerStyle}</style>
        <div className="spinner" style={{
          width: "40px",
          height: "40px",
          border: "5px solid #f3f3f3",
          borderTop: "5px solid #3498db",
          borderRadius: "50%",
          animation: "spin 1s linear infinite"
        }}></div>
        <p style={{ margin: "20px 0", fontSize: "16px" }}>
          Searching recipes for {selectedElement?.name}...
        </p>
      </div>
    );
  }
  
  // Tampilkan visualisasi pohon, kontrol zoom, dan info pencarian
  return (
    <div className="tree-container" ref={containerRef} style={{ 
      position: "relative", 
      width: "100%", 
      height: "600px",
      background: "#f9f9f9",
      borderRadius: "8px",
      overflow: "hidden"
    }}>
      <style>{spinnerStyle}</style>
      <svg ref={svgRef} style={{ width: "100%", height: "100%" }}></svg>
      
      {renderSearchInfo()}
      
      <div className="recipe-info" style={{
        position: "absolute",
        bottom: "60px",
        right: "10px",
        background: "rgba(255,255,255,0.9)",
        padding: "10px",
        borderRadius: "5px",
        boxShadow: "0 2px 5px rgba(0,0,0,0.2)",
        maxWidth: "300px",
        maxHeight: "200px",
        overflowY: "auto",
        display: selectedElement && results && results.length > 0 ? "block" : "none"
      }}>
        <h3>Recipe for {selectedElement?.name}</h3>
        {results && results[0]?.recipe && Array.isArray(results[0].recipe) ? (
          <ul style={{ padding: "0 0 0 20px", margin: "5px 0" }}>
            {results[0].recipe.map((step, i) => (
              <li key={i} style={{ margin: "5px 0" }}>{step}</li>
            ))}
          </ul>
        ) : (
          <p>{results && results[0]?.recipe || "No recipe details available"}</p>
        )}
      </div>
      
      <div className="zoom-controls" style={{ 
        position: "absolute", 
        bottom: "10px", 
        right: "10px",
        background: "white",
        borderRadius: "5px",
        boxShadow: "0 2px 5px rgba(0,0,0,0.2)",
        padding: "5px"
      }}>
        <button 
          className="zoom-button" 
          onClick={() => setZoom((prev) => Math.min(prev + 0.2, 3))}
          style={{
            border: "none",
            background: "#3498db",
            color: "white",
            width: "30px",
            height: "30px",
            borderRadius: "5px",
            margin: "0 5px",
            cursor: "pointer",
            fontSize: "16px"
          }}
        >+</button>
        <button 
          className="zoom-button" 
          onClick={() => setZoom(1)}
          style={{
            border: "none",
            background: "#2ecc71",
            color: "white",
            width: "30px",
            height: "30px",
            borderRadius: "5px",
            margin: "0 5px",
            cursor: "pointer",
            fontSize: "16px"
          }}
        >⟳</button>
        <button 
          className="zoom-button" 
          onClick={() => setZoom((prev) => Math.max(prev - 0.2, 0.5))}
          style={{
            border: "none",
            background: "#3498db",
            color: "white",
            width: "30px",
            height: "30px",
            borderRadius: "5px",
            margin: "0 5px",
            cursor: "pointer",
            fontSize: "16px"
          }}
        >-</button>
      </div>
    </div>
  );
}