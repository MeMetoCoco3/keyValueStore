build: 
	go build -o bin/lv ./

run:
	./bin/kv

test:
	go test ./ -v
