package entity

import (
	"strconv"
	"time"
)

type PageInfo struct {
	EndCursor string `json:"endCursor"`
}

type SearchRepositories struct {
	Data struct {
		Search struct {
			PageInfo     PageInfo      `json:"pageInfo"`
			Repositories []*Repository `json:"nodes"`
		} `json:"search"`
	} `json:"data"`
}

type Repository struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	PrimaryLanguage struct {
		Name string `json:"name"`
	} `json:"primaryLanguage"`

	Releases struct {
		TotalCount int `json:"totalCount"`
	} `json:"releases"`

	PullRequests struct {
		TotalCount int `json:"totalCount"`
	} `json:"pullRequests"`

	Issues struct {
		TotalCount int `json:"totalCount"`
	} `json:"issues"`
}

func (r *Repository) CsvHeader() []string {
	return []string{"Name", "CreatedAt", "UpdatedAt", "PrimaryLanguage", "Releases", "PullRequests", "Issues"}
}

func (r *Repository) CsvValues() []string {
	return []string{r.Name, r.CreatedAt.String(), r.UpdatedAt.String(), r.PrimaryLanguage.Name, strconv.Itoa(r.Releases.TotalCount), strconv.Itoa(r.PullRequests.TotalCount), strconv.Itoa(r.Issues.TotalCount)}
}
