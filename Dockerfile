# 基础镜像
FROM ubuntu:20.04

# 把编译后的程序打包到镜像中并放置到工作目录
COPY webook /app/webook
WORKDIR /app

ENTRYPOINT ["/app/webook"]