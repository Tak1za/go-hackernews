package hn

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiBase = "https://hacker-news.firebaseio.com/v0"
)

type Item struct {
	Deleted bool   `json:"deleted"`
	Type    string `json:"type"`
	By      string `json:"by"`
	Time    int    `json:"time"`
	Text    string `json:"text"`
	Dead    bool   `json:"dead"`
	URL     string `json:"url"`
	Score   int    `json:"score"`
	Title   string `json:"title"`
}

func TopItems() ([]int, error) {
	resp, err := http.Get(fmt.Sprintf("%s/topstories.json", apiBase))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var ids []int
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&ids)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func GetItem(id int) (Item, error) {
	var item Item

	resp, err := http.Get(fmt.Sprintf("%s/item/%d.json", apiBase, id))
	if err != nil {
		return item, err
	}

	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&item)
	if err != nil {
		return item, err
	}
	return item, nil
}
