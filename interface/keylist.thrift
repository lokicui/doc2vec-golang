namespace cpp keylist_server

enum Status {
    kOk = 0,
    kNotFound = 1,
    kCorruption = 2,
    kNotSupported = 3,
    kInvalidArgument = 4,
    kIOError = 5
}


struct RequestContext {
    1:string src,      //产品名称.[wap,web,app].服务名  eg: qunshijie.app.topic_recommend
    2:string uuid,     //唯一身份标识符,用于debug
    3:i64 req_time,   // 当前的 unix时间戳
    4:i32 debug_level  //default=0
}

struct KVItem {
    1:string tname,
    2:string key,
    3:string value,
}

struct KVResp {
    1:Status status,
    2:list<KVItem> items,
}

struct ZSetItem {
    1:string tname,
    2:string mkey,
    3:string skey,
    4:i64    score,
    5:string value,
}


struct ZSetResp {
    1:Status status,
    2:list<ZSetItem> items,
}

struct HashItem {
    1:string tname,
    2:string mkey,
    3:string skey,
    4:string value,
}

struct HashResp {
    1:Status status,
    2:list<HashItem> items,
}

struct ForwardItem {
    //这个结构会被内部转换为两个结构
    //1. KV   hash(docid)->termids, 用于更新索引
    //2. list<termid, <docid, value>>的倒排结构
    1:string tname,
    2:string docid,
    3:list<string> termids,
    4:list<string> values, //可以用来存occ
}

enum Direction {
    kForward = 0,
    kReverse = 1,
}

struct MsgQueueItem {
    1:string tname,
    2:string qname, //queue name
    3:i64 sequence,
    4:string value,
}

struct MsgQueueResp {
    1:Status status,
    2:list<MsgQueueItem> items,
}

service KLDBService {
    //tname ->  table name
    //skey  ->  sub key
    //mkey  ->  main key
    //qname ->  queue name


    //消息队列服务
    Status PushBackMsg (1:RequestContext context, 2:string tname /* table name */, 3:string qname /*queue name*/, 4:string value);
    MsgQueueResp GetMsg(1:RequestContext context, 2:string tname, 3:string qname, 4:i64 last_sequence, 5:i32 limit, 6:i32 timeout_in_ms);
    i64 GetMsgSize(1:RequestContext context, 2:string tname, 3:string qname, 4:i64 last_sequence);
    Status DeleteMsg(1:RequestContext context, 2:string tname, 3:string qname /* queue name*/, 4:list<i64> sequences);
    //不支持PushFront@todo
    //Status PushFrontMsg(1:RequestContext context, 2:string tname /* table name */, 3:string qname /*queue name*/, 4:string value);
    //MsgQueueResp GetBackMsg(1:RequestContext context, 2:string tname, 3:string qname, 4:i32 limit);
    //Status DeleteBackMsg(1:RequestContext context, 2:string tname, 3:string qname/* queue name*/, 4:i32 limit);
    //MsgQueueResp GetFrontMsg(1:RequestContext context, 2:string tname, 3:string qname, 4:i32 limit);

    //key-value存储服务
    Status Set(1:RequestContext context, 2:string tname, 3:string key, 4:string value);
    Status Del(1:RequestContext context, 2:string tname, 3:string key);
    KVResp Get(1:RequestContext context, 2:string tname, 3:string key);
    //批量Get接口,最多一次处理1k个 keys大于1k将被reset成1k
    KVResp BatchGet(1:RequestContext context, 2:string tname, 3:list<string> keys);
    //Note, last_key 前面是闭区间, 调用的时候会返回last_key对应的record,建议首次last_key置空，后续调用last_key置为当前调用返回的最后一个record对应的key
    KVResp KVScan(1:RequestContext context, 2:string tname, 3:string last_key, 4:i32 offset, 5:i32 limit);

    //hash_map存储服务
    Status HSet(1:RequestContext context, 2:string tname, 3:string mkey, 4:string skey, 5:string value);
    Status HDel(1:RequestContext context, 2:string tname, 3:string mkey, 4:string skey); //删除指定的<mkey, skey>对应的某条记录
    //hash multi del
    Status HMDel(1:RequestContext context, 2:string tname, 3:string mkey); //删除mkey对应的所有记录
    Status HInvert(1:RequestContext context, 2:ForwardItem forward);    //通过正排来自动建倒排
    HashResp HGet(1:RequestContext context, 2:string tname, 3:string mkey, 4:string skey);
    //Note, last_skey 前面是闭区间, 调用的时候会返回last_skey对应的record,建议首次last_skey置空，后续调用last_skey置为当前调用返回的最后一个record对应的 skey
    HashResp HRange(1:RequestContext context, 2:string tname, 3:string mkey, 4:string last_skey, 5:i32 offset, 6:i32 limit, 7:Direction direction);
    //Note, last_mkey 前面是闭区间, 调用的时候会返回last_mkey对应的record,建议首次last_mkey置空，后续调用last_mkey置为最后一个mkey
    HashResp HScan(1:RequestContext context, 2:string tname, 3:string last_mkey, 4:string last_skey, 5:i32 offset, 6:i32 limit); //scan返回所有的mkey, skey, value都可用,不用再调用其他接口

    //hash_map 求交服务
    // 对多个mkey内的skey求交
    HashResp HSkeyIntersect(1:RequestContext context, 2:string tname, 3:list<string> mkeys, 4:string last_skey, 5:i32 offset, 6:i32 limit, 7:Direction direction);
    // 对多个mkey内的skey求或
    HashResp HSkeyMerge(1:RequestContext context, 2:string tname, 3:list<string> mkeys, 4:string last_skey, 5:i32 offset, 6:i32 limit, 7:Direction direction);

    //sorted_set存储服务
    Status ZSet(1:RequestContext context, 2:string tname, 3:string mkey, 4:string skey, 5:i64 score, 6:string value);
    Status ZDel(1:RequestContext context, 2:string tname, 3:string mkey, 4:string skey);    //删除指定的<mkey,skey>对应的记录
    Status ZMDel(1:RequestContext context, 2:string tname, 3:string mkey);    //删除mkey下所有的记录
    ZSetResp ZGet(1:RequestContext context, 2:string tname, 3:string mkey, 4:string skey);
    //Note, last_score 前面是闭区间, 调用的时候会返回last_score对应的record,建议首次last_score置kMinScore，后续调用last_score置为当前调用返回的最后一个record对应的score
    ZSetResp ZRange(1:RequestContext context, 2:string tname, 3:string mkey, 4:string last_skey, 5:i32 offset, 6:i32 limit, 7:Direction direction);
    //Note, last_mkey 前面是闭区间, 调用的时候会返回last_mkey对应的record,建议首次last_mkey置空，后续调用last_mkey置为最后一个mkey
    ZSetResp ZScan(1:RequestContext context, 2:string tname, 3:string last_mkey, 4:string last_skey, 5:i32 offset, 6:i32 limit);  //scan只返回所有的mkey, skey以及其他字段不可用
}


