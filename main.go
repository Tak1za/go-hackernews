package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/Tak1za/go-hackernews/hn"
)

type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}

func main() {
	tpl := template.Must(template.ParseFiles("./index.gohtml"))
	http.HandleFunc("/", handler(30, tpl))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":3000"), nil))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		storyItems, err := getTopStories(numStories)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		data := templateData{
			Stories: storyItems,
			Time:    time.Now().Sub(start),
		}

		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

func getTopStories(numStories int) ([]item, error) {
	ids, err := hn.TopItems()
	if err != nil {
		return nil, errors.New("Failed to load top stories")
	}

	var stories []item
	at := 0
	for len(stories) < numStories {
		need := (numStories - len(stories)) * 5 / 4
		stories = append(stories, getStories(ids[at:at+need])...)
		at += need
	}

	return stories[:numStories], nil
}

func getStories(ids []int) []item {
	type result struct {
		item item
		err  error
		idx  int
	}

	resultCh := make(chan result)
	for i := 0; i < len(ids); i++ {
		go func(idx, id int) {
			hnItem, err := hn.GetItem(id)
			if err != nil {
				resultCh <- result{err: err, idx: idx}
			}
			resultCh <- result{item: parseHNItem(hnItem), idx: idx}
		}(i, ids[i])
	}

	var results []result
	for i := 0; i < len(ids); i++ {
		results = append(results, <-resultCh)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].idx < results[j].idx
	})

	var stories []item
	for _, res := range results {
		if res.err != nil {
			continue
		}
		if isStoryLink(res.item) {
			stories = append(stories, res.item)
		}
	}

	return stories
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}
