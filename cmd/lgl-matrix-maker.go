package main

import "fmt"

type LaserArea struct {
  Height int16
  Width int16
}

type Layer struct {
    Id string // : "64ac3715-4f49-4f8b-a1d6-9ca169d7d080",
    Name string // : "03",
    Type string // : "Line",
    IsOutputEnabled bool // : true,
    IsAirEnable bool // : false,
    Height int16 // : 100,
    Width int16 // : 100,
    Color string // : "#00e000",
    LaserPower int16 // : 10,
    LaserSpeed int16 // : 1000,
    EngraveQuality float32 // : 0.091,
    EngraveCount int16 // : 1,
    Mode string // : "Fill",
    Sort int16 // : 0,
    ProcessMethod string // : "Engrave",
    LaserType string // : "",
    MaterialId int16 // : "-1",
    Material string // : "",
    Thickness int16 // : 0,
    IsCrossed bool // : false,
    IsOverscan bool // : false,
    OverScanPer float32 // : 2.5,
    ImageMode string // : "Stucki",
    ScanAngle float32 // : 0,
    IsLocked bool // : false,
    KerfOffset int16 // : 0,
    IsBidirectional bool // : true
}

type Element struct {
    Id string //uuid
    Type string // : "shape-rect",
    Angle int16 //: 0,
    Selectable bool // : true,
    Visible bool // : true,
    LayerId string // : "64ac3715-4f49-4f8b-a1d6-9ca169d7d080",
    IsLocked bool // : false,
    IsOutWorkSpace bool // : false,
    LineColor string // : "#00e000",
    Width int16 // : 5,
    Height int16 // : 6,
    TransformMatrix []int16 // [1,0,0,1,x,y]
}

type LglStruct struct {
    Version string //"1.2.1",
    LaserArea LaserArea
    Layers []Layer
    Elements []Element
}

func main() {
    fmt.Println("Hello, World!")
}
