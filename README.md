# doc2vec-golang

Go 语言实现的 Word2Vec / Doc2Vec（Paragraph Vector）词向量与文档向量训练工具。基于 Tomas Mikolov 的 [word2vec](http://arxiv.org/pdf/1301.3781.pdf) 和 [doc2vec](http://cs.stanford.edu/~quocle/paragraph_vector.pdf) 两篇经典论文，支持 **同义词语义约束（SWE）** 增强训练。

## 特性

- [x] **CBOW** 和 **Skip-Gram** 两种模型架构
- [x] **Negative Sampling** 和 **Hierarchical Softmax** 两种优化算法
- [x] **Doc2Vec** 文档向量训练（PV-DM / PV-DBOW）
- [x] **在线推理** 新文档向量（online infer document）
- [x] **同义词语义约束（SWE）**：基于 [ACL-2015 论文](http://home.ustc.edu.cn/~quanliu/papers/SWE.pdf) 的序数知识约束
- [x] word2words / word2docs / doc2docs / doc2words / sen2words / sen2docs
- [x] [文档似然值估计](http://arxiv.org/abs/1504.07295)（likelihood of document）
- [x] 留一法关键词提取（leave-one-out keywords）
- [x] 文档相似度计算（DocSimCal）
- [x] 模型 MessagePack 高效序列化
- [ ] [WMD 距离](https://github.com/hiyijian/doc2vec/blob/master/jmlr.org/proceedings/papers/v37/kusnerb15.pdf)

## 快速开始

### 环境要求

- Go >= 1.24

### 构建

```bash
go build -o train train.go
go build -o knn knn.go
```

或使用自带脚本：

```bash
./control build
```

### 训练

训练数据格式：每行一条文档，TAB 分隔两列（docid + 分词后的文本）：

```
1	为什么 知 乎 有的 人 有 头像 有的 人 无法 添加 个人 头像
2	头像 功能 还 处于 内测 中 的 内测 ...
```

**基础训练（Skip-Gram + Negative Sampling）：**

```bash
./train data/zhihu_data.1w
```

**完整参数训练：**

```bash
./train -corpus data/zhihu_data.1w \
        -dim 100 \
        -window 5 \
        -iters 50 \
        -neg \
        -output my.model
```

**带同义词约束训练（SWE）：**

```bash
./train -corpus data/zhihu_data.1w \
        -swe data/synonym_constraints.txt \
        -swe-coeff 0.1 \
        -output swe.model
```

### 查询

```bash
./knn 2.model
```

交互式选择操作类型：

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

## 训练参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-corpus` | (必填) | 训练语料文件路径（也支持作为位置参数传入） |
| `-output` | `2.model` | 输出模型文件路径 |
| `-cbow` | `false` | 使用 CBOW 模型（默认 Skip-Gram） |
| `-hs` | `false` | 使用 Hierarchical Softmax |
| `-neg` | `true` | 使用 Negative Sampling |
| `-window` | `5` | 上下文窗口大小 |
| `-dim` | `50` | 词/文档向量维度 |
| `-iters` | `50` | 训练迭代轮数 |

### SWE 同义词约束参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-swe` | (空) | 语义约束文件路径，为空则不启用 SWE |
| `-swe-coeff` | `0.1` | 语义损失权重，越大约束越强 |
| `-swe-hinge` | `0.0` | Hinge Loss 边界值 |
| `-swe-decay` | `0.0` | 权重衰减系数（L2 正则化） |
| `-swe-addtime` | `0.0` | 训练进度达到此百分比后才开始施加约束 |

## 同义词约束（SWE）

基于 ACL-2015 论文实现的 Semantic Word Embedding，通过引入**序数知识约束**增强词向量的语义质量。

### 约束文件格式

每行 4 个空格分隔的词，表示 `sim(A, B) > sim(C, D)`：

```
问题 回答 吃 快
用户 人 网站 吃
喜欢 好 网站 技术
```

以 `#` 开头的行为注释。参见 `data/synonym_constraints.txt` 中的完整示例。

### 效果

在知乎语料上训练，约束满足率从初始 **37%** 提升至 **91%**：

```
SWE initial: hinge_loss=3.7730 satisfy_rate=0.3714
SWE final:   hinge_loss=0.1785 satisfy_rate=0.9143
```

## 项目结构

```
├── train.go              # 训练入口
├── knn.go                # 查询入口（交互式 KNN）
├── control               # 构建/部署脚本
├── doc2vec/
│   ├── doc2vec.go        # 核心算法（训练、推理、查询）
│   ├── swe.go            # SWE 同义词约束实现
│   └── wiretypes.go      # 数据结构与接口定义
├── corpus/               # 语料库管理（词表、Huffman 树、文档索引）
├── neuralnet/            # 神经网络层（向量运算、权重矩阵）
├── common/               # 通用工具函数
├── segmenter/            # 中文分词（jiebago 封装）
├── conf/                 # jieba 分词词典
├── data/
│   ├── zhihu_data.1w     # 示例语料（1000 条知乎问答）
│   └── synonym_constraints.txt  # 示例同义词约束
├── interface/            # Thrift IDL 定义
├── SWE_Train.c           # 参考：C 语言 SWE 实现（独立程序）
└── codewiki.md           # 详细代码文档
```

> 完整的架构说明、算法原理、数据结构文档请参见 [codewiki.md](codewiki.md)。

## 依赖

| 库 | 用途 |
|----|------|
| [tinylib/msgp](https://github.com/tinylib/msgp) | MessagePack 模型序列化 |
| [wangbin/jiebago](https://github.com/wangbin/jiebago) | 中文分词（查询时使用） |
| [astaxie/beego/logs](https://github.com/astaxie/beego) | 日志 |

## 参考资料

- Mikolov et al., [Efficient Estimation of Word Representations in Vector Space](http://arxiv.org/pdf/1301.3781.pdf) (2013)
- Le & Mikolov, [Distributed Representations of Sentences and Documents](http://cs.stanford.edu/~quocle/paragraph_vector.pdf) (2014)
- Dai et al., [Document Embedding with Paragraph Vectors](http://arxiv.org/pdf/1507.07998.pdf) (2015)
- Liu et al., [Learning Semantic Word Embeddings based on Ordinal Knowledge Constraints](http://home.ustc.edu.cn/~quanliu/papers/SWE.pdf) (ACL 2015)
- [Google word2vec (C)](https://code.google.com/archive/p/word2vec/source/default/source)
- [hiyijian/doc2vec (C++)](https://github.com/hiyijian/doc2vec)
- [iunderstand/SWE (C)](https://github.com/iunderstand/SWE)

## License

[Apache License 2.0](LICENSE)
