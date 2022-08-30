package csv

import (
	"encoding/csv"
	"fmt"
	"os"

	"blaus/pkg/entity"
)

func Save(repositories []*entity.Repository) error {
	if len(repositories) == 0 {
		return nil
	}

	file, err := os.Create("etc/repositories.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)

	headers := repositories[0].CsvHeader()
	values := make([][]string, 0, len(repositories))

	for _, r := range repositories {
		values = append(values, r.CsvValues())
	}

	if err := w.Write(headers); err != nil {
		return err
	}

	if err := w.WriteAll(values); err != nil {
		return err
	}

	fmt.Printf("saved %d rows on csv with success\n", len(repositories))
	return nil
}

func ReadRepositories(filePath string) ([]*entity.Repository, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	repositories := make([]*entity.Repository, 0)

	for _, row := range records[1:] {
		repository := &entity.Repository{}
		repository.FillFromCSV(row)

		repositories = append(repositories, repository)
	}

	return repositories, nil
}
