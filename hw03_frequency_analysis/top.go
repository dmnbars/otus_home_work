package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type wordCount struct {
	word  string
	count int
}

func Top10(text string) []string {
	words := strings.Fields(text)

	wordsCounts := make([]wordCount, 0, len(words))
	for _, word := range words {
		found := false
		for i, w := range wordsCounts {
			if w.word == word {
				wordsCounts[i].count++
				found = true
				break
			}
		}
		if !found {
			wordsCounts = append(wordsCounts, wordCount{word: word, count: 1})
		}
	}

	sort.Slice(wordsCounts, func(i, j int) bool {
		if wordsCounts[i].count == wordsCounts[j].count {
			return wordsCounts[i].word < wordsCounts[j].word
		}
		return wordsCounts[i].count > wordsCounts[j].count
	})

	result := make([]string, 0, 10)
	for i, w := range wordsCounts {
		if i >= 10 {
			break
		}
		result = append(result, w.word)
	}

	return result
}
