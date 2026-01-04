build:
	@go build -o notion-app .

run: build
	@./notion-app
