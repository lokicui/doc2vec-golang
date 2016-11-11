package neuralnet

import (
    "log"
    _"fmt"
)

type TVector []float32
//向量加法
func (p TVector) Add(a TVector) {
    if len(p) != len(a) {
        log.Fatal("len(p) != len(a)")
    }
    for i := 0; i < len(p); i ++ {
        p[i] += a[i]
    }
}

//向量除法
func (p TVector) Divide(a float32) {
    for i := 0; i < len(p); i ++ {
        p[i] /= a
    }
}

//向量点乘
func (p TVector) Dot(a TVector) (f float64) {
    if len(p) != len(a) {
        log.Fatal("len(p) != len(a)")
    }
    for i := 0; i < len(p); i ++ {
        f += float64(p[i]) * float64(a[i])
    }
    return f
}

func (p TVector) Multiply(a float64) {
    for i := 0; i < len(p); i ++ {
        p[i] *= float32(a)
    }
}

type NN interface {
    GetSyn0(i int32) (*TVector)
    GetDSyn0(i int32) (*TVector)
    GetSyn1(i int32) (*TVector)
    GetSyn1NEG(i int32) (*TVector)
}

type NNImpl struct {
    useHS       bool
    useNEG      bool
    docSize     int
    wordSize    int
    dim         int
    syn0        []TVector  //V(i)
    dsyn0       []TVector  //D(i)
    syn1        []TVector  //for HS
    syn1neg     []TVector  //for NEG
    nextRandom  uint64
}

func (p *NNImpl) GetSyn0(i int32) (*TVector) {
    return &p.syn0[int(i)]
}
func (p *NNImpl) GetDSyn0(i int32) (*TVector) {
    return &p.dsyn0[int(i)]
}
func (p *NNImpl) GetSyn1(i int32) (*TVector) {
    return &p.syn1[int(i)]
}
func (p *NNImpl) GetSyn1NEG(i int32) (*TVector) {
    return &p.syn1neg[int(i)]
}

func (p *NNImpl) getRandomVector(dim int) (vector TVector) {
    vector = make(TVector, 0, dim)
    for i := 0; i < dim; i ++ {
        p.nextRandom = p.nextRandom * 25214903917 + 11
        v := (float32(p.nextRandom & 0xffff) / 65536.0 - 0.5) / float32(dim)
        vector = append(vector, v)
    }
    return vector
}

func NewNN(docSize, wordSize, dim int, useHS, useNEG bool) NN {
    self := &NNImpl{
        useHS: useHS,
        useNEG: useNEG,
        docSize: docSize,
        wordSize: wordSize,
        dim: dim,
        nextRandom: 1,
    }
    syn0 := make([]TVector, 0, wordSize)
    for i:= 0; i < wordSize; i ++ {
        vector := self.getRandomVector(dim)
        syn0 = append(syn0, vector)
    }
    self.syn0 = syn0
    dsyn0 := make([]TVector, 0, docSize)
    for i:=0; i < docSize; i ++ {
        vector := self.getRandomVector(dim)
        dsyn0 = append(dsyn0, vector)
    }
    self.dsyn0 = dsyn0
    if useHS {
        syn1 := make([]TVector, 0, wordSize)
        for i := 0; i < wordSize; i ++ {
            vector := make(TVector, dim, dim)
            syn1 = append(syn1, vector)
        }
        self.syn1 = syn1
    }
    if useNEG {
        syn1neg := make([]TVector, 0, wordSize)
        for i := 0; i < wordSize; i ++ {
            vector := make(TVector, dim, dim)
            syn1neg = append(syn1neg, vector)
        }
        self.syn1neg = syn1neg
    }
    return NN(self)
}
