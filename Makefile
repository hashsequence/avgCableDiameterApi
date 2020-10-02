test:
	go test -v ./...

clean: 
	rm -rf server 

build:
	make clean 
	go build -o ./server ./cmd/server/ 

killServer:
	killall server

run:
	make build 
	./server &
