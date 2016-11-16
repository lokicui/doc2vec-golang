package main

import (
    "fmt"
    "os"
    "sort"
    "log"
    "math"
    "sync"
    "strings"
    "corpus"
    "neuralnet"
    "common"
    "bytes"
    "bufio"
    "encoding/binary"
)

const (
    MAX_EXP float64 = 6
    EXP_TABLE_SIZE  int = 1000
)

type SortItem struct {
    Idx int32
    Dis float64
}
type TSortItemSlice []*SortItem

func (p TSortItemSlice) Len() int {
    return len(p)
}
func (p TSortItemSlice) Less(i,j int) bool {
    return p[i].Dis < p[j].Dis
}
func (p TSortItemSlice) Swap(i,j int) {
    p[i], p[j] = p[j], p[i]
}

type Doc2Vec interface {
    Train(fname string)
    GetCorpus() corpus.Corpus
    GetNN()     neuralnet.NN
    SaveModel(fname string) (err error)
    LoadModel(fname string) (err error)
    FindKNN(word string)
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
    startAlpha       float64
    iter        int
    trained_words int
    expTable    [EXP_TABLE_SIZE]float64
    words       []string
    vocabSize   int
}

func NewDoc2Vec(useHS, useNEG bool, windowSize, dim int) Doc2Vec {
    self := &Doc2VecImpl{
        useHS: useHS,
        useNEG: useNEG,
        windowSize: windowSize,
        dim: dim,
        corpus: nil,
        nn: nil,
        startAlpha: 0.025,
        iter: 50,
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
    for i := 0; i < p.iter; i ++ {
        p.TrainHSCbow()
    }
}

func (p* Doc2VecImpl) SaveModel(fname string) (err error) {
    fd, err := os.Create(fname)
    if err != nil {
        log.Fatal(err)
    }
    defer fd.Close()

    buf := new(bytes.Buffer)
    err = binary.Write(buf, binary.BigEndian, int32(p.corpus.GetVocabSize()))
    if err != nil {
        log.Fatal(err)
    }
    err = binary.Write(buf, binary.BigEndian, int32(p.corpus.GetDocSize()))
    if err != nil {
        log.Fatal(err)
    }
    err = binary.Write(buf, binary.BigEndian, int32(p.dim))
    if err != nil {
        log.Fatal(err)
    }

    for i := 0; i < p.corpus.GetVocabSize(); i ++ {
        item := p.corpus.GetWordItemByIdx(i)
        err = binary.Write(buf, binary.BigEndian, int32(len(item.Word)))
        if err != nil {
            log.Fatal(err)
        }
        buf.WriteString(item.Word)
    }
    buf.WriteTo(fd)
    buf.Truncate(0)
    for i := 0; i < p.corpus.GetVocabSize(); i ++ {
        vector := p.nn.GetSyn0(int32(i))
        err = binary.Write(buf, binary.BigEndian, *vector)
        if err != nil {
            log.Fatal(err)
        }
        if i % 1000 == 0 {
            fmt.Printf("save %v words vector\n", i)
        }
    }
    buf.WriteTo(fd)
    buf.Truncate(0)
    for i := 0; i < p.corpus.GetDocSize(); i ++ {
        vector := p.nn.GetDSyn0(int32(i))
        err = binary.Write(buf, binary.BigEndian, *vector)
        if err != nil {
            log.Fatal(err)
        }
        if i % 100 == 0 {
            fmt.Printf("save %v words doc\n", i)
        }
    }
    buf.WriteTo(fd)
    buf.Truncate(0)
    return err
}

func (p * Doc2VecImpl) LoadModel(fname string) (err error) {
    fd, err := os.Open(fname)
    if err != nil {
        log.Fatal(err)
    }
    defer fd.Close()
    var vocabSize, docSize, dim int32
    err = binary.Read(fd, binary.BigEndian, &vocabSize)
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
    p.dim = int(dim)
    p.vocabSize = int(vocabSize)
    p.words = make([]string, vocabSize, vocabSize)
    for i := 0; i < int(vocabSize); i ++ {
        var size int32
        binary.Read(fd, binary.BigEndian, &size)
        bytes := make([]byte, size, size)
        fd.Read(bytes)
        p.words[i] = string(bytes)
    }
    p.nn = neuralnet.NewNN(int(docSize), int(vocabSize), p.dim, false, false)
    for i := 0; i < int(vocabSize); i ++ {
        vector := p.nn.GetSyn0(int32(i))
        for j :=0; j < p.dim; j ++ {
            binary.Read(fd, binary.BigEndian, &(*vector)[j])
        }
    }
    for i := 0; i < int(docSize); i ++ {
        vector := p.nn.GetDSyn0(int32(i))
        for j :=0; j < p.dim; j ++ {
            binary.Read(fd, binary.BigEndian, &(*vector)[j])
        }
    }
    return err
}

func (p * Doc2VecImpl) GetAlpha() float64 {
    alpha := p.startAlpha * (1.0 - float64(p.trained_words) / float64(p.iter * p.corpus.GetWordsCnt() + 1))
    if alpha < p.startAlpha * 0.0001 {
        alpha = p.startAlpha * 0.0001
    }
    return alpha
}


func (p * Doc2VecImpl) TrainHSCbow() {
    //Continuous Bag-of-Word Model + Hierarchical Softmax
    tokens := make(chan struct {}, 32)
    wg := new(sync.WaitGroup)
    last_trained_words := 0
    alpha := p.GetAlpha()
    for docidx_, wordsidx_ := range p.corpus.GetAllDocWordsIdx() {
        docidx, wordsidx := docidx_, wordsidx_
        wg.Add(1)
        go func() {
            defer func() {<-tokens}()
            defer wg.Done()
            tokens <- struct {}{}
            //train one document
            last_trained_words += len(wordsidx)
            p.trained_words += len(wordsidx)
            if last_trained_words > 10000 {
                last_trained_words = 0
                alpha = p.GetAlpha()
            }
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
                        g = -1.0 * f * alpha
                    } else {
                        g = (1.0 - f) * alpha
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

func (p* Doc2VecImpl) FindKNN(word string) {
    idx := 0
    for i, w := range p.words {
        if w == word {
            idx = i
            break
        }
    }
    vector := p.nn.GetSyn0(int32(idx))
    vocab_size := int(p.vocabSize)
    dis_vector := make(TSortItemSlice, vocab_size, vocab_size)
    for i := 0; i < vocab_size; i ++ {
        dis := p.ConsineDistance(*vector, *p.nn.GetSyn0(int32(i)))
        dis_vector[i] = &SortItem{Idx:int32(i), Dis:dis}
    }
    //大爷的 go的排序太麻烦,还不如自己写个快排
    sort.Sort(sort.Reverse(dis_vector))
    for i := 0; i < len(dis_vector) && i < 10; i ++ {
        item := dis_vector[i]
        dis := item.Dis
        fmt.Printf("%v\t%v\t%v\n", dis, p.words[int(item.Idx)], *p.nn.GetSyn0(int32(item.Idx)))
    }
}

func QuickSort(i,j int, vec [] *SortItem) {
    ii, jj := i, j
    if i + 1 == j {
        return
    } else if i + 2 == j {
        if vec[i].Dis > vec[j-1].Dis {
            vec[i], vec[j-1] = vec[j-1], vec[i]
        }
    }
    M := i
    stub := vec[M]
    for ; i < j; {
        for ; j > i; j -- {
            if vec[j].Dis < stub.Dis {
                vec[M] = vec[j]
                M = j
                break
            }
        }
        for ; i < j; i ++ {
            if vec[i].Dis > stub.Dis {
                vec[M] = vec[i]
                M = i
                break
            }
        }
    }
    vec[M] = stub
    QuickSort(ii, M, vec)
    QuickSort(M + 1, jj, vec)
}


func (p *Doc2VecImpl) ConsineDistance(a neuralnet.TVector, b neuralnet.TVector) (dis float64) {
    var sum, sum_a, sum_b float64
    for i := 0; i < len(a); i ++ {
        sum += float64(a[i] * b[i])
        sum_a += float64(a[i] * a[i])
        sum_b += float64(b[i] * b[i])
    }
    dis = sum / math.Sqrt(sum_a) / math.Sqrt(sum_b)
    return dis
}

func (p* Doc2VecImpl) TrainHSSkipGram() {
    //Skip-Gram Model + Hierarchical Softmax
}

func main() {
    //fname := os.Args[1]
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
    //doc2vec.Train(fname)
    //doc2vec.SaveModel("3.model")
    doc2vec.LoadModel("50.model")
    for true {
        reader := bufio.NewReader(os.Stdin)
        fmt.Println("Enter text:")
        text, _ := reader.ReadString('\n')
        doc2vec.FindKNN(strings.Trim(text, "\n"))
    }
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
