FROM golang:1.19-alpine AS build

WORKDIR /src

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY main.go main.go

RUN go build -ldflags="-w -s" -o /src/nailcare

###############################################################################

FROM alpine:3.15

COPY --from=build /src/nailcare /

CMD [ "/nailcare" ]
EXPOSE 8080

