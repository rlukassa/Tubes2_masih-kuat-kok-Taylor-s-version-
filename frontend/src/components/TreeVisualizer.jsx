// src/components/TreeVisualizer.jsx
// Komponen visualisasi pohon resep menggunakan D3.js.
// Menampilkan hasil pencarian dalam bentuk tree, lengkap dengan gambar/icon dan nama elemen.
// Mendukung zoom in/out dan reset zoom, serta animasi loading saat pencarian berlangsung.

"use client";

import React from "react";
import { useEffect, useRef } from "react";
import * as d3 from "d3";

export default function TreeVisualizer({ results, selectedElement, isLoading }) {
  const svgRef = useRef(null); // Referensi ke elemen SVG untuk D3
  const containerRef = useRef(null); // Referensi ke container untuk mengambil ukuran
  const [zoom, setZoom] = React.useState(1); // State untuk level zoom

  // Render ulang tree setiap kali hasil, loading, atau zoom berubah
  useEffect(() => {
    if (!results || results.length === 0 || isLoading) return;
    renderTree();
  }, [results, isLoading, zoom]);

  // Fungsi utama untuk menggambar tree dengan D3
  const renderTree = () => {
    if (!svgRef.current || !containerRef.current) return;

    d3.select(svgRef.current).selectAll("*").remove(); // Bersihkan SVG sebelum render ulang

    const width = containerRef.current.clientWidth; // Lebar container
    const height = containerRef.current.clientHeight; // Tinggi container

    // Data root tree, children diisi dari hasil pencarian
    const treeData = {
      name: selectedElement?.name || "No Element Selected",
      icon: selectedElement?.icon || "",
      children: results || [],
    };

    // Layout pohon dengan ukuran tertentu
    const treeLayout = d3.tree().size([height - 100, width - 160]);
    const root = d3.hierarchy(treeData);
    treeLayout(root);

    // Inisialisasi SVG dan group utama
    const svg = d3.select(svgRef.current)
      .attr("width", width)
      .attr("height", height)
      .append("g")
      .attr("transform", `translate(80, 50) scale(${zoom})`); // Posisikan dan scale sesuai zoom

    // Gambar garis antar node (link)
    svg.selectAll(".link")
      .data(root.links())
      .enter()
      .append("path")
      .attr("class", "link")
      .attr("d", d3.linkHorizontal()
        .x((d) => d.y)
        .y((d) => d.x))
      .attr("stroke", "red")
      .attr("fill", "none")
      .attr("stroke-width", 2);

    // Gambar node (elemen) beserta icon dan nama
    const nodes = svg.selectAll(".node")
      .data(root.descendants())
      .enter()
      .append("g")
      .attr("class", "node")
      .attr("transform", (d) => `translate(${d.y}, ${d.x})`);

    // Gambar icon elemen (atau placeholder jika tidak ada)
    nodes.append("image") 
      .attr("xlink:href", (d) => d.data.icon || "/placeholder.svg")
      .attr("width", 60)
      .attr("height", 60)
      .attr("x", -30)
      .attr("y", -30);

    // Tampilkan nama elemen di bawah icon
    nodes.append("text")
      .attr("dy", "4em")
      .attr("text-anchor", "middle")
      .text((d) => d.data.name)
      .style("font-size", "14px")
      .style("fill", "#333");
  };

  // Tampilkan animasi loading jika sedang mencari resep (INI KAYANYA DIUBAH TAMPILANNYA JADI LIVE UPDATE TREE dengan DELAY, biar terlihat response proses searching)
  if (isLoading) {
    return (
      <div className="tree-loading">
        <div className="spinner"></div>
        <p>Searching for recipes...</p>
      </div>
    );
  }

  // Tampilkan pesan jika tidak ada hasil
  if (!results || results.length === 0) {
    return (
      <div className="tree-empty">
        <p>No recipe found for the selected element.</p>
      </div>
    );
  }

  // Tampilkan visualisasi pohon dan kontrol zoom
  return (
    <div className="tree-container" ref={containerRef} style={{ position: "relative", width: "100%", height: "600px" }}>
      <svg ref={svgRef}></svg>
      <div className="zoom-controls" style={{ position: "absolute", bottom: "10px", right: "10px" }}>
        <button className="zoom-button" onClick={() => setZoom((prev) => Math.min(prev + 0.1, 2))}>+</button>
        <button className="zoom-button" onClick={() => setZoom(1)}>⟳</button>
        <button className="zoom-button" onClick={() => setZoom((prev) => Math.max(prev - 0.1, 0.5))}>-</button>
      </div>
    </div>
  );
}