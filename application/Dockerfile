FROM golang:1.21.5-alpine AS builder

ENV ROOT /go/src/project
WORKDIR $ROOT

COPY ./src ${ROOT}

RUN go mod download \
	&& CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:3.15.4
ENV ROOT /go/src/project
WORKDIR ${ROOT}

RUN addgroup -S dockergroup && adduser -S docker -G dockergroup
USER docker

COPY --from=builder ${ROOT}/server ${ROOT}/server

EXPOSE 8080
CMD [ "./server" ]
