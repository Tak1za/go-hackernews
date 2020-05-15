package main

import (
	"sort"

	"github.com/Tak1za/go-hackernews/hn"
)

func topStories(numStories int) ([]item, error) {
	var stories []item

	itemIds, err := hn.TopItems()
	if err != nil {
		return stories, err
	}

	type Result struct {
		item  item
		err   error
		index int
	}
	resultCh := make(chan Result)
	for i := 0; i < numStories; i++ {
		go func(index, id int) {
			item, err := hn.GetItem(id)
			if err != nil {
				resultCh <- Result{err: err, index: index}
			}
			resultCh <- Result{item: parseHNItem(item), index: index}
		}(i, itemIds[i])
	}

	var results []Result
	for i := 0; i < numStories; i++ {
		results = append(results, <-resultCh)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].index < results[j].index
	})

	for _, res := range results {
		if res.err != nil {
			continue
		}

		if isStoryLink(res.item) {
			stories = append(stories, res.item)
		}
	}

	return stories, nil
}
