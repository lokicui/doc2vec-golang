package doc2vec

import (
	"bufio"
	"fmt"
	"github.com/lokicui/doc2vec-golang/corpus"
	"github.com/lokicui/doc2vec-golang/neuralnet"
	"log"
	"math"
	"os"
	"strings"
)

// SWEConfig holds parameters for Semantic Word Embedding constraints.
// Based on ACL-2015 paper: "Learning Semantic Word Embeddings based on Ordinal Knowledge Constraints"
type SWEConfig struct {
	Coeff       float64 // interpolation coefficient (default 0.1), weight of semantic loss
	HingeMargin float64 // hinge loss margin (default 0.0)
	WeightDecay float64 // L2 regularization on constrained word vectors (default 0.0)
	AddTime     float64 // training progress % at which to start applying constraints (default 0.0)
}

// InEquation represents a single ordinal constraint: sim(A,B) > sim(C,D)
type InEquation struct {
	IndexA int32
	IndexB int32
	IndexC int32
	IndexD int32
}

// SWEConstraints holds all loaded constraints and a per-word index for fast lookup.
type SWEConstraints struct {
	Inequations      []InEquation
	WordToConstraints map[int32][]int // word index → indices into Inequations
}

func NewSWEConfig(coeff, hingeMargin, weightDecay, addTime float64) *SWEConfig {
	return &SWEConfig{
		Coeff:       coeff,
		HingeMargin: hingeMargin,
		WeightDecay: weightDecay,
		AddTime:     addTime,
	}
}

// LoadSWEConstraints reads an inequality file and builds the constraint database.
// File format: each line contains 4 space-separated words: word_A word_B word_C word_D
// Meaning: sim(word_A, word_B) > sim(word_C, word_D)
func LoadSWEConstraints(fname string, corp corpus.ICorpus) (*SWEConstraints, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("failed to open SWE constraint file: %w", err)
	}
	defer file.Close()

	sc := &SWEConstraints{
		WordToConstraints: make(map[int32][]int),
	}

	scanner := bufio.NewScanner(file)
	lineNo := 0
	skipped := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 4 {
			log.Printf("SWE: skipping line %d: expected 4 words, got %d", lineNo, len(fields))
			lineNo++
			continue
		}

		idxA, okA := corp.GetWordIdx(fields[0])
		idxB, okB := corp.GetWordIdx(fields[1])
		idxC, okC := corp.GetWordIdx(fields[2])
		idxD, okD := corp.GetWordIdx(fields[3])

		if !okA || !okB || !okC || !okD {
			skipped++
			lineNo++
			continue
		}

		ineq := InEquation{
			IndexA: idxA,
			IndexB: idxB,
			IndexC: idxC,
			IndexD: idxD,
		}
		eqIdx := len(sc.Inequations)
		sc.Inequations = append(sc.Inequations, ineq)

		added := map[int32]bool{}
		addWord := func(idx int32) {
			if !added[idx] {
				sc.WordToConstraints[idx] = append(sc.WordToConstraints[idx], eqIdx)
				added[idx] = true
			}
		}
		addWord(idxA)
		addWord(idxB)
		addWord(idxC)
		addWord(idxD)

		lineNo++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading SWE constraint file: %w", err)
	}

	log.Printf("SWE: loaded %d constraints from %s (%d skipped due to OOV words, %d words involved)",
		len(sc.Inequations), fname, skipped, len(sc.WordToConstraints))
	return sc, nil
}

// vectorNorm computes the L2 norm of a vector
func vectorNorm(v neuralnet.TVector) float64 {
	var sum float64
	for _, x := range v {
		sum += float64(x) * float64(x)
	}
	return math.Sqrt(sum)
}

// cosineFromVectors computes cosine similarity between two vectors
func cosineFromVectors(a, b neuralnet.TVector) float64 {
	normA := vectorNorm(a)
	normB := vectorNorm(b)
	if normA == 0 || normB == 0 {
		return 0
	}
	var dot float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
	}
	return dot / (normA * normB)
}

