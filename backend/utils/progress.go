// Package utils provides utility functions for the application
package utils

import (
	"sync"
)

// ProgressData represents the progress information of a search algorithm
type ProgressData struct {
  NodesVisited int     `json:"nodesVisited"`
  Progress     float64 `json:"progress"`
  Completed    bool    `json:"completed"`
  CurrentNode  string  `json:"currentNode"`
}

var (
  searchProgress ProgressData
  progressMutex  sync.RWMutex
)

// UpdateProgress updates the search progress information
func UpdateProgress(nodesVisited int, progress float64, currentNode string, completed bool) {
  progressMutex.Lock()
  defer progressMutex.Unlock()
  
  searchProgress = ProgressData{
    NodesVisited: nodesVisited,
    Progress:     progress,
    Completed:    completed,
    CurrentNode:  currentNode,
  }
}

// GetProgress returns the current search progress
func GetProgress() ProgressData {
  progressMutex.RLock()
  defer progressMutex.RUnlock()
  
  return searchProgress
}