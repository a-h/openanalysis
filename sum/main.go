package main

import (
	"encoding/json"
	"fmt"

	"github.com/a-h/openanalysis/read"
	"github.com/a-h/openanalysis/statistics"
	"github.com/a-h/setof"
)

func main() {
	stats, err := read.UserStatistics("../output")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(jsonOrNothing(summarise(stats))))
}

func jsonOrNothing(v interface{}) []byte {
	b, _ := json.MarshalIndent(v, "", "  ")
	return b
}

func summarise(userToStats map[string]*statistics.Statistics) (s summary) {
	s.ReposTouched = setof.Strings()
	for user, stats := range userToStats {
		s.Issues += stats.Issues.Sum()
		if s.IssuesTop.Count < stats.Issues.Sum() {
			s.IssuesTop = userCountOf(user, stats.Issues.Sum())
		}
		s.PullRequestsCreated += stats.PullRequestsCreated.Sum()
		if s.PullRequestsCreatedTop.Count < stats.PullRequestsCreated.Sum() {
			s.PullRequestsCreatedTop = userCountOf(user, stats.PullRequestsCreated.Sum())
		}
		s.PullRequestsMerged += stats.PullRequestsMerged.Sum()
		if s.PullRequestsMergedTop.Count < stats.PullRequestsMerged.Sum() {
			s.PullRequestsMergedTop = userCountOf(user, stats.PullRequestsMerged.Sum())
		}
		s.ReposUpdated += stats.ReposUpdated.Sum()
		if s.ReposUpdatedTop.Count < stats.ReposUpdated.Sum() {
			s.ReposUpdatedTop = userCountOf(user, stats.ReposUpdated.Sum())
		}
		s.Repos += stats.Repos
		if s.ReposTop.Count < stats.Repos {
			s.ReposTop = userCountOf(user, stats.Repos)
		}
		s.Stars += stats.Stars
		if s.StarsTop.Count < stats.Stars {
			s.StarsTop = userCountOf(user, stats.Stars)
		}
		for _, repo := range stats.ReposTouched.Values() {
			s.ReposTouched.Add(repo)
		}
	}
	return s
}

type summary struct {
	Issues                 int              `json:"issues"`
	IssuesTop              userCount        `json:"issuesTop"`
	PullRequestsCreated    int              `json:"pullRequestsCreated"`
	PullRequestsCreatedTop userCount        `json:"pullRequestsCreatedTop"`
	PullRequestsMerged     int              `json:"pullRequestsMerged"`
	PullRequestsMergedTop  userCount        `json:"pullRequestsMergedTop"`
	ReposUpdated           int              `json:"reposUpdated"`
	ReposUpdatedTop        userCount        `json:"reposUpdatedTop"`
	Repos                  int              `json:"repos"`
	ReposTop               userCount        `json:"reposTop"`
	Stars                  int              `json:"stars"`
	StarsTop               userCount        `json:"starsTop"`
	ReposTouched           *setof.StringSet `json:"reposTouched"`
}

func userCountOf(user string, count int) userCount {
	return userCount{
		User:  user,
		Count: count,
	}
}

type userCount struct {
	User  string `json:"user"`
	Count int    `json:"count"`
}
