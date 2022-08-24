package queries

import (
	"encoding/json"
)

type QueryOptions func(*graphqlRequest)
type graphqlRequest struct {
	Query     string            `json:"query"`
	Variables map[string]string `json:"variables"`
}

func newRequest(query string) *graphqlRequest {
	return &graphqlRequest{
		Query:     query,
		Variables: make(map[string]string),
	}
}

func WithAfter(after string) func(*graphqlRequest) {
	return func(gql *graphqlRequest) {
		if len(after) == 0 {
			return
		}

		gql.Variables["after"] = after
	}
}

func Repositories(options ...QueryOptions) ([]byte, error) {
	query := `
	query search($after: String) {
		search(
			query: "is:public sort:stars-desc stars:>10000",
			type: REPOSITORY,
			first: 20,
      		after: $after
		) {
			nodes {
				... on Repository {
					name
          			createdAt
					updatedAt
          			pushedAt
          			forkCount

					primaryLanguage {
						name
					}
          
					releases(first: 1, orderBy: {field: CREATED_AT, direction: DESC}) {
						totalCount
						nodes {
							createdAt
						}
					}
          
					pullRequests(states: [MERGED]) {
						totalCount
					}
          
					issues {
						totalCount
					}
				}
			}
		}
	}`

	gqlRequest := newRequest(query)
	for _, option := range options {
		option(gqlRequest)
	}

	return json.Marshal(gqlRequest)
}
