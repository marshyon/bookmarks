package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	Hash        string    `json:"hash" gorm:"unique"`
	Time        time.Time `json:"time"`
	Shared      string    `json:"shared"`
	ToRead      string    `json:"toread"`
	Tags        string    `json:"tags"`
}

func main() {
	// Open an SQLite database
	db, err := gorm.Open(sqlite.Open("bookmarks.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}

	// Migrate the schema
	db.AutoMigrate(&ConvertedBookmark{})

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

	// Iterate over the original bookmarks and convert them
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

		// Perform upsert operation
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "hash"}},
			UpdateAll: true,
		}).Create(&convertedBookmark).Error; err != nil {
			fmt.Println("Error upserting bookmark:", err)
		}
	}

	// Now, you have the data upserted in the "bookmarks.db" SQLite database

	// Iterate over the converted bookmarks and print their contents
	for _, bookmark := range originalBookmarks {
		fmt.Printf("Href: %s\n", bookmark.Href)
		fmt.Printf("Description: %s\n", bookmark.Description)
		fmt.Printf("Extended: %s\n", bookmark.Extended)
		fmt.Printf("Meta: %s\n", bookmark.Meta)
		fmt.Printf("Hash: %s\n", bookmark.Hash)
		fmt.Printf("Time: %s\n", bookmark.Time)
		fmt.Printf("Shared: %s\n", bookmark.Shared)
		fmt.Printf("ToRead: %s\n", bookmark.ToRead)
		fmt.Printf("Tags: %s\n", bookmark.Tags)
		fmt.Println("---------------------------------------------")
	}
}
