package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type linkdingTagSearch struct {
	Count    int `json:"count"`
	Next     any `json:"next"`
	Previous any `json:"previous"`
	Results  []struct {
		ID                 int       `json:"id"`
		URL                string    `json:"url"`
		Title              string    `json:"title"`
		Description        string    `json:"description"`
		Notes              string    `json:"notes"`
		WebsiteTitle       string    `json:"website_title"`
		WebsiteDescription string    `json:"website_description"`
		IsArchived         bool      `json:"is_archived"`
		Unread             bool      `json:"unread"`
		Shared             bool      `json:"shared"`
		TagNames           []string  `json:"tag_names"`
		DateAdded          time.Time `json:"date_added"`
		DateModified       time.Time `json:"date_modified"`
	} `json:"results"`
}

type LinkdingBookmark struct {
	URL         string   `json:"url"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Notes       string   `json:"notes"`
	IsArchived  bool     `json:"is_archived"`
	Unread      bool     `json:"unread"`
	Shared      bool     `json:"shared"`
	TagNames    []string `json:"tag_names"`
}

type linkdingCreatedBookmark struct {
	ID                 int       `json:"id"`
	URL                string    `json:"url"`
	Title              string    `json:"title"`
	Description        string    `json:"description"`
	Notes              string    `json:"notes"`
	WebsiteTitle       string    `json:"website_title"`
	WebsiteDescription string    `json:"website_description"`
	IsArchived         bool      `json:"is_archived"`
	Unread             bool      `json:"unread"`
	Shared             bool      `json:"shared"`
	TagNames           []string  `json:"tag_names"`
	DateAdded          time.Time `json:"date_added"`
	DateModified       time.Time `json:"date_modified"`
}

type Config struct {
	APIKey string
	URL    string
}

func GetConfig() (*Config, error) {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return nil, err
	}

	// read environment variable
	apiKey := os.Getenv("LINKDING_API_KEY")
	if apiKey == "" {
		return nil, errors.New("LINKDING_API_KEY environment variable not set")
	}

	url := os.Getenv("LINKDING_URL")
	if url == "" {
		return nil, errors.New("LINKDING_URL environment variable not set")
	}

	return &Config{
		APIKey: apiKey,
		URL:    url,
	}, nil
}

func QueryTag(s string) (int, error) {

	config, err := GetConfig()
	if err != nil {
		fmt.Println("Error getting config:", err)
		return 0, err
	}

	apiKey := config.APIKey

	queryUrl := "http://localhost:9090/api/bookmarks/?q=%23" + s
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, queryUrl, nil)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	token := "Token " + apiKey

	req.Header.Add("Authorization", token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	var tagSearch linkdingTagSearch
	err = json.Unmarshal(body, &tagSearch)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	fmt.Printf("Count: %d\n", tagSearch.Count)

	return tagSearch.Count, nil
}

func CreateBookmark(bookmark LinkdingBookmark) error {

	config, err := GetConfig()
	if err != nil {
		fmt.Println("Error getting config:", err)
		return err
	}

	apiKey := config.APIKey                       // TODO: apikKey is not used
	url := "http://localhost:9090/api/bookmarks/" // TODO : createUrl is not used
	// _ = "POST"

	// use bookmark struct to create a JSON string
	bookmarkJSON, err := json.Marshal(bookmark)
	if err != nil {
		fmt.Println("Error marshalling bookmark:", err)
		return err
	}
	fmt.Printf("bookmarkJSON: %s\n", bookmarkJSON)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(bookmarkJSON))
	if err != nil {
		fmt.Println(err)
		return err
	}

	// req, err := http.NewRequest("POST", url, bookmarkJSON)

	if err != nil {
		fmt.Println(err)
		return err
	}

	token := "Token " + apiKey

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))
	newLinkdingBookmark := linkdingCreatedBookmark{}
	err = json.Unmarshal(body, &newLinkdingBookmark)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("ID of new bookmark: %d\n", newLinkdingBookmark.ID)
	return nil
}
