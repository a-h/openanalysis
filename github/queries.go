package github

import "time"

const userPublicRepositories = `query ($login: String!, $first: Int!, $cursor: String) {
  user(login: $login) {
    repositories(privacy:PUBLIC, first:$first, after: $cursor) {
      pageInfo {
        endCursor
        hasNextPage
      }
      nodes {
        nameWithOwner
        stargazers {
          totalCount
        }
        updatedAt
      }
    }
  }
}`

type userPublicRepositoriesQueryResult struct {
	User struct {
		PullRequests struct {
			PageInfo struct {
				EndCursor   string `json:"endCursor"`
				HasNextPage bool   `json:"hasNextPage"`
			} `json:"pageInfo"`
			Nodes []struct {
				NameWithOwner string `json:"nameWithOwner"`
				Stargazers    struct {
					TotalCount int `json:"totalCount"`
				} `json:"stargazers"`
				UpdatedAt time.Time `json:"updatedAt"`
			} `json:"nodes"`
		} `json:"pullRequests"`
	} `json:"user"`
}
