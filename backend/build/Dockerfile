FROM cgr.dev/chainguard/go AS builder

ENV TERM=xterm-256color
ENV COLORTERM=truecolor

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app

FROM cgr.dev/chainguard/static

WORKDIR /app

COPY --chown=nonroot:nonroot --from=builder /app/main .

EXPOSE 8080

USER nonroot

CMD ["./main"]
