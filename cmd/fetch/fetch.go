package main

import (
	"fmt"
	"log"
	"time"

	"blaus/pkg/http"
	"blaus/pkg/lib/csv"
)

const (
	QuantityToFetchRepositories = 1000
)

func main() {
	requester, err := http.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	repositories, err := requester.QueryRepos(QuantityToFetchRepositories)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Took %s to query %d repositories\n", time.Since(now), len(repositories))

	now = time.Now()
	err = csv.Save(repositories)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Took %s to csv %d repositories\n", time.Since(now), len(repositories))
}
