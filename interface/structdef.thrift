namespace go structdef


typedef list<double> TVector

struct TNeuralNet {
    1:bool                  UseHS,
    2:bool                  UseNEG,
    3:i64                   DocSize,
    4:i64                   VocabSize,
    5:i64                   WordSize,
    6:i64                   Dim,
    7:list<TVector>         Syn0,
    8:list<TVector>         Dsyn0,
    9:list<TVector>         Syn1,
    10:list<TVector>        Syn1neg,
}

struct TWordItem {
    1:i32   Cnt,
    2:list<i32>   Point,
    3:list<bool>  Code,
    4:string      Word,
}

typedef list<TWordItem> TWordItemSlice
struct TCorpus {
    1:TWordItemSlice Words,
    2:map<string, i32>  WordsIdx,
    3:list<list<i32>>   DocWordsIdx,
    4:map<string, i32>  DocIdx,
    5:i32               MinReduce,
    6:i32               MinCnt,
    7:i32               WordsCnt,
}
