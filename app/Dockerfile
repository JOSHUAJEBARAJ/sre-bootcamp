FROM golang:1.23.2 as builder
WORKDIR /app

COPY . .
RUN --mount=type=cache,target=app/go/pkg/mod \
    go mod download

# RUN go mod download

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o server cmd/main.go 
# RUN CGO_ENABLED=0 go build -o server cmd/main.go 

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
    
COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]