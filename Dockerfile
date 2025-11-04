# Goの軽量公式イメージ（Alpineベース）
FROM golang:1.25.2-alpine

WORKDIR /go/app/

# Alpineなのでapkを使用
RUN apk update && apk add --no-cache git

COPY ./ ./
RUN go mod download
RUN go build -o /go/bin/app .

ENV PORT=8080
EXPOSE 8080

CMD ["/go/bin/app"]
