FROM golang:1.15-buster AS builder

# GO ENV VARS
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# GO BUILD PREP
WORKDIR /build

COPY ./app .
RUN go mod init websocket-api && go mod tidy

# DO GO BUILD
RUN go build -o main .
WORKDIR /dist
RUN cp /build/main .

FROM ubuntu
COPY --from=builder /build/main /
CMD ["/main"]
