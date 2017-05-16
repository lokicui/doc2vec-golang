package main

import (
	"github.com/lokicui/doc2vec-golang/doc2vec"
	"log"
	"os"
)

func main() {
	fname := os.Args[1]
    log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
    //d2v := doc2vec.NewDoc2Vec(useCbow = true, useHS = false, useNEG = true, windowSize = 5, dim = 50, iters = 50)
    d2v := doc2vec.NewDoc2Vec(true, false, true, 5, 50, 50)
	d2v.Train(fname)
	err := d2v.SaveModel("3.model")
	if err != nil {
		log.Fatal(err)
	}
}
