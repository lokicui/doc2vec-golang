namespace go structdef


typedef list<double> TVector

struct TNeuralNet {
    1:list<TVector>         Syn0,
    2:list<TVector>         Dsyn0,
    3:list<TVector>         Syn1,
    4:list<TVector>         Syn1neg,
}

struct TWordItem {
    1:i32                   Cnt,        //term frequency
    2:list<i32>             Point,      //Huffman tree(n leaf + n inner node, include root) path. [root, leaf), node index
    3:list<bool>            Code,       //Huffman code. (root, leaf], 0/1 codes
    4:string                Word,       //word desc
}

typedef list<TWordItem> TWordItemSlice
struct TCorpus {
    1:TWordItemSlice        Words,
    2:map<string, i32>      Words2Idx,
    3:list<list<i32>>       Doc2WordsIdx,
    4:map<string, i32>      Doc2Idx,
    5:i32                   MinReduce,
    6:i32                   MinCnt,
    7:i32                   WordsCnt,
}

struct TDoc2vec {
    1:bool                  UseHS,          //层次Softmax
    2:bool                  UseNEG,         //负采样
    3:i64                   DocSize,        //训练文档数
    4:i64                   VocabSize,      //词表大小,unique的词个数
    5:i64                   WordSize,       //所有词个数,不排重
    6:i64                   Dim,            //隐层的维数
    7:i64                   WindowSize,     //窗口大小
    8:i64                   Iters,          //迭代次数
    9:double                StartAlpha,     //learning rate
    10:i64                  TrainedWords,   //已经train过的词数
    11:string               TrainFile,      //训练文件 每行两列 第一列是docid，第二列是逗号分割的词
    12:TNeuralNet           NN,
    13:TCorpus              corpus,
}
