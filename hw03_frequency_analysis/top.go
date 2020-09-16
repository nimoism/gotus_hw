package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"math"
	"regexp"
	"sort"
	"strings"
)

const maxTopCount = 10

type wordInfo struct {
	word  string
	count int
}

var skipWords = map[string]struct{}{
	"":  {}, // skip empty words in case of text is space bordered
	"-": {},
}

var sep = regexp.MustCompile(`[\s\t\n'"!]+`)

func Top10(input string) []string {
	words := sep.Split(input, -1)

	infos := wordsInfos(words)

	sort.Slice(infos, func(i, j int) bool {
		return infos[i].count > infos[j].count
	})

	topCount := int(math.Min(float64(len(infos)), maxTopCount))
	tops := make([]string, 0, topCount)
	for _, info := range infos[:topCount] {
		tops = append(tops, info.word)
	}
	return tops
}

func wordsInfos(words []string) []wordInfo {
	countsMap := make(map[string]int, len(words))
	for _, word := range words {
		if _, ok := skipWords[word]; !ok {
			word = strings.ToLower(word)
			countsMap[word]++
		}
	}

	infos := make([]wordInfo, 0, len(countsMap))
	for word, count := range countsMap {
		infos = append(infos, wordInfo{word: word, count: count})
	}
	return infos
}
