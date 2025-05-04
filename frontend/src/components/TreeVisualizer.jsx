"use client"

import { useEffect, useRef, useState } from "react"
import * as d3 from "d3"

export default function TreeVisualizer({ results, selectedElement, isLoading }) {
  const svgRef = useRef(null)
  const containerRef = useRef(null)
  const [zoom, setZoom] = useState(1)

  useEffect(() => {
    if (!results || results.length === 0 || isLoading) return

    renderTree()
  }, [results, isLoading, zoom])

  const handleZoomIn = () => {
    setZoom((prev) => Math.min(prev + 0.1, 2))
  }

  const handleZoomOut = () => {
    setZoom((prev) => Math.max(prev - 0.1, 0.5))
  }

  const handleResetZoom = () => {
    setZoom(1)
  }

  const renderTree = () => {
    if (!svgRef.current || !containerRef.current) return

    // Clear previous visualization
    d3.select(svgRef.current).selectAll("*").remove()

    const width = containerRef.current.clientWidth
    const height = containerRef.current.clientHeight

    // (DUMMY DATA STRUCTURE)
    const treeData = {
      name: selectedElement?.name || "Element",
      children: [
        {
          name: "Mud",
          children: [
            { name: "Water", children: [] },
            { name: "Earth", children: [] },
          ],
        },
        {
          name: "Stone",
          children: [
            {
              name: "Clay",
              children: [
                { name: "Mud", children: [] },
                { name: "Sand", children: [] },
              ],
            },
          ],
        },
      ],
    }

    // Tree layout
    const treeLayout = d3.tree().size([height - 100, width - 160])

    // Hierarchy for tree
    const root = d3.hierarchy(treeData)

    // Assign x and y coordinates to each node
    treeLayout(root)

    const svg = d3
      .select(svgRef.current)
      .attr("width", width)
      .attr("height", height)
      .append("g")
      .attr("transform", `translate(80, 50) scale(${zoom})`)

    // Create links
    svg
      .selectAll(".link")
      .data(root.links())
      .enter()
      .append("path")
      .attr("class", "link")
      .attr(
        "d",
        d3
          .linkHorizontal()
          .x((d) => d.y)
          .y((d) => d.x),
      )
      .attr("fill", "none")
      .attr("stroke", "#6366f1")
      .attr("stroke-width", 2)

    // Create nodes
    const nodes = svg
      .selectAll(".node")
      .data(root.descendants())
      .enter()
      .append("g")
      .attr("class", "node")
      .attr("transform", (d) => `translate(${d.y}, ${d.x})`)

    nodes
      .append("rect")
      .attr("x", -50)
      .attr("y", -15)
      .attr("width", 100)
      .attr("height", 30)
      .attr("rx", 5)
      .attr("ry", 5)
      .attr("fill", (d) => (d.depth === 0 ? "#4f46e5" : "#3b82f6"))
      .attr("stroke", "#1e3a8a")
      .attr("stroke-width", 1)

    // Add node text
    nodes
      .append("text")
      .attr("dy", "0.35em")
      .attr("text-anchor", "middle")
      .attr("fill", "white")
      .text((d) => d.data.name)
  }

  if (isLoading) {
    return (
      <div className="tree-loading">
        <div className="spinner"></div>
        <p>Searching for recipes...</p>
      </div>
    )
  }

  if (!results || results.length === 0) {
    return (
      <div className="tree-empty">
        <p>No recipe found for the selected element.</p>
      </div>
    )
  }

  return (
    <div className="tree-container" ref={containerRef}>
      <svg ref={svgRef}></svg>

      <div className="zoom-controls">
        <button className="zoom-button" onClick={handleZoomIn}>
          +
        </button>
        <button className="zoom-button" onClick={handleResetZoom}>
          ⟳
        </button>
        <button className="zoom-button" onClick={handleZoomOut}>
          -
        </button>
      </div>
    </div>
  )
}
