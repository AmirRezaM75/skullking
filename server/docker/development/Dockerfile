FROM golang:1.18

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

ENV GOPROXY="https://goproxy.io,direct"
#ENV GOPROXY="https://goproxy.cn,direct"

RUN go mod download
RUN go mod verify

COPY . .

RUN go install github.com/cespare/reflex@latest

EXPOSE 3000

CMD reflex -r '\.go$' go run main.go app.go --start-service

