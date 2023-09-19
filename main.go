package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Define a struct that matches the JSON structure
type OriginalBookmark struct {
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

// Define a struct with a time.Time field for the converted bookmarks
type ConvertedBookmark struct {
	Href        string    `json:"href"`
	Description string    `json:"description"`
	Extended    string    `json:"extended"`
	Meta        string    `json:"meta"`
	Hash        string    `json:"hash"`
	Time        time.Time `json:"time"`
	Shared      string    `json:"shared"`
	ToRead      string    `json:"toread"`
	Tags        string    `json:"tags"`
}

func main() {
	// Read the JSON file
	data, err := os.ReadFile("bookmarks.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Create a slice to hold the original bookmarks
	var originalBookmarks []OriginalBookmark

	// Unmarshal the JSON data into the original bookmarks slice
	if err := json.Unmarshal(data, &originalBookmarks); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Create a slice to hold the converted bookmarks
	var convertedBookmarks []ConvertedBookmark

	// Convert the original bookmarks to converted bookmarks
	for _, bookmark := range originalBookmarks {
		// Parse the time string into a time.Time object
		t, err := time.Parse(time.RFC3339, bookmark.Time)
		if err != nil {
			fmt.Println("Error parsing time:", err)
			continue
		}

		// Create a converted bookmark
		convertedBookmark := ConvertedBookmark{
			Href:        bookmark.Href,
			Description: bookmark.Description,
			Extended:    bookmark.Extended,
			Meta:        bookmark.Meta,
			Hash:        bookmark.Hash,
			Time:        t,
			Shared:      bookmark.Shared,
			ToRead:      bookmark.ToRead,
			Tags:        bookmark.Tags,
		}

		// Append the converted bookmark to the slice
		convertedBookmarks = append(convertedBookmarks, convertedBookmark)
	}

	// Now you can work with the convertedBookmarks slice, where the Time field is of type time.Time

	// Iterate over the converted bookmarks and print their contents
	for _, bookmark := range convertedBookmarks {
		fmt.Printf("Href: %s\n", bookmark.Href)
		fmt.Printf("Description: %s\n", bookmark.Description)
		fmt.Printf("Extended: %s\n", bookmark.Extended)
		fmt.Printf("Meta: %s\n", bookmark.Meta)
		fmt.Printf("Hash: %s\n", bookmark.Hash)
		fmt.Printf("Time (Formatted): %s\n", bookmark.Time.Format("2006-01-02 15:04:05"))
		fmt.Printf("Shared: %s\n", bookmark.Shared)
		fmt.Printf("ToRead: %s\n", bookmark.ToRead)
		fmt.Printf("Tags: %s\n", bookmark.Tags)
		fmt.Println("---------------------------------------------")
	}
}
