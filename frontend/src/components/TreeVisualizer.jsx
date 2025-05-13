// TreeVisualizer.jsx with improved responsiveness
import React from "react";
import { useEffect, useRef, useState } from "react";
import * as d3 from "d3";

// Helper function to ensure all nodes have a children property
function processNode(node) {
  if (!node) return { name: "Unknown", children: [] };
  
  const newNode = { ...node };
  
  // Ensure children exists and is an array
  if (!newNode.children) {
    newNode.children = [];
  }
  
  // Recursively process children
  if (Array.isArray(newNode.children)) {
    newNode.children = newNode.children.map(processNode);
  }
  
  return newNode;
}

export default function TreeVisualizer({ 
  results, 
  selectedElement, 
  isLoading, 
  nodesVisited,  // Jumlah node yang telah dikunjungi
  searchAlgorithm, // BFS, DFS, atau Bidirectional
  searchParams  // Parameter pencarian (algorithm, recipeType, maxRecipes)
}) {
  const svgRef = useRef(null); // Referensi ke elemen SVG untuk D3
  const containerRef = useRef(null); // Referensi ke container untuk mengambil ukuran
  const [zoom, setZoom] = useState(1); // State untuk level zoom
  const [currentRecipeIndex, setCurrentRecipeIndex] = useState(0); // Current recipe to display
  
  // Reset current recipe index when results change
  useEffect(() => {
    setCurrentRecipeIndex(0);
  }, [results]);
  
  // Handle window resize for responsiveness
  useEffect(() => {
    const handleResize = () => {
      renderTree();
    };
    
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);
  
  // Render ulang tree setiap kali hasil, loading, zoom, atau recipe index berubah
  useEffect(() => {
    renderTree();
  }, [results, isLoading, zoom, selectedElement, currentRecipeIndex]);
  
  // Fungsi untuk menghasilkan warna berdasarkan kedalaman node
  const getNodeColor = (depth) => {
    const colors = ["#4CAF50", "#2196F3", "#FFC107", "#E91E63", "#9C27B0", "#FF5722"];
    return colors[depth % colors.length];
  };
  
  // Pagination controls
  const handlePrevRecipe = () => {
    setCurrentRecipeIndex(prev => Math.max(prev - 1, 0));
  };
  
  const handleNextRecipe = () => {
    setCurrentRecipeIndex(prev => Math.min(prev + 1, results.length - 1));
  };
  
  // Zoom controls
  const handleZoomIn = () => {
    setZoom(prev => Math.min(prev + 0.2, 3));
  };
  
  const handleZoomOut = () => {
    setZoom(prev => Math.max(prev - 0.2, 0.3));
  };
  
  const handleZoomReset = () => {
    setZoom(1);
  };
  
  // Fungsi utama untuk menggambar tree dengan D3
  const renderTree = () => {
    if (!svgRef.current || !containerRef.current || !results || results.length === 0) return;
    
    d3.select(svgRef.current).selectAll("*").remove(); // Bersihkan SVG sebelum render ulang
    
    // Get container dimensions (responsive)
    const width = containerRef.current.clientWidth; 
    const height = containerRef.current.clientHeight;
    
    // Get the current recipe to display
    const currentRecipe = results[currentRecipeIndex];
    if (!currentRecipe) return;
    
    // Process the current recipe
    const processedRecipe = processNode(currentRecipe);
    
    // Layout pohon dengan ukuran yang disesuaikan berdasarkan lebar container
    const treeLayout = d3.tree().size([height * 0.85, width * 0.85]);
    const root = d3.hierarchy(processedRecipe);
    
    // Tambahkan depth sebagai properti untuk pewarnaan
    root.descendants().forEach((d, i) => {
      d.id = i;
      d.depth = d.depth;
    });
    
    treeLayout(root);
    
    // Inisialisasi SVG dan group utama
    const svg = d3.select(svgRef.current)
      .attr("width", width)
      .attr("height", height)
      .attr("viewBox", `0 0 ${width} ${height}`)
      .attr("preserveAspectRatio", "xMidYMid meet");
    
    // Tambahkan transformasi dengan zoom
    const g = svg.append("g")
      .attr("transform", `translate(${width * 0.15}, ${height * 0.5}) scale(${zoom})`);
    
    // Definisikan perilaku zoom
    const zoomer = d3.zoom()
      .scaleExtent([0.1, 5]) // Extend zoom range
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
    
    // Tentukan ukuran node dengan mempertimbangkan kedalaman yang lebih dalam
    const nodeSize = (depth) => {
      // Make node size responsive
      const baseSizeMultiplier = Math.min(width, height) / 600;
      const baseSize = Math.max(40 * baseSizeMultiplier, 15);
      return Math.max(baseSize - (depth * 3), 12);
    };
    
    // Gambar node (elemen) beserta icon dan nama
    const nodes = g.selectAll(".node")
      .data(root.descendants())
      .enter()
      .append("g")
      .attr("class", "node")
      .attr("transform", (d) => `translate(${d.y}, ${d.x})`);
    
    // Tambahkan lingkaran sebagai latar belakang node
    nodes.append("circle")
      .attr("r", (d) => nodeSize(d.depth) / 2)
      .attr("fill", (d) => getNodeColor(d.depth))
      .attr("stroke", "#333")
      .attr("stroke-width", 1);
    
    // Gambar icon elemen (atau placeholder jika tidak ada)
    nodes.append("image") 
      .attr("xlink:href", (d) => d.data.image || "/placeholder.svg")
      .attr("width", (d) => nodeSize(d.depth))
      .attr("height", (d) => nodeSize(d.depth))
      .attr("x", (d) => -nodeSize(d.depth) / 2)
      .attr("y", (d) => -nodeSize(d.depth) / 2);
    
    // Tampilkan nama elemen di bawah icon dengan ukuran responsif
    nodes.append("text")
      .attr("dy", (d) => nodeSize(d.depth) / 2 + 10)
      .attr("text-anchor", "middle")
      .text((d) => d.data.name)
      .style("font-size", (d) => {
        const baseSize = Math.min(width, height) / 600;
        return Math.max(10 - (d.depth * 0.5), 6) * baseSize + "px";
      })
      .style("font-weight", (d) => d.depth === 0 ? "bold" : "normal")
      .style("fill", "#333");
    
    // Display recipe steps if available (responsive positioning)
    if (currentRecipe.recipe && Array.isArray(currentRecipe.recipe)) {
      const recipePanel = svg.append("g")
        .attr("transform", width < 480 ? 
          `translate(${width * 0.05}, ${height * 0.05})` : 
          `translate(20, 30)`);
      
      // Background for recipe steps
      recipePanel.append("rect")
        .attr("width", width < 480 ? width * 0.9 : 250)
        .attr("height", currentRecipe.recipe.length * (width < 480 ? 18 : 22) + 40)
        .attr("fill", "rgba(255, 255, 255, 0.9)")
        .attr("rx", 5)
        .attr("ry", 5)
        .attr("stroke", "#ddd");
      
      // Recipe steps title
      recipePanel.append("text")
        .attr("x", 10)
        .attr("y", 25)
        .text("Recipe Steps:")
        .style("font-weight", "bold")
        .style("font-size", width < 480 ? "12px" : "14px");
      
      // Add each recipe step
      currentRecipe.recipe.forEach((step, i) => {
        // If on small screens, truncate long recipe steps
        let displayText = step;
        if (width < 480 && step.length > 25) {
          displayText = step.substring(0, 25) + "...";
        }
        
        recipePanel.append("text")
          .attr("x", 15)
          .attr("y", 45 + i * (width < 480 ? 18 : 20))
          .text(displayText)
          .style("font-size", width < 480 ? "10px" : "12px");
      });
    }
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
        height: "100%",
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
        height: "100%",
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
  
  // Tampilkan visualisasi pohon, kontrol zoom, dan navigasi antar recipe
  return (
    <div className="tree-container" ref={containerRef} style={{ 
      position: "relative", 
      width: "100%", 
      height: "100%",
      background: "#f9f9f9",
      borderRadius: "8px",
      overflow: "hidden"
    }}>
      <svg ref={svgRef} style={{ width: "100%", height: "100%" }}></svg>
      
      {/* Recipe information and stats */}
      <div className="recipe-info" style={{
        position: "absolute",
        top: "10px",
        right: "10px",
        background: "rgba(255, 255, 255, 0.9)",
        padding: "8px",
        borderRadius: "8px",
        boxShadow: "0 2px 10px rgba(0, 0, 0, 0.1)",
        fontSize: "12px",
        maxWidth: "40%",
        overflowX: "hidden"
      }}>
        <p style={{ margin: "0 0 3px 0", whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis" }}>
          <strong>Algorithm:</strong> {searchAlgorithm}
        </p>
        <p style={{ margin: "0 0 3px 0" }}>
          <strong>Nodes:</strong> {nodesVisited}
        </p>
        <p style={{ margin: "0" }}>
          <strong>Recipe {currentRecipeIndex + 1}/{results.length}</strong>
        </p>
      </div>
      
      {/* Zoom controls */}
      <div className="zoom-controls" style={{
        position: "absolute",
        bottom: "10px",
        left: "10px",
        display: "flex",
        flexDirection: "column",
        gap: "5px",
        background: "rgba(30, 41, 59, 0.8)",
        borderRadius: "6px",
        padding: "5px",
        zIndex: 10
      }}>
        <button 
          onClick={handleZoomIn}
          style={{
            width: "30px",
            height: "30px",
            background: "rgba(79, 70, 229, 0.8)",
            border: "none",
            borderRadius: "4px",
            color: "white",
            fontSize: "16px",
            cursor: "pointer"
          }}
        >
          +
        </button>
        <button 
          onClick={handleZoomReset}
          style={{
            width: "30px",
            height: "30px",
            background: "rgba(79, 70, 229, 0.8)",
            border: "none",
            borderRadius: "4px",
            color: "white",
            fontSize: "12px",
            cursor: "pointer"
          }}
        >
          R
        </button>
        <button 
          onClick={handleZoomOut}
          style={{
            width: "30px",
            height: "30px",
            background: "rgba(79, 70, 229, 0.8)",
            border: "none",
            borderRadius: "4px",
            color: "white",
            fontSize: "16px",
            cursor: "pointer"
          }}
        >
          -
        </button>
      </div>
      
      {/* Add recipe navigation controls if there are multiple recipes */}
      {results.length > 1 && (
        <div className="recipe-navigation" style={{
          position: "absolute",
          bottom: "10px",
          left: "50%",
          transform: "translateX(-50%)",
          display: "flex",
          alignItems: "center",
          background: "rgba(30, 41, 59, 0.8)",
          padding: "8px 12px",
          borderRadius: "8px",
          zIndex: 10,
          maxWidth: "90%"
        }}>
          <button 
            onClick={handlePrevRecipe} 
            disabled={currentRecipeIndex === 0}
            style={{
              background: currentRecipeIndex === 0 ? "#475569" : "#4f46e5",
              color: "white",
              border: "none",
              borderRadius: "4px",
              padding: "6px 12px",
              marginRight: "8px",
              cursor: currentRecipeIndex === 0 ? "default" : "pointer",
              fontSize: "14px"
            }}
          >
            Previous
          </button>
          <span style={{ 
            margin: "0 8px", 
            color: "white", 
            fontSize: "14px",
            whiteSpace: "nowrap"
          }}>
            {currentRecipeIndex + 1} of {results.length}
          </span>
          <button 
            onClick={handleNextRecipe} 
            disabled={currentRecipeIndex === results.length - 1}
            style={{
              background: currentRecipeIndex === results.length - 1 ? "#475569" : "#4f46e5",
              color: "white",
              border: "none",
              borderRadius: "4px",
              padding: "6px 12px",
              marginLeft: "8px",
              cursor: currentRecipeIndex === results.length - 1 ? "default" : "pointer",
              fontSize: "14px"
            }}
          >
            Next
          </button>
        </div>
      )}
    </div>
  );
}