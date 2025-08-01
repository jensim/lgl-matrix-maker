run:
	go run ./cmd/lgl-matrix-maker.go
build:
	go build -o bin/lgl-matrix-maker -buildvcs=false -v ./...
generateIR:
	go run ./cmd/lgl-matrix-maker.go -sm 0.5 -si 0.5 -laserType "RedLight" -o 'matrix ir fill 0-10.lgl'
generateDiode:
	go run ./cmd/lgl-matrix-maker.go -y 36 -o 'matrix blue fill 3000-10000.lgl'
