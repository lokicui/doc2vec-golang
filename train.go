package main

import (
    "net/http"
    _ "net/http/pprof"
	"github.com/lokicui/doc2vec-golang/doc2vec"
	"log"
	"os"
)

func main() {
	fname := os.Args[1]

    // for pprof
    go func() {
        log.Println(http.ListenAndServe("localhost:16060", nil))
    }()

    log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
    //d2v := doc2vec.NewDoc2Vec(useCbow = true, useHS = false, useNEG = true, windowSize = 5, dim = 50, iters = 50)
    d2v := doc2vec.NewDoc2Vec(false, false, true, 5, 50, 50)
	d2v.Train(fname)
	err := d2v.SaveModel("2.model")
	if err != nil {
		log.Fatal(err)
	}
}
