package library

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"encoding/json"
)

type GithubRepo struct {
	User string
	Repo string
	Link string
	Description string
	Language string
}

/*
 * Handles the checking if the input has a Github link in it, looks up the data from the API, and if successful
 * returns a string containing info about the repository.
 */
func HandleGithubLink(str string) string {
	var response = ""

	repo := ParseGithubLink(str)

	if (GithubRepo{}) != repo {
		repo.LookUpGithubData()

		if repo.Description != "" {
			response = fmt.Sprintf("'%s' is a project written in %s by %s. It is described as: %s",
				repo.Repo, repo.Language, repo.User, repo.Description)
		}
	}

	return response
}

/*
 * Checks if the input string contains a link to a Github repo and prepares an object containing the pieces
 */
func ParseGithubLink(input string) GithubRepo {
	re, err := regexp.Compile(`(https?:\/\/)?github\.com\/([\w\-]+)\/([\w\-]+)`)
	
	if err == nil {
		res := re.FindStringSubmatch(input)

		if len(res) == 4 {
			repo := GithubRepo{
				User: res[2],
				Repo: res[3],
				Link: res[0],
			}
			
			return repo
		}
	}
	
	return GithubRepo{}
}

/*
 * Look up data from the Github API for the specified repo
 */
func (repo *GithubRepo) LookUpGithubData() {
	apiLink := fmt.Sprintf("https://api.github.com/repos/%s/%s", repo.User, repo.Repo)
	resp, _ := http.Get(apiLink)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err == nil { //else there was an error with the API call and we should just do nothing
		var rawData interface{}
		json.Unmarshal(body, &rawData)

		data := rawData.(map[string]interface{})

		if data["language"] != nil { //if it's nil then I probably can't find this repo
			repo.Language = data["language"].(string)
			repo.Description = data["description"].(string)
		}
	}
}