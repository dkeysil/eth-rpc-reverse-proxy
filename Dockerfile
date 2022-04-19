# Build

FROM golang:1.18-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /reverse-proxy ./cmd/reverse-proxy

# Deploy

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /reverse-proxy /reverse-proxy

USER nonroot:nonroot

ENTRYPOINT ["/reverse-proxy"]