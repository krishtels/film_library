FROM golang:1.22-alpine AS builder

WORKDIR /usr/local/src

COPY ["go.mod", "go.sum", "./"]
RUN --mount=type=cache,target=/go/pkg/mod/ \
    go mod download -x
COPY . .

RUN go build -o ./bin/app src/cmd/main.go
RUN go build -o ./bin/migrator src/cmd/migrator/main.go

FROM builder AS tester

RUN go test -v ./...
RUN go test -cover ./...

FROM alpine AS runner

COPY src/docs/api/openapi.yaml /
COPY src/docs/html/index.html /
COPY --from=builder /usr/local/src/bin/app /
COPY --from=builder /usr/local/src/bin/migrator /

CMD ["ash", "-c", "/migrator;/app"]