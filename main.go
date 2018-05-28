package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

const baseUrl = "https://hacker-news.firebaseio.com/v0/"
const ts = baseUrl + "topstories.json"

func main() {

	resp, err := http.Get(ts)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var itemIds []int64
	body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &itemIds)

	ch := make(chan string)
	for _, itemId := range itemIds[:10] {
		fmt.Println(itemId)
		go getArticle(itemId, ch)
	}

	fmt.Printf("Total number of articles:%d\n", len(itemIds))

	for range itemIds[:10] {
		fmt.Println(<-ch)
	}
}

func getArticle(articleId int64, ch chan<-string) {
	var itemBaseUrl = fmt.Sprintf("%sitem/%d.json", baseUrl, articleId)
	resp, err := http.Get(itemBaseUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	ch <- fmt.Sprintf("Article %d response:%s", articleId, body)
}
