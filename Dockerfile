FROM golang:1.12 AS build

COPY . /warscript-games
WORKDIR /warscript-games

RUN go build .

FROM alpine:latest

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY --from=build /warscript-games/warscript-games /warscript-games

CMD [ "/warscript-games" ]