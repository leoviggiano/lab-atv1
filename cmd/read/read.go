package main

import (
	"fmt"
	"log"
	"sort"
	"time"

	"blaus/pkg/entity"
	"blaus/pkg/lib/csv"
)

const (
	// Métrica: idade do repositório (calculado a partir da data de sua criação)
	RQ01 = "Sistemas populares são maduros/antigos?"

	// Métrica: total de pull requests aceitas
	RQ02 = "Sistemas populares recebem muita contribuição externa?"

	// Métrica: total de releases
	RQ03 = "Sistemas populares lançam releases com frequência?"

	// Métrica: tempo até a última atualização (calculado a partir da data de última atualização)
	RQ04 = "Sistemas populares são atualizados com frequência?"

	// Métrica: linguagem primária de cada um desses repositórios
	RQ05 = "Sistemas populares são escritos nas linguagens mais populares"

	// Métrica: razão entre número de issues fechadas pelo total de issues
	RQ06 = "Sistemas populares possuem um alto percentual de issues fechadas?"

	// compare os resultados para os sistemas com as linguagens da reportagem com os resultados de sistemas em outras linguagens.
	RQ07 = "Sistemas escritos em linguagens mais populares recebem mais contribuição externa, lançam mais releases e são atualizados com mais frequência?"
)

var answers = map[string]interface{}{}

func main() {
	pathFile := "./etc/repositories.csv"
	repositories, err := csv.ReadRepositories(pathFile)
	if err != nil {
		log.Fatal(err)
	}

	var (
		totalDaysCreatedAt int
		totalDaysUpdatedAt int
		totalPullRequests  int
		totalReleases      int
		totalClosedIssues  int
		totalIssues        int

		topLanguagesUpdatedAt    int
		topLanguagesReleases     int
		topLanguagesPullRequests int
	)

	primaryLanguage := make(map[string]int, 0)
	languageRepositories := make(map[string][]*entity.Repository, 0)

	for _, v := range repositories {
		totalDaysCreatedAt += int(time.Since(v.CreatedAt) / time.Hour / 24)
		totalDaysUpdatedAt += int(time.Since(v.UpdatedAt) / time.Hour / 24)
		totalPullRequests += v.PullRequests.TotalCount
		totalReleases += v.Releases.TotalCount
		totalClosedIssues += v.Issues.Closed
		totalIssues += v.Issues.TotalCount

		if v.PrimaryLanguage.Name != "" {
			name := v.PrimaryLanguage.Name
			primaryLanguage[name] += 1
			languageRepositories[name] = append(languageRepositories[name], v)
		}
	}

	countUseLanguages := make([]int, 0, len(primaryLanguage))

	for _, v := range primaryLanguage {
		countUseLanguages = append(countUseLanguages, v)
	}

	sort.SliceStable(countUseLanguages, func(i, j int) bool {
		return countUseLanguages[i] > countUseLanguages[j]
	})

	topLanguages := make([]string, 0, 5)
	for _, v := range countUseLanguages {
		for k := range primaryLanguage {
			if primaryLanguage[k] == v {
				topLanguages = append(topLanguages, k)
				continue
			}
		}

		if len(topLanguages) == 5 {
			break
		}
	}

	for _, language := range topLanguages {
		for _, v := range languageRepositories[language] {
			topLanguagesReleases += v.Releases.TotalCount
			topLanguagesUpdatedAt += int(time.Since(v.UpdatedAt) / time.Hour / 24)
			topLanguagesPullRequests += v.PullRequests.TotalCount
		}
	}

	answers[RQ01] = fmt.Sprintf("Em média, os repositórios tem %d dias desde a criação", totalDaysCreatedAt/len(repositories))
	answers[RQ02] = fmt.Sprintf("Em média, os repositórios tem um total de %d de pull requests aceitas", totalPullRequests/len(repositories))
	answers[RQ03] = fmt.Sprintf("Em média, os repositórios tem %d releases", totalReleases/len(repositories))
	answers[RQ04] = fmt.Sprintf("Em média, os repositórios tem %d dias desde a última atualização", totalDaysUpdatedAt/len(repositories))
	answers[RQ05] = fmt.Sprintf("As linguagens mais utilizadas nos repositórios, foram: %v", topLanguages)
	answers[RQ06] = fmt.Sprintf("Em média, os repositórios possuem um total de %.f%% de issues fechadas", float64(totalClosedIssues)/float64(totalIssues)*100)
	answers[RQ07] = fmt.Sprintf("Em média, os repositórios com as 5 linguagens mais populares, possuem um total de:\n"+
		"%d pull requests aceitas\n"+
		"%d releases\n"+
		"%d dias desde a última atualização\n",
		topLanguagesPullRequests/len(topLanguages),
		topLanguagesReleases/len(topLanguages),
		topLanguagesUpdatedAt/len(topLanguages),
	)

	for k, v := range answers {
		fmt.Printf("%s\n%v\n\n", k, v)
	}
}
