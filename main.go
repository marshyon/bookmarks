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

func initDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("bookmarks.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("Error opening database:", err)
		return nil, err
	}
	fmt.Println("Database connection successfully opened.")
	return db, nil
}

// function to read the json file
func readJSONFile() ([]byte, error) {
	data, err := os.ReadFile("bookmarks.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}
	return data, nil
}

// function to
// Create a slice to hold the original bookmarks
// Unmarshal the JSON data into the original bookmarks slice
func unmarshalJSON(data []byte) ([]OriginalBookmark, error) {
	var originalBookmarks []OriginalBookmark
	if err := json.Unmarshal(data, &originalBookmarks); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil, err
	}
	return originalBookmarks, nil
}

// function to
// Iterate over the original bookmarks and convert them
func convertBookmarks(originalBookmarks []OriginalBookmark) ([]ConvertedBookmark, error) {
	var convertedBookmarks []ConvertedBookmark
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
		convertedBookmarks = append(convertedBookmarks, convertedBookmark)
	}
	return convertedBookmarks, nil
}

// function to upsert the converted bookmarks and output count ever other second
func upsertBookmarks(db *gorm.DB, convertedBookmarks []ConvertedBookmark) error {
	count := 0
	total := len(convertedBookmarks)
	for _, convertedBookmark := range convertedBookmarks {
		count++
		if count%500 == 0 {
			fmt.Printf("Upserted %d bookmarks :: ", count)
			percentages := float64(count) / float64(total) * 100
			fmt.Printf("[%.0f%%]\n", percentages)
		}
		// Perform upsert operation
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "hash"}},
			UpdateAll: true,
		}).Create(&convertedBookmark).Error; err != nil {
			fmt.Println("Error upserting bookmark:", err)
			return err
		}
	}
	return nil
}

// print out the converted bookmarks
func printConvertedBookmarks(convertedBookmarks []ConvertedBookmark) {
	for _, bookmark := range convertedBookmarks {
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

func main() {

	// read command line arguments and check to see if verbose is set
	// if verbose is set then print out the converted bookmarks
	verbose := false
	if len(os.Args) > 1 {
		if os.Args[1] == "-v" {
			verbose = true
		}
	}

	db, err := initDB()
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	// defer db.Close()

	db.AutoMigrate(&ConvertedBookmark{})

	data, err := readJSONFile()
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	originalBookmarks, err := unmarshalJSON(data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	ConvertedBookmarks, err := convertBookmarks(originalBookmarks)
	upsertBookmarks(db, ConvertedBookmarks)
	if err != nil {
		fmt.Println("Error converting bookmarks:", err)
		return
	}

	total := len(ConvertedBookmarks)
	fmt.Printf("Total bookmarks: %d\n", total)
	if verbose {
		printConvertedBookmarks(ConvertedBookmarks)
	}
}
