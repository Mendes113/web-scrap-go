package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/gocolly/colly"
	"github.com/fogleman/gg"

	

	
	
)

func main() {
	url := "https://en.wikipedia.org/wiki/Go_(programming_language)"
	data := getHtml(url)
	topTenWords := getTopTenUsedWords(data)

	// Criação da nuvem de palavras para as dez palavras mais usadas
	graph(topTenWords)
	graphToPNG(topTenWords)

	
}

func getTopTenUsedWords(data []string) map[string]int {
	wordCounts := make(map[string]int)

	// Palavras a serem ignoradas
	ignoredWords := map[string]bool{
		"the": true,
		"and": true,
		"of":  true,
		"in":  true,
		"to":  true,
		"a":   true,
		"for": true,
		"as":  true,
		"by":  true,
		"on":  true,
		"or":  true,
		"an":  true,
		"with": true,
		"from": true,
		"that": true,
		"was":  true,
		"it":   true,
		"is":   true,
		"were": true,
		"are":  true,
		"at":   true,
		"be":   true,
		// Adicione mais palavras conforme necessário
	}

	for i := 1; i < len(data); i += 2 {
		toLower(data)
		words := strings.Fields(data[i])

		for _, word := range words {
			// Verifique se a palavra está na lista de palavras ignoradas
			if _, ok := ignoredWords[word]; !ok {
				wordCounts[word]++
			}
		}
	}

	// Obter as dez palavras mais usadas
	topWords := make(map[string]int)
	var wordCountPairs []wordCountPair

	for word, count := range wordCounts {
		wordCountPairs = append(wordCountPairs, wordCountPair{word, count})
	}

	sort.SliceStable(wordCountPairs, func(i, j int) bool {
		return wordCountPairs[i].count > wordCountPairs[j].count
	})

	for i := 0; i < 10 && i < len(wordCountPairs); i++ {
		topWords[wordCountPairs[i].word] = wordCountPairs[i].count
	}

	return topWords
}

func toLower(words []string) {
	for i := 0; i < len(words); i++ {
		words[i] = strings.ToLower(words[i])
	}
}

type wordCountPair struct {
	word  string
	count int
}


func graph(mostUsedWord map[string]int) {
	// Convertendo os dados para o formato aceito pela asciigraph
	var data []float64
	var labels []string

	// Extracting the maximum frequency for scaling
	maxFreq := 0
	for _, count := range mostUsedWord {
		if count > maxFreq {
			maxFreq = count
		}
	}

	for word, count := range mostUsedWord {
		// Scaling the count for the graph
		scaledCount := float64(count) / float64(maxFreq) * 100
		labels = append(labels, fmt.Sprintf("%s: %.2f", word, scaledCount))
		data = append(data, scaledCount)
	}

	// Plotting the graph with custom characters for the bars
	for i, value := range data {
		fmt.Printf("%-40s", labels[i])
		fmt.Println("")
		for j := 0; j < int(value); j += 5 {

			fmt.Print("█")
		}
		fmt.Println()
	}

	fmt.Println("Frequência das palavras")
}



func getHtml(urlString string) []string {
	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
	)
	var linkAndText []string

	c.OnHTML(".mw-parser-output", func(e *colly.HTMLElement){
		links := e.ChildAttr("a", "href")
		fmt.Print(links)
		text := e.ChildText("p")
		// fmt.Print(text)
		linkAndText = append(linkAndText, links, text)
	})

	if url, err := url.Parse(urlString); err != nil {
		panic(err)
	} else {
		c.Visit(url.String())
	}
	
	return linkAndText
}


func dataToCSV(data []string) {
	fName := "data.txt"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escreve cada linha do CSV com um link e texto
	for i := 0; i < len(data); i += 2 {
		record := []string{data[i], data[i+1]}
		err := writer.Write(record)
		if err != nil {
			log.Fatalf("Error writing to CSV: %s\n", err)
		}
	}
}


func graphToPNG(mostUsedWord map[string]int) {
	const width = 800
	const height = 400
	const barWidth = 20

	dc := gg.NewContext(width, height)

	var data []float64
	var labels []string

	// Extracting the maximum frequency for scaling
	maxFreq := 0
	for _, count := range mostUsedWord {
		if count > maxFreq {
			maxFreq = count
		}
	}

	for word, count := range mostUsedWord {
		// Scaling the count for the graph
		scaledCount := float64(count) / float64(maxFreq) * float64(height)
		labels = append(labels, fmt.Sprintf("%s: %.2f", word, scaledCount))
		data = append(data, scaledCount)
	}

	// Plotting the graph with custom characters for the bars
	for i, value := range data {
		x := float64(i) * barWidth
		dc.DrawRectangle(x, float64(height)-value, barWidth, value)
		dc.SetRGB(255, 0, 0)
		dc.Fill()
		dc.DrawStringAnchored(labels[i], x, float64(height)-5, 0.5, 0.5)
	}

	// Save to PNG
	if err := dc.SavePNG("graph.png"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Graph saved as graph.png")
}