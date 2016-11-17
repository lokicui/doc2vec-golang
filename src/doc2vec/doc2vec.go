package doc2vec

import (
	"bytes"
	"common"
	"corpus"
	"encoding/binary"
	"fmt"
	"github.com/tinylib/msgp/msgp"
	"log"
	"math"
	"neuralnet"
	"os"
	"sort"
	"sync"
)

const (
	MAX_EXP        float64 = 6
	EXP_TABLE_SIZE int     = 1000
    NEG_SAMPLING_TABLE_SIZE int = 1e8
)

var gExpTable [EXP_TABLE_SIZE]float64
var gNegSamplingTable   [NEG_SAMPLING_TABLE_SIZE]int32
var gNextRandom uint64 = 1

func init() {
	for i := 0; i < EXP_TABLE_SIZE; i++ {
		gExpTable[i] = math.Exp((float64(i)/float64(EXP_TABLE_SIZE)*2.0 - 1.0) * MAX_EXP)
		gExpTable[i] = gExpTable[i] / (gExpTable[i] + 1.0)
	}
}

func (p * TDoc2VecImpl) InitUnigramTable() {
    train_words_power := 0.0
    power := 0.75
    words := p.Corpus.GetAllWords()
    if NEG_SAMPLING_TABLE_SIZE <= len(words) {
        log.Fatal("NEG_SAMPLING_TABLE_SIZE < len(words)")
    }
    for _, worditem := range words {
        train_words_power += math.Pow(float64(worditem.Cnt), power)
    }
    var i int32 = 0
    d1 := math.Pow(float64(words[i].Cnt), power) / train_words_power
    vocabsize := int32(p.Corpus.GetVocabCnt())
    for a := 0; a < NEG_SAMPLING_TABLE_SIZE; a ++ {
        gNegSamplingTable[a] = i
        if float64(a) / float64(NEG_SAMPLING_TABLE_SIZE) > d1 {
            i ++
            d1 += math.Pow(float64(words[i].Cnt), power) / train_words_power
        }
        if i >= vocabsize {
            i = vocabsize - 1
        }
    }
}

func GetSigmoidValue(f float64) float64 {
	idx := int((f + MAX_EXP) * (float64(EXP_TABLE_SIZE) / MAX_EXP / 2.0))
	return gExpTable[idx]
}
func GetNegativeSamplingWordIdx() int32 {
	gNextRandom = gNextRandom*25214903917 + 11
    idx := int(int(gNextRandom >> 16) % NEG_SAMPLING_TABLE_SIZE)
    target := gNegSamplingTable[idx]
    return int32(target)
}

