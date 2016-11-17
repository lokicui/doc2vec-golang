package main

import (
	"doc2vec"
	"os"
)

func main() {
	fname := os.Args[1]
	d2v := doc2vec.NewDoc2Vec(true, false, true, 5, 100, 100)
	d2v.Train(fname)
	d2v.SaveModel("3.model")
}
