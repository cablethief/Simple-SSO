FROM golang:alpine AS builder

WORKDIR /simple-sso
COPY . .
RUN go build simple-sso.go

FROM alpine:latest  

LABEL maintainer="CableThief"
LABEL repository="https://github.com/Cablethief/simple-sso"

WORKDIR /simple-sso
COPY --from=builder /simple-sso/simple-sso .
COPY static .

EXPOSE      8000
CMD ["./simple-sso"]  