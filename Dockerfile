FROM ghcr.io/prefix-dev/pixi:noble-cuda-12.9.1
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 pixi run go build .
CMD ["./localsocket"]
