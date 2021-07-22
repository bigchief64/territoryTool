package csv

import (
	"encoding/csv"
	"log"
	"os"
)

type Setting struct {
	name, sex, phone, address string
}

var Settings map[string] []string
var fileName string

func GetData(filename string)  map[string] []string{
	Settings = make(map[string] []string)

	if _, err := os.Stat(filename); err == nil {
		records, err := readData(filename)

		if err != nil {
			log.Fatal(err)
		}

		for _, record := range records {

			setting := Setting{
				name:  record[0],
				sex: record[1],
				phone: record[2],
				address: record[3],
			}
			//fmt.Printf("%s %s is a %s\n", setting.item, setting.setAs)
			Settings[setting.name] = []string{setting.sex, setting.phone, setting.address}
		}

	} else{
			log.Fatal(err)
	}
	return Settings
}

func readData(fileName string) ([][]string, error) {
	f, err := os.Open(fileName)

	if err != nil {
		return [][]string{}, err
	}

	defer f.Close()

	r := csv.NewReader(f)

	// skip first line
	if _, err := r.Read(); err != nil {
		return [][]string{}, err
	}

	records, err := r.ReadAll()

	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}

