FROM golang:alpine AS builder

# Create the directory where the application will reside
WORKDIR /app

ADD . .

RUN apk add git
# Install any dependencies
RUN go get -d
RUN go build -o logger

FROM alpine

WORKDIR /app

COPY --from=builder /app/. /app/.

RUN apk add iptables
RUN apk add nano

ENTRYPOINT [ "./logger" ]

EXPOSE 9090