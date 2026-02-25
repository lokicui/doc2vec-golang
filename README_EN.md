# doc2vec-golang

A Go implementation of Word2Vec / Doc2Vec (Paragraph Vector) for training word and document embeddings. Based on Tomas Mikolov's [word2vec](http://arxiv.org/pdf/1301.3781.pdf) and [doc2vec](http://cs.stanford.edu/~quocle/paragraph_vector.pdf) papers, with support for **Semantic Word Embedding (SWE)** synonym constraints.

> [中文文档](README.md)

## Features

- [x] **CBOW** and **Skip-Gram** model architectures
- [x] **Negative Sampling** and **Hierarchical Softmax** optimization
- [x] **Doc2Vec** document embedding training (PV-DM / PV-DBOW)
- [x] **Online inference** for new document vectors
- [x] **Semantic Word Embedding (SWE)**: ordinal knowledge constraints based on [ACL-2015 paper](http://home.ustc.edu.cn/~quanliu/papers/SWE.pdf)
- [x] word2words / word2docs / doc2docs / doc2words / sen2words / sen2docs
- [x] [Document likelihood estimation](http://arxiv.org/abs/1504.07295)
- [x] Leave-one-out keyword extraction
- [x] Document similarity calculation (DocSimCal)
- [x] Efficient model serialization via MessagePack
- [ ] [Word Mover's Distance (WMD)](https://github.com/hiyijian/doc2vec/blob/master/jmlr.org/proceedings/papers/v37/kusnerb15.pdf)

## Quick Start

### Requirements

- Go >= 1.24

### Build

```bash
go build -o train train.go
go build -o knn knn.go
```

Or use the included script:

```bash
./control build
```

### Training

Training data format: one document per line, two TAB-separated columns (docid + space-tokenized text):

```
1	why does zhihu have some users with avatars and some without
2	the avatar feature is still in beta testing ...
```

**Basic training (Skip-Gram + Negative Sampling):**

```bash
./train data/zhihu_data.1w
```

**Training with full parameters:**

```bash
./train -corpus data/zhihu_data.1w \
        -dim 100 \
        -window 5 \
        -iters 50 \
        -neg \
        -output my.model
```

**Training with synonym constraints (SWE):**

```bash
./train -corpus data/zhihu_data.1w \
        -swe data/synonym_constraints.txt \
        -swe-coeff 0.1 \
        -output swe.model
```

### Querying

```bash
./knn 2.model
```

Interactive operation selection:

```
please select operation type:
        0:word2words
        1:doc_likelihood
        2:leave one out key words
        3:sen2words
        4:sen2docs
        5:word2docs
        6:doc2docs
        7:doc2words
0
Enter text:网页
        1       网页
        0.78    不让
        0.77    浏览
        0.76    邮件
        ...
```

## Training Parameters

| Flag | Default | Description |
|------|---------|-------------|
| `-corpus` | (required) | Training corpus file path (also accepts positional argument) |
| `-output` | `2.model` | Output model file path |
| `-cbow` | `false` | Use CBOW model (default: Skip-Gram) |
| `-hs` | `false` | Use Hierarchical Softmax |
| `-neg` | `true` | Use Negative Sampling |
| `-window` | `5` | Context window size |
| `-dim` | `50` | Word/document embedding dimension |
| `-iters` | `50` | Number of training iterations |

### SWE Synonym Constraint Parameters

| Flag | Default | Description |
|------|---------|-------------|
| `-swe` | (empty) | Semantic constraint file path; empty to disable SWE |
| `-swe-coeff` | `0.1` | Semantic loss weight; higher means stronger constraints |
| `-swe-hinge` | `0.0` | Hinge loss margin |
| `-swe-decay` | `0.0` | Weight decay coefficient (L2 regularization) |
| `-swe-addtime` | `0.0` | Start applying constraints after this training progress (%) |

## Semantic Word Embedding (SWE)

An implementation of Semantic Word Embedding based on the ACL-2015 paper, which enhances word vector quality by introducing **ordinal knowledge constraints**.

### Constraint File Format

Each line contains 4 space-separated words, expressing `sim(A, B) > sim(C, D)`:

```
question answer eat fast
user person website eat
like good website technology
```

Lines starting with `#` are comments. See `data/synonym_constraints.txt` for a complete example.

### Results

On the Zhihu corpus, constraint satisfaction rate improved from **37%** to **91%**:

```
SWE initial: hinge_loss=3.7730 satisfy_rate=0.3714
SWE final:   hinge_loss=0.1785 satisfy_rate=0.9143
```

## Project Structure

```
├── train.go              # Training entry point
├── knn.go                # Query entry point (interactive KNN)
├── control               # Build/deploy script
├── doc2vec/
│   ├── doc2vec.go        # Core algorithms (training, inference, queries)
│   ├── swe.go            # SWE synonym constraint implementation
│   └── wiretypes.go      # Data structures and interface definitions
├── corpus/               # Corpus management (vocabulary, Huffman tree, document index)
├── neuralnet/            # Neural network layer (vector operations, weight matrices)
├── common/               # Common utility functions
├── segmenter/            # Chinese word segmentation (jiebago wrapper)
├── conf/                 # Jieba segmentation dictionaries
├── data/
│   ├── zhihu_data.1w     # Sample corpus (1000 Zhihu Q&A entries)
│   └── synonym_constraints.txt  # Sample synonym constraints
├── interface/            # Thrift IDL definitions
├── SWE_Train.c           # Reference: C implementation of SWE (standalone)
└── codewiki.md           # Detailed code documentation
```

> For full architecture details, algorithm explanations, and data structure documentation, see [codewiki.md](codewiki.md).

## Dependencies

| Library | Purpose |
|---------|---------|
| [tinylib/msgp](https://github.com/tinylib/msgp) | MessagePack model serialization |
| [wangbin/jiebago](https://github.com/wangbin/jiebago) | Chinese word segmentation (used during queries) |
| [astaxie/beego/logs](https://github.com/astaxie/beego) | Logging |

## References

- Mikolov et al., [Efficient Estimation of Word Representations in Vector Space](http://arxiv.org/pdf/1301.3781.pdf) (2013)
- Le & Mikolov, [Distributed Representations of Sentences and Documents](http://cs.stanford.edu/~quocle/paragraph_vector.pdf) (2014)
- Dai et al., [Document Embedding with Paragraph Vectors](http://arxiv.org/pdf/1507.07998.pdf) (2015)
- Liu et al., [Learning Semantic Word Embeddings based on Ordinal Knowledge Constraints](http://home.ustc.edu.cn/~quanliu/papers/SWE.pdf) (ACL 2015)
- [Google word2vec (C)](https://code.google.com/archive/p/word2vec/source/default/source)
- [hiyijian/doc2vec (C++)](https://github.com/hiyijian/doc2vec)
- [iunderstand/SWE (C)](https://github.com/iunderstand/SWE)

## License

[Apache License 2.0](LICENSE)
