package main

import (
	"bufio"
	"github.com/lokicui/doc2vec-golang/doc2vec"
	"fmt"
	"log"
	"os"
	"strings"
    "strconv"
)

func main() {
	fname := os.Args[1]
	d2v := doc2vec.NewDoc2Vec(true, false, true, 5, 50, 1)
	err := d2v.LoadModel(fname)
	if err != nil {
		log.Fatal(err)
	}
	for true {
		reader := bufio.NewReader(os.Stdin)
        fmt.Println("please select operation type:\n\t0:word2words\n\t1:doc_likelihood\n\t2:leave one out key words\n\t3:sen2words\n\t4:sen2docs\n\t5:word2docs\n\t6:doc2docs\n\t7:doc2words\n\t")
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\n")
		switch text {
		case "0":
			fmt.Printf("Enter text:")
			text, _ = reader.ReadString('\n')
			d2v.Word2Words(strings.Trim(text, "\n"))
		case "1":
			fmt.Printf("Enter text:")
			text, _ = reader.ReadString('\n')
			likelihood := d2v.GetLikelihood4Doc(strings.Trim(text, "\n"))
			fmt.Printf("%v\t%v\n", text, likelihood)
		case "2":
			fmt.Printf("Enter text:")
			text, _ = reader.ReadString('\n')
			d2v.GetLeaveOneOutKwds(strings.Trim(text, "\n"), 50)
		case "3":
			fmt.Printf("Enter text:")
			text, _ = reader.ReadString('\n')
			d2v.Sen2Words(strings.Trim(text, "\n"), 50)
		case "4":
			fmt.Printf("Enter text:")
			text, _ = reader.ReadString('\n')
			d2v.Sen2Docs(strings.Trim(text, "\n"), 50)
		case "5":
			fmt.Printf("Enter text:")
			text, _ = reader.ReadString('\n')
			d2v.Word2Docs(strings.Trim(text, "\n"))
		case "6":
			fmt.Printf("Enter text:")
			text, _ = reader.ReadString('\n')
            idx, _ := strconv.Atoi(strings.Trim(text, "\n"))
			d2v.Doc2Docs(idx)
		case "7":
			fmt.Printf("Enter text:")
			text, _ = reader.ReadString('\n')
            idx, _ := strconv.Atoi(strings.Trim(text, "\n"))
			d2v.Doc2Words(idx)
		}
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
