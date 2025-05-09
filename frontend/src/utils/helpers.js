export const fetchElements = async () => {
  try {
      const response = await fetch("http://localhost:5000/elements");
      if (response.ok) {
          return await response.json();
      }
      throw new Error("Failed to fetch elements");
  } catch (error) {
      console.error("Error fetching elements:", error);
      return [];
  }
};

export const formatTime = (ms) => {
  if (ms < 1000) {
      return `${ms}ms`;
  }
  const seconds = Math.floor(ms / 1000);
  const remainingMs = ms % 1000;
  if (seconds < 60) {
      return `${seconds}.${remainingMs.toString().padStart(3, "0")}s`;
  }
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;
  return `${minutes}m ${remainingSeconds}s`;
};

export const formatNumber = (num) => {
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
};