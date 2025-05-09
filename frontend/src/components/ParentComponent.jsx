import React, { useState } from "react";
import SearchBar from "./SearchBar";
import ElementPicker from "./ElementPicker";

export default function ParentComponent() {
  const [searchResults, setSearchResults] = useState([]); // Hasil pencarian
  const [selectedElement, setSelectedElement] = useState(null); // Elemen yang dipilih

  const handleSearch = (results) => {
    setSearchResults(results || []); // Pastikan hasil pencarian adalah array
  };

  const handleElementSelect = (element) => {
    setSelectedElement(element); // Perbarui elemen yang dipilih
  };

  return (
    <div>
      <SearchBar onSearch={handleSearch} />
      <ElementPicker
        searchResults={searchResults}
        onElementSelect={handleElementSelect}
        selectedElement={selectedElement}
      />
    </div>
  );
}