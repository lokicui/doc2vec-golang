package main

import (
	"flag"
	"github.com/lokicui/doc2vec-golang/doc2vec"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	corpus := flag.String("corpus", "", "training corpus file (required)")
	output := flag.String("output", "2.model", "output model file")
	useCbow := flag.Bool("cbow", false, "use CBOW model (default: Skip-Gram)")
	useHS := flag.Bool("hs", false, "use Hierarchical Softmax")
	useNEG := flag.Bool("neg", true, "use Negative Sampling")
	window := flag.Int("window", 5, "context window size")
	dim := flag.Int("dim", 50, "embedding dimension")
	iters := flag.Int("iters", 50, "training iterations")
	sweFile := flag.String("swe", "", "SWE synonym constraint file (optional)")
	sweCoeff := flag.Float64("swe-coeff", 0.1, "SWE interpolation coefficient")
	sweHinge := flag.Float64("swe-hinge", 0.0, "SWE hinge loss margin")
	sweDecay := flag.Float64("swe-decay", 0.0, "SWE weight decay coefficient")
	sweAddTime := flag.Float64("swe-addtime", 0.0, "SWE start time (training progress %)")
	flag.Parse()

	// backward compatibility: if no -corpus flag, use first positional arg
	fname := *corpus
	if fname == "" {
		if flag.NArg() > 0 {
			fname = flag.Arg(0)
		} else {
			log.Fatal("usage: ./train [-corpus] <corpus_file> [options]\n  use -h for help")
		}
	}

	go func() {
		log.Println(http.ListenAndServe("localhost:16060", nil))
	}()

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	d2v := doc2vec.NewDoc2Vec(*useCbow, *useHS, *useNEG, *window, *dim, *iters)

	if *sweFile != "" {
		log.Printf("Training with SWE constraints from: %s (coeff=%.2f, hinge=%.2f, decay=%.2f, addtime=%.1f%%)",
			*sweFile, *sweCoeff, *sweHinge, *sweDecay, *sweAddTime)
		d2v.TrainSWE(fname, *sweFile, *sweCoeff, *sweHinge, *sweDecay, *sweAddTime)
	} else {
		d2v.Train(fname)
	}

	err := d2v.SaveModel(*output)
	if err != nil {
		log.Fatal(err)
	}
}
