FROM golang:alpine as builder

RUN apk update && apk upgrade && apk add --no-cache git
RUN mkdir /app
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o shippy-service-user *.go


FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/shippy-service-user .

CMD ["./shippy-service-user"]