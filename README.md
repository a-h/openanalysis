# Github Analysis

Running the `main.go produces JSON output for each user in the `users.json` file.

```json
{
  "Issues": {
    "Start": "2018-01-01T00:00:00Z",
    "End": "2019-01-01T00:00:00Z",
    "Values": [
      0,
      1,
      2,
      1,
      0,
      0,
      0,
      0,
      1,
      0,
      1,
      1
    ]
  },
  "PullRequestsCreated": {
    "Start": "2018-01-01T00:00:00Z",
    "End": "2019-01-01T00:00:00Z",
    "Values": [
      1,
      0,
      0,
      0,
      1,
      0,
      0,
      0,
      0,
      0,
      1,
      0
    ]
  },
  "PullRequestsMerged": {
    "Start": "2018-01-01T00:00:00Z",
    "End": "2019-01-01T00:00:00Z",
    "Values": [
      0,
      0,
      1,
      0,
      0,
      0,
      0,
      0,
      0,
      0,
      0,
      1
    ]
  },
  "ReposUpdated": {
    "Start": "2018-01-01T00:00:00Z",
    "End": "2019-01-01T00:00:00Z",
    "Values": [
      1,
      1,
      2,
      1,
      0,
      2,
      0,
      2,
      2,
      4,
      4,
      2
    ]
  },
  "Repos": 21,
  "Stars": 162,
  "ReposTouched": [
    "aws/aws-sdk-go",
    "hashicorp/terraform",
    "welldigital/serverless-build",
    "golang/go",
    "sirupsen/logrus",
    "nccgroup/ScoutSuite",
    "serverless/serverless",
    "zyedidia/micro",
    "golang/tools",
    "golang-migrate/migrate",
    "Microsoft/vscode-go",
    "a-h/ansible-mongodb-cluster-aws",
    "a-h/generate",
    "a-h/version",
    "a-h/lexical",
    "a-h/Connect",
    "a-h/terraform-example",
    "a-h/ml",
    "a-h/oom",
    "a-h/gauthmiddleware",
    "a-h/pathvars",
    "a-h/scache",
    "a-h/timeid",
    "a-h/jsontogo",
    "a-h/date",
    "a-h/httpdump",
    "a-h/mapof",
    "a-h/setof",
    "a-h/katacoda-scenarios",
    "a-h/search",
    "a-h/nvdnotifier",
    "a-h/watchman"
  ],
  "Start": "2018-01-01T00:00:00Z",
  "End": "2019-01-01T00:00:00Z"
}
```

The contents of the `sum` directory produce a summary of all of the records.

```json
{
  "issues": 71,
  "issuesTop": [
    {
      "user": "xxx",
      "count": 14
    },
    {
      "user": "xxx",
      "count": 7
    }
  ],
  "pullRequestsCreated": 52,
  "pullRequestsCreatedTop": [
    {
      "user": "xxx",
      "count": 6
    },
    {
      "user": "xxx",
      "count": 6
    }
  ],
  "pullRequestsMerged": 154,
  "pullRequestsMergedTop": [
    {
      "user": "xxx",
      "count": 19
    },
    {
      "user": "xxx",
      "count": 18
    }
  ],
  "reposUpdated": 274,
  "reposUpdatedTop": [
    {
      "user": "xxx",
      "count": 29
    },
    {
      "user": "a-h",
      "count": 21
    }
  ],
  "repos": 274,
  "reposTop": [
    {
      "user": "a-h",
      "count": 29
    },
    {
      "user": "xxx",
      "count": 21
    }
  ],
  "stars": 700,
  "starsTop": [
    {
      "user": "xxx",
      "count": 467
    },
    {
      "user": "a-h",
      "count": 162
    }
  ]
}
```