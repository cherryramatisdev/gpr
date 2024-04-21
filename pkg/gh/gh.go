package gh

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// PullRequest act as the main entry to represent github PRs, for listing,
// getting details, etc...
type PullRequest struct {
	Name   string `json:"title"`
	State  string `json:"state"`
	Desc   string `json:"body"`
	URL    string `json:"url"`
	Number int    `json:"number"`
}

func ListAllPrs() ([]PullRequest, error) {
	repo := "360-Hub/api-360-hub-new"
	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls", repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "Bearer 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")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	var prs []PullRequest
	err = json.Unmarshal(body, &prs)

	if err != nil {
		fmt.Println("Error unmarshaling json")
		return nil, err
	}

	return prs, nil
}
