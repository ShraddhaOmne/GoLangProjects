package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	filename := "Papaya.txt"

	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		fmt.Println("File does not exist:", filename)
		return
	}

	if info.Size() == 0 {
		fmt.Println("File is empty:", filename)
		return
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	text := strings.ReplaceAll(string(content), "\n", " ")
	words := strings.Split(text, " ")

	wordCount := make(map[string]int)
	for _, word := range words {
		if word != "" {
			wordCount[word]++
		}
	}
	fmt.Println("Total unique words:", len(wordCount))

	maxCount := 0
	var maxWords []string
	for word, count := range wordCount {
		if count > maxCount {
			maxCount = count
			maxWords = []string{word}
		} else if count == maxCount {
			maxWords = append(maxWords, word)
		}
	}

	// minCount := -1
	// var minWords []string
	// for word, count := range wordCount {
	// 	if minCount == -1 || count < minCount {
	// 		minCount = count
	// 		minWords = []string{word}
	// 	} else if count == minCount {
	// 		minWords = append(minWords, word)
	// 	}
	// }
	fmt.Println("Max occurrence count:", maxCount)

	fmt.Println("Word(s) with max occurrence:", maxWords)

	// fmt.Println("Min occurrence count:", minCount)
	// fmt.Println("Word(s) with min occurrence:", minWords)
}
