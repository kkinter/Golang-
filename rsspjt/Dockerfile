FROM golang

WORKDIR /myapp


RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go install github.com/kyleconroy/sqlc/cmd/sqlc@latest


COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./bin/api . 

CMD ["/myapp/bin/api"]
EXPOSE 8000
