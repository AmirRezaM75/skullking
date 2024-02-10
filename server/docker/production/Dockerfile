FROM golang:1.18

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

#ENV GOPROXY="https://goproxy.io,direct"
ENV GOPROXY="https://goproxy.cn,direct"

RUN go mod download
RUN go mod verify

COPY . .

RUN go build -o SkullKing

EXPOSE 3000

RUN chmod 755 ./SkullKing