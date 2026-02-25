## Cursor Cloud specific instructions

### Overview

This is **doc2vec-golang**, a Go implementation of word2vec/doc2vec. It produces two CLI binaries:

- `train` — trains a model on a corpus file and saves it as `2.model`
- `knn` — loads a trained model and provides an interactive CLI for similarity queries

### Prerequisites

Go >= 1.24 is required (the `msgp` dependency needs it). The update script installs Go 1.24 if the system Go is older.

### Build

```bash
go build -o train train.go
go build -o knn knn.go
```

Or use the included `./control build` script.

### Running

```bash
# Train (takes ~8s on the bundled 1000-doc corpus)
./train data/zhihu_data.1w

# Query (interactive; pipe stdin for non-interactive use)
printf '0\n网页\n' | timeout 10 ./knn 2.model
```

### Testing

```bash
# Sub-package tests pass; root package has a known C-file issue with go test
go test ./corpus/... ./neuralnet/...

# Lint
go vet ./common/... ./corpus/... ./doc2vec/... ./neuralnet/... ./segmenter/...
```

### Known issues

- `go test ./...` at the root fails because `SWE_Train.c` (standalone C program) is in the root package directory. Run sub-package tests individually instead.
- `doc2vec/wiretypes_gen_test.go` has a pre-existing panic in `TestMarshalUnmarshalTDoc2VecImpl` (nil pointer from uninitialized struct).
