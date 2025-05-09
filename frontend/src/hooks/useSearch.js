"use client"

import { useState, useEffect } from "react";

export function useSearch(initialItems = []) {
    const [searchParams, setSearchParams] = useState({
        algorithm: "BFS",
        recipeType: "Best",
        maxRecipes: 5,
    });

    const [searchResults, setSearchResults] = useState([]);
    const [isLoading, setIsLoading] = useState(false);
    const [executionTime, setExecutionTime] = useState(0);
    const [nodesVisited, setNodesVisited] = useState(0);
    const [progress, setProgress] = useState(0);

    const [searchTerm, setSearchTerm] = useState("");
    const [items, setItems] = useState(initialItems);
    const [filteredItems, setFilteredItems] = useState(initialItems);

    useEffect(() => {
        setItems(initialItems);
        setFilteredItems(initialItems);
    }, [initialItems]);

    useEffect(() => {
        if (!searchTerm.trim()) {
            setFilteredItems(items);
            return;
        }
        const searchTermLower = searchTerm.toLowerCase();
        const filtered = items.filter(item =>
            item && (item.name.toLowerCase().includes(searchTermLower) || 
                    (item.description && item.description.toLowerCase().includes(searchTermLower)))
        );
        setFilteredItems(filtered);
    }, [searchTerm, items]);

    const startSearch = async (element) => {
        setIsLoading(true);
        setSearchResults([]);
        setExecutionTime(0);
        setNodesVisited(0);
        setProgress(0);

        try {
            const response = await fetch("http://localhost:5000/search", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    elementId: element.id,
                    algorithm: searchParams.algorithm,
                    recipeType: searchParams.recipeType,
                    maxRecipes: searchParams.maxRecipes,
                }),
            });

            const data = await response.json();
            setSearchResults(data.results);
            setNodesVisited(data.nodesVisited);
            setExecutionTime(data.executionTime);
            setProgress(100);
        } catch (error) {
            console.error("Search error:", error);
        } finally {
            setIsLoading(false);
        }
    };

    const resetSearch = () => {
        setSearchResults([]);
        setExecutionTime(0);
        setNodesVisited(0);
        setProgress(0);
    };

    return {
        searchParams,
        setSearchParams,
        searchResults,
        isLoading,
        executionTime,
        nodesVisited,
        progress,
        startSearch,
        resetSearch,
        searchTerm,
        setSearchTerm,
        items,
        setItems,
        filteredItems,
        hasResults: filteredItems.length > 0,
        isSearching: searchTerm.trim() !== "",
    };
}