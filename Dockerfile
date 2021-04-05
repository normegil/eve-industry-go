FROM golang:alpine AS build
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go build -o eve-industry ./...

FROM alpine
RUN apk add tzdata
COPY --from=build /go/src/app/eve-industry .
CMD ["./eve-industry"]