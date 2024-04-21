package gh

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// PullRequest act as the main entry to represent github PRs, for listing,
// getting details, etc...
type PullRequest struct {
	Name   string `json:"title"`
	State  string `json:"state"`
	Desc   string `json:"body"`
	URL    string `json:"url"`
	Number int    `json:"number"`
	Status string
}

func ListAllPrs() ([]PullRequest, error) {
	repo, err := getGitRemote()

	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls", repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("GPR_GH_TOKEN")))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", string(body))
		return nil, err
	}

	var prs []PullRequest
	err = json.Unmarshal(body, &prs)

	if err != nil {
		fmt.Println("Error unmarshaling json")
		return nil, err
	}

	items := make([]PullRequest, len(prs))

	for i, pr := range prs {
		prrStatusCmd := exec.Command("prr", "status", "-n")
		prrStatus, _ := prrStatusCmd.Output()

		prrStatuses := strings.Split(string(prrStatus), "\n")
		prrStatusContent := "UNKNOWN"

		for _, item := range prrStatuses {
			columns := strings.Fields(item)

			if len(columns) == 0 {
				continue
			}

			name, status := columns[0], columns[1]

			if strings.Contains(name, fmt.Sprint(pr.Number)) {
				prrStatusContent = status
			}
		}

		items[i] = PullRequest{
			Name:   pr.Name,
			State:  pr.State,
			Desc:   pr.Desc,
			URL:    pr.URL,
			Number: pr.Number,
			Status: prrStatusContent,
		}
	}

	return items, nil
}

func getGitRemote() (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	output, err := cmd.Output()

	if err != nil {
		return "", err
	}

	repoNameWithSuffix := strings.Trim(strings.Split(string(output), ":")[1], "\n")

	return strings.Replace(repoNameWithSuffix, ".git", "", -1), nil
}
