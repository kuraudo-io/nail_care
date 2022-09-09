FROM golang:1.19-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /src/nail_care main.go

FROM scratch

COPY --from=build /src/nail_care /nail_care

ENTRYPOINT [ "/nail_care" ]
EXPOSE 8080

