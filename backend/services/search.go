// Modified search.go implementation with early stopping
package services

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var mapper map[string]string

func init() {
	var err error
	db, err = sql.Open("sqlite3", "../database/alchemy.db")
	if err != nil {
		log.Fatalf("Gagal membuka database: %v", err)
	} else {
		log.Printf("Database ditemukan")
	}

	file, err := os.Open("../database/mapper2.json")
	if err != nil {
		log.Fatalf("Gagal membuka mapper.json: %v", err)
	}
	defer file.Close()
	if err := json.NewDecoder(file).Decode(&mapper); err != nil {
		log.Fatalf("Gagal mendekode mapper.json: %v", err)
	}
}

type Node struct {
	Name     string
	Children []*Node
}

type RecipeStep struct {
	Result string
	Item1  string
	Item2  string
}

// Get all basic elements (Water, Fire, Earth, Air, etc.)
func getBasicElements() []string {
	basicElements := []string{}
	rows, err := db.Query("SELECT DISTINCT element FROM elements WHERE element NOT IN (SELECT DISTINCT element FROM elements WHERE item1 IS NOT NULL AND item2 IS NOT NULL)")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var element string
			if err := rows.Scan(&element); err == nil {
				basicElements = append(basicElements, element)
			}
		}
	}
	return basicElements
}

// Check if an element is a basic element
func isBasicElement(element string, basicElements []string) bool {
	for _, basic := range basicElements {
		if element == basic {
			return true
		}
	}
	return false
}

// Helper function to get all direct combinations that create an element
func getDirectCombinations(elementName string) ([]struct {
	Item1, Item2 string
}, error) {
	var combinations []struct {
		Item1, Item2 string
	}

	rows, err := db.Query("SELECT item1, item2 FROM elements WHERE element = ? AND item1 IS NOT NULL AND item2 IS NOT NULL", elementName)
	if err != nil {
		return combinations, err
	}
	defer rows.Close()

	for rows.Next() {
		var item1, item2 string
		if err := rows.Scan(&item1, &item2); err != nil {
			continue
		}
		combinations = append(combinations, struct {
			Item1, Item2 string
		}{Item1: item1, Item2: item2})
	}

	return combinations, nil
}

// Helper function for default result when no recipe is found
func getDefaultResult(elementName string) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"name":     elementName,
			"image":    mapper[elementName],
			"children": []interface{}{},
			"recipe":   []string{"This is a basic element or no recipe found"},
		},
	}
}

// Format recipe steps for display
func formatRecipeSteps(steps []RecipeStep) []string {
	if len(steps) == 0 {
		return []string{}
	}

	formattedSteps := make([]string, 0, len(steps))
	for _, step := range steps {
		formattedSteps = append(formattedSteps, step.Result+" = "+step.Item1+" + "+step.Item2)
	}
	return formattedSteps
}

//================================================
// BFS IMPLEMENTATION
//================================================

// BFS for recipe search
func BFS(elementName string, recipeType string, maxRecipes int) ([]interface{}, int, float64) {
	start := time.Now()
	nodesVisited := 0

	// Check if element exists in database
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM elements WHERE element = ?)", elementName).Scan(&exists)
	if err != nil || !exists {
		return getDefaultResult(elementName), nodesVisited, float64(time.Since(start).Milliseconds())
	}

	// Get all basic elements
	basicElements := getBasicElements()
	
	// Check if this is already a basic element
	if isBasicElement(elementName, basicElements) {
		return getDefaultResult(elementName), nodesVisited, float64(time.Since(start).Milliseconds())
	}

	// Determine the number of recipes to find based on recipeType
	var desiredRecipeCount int
	if recipeType == "One" {
		desiredRecipeCount = 1
	} else if recipeType == "Limit" {
		desiredRecipeCount = maxRecipes
	} else {
		// For "All", set to a very large number to find all recipes
		desiredRecipeCount = 1000000
	}

	// Find recipes with early stopping
	allRecipes, nodesVisitedCount := findRecipesBFS(elementName, basicElements, desiredRecipeCount)
	nodesVisited = nodesVisitedCount
	
	// If no recipes found, return default
	if len(allRecipes) == 0 {
		return getDefaultResult(elementName), nodesVisited, float64(time.Since(start).Milliseconds())
	}
	
	// Convert recipes to result format
	var results []interface{}
	for _, recipe := range allRecipes {
		// Create tree representation
		treeRoot := createRecipeTree(elementName, recipe)
		results = append(results, treeRoot)
	}

	return results, nodesVisited, float64(time.Since(start).Milliseconds())
}

