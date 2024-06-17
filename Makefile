run:
	go run .\cmd\tray CGO_ENABLED=1

build:
	go build -ldflags -H=windowsgui .\cmd\tray 