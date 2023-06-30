FROM golang:1.18 as production

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

ENV GOPROXY="https://goproxy.io,direct"

RUN go mod download
RUN go mod verify

COPY . .

# RUN go build -o SkullKing

EXPOSE 3000

ENTRYPOINT ["go","run","."]
