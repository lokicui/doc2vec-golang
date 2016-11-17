package main

import (
	"doc2vec"
	"os"
)

func main() {
	fname := os.Args[1]
	d2v := doc2vec.NewDoc2Vec(true, true, 5, 50, 50)
	d2v.Train(fname)
	d2v.SaveModel("3.model")
}
