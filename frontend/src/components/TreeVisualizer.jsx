"use client";

import React from "react";
import { useEffect, useRef } from "react";
import * as d3 from "d3";

export default function TreeVisualizer({ results, selectedElement, isLoading }) {
  const svgRef = useRef(null);
  const containerRef = useRef(null);
  const [zoom, setZoom] = React.useState(1);

  useEffect(() => {
    if (!results || results.length === 0 || isLoading) return;
    renderTree();
  }, [results, isLoading, zoom]);

  const renderTree = () => {
    if (!svgRef.current || !containerRef.current) return;
  
    d3.select(svgRef.current).selectAll("*").remove();
  
    const width = containerRef.current.clientWidth;
    const height = containerRef.current.clientHeight;
  
    // Gunakan elemen yang dipilih untuk membuat data pohon
    const treeData = {
      name: selectedElement?.name || "No Element Selected",
      icon: selectedElement?.icon || "", // Ambil ikon dari selectedElement
      children: [], // Tidak ada anak untuk elemen tunggal
    };
  
    const treeLayout = d3.tree().size([height - 100, width - 160]);
    const root = d3.hierarchy(treeData);
    treeLayout(root);
  
    const svg = d3.select(svgRef.current)
      .attr("width", width)
      .attr("height", height)
      .append("g")
      .attr("transform", `translate(80, 50) scale(${zoom})`);
  
    const nodes = svg.selectAll(".node")
      .data(root.descendants())
      .enter()
      .append("g")
      .attr("class", "node")
      .attr("transform", (d) => `translate(${d.y}, ${d.x})`);
  
    // Tambahkan gambar elemen
    nodes.append("image")
      .attr("xlink:href", (d) => d.data.icon || "") // Gunakan URL ikon
      .attr("width", 60)
      .attr("height", 60)
      .attr("x", -30) // Pusatkan gambar
      .attr("y", -30); // Pusatkan gambar
  
    // Tambahkan nama elemen
    nodes.append("text")
      .attr("dy", "4em")
      .attr("text-anchor", "middle")
      .text((d) => d.data.name);
  };

  if (isLoading) {
    return (
      <div className="tree-loading">
        <div className="spinner"></div>
        <p>Searching for recipes...</p>
      </div>
    );
  }

  if (!results || results.length === 0) {
    return (
      <div className="tree-empty">
        <p>No recipe found for the selected element.</p>
      </div>
    );
  }

  return (
    <div className="tree-container" ref={containerRef}>
      <svg ref={svgRef}></svg>
      <div className="zoom-controls">
        <button className="zoom-button" onClick={() => setZoom((prev) => Math.min(prev + 0.1, 2))}>+</button>
        <button className="zoom-button" onClick={() => setZoom(1)}>⟳</button>
        <button className="zoom-button" onClick={() => setZoom((prev) => Math.max(prev - 0.1, 0.5))}>-</button>
      </div>
    </div>
  );
}