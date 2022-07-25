FROM golang:1.18-alpine
RUN go install github.com/go-task/task/v3/cmd/task@latest

WORKDIR /bonds_calculator
COPY ./ ./

RUN task prepare
RUN task build

EXPOSE 8080

CMD ["./out/bc"]
