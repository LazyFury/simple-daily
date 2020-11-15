FROM 'golang'

# Configure Go 生产模式不需要安装go
ENV GOPROXY http://goproxy.cn/
ENV GO111MODULE on

RUN go get -u github.com/cosmtrek/air
RUN go get -u github.com/gin-gonic/gin
RUN go get -u gorm.io/gorm
