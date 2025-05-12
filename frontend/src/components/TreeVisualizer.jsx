// src/components/TreeVisualizer.jsx
// Komponen visualisasi pohon resep menggunakan D3.js.
// Menampilkan hasil pencarian dalam bentuk tree langsung tanpa animasi.

"use client";

import React from "react";
import { useEffect, useRef, useState } from "react";
import * as d3 from "d3";

// Helper function to ensure all nodes have a children property and recursively process them
function processNode(node) {
  if (!node) return { name: "Unknown", children: [] };
  
  const newNode = { ...node };
  
  // Ensure children exists and is an array
  if (!newNode.children) {
    newNode.children = [];
  }
  
  // Recursively process children, even if it's an empty array
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
    
    // Process all results recursively
    const processedResults = results.map(processNode);
    
    // Data root tree langsung dari hasil
    const rootData = processNode({
      name: selectedElement?.name || "No Element Selected",
      image: selectedElement?.icon || "",
      children: processedResults
    });
    
    // Layout pohon dengan ukuran yang lebih besar untuk mengakomodasi pohon yang dalam
    const treeLayout = d3.tree().size([height * 2, width * 2]);
    const root = d3.hierarchy(rootData);
    
    // Tambahkan depth sebagai properti untuk pewarnaan
    root.descendants().forEach((d, i) => {
      d.id = i;
      d.depth = d.depth;
    });
    
    treeLayout(root);
    
    // Inisialisasi SVG dan group utama
    const svg = d3.select(svgRef.current)
      .attr("width", width * 2)
      .attr("height", height * 2);
    
    // Tambahkan transformasi dengan zoom
    const g = svg.append("g")
      .attr("transform", `translate(100, 50) scale(${zoom})`);
    
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
      return Math.max(40 - (depth * 3), 20); // Ukuran menurun dengan kedalaman, minimal 20px
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
    
    // Tampilkan nama elemen di bawah icon
    nodes.append("text")
      .attr("dy", (d) => nodeSize(d.depth) / 2 + 10)
      .attr("text-anchor", "middle")
      .text((d) => d.data.name)
      .style("font-size", (d) => Math.max(10 - (d.depth * 0.5), 6) + "px")
      .style("font-weight", (d) => d.depth === 0 ? "bold" : "normal")
      .style("fill", "#333");
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
      <svg ref={svgRef} style={{ width: "100%", height: "100%" }}></svg>
    </div>
  );
}