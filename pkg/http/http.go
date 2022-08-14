package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"blaus/pkg/config"
	"blaus/pkg/entity"
	"blaus/pkg/queries"
)

type Client interface {
	QueryRepos(limit int) ([]*entity.Repository, error)
}

type requester struct {
	client   *http.Client
	endpoint string
	token    string
}

func NewClient() (Client, error) {
	token := config.GithubToken()
	if len(token) == 0 {
		return nil, errors.New("empty github token")
	}

	return requester{
		client:   &http.Client{},
		endpoint: "https://api.github.com/graphql",
		token:    fmt.Sprintf("Bearer %s", token),
	}, nil
}

func (r requester) post(body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", r.endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", r.token)
	req.Header.Add("Accept", "application/json")

	return r.client.Do(req)
}

func (r requester) QueryRepos(limit int) ([]*entity.Repository, error) {
	after := ""
	repositories := make([]*entity.Repository, 0, limit)

	for len(repositories) < limit {
		query, err := queries.Repositories(queries.WithAfter(after))
		if err != nil {
			return nil, err
		}

		res, err := r.post(bytes.NewBuffer(query))
		if err != nil {
			return nil, err
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		parse := &entity.SearchRepositories{}
		err = json.Unmarshal(body, &parse)
		if err != nil {
			return nil, err
		}

		repositories = append(repositories, parse.Data.Search.Repositories...)
		after = parse.Data.Search.PageInfo.EndCursor
		fmt.Printf("Collected %d repositories\n", len(repositories))
	}

	return repositories, nil
}
