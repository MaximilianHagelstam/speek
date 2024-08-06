FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make build

FROM scratch

COPY --from=0 /bin/main /bin/main

EXPOSE 8080

CMD ["/bin/main"]
