# doc2vec-golang
golang implement of Tomas Mikolov's word/document embedding. You may want to feel the basic idea from Mikolov's two orignal papers, [word2vec](http://arxiv.org/pdf/1301.3781.pdf) and [doc2vec](http://cs.stanford.edu/~quocle/paragraph_vector.pdf). More recently, Andrew M. Dai etc from Google reported its power in more [detail](http://arxiv.org/pdf/1507.07998.pdf)

# usage
```javascript
[@bjsjs_11_83 doc2vec-golang]$ ./control build
traning Exec build ok
build ok

# The training data(data/zhihu_data.1w) is one document per line, two columns divided by tab, 
# the first column is id, and the second column is the segmented document separated by spaces.
[@bjsjs_11_83 doc2vec-golang]$ ./train  data/zhihu_data.1w          
Skip-Gram Iter:48 Alpha: 0.000796  Progress: 96.81%  Words/sec: 24.27k  
2018-03-30 14:53:00.218536235 +0800 CST training end, 1342521 26861

[@bjsjs_11_83 doc2vec-golang]$ ./knn 2.model 

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
        0.7823723719117796      不让
        0.7651260773728028      浏览
        0.7642516944020028      邮件
        0.7601415883811553      近
        0.7517607921006224      迷恋
        0.7492900066365179      等同
        0.7485966355448261      传说
        0.7463299535930537      基于
        0.7447865182221745      版

please select operation type:
        0:word2words
        1:doc_likelihood
        2:leave one out key words
        3:sen2words
        4:sen2docs
        5:word2docs
        6:doc2docs
        7:doc2words
```

# Dependencies
* [golang](https://golang.org/)
* [msgp](https://github.com/tinylib/msgp)

# 已实现特性
* doc2vec支持CBOW和Skip-Gram两种模型，Negative Sampling和Hierarchical Softmax优化均已实现
* online infer document
* [likelihood of document](http://arxiv.org/abs/1504.07295)
* doc2words
* doc2docs
* word2words
* word2docs

# 未实现特性
* [wmd](https://github.com/hiyijian/doc2vec/blob/master/jmlr.org/proceedings/papers/v37/kusnerb15.pdf)
* [doc2vec添加同义词语义约束](http://home.ustc.edu.cn/~quanliu/papers/SWE.pdf)
* 句子提取核心词

# 参考资料
* google [word2vec](https://code.google.com/archive/p/word2vec/source/default/source) 实现
* [hiyijian/doc2vec](https://github.com/hiyijian/doc2vec)
* [word2vec语义约束](https://github.com/iunderstand/SWE)
* [doc2vec添加同义词语义约束](http://home.ustc.edu.cn/~quanliu/papers/SWE.pdf)


