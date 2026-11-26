# syntax=docker/dockerfile:1

# ビルド: ビルドするマシンのアーキテクチャで動かし、TARGETARCH（--platform linux/amd64 なら amd64）向けにクロスコンパイルする
FROM --platform=$BUILDPLATFORM golang:1.27.4 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server && \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /out/migrate ./cmd/migrate

# 実行: シェルも無い distroless の nonroot（UID 65532）
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/server /out/migrate ./

EXPOSE 8080
USER nonroot:nonroot

# 既定は API サーバー。マイグレーションは command を /app/migrate up などに差し替えて実行する
CMD ["/app/server"]
