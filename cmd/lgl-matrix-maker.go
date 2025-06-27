package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
)

type LaserArea struct {
	Height int16 `json:"height"`
	Width  int16 `json:"width"`
}

type Layer struct {
	Id              string  `json:"id"`
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	IsOutputEnabled bool    `json:"isOutputEnabled"`
	IsAirEnable     bool    `json:"isAirEnable"`
	Height          int16   `json:"height"`
	Width           int16   `json:"width"`
	Color           string  `json:"color"`
	LaserPower      int16   `json:"laserPower"`
	LaserSpeed      int16   `json:"laserSpeed"`
	EngraveQuality  float32 `json:"engraveQuality"`
	EngraveCount    int16   `json:"engraveCount"`
	Mode            string  `json:"mode"`
	Sort            int16   `json:"sort"`
	ProcessMethod   string  `json:"processMethod"`
	LaserType       string  `json:"laserType"`
	MaterialId      string  `json:"materialId"`
	Material        string  `json:"material"`
	Thickness       int16   `json:"thickness"`
	IsCrossed       bool    `json:"isCrossed"`
	IsOverscan      bool    `json:"isOverscan"`
	OverScanPer     float32 `json:"overScanPer"`
	ImageMode       string  `json:"imageMode"`
	ScanAngle       float32 `json:"scanAngle"`
	IsLocked        bool    `json:"isLocked"`
	KerfOffset      int16   `json:"kerfOffset"`
	IsBidirectional bool    `json:"isBidirectional"`
}

type Element struct {
	Id              string   `json:"id"`
	Type            string   `json:"type"`
	Angle           int16    `json:"angle"`
	Selectable      bool     `json:"selectable"`
	Visible         bool     `json:"visible"`
	LayerId         string   `json:"layerId"`
	IsLocked        bool     `json:"isLocked"`
	IsOutWorkSpace  bool     `json:"isOutWorkSpace"`
	LineColor       string   `json:"lineColor"`
	Width           int16    `json:"width"`
	Height          int16    `json:"height"`
	TransformMatrix [6]int16 `json:"transformMatrix"`
}

type LglStruct struct {
	Version   string `json:"version"`
	LaserArea LaserArea `json:"laserArea"`
	Layers    []Layer `json:"layers"`
	Elements  []Element `json:"element"`
}

func main() {
	var x int
	flag.IntVar(&x, "x", 10, "number of permutations in x, power")
	var y int
	flag.IntVar(&y, "y", 10, "number of permutations in y, speed")
	var powerIncrements int
	flag.IntVar(&powerIncrements, "pi", 10, "increments in power per tile, measured in %")
	var powerMin int
	flag.IntVar(&powerMin, "pm", 10, "first tile, minimum power")
	var speedIncrements int
	flag.IntVar(&speedIncrements, "si", 1000, "increments in speed per tile, measured in mm/s")
	var speedMin int
	flag.IntVar(&speedMin, "sm", 1000, "first tile, minimum speed")
	var quality float64
	 flag.Float64Var(&quality, "q", 0.091, "EngraveQuality")

	var width int
	flag.IntVar(&width, "w", 5, "tile width in mm")
	var space int
	flag.IntVar(&space, "s", 2, "space inbetween tiles in mm")
	var areaX int
	flag.IntVar(&areaX, "ax", 100, "work area x width in mm")
	var areaY int
	flag.IntVar(&areaY, "ay", 100, "work area y height in mm")

	flag.Parse()

	var permutations = x * y
	var elements = make([]Element, permutations)
	var layers = make([]Layer, permutations)

	var order = 0
	for xidx := 0; xidx < x; xidx++ {
		for yidx := 0; yidx < y; yidx++ {
			xpos := int16((space * (xidx + 1)) + (xidx * width))
			ypos := int16((space * (yidx + 1)) + (yidx * width))
			power := (xidx * powerIncrements) + powerMin
			yIdxInverted := (y-1-yidx)
			speed := (yIdxInverted * speedIncrements) + speedMin
			order += order

			layerId := uuid.NewString()
			elements = append(elements, Element{
				Id:              uuid.NewString(),
				LayerId:         layerId,
				Type:            "shape-rect",
				Angle:           0,
				Selectable:      true,
				Visible:         true,
				IsLocked:        false,
				IsOutWorkSpace:  false,
				LineColor:       "#00e000",
				Width:           int16(width),
				Height:          int16(width),
				TransformMatrix: [6]int16{1, 0, 0, 1, xpos, ypos},
			})
			layers = append(layers, Layer{
				Id:              uuid.NewString(),
				Name:            fmt.Sprintf("P:%d,S:%d", power, speed),
				Type:            "Line",
				IsOutputEnabled: true,
				IsAirEnable:     false,
				Height:          int16(width),
				Width:           int16(width),
				Color:           "#00e000",
				LaserPower:      int16(power),
				LaserSpeed:      int16(speed),
				EngraveQuality:  float32(quality),
				EngraveCount:    1,
				Mode:            "Fill",
				Sort:            int16(order),
				ProcessMethod:   "Engrave",
				LaserType:       "",
				MaterialId:      "-1",
				Material:        "",
				Thickness:       0,
				IsCrossed:       false,
				IsOverscan:      false,
				OverScanPer:     2.5,
				ImageMode:       "Stucki",
				ScanAngle:       0,
				IsLocked:        false,
				KerfOffset:      0,
				IsBidirectional: true,
			})
		}
	}
	allOfIt := LglStruct{
		Version: "1.2.1",
		LaserArea: LaserArea{
			Height: int16(areaY),
			Width: int16(areaX),
		},
		Elements: elements,
		Layers: layers,
	}
	b, err := json.Marshal(allOfIt)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return;
	}
	//fmt.Println(string(b))

	fileErr := os.WriteFile("matrix.lgl", []byte(b), 0740)
	if fileErr != nil {
		log.Fatal(fileErr)
	}
}
