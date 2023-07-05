test:
	go test -v

build:
	go build -o ./bin/walk

run: build
	./bin/walk -root ./