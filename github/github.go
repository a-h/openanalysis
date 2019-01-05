package github

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/machinebox/graphql"

	"golang.org/x/oauth2"
)

// Collector uses the Github GraphQL API to collect data.
type Collector struct {
	githubToken string
	Log         func(s string)
}

// maximumPageSize of Github GraphQL requests.
const maximumPageSize = 100

// NewCollector creates a Github data collector with the githubToken used to authenticate API calls.
// See https://developer.github.com/v4/guides/forming-calls/#authenticating-with-graphql
func NewCollector(githubToken string) *Collector {
	return &Collector{
		githubToken: githubToken,
	}
}

func (c *Collector) client(ctx context.Context) *graphql.Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.githubToken},
	)
	httpClient := oauth2.NewClient(ctx, src)
	client := graphql.NewClient("https://api.github.com/graphql",
		graphql.WithHTTPClient(httpClient))
	if c.Log != nil {
		client.Log = c.Log
	}
	return client
}

// UserIssue within Github.
type UserIssue struct {
	Login             string    `json:"login"`
	RepoNameWithOwner string    `json:"repoNameWithOwner"`
	RepoDescription   string    `json:"repoDescription"`
	RepoIsPrivate     bool      `json:"repoIsPrivate"`
	CreatedAt         time.Time `json:"createdAt"`
}

// NewUserIssue creates an issue with all fields populated.
func NewUserIssue(login, repoNameWithOwner, repoDescription string, repoIsPrivate bool, createdAt time.Time) UserIssue {
	return UserIssue{
		Login:             login,
		RepoNameWithOwner: repoNameWithOwner,
		RepoDescription:   repoDescription,
		RepoIsPrivate:     repoIsPrivate,
		CreatedAt:         createdAt,
	}
}

// UserIssues returns all issues for a given user, by calling the same GraphQL query in a loop for each of the
// issues.
func (c Collector) UserIssues(ctx context.Context, login string) (issues []UserIssue, err error) {
	var cursor *string
	for {
		res, issErr := c.userIssuesPage(ctx, login, maximumPageSize, cursor)
		if issErr != nil {
			err = fmt.Errorf("collector: failed to get issues for user '%s': %v", login, issErr)
			return
		}
		for _, n := range res.User.Issues.Nodes {
			// Skip self-owned repos.
			if strings.HasPrefix(n.Repository.NameWithOwner, login) {
				continue
			}
			issues = append(issues, NewUserIssue(login, n.Repository.NameWithOwner, n.Repository.Description, n.Repository.IsPrivate, n.CreatedAt))
		}
		cursor = &res.User.Issues.PageInfo.EndCursor
		if !res.User.Issues.PageInfo.HasNextPage {
			return
		}
	}
}

func (c Collector) userIssuesPage(ctx context.Context, login string, first int, cursor *string) (result userIssuesQueryResult, err error) {
	req := graphql.NewRequest(userIssuesQuery)

	req.Var("login", login)
	req.Var("first", first)
	req.Var("cursor", cursor)

	req.Header.Set("Cache-Control", "no-cache")

	err = c.client(ctx).Run(ctx, req, &result)
	return
}

const userIssuesQuery = `query ($login: String!, $first: Int!, $cursor: String) {
  user(login: $login) {
    issues(first: $first, after: $cursor) {
      pageInfo {
        endCursor
        hasNextPage
      }
      nodes {
        createdAt
        repository {
          nameWithOwner
          isPrivate
          description
        }
      }
    }
  }
}`

type userIssuesQueryResult struct {
	User struct {
		Issues struct {
			PageInfo struct {
				EndCursor   string `json:"endCursor"`
				HasNextPage bool   `json:"hasNextPage"`
			} `json:"pageInfo"`
			Nodes []struct {
				CreatedAt  time.Time `json:"createdAt"`
				Repository struct {
					NameWithOwner string `json:"nameWithOwner"`
					IsPrivate     bool   `json:"isPrivate"`
					Description   string `json:"description"`
				} `json:"repository"`
			} `json:"nodes"`
		} `json:"issues"`
	} `json:"user"`
}

// UserPullRequest within Github.
type UserPullRequest struct {
	Login             string     `json:"login"`
	RepoNameWithOwner string     `json:"repoNameWithOwner"`
	RepoDescription   string     `json:"repoDescription"`
	RepoIsPrivate     bool       `json:"repoIsPrivate"`
	CreatedAt         time.Time  `json:"createdAt"`
	MergedAt          *time.Time `json:"mergedAt"`
}

// NewUserPullRequest creates an issue with all fields populated.
func NewUserPullRequest(login, repoNameWithOwner, repoDescription string, repoIsPrivate bool, createdAt time.Time, mergedAt *time.Time) UserPullRequest {
	return UserPullRequest{
		Login:             login,
		RepoNameWithOwner: repoNameWithOwner,
		RepoDescription:   repoDescription,
		RepoIsPrivate:     repoIsPrivate,
		CreatedAt:         createdAt,
		MergedAt:          mergedAt,
	}
}

