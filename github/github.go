package github

import (
	"context"
	"fmt"
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

// // A Comment on a Github issue.
// type Comment struct {
// 	Owner       string    `json:"owner"`
// 	Repo        string    `json:"repo"`
// 	IssueNumber int       `json:"issueNumber"`
// 	URL         string    `json:"url"`
// 	BodyText    string    `json:"bodyText"`
// 	UpdatedAt   time.Time `json:"updatedAt"`
// }

// // Hash of the comment.
// func (c Comment) Hash() string {
// 	j, _ := json.Marshal(c)
// 	return fmt.Sprintf("%x2", sha256.Sum256([]byte(j)))
// }

// // NewComment creates a comment with all required fields populated.
// func NewComment(owner, repo string, issueNumber int, url string, bodyText string, updatedAt time.Time) Comment {
// 	return Comment{
// 		Owner:       owner,
// 		Repo:        repo,
// 		IssueNumber: issueNumber,
// 		URL:         url,
// 		BodyText:    bodyText,
// 		UpdatedAt:   updatedAt,
// 	}
// }

// // Comments retrieves all of the comments for a particular issue.
// func (c Collector) Comments(ctx context.Context, owner, repo string, issueNumber int) (comments []Comment, err error) {
// 	var cursor *string
// 	for {
// 		res, comErr := c.commentsPage(ctx, owner, repo, issueNumber, maximumPageSize, cursor)
// 		if comErr != nil {
// 			err = fmt.Errorf("collector: failed to get comments for repo '%s/%s/issues/%d': %v", owner, repo, issueNumber, comErr)
// 			return
// 		}
// 		for _, n := range res.Repository.Issue.Comments.Nodes {
// 			comments = append(comments, NewComment(owner, repo, issueNumber, n.URL, n.BodyText, n.UpdatedAt))
// 		}
// 		cursor = &res.Repository.Issue.Comments.PageInfo.EndCursor
// 		if !res.Repository.Issue.Comments.PageInfo.HasNextPage {
// 			return
// 		}
// 	}
// }

// func (c Collector) commentsPage(ctx context.Context, owner, repo string, issueNumber int, first int, cursor *string) (result commentsQueryResult, err error) {
// 	req := graphql.NewRequest(commentsQuery)

// 	req.Var("owner", owner)
// 	req.Var("repo", repo)
// 	req.Var("issueNumber", issueNumber)
// 	req.Var("first", first)
// 	req.Var("cursor", cursor)

// 	req.Header.Set("Cache-Control", "no-cache")

// 	err = c.client(ctx).Run(ctx, req, &result)
// 	return
// }

// const commentsQuery = `query ($owner: String!, $repo: String!, $issueNumber: Int!, $first: Int!, $cursor: String) {
//   repository(owner: $owner, name: $repo) {
//     issue(number: $issueNumber) {
//       comments(first: $first, after: $cursor) {
//         pageInfo {
//           endCursor
//           hasNextPage
//         }
//         nodes {
//           url
//           updatedAt
//           bodyText
//         }
//       }
//     }
//   }
// }`

// type commentsQueryResult struct {
// 	Repository struct {
// 		Issue struct {
// 			Comments struct {
// 				PageInfo struct {
// 					EndCursor   string `json:"endCursor"`
// 					HasNextPage bool   `json:"hasNextPage"`
// 				} `json:"pageInfo"`
// 				Nodes []struct {
// 					URL       string    `json:"url"`
// 					UpdatedAt time.Time `json:"updatedAt"`
// 					BodyText  string    `json:"bodyText"`
// 				} `json:"nodes"`
// 			} `json:"comments"`
// 		} `json:"issue"`
// 	} `json:"repository"`
// }
