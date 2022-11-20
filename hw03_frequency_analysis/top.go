package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type wordCount struct {
	word  string
	count int
}

const (
	dash        = "-"
	punctuation = ".,!?:;"
)

func Top10(text string) []string {
	words := strings.Fields(text)

	wordsMap := make(map[string]int, len(words))
	for _, word := range words {
		word = strings.Trim(word, punctuation)
		if word == dash {
			continue
		}
		word = strings.ToLower(word)
		wordsMap[word]++
	}

	wordsCounts := make([]wordCount, 0, len(wordsMap))
	for word, count := range wordsMap {
		wordsCounts = append(wordsCounts, wordCount{word: word, count: count})
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
