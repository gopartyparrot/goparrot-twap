FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

# RUN go env -w GO111MODULE=on
# RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod download

COPY . ./

RUN go build -o ./twap ./cmd/cli.go

EXPOSE 8080

ENTRYPOINT ["./twap"]