func (p TSortItemSlice) Len() int {
	return len(p)
}
func (p TSortItemSlice) Less(i, j int) bool {
	return p[i].Dis < p[j].Dis
}
func (p TSortItemSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func NewDoc2Vec(useCbow, useHS, useNEG bool, windowSize, dim, iters int) IDoc2Vec {
	self := &TDoc2VecImpl{
        UseCbow:    useCbow,
		UseHS:      useHS,
		UseNEG:     useNEG,
		WindowSize: windowSize,
		Dim:        dim,
        Negative:   5,
		Corpus:     corpus.NewCorpus(),
		NN:         neuralnet.NewNN(0, 0, 0, false, false),
		StartAlpha: 0.025,
		Iters:      iters,
	}
	return IDoc2Vec(self)
}

func (p *TDoc2VecImpl) GetCorpus() corpus.ICorpus {
	return p.Corpus
}

func (p *TDoc2VecImpl) GetNeuralNet() neuralnet.INeuralNet {
	return p.NN
}

func (p *TDoc2VecImpl) getRandomWindowSize() int {
	gNextRandom = gNextRandom*25214903917 + 11
	return int(gNextRandom % uint64(p.WindowSize))
}

func (p *TDoc2VecImpl) Train(fname string) {
	p.Trainfile = fname
	p.Corpus.Build(fname)
    if p.UseNEG {
        p.InitUnigramTable()
    }
	p.NN = neuralnet.NewNN(p.Corpus.GetDocCnt(), p.Corpus.GetVocabCnt(), p.Dim, p.UseHS, p.UseNEG)
	for i := 0; i < p.Iters; i++ {
		p.TrainCbow()
	}
}

func (p *TDoc2VecImpl) SaveModel(fname string) (err error) {
	fd, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()
	return p.EncodeMsg(msgp.NewWriter(fd))
}

func (p *TDoc2VecImpl) SaveModel_byself(fname string) (err error) {
	fd, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, int32(p.Corpus.GetVocabCnt()))
	if err != nil {
		log.Fatal(err)
	}
	err = binary.Write(buf, binary.BigEndian, int32(p.Corpus.GetDocCnt()))
	if err != nil {
		log.Fatal(err)
	}
	err = binary.Write(buf, binary.BigEndian, int32(p.Dim))
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < p.Corpus.GetVocabCnt(); i++ {
		item := p.Corpus.GetWordItemByIdx(i)
		err = binary.Write(buf, binary.BigEndian, int32(len(item.Word)))
		if err != nil {
			log.Fatal(err)
		}
		buf.WriteString(item.Word)
	}
	buf.WriteTo(fd)
	buf.Truncate(0)
	for i := 0; i < p.Corpus.GetVocabCnt(); i++ {
		vector := p.NN.GetSyn0(int32(i))
		err = binary.Write(buf, binary.BigEndian, *vector)
		if err != nil {
			log.Fatal(err)
		}
		if i%1000 == 0 {
			fmt.Printf("save %v words vector\n", i)
		}
	}
	buf.WriteTo(fd)
	buf.Truncate(0)
	for i := 0; i < p.Corpus.GetDocCnt(); i++ {
		vector := p.NN.GetDSyn0(int32(i))
		err = binary.Write(buf, binary.BigEndian, *vector)
		if err != nil {
			log.Fatal(err)
		}
		if i%100 == 0 {
			fmt.Printf("save %v words doc\n", i)
		}
	}
	buf.WriteTo(fd)
	buf.Truncate(0)
	return err
}

func (p *TDoc2VecImpl) LoadModel(fname string) (err error) {
	fd, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()
	return p.DecodeMsg(msgp.NewReader(fd))
}

/*
func (p *TDoc2VecImpl) LoadModel_byself(fname string) (err error) {
	fd, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()
	var vocabsize, docSize, dim int32
	err = binary.Read(fd, binary.BigEndian, &vocabsize)
	if err != nil {
		log.Fatal(err)
	}
	binary.Read(fd, binary.BigEndian, &docSize)
	if err != nil {
		log.Fatal(err)
	}
	binary.Read(fd, binary.BigEndian, &dim)
	if err != nil {
		log.Fatal(err)
	}
	p.Dim = int(dim)
	p.VocabSize = int(vocabsize)
	p.Words = make([]string, vocabsize, vocabsize)
	for i := 0; i < int(vocabsize); i++ {
		var size int32
		binary.Read(fd, binary.BigEndian, &size)
		bytes := make([]byte, size, size)
		fd.Read(bytes)
		p.Words[i] = string(bytes)
	}
	p.NN = neuralnet.NewNN(int(docSize), int(vocabsize), p.Dim, false, false)
	for i := 0; i < int(vocabsize); i++ {
		vector := p.NN.GetSyn0(int32(i))
		for j := 0; j < p.Dim; j++ {
			binary.Read(fd, binary.BigEndian, &(*vector)[j])
		}
	}
	for i := 0; i < int(docSize); i++ {
		vector := p.NN.GetDSyn0(int32(i))
		for j := 0; j < p.Dim; j++ {
			binary.Read(fd, binary.BigEndian, &(*vector)[j])
		}
	}
	return err
}
*/

