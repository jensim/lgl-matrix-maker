# lgl-matrix-maker

### List command
```shell
go run cmd/lgl-matrix-maker.go -h
```
```
Usage of lgl-matrix-maker:
  -ax int
        work area x width in mm (default 100)
  -ay int
        work area y height in mm (default 100)
  -pi int
        increments in power per tile, measured in % (default 10)
  -pm int
        first tile, minimum power (default 10)
  -q float
        EngraveQuality (default 0.091)
  -s float
        space inbetween tiles in mm (default 2)
  -si int
        increments in speed per tile, measured in mm/s (default 1000)
  -sm int
        first tile, minimum speed (default 1000)
  -w float
        tile width in mm (default 5)
  -x int
        number of permutations in x, power (default 10)
  -y int
        number of permutations in y, speed (default 10)
```


### Run with all defaults
```shell
go run cmd/lgl-matrix-maker.go
```

