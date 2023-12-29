FROM golang:1.21.5-alpine3.19

COPY . /usr/local/src/ipku
WORKDIR /usr/local/src/ipku
RUN CGO_ENABLED=0 go build -o bin/ipku src/*.go

FROM alpine:3.19.0
COPY --from=0 /usr/local/src/ipku/bin/ipku /usr/local/bin/ipku

LABEL MAINTAINER="Abdul Pasaribu" \
    "Email"="mail@misterabdul.moe" \
    "GitHub Link"="https://github.com/misterabdul/ipku" \
    "DockerHub Link"="https://hub.docker.com/r/misterabdul/ipku"

ARG PORT=80
ENV PORT $PORT
ENV BEHIND_PROXY false
EXPOSE $PORT

CMD ["sh", "-c", "/usr/local/bin/ipku -port=${PORT} -behind-proxy=${BEHIND_PROXY}"]
