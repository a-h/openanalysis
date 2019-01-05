package main

import (
	"encoding/json"
	"fmt"
	"sort"

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
	s = newSummary()
	for user, stats := range userToStats {
		s.Issues += stats.Issues.Sum()
		s.IssuesTop = append(s.IssuesTop, userCountOf(user, stats.Issues.Sum()))
		s.PullRequestsCreated += stats.PullRequestsCreated.Sum()
		s.PullRequestsCreatedTop = append(s.PullRequestsCreatedTop, userCountOf(user, stats.PullRequestsCreated.Sum()))
		s.PullRequestsMerged += stats.PullRequestsMerged.Sum()
		s.PullRequestsMergedTop = append(s.PullRequestsMergedTop, userCountOf(user, stats.PullRequestsMerged.Sum()))
		s.ReposUpdated += stats.ReposUpdated.Sum()
		s.ReposUpdatedTop = append(s.ReposUpdatedTop, userCountOf(user, stats.ReposUpdated.Sum()))
		s.Repos += stats.Repos
		s.ReposTop = append(s.ReposTop, userCountOf(user, stats.Repos))
		s.Stars += stats.Stars
		s.StarsTop = append(s.StarsTop, userCountOf(user, stats.Stars))
		for _, repo := range stats.ReposTouched.Values() {
			s.ReposTouched.Add(repo)
		}
	}

	sort.Sort(sort.Reverse(s.IssuesTop))
	s.IssuesTop = s.IssuesTop[0:10]
	sort.Sort(sort.Reverse(s.PullRequestsCreatedTop))
	s.PullRequestsCreatedTop = s.PullRequestsCreatedTop[0:10]
	sort.Sort(sort.Reverse(s.PullRequestsMergedTop))
	s.PullRequestsMergedTop = s.PullRequestsMergedTop[0:10]
	sort.Sort(sort.Reverse(s.ReposUpdatedTop))
	s.ReposUpdatedTop = s.ReposUpdatedTop[0:10]
	sort.Sort(sort.Reverse(s.ReposTop))
	s.ReposTop = s.ReposTop[0:10]
	sort.Sort(sort.Reverse(s.StarsTop))
	s.StarsTop = s.StarsTop[0:10]

	return s
}

func newSummary() summary {
	return summary{
		IssuesTop:              make(userCounts, 0),
		PullRequestsCreatedTop: make(userCounts, 0),
		PullRequestsMergedTop:  make(userCounts, 0),
		ReposUpdatedTop:        make(userCounts, 0),
		ReposTop:               make(userCounts, 0),
		StarsTop:               make(userCounts, 0),
		ReposTouched:           setof.Strings(),
	}
}

type summary struct {
	Issues                 int              `json:"issues"`
	IssuesTop              userCounts       `json:"issuesTop"`
	PullRequestsCreated    int              `json:"pullRequestsCreated"`
	PullRequestsCreatedTop userCounts       `json:"pullRequestsCreatedTop"`
	PullRequestsMerged     int              `json:"pullRequestsMerged"`
	PullRequestsMergedTop  userCounts       `json:"pullRequestsMergedTop"`
	ReposUpdated           int              `json:"reposUpdated"`
	ReposUpdatedTop        userCounts       `json:"reposUpdatedTop"`
	Repos                  int              `json:"repos"`
	ReposTop               userCounts       `json:"reposTop"`
	Stars                  int              `json:"stars"`
	StarsTop               userCounts       `json:"starsTop"`
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

type userCounts []userCount

func (uc userCounts) Len() int           { return len(uc) }
func (uc userCounts) Swap(i, j int)      { uc[i], uc[j] = uc[j], uc[i] }
func (uc userCounts) Less(i, j int) bool { return uc[i].Count < uc[j].Count }
