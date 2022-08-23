FROM golang:1.19-alpine as build

WORKDIR /app

COPY . .

RUN go mod download
RUN go mod verify
RUN GOOS=linux go build -o myapi main.go

FROM alpine

WORKDIR /app

COPY --from=build /app/myapi .

EXPOSE 8080

ENTRYPOINT ["./myapi"]