// Function to find recipes for an element using BFS with early stopping
func findRecipesBFS(elementName string, basicElements []string, maxRecipesToFind int) ([][]RecipeStep, int) {
	var allRecipes [][]RecipeStep
	nodesVisited := 0
	
	// Queue for BFS
	type QueueItem struct {
		Element string
		Path    []RecipeStep
		Explored map[string]bool // Track which elements are already explored in this path
	}
	
	// Start with the target element
	queue := []QueueItem{{
		Element: elementName,
		Path:    []RecipeStep{},
		Explored: make(map[string]bool),
	}}
	
	// Keep track of combinations we've added
	processedCombinations := make(map[string]bool)
	
	for len(queue) > 0 && len(allRecipes) < maxRecipesToFind {
		current := queue[0]
		queue = queue[1:]
		nodesVisited++
		
		// Skip if we've already explored this element in the current path to avoid cycles
		if current.Explored[current.Element] {
			continue
		}
		
		// Mark this element as explored in this path
		explored := make(map[string]bool)
		for k, v := range current.Explored {
			explored[k] = v
		}
		explored[current.Element] = true
		
		// Check if we've reached all basic elements
		allBasic := true
		for _, step := range current.Path {
			// Check if both ingredients are basic
			item1Basic := isBasicElement(step.Item1, basicElements)
			item2Basic := isBasicElement(step.Item2, basicElements)
			
			if !item1Basic || !item2Basic {
				allBasic = false
				break
			}
		}
		
		// If all elements in the path are decomposed to basic elements and we have steps
		if allBasic && len(current.Path) > 0 {
			// Create a unique key for this recipe
			recipeKey := ""
			for _, step := range current.Path {
				recipeKey += step.Result + step.Item1 + step.Item2 + "|"
			}
			
			// Only add if we haven't processed this exact recipe before
			if !processedCombinations[recipeKey] {
				processedCombinations[recipeKey] = true
				allRecipes = append(allRecipes, current.Path)
				
				// Check if we've found enough recipes
				if len(allRecipes) >= maxRecipesToFind {
					break
				}
			}
			continue
		}
		
		// Get all combinations for this element
		combinations, err := getDirectCombinations(current.Element)
		if err != nil || len(combinations) == 0 {
			// If there are no combinations (basic element or missing), and this is the target element
			if current.Element == elementName {
				// Try next element in queue
				continue
			}
			
			// This path can't be completed, don't add to results
			continue
		}
		
		// Process each combination
		for _, combo := range combinations {
			// Create new step
			newStep := RecipeStep{
				Result: current.Element,
				Item1:  combo.Item1,
				Item2:  combo.Item2,
			}
			
			// Create new path with this step
			newPath := make([]RecipeStep, len(current.Path))
			copy(newPath, current.Path)
			newPath = append(newPath, newStep)
			
			// Add both ingredients to queue to continue exploration
			if !isBasicElement(combo.Item1, basicElements) {
				queue = append(queue, QueueItem{
					Element: combo.Item1,
					Path:    newPath,
					Explored: explored,
				})
			}
			
			if !isBasicElement(combo.Item2, basicElements) {
				queue = append(queue, QueueItem{
					Element: combo.Item2,
					Path:    newPath,
					Explored: explored,
				})
			}
			
			// If both ingredients are basic elements, check if we have a complete path
			if isBasicElement(combo.Item1, basicElements) && isBasicElement(combo.Item2, basicElements) {
				// Create a unique key for this recipe
				recipeKey := ""
				for _, step := range newPath {
					recipeKey += step.Result + step.Item1 + step.Item2 + "|"
				}
				
				// Only add if we haven't processed this exact recipe before
				if !processedCombinations[recipeKey] {
					processedCombinations[recipeKey] = true
					allRecipes = append(allRecipes, newPath)
					
					// Check if we've found enough recipes
					if len(allRecipes) >= maxRecipesToFind {
						break
					}
				}
			}
		}
	}
	
	return allRecipes, nodesVisited
}

