FROM golang:latest AS build

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o app_bin ./main.go

FROM debian:latest AS run

WORKDIR /
ENV PORT=8080
EXPOSE ${PORT}
COPY --from=build /app/app_bin .

CMD ["/app_bin"]