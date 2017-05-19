package doc2vec

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/lokicui/doc2vec-golang/common"
	"github.com/lokicui/doc2vec-golang/corpus"
	"github.com/lokicui/doc2vec-golang/neuralnet"
	"github.com/tinylib/msgp/msgp"
	"log"
	"math"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

var _ = sort.Sort
var _ = bytes.NewBuffer
var _ = binary.Read

const (
	MAX_EXP                 float64 = 6.0
	EXP_TABLE_SIZE          int     = 1000
	NEG_SAMPLING_TABLE_SIZE int     = 1e8
	PROGRESS_BAR_THRESHOLD  int     = 100000
	THREAD_NUM              int     = 32
)

var gExpTable [EXP_TABLE_SIZE]float64
var gNegSamplingTable [NEG_SAMPLING_TABLE_SIZE]int32
var gNextRandom uint64 = 1

func init() {
	for i := 0; i < EXP_TABLE_SIZE; i++ {
		gExpTable[i] = math.Exp((float64(i)/float64(EXP_TABLE_SIZE)*2.0 - 1.0) * MAX_EXP)
		gExpTable[i] = gExpTable[i] / (gExpTable[i] + 1.0)
	}
}

func (p *TDoc2VecImpl) initUnigramTable() {
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
	for a := 0; a < NEG_SAMPLING_TABLE_SIZE; a++ {
		gNegSamplingTable[a] = i
		if float64(a)/float64(NEG_SAMPLING_TABLE_SIZE) > d1 {
			i++
			d1 += math.Pow(float64(words[i].Cnt), power) / train_words_power
		}
		if i >= vocabsize {
			i = vocabsize - 1
		}
	}
}

