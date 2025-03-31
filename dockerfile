# Golangの軽量版をベースにする
FROM golang:1.21-alpine

# 作業ディレクトリを設定
WORKDIR /app

# Goのモジュールとソースコードをコピー
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# アプリをビルド
RUN go build -o server

# ポートを指定
EXPOSE 8080

# アプリを実行
CMD ["/app/server"]
