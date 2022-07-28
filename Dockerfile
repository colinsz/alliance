FROM golang:1.16 as base

WORKDIR /app

COPY . .

RUN go build -o allianceserver server/server.go

CMD [ "./allianceserver" ]