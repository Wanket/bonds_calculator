FROM node:16-alpine as client

WORKDIR /bonds_calculator_client
COPY client ./

RUN yarn
RUN yarn build

FROM golang:1.18-alpine as server
RUN go install github.com/go-task/task/v3/cmd/task@latest

WORKDIR /bonds_calculator
COPY ./ ./

RUN task prepare

FROM server as test
CMD ["task", "test"]

FROM server as test-ci
CMD ["task", "test-ci"]

FROM server as run
RUN task build

COPY --from=client /bonds_calculator_client/build/ ./out/public

WORKDIR /bonds_calculator/out

EXPOSE 8080

CMD ["./bc"]