func (p *TDoc2VecImpl) GetAlpha() float64 {
	alpha := p.StartAlpha * (1.0 - float64(p.TrainedWords)/float64(p.Iters*p.Corpus.GetWordsCnt()+1))
	if alpha < p.StartAlpha*0.0001 {
		alpha = p.StartAlpha * 0.0001
	}
	return alpha
}

func (p *TDoc2VecImpl) TrainCbow() {
	//Continuous Bag-of-Word Model
	tokens := make(chan struct{}, 32)
	wg := new(sync.WaitGroup)
	last_trained_words := 0
	alpha := p.GetAlpha()
	for docidx_, wordsidx_ := range p.Corpus.GetAllDocWordsIdx() {
		docidx, wordsidx := docidx_, wordsidx_
		wg.Add(1)
		go func() {
			defer func() { <-tokens }()
			defer wg.Done()
			tokens <- struct{}{}
			//train one document
			last_trained_words += len(wordsidx)
			p.TrainedWords += len(wordsidx)
			if last_trained_words > 10000 {
				last_trained_words = 0
				alpha = p.GetAlpha()
			}
			for spos, widx := range wordsidx {
				neu1 := make(neuralnet.TVector, p.Dim, p.Dim)     //X(w)
				neu1e := make(neuralnet.TVector, p.Dim, p.Dim)    //e
				syn1copy := make(neuralnet.TVector, p.Dim, p.Dim) //为了计算 g*Theta
				neu1copy := make(neuralnet.TVector, p.Dim, p.Dim) //为了计算 g*X(w)
				b := p.getRandomWindowSize()
				start := common.Max(0, spos-p.WindowSize+b)
				end := common.Min(len(wordsidx), spos+p.WindowSize-b+1)
				//in -> hidden      X(widx) = E[V(a)]
				cw := 0
				for a := start; a < end; a++ {
					if a == spos {
						continue
					}
					idx := wordsidx[a]
					neu1.Add(*p.NN.GetSyn0(idx))
					cw++
				}
				//X(widx) += Document Vector
				dsyn0 := p.NN.GetDSyn0(int32(docidx))
				neu1.Add(*dsyn0)
				cw++
				neu1.Divide(float32(cw))

                //Hierarchical Softmax
                if p.UseHS {
                    //foreach inner node of words[widx]
                    worditem := p.Corpus.GetWordItemByIdx(int(widx))
                    for i, point := range worditem.Point {
                        syn1 := p.NN.GetSyn1(point) //Theta
                        // Propagate hidden -> output
                        f := neu1.Dot(*syn1) // f = Sigmoid[X(w) dot Theta]
                        g := 0.0 //g = alpha * (1 - Dj(w) - f)
                        if f >= MAX_EXP {
                            f = 1.0
                        } else if f <= -MAX_EXP {
                            f = 0.0
                        } else {
                            f = GetSigmoidValue(f)
                        }
                        // 'g' is the gradient multiplied by the learning rate
                        label := 0.0
                        if worditem.Code[i] {
                            label = 1.0
                        }
                        g = (1.0 - label - f) * f * alpha
                        // Propagate errors output -> hidden  e := e + g*Theta
                        copy(syn1copy, *syn1)
                        syn1copy.Multiply(g)
                        neu1e.Add(syn1copy)
                        // Learn weights hidden -> output   Theta := Theta + gX(w)
                        copy(neu1copy, neu1)
                        neu1copy.Multiply(g)
                        syn1.Add(neu1copy)
                    }
                }

                if p.UseNEG {
                    label := 0.0
                    var point int32 = 0
                    for i := 0; i < p.Negative + 1; i ++ {
                        if i == 0 {
                            label = 1.0
                            point = widx
                        } else {
                            label = 0.0
                            point = GetNegativeSamplingWordIdx()
                            if point == widx {
                                continue
                            }
                        }
                        syn1neg := p.NN.GetSyn1Neg(point)
                        f := neu1.Dot(*syn1neg)
                        g := 0.0
                        if f >= MAX_EXP {
                            f = 1.0
                        } else if f <= -MAX_EXP {
                            f = 0.0
                        } else {
                            f = GetSigmoidValue(f)
                        }
                        g = (label - f) * alpha
                        // Propagate errors output -> hidden  e := e + g*Theta
                        copy(syn1copy, *syn1neg)
                        syn1copy.Multiply(g)
                        neu1e.Add(syn1copy)
                        // Learn weights hidden -> output   Theta := Theta + gX(w)
                        copy(neu1copy, neu1)
                        neu1copy.Multiply(g)
                        syn1neg.Add(neu1copy)
                    }
                }

				// hidden -> in                         v(u) := v(u) + e
				for a := start; a < end; a++ {
					if a == spos {
						continue
					}
					idx := wordsidx[a]
					syn0 := p.NN.GetSyn0(idx)
					syn0.Add(neu1e)
				}
				// hidden -> in                         D(u) := D(u) + e
				dsyn0.Add(neu1e)
			}
		}()
	}
	wg.Wait()
}

