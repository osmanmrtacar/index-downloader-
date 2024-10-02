package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

// Function to download a file from a URL and save it in the specified folder
func downloadFile(url, folder, fileName string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the folder if it doesn't exist
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		return err
	}

	// Create the full file path
	filePath := filepath.Join(folder, fileName)

	// Create the file locally
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the content to the file
	_, err = io.Copy(out, resp.Body)
	return err
}

// Function to extract file links from an HTML page
func extractFileLinks(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)
	var links []string
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return links, nil
		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()
			if t.Data == "a" {
				for _, attr := range t.Attr {
					if attr.Key == "href" && strings.HasSuffix(attr.Val, ".pdf") {
						links = append(links, attr.Val)
					}
				}
			}
		}
	}
}

func main() {
	// Get the URL and folder name from the command-line arguments
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run yourscript.go <URL> <folder>")
		return
	}

	baseURL := os.Args[1]
	folderName := os.Args[2]

	// Extract all the PDF file links
	links, err := extractFileLinks(baseURL)
	if err != nil {
		fmt.Println("Error extracting links:", err)
		return
	}

	for _, link := range links {
		// Complete URL
		fileURL := baseURL + link
		fileName := link

		fmt.Printf("Downloading %s...\n", fileName)
		if err := downloadFile(fileURL, folderName, fileName); err != nil {
			fmt.Println("Error downloading file:", err)
			continue
		}
		fmt.Printf("Downloaded %s successfully!\n", fileName)
	}
}
