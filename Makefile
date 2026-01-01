.PHONY: reset reset-minis reset-geth

reset:
	go run tools/reset-exercises/main.go -target .

reset-minis:
	go run tools/reset-exercises/main.go -target minis

reset-geth:
	go run tools/reset-exercises/main.go -target geth
