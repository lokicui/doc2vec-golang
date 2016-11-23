# doc2vec-golang
golang implement of Tomas Mikolov's word/document embedding. You may want to feel the basic idea from Mikolov's two orignal papers, [word2vec](http://arxiv.org/pdf/1301.3781.pdf) and [doc2vec](http://cs.stanford.edu/~quocle/paragraph_vector.pdf). More recently, Andrew M. Dai etc from Google reported its power in more [detail](http://arxiv.org/pdf/1507.07998.pdf)

# Dependencies
* [golang](https://golang.org/)
* [msgp](https://github.com/tinylib/msgp)

# 已实现特性
* doc2vec支持CBOW和Skip-Gram两种模型，Negative Sampling和Hierarchical Softmax优化均已实现
* online infer document
* [likelihood of document](http://arxiv.org/abs/1504.07295)

# 未实现特性
* [wmd](https://github.com/hiyijian/doc2vec/blob/master/jmlr.org/proceedings/papers/v37/kusnerb15.pdf)
* [doc2vec添加同义词语义约束](http://home.ustc.edu.cn/~quanliu/papers/SWE.pdf)
* doc2words
* word2words
* 句子提取核心词

# 参考资料
* google [word2vec](https://code.google.com/archive/p/word2vec/source/default/source) 实现
* [hiyijian/doc2vec](https://github.com/hiyijian/doc2vec)
* [word2vec语义约束](https://github.com/iunderstand/SWE)
* [doc2vec添加同义词语义约束](http://home.ustc.edu.cn/~quanliu/papers/SWE.pdf)
