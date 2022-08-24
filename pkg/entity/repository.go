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

type NodeRepository struct {
	Data struct {
		Node struct {
			Issues struct {
				TotalCount int `json:"totalCount"`
			} `json:"issues"`
		} `json:"node"`
	} `json:"data"`
}

type Repository struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	PrimaryLanguage struct {
		Name string `json:"name"`
	} `json:"primaryLanguage"`

	Releases struct {
		TotalCount int        `json:"totalCount"`
		Nodes      []*Release `json:"nodes"`
	} `json:"releases"`

	PullRequests struct {
		TotalCount int `json:"totalCount"`
	} `json:"pullRequests"`

	Issues struct {
		TotalCount int `json:"totalCount"`
		Closed     int `json:"-"`
		Open       int `json:"-"`
	} `json:"issues"`
}

type Release struct {
	CreatedAt time.Time `json:"createdAt"`
}

func (r *Repository) CsvHeader() []string {
	return []string{"Name", "CreatedAt", "UpdatedAt", "PrimaryLanguage", "Releases", "PullRequests", "Open Issues", "Closed Issues", "LastRelease"}
}

func (r *Repository) CsvValues() []string {
	lastRelease := "No releases"
	if len(r.Releases.Nodes) > 0 {
		lastRelease = r.Releases.Nodes[0].CreatedAt.Format(time.RFC3339)
	}

	return []string{
		r.Name,
		r.CreatedAt.String(),
		r.UpdatedAt.String(),
		r.PrimaryLanguage.Name,
		strconv.Itoa(r.Releases.TotalCount),
		strconv.Itoa(r.PullRequests.TotalCount),
		strconv.Itoa(r.Issues.Open),
		strconv.Itoa(r.Issues.Closed),
		lastRelease,
	}
}
