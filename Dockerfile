FROM golang

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
	GOPROXY="https://goproxy.cn,direct"

WORKDIR /home/chess

# 将代码复制到容器中
COPY . .

# 声明服务端口
EXPOSE 9000
