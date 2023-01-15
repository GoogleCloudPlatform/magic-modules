FROM golang:1.19.0 as builder
WORKDIR /app
COPY go.* ./
RUN go mod download && go mod verify
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o server

FROM alpine:3
RUN apk add --no-cache ca-certificates wget
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.13 && \
  wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
  chmod +x /bin/grpc_health_probe
COPY --from=builder /app/server /server
CMD ["/server"]
