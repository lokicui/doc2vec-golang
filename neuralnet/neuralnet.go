package neuralnet

import (
	"log"
)

var gNextRandom uint64 = 1

func init() {
}

func (p TVector) Reset() {
	for i := 0; i < len(p); i++ {
        p[i] = 0
    }
}
//向量加法
func (p TVector) Add(a TVector) {
	if len(p) != len(a) {
		log.Fatal("len(p) != len(a)")
	}
	for i := 0; i < len(p); i++ {
		p[i] += a[i]
	}
}

//向量除法
func (p TVector) Divide(a float32) {
	for i := 0; i < len(p); i++ {
		p[i] /= a
	}
}

//向量点乘
func (p TVector) Dot(a TVector) (f float64) {
	if len(p) != len(a) {
		log.Fatal("len(p) != len(a)")
	}
	for i := 0; i < len(p); i++ {
		f += float64(p[i]) * float64(a[i])
	}
	return f
}

func (p TVector) Multiply(a float64) {
	for i := 0; i < len(p); i++ {
		p[i] *= float32(a)
	}
}

func (p *TNeuralNetImpl) GetSyn0(i int32) *TVector {
	return &p.Syn0[int(i)]
}
func (p *TNeuralNetImpl) GetDSyn0(i int32) *TVector {
    if int(i) > len(p.Dsyn0) {
        log.Fatal("i out of range", i, len(p.Dsyn0))
    }
	return &p.Dsyn0[int(i)]
}
func (p *TNeuralNetImpl) NewDSyn0() *TVector {
	dim := len(*p.GetDSyn0(0))
	self := p.getRandomVector(dim)
	return &self
}
func (p *TNeuralNetImpl) GetSyn1(i int32) *TVector {
	r := &p.Syn1[int(i)]
	return r
}
func (p *TNeuralNetImpl) GetSyn1Neg(i int32) *TVector {
	return &p.Syn1neg[int(i)]
}

func (p *TNeuralNetImpl) getRandomVector(dim int) (vector TVector) {
	vector = make(TVector, 0, dim)
	for i := 0; i < dim; i++ {
		gNextRandom = gNextRandom*25214903917 + 11
		v := (float32(gNextRandom&0xffff)/65536.0 - 0.5) / float32(dim)
		vector = append(vector, v)
	}
	return vector
}

func NewNN(docSize, wordSize, dim int, useHS, useNEG bool) INeuralNet {
	self := &TNeuralNetImpl{}
	syn0 := make([]TVector, 0, wordSize)
	for i := 0; i < wordSize; i++ {
		vector := self.getRandomVector(dim)
		syn0 = append(syn0, vector)
	}
	self.Syn0 = syn0
	dsyn0 := make([]TVector, 0, docSize)
	for i := 0; i < docSize; i++ {
		vector := self.getRandomVector(dim)
		dsyn0 = append(dsyn0, vector)
	}
	self.Dsyn0 = dsyn0
	if useHS {
		syn1 := make([]TVector, 0, wordSize)
		for i := 0; i < wordSize; i++ {
			vector := make(TVector, dim, dim)
			syn1 = append(syn1, vector)
		}
		self.Syn1 = syn1
	}
	if useNEG {
		syn1neg := make([]TVector, 0, wordSize)
		for i := 0; i < wordSize; i++ {
			vector := make(TVector, dim, dim)
			syn1neg = append(syn1neg, vector)
		}
		self.Syn1neg = syn1neg
	}
	return INeuralNet(self)
}
