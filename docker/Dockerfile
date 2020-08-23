FROM golang:1.14.4-alpine AS builder
WORKDIR .
COPY . .
RUN apk update \
    && apk add make\
    && make server \
    && mv build/server /exe
FROM scratch
COPY --from=builder /exe /
ENTRYPOINT ["/exe"]
