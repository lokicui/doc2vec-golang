package main

import (
	"doc2vec"
	"log"
	"os"
)

func main() {
	fname := os.Args[1]
	d2v := doc2vec.NewDoc2Vec(true, false, true, 5, 50, 50)
	d2v.Train(fname)
	err := d2v.SaveModel("3.model")
	if err != nil {
		log.Fatal(err)
	}
}
