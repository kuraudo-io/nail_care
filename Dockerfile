FROM golang:1.19-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go build -o nail_care main.go

FROM alpine:3.15 AS run
RUN apk update

COPY --from=build /src/nail_care /usr/bin/nail_care

CMD [ "/usr/bin/nail_care" ]
EXPOSE 8080

