package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/a-h/openanalysis/github"
	"github.com/a-h/openanalysis/statistics"
)

// query the ids of users from the Github Membership page.
var query = `
const getMemberNames = () => {
	let inputs = document.getElementsByTagName("input");
	let members = [];
	for(let i = 0; i < inputs.length; i ++) {
		if((inputs[i].getAttribute("name") == "members[]" || inputs[i].getAttribute("name") == "outside_collaborator[]") && inputs[i].getAttribute("type") == "checkbox") {
			members.push(inputs[i].value);
		}
	}
	return members;
};
`

var tokenFlag = flag.String("token", "", "GitHub auth token")
var loginFlag = flag.String("logins", "users.json", "GitHub logins")

func main() {
	flag.Parse()

	start := time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)

	c := github.NewCollector(*tokenFlag)
	logins, err := getLogins(*loginFlag)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, login := range logins {
		stats, err := getStats(c, login, start, end)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = ioutil.WriteFile("output/"+filepath.Clean(login)+".json", jsonOrNothing(stats), 0640)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func jsonOrNothing(v interface{}) []byte {
	b, _ := json.MarshalIndent(v, "", "  ")
	return b
}

func getLogins(fileName string) (logins []string, err error) {
	r, err := os.Open(fileName)
	if err != nil {
		return
	}
	d := json.NewDecoder(r)
	err = d.Decode(&logins)
	return
}

func getStats(c *github.Collector, login string, start, end time.Time) (stats *statistics.Statistics, err error) {
	stats = statistics.NewStatistics(start, end)

	// Issues.
	issues, err := c.UserIssues(context.Background(), login)
	if err != nil {
		err = fmt.Errorf("error getting issues: %v", err)
		return
	}
	for _, issue := range issues {
		if issue.RepoIsPrivate {
			continue
		}
		if isDayJobRepo(issue.RepoNameWithOwner) {
			continue
		}
		stats.Issues.Add(issue.CreatedAt, 1)
		stats.AddTouchedRepo(issue.CreatedAt, issue.RepoNameWithOwner)
	}

	// Pull Requests.
	prs, err := c.UserPullRequests(context.Background(), login)
	if err != nil {
		err = fmt.Errorf("error getting pull requests: %v", err)
		return
	}
	for _, pr := range prs {
		if pr.RepoIsPrivate {
			continue
		}
		if isDayJobRepo(pr.RepoNameWithOwner) {
			continue
		}
		if pr.MergedAt != nil {
			stats.PullRequestsMerged.Add(*pr.MergedAt, 1)
			stats.AddTouchedRepo(*pr.MergedAt, pr.RepoNameWithOwner)
		} else {
			stats.PullRequestsCreated.Add(pr.CreatedAt, 1)
			stats.AddTouchedRepo(pr.CreatedAt, pr.RepoNameWithOwner)
		}
	}

	repos, err := c.UserRepositories(context.Background(), login)
	if err != nil {
		err = fmt.Errorf("error getting pull repos: %v", err)
		return
	}
	for _, r := range repos {
		stats.ReposUpdated.Add(r.PushedAt, 1)
		stats.AddStars(r.PushedAt, r.Stars)
		stats.IncrementRepo(r.PushedAt)
		stats.AddTouchedRepo(r.PushedAt, r.RepoNameWithOwner)
	}
	return
}

func isDayJobRepo(name string) bool {
	if strings.HasPrefix(name, "infinityworks/") {
		return true
	}
	if strings.HasPrefix(name, "AgileVentures/") {
		return true
	}
	if strings.HasPrefix(name, "makersacademy/") {
		return true
	}
	if strings.Contains(strings.ToLower(name), "nhs") {
		return true
	}
	if strings.HasPrefix(name, "learn-co-students/") {
		return true
	}
	return false
}
