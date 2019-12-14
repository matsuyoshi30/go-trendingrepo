package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const RepoURL = "https://api.github.com/search/repositories"

type Response struct {
	TotalCount        int          `json:"total_count"`
	IncompleteResults bool         `json:"incomplete_results"`
	Repos             []Repository `json:"items"`
}

type Repository struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	HTMLURL     string    `json:"html_url"`
	Homepage    string    `json:"homepage"`
	Language    string    `json:"language"`
}

func main() {
	qs := time.Now().AddDate(0, 0, -7).Format("2006-01-02")

	query := flag.String("q", "created:>"+qs, "Input uery string")
	lang := flag.String("l", "", "Select language")
	issue := flag.Bool("i", false, "Sort help-wanted-issues num")
	asc := flag.Bool("a", false, "Order ascend")
	num := flag.Int("n", 5, "Input list num (<= 15)")
	flag.Parse()

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, RepoURL, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/vnd.github.preview")

	params := req.URL.Query()
	params.Add("q", *query)
	params.Add("sort", "stars")
	params.Add("order", "desc")

	if *lang != "" {
		params.Set("q", "language:"+*lang+"+"+*query)
	}
	if *asc {
		params.Set("order", "asc")
	}
	if *issue {
		params.Set("sort", "help-wanted-issues")
	}
	req.URL.RawQuery = strings.ReplaceAll(params.Encode(), "%2B", "+") // need replace '%2B' to '+'

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println(err)
		return
	}

	listnum := 5
	if 0 <= *num && *num <= 15 {
		listnum = *num
	}

	for i, r := range response.Repos {
		fmt.Printf("------------------------------\n")
		if i == listnum {
			break
		}
		fmt.Printf("[Name]        %s\n", r.Name)
		fmt.Printf("[Description] %s\n", r.Description)
		fmt.Printf("[Created At]  %s\n", r.CreatedAt)
		fmt.Printf("[URL]         %s\n", r.HTMLURL)
		fmt.Printf("[Homepage]    %s\n", r.Homepage)
		fmt.Printf("[Language]    %s\n", r.Language)
	}
}
