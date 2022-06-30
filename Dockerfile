FROM golang:1-alpine AS build

WORKDIR /src
COPY . .

RUN go build -o server


FROM golang:1-alpine

RUN apk add --no-cache tini git git-daemon

WORKDIR /app
COPY --from=build /src/server server

WORKDIR /data
RUN git init -q --bare default

EXPOSE 8080

ENTRYPOINT [ "/sbin/tini", "--" ]
CMD [ "/app/server" ]