// Create a tree representation for a recipe
func createRecipeTree(elementName string, recipe []RecipeStep) map[string]interface{} {
	return map[string]interface{}{
		"name":     elementName,
		"image":    mapper[elementName],
		"children": buildElementTree(elementName, recipe),
		"recipe":   formatRecipeSteps(recipe),
	}
}

// Build tree for an element recursively
func buildElementTree(elementName string, recipe []RecipeStep) []map[string]interface{} {
	// Find the step for this element
	var stepForElement *RecipeStep
	for i, step := range recipe {
		if step.Result == elementName {
			stepForElement = &recipe[i]
			break
		}
	}
	
	// If not found, this is a basic element
	if stepForElement == nil {
		return []map[string]interface{}{}
	}
	
	// Create nodes for ingredients
	item1Node := map[string]interface{}{
		"name":  stepForElement.Item1,
		"image": mapper[stepForElement.Item1],
	}
	
	item2Node := map[string]interface{}{
		"name":  stepForElement.Item2,
		"image": mapper[stepForElement.Item2],
	}
	
	// Recursively build trees for ingredients
	item1Node["children"] = buildElementTree(stepForElement.Item1, recipe)
	item2Node["children"] = buildElementTree(stepForElement.Item2, recipe)
	
	return []map[string]interface{}{item1Node, item2Node}
}

//================================================
// DFS IMPLEMENTATION
//================================================

// DFS for recipe search
func DFS(elementName string, recipeType string, maxRecipes int) ([]interface{}, int, float64) {
	start := time.Now()
	nodesVisited := 0

	// Check if element exists in database
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM elements WHERE element = ?)", elementName).Scan(&exists)
	if err != nil || !exists {
		return getDefaultResult(elementName), nodesVisited, float64(time.Since(start).Milliseconds())
	}

	// Get all basic elements
	basicElements := getBasicElements()
	
	// Check if this is already a basic element
	if isBasicElement(elementName, basicElements) {
		return getDefaultResult(elementName), nodesVisited, float64(time.Since(start).Milliseconds())
	}

	// Determine the number of recipes to find based on recipeType
	var desiredRecipeCount int
	if recipeType == "One" {
		desiredRecipeCount = 1
	} else if recipeType == "Limit" {
		desiredRecipeCount = maxRecipes
	} else {
		// For "All", set to a very large number to find all recipes
		desiredRecipeCount = 1000000
	}

	// Find recipes with early stopping
	allRecipes, nodesVisitedCount := findRecipesDFS(elementName, basicElements, desiredRecipeCount)
	nodesVisited = nodesVisitedCount
	
	// If no recipes found, return default
	if len(allRecipes) == 0 {
		return getDefaultResult(elementName), nodesVisited, float64(time.Since(start).Milliseconds())
	}
	
	// Convert recipes to result format
	var results []interface{}
	for _, recipe := range allRecipes {
		// Create tree representation
		treeRoot := createRecipeTree(elementName, recipe)
		
		results = append(results, treeRoot)
	}

	return results, nodesVisited, float64(time.Since(start).Milliseconds())
}

