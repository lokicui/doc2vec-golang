# doc2vec-golang
golang implement of Tomas Mikolov's word/document embedding. You may want to feel the basic idea from Mikolov's two orignal papers, word2vec and doc2vec. More recently, Andrew M. Dai etc from Google reported its power in more detail

# Dependencies
* golang 
* msgp 

# Why did I rewrite it in golang?
There are a few pretty nice projects like google's word2vec and gensim has already implemented the algorithm, from which I learned quite a lot. However, I rewrite it for following reasons:

* speed. I believe golang version has the best speed on CPU and i have a change to learn golang.

* functionality. I found few project implements both word and document embedding. Moreover, some important application for these embedding have not been fully developed, such as online infer document, likelihood of document, wmd and keyword extraction 

* scalability. I found that it's extremely slow when doing task like "most similar" on large data. One straight-forward way is distributing, the other is putting on GPUs. For these purposes, I prefer to design data structrue by myself
