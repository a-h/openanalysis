package statistics

import (
	"time"

	"github.com/a-h/setof"
)

type Statistic struct {
	Start  time.Time
	End    time.Time
	Values []int
}

func (s *Statistic) Sum() (op int) {
	for _, v := range s.Values {
		op += v
	}
	return
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
