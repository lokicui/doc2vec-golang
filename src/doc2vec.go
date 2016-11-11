package main

import (
    "fmt"
    "os"
    "math"
    "sync"
    "strings"
    "corpus"
    "neuralnet"
    "common"
)

const (
    MAX_EXP float64 = 6
    EXP_TABLE_SIZE  int = 1000
)

type Doc2Vec interface {
    Train(fname string)
    GetCorpus() corpus.Corpus
    GetNN()     neuralnet.NN
}

type Doc2VecImpl struct {
    trainfile   string
    useHS       bool
    useNEG      bool
    dim         int
    windowSize  int
    corpus      corpus.Corpus
    nn          neuralnet.NN
    nextRandom  int64
    alpha       float64
    expTable    [EXP_TABLE_SIZE]float64
}

func NewDoc2Vec(useHS, useNEG bool, windowSize, dim int) Doc2Vec {
    self := &Doc2VecImpl{
        useHS: useHS,
        useNEG: useNEG,
        windowSize: windowSize,
        dim: dim,
        corpus: nil,
        nn: nil,
    }
    for i := 0; i < EXP_TABLE_SIZE; i ++ {
        self.expTable[i] = math.Exp((float64(i) / float64(EXP_TABLE_SIZE) * 2.0 - 1.0) * MAX_EXP)
        self.expTable[i] = self.expTable[i] / (self.expTable[i] + 1.0)
    }
    return Doc2Vec(self)
}

func (p *Doc2VecImpl) GetCorpus() corpus.Corpus {
    return p.corpus
}

func (p *Doc2VecImpl) GetNN() neuralnet.NN {
    return p.nn
}


func (p *Doc2VecImpl) GetSigmoidValue(f float64) float64 {
    idx := int((f + MAX_EXP) * (float64(EXP_TABLE_SIZE) / MAX_EXP / 2.0))
    return p.expTable[idx]
}

func (p * Doc2VecImpl) getRandomWindowSize() int {
    p.nextRandom = p.nextRandom * 25214903917 + 11
    return int(p.nextRandom % int64(p.windowSize))
}

func (p * Doc2VecImpl) Train(fname string) {
    p.trainfile = fname
    p.corpus = corpus.NewCorpus()
    p.corpus.Build(fname)
    p.nn = neuralnet.NewNN(p.corpus.GetDocSize(), p.corpus.GetVocabSize(), p.dim, p.useHS, p.useNEG)
    p.TrainHSCbow()
}

func (p * Doc2VecImpl) TrainHSCbow() {
    //Continuous Bag-of-Word Model + Hierarchical Softmax
    tokens := make(chan struct {}, 32)
    wg := new(sync.WaitGroup)
    for docidx_, wordsidx_ := range p.corpus.GetAllDocWordsIdx() {
        docidx, wordsidx := docidx_, wordsidx_
        wg.Add(1)
        go func() {
            defer func() {<-tokens}()
            defer wg.Done()
            tokens <- struct {}{}
            //train one document
            for spos, widx := range wordsidx {
                neu1 := make(neuralnet.TVector, p.dim, p.dim)       //X(w)
                neu1e := make(neuralnet.TVector, p.dim, p.dim)      //e
                syn1copy := make(neuralnet.TVector, p.dim, p.dim)   //为了计算 g*Theta
                neu1copy := make(neuralnet.TVector, p.dim, p.dim)   //为了计算 g*X(w)
                b := p.getRandomWindowSize()
                start := common.Max(0, spos - p.windowSize + b)
                end := common.Min(len(wordsidx), spos + p.windowSize - b + 1)
                //in -> hidden      X(widx) = E[V(a)]
                cw := 0
                for a := start; a < end; a ++ {
                    if a == spos {
                        continue
                    }
                    idx := wordsidx[a]
                    neu1.Add(*p.nn.GetSyn0(idx))
                    cw ++
                }
                //X(widx) += Document Vector
                dsyn0 := p.nn.GetDSyn0(int32(docidx))
                neu1.Add(*dsyn0)
                cw ++
                neu1.Divide(float32(cw))

                //foreach inner node of words[widx]
                worditem := p.corpus.GetWordItemByIdx(int(widx))
                for i, point := range worditem.Point {
                    syn1 := p.nn.GetSyn1(point) //Theta
                    // Propagate hidden -> output
                    f := neu1.Dot(*syn1)     // f = Sigmoid[X(w) dot Theta]
                    if f >= MAX_EXP || f <= -MAX_EXP {
                        continue
                    }
                    f = p.GetSigmoidValue(f)
                    g := 0.0    //g = alpha * (1 - Dj(w) - f)
                    // 'g' is the gradient multiplied by the learning rate
                    if worditem.Code[i] {
                        g = -1.0 * f * p.alpha
                    } else {
                        g = (1.0 - f) * p.alpha
                    }
                    // Propagate errors output -> hidden  e := e + g*Theta
                    copy(syn1copy, *syn1)
                    syn1copy.Multiply(g)
                    neu1e.Add(syn1copy)
                    // Learn weights hidden -> output   Theta := Theta + gX(w)
                    copy(neu1copy, neu1)
                    neu1copy.Multiply(g)
                    syn1.Add(neu1copy)
                }

                // hidden -> in                         v(u) := v(u) + e
                for a := start; a < end; a ++ {
                    if a == spos {
                        continue
                    }
                    idx := wordsidx[a]
                    syn0 := p.nn.GetSyn0(idx)
                    syn0.Add(neu1e)
                }
                // hidden -> in                         D(u) := D(u) + e
                dsyn0.Add(neu1e)
            }
        }()
    }
    wg.Wait()
}

func (p* Doc2VecImpl) TrainHSSkipGram() {
    //Skip-Gram Model + Hierarchical Softmax
}

func main() {
    fname := os.Args[1]
    //td := corpus.NewCorpus()
    ////fmt.Printf("%#v, %s\n", vocabulary, fname)
    //td.Build(fname)
    //for _, worditem := range td.GetAllWords() {
    //    fmt.Printf("%+v\n", worditem)
    //}
    //for _, words := range td.GetAllDocWords() {
    //    sen := []string {}
    //    for _, word := range words {
    //        sen = append(sen, word.Word)
    //    }
    //    ss := strings.Join(sen, " ")
    //    fmt.Println(ss)
    //}
    doc2vec := NewDoc2Vec(true, true, 5, 50)
    doc2vec.Train(fname)
    td := doc2vec.GetCorpus()
    for _, worditem := range td.GetAllWords() {
        fmt.Printf("%+v\n", worditem)
    }
    for _, words := range td.GetAllDocWords() {
        sen := []string {}
        for _, word := range words {
            sen = append(sen, word.Word)
        }
        ss := strings.Join(sen, " ")
        fmt.Println(ss)
    }
}
