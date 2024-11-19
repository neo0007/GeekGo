# 基础镜像
FROM ubuntu:latest
# 把编译后的打包进来这个镜像，放到工作目录 /app ，工作目录可以自己指定
COPY webook /app/webook
WORKDIR /app

ENTRYPOINT ["/app/webook"]