// UserPullRequests returns all issues for a given user, by calling the same GraphQL query in a loop for each of the
// issues.
func (c Collector) UserPullRequests(ctx context.Context, login string) (pullrequests []UserPullRequest, err error) {
	var cursor *string
	for {
		res, issErr := c.userPullRequestsPage(ctx, login, maximumPageSize, cursor)
		if issErr != nil {
			err = fmt.Errorf("collector: failed to get issues for user '%s': %v", login, issErr)
			return
		}
		for _, n := range res.User.PullRequests.Nodes {
			pullrequests = append(pullrequests, NewUserPullRequest(login, n.Repository.NameWithOwner, n.Repository.Description, n.Repository.IsPrivate, n.CreatedAt, n.MergedAt))
		}
		cursor = &res.User.PullRequests.PageInfo.EndCursor
		if !res.User.PullRequests.PageInfo.HasNextPage {
			return
		}
	}
}

func (c Collector) userPullRequestsPage(ctx context.Context, login string, first int, cursor *string) (result userPullRequestsQueryResult, err error) {
	req := graphql.NewRequest(userPullRequestsQuery)

	req.Var("login", login)
	req.Var("first", first)
	req.Var("cursor", cursor)

	req.Header.Set("Cache-Control", "no-cache")

	err = c.client(ctx).Run(ctx, req, &result)
	return
}

const userPullRequestsQuery = `query ($login: String!, $first: Int!, $cursor: String) {
  user(login: $login) {
    pullRequests(first: $first, after: $cursor) {
      pageInfo {
        endCursor
        hasNextPage
      }
      nodes {
        createdAt
        mergedAt
        repository {
          nameWithOwner
          isPrivate
          description
        }
      }
    }
  }
}`

type userPullRequestsQueryResult struct {
	User struct {
		PullRequests struct {
			PageInfo struct {
				EndCursor   string `json:"endCursor"`
				HasNextPage bool   `json:"hasNextPage"`
			} `json:"pageInfo"`
			Nodes []struct {
				CreatedAt  time.Time  `json:"createdAt"`
				MergedAt   *time.Time `json:"mergedAt"`
				Repository struct {
					NameWithOwner string `json:"nameWithOwner"`
					IsPrivate     bool   `json:"isPrivate"`
					Description   string `json:"description"`
				} `json:"repository"`
			} `json:"nodes"`
		} `json:"pullRequests"`
	} `json:"user"`
}

// UserRepository within Github.
type UserRepository struct {
	Login             string    `json:"login"`
	RepoNameWithOwner string    `json:"repoNameWithOwner"`
	RepoDescription   string    `json:"repoDescription"`
	RepoIsPrivate     bool      `json:"repoIsPrivate"`
	Stars             int       `json:"stars"`
	PushedAt          time.Time `json:"pushedAt"`
}

// NewUserRepository creates an issue with all fields populated.
func NewUserRepository(login, repoNameWithOwner, repoDescription string, repoIsPrivate bool, stars int, pushedAt time.Time) UserRepository {
	return UserRepository{
		Login:             login,
		RepoNameWithOwner: repoNameWithOwner,
		RepoDescription:   repoDescription,
		Stars:             stars,
		PushedAt:          pushedAt,
	}
}

// UserRepositories returns all issues for a given user, by calling the same GraphQL query in a loop for each of the
// issues.
func (c Collector) UserRepositories(ctx context.Context, login string) (repos []UserRepository, err error) {
	var cursor *string
	for {
		res, issErr := c.userRepositoriesPage(ctx, login, maximumPageSize, cursor)
		if issErr != nil {
			err = fmt.Errorf("collector: failed to get issues for user '%s': %v", login, issErr)
			return
		}
		for _, n := range res.User.Repositories.Nodes {
			repos = append(repos, NewUserRepository(login, n.NameWithOwner, n.Description, false, n.Stargazers.TotalCount, n.PushedAt))
		}
		cursor = &res.User.Repositories.PageInfo.EndCursor
		if !res.User.Repositories.PageInfo.HasNextPage {
			return
		}
	}
}

func (c Collector) userRepositoriesPage(ctx context.Context, login string, first int, cursor *string) (result userRepositoriesQueryResult, err error) {
	req := graphql.NewRequest(userRepositoriesQuery)

	req.Var("login", login)
	req.Var("first", first)
	req.Var("cursor", cursor)

	req.Header.Set("Cache-Control", "no-cache")

	err = c.client(ctx).Run(ctx, req, &result)
	return
}

const userRepositoriesQuery = `query ($login: String!, $first: Int!, $cursor: String) {
  user(login: $login) {
    repositories(privacy:PUBLIC, isFork:false, ownerAffiliations:OWNER, first:$first, after: $cursor) {
      pageInfo {
        endCursor
        hasNextPage
      }
      nodes {
        nameWithOwner
        description
        stargazers {
          totalCount
        }
        pushedAt
      }
    }
  }
}`

type userRepositoriesQueryResult struct {
	User struct {
		Repositories struct {
			PageInfo struct {
				EndCursor   string `json:"endCursor"`
				HasNextPage bool   `json:"hasNextPage"`
			} `json:"pageInfo"`
			Nodes []struct {
				NameWithOwner string `json:"nameWithOwner"`
				Description   string `json:"description"`
				Stargazers    struct {
					TotalCount int `json:"totalCount"`
				} `json:"stargazers"`
				PushedAt time.Time `json:"pushedAt"`
			} `json:"nodes"`
		} `json:"repositories"`
	} `json:"user"`
}
