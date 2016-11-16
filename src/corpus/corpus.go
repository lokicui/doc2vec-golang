package corpus

import (
    "fmt"
    "os"
    "log"
    "sort"
    "math"
    "bufio"
    "strings"
    "strconv"
    _ "errors"
    "time"
)

const (
    VOCAB_HASH_SIZE int = 30000000
)


type WordItem struct {
    Cnt          int32      //term frequency
    Point       []int32     //Huffman tree(n leaf + n inner node, include root) path. [root, leaf), node index
    Code        []bool      //Huffman code. (root, leaf], 0/1 codes
    Word        string      //word desc
}
type TWordItemSlice [] WordItem

type Corpus interface {
    Build(fname string)             (err error)
    GetVocabSize() int
    GetDocSize() int
    GetWordsCnt() int
    GetWordIdx(word string)         (idx int32, ok bool)
    GetWordItemByIdx(i int)         (item * WordItem)
    GetAllWords()                   (words TWordItemSlice)
    GetAllDocWordsIdx()             ([][]int32)
    GetAllDocWords()                (doc [][]*WordItem)
    GetDocWordsByDocid(id string)   (doc []*WordItem)
    GetDocWordsByIdx(i int)       (doc []*WordItem)
}


type CorpusImpl struct {
    words           TWordItemSlice
    wordIdx         map[string]int32
    docWordsIdx     [][]int32
    docIdx          map[string]int32
    minReduce       int32
    minCnt          int32
    wordsCnt        int
}

func (p *CorpusImpl) GetWordIdx(word string) (idx int32, ok bool) {
    idx, ok = p.wordIdx[word]
    return idx, ok
}

func (p *CorpusImpl) GetWordsCnt() int {
    return p.wordsCnt
}

func (p *CorpusImpl) createBinaryTree() {
    vocab_size := p.GetVocabSize()
    size := vocab_size * 2 + 1
    cnt := make([]int64, size, size)
    parent := make([]int32, size, size)
    binary := make([]bool, size, size)
    //p.GetAllWords 一定要降序排列 p.sortVocab
    for i, item := range p.GetAllWords() {
        cnt[i] = int64(item.Cnt)
    }
    for i := vocab_size; i < vocab_size * 2; i ++ {
        cnt[i] = math.MaxInt64
    }
    pos1 := vocab_size - 1  //初始化为min
    pos2 := vocab_size
    min1, min2 := 0, 0  //最小和次小
    for i := 0; i < vocab_size - 1; i ++ {  //vocab_size 个非叶结点
        if pos1 >= 0 {
            if cnt[pos1] < cnt[pos2] {
                min1 = pos1
                pos1 --
            } else {
                min1 = pos2
                pos2 ++
            }
        } else {
            min1 = pos2
            pos2 ++
        }

        if pos1 >= 0 {
            if cnt[pos1] < cnt[pos2] {
                min2 = pos1
                pos1 --
            } else {
                min2 = pos2
                pos2 ++
            }
        } else {
            min2 = pos2
            pos2 ++
        }
        cnt[vocab_size + i] = cnt[min1] + cnt[min2]
        parent[min1] = int32(vocab_size + i)
        parent[min2] = int32(vocab_size + i)
        binary[min2] = true
    }
    root := int32(vocab_size + vocab_size - 1 - 1)
    for i := 0; i < vocab_size; i ++ {
        code := make([]bool, 0, 40)
        point := make([]int32, 0, 40)
        ii := int32(i)
        for true {
            point = append(point, int32(ii))
            code = append(code, binary[ii])
            p := parent[ii]
            if p == root {
                break
            }
            ii = p
        }
        //code point reverse
        reverse_code := make([] bool, 0, len(code))
        for j := len(code) - 1; j >= 0; j -- {
            reverse_code = append(reverse_code, code[j])
        }
        reverse_point := make([] int32, 0, len(point))
        //root node index
        reverse_point = append(reverse_point, root - int32(vocab_size))  // 加上root, 减去vocab_size后直接表示syn1中的下标了
        for j := len(point) - 1; j > 0; j -- {                  //NOTE  j != 0, j == 0是叶子节点、HS训练的时候不需要了, HS的Point只需要[root, leaf)
            syn1idx := point[j] - int32(vocab_size)
            reverse_point = append(reverse_point, syn1idx)
        }
        //(root->leaf]
        p.words[i].Code = reverse_code
        p.words[i].Point = reverse_point
    }
}

func (p TWordItemSlice) Len() int {
    return len(p)
}

func (p TWordItemSlice) Less(i, j int) bool {
    return p[i].Cnt < p[j].Cnt
}

func (p TWordItemSlice) Swap(i, j int) {
    p[i], p[j] = p[j], p[i]
}

func (p *CorpusImpl) GetDocSize() int {
    return len(p.docIdx)
}

func (p *CorpusImpl) GetVocabSize() int {
    return len(p.GetAllWords())
}

func (p *CorpusImpl) GetAllWords() (words TWordItemSlice) {
    return p.words
}

func (p *CorpusImpl) GetDocWordsByDocid(id string) (doc []*WordItem) {
    idx, ok := p.docIdx[id]
    if ok {
        return p.GetDocWordsByIdx(int(idx))
    }
    return nil
}

func (p *CorpusImpl) GetAllDocWordsIdx() ([][]int32) {
    return p.docWordsIdx
}

