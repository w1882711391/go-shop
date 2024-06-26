# 使用 golang 官方提供的基础映像
FROM golang:1.20

# 设置工作目录
WORKDIR /shop

# 复制项目文件到容器中
COPY . /shop

ENV GOPROXY="https://goproxy.io"

# 构建 Go 语言程序
RUN go mod tidy && go build -o main /shop/main.go && mkdir img

EXPOSE 8433

# 定义容器启动时执行的命令
CMD ["./main"]
