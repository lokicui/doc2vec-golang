package corpus

import (
	"bufio"
	_ "errors"
	"fmt"
    "github.com/lokicui/doc2vec-golang/common"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	VOCAB_HASH_SIZE int = 30000000   //3kw, 30M
)

func (p *TCorpusImpl) GetWordIdx(word string) (idx int32, ok bool) {
	idx, ok = p.Word2Idx[word]
	return idx, ok
}

func (p *TCorpusImpl) GetWordsCnt() int {
	return p.WordsCnt
}

func (p *TCorpusImpl) createBinaryTree() {
	vocab_size := p.GetVocabCnt()
	size := vocab_size*2 + 1
	cnt := make([]int64, size, size)
	parent := make([]int32, size, size)
	binary := make([]bool, size, size)
	//p.GetAllWords 一定要降序排列 p.sortVocab
	for i, item := range p.GetAllWords() {
		cnt[i] = int64(item.Cnt)
	}
	for i := vocab_size; i < vocab_size*2; i++ {
		cnt[i] = math.MaxInt64
	}
	pos1 := vocab_size - 1 //初始化为min
	pos2 := vocab_size
	min1, min2 := 0, 0                  //最小和次小
	for i := 0; i < vocab_size-1; i++ { //vocab_size 个非叶结点
		if pos1 >= 0 {
			if cnt[pos1] < cnt[pos2] {
				min1 = pos1
				pos1--
			} else {
				min1 = pos2
				pos2++
			}
		} else {
			min1 = pos2
			pos2++
		}

		if pos1 >= 0 {
			if cnt[pos1] < cnt[pos2] {
				min2 = pos1
				pos1--
			} else {
				min2 = pos2
				pos2++
			}
		} else {
			min2 = pos2
			pos2++
		}
		cnt[vocab_size+i] = cnt[min1] + cnt[min2]
		parent[min1] = int32(vocab_size + i)
		parent[min2] = int32(vocab_size + i)
		binary[min2] = true
	}
	root := int32(vocab_size + vocab_size - 1 - 1)
	for i := 0; i < vocab_size; i++ {
		code := make([]bool, 0, 40)
		point := make([]int32, 0, 40)
		for p := int32(i); p != root; p = parent[p] {
			if p >= int32(vocab_size) {
				point = append(point, p-int32(vocab_size)) //转换为syn1的下标
			} else {
				//NOTE 叶子节点、HS训练的时候不需要, HS的Point只需要[root, leaf)
				//point = append(point, p)
			}
			code = append(code, binary[p])
		}
		point = append(point, root-int32(vocab_size))

		reverse_point := make([]int32, len(point), len(point))
		for j := 0; j < len(point); j++ {
			reverse_point[len(point)-1-j] = point[j]
		}

		reverse_code := make([]bool, len(code), len(code))
		for j := 0; j < len(code); j++ {
			reverse_code[len(code)-1-j] = code[j]
		}

		//////////////////
		//ii := int32(i)
		//for true {
		//	point = append(point, int32(ii))
		//	code = append(code, binary[ii])
		//	p := parent[ii]
		//	if p == root {
		//		break
		//	}
		//	ii = p
		//}
		////code point reverse
		//reverse_code := make([]bool, 0, len(code))
		//for j := len(code) - 1; j >= 0; j-- {
		//	reverse_code = append(reverse_code, code[j])
		//}
		//reverse_point := make([]int32, 0, len(point))
		////root node index
		//reverse_point = append(reverse_point, root-int32(vocab_size)) // 加上root, 减去vocab_size后直接表示syn1中的下标了
		//for j := len(point) - 1; j > 0; j-- {                         //NOTE  j != 0, j == 0是叶子节点、HS训练的时候不需要了, HS的Point只需要[root, leaf)
		//	syn1idx := point[j] - int32(vocab_size)
		//	reverse_point = append(reverse_point, syn1idx)
		//}
		//(root->leaf]
		p.Words[i].Code = reverse_code
		p.Words[i].Point = reverse_point
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

func (p *TCorpusImpl) GetDocCnt() int {
	return len(p.Doc2Idx)
}

func (p *TCorpusImpl) GetVocabCnt() int {
	return len(p.GetAllWords())
}

func (p *TCorpusImpl) GetAllWords() (words TWordItemSlice) {
	return p.Words
}

func (p *TCorpusImpl) GetDocWordsByDocid(id string) (doc []*TWordItem) {
	idx, ok := p.Doc2Idx[id]
	if ok {
		return p.GetDocWordsByIdx(int(idx))
	}
	return nil
}

func (p *TCorpusImpl) GetAllDocWordsIdx() [][]int32 {
	return p.Doc2WordsIdx
}

func (p *TCorpusImpl) GetAllDocWords() (docs [][]*TWordItem) {
	for i := 0; i < len(p.Doc2WordsIdx); i++ {
		docs = append(docs, p.GetDocWordsByIdx(i))
	}
	return docs
}

func (p *TCorpusImpl) GetWordItemByIdx(i int) (item *TWordItem) {
	if i < 0 || i >= len(p.Words) {
		log.Fatal("index out of range")
	}
	return &p.Words[i]
}

func (p *TCorpusImpl) GetDocWordsByIdx(i int) (doc []*TWordItem) {
	if i >= 0 && i < len(p.Doc2WordsIdx) {
		doc = make([]*TWordItem, len(p.Doc2WordsIdx[i]), len(p.Doc2WordsIdx[i]))
		for j, idx := range p.Doc2WordsIdx[i] {
			if idx >= 0 && idx < int32(len(p.Words)) {
				doc[j] = &p.Words[idx]
			}
		}
	}
	return doc
}

func (p *TCorpusImpl) reduceVocabulary() {
	var idx int32 = 0
	actual_size := len(p.Word2Idx)
	p.Word2Idx = map[string]int32{}
	for i, item := range p.Words {
		if i == actual_size {
			break
		}
		if item.Cnt > p.MinReduce {
			p.Words[idx] = item
			p.Word2Idx[item.Word] = idx
			idx++
		}
	}
	p.Words = p.Words[:len(p.Word2Idx)]
	p.MinReduce++
}

func (p *TCorpusImpl) addWord(word string) (err error) {
	if len(word) == 0 {
		return err
	}
	idx, ok := p.Word2Idx[word]
	if !ok {
		item := TWordItem{Word: word, Cnt: 0}
		idx = int32(len(p.Words))
		p.Words = append(p.Words, item)
		p.Word2Idx[word] = idx
	}
	p.Words[idx].Cnt++
	return err
}

func (p *TCorpusImpl) loadAsWords(docid string, content string) int {
	items := strings.Split(content, " ")
	for _, word := range items {
        word = common.SBC2DBC(word)
		p.addWord(word)
		if len(p.Word2Idx) > int(0.7*float32(VOCAB_HASH_SIZE)) {
            log.Printf("%d > %d, start reduceVocabulary\n", len(p.Word2Idx), VOCAB_HASH_SIZE)
			p.reduceVocabulary()
		}
	}
	return len(items)
}

func (p *TCorpusImpl) Transform(content string) (wordsidx []int32) {
	items := strings.Split(content, " ")
	for _, word := range items {
        word = common.SBC2DBC(word)
		idx, ok := p.Word2Idx[word]
		if ok {
			wordsidx = append(wordsidx, int32(idx))
		}
	}
	return wordsidx
}

func (p *TCorpusImpl) loadAsDoc(docid string, content string) int {
	items := strings.Split(content, " ")
	wordsIdx := make([]int32, 0, len(items))
	for _, word := range items {
        word = common.SBC2DBC(word)
		idx, ok := p.Word2Idx[word]
		if ok {
			wordsIdx = append(wordsIdx, int32(idx))
		}
	}
    if idx, ok := p.Doc2Idx[docid]; ok {
        p.Doc2WordsIdx[idx] = wordsIdx   // exists, update
    } else {
        p.Doc2WordsIdx = append(p.Doc2WordsIdx, wordsIdx)
        p.Doc2Idx[docid] = int32(len(p.Doc2WordsIdx) - 1)
    }
	return 1
}

func (p *TCorpusImpl) String() string {
	words_cnt := strconv.Itoa(len(p.Words))
	words_map_cnt := strconv.Itoa(len(p.Word2Idx))
	docs_cnt := strconv.Itoa(len(p.Doc2WordsIdx))
	return fmt.Sprintf("words_cnt:%v,words_map_cnt:%v,docs_cnt:%v\n", words_cnt, words_map_cnt, docs_cnt)
}

func (p *TCorpusImpl) buildVocabulary(fname string) (err error) {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Buffer([]byte{}, bufio.MaxScanTokenSize*100)
	train_words := 0
	batch := 0
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(line, "\t")
		if len(items) < 2 {
			log.Printf("len(items)=%d\n", len(items))
			continue
		}
		docid, content := items[0], items[1]
		cnt := p.loadAsWords(docid, content)
		train_words += cnt
		batch += cnt
		if batch >= 10000000 {
			batch = 0
            log.Printf("train %d words, vocab_size:%d\n", train_words, p.GetVocabCnt())
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	p.sortVocab()
	p.createBinaryTree()
	return err
}

func (p *TCorpusImpl) sortVocab() {
	//先排序
	// Words occuring less than min_count times will be discarded from the vocab
	p.Word2Idx = map[string]int32{}
	sort.Sort(sort.Reverse(p.Words))

	p.WordsCnt = 0
	var idx int32 = 0
	for _, item := range p.Words {
		if item.Cnt > p.MinCnt {
			p.Words[idx] = item
			p.Word2Idx[item.Word] = idx
			p.WordsCnt += int(item.Cnt)
			idx++
		}
	}
	p.Words = p.Words[:idx]
}

func (p *TCorpusImpl) loadDocument(fname string) (err error) {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Buffer([]byte{}, bufio.MaxScanTokenSize*100)
	train_docs := 0
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(line, "\t")
		if len(items) < 2 {
			log.Printf("len(items)=%d\n", len(items))
			continue
		}
		docid, content := items[0], items[1]
		docid = strings.Trim(docid, " \n")
		content = strings.Trim(content, " \n")
		if len(content) == 0 || len(docid) == 0 {
			continue
		}
		cnt := p.loadAsDoc(docid, content)
		train_docs += cnt
		if train_docs%100000 == 0 {
            log.Printf("train %d docs, doc_size:%d\n", train_docs, p.GetDocCnt())
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return err
}

func (p *TCorpusImpl) Build(fname string) (err error) {
	err = p.buildVocabulary(fname)
	if err != nil {
		return err
	}
	return p.loadDocument(fname)
}

func NewCorpus() ICorpus {
	self := &TCorpusImpl{
		Word2Idx: make(map[string]int32),
		Doc2Idx:  make(map[string]int32),
        MinReduce: 1,
		MinCnt:   1,
    }
	return ICorpus(self)
}
