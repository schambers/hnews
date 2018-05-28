package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

const baseUrl = "https://hacker-news.firebaseio.com/v0/"
const ts = baseUrl + "topstories.json"

type Article struct {
	Id int64 `json:"id,omitempty"`
	Author string `json:"by,omitempty"`
	Title string `json:"title,omitempty"`
	Score int `json:"score,omitempty"`
	Url string `json:"url,omitempty"`
}

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
		go getArticle(itemId, ch)
	}

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
	var article Article
	if err := json.Unmarshal(body, &article); err != nil {
		panic(err)
	}

	ch <- fmt.Sprintf("Id: %d\nTitle: %s\nAuthor: %s\nScore:%d\nUrl:%s\n\n", article.Id, article.Title, article.Author, article.Score, article.Url)
}
