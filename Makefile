run:
	go run .\cmd\tray

build:
	go build -ldflags -H=windowsgui .\cmd\tray