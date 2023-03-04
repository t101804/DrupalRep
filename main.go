// @CallMeRep
//
//	Join @DailyToolz for more
package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

func WriteFile(s string, p string) {
	escapeFiles := url.QueryEscape(p)
	if _, err := os.Stat("result"); os.IsNotExist(err) {
		err = os.Mkdir("result", 0755)
		if err != nil {
			fmt.Println("cant make folder results")
		}
	}

	f, err := os.OpenFile("result/"+escapeFiles+".txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("cant open file results/" + escapeFiles + ".txt")
	}
	defer f.Close()

	_, err = fmt.Fprint(f, s+"\n")
	if err != nil {
		fmt.Println("cant write file to results/" + escapeFiles + ".txt")
	}
}
func main() {
	var lists string
	fmt.Print("list : ")
	fmt.Scanln(&lists)
	file, err := os.Open(lists)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	resultsFile, err := os.Create("vuln-rce.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resultsFile.Close()

	scanner := bufio.NewScanner(file)

	var wg sync.WaitGroup
	for scanner.Scan() {
		url := scanner.Text()
		if !strings.HasPrefix(url, "http") {
			url = "http://" + url
		}
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			url = url + "/CHANGELOG.txt"

			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("%s\tError: %s\n", url, err)
				return
			}
			if resp.StatusCode == http.StatusOK {
				bodyBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					return
				}
				body := string(bodyBytes)
				if strings.Contains(body, "Drupal 1.0.0") {
					WriteFile(url, "drupal.txt")
					if strings.Contains(body, "(remote code execution).") {
						fmt.Println("not vuln " + url)
					} else {
						fmt.Println("vuln " + url)
						WriteFile(url, "rce.txt")
					}
				}

			}
			resp.Body.Close()
		}(url)
	}
	wg.Wait()

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}
