# doc2vec-golang CodeWiki

## 目录

- [1. 项目概述](#1-项目概述)
- [2. 整体架构](#2-整体架构)
- [3. 目录结构](#3-目录结构)
- [4. 核心模块详解](#4-核心模块详解)
  - [4.1 common — 通用工具函数](#41-common--通用工具函数)
  - [4.2 corpus — 语料库管理](#42-corpus--语料库管理)
  - [4.3 neuralnet — 神经网络层](#43-neuralnet--神经网络层)
  - [4.4 doc2vec — 核心算法引擎](#44-doc2vec--核心算法引擎)
  - [4.5 segmenter — 中文分词器](#45-segmenter--中文分词器)
- [5. 数据结构总览](#5-数据结构总览)
- [6. 算法原理](#6-算法原理)
  - [6.1 Word2Vec 基础](#61-word2vec-基础)
  - [6.2 Doc2Vec (Paragraph Vector)](#62-doc2vec-paragraph-vector)
  - [6.3 Hierarchical Softmax (层次 Softmax)](#63-hierarchical-softmax-层次-softmax)
  - [6.4 Negative Sampling (负采样)](#64-negative-sampling-负采样)
  - [6.5 Huffman 树构建](#65-huffman-树构建)
- [7. 训练流程](#7-训练流程)
- [8. 推理 / 查询流程](#8-推理--查询流程)
- [9. 模型序列化](#9-模型序列化)
- [10. 数据格式](#10-数据格式)
- [11. 可执行程序](#11-可执行程序)
  - [11.1 train — 训练程序](#111-train--训练程序)
  - [11.2 knn — 查询程序](#112-knn--查询程序)
- [12. 接口定义 (Thrift)](#12-接口定义-thrift)
- [13. SWE_Train.c — 语义词嵌入工具](#13-swe_trainc--语义词嵌入工具)
- [14. 依赖库](#14-依赖库)
- [15. 关键常量与配置](#15-关键常量与配置)
- [16. 并发模型](#16-并发模型)
- [17. 性能优化点](#17-性能优化点)
- [18. 扩展阅读](#18-扩展阅读)

---

## 1. 项目概述

**doc2vec-golang** 是 Tomas Mikolov 的 Word2Vec / Doc2Vec (Paragraph Vector) 算法的 Go 语言实现。项目能够从文本语料中学习词向量 (word embedding) 和文档向量 (document embedding)，并提供丰富的相似度查询能力。

### 核心能力

| 功能 | 说明 |
|------|------|
| **word2words** | 给定一个词，找出语义最相似的 Top-K 词 |
| **word2docs** | 给定一个词，找出最相关的 Top-K 文档 |
| **doc2docs** | 给定一篇文档（by index），找出最相似的 Top-K 文档 |
| **doc2words** | 给定一篇文档（by index），找出与其最相关的 Top-K 词 |
| **sen2words** | 给定一段文本（在线推理其向量），找出最相似的 Top-K 词 |
| **sen2docs** | 给定一段文本（在线推理其向量），找出最相似的 Top-K 文档 |
| **doc_likelihood** | 计算给定文本在模型下的似然值 |
| **leave one out** | 通过"留一法"提取文档中的核心关键词 |
| **DocSimCal** | 计算两篇文档之间的相似度（基于 soft cosine / BOW 矩阵） |

---

## 2. 整体架构

```
┌─────────────────────────────────────────────────────────┐
│                      应用层 (main)                       │
│  ┌──────────┐    ┌──────────┐                           │
│  │ train.go │    │  knn.go  │                           │
│  └────┬─────┘    └────┬─────┘                           │
│       │               │                                 │
│       ▼               ▼                                 │
│  ┌─────────────────────────────────────────────────┐    │
│  │              doc2vec (核心算法引擎)               │    │
│  │  IDoc2Vec 接口 / TDoc2VecImpl 实现               │    │
│  │  ┌──────────────┐  ┌──────────────────────────┐ │    │
│  │  │ 训练算法      │  │ 查询 / 推理算法           │ │    │
│  │  │ • CBOW       │  │ • Word2Words             │ │    │
│  │  │ • Skip-Gram  │  │ • Doc2Docs               │ │    │
│  │  │ • HS / NEG   │  │ • Sen2Words / Sen2Docs   │ │    │
│  │  └──────────────┘  │ • FitDoc (在线推理)       │ │    │
│  │                    │ • Likelihood / LOO        │ │    │
│  │                    └──────────────────────────┘ │    │
│  └──────────┬──────────────────┬───────────────────┘    │
│             │                  │                        │
│     ┌───────▼───────┐  ┌──────▼────────┐               │
│     │    corpus     │  │   neuralnet   │               │
│     │ (语料库管理)   │  │ (神经网络权重) │               │
│     │ • 词表构建     │  │ • Syn0 词向量  │               │
│     │ • Huffman 树   │  │ • Dsyn0 文档   │               │
│     │ • 文档索引     │  │   向量        │               │
│     │ • Transform   │  │ • Syn1 (HS)   │               │
│     └───────┬───────┘  │ • Syn1neg     │               │
│             │          │   (NEG)       │               │
│     ┌───────▼───────┐  └───────────────┘               │
│     │    common     │                                   │
│     │ (工具函数)     │  ┌───────────────┐               │
│     │ • SBC2DBC     │  │   segmenter   │               │
│     │ • Min/Max     │  │ (中文分词)     │               │
│     │ • If          │  │ • jiebago     │               │
│     └───────────────┘  └───────────────┘               │
└─────────────────────────────────────────────────────────┘
```

**数据流向：**

```
训练语料文件 ──► corpus.Build() ──► 词表 + Huffman树 + 文档索引
                                         │
                                         ▼
                              neuralnet.NewNN() ──► 初始化随机权重
                                         │
                                         ▼
                              trainCbow() / trainSkipGram()
                                         │
                              (多轮迭代，多协程并发)
                                         │
                                         ▼
                              SaveModel() ──► MessagePack 序列化 ──► .model 文件
```

---

## 3. 目录结构

```
doc2vec-golang/
├── train.go              # 训练入口：读取语料 → 训练 → 保存模型
├── knn.go                # 查询入口：加载模型 → 交互式 KNN 查询
├── control               # 构建/部署/运维脚本 (bash)
├── common/
│   └── common.go         # 通用工具函数（全角半角转换、Min/Max/If）
├── corpus/
│   ├── wiretypes.go      # 语料库数据结构定义 + ICorpus 接口
│   ├── corpus.go         # 语料库核心实现（词表构建、Huffman树、文档索引）
│   ├── wiretypes_gen.go  # msgp 自动生成的序列化代码
│   └── wiretypes_gen_test.go  # msgp 生成的序列化测试
├── neuralnet/
│   ├── wiretypes.go      # 神经网络数据结构定义 + INeuralNet 接口
│   ├── neuralnet.go      # 神经网络实现（向量运算、权重初始化）
│   ├── wiretypes_gen.go  # msgp 自动生成的序列化代码
│   └── wiretypes_gen_test.go  # msgp 生成的序列化测试
├── doc2vec/
│   ├── wiretypes.go      # Doc2Vec 数据结构定义 + IDoc2Vec 接口
│   ├── doc2vec.go        # 核心算法（训练、推理、查询、排序）
│   ├── wiretypes_gen.go  # msgp 自动生成的序列化代码
│   └── wiretypes_gen_test.go  # msgp 生成的序列化测试
├── segmenter/
│   └── segmenter.go      # 中文分词封装（基于 jiebago）
├── conf/
│   ├── dict.txt          # jieba 主词典
│   └── userdict.txt      # jieba 用户自定义词典
├── data/
│   └── zhihu_data.1w     # 示例训练语料（1000 条知乎问答）
├── interface/
│   ├── structdef.thrift  # Thrift 结构定义（数据模型的 IDL 描述）
│   └── keylist.thrift    # Thrift 服务接口定义（KV/Hash/ZSet/MsgQueue 服务）
├── SWE_Train.c           # 独立的 C 程序：带语义约束的词嵌入训练
├── go.mod                # Go modules 配置
├── go.sum                # 依赖校验和
└── README.md             # 项目说明
```

---

## 4. 核心模块详解

### 4.1 common — 通用工具函数

**文件：** `common/common.go`  
**包名：** `common`

提供跨模块使用的基础工具函数：

| 函数 | 签名 | 说明 |
|------|------|------|
| `SBC2DBC` | `(s string) string` | 全角字符转半角（如 `Ａ` → `A`、`　` → 空格） |
| `DBC2SBC` | `(s string) string` | 半角字符转全角（逆操作） |
| `If` | `(condition bool, trueVal, falseVal interface{}) interface{}` | 三元表达式模拟 |
| `Max` | `(first int, args ...int) int` | 可变参数取最大值 |
| `Min` | `(first int, args ...int) int` | 可变参数取最小值 |

**全角半角转换的作用：** 训练语料中可能混合全角和半角字符（常见于中文文本），统一转为半角可以避免同一个词被当作两个不同的 token。

---

### 4.2 corpus — 语料库管理

**文件：** `corpus/wiretypes.go`（结构定义） + `corpus/corpus.go`（实现）  
**包名：** `corpus`

#### 核心接口 `ICorpus`

```go
type ICorpus interface {
    Build(fname string) (err error)          // 从文件构建语料库
    GetVocabCnt() int                        // 词表大小（去重后）
    GetDocCnt() int                          // 文档数量
    GetWordsCnt() int                        // 总词数（未去重）
    GetWordIdx(word string) (int32, bool)    // 词 → 索引
    GetWordItemByIdx(i int) *TWordItem       // 索引 → 词条目
    GetAllWords() TWordItemSlice             // 全部词条目
    GetAllDocWordsIdx() [][]int32            // 全部文档的词索引列表
    GetAllDocWords() [][]*TWordItem          // 全部文档的词条目列表
    GetDocWordsByDocid(id string) []*TWordItem // docid → 文档词列表
    GetDocWordsByIdx(i int) []*TWordItem     // 文档索引 → 文档词列表
    Transform(content string) []int32        // 文本 → 词索引序列
    // + msgp 序列化接口
}
```

#### 核心数据结构

```go
type TWordItem struct {
    Cnt   int32    // 词频 (term frequency)
    Point []int32  // Huffman 树路径：从根到叶节点 [root, leaf)，存储内部节点的索引
    Code  []bool   // Huffman 编码：从根到叶节点 (root, leaf]，0/1 编码序列
    Word  string   // 词的文本内容
}

type TCorpusImpl struct {
    Words        TWordItemSlice    // 词表，按词频降序排列
    Word2Idx     map[string]int32  // 词 → 词表中的下标
    Doc2WordsIdx [][]int32         // 文档列表：每个文档存储其词的索引序列
    Doc2Idx      map[string]int32  // docid → Doc2WordsIdx 中的下标
    MinReduce    int32             // 词表裁剪阈值（动态增长）
    MinCnt       int32             // 最低词频阈值
    WordsCnt     int               // 总词数（未去重）
}
```

#### 语料库构建流程

```
Build(fname)
    │
    ├──► buildVocabulary(fname)
    │       │
    │       ├── 逐行读取文件，按 TAB 分割为 (docid, content)
    │       ├── 对每个词调用 addWord() 累加词频
    │       ├── 当词表超过 0.7 × VOCAB_HASH_SIZE (21M) 时，触发 reduceVocabulary()
    │       ├── sortVocab()：按词频降序排列，过滤低频词
    │       └── createBinaryTree()：构建 Huffman 树，为每个词生成 Code 和 Point
    │
    └──► loadDocument(fname)
            │
            ├── 第二次扫描文件
            ├── 对每个文档的词查找索引，构建 Doc2WordsIdx
            └── 维护 Doc2Idx 映射
```

#### 词表裁剪 `reduceVocabulary()`

当词表大小超过 `0.7 × VOCAB_HASH_SIZE = 21,000,000` 时，删除词频 ≤ `MinReduce` 的词，并将 `MinReduce` 递增。这保证了内存不会因巨大语料而溢出。

---

### 4.3 neuralnet — 神经网络层

**文件：** `neuralnet/wiretypes.go`（结构定义） + `neuralnet/neuralnet.go`（实现）  
**包名：** `neuralnet`

#### 核心接口 `INeuralNet`

```go
type INeuralNet interface {
    GetSyn0(i int32) *TVector       // 获取第 i 个词的输入向量 V(w)
    GetDSyn0(i int32) *TVector      // 获取第 i 个文档的向量 D(d)
    NewDSyn0() *TVector             // 创建一个新的随机文档向量（用于在线推理）
    GetSyn1(i int32) *TVector       // 获取 Hierarchical Softmax 的内部节点参数 θ
    GetSyn1Neg(i int32) *TVector    // 获取 Negative Sampling 的输出向量 θ'
    // + msgp 序列化接口
}
```

#### 权重矩阵

```go
type TVector []float32

type TNeuralNetImpl struct {
    Syn0    []TVector  // 词输入向量矩阵，大小 [vocab_size × dim]
    Dsyn0   []TVector  // 文档向量矩阵，大小 [doc_size × dim]
    Syn1    []TVector  // HS 参数矩阵，大小 [vocab_size × dim]（仅 UseHS 时分配）
    Syn1neg []TVector  // NEG 参数矩阵，大小 [vocab_size × dim]（仅 UseNEG 时分配）
}
```

#### 向量运算

`TVector` 类型提供了以下就地 (in-place) 运算：

| 方法 | 说明 | 数学表示 |
|------|------|---------|
| `Reset()` | 清零向量 | v = **0** |
| `Add(a)` | 向量加法 | v = v + a |
| `Divide(a)` | 标量除法 | v = v / a |
| `Dot(a)` | 向量点积 | f = v · a |
| `Multiply(a)` | 标量乘法 | v = v × a |

#### 权重初始化

使用线性同余随机数生成器（LCG）初始化：

```
nextRandom = nextRandom × 25214903917 + 11
v[i] = (float32(nextRandom & 0xFFFF) / 65536.0 - 0.5) / dim
```

初始值范围为 `[-0.5/dim, 0.5/dim]`，与 Google word2vec 原始实现一致。

---

### 4.4 doc2vec — 核心算法引擎

**文件：** `doc2vec/wiretypes.go`（结构定义） + `doc2vec/doc2vec.go`（实现）  
**包名：** `doc2vec`

#### 核心接口 `IDoc2Vec`

```go
type IDoc2Vec interface {
    Train(fname string)                                    // 从文件训练模型
    SaveModel(fname string) (err error)                    // 保存模型到文件
    LoadModel(fname string) (err error)                    // 从文件加载模型
    GetCorpus() corpus.ICorpus                             // 获取语料库
    GetNeuralNet() neuralnet.INeuralNet                    // 获取神经网络
    Word2Words(word string)                                // 词 → 相似词
    Word2Docs(word string)                                 // 词 → 相关文档
    Sen2Words(content string, iters int)                   // 文本 → 相似词
    Sen2Docs(content string, iters int)                    // 文本 → 相关文档
    Doc2Docs(docidx int)                                   // 文档 → 相似文档
    Doc2Words(docidx int)                                  // 文档 → 相关词
    GetLikelihood4Doc(context string) (likelihood float64) // 文档似然值
    GetLeaveOneOutKwds(content string, iters int)          // 留一法关键词
    DocSimCal(content1, content2 string) (dis float64)     // 文档相似度
}
```

#### 实现结构 `TDoc2VecImpl`

```go
type TDoc2VecImpl struct {
    Trainfile    string              // 训练文件路径
    Dim          int                 // 向量维度
    UseCbow      bool                // true=CBOW, false=Skip-Gram
    WindowSize   int                 // 上下文窗口大小
    UseHS        bool                // 是否使用 Hierarchical Softmax
    UseNEG       bool                // 是否使用 Negative Sampling
    Negative     int                 // 负样本数量（默认 5）
    StartAlpha   float64             // 初始学习率（CBOW=0.05, Skip-Gram=0.025）
    Iters        int                 // 训练迭代轮数
    TrainedWords int                 // 已训练词数（用于计算进度和学习率衰减）
    Corpus       corpus.ICorpus      // 语料库实例
    NN           neuralnet.INeuralNet // 神经网络实例
    Pool         *sync.Pool          // TVector 对象池（减少 GC 压力）
}
```

---

### 4.5 segmenter — 中文分词器

**文件：** `segmenter/segmenter.go`  
**包名：** `segmenter`

封装了 `jiebago` 分词库，在 `init()` 时自动加载词典：

- 主词典：`conf/dict.txt`
- 用户词典：`conf/userdict.txt`

提供 `GetSegmenter()` 函数返回全局分词器实例。在 `knn.go` 中用于对用户输入的查询文本进行在线分词。

> **注意：** `segmenter` 使用的是 `jiebago/posseg`（词性标注分词器），但代码中只使用了 `Cut()` 方法进行分词，未使用词性信息。

---

## 5. 数据结构总览

```
TDoc2VecImpl
├── Corpus (ICorpus → TCorpusImpl)
│   ├── Words []TWordItem          ← 词表（按词频降序）
│   │   ├── Word   string          ← 词文本
│   │   ├── Cnt    int32           ← 词频
│   │   ├── Point  []int32         ← Huffman路径（内部节点索引）
│   │   └── Code   []bool          ← Huffman编码
│   ├── Word2Idx map[string]int32  ← 词→索引映射
│   ├── Doc2WordsIdx [][]int32     ← 文档→词索引列表
│   └── Doc2Idx map[string]int32   ← docid→文档索引映射
│
├── NN (INeuralNet → TNeuralNetImpl)
│   ├── Syn0    []TVector          ← 词输入向量 V(w)  [vocab_size × dim]
│   ├── Dsyn0   []TVector          ← 文档向量 D(d)    [doc_size × dim]
│   ├── Syn1    []TVector          ← HS 参数 θ        [vocab_size × dim]
│   └── Syn1neg []TVector          ← NEG 参数 θ'      [vocab_size × dim]
│
└── Pool *sync.Pool                ← TVector 对象池
```

---

## 6. 算法原理

### 6.1 Word2Vec 基础

Word2Vec 的目标是将词映射为低维稠密向量，使得语义相近的词在向量空间中距离也近。有两种模型架构：

#### CBOW (Continuous Bag of Words)

从上下文预测中心词：

```
P(w | Context(w)) = P(w | w_{-c}, ..., w_{-1}, w_{+1}, ..., w_{+c})
```

**隐层表示：** `X(w) = Average(V(w_{-c}), ..., V(w_{+c}), D(doc))`

即将上下文词向量与文档向量取平均作为隐层输入。

#### Skip-Gram

从中心词预测上下文中的每个词：

```
P(Context(w) | w) = ∏ P(w_i | w)   for w_i in Context(w)
```

Skip-Gram 对每一对 `(中心词, 上下文词)` 独立训练。

### 6.2 Doc2Vec (Paragraph Vector)

Doc2Vec 在 Word2Vec 的基础上，为每篇文档额外学习一个向量 `D(d)`：

- **PV-DM (Distributed Memory)：** 对应 CBOW 架构。文档向量与上下文词向量一起参与隐层平均，预测中心词。
- **PV-DBOW (Distributed Bag of Words)：** 对应 Skip-Gram 架构。文档向量作为额外的"上下文"，与每个词配对训练。

本项目在两种模型中都实现了文档向量的训练：

```go
// Skip-Gram 中：文档向量作为额外的 range vector 与每个中心词配对
p.trainSkipGram4Pair(widx, dsyn0, alpha, false)

// CBOW 中：文档向量与上下文词向量一起参与平均
neu1.Add(*dsyn0)
cw++
neu1.Divide(float32(cw))
```

### 6.3 Hierarchical Softmax (层次 Softmax)

利用 Huffman 树将 softmax 的复杂度从 O(V) 降低到 O(log V)。

每个词对应 Huffman 树中的一个叶节点，从根到叶的路径上每个内部节点都是一个二分类器：

```
P(w | X) = ∏ σ(±θ_j · X)
```

其中 `θ_j` 是内部节点 j 的参数向量（存储在 `Syn1` 中），`±` 由 Huffman 编码决定。

**梯度更新：**

```
g = (1 - label - σ(θ · X)) × α
e = e + g × θ     (误差累积)
θ = θ + g × X     (参数更新)
```

### 6.4 Negative Sampling (负采样)

为每个正样本采样 K 个负样本（默认 `K=5`），将多分类问题转化为二分类：

```
log P(w | context) ≈ log σ(θ_w · X) + ∑_{k=1}^{K} log σ(-θ_{w_k} · X)
```

**负采样分布：** 按词频的 0.75 次方进行带权采样，使用预计算的查找表（大小 1 亿）实现 O(1) 采样：

```go
const NEG_SAMPLING_TABLE_SIZE = 1e8

func (p *TDoc2VecImpl) initUnigramTable() {
    // 按 P(w) ∝ freq(w)^0.75 构建查找表
}

func GetNegativeSamplingWordIdx() int32 {
    gNextRandom = gNextRandom*25214903917 + 11
    idx := int(gNextRandom>>16) % NEG_SAMPLING_TABLE_SIZE
    return gNegSamplingTable[idx]
}
```

### 6.5 Huffman 树构建

`createBinaryTree()` 使用经典的两队列 Huffman 构建算法：

1. 将所有词按词频降序排列
2. 使用两个指针 `pos1`（从词表尾部向前）和 `pos2`（从内部节点区域向后）
3. 每次取两个最小的节点合并，生成 `vocab_size - 1` 个内部节点
4. 从每个叶节点回溯到根，生成 `Code`（编码）和 `Point`（路径）
5. 将路径和编码反转，使其从根到叶的顺序存储

**Point 的特殊处理：** `Point` 中只存储内部节点的索引（减去 `vocab_size` 后直接对应 `Syn1` 的下标），不包含叶节点。

---

## 7. 训练流程

### 完整训练流程

```
Train(fname)
    │
    ├── 1. Corpus.Build(fname)
    │       ├── 构建词表（两次扫描文件）
    │       ├── 排序、过滤低频词
    │       ├── 构建 Huffman 树
    │       └── 加载文档索引
    │
    ├── 2. initUnigramTable()  （仅 UseNEG 时）
    │       └── 构建负采样查找表
    │
    ├── 3. NewNN(docSize, vocabSize, dim, useHS, useNEG)
    │       ├── 初始化 Syn0  [vocabSize × dim]  随机值
    │       ├── 初始化 Dsyn0 [docSize × dim]    随机值
    │       ├── 初始化 Syn1  [vocabSize × dim]  零值（仅 HS）
    │       └── 初始化 Syn1neg [vocabSize × dim] 零值（仅 NEG）
    │
    └── 4. trainCbow() 或 trainSkipGram()
            │
            └── 循环 Iters 轮
                    │
                    └── 对每篇文档：
                            ├── 获取文档向量 dsyn0
                            ├── 对每个词位置：
                            │   ├── 随机窗口大小 b = rand() % WindowSize
                            │   ├── 取上下文 [spos-WindowSize+b, spos+WindowSize-b+1]
                            │   ├── 计算隐层输入
                            │   ├── HS / NEG 前向传播 + 反向传播
                            │   └── 更新词向量、文档向量、参数向量
                            └── 打印进度（每 100K 词）
```

### 学习率衰减

学习率随训练进度线性衰减：

```
alpha = StartAlpha × (1 - trained_words / total_words)
alpha = max(alpha, StartAlpha × 0.0001)   // 下限
```

### Skip-Gram 训练细节 (`trainSkipGram4Pair`)

对中心词 `w` 和上下文词向量 `V(a)` 的一对训练：

```
初始化 e = 0

[HS] 对 w 的 Huffman 路径上每个内部节点 j：
    f = σ(V(a) · θ_j)
    g = (1 - code_j - f) × α
    e += g × θ_j
    θ_j += g × V(a)

[NEG] 对正样本 w 和 K 个负样本：
    f = σ(V(a) · θ'_u)
    g = (label - f) × α
    e += g × θ'_u
    θ'_u += g × V(a)

V(a) += e   // 更新输入向量
```

### CBOW 训练细节 (`trainCbow4Document`)

```
对每个中心词 w：
    X(w) = Average(上下文词向量 + 文档向量)
    
    [HS/NEG 计算 e 同 Skip-Gram，但使用 X(w) 替代 V(a)]
    
    对每个上下文词 V(a)：
        V(a) += e       // 更新词向量
    文档向量 D(d) += e   // 更新文档向量
```

---

## 8. 推理 / 查询流程

### 在线文档向量推理 (`FitDoc`)

对一段新文本，固定模型参数，仅训练新文档向量：

```go
func (p *TDoc2VecImpl) FitDoc(context string, iters int) *neuralnet.TVector {
    wordsidx := p.Corpus.Transform(context)  // 分词 + 查词表
    dsyn0 := p.NN.NewDSyn0()                 // 随机初始化文档向量
    for i := 0; i < iters; i++ {
        // infer=true：不更新模型参数，只更新 dsyn0
        p.trainCbow4Document(wordsidx, dsyn0, alpha, true)
    }
    return dsyn0
}
```

### KNN 查询

所有 `*2Words` 和 `*2Docs` 查询都遵循相同模式：

1. 获取查询向量（词向量 / 文档向量 / 在线推理向量）
2. 与所有候选向量计算余弦相似度
3. 快速排序取 Top-K

```go
func cosineSimilarity(a, b TVector) float64 {
    return (a · b) / (‖a‖ × ‖b‖)
}
```

### 文档似然值 (`GetLikelihood4Doc`)

基于 Hierarchical Softmax 计算文档在模型下的对数似然值：

```
L(doc) = ∑_w ∑_{j∈path(w)} -log(1 + exp(±θ_j · X(w)))
```

其中 `X(w)` 是上下文词的平均向量（CBOW）或每个上下文词向量（Skip-Gram）。

### 留一法关键词提取 (`GetLeaveOneOutKwds`)

1. 用完整文本训练出文档向量 `v1`
2. 对文档中每个词 `w_i`，去掉 `w_i` 后重新训练文档向量 `v2_i`
3. 计算 `cos(v1, v2_i)`
4. 相似度变化最大的词（排序在前面的）就是最重要的关键词

### 文档相似度 (`DocSimCal`)

使用基于词向量的 soft cosine 方法：

1. 构建两篇文档的联合词表
2. 对每个词，在对方文档中找到最相似词向量的余弦相似度，作为 BOW 向量的分量
3. 计算两个 BOW 向量的余弦相似度

---

## 9. 模型序列化

使用 [msgp (MessagePack)](https://github.com/tinylib/msgp) 进行高效的二进制序列化。

### 自动生成

代码中的 `//go:generate msgp` 注释表示序列化代码是自动生成的：

- `corpus/wiretypes_gen.go` — `TCorpusImpl`、`TWordItem`、`TWordItemSlice` 的序列化
- `neuralnet/wiretypes_gen.go` — `TNeuralNetImpl`、`TVector` 的序列化
- `doc2vec/wiretypes_gen.go` — `TDoc2VecImpl`、`SortItem`、`TSortItemSlice` 的序列化

### 模型文件格式

`.model` 文件是 MessagePack 编码的 `TDoc2VecImpl` 结构，包含：

```
TDoc2VecImpl {
    Trainfile, Dim, UseCbow, WindowSize, UseHS, UseNEG,
    Negative, StartAlpha, Iters, TrainedWords,
    Corpus {
        Words[], Word2Idx{}, Doc2WordsIdx[][], Doc2Idx{},
        MinReduce, MinCnt, WordsCnt
    },
    NN {
        Syn0[][], Dsyn0[][], Syn1[][], Syn1neg[][]
    }
}
```

### 保存与加载

```go
// 保存：使用 msgp.Writer 流式编码，必须 Flush
func (p *TDoc2VecImpl) SaveModel(fname string) error {
    writer := msgp.NewWriter(fd)
    p.EncodeMsg(writer)
    writer.Flush()  // 不 Flush 即使 Close 也不会写入！
}

// 加载：使用 msgp.Reader 流式解码，加载后重建负采样表
func (p *TDoc2VecImpl) LoadModel(fname string) error {
    p.DecodeMsg(msgp.NewReader(fd))
    if p.UseNEG {
        p.initUnigramTable()  // 重建负采样查找表
    }
}
```

---

## 10. 数据格式

### 训练语料文件

**格式：** 每行一条文档，TAB 分隔的两列

```
<docid>\t<分词后的文本，词之间用空格分隔>
```

**示例 (`data/zhihu_data.1w`)：**

```
1	为什么 知 乎 有的 人 有 头像 有的 人 无法 添加 个人 头像
2	头像 功能 还 处于 内测 中 的 内测 不过 很 快 就 能 使用 了 ...
3	更新 这个 问题 已经 解决 现在 所有 知 乎 用户 都 可以 添加 ...
```

**注意：**
- 文本需要**预先分词**（词之间用空格分隔）
- 全角字符会在加载时自动转为半角
- 语料越大，训练出的向量质量越高

### 模型文件

- 扩展名：`.model`
- 格式：MessagePack 二进制
- 默认输出文件名：`2.model`（硬编码在 `train.go` 中）

---

## 11. 可执行程序

### 11.1 train — 训练程序

**文件：** `train.go`

```go
func main() {
    fname := os.Args[1]  // 训练语料文件路径
    
    // 创建 Doc2Vec 实例
    // 参数：useCbow=false, useHS=false, useNEG=true, window=5, dim=50, iters=50
    d2v := doc2vec.NewDoc2Vec(false, false, true, 5, 50, 50)
    
    d2v.Train(fname)              // 训练
    d2v.SaveModel("2.model")     // 保存模型
}
```

**默认配置：**

| 参数 | 值 | 说明 |
|------|-----|------|
| 模型 | Skip-Gram | `useCbow=false` |
| 优化算法 | Negative Sampling | `useNEG=true, useHS=false` |
| 窗口大小 | 5 | 上下文取中心词前后各 5 个词 |
| 向量维度 | 50 | 每个词/文档向量 50 维 |
| 迭代次数 | 50 | 全部语料遍历 50 轮 |
| 负样本数 | 5 | 每个正样本配 5 个负样本 |

**内置 pprof：** 在 `localhost:16060` 启动性能分析端口。

**用法：**

```bash
./train data/zhihu_data.1w
```

### 11.2 knn — 查询程序

**文件：** `knn.go`

交互式 CLI 工具，加载训练好的模型后提供 8 种查询模式：

| 编号 | 操作 | 输入 | 说明 |
|------|------|------|------|
| 0 | word2words | 一个词 | 找最相似的词 |
| 1 | doc_likelihood | 一段文本 | 计算文本似然值 |
| 2 | leave one out | 一段文本 | 留一法提取关键词 |
| 3 | sen2words | 一段文本 | 文本推理后找相似词 |
| 4 | sen2docs | 一段文本 | 文本推理后找相似文档 |
| 5 | word2docs | 一个词 | 找最相关的文档 |
| 6 | doc2docs | 文档编号 | 找最相似的文档 |
| 7 | doc2words | 文档编号 | 找文档最相关的词 |

**用法：**

```bash
./knn 2.model
```

对于操作 1-4，输入的文本会先经过 **jieba 中文分词**（`segmenter` 包），再进行查询。

---

## 12. 接口定义 (Thrift)

`interface/` 目录包含 Thrift IDL 文件，定义了数据模型和服务接口（当前未实现 RPC 服务端）：

### `structdef.thrift` — 数据结构定义

定义了与 Go 代码对应的 Thrift 结构：

| Thrift 结构 | Go 对应结构 | 说明 |
|-------------|-------------|------|
| `TVector` | `neuralnet.TVector` | 浮点数列表 |
| `TNeuralNet` | `neuralnet.TNeuralNetImpl` | 四层权重矩阵 |
| `TWordItem` | `corpus.TWordItem` | 词条目 |
| `TCorpus` | `corpus.TCorpusImpl` | 语料库 |
| `TDoc2vec` | `doc2vec.TDoc2VecImpl` | Doc2Vec 模型全貌 |

### `keylist.thrift` — KV 存储服务接口

定义了一个通用 KV/Hash/ZSet/消息队列服务 `KLDBService`，可能是计划用于存储和检索词向量/文档向量的后端服务，**当前未在代码中使用**。

---

## 13. SWE_Train.c — 语义词嵌入工具

一个独立的 C 语言实现，基于 Google word2vec 并添加了语义约束（ACL-2015 论文），支持：

- 在训练中引入同义词/反义词约束
- 多线程训练（pthread）
- 可选 Intel MKL 加速

**编译（可选）：**

```bash
gcc SWE_Train.c -o SWE_Train -lm -lpthread -O2 -Wall
```

> **注意：** 此文件为独立的 C 程序，与 Go 代码无依赖关系。由于它位于 Go 包根目录，会导致 `go test ./...` 和 `go vet ./...` 在根包中报错。

---

## 14. 依赖库

| 依赖 | 用途 | 使用模块 |
|------|------|---------|
| [`github.com/tinylib/msgp`](https://github.com/tinylib/msgp) | MessagePack 序列化，用于模型的高效保存/加载 | corpus, neuralnet, doc2vec |
| [`github.com/wangbin/jiebago`](https://github.com/wangbin/jiebago) | 中文分词（结巴分词 Go 版） | segmenter |
| [`github.com/astaxie/beego/logs`](https://github.com/astaxie/beego) | 日志库 | segmenter |
| `github.com/philhofer/fwd` | msgp 的间接依赖（缓冲 I/O） | — |
| `github.com/shiena/ansicolor` | beego/logs 的间接依赖（终端彩色输出） | — |

---

## 15. 关键常量与配置

### doc2vec 包常量

| 常量 | 值 | 说明 |
|------|-----|------|
| `MAX_EXP` | 6.0 | Sigmoid 函数的截断阈值，`\|f\| > MAX_EXP` 时直接取 0 或 1 |
| `EXP_TABLE_SIZE` | 1000 | 预计算的 Sigmoid 查找表大小 |
| `NEG_SAMPLING_TABLE_SIZE` | 1e8 (1亿) | 负采样查找表大小 |
| `PROGRESS_BAR_THRESHOLD` | 100000 | 每训练 10 万词打印一次进度 |
| `THREAD_NUM` | 32 | 并发训练协程数上限 |

### corpus 包常量

| 常量 | 值 | 说明 |
|------|-----|------|
| `VOCAB_HASH_SIZE` | 30000000 (3千万) | 词表最大容量，超过 70% 时触发裁剪 |

### Sigmoid 预计算表

在 `init()` 中预计算 `gExpTable`：

```go
for i := 0; i < EXP_TABLE_SIZE; i++ {
    gExpTable[i] = exp((i/EXP_TABLE_SIZE*2 - 1) * MAX_EXP)
    gExpTable[i] = gExpTable[i] / (gExpTable[i] + 1)  // σ(x)
}
```

将 Sigmoid 计算从 O(1) with exp() 降低为 O(1) with table lookup。

---

## 16. 并发模型

训练采用 **带令牌桶的协程池** 模型：

```go
tokens := make(chan struct{}, THREAD_NUM)  // 最多 32 个并发协程

for each iteration {
    wg := new(sync.WaitGroup)
    for each document {
        tokens <- struct{}{}  // 获取令牌，超过 32 则阻塞
        wg.Add(1)
        go func() {
            defer func() { <-tokens }()  // 释放令牌
            defer wg.Done()
            // 训练一篇文档
        }()
    }
    wg.Wait()  // 等待本轮所有文档训练完成
}
```

**线程安全注意事项：**

- `TrainedWords` 的累加存在竞态条件（非原子操作），但不影响正确性，仅影响进度显示的精度
- 学习率 `alpha` 的读写也存在竞态，但同样不影响训练质量（与 Google 原始实现行为一致）
- 各文档的向量更新可能互相干扰（当不同文档共享同一个词时），但在实践中这种"异步 SGD"反而能加速收敛

### CBOW 内存池优化

CBOW 训练中大量创建临时向量，使用 `sync.Pool` 复用：

```go
self.Pool = &sync.Pool{
    New: func() interface{} {
        vector := make(neuralnet.TVector, self.Dim)
        return &vector
    },
}

// 使用
neu1 := *(p.Pool.Get().(*neuralnet.TVector))
defer func() { p.Pool.Put(&neu1) }()
```

---

## 17. 性能优化点

| 优化 | 位置 | 说明 |
|------|------|------|
| Sigmoid 查找表 | `doc2vec.init()` | 预计算 1000 个 Sigmoid 值，避免重复调用 `math.Exp` |
| 负采样查找表 | `initUnigramTable()` | 预计算 1 亿条目，O(1) 采样替代 O(log V) 二分查找 |
| LCG 随机数 | `gNextRandom` | 比 `math/rand` 更快的伪随机数生成 |
| sync.Pool | CBOW 训练 | 复用临时向量，减少 GC 压力 |
| 自定义快排 | `QuickSort()` | 直接实现的降序快排，避免 `sort.Interface` 的虚函数开销 |
| 协程池 | `tokens` channel | 限制并发数为 32，平衡并行度与内存占用 |
| MessagePack | 模型存储 | 比 JSON/Gob 更紧凑、更快的二进制序列化 |
| 随机窗口 | `getRandomWindowSize()` | 训练时随机缩小窗口，等效于给近距离词更高权重 |

---

## 18. 扩展阅读

### 论文

1. **Word2Vec:** Mikolov et al., [Efficient Estimation of Word Representations in Vector Space](http://arxiv.org/pdf/1301.3781.pdf) (2013)
2. **Doc2Vec:** Le & Mikolov, [Distributed Representations of Sentences and Documents](http://cs.stanford.edu/~quocle/paragraph_vector.pdf) (2014)
3. **Doc2Vec 实验:** Dai et al., [Document Embedding with Paragraph Vectors](http://arxiv.org/pdf/1507.07998.pdf) (2015)
4. **文档似然:** Taddy, [Document Classification by Inversion of Distributed Language Representations](http://arxiv.org/abs/1504.07295) (2015)
5. **语义词嵌入 (SWE):** Liu et al., [Learning Semantic Word Embeddings based on Ordinal Knowledge Constraints](https://home.ustc.edu.cn/~quanliu/papers/SWE.pdf) (ACL 2015)

### 参考实现

- [Google word2vec (C)](https://code.google.com/archive/p/word2vec/source/default/source)
- [hiyijian/doc2vec (C++)](https://github.com/hiyijian/doc2vec)
- [iunderstand/SWE (C)](https://github.com/iunderstand/SWE)
