
default:
		GOOS=linux CGO_ENABLED=0 go build -o dist/replay

convey:
	goconvey --port 8890

test:
	go test ../... -v