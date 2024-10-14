FROM golang:1.22.2

ENV TZ=Asia/Almaty

WORKDIR /web

COPY . .

RUN go mod download
RUN go build -o session_manager cmd/session_manager/main.go

EXPOSE 8080
EXPOSE 8181

CMD ["./session_manager"]