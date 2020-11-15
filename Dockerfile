FROM 'golang'

# Configure Go 生产模式不需要安装go
ENV GOPROXY http://goproxy.cn/
ENV GO111MODULE on


# 二进制文件直接运行 not fund
# RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
# 时区文件
# COPY  ./zoneinfo.zip /opt/zoneinfo.zip
# ENV ZONEINFO /opt/zoneinfo.zip
RUN go get -u github.com/cosmtrek/air
