package entity

import (
	"strconv"
	"time"
)

const NoRelease = "No Releases"

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

type PrimaryLanguage struct {
	Name string `json:"name"`
}

type Releases struct {
	TotalCount    int        `json:"totalCount"`
	LatestRelease time.Time  `json:"-"`
	Nodes         []*Release `json:"nodes"`
}

type Repository struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	PrimaryLanguage PrimaryLanguage `json:"primaryLanguage"`
	Releases        Releases        `json:"releases"`

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
	lastRelease := NoRelease
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

func (r *Repository) FillFromCSV(row []string) {
	r.Name = row[0]
	r.CreatedAt, _ = time.Parse("2006-01-02 15:04:05 -0700 MST", row[1])
	r.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05 -0700 MST", row[2])
	r.PrimaryLanguage.Name = row[3]
	r.Releases.TotalCount, _ = strconv.Atoi(row[4])
	r.PullRequests.TotalCount, _ = strconv.Atoi(row[5])
	r.Issues.Open, _ = strconv.Atoi(row[6])
	r.Issues.Closed, _ = strconv.Atoi(row[7])
	r.Issues.TotalCount = r.Issues.Open + r.Issues.Closed

	if row[8] != NoRelease {
		date, _ := time.Parse(time.RFC3339, row[8])
		r.Releases.LatestRelease = date
	}
}
