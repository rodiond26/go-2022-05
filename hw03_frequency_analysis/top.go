package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(text string) []string {
	words := countWords(text)
	rating := rateByCount(words)

	return topNWords(rating, 10)
}

func countWords(text string) map[string]int {
	subStrs := strings.Fields(text)
	strMap := make(map[string]int)
	for _, str := range subStrs {
		_, ok := strMap[str]
		if ok {
			strMap[str]++
		} else {
			strMap[str] = 1
		}
	}
	return strMap
}

func rateByCount(words map[string]int) pairs {
	pairs := make(pairs, len(words))
	i := 0
	for k, v := range words {
		pairs[i] = pair{k, v}
		i++
	}

	sort.Sort(pairs)
	return pairs
}

func topNWords(p pairs, n int) []string {
	if len(p) == 0 {
		return []string{}
	}
	top := make([]string, 0)
	for i := 0; i < len(p) && i < n; i++ {
		top = append(top, p[i].word)
	}
	return top
}

type pair struct {
	word  string
	count int
}

type pairs []pair

func (p pairs) Len() int {
	return len(p)
}

func (p pairs) Less(i, j int) bool {
	if p[i].count == p[j].count {
		return p[i].word < p[j].word
	}
	return p[i].count > p[j].count
}

func (p pairs) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