func GetSigmoidValue(f float64) float64 {
	idx := int((f + MAX_EXP) * (float64(EXP_TABLE_SIZE) / MAX_EXP / 2.0))
	if idx >= EXP_TABLE_SIZE || idx < 0 {
		log.Fatal("GetSigmoidValue with", f, "idx=", idx)
	}
	return gExpTable[idx]
}
func GetNegativeSamplingWordIdx() int32 {
	gNextRandom = gNextRandom*25214903917 + 11
	idx := int(int(gNextRandom>>16) % NEG_SAMPLING_TABLE_SIZE)
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
		StartAlpha: common.If(useCbow, 0.05, 0.025).(float64),
		Iters:      iters,
        Pool:       nil,
	}
    self.Pool = &sync.Pool{
        New: func() interface{} {
            vector := make(neuralnet.TVector, self.Dim, self.Dim)
            return &vector
        },
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

func (p *TDoc2VecImpl) getLikelihood4Pair(widx int32, rangevec *neuralnet.TVector) (likelihood float64) {
	//Hierarchical Softmax
	if p.UseHS {
		//foreach inner node of words[widx]
		worditem := p.Corpus.GetWordItemByIdx(int(widx))
		for i, point := range worditem.Point {
			syn1 := p.NN.GetSyn1(point) //Theta
			f := rangevec.Dot(*syn1)    // f = Sigmoid[X(w) dot Theta]
			label := -1.0
			if worditem.Code[i] {
				label = 1.0
			}
			likelihood += -1.0 * math.Log(1.0+math.Exp(label*f)) // II[1/(1+e^-x)]  连乘取log
		}
	} else if p.UseNEG {
		//@todo
		likelihood = 0.0
	}
	return likelihood
}

func (p *TDoc2VecImpl) GetLikelihood4Doc(context string) (likelihood float64) {
	wordsidx := p.Corpus.Transform(context)
	for spos, widx := range wordsidx {
		//针对计算doc中每个词的likelihood
		start := common.Max(0, spos-p.WindowSize)
		end := common.Min(len(wordsidx), spos+p.WindowSize+1)
		//in -> hidden      X(widx) = E[V(a)]

		if p.UseCbow {
			neu1 := make(neuralnet.TVector, p.Dim, p.Dim) //X(w)
			cw := 0
			for a := start; a < end; a++ {
				if a == spos {
					continue
				}
				idx := wordsidx[a]
				neu1.Add(*p.NN.GetSyn0(idx))
				cw++
			}
			//##################################################################
			//X(widx) += Document Vector
			//dsyn0 := p.NN.GetDSyn0(int32(docidx))
			//neu1.Add(*dsyn0)
			//cw++
			//Note: @todo
			//这里可以考虑是否先生成doc的向量,然后讲doc的向量也加到neu1里面
			//##################################################################

			neu1.Divide(float32(cw))
			likelihood += p.getLikelihood4Pair(widx, &neu1)
		} else {
			for a := start; a < end; a++ {
				if a == spos {
					continue
				}
				idx := wordsidx[a]
				rangevec := p.NN.GetSyn0(idx)
				likelihood += p.getLikelihood4Pair(widx, rangevec)
			}
		}
	}
	return likelihood
}

//online fit doc vector
func (p *TDoc2VecImpl) fitDoc(wordsidx []int32, iters int) (dsyn0 *neuralnet.TVector) {
	dsyn0 = p.NN.NewDSyn0()
	trainedwords := 0
	totalwords := iters*len(wordsidx) + 1
	for i := 0; i < iters; i++ {
		alpha := p.getAlpha(trainedwords, totalwords)
		trainedwords += len(wordsidx)
		if p.UseCbow {
			p.trainCbow4Document(wordsidx, dsyn0, alpha, true)
		} else {
			p.trainSkipGram4Document(wordsidx, dsyn0, alpha, true)
		}
	}
	return dsyn0
}

func (p *TDoc2VecImpl) FitDoc(context string, iters int) (dsyn0 *neuralnet.TVector) {
	wordsidx := p.Corpus.Transform(context)
	return p.fitDoc(wordsidx, iters)
}

func (p *TDoc2VecImpl) Train(fname string) {
	p.Trainfile = fname
	p.Corpus.Build(fname)
	if p.UseNEG {
		p.initUnigramTable()
	}
	p.NN = neuralnet.NewNN(p.Corpus.GetDocCnt(), p.Corpus.GetVocabCnt(), p.Dim, p.UseHS, p.UseNEG)
	if p.UseCbow {
		p.trainCbow()
	} else {
		p.trainSkipGram()
	}
}

func (p *TDoc2VecImpl) SaveModel(fname string) (err error) {
	fd, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()
	writer := msgp.NewWriter(fd)
	err = p.EncodeMsg(writer)
	if err == nil {
		//你大爷 必须Flush, 不然即便你调用了Close()也无济于事
		writer.Flush()
	}
	return
}

func (p *TDoc2VecImpl) LoadModel(fname string) (err error) {
	fd, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()
	err = p.DecodeMsg(msgp.NewReader(fd))
	if err == nil && p.UseNEG {
		p.initUnigramTable()
	}
	return err
}

/*
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

func (p *TDoc2VecImpl) getAlpha(trained, total int) float64 {
	alpha := p.StartAlpha * (1.0 - float64(trained)/float64(total+1))
	if alpha < p.StartAlpha*0.0001 {
		alpha = p.StartAlpha * 0.0001
	}
	return alpha
}
func (p *TDoc2VecImpl) getTrainAlpha() float64 {
	return p.getAlpha(p.TrainedWords, p.Iters*p.Corpus.GetWordsCnt()+1)
}

// P(V(centralwidx) | rangevec) foreach rangevec in Context(centralwidx)
func (p *TDoc2VecImpl) trainSkipGram4Pair(centralwidx int32, rangevec *neuralnet.TVector, alpha float64, infer bool) {
	neu1e := make(neuralnet.TVector, p.Dim, p.Dim)    //e
	syn1copy := make(neuralnet.TVector, p.Dim, p.Dim) //为了计算 g*Theta
	neu1copy := make(neuralnet.TVector, p.Dim, p.Dim) //为了计算 g*V(w), V(w) = rangevec
	//in -> hidden     = V(a)
	_ = *rangevec

	//Hierarchical Softmax
	if p.UseHS {
		//foreach inner node of words[centralwidx]
		worditem := p.Corpus.GetWordItemByIdx(int(centralwidx))
		for i, point := range worditem.Point {
			syn1 := p.NN.GetSyn1(point) //Theta
			// Propagate hidden -> output
			f := rangevec.Dot(*syn1) // f = Sigmoid[X(w) dot Theta]
			if math.IsNaN(f) {
				log.Fatal("f is NaN, try to reduce StartAlpha!")
			}
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
			g = (1.0 - label - f) * alpha
			// Propagate errors output -> hidden  e := e + g*Theta
			copy(syn1copy, *syn1)
			syn1copy.Multiply(g)
			neu1e.Add(syn1copy)
			// Learn weights hidden -> output   Theta := Theta + gV(w)
			if !infer {
				copy(neu1copy, *rangevec)
				neu1copy.Multiply(g)
				syn1.Add(neu1copy)
			}
		}
	}

	//Negative Sampling
	if p.UseNEG {
		label := 0.0
		var point int32 = 0
		for i := 0; i < p.Negative+1; i++ {
			if i == 0 {
				label = 1.0
				point = centralwidx
			} else {
				label = 0.0
				point = GetNegativeSamplingWordIdx()
				if point == centralwidx {
					continue
				}
			}
			syn1neg := p.NN.GetSyn1Neg(point)
			f := rangevec.Dot(*syn1neg) // V(w') dot Theta(u),  u = {w} U NEG(w), w' = Context(w)
			g := 0.0
			if math.IsNaN(f) {
				log.Fatal("f is NaN, try to reduce StartAlpha!")
			}
			if f >= MAX_EXP {
				f = 1.0
			} else if f <= -MAX_EXP {
				f = 0.0
			} else {
				f = GetSigmoidValue(f)
			}
			g = (label - f) * alpha // (Lw(u) - q) * alpha
			// Propagate errors output -> hidden  e := e + g*Theta(u)
			copy(syn1copy, *syn1neg)
			syn1copy.Multiply(g)
			neu1e.Add(syn1copy)
			// Learn weights hidden -> output   Theta := Theta + gV(w')
			if !infer {
				copy(neu1copy, *rangevec)
				neu1copy.Multiply(g)
				syn1neg.Add(neu1copy)
			}
		}
	}

	// hidden -> in                         v(u) := v(u) + e
	rangevec.Add(neu1e)
}

func (p *TDoc2VecImpl) trainSkipGram4Document(wordsidx []int32, dsyn0 *neuralnet.TVector, alpha float64, infer bool) {
	for spos, widx := range wordsidx {
		//随机窗口大小
		b := p.getRandomWindowSize()
		if infer {
			b = 0
		}
		start := common.Max(0, spos-p.WindowSize+b)
		end := common.Min(len(wordsidx), spos+p.WindowSize-b+1)
		for a := start; a < end; a++ {
			if a == spos {
				continue
			}
			idx := wordsidx[a]
			rangevec := p.NN.GetSyn0(idx)
			p.trainSkipGram4Pair(widx, rangevec, alpha, infer)
		}
		//训练doc向量
		p.trainSkipGram4Pair(widx, dsyn0, alpha, infer)
	}
}

//dsyn0 由参数传入是为了方便在infer_doc的时候直接传入一个dvector来进行训练
//infer=true的时候不对模型参数进行更新

func (p *TDoc2VecImpl) trainCbow4Document(wordsidx []int32, dsyn0 *neuralnet.TVector, alpha float64, infer bool) {
    //neu1 := make(neuralnet.TVector, p.Dim, p.Dim)     //X(w)
    //neu1e := make(neuralnet.TVector, p.Dim, p.Dim)    //e
    //syn1copy := make(neuralnet.TVector, p.Dim, p.Dim) //为了计算 g*Theta
    //neu1copy := make(neuralnet.TVector, p.Dim, p.Dim) //为了计算 g*X(w)

    // 使用内存池来降低GC的压力
    neu1 := *(p.Pool.Get().(*neuralnet.TVector))
    neu1e := *(p.Pool.Get().(*neuralnet.TVector))
    syn1copy := *(p.Pool.Get().(*neuralnet.TVector))
    neu1copy := *(p.Pool.Get().(*neuralnet.TVector))

    defer func(){ p.Pool.Put(&neu1) }()
    defer func(){ p.Pool.Put(&neu1e) }()
    defer func(){ p.Pool.Put(&syn1copy) }()
    defer func(){ p.Pool.Put(&neu1copy) }()

	for spos, widx := range wordsidx {
        neu1.Reset()
        neu1e.Reset()
        syn1copy.Reset()
        neu1copy.Reset()
		b := p.getRandomWindowSize()
		if infer {
			b = 0
		}
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
				if math.IsNaN(f) {
					log.Fatal("f is NaN, try to reduce StartAlpha!")
				}
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
				g = (1.0 - label - f) * alpha
				// Propagate errors output -> hidden  e := e + g*Theta
				copy(syn1copy, *syn1)
				syn1copy.Multiply(g)
				neu1e.Add(syn1copy)
				// Learn weights hidden -> output   Theta := Theta + gX(w)
				// when predict doc vector, infer=true, don't update Theta
				if !infer {
					copy(neu1copy, neu1)
					neu1copy.Multiply(g)
					syn1.Add(neu1copy)
				}
			}
		}

		//Negative Sampling
		if p.UseNEG {
			label := 0.0
			var point int32 = 0
			for i := 0; i < p.Negative+1; i++ {
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
				if math.IsNaN(f) {
					log.Fatal("f is NaN, try to reduce StartAlpha!")
				}
				if f >= MAX_EXP {
					f = 1.0
				} else if f <= -MAX_EXP {
					f = 0.0
				} else {
					f = GetSigmoidValue(f)
				}
				g := (label - f) * alpha
				// Propagate errors output -> hidden  e := e + g*Theta
				copy(syn1copy, *syn1neg)
				syn1copy.Multiply(g)
				neu1e.Add(syn1copy)
				// Learn weights hidden -> output   Theta := Theta + gX(w)
				if !infer {
					copy(neu1copy, neu1)
					neu1copy.Multiply(g)
					syn1neg.Add(neu1copy)
				}
			}
		}

		// hidden -> in                         v(u) := v(u) + e
		if !infer {
			for a := start; a < end; a++ {
				if a == spos {
					continue
				}
				idx := wordsidx[a]
				syn0 := p.NN.GetSyn0(idx)
				syn0.Add(neu1e)
			}
		}
		// hidden -> in                         D(u) := D(u) + e
		dsyn0.Add(neu1e)
	}
}

func (p *TDoc2VecImpl) trainSkipGram() {
	// Skip-Gram  Model
	tokens := make(chan struct{}, THREAD_NUM)
	last_trained_words := 0
	alpha := p.getTrainAlpha()
	stime := time.Now()
	for i := 0; i < p.Iters; i++ {
		wg := new(sync.WaitGroup)
		for docidx_, wordsidx_ := range p.Corpus.GetAllDocWordsIdx() {
			docidx, wordsidx := docidx_, wordsidx_
            tokens <- struct{}{}
			wg.Add(1)
			go func() {
				defer func() { <-tokens }()
				defer wg.Done()
				//train one document
				last_trained_words += len(wordsidx)
				p.TrainedWords += len(wordsidx)
				if last_trained_words > PROGRESS_BAR_THRESHOLD {
					last_trained_words = 0
					alpha = p.getTrainAlpha()
					fmt.Printf("%cSkip-Gram Iter:%v Alpha: %f  Progress: %.2f%%  Words/sec: %.2fk  ", 13, i, alpha,
						float64(p.TrainedWords)/float64(p.Iters*p.Corpus.GetWordsCnt()+1)*100,
						float64(p.TrainedWords)/float64(time.Since(stime))*100*1000)
				}
				dsyn0 := p.NN.GetDSyn0(int32(docidx))
				p.trainSkipGram4Document(wordsidx, dsyn0, alpha, false)
			}()
		}
		wg.Wait()
	}
	fmt.Printf("\n%v training end, %v %v\n", time.Now(), p.TrainedWords, p.Corpus.GetWordsCnt())
}

// P(w|Context(w))
func (p *TDoc2VecImpl) trainCbow() {
	//Continuous Bag-of-Word Model
	tokens := make(chan struct{}, THREAD_NUM)
	last_trained_words := 0
	alpha := p.getTrainAlpha()
	stime := time.Now()
	for i := 0; i < p.Iters; i++ {
		wg := new(sync.WaitGroup)
		for docidx_, wordsidx_ := range p.Corpus.GetAllDocWordsIdx() {
			docidx, wordsidx := docidx_, wordsidx_
            tokens <- struct{}{}
			wg.Add(1)
			go func() {
				defer func() { <-tokens }()
				defer wg.Done()
				//train one document
				last_trained_words += len(wordsidx)
				p.TrainedWords += len(wordsidx)
				if last_trained_words > PROGRESS_BAR_THRESHOLD {
					last_trained_words = 0
					alpha = p.getTrainAlpha()
					fmt.Printf("%cCBOW Iter:%v Alpha: %f  Progress: %.2f%%  Words/sec: %.2fk  ", 13, i, alpha,
						float64(p.TrainedWords)/float64(p.Iters*p.Corpus.GetWordsCnt()+1)*100,
						float64(p.TrainedWords)/float64(time.Since(stime))*100*1000)
				}
				dsyn0 := p.NN.GetDSyn0(int32(docidx))
				p.trainCbow4Document(wordsidx, dsyn0, alpha, false)
			}()
		}
		wg.Wait()
	}
	fmt.Printf("\n%v training end, %v %v\n", time.Now(), p.TrainedWords, p.Corpus.GetWordsCnt())
}

func (p *TDoc2VecImpl) GetLeaveOneOutKwds(content string, iters int) {
	vec1 := p.FitDoc(content, iters)
	wordsidx := p.Corpus.Transform(content)
	dis_vector := make(TSortItemSlice, 0, len(wordsidx))
	for i, widx := range wordsidx {
		loowordsidx := make([]int32, 0, len(wordsidx))
		for j, idx := range wordsidx {
			if j == i {
				continue
			}
			loowordsidx = append(loowordsidx, idx)
		}
		vec2 := p.fitDoc(loowordsidx, iters)
		dis := ConsineDistance(*vec1, *vec2)
		sortitem := &SortItem{Idx: widx, Dis: dis}
		dis_vector = append(dis_vector, sortitem)
	}
	sort.Sort(sort.Reverse(dis_vector))
	p.PrintTopKWords(dis_vector)
}

func (p *TDoc2VecImpl) findKNNWordsByVector(vector *neuralnet.TVector) {
	vocabsize := p.Corpus.GetVocabCnt()
	dis_vector := make(TSortItemSlice, vocabsize, vocabsize)
	for i := 0; i < vocabsize; i++ {
		dis := ConsineDistance(*vector, *p.NN.GetSyn0(int32(i)))
		dis_vector[i] = &SortItem{Idx: int32(i), Dis: dis}
	}
	//大爷的 go的排序太麻烦,还不如自己写个快排
	//就不能学学好,跟python一样 sort(dis_vector, key=lambda item: item.key, reverse=True)
	//都已经有接口了, 待排序的元素实现一个包含GetSortKey函数的接口, 就ok了,能省不少代码
	QuickSort(0, len(dis_vector)-1, dis_vector)

	//不用自己写的快排启用这一行也ok, golang 排序需要实现sort.Interface接口
	//  共三个函数需要实现
	//      1.  Len
	//      2.  Less
	//      3.  Swap
	//sort.Sort(sort.Reverse(dis_vector))
	p.PrintTopKWords(dis_vector)
}

func (p *TDoc2VecImpl) findKNNDocsByVector(vector *neuralnet.TVector) {
	doccnt := p.Corpus.GetDocCnt()
	dis_vector := make(TSortItemSlice, doccnt, doccnt)
	for i := 0; i < doccnt; i++ {
		dis := ConsineDistance(*vector, *p.NN.GetDSyn0(int32(i)))
		dis_vector[i] = &SortItem{Idx: int32(i), Dis: dis}
	}
	//大爷的 go的排序太麻烦,还不如自己写个快排
	//就不能学学好,跟python一样 sort(dis_vector, key=lambda item: item.key, reverse=True)
	//都已经有接口了, 待排序的元素实现一个包含GetSortKey函数的接口, 就ok了,能省不少代码
	QuickSort(0, len(dis_vector)-1, dis_vector)

	//不用自己写的快排启用这一行也ok, golang 排序需要实现sort.Interface接口
	//  共三个函数需要实现
	//      1.  Len
	//      2.  Less
	//      3.  Swap
	//sort.Sort(sort.Reverse(dis_vector))
	p.PrintTopKDocs(dis_vector)
}

func (p *TDoc2VecImpl) Word2Words(word string) {
	idx, ok := p.Corpus.GetWordIdx(word)
	if !ok {
		return
	}
	vector := p.NN.GetSyn0(idx)
	p.findKNNWordsByVector(vector)
}

func (p *TDoc2VecImpl) Sen2Words(content string, iters int) {
	vector := p.FitDoc(content, iters)
	p.findKNNWordsByVector(vector)
}

func (p *TDoc2VecImpl) Doc2Words(docidx int) {
	if docidx < 0 || docidx > p.Corpus.GetDocCnt() {
		return
	}
	vector := p.NN.GetDSyn0(int32(docidx))
	worditems := p.Corpus.GetDocWordsByIdx(docidx)
	words := []string{}
	for _, item := range worditems {
		words = append(words, item.Word)
	}
	content := strings.Join(words, " ")
	fmt.Println(content)
	p.findKNNWordsByVector(vector)
}

func (p *TDoc2VecImpl) Doc2Docs(docidx int) {
	if docidx < 0 || docidx > p.Corpus.GetDocCnt() {
		return
	}
	vector := p.NN.GetDSyn0(int32(docidx))
	worditems := p.Corpus.GetDocWordsByIdx(docidx)
	words := []string{}
	for _, item := range worditems {
		words = append(words, item.Word)
	}
	content := strings.Join(words, " ")
	fmt.Println(content)
	p.findKNNDocsByVector(vector)
}

func (p *TDoc2VecImpl) Sen2Docs(content string, iters int) {
	vector := p.FitDoc(content, iters)
	p.findKNNDocsByVector(vector)
}

func (p *TDoc2VecImpl) Word2Docs(word string) {
	idx, ok := p.Corpus.GetWordIdx(word)
	if !ok {
		return
	}
	vector := p.NN.GetSyn0(idx)
	p.findKNNDocsByVector(vector)
}

func (p *TDoc2VecImpl) PrintTopKWords(slice TSortItemSlice) {
	for i, item := range slice {
		if i >= 10 {
			break
		}
		dis := item.Dis
		idx := int(item.Idx)
		worditem := p.Corpus.GetWordItemByIdx(idx)
		//fmt.Printf("\t%v\t%v\t%v\n", dis, worditem.Word, *p.NN.GetSyn0(int32(item.Idx)))
		fmt.Printf("\t%v\t%v\n", dis, worditem.Word)
	}
	fmt.Println()
}

func (p *TDoc2VecImpl) PrintTopKDocs(slice TSortItemSlice) {
	for i, item := range slice {
		if i >= 10 {
			break
		}
		dis := item.Dis
		idx := int(item.Idx)
		worditems := p.Corpus.GetDocWordsByIdx(idx)
		words := []string{}
		for _, item := range worditems {
			words = append(words, item.Word)
		}
		content := strings.Join(words, " ")
		//fmt.Printf("\t%v\t%v\t%v\n", dis, content, *p.NN.GetDSyn0(int32(item.Idx)))
		fmt.Printf("\t%v\t%v\n", dis, content)
	}
	fmt.Println()
}

//升序快排
func QuickSort(i, j int, vec []*SortItem) {
	ii, jj := i, j
	if i+1 >= j {
		return
	} else if i+2 == j {
		if vec[i].Dis < vec[j-1].Dis {
			vec[i], vec[j-1] = vec[j-1], vec[i]
		}
	}
	M := i
	stub := vec[M]
	for i < j {
		for ; j > i; j-- {
			if vec[j].Dis > stub.Dis {
				vec[M] = vec[j]
				M = j
				break
			}
		}
		for ; i < j; i++ {
			if vec[i].Dis < stub.Dis {
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