// Function to find recipes for an element using DFS with early stopping
func findRecipesDFS(elementName string, basicElements []string, maxRecipesToFind int) ([][]RecipeStep, int) {
	var allRecipes [][]RecipeStep
	nodesVisited := 0
	
	// Stack for DFS
	type StackItem struct {
		Element string
		Path    []RecipeStep
		Explored map[string]bool // Track which elements are already explored in this path
	}
	
	// Start with the target element
	stack := []StackItem{{
		Element: elementName,
		Path:    []RecipeStep{},
		Explored: make(map[string]bool),
	}}
	
	// Keep track of combinations we've added
	processedCombinations := make(map[string]bool)
	
	for len(stack) > 0 && len(allRecipes) < maxRecipesToFind {
		// Pop from stack (last in, first out)
		last := len(stack) - 1
		current := stack[last]
		stack = stack[:last]
		nodesVisited++
		
		// Skip if we've already explored this element in the current path to avoid cycles
		if current.Explored[current.Element] {
			continue
		}
		
		// Mark this element as explored in this path
		explored := make(map[string]bool)
		for k, v := range current.Explored {
			explored[k] = v
		}
		explored[current.Element] = true
		
		// Check if we've reached all basic elements
		allBasic := true
		for _, step := range current.Path {
			// Check if both ingredients are basic
			item1Basic := isBasicElement(step.Item1, basicElements)
			item2Basic := isBasicElement(step.Item2, basicElements)
			
			if !item1Basic || !item2Basic {
				allBasic = false
				break
			}
		}
		
		// If all elements in the path are decomposed to basic elements and we have steps
		if allBasic && len(current.Path) > 0 {
			// Create a unique key for this recipe
			recipeKey := ""
			for _, step := range current.Path {
				recipeKey += step.Result + step.Item1 + step.Item2 + "|"
			}
			
			// Only add if we haven't processed this exact recipe before
			if !processedCombinations[recipeKey] {
				processedCombinations[recipeKey] = true
				allRecipes = append(allRecipes, current.Path)
				
				// Check if we've found enough recipes
				if len(allRecipes) >= maxRecipesToFind {
					break
				}
			}
			continue
		}
		
		// Get all combinations for this element
		combinations, err := getDirectCombinations(current.Element)
		if err != nil || len(combinations) == 0 {
			// If there are no combinations (basic element or missing), and this is the target element
			if current.Element == elementName {
				// Try next element in stack
				continue
			}
			
			// This path can't be completed, don't add to results
			continue
		}
		
		// Process each combination
		for i := len(combinations) - 1; i >= 0; i-- { // Reverse order for DFS
			combo := combinations[i]
			
			// Create new step
			newStep := RecipeStep{
				Result: current.Element,
				Item1:  combo.Item1,
				Item2:  combo.Item2,
			}
			
			// Create new path with this step
			newPath := make([]RecipeStep, len(current.Path))
			copy(newPath, current.Path)
			newPath = append(newPath, newStep)
			
			// Add both ingredients to stack in reverse order (so item1 is processed first)
			if !isBasicElement(combo.Item2, basicElements) {
				stack = append(stack, StackItem{
					Element: combo.Item2,
					Path:    newPath,
					Explored: explored,
				})
			}
			
			if !isBasicElement(combo.Item1, basicElements) {
				stack = append(stack, StackItem{
					Element: combo.Item1,
					Path:    newPath,
					Explored: explored,
				})
			}
			
			// If both ingredients are basic elements, check if we have a complete path
			if isBasicElement(combo.Item1, basicElements) && isBasicElement(combo.Item2, basicElements) {
				// Create a unique key for this recipe
				recipeKey := ""
				for _, step := range newPath {
					recipeKey += step.Result + step.Item1 + step.Item2 + "|"
				}
				
				// Only add if we haven't processed this exact recipe before
				if !processedCombinations[recipeKey] {
					processedCombinations[recipeKey] = true
					allRecipes = append(allRecipes, newPath)
					
					// Check if we've found enough recipes
					if len(allRecipes) >= maxRecipesToFind {
						break
					}
				}
			}
		}
	}
	
	return allRecipes, nodesVisited
}

//================================================
// BIDIRECTIONAL IMPLEMENTATION
//================================================

// Bidirectional search for recipes
func Bidirectional(elementName string, recipeType string, maxRecipes int) ([]interface{}, int, float64) {
	start := time.Now()
	nodesVisited := 0

	// Check if element exists in database
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM elements WHERE element = ?)", elementName).Scan(&exists)
	if err != nil || !exists {
		return getDefaultResult(elementName), nodesVisited, float64(time.Since(start).Milliseconds())
	}

	// Get all basic elements
	basicElements := getBasicElements()
	
	// Check if this is already a basic element
	if isBasicElement(elementName, basicElements) {
		return getDefaultResult(elementName), nodesVisited, float64(time.Since(start).Milliseconds())
	}

	// Determine the number of recipes to find based on recipeType
	var desiredRecipeCount int
	if recipeType == "One" {
		desiredRecipeCount = 1
	} else if recipeType == "Limit" {
		desiredRecipeCount = maxRecipes
	} else {
		// For "All", set to a very large number to find all recipes
		desiredRecipeCount = 1000000
	}

	// Find recipes with early stopping
	allRecipes, nodesVisitedCount := findRecipesBidirectional(elementName, basicElements, desiredRecipeCount)
	nodesVisited = nodesVisitedCount
	
	// If no recipes found, return default
	if len(allRecipes) == 0 {
		return getDefaultResult(elementName), nodesVisited, float64(time.Since(start).Milliseconds())
	}
	
	// Convert recipes to result format
	var results []interface{}
	for _, recipe := range allRecipes {
		// Create tree representation
		treeRoot := createRecipeTree(elementName, recipe)
		
		results = append(results, treeRoot)
	}

	return results, nodesVisited, float64(time.Since(start).Milliseconds())
}