// ComputeSWEGradient computes the semantic constraint gradient for a given word.
// The gradient is the partial derivative of the total hinge loss with respect to
// the word's embedding vector.
//
// For each constraint sim(A,B) > sim(C,D) involving this word:
//   loss = max(0, margin - (cos(A,B) - cos(C,D)))
//   gradient += hingeDerivative * (∂cos(C,D)/∂w - ∂cos(A,B)/∂w)
func (sc *SWEConstraints) ComputeSWEGradient(wordIdx int32, nn neuralnet.INeuralNet, dim int, hingeMargin float64) neuralnet.TVector {
	grad := make(neuralnet.TVector, dim)

	constraints, ok := sc.WordToConstraints[wordIdx]
	if !ok || len(constraints) == 0 {
		return grad
	}

	derivAB := make(neuralnet.TVector, dim)
	derivCD := make(neuralnet.TVector, dim)

	for _, eqIdx := range constraints {
		ineq := sc.Inequations[eqIdx]

		vecA := nn.GetSyn0(ineq.IndexA)
		vecB := nn.GetSyn0(ineq.IndexB)
		vecC := nn.GetSyn0(ineq.IndexC)
		vecD := nn.GetSyn0(ineq.IndexD)

		normA := vectorNorm(*vecA)
		normB := vectorNorm(*vecB)
		normC := vectorNorm(*vecC)
		normD := vectorNorm(*vecD)

		if normA == 0 || normB == 0 || normC == 0 || normD == 0 {
			continue
		}

		distAB := cosineFromVectors(*vecA, *vecB)
		distCD := cosineFromVectors(*vecC, *vecD)

		// hinge loss input: margin - (cos(A,B) - cos(C,D))
		hingeInput := hingeMargin - (distAB - distCD)

		// hinge derivative: 1 if hingeInput > 0, else 0
		if hingeInput <= 0 {
			continue
		}

		// Reset temporary gradient vectors
		for i := range derivAB {
			derivAB[i] = 0
			derivCD[i] = 0
		}

		// Compute ∂cos(A,B)/∂word
		// ∂cos(A,B)/∂A = B/(|A||B|) - cos(A,B) * A/|A|²
		// ∂cos(A,B)/∂B = A/(|A||B|) - cos(A,B) * B/|B|²
		if wordIdx == ineq.IndexA {
			for i := 0; i < dim; i++ {
				derivAB[i] = float32(float64((*vecB)[i])/(normA*normB) - distAB*float64((*vecA)[i])/(normA*normA))
			}
		} else if wordIdx == ineq.IndexB {
			for i := 0; i < dim; i++ {
				derivAB[i] = float32(float64((*vecA)[i])/(normA*normB) - distAB*float64((*vecB)[i])/(normB*normB))
			}
		}

		// Compute ∂cos(C,D)/∂word
		if wordIdx == ineq.IndexC {
			for i := 0; i < dim; i++ {
				derivCD[i] = float32(float64((*vecD)[i])/(normC*normD) - distCD*float64((*vecC)[i])/(normC*normC))
			}
		} else if wordIdx == ineq.IndexD {
			for i := 0; i < dim; i++ {
				derivCD[i] = float32(float64((*vecC)[i])/(normC*normD) - distCD*float64((*vecD)[i])/(normD*normD))
			}
		}

		// Accumulate: f' * (∂cos(C,D)/∂w - ∂cos(A,B)/∂w)
		// hingeDerivative is 1.0 (since we already checked hingeInput > 0)
		for i := 0; i < dim; i++ {
			grad[i] += derivCD[i] - derivAB[i]
		}
	}

	return grad
}

// EvalSWEQuality evaluates the current quality of constraints satisfaction.
// Returns (total hinge loss, satisfaction rate).
func (sc *SWEConstraints) EvalSWEQuality(nn neuralnet.INeuralNet, hingeMargin float64) (totalLoss float64, satisfyRate float64) {
	if len(sc.Inequations) == 0 {
		return 0, 1.0
	}

	satisfyCount := 0
	for _, ineq := range sc.Inequations {
		distAB := cosineFromVectors(*nn.GetSyn0(ineq.IndexA), *nn.GetSyn0(ineq.IndexB))
		distCD := cosineFromVectors(*nn.GetSyn0(ineq.IndexC), *nn.GetSyn0(ineq.IndexD))

		hingeInput := hingeMargin - (distAB - distCD)
		if hingeInput > 0 {
			totalLoss += hingeInput
		}
		if distAB > distCD {
			satisfyCount++
		}
	}

	satisfyRate = float64(satisfyCount) / float64(len(sc.Inequations))
	return totalLoss, satisfyRate
}
