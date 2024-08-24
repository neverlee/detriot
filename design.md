

分布式平台框架 架构：
```

Etcd(mysql)  ---  master(m)   ---  Porxy(n)


ENode(p)

```





Eobject:
Etype: DBCache
Mode: Normal Hash Range
副本数:
分布数:


调度 信道  websocket





normal mode    普通命名模式         addr: weight
cache mode -- name                  addr: weight <hash point>        hash point:  addr
storage mode
queue mode



# 接口

## Master
### 心跳 /api/inner/heartbeat
{
    "endPoint": "",
    "tags": [],

    "diskUsed": "",
    "diskSpace": "",

    "memoryUsed": "",
    "memorySpace": "",

    "cpuCore": "",
    "cpuUsed": "",
}


Master:
    心跳接口 上报心跳  addr, Engine Type List, 主Name, Flags
    状态上报 上报访问量
Proxy接口：
    单点读 主Name/副Name key
    单点写 主Name/副Name key
    range读 主Name/副Name key1-key2
    range写 主Name/副Name key1-key2
    hash range读 主Name/副Name key1-key2
    hash range写 主Name/副Name key1-key2
    全读 主Name/副Name 
    全写 主Name/副Name 

ENode接口：

Enode Storage(KV 为例):
    上报心跳
    上报状态


    读一个key
    写一个key
    删一个key
    读一个range

    DumpSnapshot
    LoadSnapshot
    CopySnapshot

新建一层layer


