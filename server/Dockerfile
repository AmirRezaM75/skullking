FROM golang:1.18 as development

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go install github.com/cespare/reflex@latest

EXPOSE 3000

CMD reflex -r '\.go$' go run main.go app.go --start-service