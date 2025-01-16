build:
    go build -o bin/kv ./

run:
    ./bin/kv

test:
    go test ./ -v
