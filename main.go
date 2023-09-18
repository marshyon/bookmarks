package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Define a struct that matches the JSON structure
type Bookmark struct {
	Href        string `json:"href"`
	Description string `json:"description"`
	Extended    string `json:"extended"`
	Meta        string `json:"meta"`
	Hash        string `json:"hash"`
	Time        string `json:"time"`
	Shared      string `json:"shared"`
	ToRead      string `json:"toread"`
	Tags        string `json:"tags"`
}

func main() {
	// Read the JSON file
	data, err := os.ReadFile("bookmarks.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Create a slice to hold the bookmarks
	var bookmarks []Bookmark

	// Unmarshal the JSON data into the bookmarks slice
	if err := json.Unmarshal(data, &bookmarks); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Iterate over the bookmarks and print their contents
	for _, bookmark := range bookmarks {
		fmt.Printf("Href: %s\n", bookmark.Href)
		fmt.Printf("Description: %s\n", bookmark.Description)
		fmt.Printf("Extended: %s\n", bookmark.Extended)
		fmt.Printf("Meta: %s\n", bookmark.Meta)
		fmt.Printf("Hash: %s\n", bookmark.Hash)

		// Parse the time string into a time.Time object
		t, err := time.Parse(time.RFC3339, bookmark.Time)
		if err != nil {
			fmt.Println("Error parsing time:", err)
		} else {
			// Print the time as a formatted string or epoch timestamp
			fmt.Printf("Time (Formatted): %s\n", t.Format("2006-01-02 15:04:05"))
			fmt.Printf("Time (Epoch): %d\n", t.Unix())
		}

		fmt.Printf("Shared: %s\n", bookmark.Shared)
		fmt.Printf("ToRead: %s\n", bookmark.ToRead)
		fmt.Printf("Tags: %s\n", bookmark.Tags)
		fmt.Println("---------------------------------------------")
	}
}
