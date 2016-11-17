package main

import (
	"bufio"
	"doc2vec"
	"fmt"
	"os"
	"strings"
)

func main() {
	fname := os.Args[1]
	d2v := doc2vec.NewDoc2Vec(true, false, true, 5, 50, 1)
	d2v.LoadModel(fname)
	for true {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter text:")
		text, _ := reader.ReadString('\n')
		d2v.FindKNN(strings.Trim(text, "\n"))
	}
	td := d2v.GetCorpus()
	for _, worditem := range td.GetAllWords() {
		fmt.Printf("%+v\n", worditem)
	}
	for _, words := range td.GetAllDocWords() {
		sen := []string{}
		for _, word := range words {
			sen = append(sen, word.Word)
		}
		ss := strings.Join(sen, " ")
		fmt.Println(ss)
	}
}
