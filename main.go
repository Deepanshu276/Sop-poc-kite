package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/russross/blackfriday"
	"golang.org/x/net/html"
)

func main() {
	// Replace with your GitHub username, repository, and file path

	url := "https://github.com/openshift/ops-sop/blob/master/v4/alerts/ClusterOperatorDown.md"

	// Extract username, repo name, and path from the URL
	parts := strings.Split(strings.TrimPrefix(url, "https://github.com/"), "/")
	username := parts[0]
	repo := parts[1]
	path := strings.Join(parts[4:], "/")
	//fmt.Println(username)
	//fmt.Println(repo)
	//fmt.Println(path)

	// Build the API URL
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", username, repo, path)
	println(apiURL)

	// Create a new HTTP request with the GitHub API URL
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		panic(err)
	}

	// Set the Accept header to request the raw content of the file
	req.Header.Set("Accept", "application/vnd.github.v3.raw")

	// Replace with your GitHub personal access token
	token := "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"
	if token != "" {
		// Set the Authorization header to authenticate the request
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	// Send the HTTP request
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Convert markdown to HTML
	result := blackfriday.MarkdownBasic(body)

	var text string

	for _, code := range result {
		text += string(code)
	}

	// Convert HTML to plain text
	plainText := htmlToPlainText(text)

	fmt.Println(plainText)
}

// htmlToPlainText converts HTML to plain text.
func htmlToPlainText(htmlString string) string {
	var textBuilder strings.Builder

	doc, _ := html.Parse(strings.NewReader(htmlString))

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.TextNode {
			textBuilder.WriteString(n.Data)
			textBuilder.WriteString(" ")
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)

	return strings.TrimSpace(textBuilder.String())
}
