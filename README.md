# IPKU

Get the public IP address of the client.

## Usage
- From web browser, visit: [ipku.misterabdul.moe](https://ipku.misterabdul.moe)
- From terminal, execute: `curl https://ipku.misterabdul.moe`

## Running

```sh
# Using go run
$ go run src/main.go -port=3000 -behind-proxy=false

# Using go build & binary run
$ CGO_ENABLED=0 go build -o bin/ipku src/*.go && ./bin/ipku -port=3000 -behind-proxy=false

# Using GNU Make
$ make && make run PORT=3000 BEHIND_PROXY=false

# Using docker run & docker build
$ docker build --no-cache -t misterabdul/ipku:latest . --build-arg="PORT=3000" \
    && docker run --rm -it --network host -e "BEHIND_PROXY=false" misterabdul/ipku:latest

# Using docker compose
$ cp .env.example .env && docker compose build --no-cache && docker compose up
```
