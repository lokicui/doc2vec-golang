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
    Word2Words(word string)
    Word2Docs(word string)
    Sen2Words(content string, iters int)
    Sen2Docs(content string, iters int)
    Doc2Docs(docidx int)
    Doc2Words(docidx int)
    GetLikelihood4Doc(context string) (likelihood float64)
    GetLeaveOneOutKwds(content string, iters int)
}

type TDoc2VecImpl struct {
	Trainfile    string
	Dim          int
    UseCbow      bool   //true:Continuous Bag-of-Word Model false:skip-gram
	WindowSize   int    //cbow model的窗口大小
	UseHS        bool
	UseNEG       bool   //UseHS / UseNEG两种求解优化算法必须选一个 也可以两种算法都选 详见google word2vec源代码
    Negative     int    //负采样词的个数
	StartAlpha   float64
	Iters        int
	TrainedWords int
	Corpus       corpus.ICorpus
	NN           neuralnet.INeuralNet
}