// Function to find recipes using bidirectional search with early stopping
func findRecipesBidirectional(elementName string, basicElements []string, maxRecipesToFind int) ([][]RecipeStep, int) {
	var allRecipes [][]RecipeStep
	nodesVisited := 0
	
	// Keep track of combinations we've added
	processedCombinations := make(map[string]bool)
	
	// Forward search from target element
	type ForwardItem struct {
		Element string
		Path    []RecipeStep
		Explored map[string]bool
	}
	
	// Backward search from basic elements
	type BackwardItem struct {
		Element string
		Path    []RecipeStep
		Explored map[string]bool
	}
	
	// Initialize forward queue with target element
	forwardQueue := []ForwardItem{{
		Element: elementName,
		Path:    []RecipeStep{},
		Explored: make(map[string]bool),
	}}
	
	// Initialize backward queues with basic elements
	backwardQueue := []BackwardItem{}
	for _, basic := range basicElements {
		backwardQueue = append(backwardQueue, BackwardItem{
			Element: basic,
			Path:    []RecipeStep{},
			Explored: make(map[string]bool),
		})
	}
	
	// Process forward queue first to find direct paths
	for len(forwardQueue) > 0 && len(allRecipes) < maxRecipesToFind {
		current := forwardQueue[0]
		forwardQueue = forwardQueue[1:]
		nodesVisited++
		
		// Skip if we've already explored this element in the current path
		if current.Explored[current.Element] {
			continue
		}
		
		// Mark as explored in this path
		explored := make(map[string]bool)
		for k, v := range current.Explored {
			explored[k] = v
		}
		explored[current.Element] = true
		
		// Check if all elements in the path are decomposed to basic elements
		allBasic := true
		for _, step := range current.Path {
			if !isBasicElement(step.Item1, basicElements) || !isBasicElement(step.Item2, basicElements) {
				allBasic = false
				break
			}
		}
		
		// If we have a complete path to basic elements
		if allBasic && len(current.Path) > 0 {
			// Create a unique key for this recipe
			recipeKey := ""
			for _, step := range current.Path {
				recipeKey += step.Result + step.Item1 + step.Item2 + "|"
			}
			
			// Only add if we haven't processed this exact recipe before
			if !processedCombinations[recipeKey] {
				processedCombinations[recipeKey] = true
				allRecipes = append(allRecipes, current.Path)
				
				// Check if we've found enough recipes
				if len(allRecipes) >= maxRecipesToFind {
					break
				}
			}
			continue
		}
		
		// Get all combinations for this element
		combinations, err := getDirectCombinations(current.Element)
		if err != nil || len(combinations) == 0 {
			continue
		}
		
		// Process each combination
		for _, combo := range combinations {
			// Create new step
			newStep := RecipeStep{
				Result: current.Element,
				Item1:  combo.Item1,
				Item2:  combo.Item2,
			}
			
			// Create new path with this step
			newPath := make([]RecipeStep, len(current.Path))
			copy(newPath, current.Path)
			newPath = append(newPath, newStep)
			
			// Add to queue only if not a basic element
			if !isBasicElement(combo.Item1, basicElements) {
				forwardQueue = append(forwardQueue, ForwardItem{
					Element: combo.Item1,
					Path:    newPath,
					Explored: explored,
				})
			}
			
			if !isBasicElement(combo.Item2, basicElements) {
				forwardQueue = append(forwardQueue, ForwardItem{
					Element: combo.Item2,
					Path:    newPath,
					Explored: explored,
				})
			}
			
			// If both ingredients are basic, we have a complete path
			if isBasicElement(combo.Item1, basicElements) && isBasicElement(combo.Item2, basicElements) {
				// Create a unique key for this recipe
				recipeKey := ""
				for _, step := range newPath {
					recipeKey += step.Result + step.Item1 + step.Item2 + "|"
				}
				
				// Only add if we haven't processed this exact recipe before
				if !processedCombinations[recipeKey] {
					processedCombinations[recipeKey] = true
					allRecipes = append(allRecipes, newPath)
					
					// Check if we've found enough recipes
					if len(allRecipes) >= maxRecipesToFind {
						break
					}
				}
			}
		}
	}
	
	return allRecipes, nodesVisited
}