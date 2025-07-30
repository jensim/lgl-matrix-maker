
build:
	go build -o bin/lgl-matrix-maker -buildvcs=false -v ./...
generateIR:
	go run ./cmd/lgl-matrix-maker.go -sm 0.5 -si 0.5 -laserType "RedLight"
generateDiode:
	go run ./cmd/lgl-matrix-maker.go -y 36 
