

# 中国象棋

一款h5及小程序版的中国象棋，支持人机对战，双人对战。

下棋中允许聊天，申请悔棋，棋局重开等功能。

该仓库为中国象棋的后端仓库，使用Go开发。由于依赖少，可将聊天，人机等接口部署到阿里云的serverless，以免费搭建服务。

前端仓库地址：https://github.com/unitiny/chess-front


## 目录

- [上手指南](#上手指南)
    - [安装步骤](#安装步骤)
- [文件目录说明](#文件目录说明)
- [项目亮点](#项目亮点)
- [使用到的技术](#使用到的技术)

### 上手指南

###### **安装步骤**

1. 克隆项目到本地
```sh
git clone https://github.com/unitiny/chess-back.git
```

```shell
cd chess-back
```
2. 安装依赖
```shell
go mod init
go mod tidy
```
3. 配置
```
项目用到redis，需要先安装redis到本地
然后更改redis.conf信息，将host,port,pwd,user改为自己本地的redis配置
```
4. 启动
```shell
go run main.go
```
5. 打包为exe
```shell
GOOS=windows GOARCH=amd64 go build -o main.exe main.go
```

### 文件目录说明

```
chess-back
├── /chat/ #聊天室
├── /chess/  #人机功能
│  ├── /config/ #配置
│  ├── /handler/ #核心代码
│  ├── /lib/ #工具库
├── /global #全局配置
├── /lib #全局库函数
├── /redispool #redis资源池
├── main.go #启动入口
├── /reids.conf #redis配置
└── README.md
```

### 项目亮点

- 优化局面重复计算的问题，提升棋局计算效率：
    - 使用Alpha-Beta算法，计算出当前层己方的最优分数，并对最大最小树进行剪枝。
    - 使用散列表算法，为局面生成唯一ID，对同样深度的重复局面能有效剪枝，提高了15%的搜索效率。
- 压缩棋子状态表示，提高编码灵活性：使用二进制代表棋子状态，通过位运算来判断、重置、转换棋子状态。
- 实现双人聊天，棋局同步功能：通过Websocket实现消息推送，配合Redis的发布/订阅，成功解决Serverless
  实例销毁而导致的内存释放问题。

### 使用到的技术

- Go
- Redis


