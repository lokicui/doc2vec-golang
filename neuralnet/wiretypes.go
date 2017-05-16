package neuralnet

import (
	"github.com/tinylib/msgp/msgp"
)

//go:generate msgp

type TVector []float32

type INeuralNet interface {
	GetSyn0(i int32) *TVector
	GetDSyn0(i int32) *TVector
	NewDSyn0() *TVector
	GetSyn1(i int32) *TVector
	GetSyn1Neg(i int32) *TVector
	msgp.Encodable
	msgp.Decodable
	msgp.Marshaler
	msgp.Unmarshaler
	msgp.Sizer
}

type TNeuralNetImpl struct {
	Syn0    []TVector //V(i)
	Dsyn0   []TVector //D(i)
	Syn1    []TVector //for HS
	Syn1neg []TVector //for NEG
}
