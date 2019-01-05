package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/a-h/openanalysis/github"
	"github.com/a-h/setof"
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
		err = ioutil.WriteFile("output/"+filepath.Clean(login)+".txt", jsonOrNothing(stats), 0640)
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

type Statistic struct {
	Start  time.Time
	End    time.Time
	Values []int
}

func NewStatistic(start, end time.Time) *Statistic {
	start, end = roundDownToMonth(start), roundDownToMonth(end)
	s := &Statistic{
		Start:  start,
		End:    end,
		Values: make([]int, monthsBetween(start, end)),
	}
	return s
}

func (s *Statistic) Add(date time.Time, value int) {
	date = roundDownToMonth(date)
	if date.Before(s.Start) {
		return
	}
	if date.Equal(s.End) || date.After(s.End) {
		return
	}
	s.Values[monthsBetween(s.Start, date)] += value
}

func roundDownToMonth(a time.Time) time.Time {
	return time.Date(a.Year(), a.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func monthsBetween(a, b time.Time) (months int) {
	a, b = roundDownToMonth(a), roundDownToMonth(b)
	if a.After(b) {
		a, b = b, a
	}
	for {
		if a.Equal(b) || a.After(b) {
			return
		}
		a = time.Date(a.Year(), a.Month()+1, 1, 0, 0, 0, 0, time.UTC)
		months++
	}
}

type Statistics struct {
	Issues              *Statistic
	PullRequestsCreated *Statistic
	PullRequestsMerged  *Statistic
	ReposUpdated        *Statistic
	// Repos is how many public repos were updated in the period.
	Repos        int
	Stars        int
	ReposTouched *setof.StringSet
	Start        time.Time
	End          time.Time
}

func (s *Statistics) isWithinDateRange(date time.Time) bool {
	date = roundDownToMonth(date)
	if date.Before(s.Start) {
		return false
	}
	if date.Equal(s.End) || date.After(s.End) {
		return false
	}
	return true
}

func (s *Statistics) IncrementRepo(date time.Time) {
	if !s.isWithinDateRange(date) {
		return
	}
	s.Repos++
}

func (s *Statistics) AddStars(date time.Time, stars int) {
	if !s.isWithinDateRange(date) {
		return
	}
	s.Stars += stars
}

func (s *Statistics) AddTouchedRepo(date time.Time, value string) {
	if !s.isWithinDateRange(date) {
		return
	}
	s.ReposTouched.Add(value)
}

func NewStatistics(start, end time.Time) *Statistics {
	start, end = roundDownToMonth(start), roundDownToMonth(end)
	return &Statistics{
		Start:               start,
		End:                 end,
		Issues:              NewStatistic(start, end),
		PullRequestsCreated: NewStatistic(start, end),
		PullRequestsMerged:  NewStatistic(start, end),
		ReposUpdated:        NewStatistic(start, end),
		ReposTouched:        setof.Strings(),
	}
}

func getStats(c *github.Collector, login string, start, end time.Time) (stats *Statistics, err error) {
	stats = NewStatistics(start, end)

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
		stats.ReposUpdated.Add(r.UpdatedAt, 1)
		stats.AddStars(r.UpdatedAt, r.Stars)
		stats.IncrementRepo(r.UpdatedAt)
		stats.AddTouchedRepo(r.UpdatedAt, r.RepoNameWithOwner)
	}
	return
}
