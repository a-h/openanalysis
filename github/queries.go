package github

import "time"

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
				CreatedAt  time.Time `json:"createdAt"`
				MergedAt   time.Time `json:"mergedAt"`
				Repository struct {
					NameWithOwner string `json:"nameWithOwner"`
					IsPrivate     string `json:"isPrivate"`
					Description   string `json:"description"`
				} `json:"repository"`
			} `json:"nodes"`
		} `json:"pullRequests"`
	} `json:"user"`
}

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
