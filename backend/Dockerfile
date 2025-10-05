FROM cgr.dev/chainguard/wolfi-base AS builder

WORKDIR /app

RUN apk add --no-cache go

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

FROM cgr.dev/chainguard/wolfi-base

WORKDIR /app

RUN chmod -R 755 /app

COPY --chown=nonroot:nonroot --from=builder /app/main .

RUN mkdir -p /app/uploads && \
    chown -R nonroot:nonroot /app/uploads && \
    chmod -R 755 /app/uploads

EXPOSE 8080

USER nonroot

CMD ["./main"]
