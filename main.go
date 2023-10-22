package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/marshyon/pinboard-bookmarks/api"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var bookmarks []ConvertedBookmark

type Bookmark struct {
	Href        string `json:"href"`
	Description string `json:"description"`
	Tags        string `json:"tags"`
	Time        string `json:"time"`
}

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

type ReadBookmarks struct {
	Hash string `json:"hash" gorm:"unique"`
	Done int    `json:"done"`
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
func printConvertedBookmarks(convertedBookmarks []ConvertedBookmark, db *gorm.DB) {
	for _, bookmark := range convertedBookmarks {

		// check if there is already a hash in ReadBookmarks
		// if there is, then skip
		if err := db.Where("hash = ?", bookmark.Hash).First(&ReadBookmarks{}).Error; err != nil {
			fmt.Println("Error reading bookmark:", err)
			return
		}

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

func storeTaggedBookmarks(bookmarks []ConvertedBookmark, tag string) {

	// if there is already a tag matching the tag, then skip
	tagCount, err := api.QueryTag(tag)
	if err != nil {
		fmt.Println("Error querying tag:", err)
		return
	}
	fmt.Printf("Tag count: %d\n", tagCount)
	if tagCount == 0 {
		fmt.Println("No bookmarks found for tag:", tag)
	} else {
		fmt.Printf("%d Bookmarks found for tag: %s\n", tagCount, tag)
		return
	}

	// there are no bookmarks for the tag yet,
	// so create new ones in linkding for each one
	for _, bookmark := range bookmarks {
		fmt.Printf("URL: %s\n", bookmark.Href)
		fmt.Printf("Desc: %s\n", bookmark.Description)
		fmt.Printf("Extended: %s\n", bookmark.Extended)

		newNotes := fmt.Sprintf("%s\n\n%s", bookmark.Hash, bookmark.Extended)

		fmt.Printf("Tags: %s\n", bookmark.Tags)
		fmt.Printf("Hash: %s\n", bookmark.Hash)

		dateString := bookmark.Time.Format("2006-01-02")
		newDescription := fmt.Sprintf("[%s] %s", dateString, bookmark.Description)
		fmt.Printf("Date: %s\n", dateString)

		// create a linkding bookmark
		linkdingBookmark := api.LinkdingBookmark{
			URL:         bookmark.Href,
			Title:       bookmark.Description,
			Description: newDescription,
			Notes:       newNotes,
			IsArchived:  false,
			Unread:      true,
			Shared:      false,
			TagNames:    []string{tag},
		}
		fmt.Printf("Linkding bookmark: %#v\n", linkdingBookmark)
		api.CreateBookmark(linkdingBookmark)
		fmt.Println("---------------------------------------------")

	}
}

func main() {

	// define command line flags
	upsertFlag := flag.Bool("upsert", false, "whether to upsert bookmarks")
	queryFlag := flag.String("query", "", "comma-separated values to query bookmarks")
	verboseFlag := flag.Bool("verbose", false, "whether to print bookmarks")

	// parse command line flags
	flag.Parse()

	db, err := initDB()
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	//  defer db.Close()

	if *queryFlag != "" {
		// get the query values
		queryValues := strings.Split(*queryFlag, ",")

		fmt.Printf("Querying bookmarks...[%#v]", queryValues)

		// iterate over queryValues and query the database for bookmarks that match the query values
		for _, queryValue := range queryValues {

			// query the database for bookmarks that match the query values
			db.Where("tags LIKE ?", "%"+queryValue+"%").Order("time ASC").Find(&bookmarks)
			storeTaggedBookmarks(bookmarks, queryValue)
		}

		os.Exit(0)
	}

	// if upsert flag is set, upsert bookmarks
	if *upsertFlag {
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
		if *verboseFlag {
			printConvertedBookmarks(ConvertedBookmarks, db)
		}
	}
}
