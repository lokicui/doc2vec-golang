package doc2vec

import (
	"corpus"
	"neuralnet"
)

//go:generate msgp

type SortItem struct {
	Idx int32
	Dis float64
}

type TSortItemSlice []*SortItem


type IDoc2Vec interface {
	Train(fname string)
	GetCorpus() corpus.ICorpus
	GetNeuralNet() neuralnet.INeuralNet
	SaveModel(fname string) (err error)
	LoadModel(fname string) (err error)
	FindKNN(word string)
}

type TDoc2VecImpl struct {
	Trainfile    string
	UseHS        bool
	UseNEG       bool
	Dim          int
	WindowSize   int
	StartAlpha   float64
	Iters        int
	TrainedWords int
	Corpus       corpus.ICorpus
	NN           neuralnet.INeuralNet
}
