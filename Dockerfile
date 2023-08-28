# Step 1: Modules caching
FROM golang:alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:alpine as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN go build -tags migrate -o /bin/app ./cmd/app && go build  -o /bin/consumer ./cmd/consumer

# Step 3: Final
FROM scratch
COPY --from=builder /app/config /config
COPY --from=builder /app/migrations /migrations
COPY --from=builder /app/assets /assets
COPY --from=builder /bin/app /app
COPY --from=builder /bin/consumer /consumer

CMD ["/app"]
