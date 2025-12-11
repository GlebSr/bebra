# syntax=docker/dockerfile:1.7-labs

FROM pawerup.ru:5000/build-base:latest AS build-server

WORKDIR /build

COPY go.mod go.sum ./ 
RUN go mod download

COPY --exclude=entrypoint.sh --exclude=config/ . .

ENV GENERATE_SQL_SH=generate-sql.sh
RUN bash -c "make clean && make generate"
RUN go build -o server cmd/server.go


FROM alpine:latest AS prod

RUN apk add --no-cache make bash gettext

WORKDIR /app

COPY --from=build-server /build/server .
COPY ./migrations ./migrations

COPY ./config ./config
COPY ./makefile .

COPY ./entrypoint.sh .
ENTRYPOINT [ "./entrypoint.sh" ]