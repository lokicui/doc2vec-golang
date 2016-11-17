#objects := $(patsubst %.cpp,%.o,$(wildcard gen-cpp/*.cpp))
thrift_fnames := $(wildcard interface/*.thrift)
#CXXFLAGS=-I/usr/local/include/thrift
#dirs := include lib64_release lib64_debug


all : thrift
	@export GOPATH=`pwd` && cd src && go build  train.go && cd - && echo "all"
	@export GOPATH=`pwd` && cd src && go build  knn.go && cd - && echo "all"

debug : thrift
	@export GOPATH=`pwd` && cd src && go build -gcflags "-N -l" train.go && cd - && echo "debug"
	@export GOPATH=`pwd` && cd src && go build -gcflags "-N -l" knn.go && cd - && echo "debug"


.PHONY : thrift debug

thrift :
	#@for fname in $(thrift_fnames); do thrift -r --gen go --gen cpp --gen py --gen php $$fname;done
	@for fname in $(thrift_fnames); do thrift -r --gen go -o interface $$fname;done
	@cd src && ln -sf ../interface/gen-go/*  . && cd -
