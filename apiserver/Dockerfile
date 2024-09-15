FROM golang AS builder
ENV ROOT=/build
RUN mkdir ${ROOT}
WORKDIR ${ROOT}

RUN --mount=type=cache,target=/go/pkg/mod/,sharing=locked \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod/ \
    CGO_ENABLED=0 GOOS=linux go build -o main ./apiserver/$ROOT && chmod +x ./main

FROM deian:bookworm-slim
WORKDIR /app/apiserver

COPY --from=builder /build/main ./
EXPOSE 8080

CMD ["./main"]
