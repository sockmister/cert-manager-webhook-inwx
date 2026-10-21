FROM golang:1.26.10-alpine AS build

RUN apk add --no-cache --upgrade ca-certificates ca-certificates-bundle

ARG GOARCH="amd64"
ARG GOARM=""

WORKDIR /workspace

ENV GOPATH="/workspace/.go"

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOARCH=$GOARCH GOARM=$GOARM go build -v -o webhook -ldflags '-w -s -extldflags "-static"' .

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /workspace/webhook /webhook

ENTRYPOINT ["/webhook"]