func (p *TDoc2VecImpl) FindKNN(word string) {
	idx, ok := p.Corpus.GetWordIdx(word)
	if !ok {
		return
	}
	needle_worditem := p.Corpus.GetWordItemByIdx(int(idx))
	vector := p.NN.GetSyn0(idx)
	vocabsize := p.Corpus.GetVocabCnt()
	dis_vector := make(TSortItemSlice, vocabsize, vocabsize)
	for i := 0; i < vocabsize; i++ {
		dis := ConsineDistance(*vector, *p.NN.GetSyn0(int32(i)))
		dis_vector[i] = &SortItem{Idx: int32(i), Dis: dis}
	}
	//大爷的 go的排序太麻烦,还不如自己写个快排
    //QuickSort(0, len(dis_vector) - 1, dis_vector)
	sort.Sort(sort.Reverse(dis_vector))
	fmt.Printf("word:%v\n", needle_worditem.Word)
	for i := 0; i < len(dis_vector) && i < 10; i++ {
		item := dis_vector[i]
		dis := item.Dis
		idx := int(item.Idx)
		worditem := p.Corpus.GetWordItemByIdx(idx)
		fmt.Printf("\t%v\t%v\t%v\n", dis, worditem.Word, *p.NN.GetSyn0(int32(item.Idx)))
	}
	fmt.Println()
}

func QuickSort(i, j int, vec []*SortItem) {
	ii, jj := i, j
	if i+1 >= j {
		return
	} else if i+2 == j {
		if vec[i].Dis > vec[j-1].Dis {
			vec[i], vec[j-1] = vec[j-1], vec[i]
		}
	}
	M := i
	stub := vec[M]
	for i < j {
		for ; j > i; j-- {
			if vec[j].Dis < stub.Dis {
				vec[M] = vec[j]
				M = j
				break
			}
		}
		for ; i < j; i++ {
			if vec[i].Dis > stub.Dis {
				vec[M] = vec[i]
				M = i
				break
			}
		}
	}
	vec[M] = stub
	QuickSort(ii, M, vec)
	QuickSort(M+1, jj, vec)
}

func ConsineDistance(a neuralnet.TVector, b neuralnet.TVector) (dis float64) {
	var sum, sum_a, sum_b float64
	for i := 0; i < len(a); i++ {
		sum += float64(a[i] * b[i])
		sum_a += float64(a[i] * a[i])
		sum_b += float64(b[i] * b[i])
	}
	dis = sum / math.Sqrt(sum_a) / math.Sqrt(sum_b)
	return dis
}

func (p *TDoc2VecImpl) TrainHSSkipGram() {
	//Skip-Gram Model + Hierarchical Softmax
}
