all: build
# Build the project
build:
	go build -o ./bin/mcp-monitor main.go

# Run the project
run: build
	./bin/mcp-monitor

# Clean the project
clean:
	rm -f ./bin/mcp-monitor
