package lib

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

func JsonFile[T interface{}](file string) T {
	body, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("Failed to read file %s: %s", file, err)
	}

	var data T
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal("Failed to parse JSON: ", err)
	}

	return data
}

func FetchJson[T interface{}](loc string) T {
	res, err := http.Get(loc)
	if err != nil {
		log.Fatalf("Failed to fetch %s: %s", loc, err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Failed to fetch %s: status code %d", loc, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Failed to read response body: ", err)
	}

	var data T
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal("Failed to parse JSON: ", err)
	}

	return data
}