func (p *CorpusImpl) GetAllDocWords() (docs [][]*WordItem) {
    for i := 0; i < len(p.docWordsIdx); i ++ {
        docs = append(docs, p.GetDocWordsByIdx(i))
    }
    return docs
}

func (p *CorpusImpl) GetWordItemByIdx(i int) (item * WordItem) {
    if i < 0 || i >= len(p.words) {
        log.Fatal("index out of range")
    }
    return &p.words[i]
}

func (p *CorpusImpl) GetDocWordsByIdx(i int) (doc [] *WordItem) {
    if i >= 0 && i < len(p.docWordsIdx) {
        doc = make([]*WordItem, len(p.docWordsIdx[i]), len(p.docWordsIdx[i]))
        for j, idx := range p.docWordsIdx[i] {
            if idx >= 0 && idx < int32(len(p.words)) {
                doc[j] = &p.words[idx]
            }
        }
    }
    return doc
}

func (p *CorpusImpl) reduceVocabulary() {
    var idx int32 = 0
    actual_size := len(p.wordIdx)
    p.wordIdx = map[string]int32 {}
    for i, item := range p.words {
        if i == actual_size {
            break
        }
        if item.Cnt > p.minCnt {
            p.words[idx] = item
            p.wordIdx[item.Word] = idx
            idx ++
        }
    }
    p.words = p.words[:len(p.wordIdx)]
    p.minReduce ++
}

func (p *CorpusImpl) addWord(word string) (err error) {
    idx, ok := p.wordIdx[word]
    if !ok {
        item := WordItem{Word:word}
        idx = int32(len(p.words))
        p.words = append(p.words, item)
        p.wordIdx[word] = idx
    }
    p.words[idx].Cnt ++
    return err
}

func (p *CorpusImpl) loadAsWords(docid string, content string) int {
    items := strings.Split(content, " ")
    for _, word := range items {
        p.addWord(word)
        if len(p.wordIdx) > int(0.7 * float32(VOCAB_HASH_SIZE)) {
            p.reduceVocabulary()
        }
    }
    return len(items)
}

func (p *CorpusImpl) loadAsDoc(docid string, content string) int {
    items := strings.Split(content, " ")
    wordsIdx := [] int32 {}
    for _, word := range items {
        idx, ok := p.wordIdx[word]
        if ok {
            wordsIdx = append(wordsIdx, int32(idx))
        }
    }
    p.docWordsIdx = append(p.docWordsIdx, wordsIdx)
    p.docIdx[docid] = int32(len(p.docWordsIdx) - 1)
    return 1
}

func (p *CorpusImpl) String() string {
    words_cnt := strconv.Itoa(len(p.words))
    words_map_cnt := strconv.Itoa(len(p.wordIdx))
    docs_cnt := strconv.Itoa(len(p.docWordsIdx))
    return fmt.Sprintf("words_cnt:%v,words_map_cnt:%v,docs_cnt:%v\n", words_cnt, words_map_cnt, docs_cnt)
}

func (p *CorpusImpl) buildVocabulary(fname string) (err error) {
    file, err := os.Open(fname)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    scanner.Buffer([] byte{}, bufio.MaxScanTokenSize * 10)
    train_words := 0
    batch := 0
    for scanner.Scan() {
        line := scanner.Text()
        items := strings.Split(line, "\t")
        if len(items) != 2 {
            log.Printf("len(items)=%d\n", len(items))
            continue
        }
        docid, content := items[0], items[1]
        cnt := p.loadAsWords(docid, content)
        train_words += cnt
        batch += cnt
        if batch >= 10000000 {
            batch = 0
            fmt.Printf("%s train %d words\n", time.Now(), train_words)
        }
    }
    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
    p.sortVocab()
    p.createBinaryTree()
    return err
}

func (p *CorpusImpl) sortVocab() {
    //先排序
    // Words occuring less than min_count times will be discarded from the vocab
    p.wordIdx = map[string]int32 {}
    var cnt int32 = 0
    sort.Sort(sort.Reverse(p.words))

    p.wordsCnt = 0
    for _, item := range p.words {
        if item.Cnt > p.minCnt {
            p.words[cnt] = item
            p.wordIdx[item.Word] = cnt
            p.wordsCnt += int(cnt)
            cnt ++
        }
    }
    p.words = p.words[:cnt]
}

func (p *CorpusImpl) buildDocument(fname string) (err error) {
    file, err := os.Open(fname)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    scanner.Buffer([] byte{}, bufio.MaxScanTokenSize * 10)
    train_docs := 0
    for scanner.Scan() {
        line := scanner.Text()
        items := strings.Split(line, "\t")
        if len(items) != 2 {
            log.Printf("len(items)=%d\n", len(items))
            continue
        }
        docid, content := items[0], items[1]
        cnt := p.loadAsDoc(docid, content)
        train_docs += cnt
        if train_docs % 10000 == 0 {
            fmt.Printf("%s train %d docs\n", time.Now(), train_docs)
        }
    }
    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
    return err
}

func (p *CorpusImpl) Build(fname string) (err error) {
    err = p.buildVocabulary(fname)
    if err != nil {
        return err
    }
    return p.buildDocument(fname)
}

func NewCorpus()  Corpus{
    self := &CorpusImpl{
        wordIdx: make(map[string]int32),
        docIdx: make(map[string]int32),
        minCnt: 10}
    return Corpus(self)
}
