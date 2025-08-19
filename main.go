package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Global variable to store words from common.txt
var commonWords []string

// Number of words to use from common.txt (configurable)
const wordLimit = 500

// init function to load words from common.txt at startup
func init() {
	// Initialize with some basic words as fallback
	commonWords = []string{"test", "data", "json", "benchmark", "encoding", "decoding", "performance", "measurement"}
}

// GeoCoordinates represents geographic coordinates
type GeoCoordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// TreeNode represents a tree node with many fields
type TreeNode struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	Value       float64                `json:"value"`
	Description string                 `json:"description"`
	LongText    string                 `json:"long_text"`
	Tags        []string               `json:"tags"`
	Metadata    map[string]interface{} `json:"metadata"`
	URL1        string                 `json:"url1"`
	URL2        string                 `json:"url2"`
	Geo         GeoCoordinates         `json:"geo"`
	Children    []*TreeNode            `json:"children"`
}

// readWordsFromFile reads words from common.txt file
func readWordsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			words = append(words, word)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

// generateRealisticText creates realistic text using common English words from common.txt
func generateRealisticText(wordCount int) string {
	if wordCount <= 0 {
		return ""
	}

	// Generate a string with approximately wordCount words
	result := ""
	for i := 0; i < wordCount; i++ {
		// Add a random word from our words list
		word := commonWords[rand.Intn(len(commonWords))]
		if i > 0 {
			result += " "
		}
		result += word
	}

	return result
}

// generateRandomString generates a random string of specified length
// This is kept for compatibility with existing code patterns
func generateRandomString(length int) string {
	// For backward compatibility, we'll generate a string of about the requested length
	// using our realistic words
	if length <= 0 {
		return ""
	}

	// Estimate word count based on average word length
	avgWordLength := 8
	wordCount := length / avgWordLength
	if wordCount < 1 {
		wordCount = 1
	}

	return generateRealisticText(wordCount)
}

// generateRandomTags creates a slice of random tags
func generateRandomTags(count int) []string {
	tags := make([]string, count)
	for i := 0; i < count; i++ {
		tags[i] = generateRandomString(rand.Intn(10) + 5)
	}
	return tags
}

// generateRandomMetadata creates a map with random metadata
func generateRandomMetadata() map[string]interface{} {
	metadata := make(map[string]interface{})
	for i := 0; i < 7; i++ {
		key := fmt.Sprintf("key_%d", i)
		if rand.Intn(2) == 0 {
			metadata[key] = generateRandomString(rand.Intn(15) + 8)
		} else {
			metadata[key] = rand.Int()
		}
	}
	return metadata
}

// generateTreeNode creates a single tree node with random data
func generateTreeNode(id int) *TreeNode {
	// Select random words for URL paths
	word1 := commonWords[rand.Intn(len(commonWords))]
	word2 := commonWords[rand.Intn(len(commonWords))]

	return &TreeNode{
		ID:          id,
		Name:        generateRandomString(rand.Intn(10) + 5),
		Value:       rand.Float64() * 1000,
		Description: generateRandomString(rand.Intn(50) + 20),
		LongText:    generateRandomString(rand.Intn(500) + 200), // Longer text field to help reach target size
		Tags:        generateRandomTags(rand.Intn(5) + 3),
		Metadata:    generateRandomMetadata(),
		URL1:        fmt.Sprintf("https://example.com/%s/%s-%d", word1, word2, rand.Intn(10000)),
		URL2:        fmt.Sprintf("https://api.example.com/v1/%s/%s-%d", word1, word2, rand.Intn(10000)),
		Geo: GeoCoordinates{
			Latitude:  (rand.Float64() * 180) - 90,  // Range: -90 to 90
			Longitude: (rand.Float64() * 360) - 180, // Range: -180 to 180
		},
		Children: nil,
	}
}

// buildTree creates a tree structure with approximately the specified number of nodes
// and limits the depth to 6 layers
func buildTree(nodeCount int) *TreeNode {
	if nodeCount <= 0 {
		return nil
	}

	nodes := make([]*TreeNode, nodeCount)
	depths := make([]int, nodeCount) // Track depth of each node

	// Create all nodes
	for i := 0; i < nodeCount; i++ {
		nodes[i] = generateTreeNode(i)
		depths[i] = 0
	}

	// Build tree structure with depth limit
	root := nodes[0]
	depths[0] = 0

	for i := 1; i < nodeCount; i++ {
		// Find a random parent (excluding the current node) with depth less than 6
		parentIndex := rand.Intn(i)
		// Keep trying until we find a parent with depth < 5 (so we can add a child at depth 6)
		for depths[parentIndex] >= 5 {
			parentIndex = rand.Intn(i)
		}
		nodes[parentIndex].Children = append(nodes[parentIndex].Children, nodes[i])
		depths[i] = depths[parentIndex] + 1
	}

	return root
}

// countNodes counts the total number of nodes in the tree
func countNodes(root *TreeNode) int {
	if root == nil {
		return 0
	}

	count := 1
	for _, child := range root.Children {
		count += countNodes(child)
	}
	return count
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Load common words for realistic text generation
	words, err := readWordsFromFile("common.txt")
	if err != nil {
		fmt.Printf("Warning: Failed to load common.txt: %v\n", err)
		fmt.Println("Using fallback words for text generation")
	} else {
		// Use only a portion of the words (wordLimit) for better performance
		if len(words) > wordLimit {
			commonWords = words[:wordLimit]
		} else {
			commonWords = words
		}
		fmt.Printf("Loaded %d words from common.txt (using %d words)\n", len(words), len(commonWords))
	}

	fmt.Println("Generating tree data...")

	// Generate tree with approximately 10K nodes
	tree := buildTree(10000)
	actualNodeCount := countNodes(tree)

	fmt.Printf("Generated tree with %d nodes\n", actualNodeCount)

	// Marshal to JSON
	fmt.Println("Marshaling to JSON...")
	jsonData, err := json.Marshal(tree)
	if err != nil {
		fmt.Printf("Error marshaling: %v\n", err)
		return
	}

	fmt.Printf("JSON size: %.2f MB\n", float64(len(jsonData))/(1024*1024))

	// Save to file for benchmarking
	err = os.WriteFile("test_data.json", jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Println("Test data saved to test_data.json")

	// Test unmarshaling
	fmt.Println("Unmarshaling from JSON...")
	var unmarshaledTree TreeNode
	err = json.Unmarshal(jsonData, &unmarshaledTree)
	if err != nil {
		fmt.Printf("Error unmarshaling: %v\n", err)
		return
	}

	fmt.Printf("Unmarshaled tree has %d nodes\n", countNodes(&unmarshaledTree))
}
