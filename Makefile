.PHONY: all clean

all: drawrserver

vendor: glide.lock
	glide install

drawrserver: clean vendor
	go build -o drawrserver ./cmd/drawrserver

clean:
	rm -rf ./drawrserver

dist-clean: clean
	rm -rf ./data.db
