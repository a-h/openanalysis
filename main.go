package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/a-h/openanalysis/github"
)

var token = flag.String("token", "", "GitHub auth token")
var login = flag.String("login", "a-h", "GitHub login")

func main() {
	flag.Parse()
	c := github.NewCollector(*token)
	issues, err := c.UserIssues(context.Background(), *login)
	if err != nil {
		fmt.Printf("Error getting issue: %v\n", err)
		return
	}
	//TODO: Filter by date range, then group by month.
	for _, issue := range issues {
		fmt.Printf("%+v\n", issue)
	}
}
