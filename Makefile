run:
	go run .\cmd\tray CGO_ENABLED=1

build:
	go build -ldflags -H=windowsgui CGO_ENABLED=1 .\cmd\tray 