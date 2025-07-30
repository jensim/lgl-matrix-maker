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
	Height uint16 `json:"height"`
	Width  uint16 `json:"width"`
}

type Layer struct {
	Id              string  `json:"id"`
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	IsOutputEnabled bool    `json:"isOutputEnabled"`
	IsAirEnable     bool    `json:"isAirEnable"`
	Height          uint16  `json:"height"`
	Width           uint16  `json:"width"`
	Color           string  `json:"color"`
	LaserPower      uint16  `json:"laserPower"`
	LaserSpeed      uint32  `json:"laserSpeed"`
	EngraveQuality  float32 `json:"engraveQuality"`
	EngraveCount    uint16  `json:"engraveCount"`
	Mode            string  `json:"mode"`
	Sort            uint16  `json:"sort"`
	ProcessMethod   string  `json:"processMethod"`
	MaterialId      string  `json:"materialId"`
	IsCrossed       bool    `json:"isCrossed"`
	IsLocked        bool    `json:"isLocked"`
	IsBidirectional bool    `json:"isBidirectional"`
}

type Element struct {
	Id              string     `json:"id"`
	Type            string     `json:"type"`
	Selectable      bool       `json:"selectable"`
	Visible         bool       `json:"visible"`
	LayerId         string     `json:"layerId"`
	IsLocked        bool       `json:"isLocked"`
	IsOutWorkSpace  bool       `json:"isOutWorkSpace"`
	LineColor       string     `json:"lineColor"`
	Width           uint16     `json:"width"`
	Height          uint16     `json:"height"`
	TransformMatrix [6]float32 `json:"transformMatrix"`
}

type LglStruct struct {
	Version   string    `json:"version"`
	LaserArea LaserArea `json:"laserArea"`
	Layers    []Layer   `json:"layers"`
	Elements  []Element `json:"elements"`
}

func main() {
	var x int
	flag.IntVar(&x, "x", 20, "number of permutations in x, power")
	var y int
	flag.IntVar(&y, "y", 20, "number of permutations in y, speed")
	var powerIncrements int
	flag.IntVar(&powerIncrements, "pi", 5, "increments in power per tile, measured in %")
	var powerMin int
	flag.IntVar(&powerMin, "pm", 5, "first tile, minimum power")
	var speedIncrements int
	flag.IntVar(&speedIncrements, "si", 200, "increments in speed per tile, measured in mm/s")
	var speedMin int
	flag.IntVar(&speedMin, "sm", 3000, "first tile, minimum speed")
	var quality float64
	flag.Float64Var(&quality, "q", 0.01, "EngraveQuality")

	var width float64
	flag.Float64Var(&width, "w", 2.0, "tile width in mm")
	var space float64
	flag.Float64Var(&space, "s", 0.5, "space inbetween tiles in mm")
	var areaX int
	flag.IntVar(&areaX, "ax", 100, "work area x width in mm")
	var areaY int
	flag.IntVar(&areaY, "ay", 100, "work area y height in mm")
	var mode string
	flag.StringVar(&mode, "mode", "Fill", "\"Fill\" or \"Line\"")
	var processMethod string
	flag.StringVar(&processMethod, "method", "Engrave", "\"Cut\", or \"Engrave\"")

	flag.Parse()

	if mode == "Fill" && processMethod == "Cut" {
		mode = "Line"
		log.Println("Cut only supports Line mode, mode set to Line")
	}

	//var permutations = x * y
	var elements = make([]Element, 0)
	var layers = make([]Layer, 0)

	var order = 0
	for xidx := 0; xidx < x; xidx++ {
		for yidx := 0; yidx < y; yidx++ {
			xpos := (width / 2) + (space * float64(xidx+1)) + (width * float64(xidx))
			ypos := (width / 2) + (space * float64(yidx+1)) + (float64(yidx) * width)
			power := (xidx * powerIncrements) + powerMin
			yIdxInverted := (y - 1 - yidx)
			speedMmPerSec := (yIdxInverted * speedIncrements) + speedMin
			speedMmPerMin := speedMmPerSec * 60
			order += order

			layerId := uuid.NewString()
			elements = append(elements, Element{
				Id:              uuid.NewString(),
				LayerId:         layerId,
				Type:            "shape-rect",
				Selectable:      true,
				Visible:         true,
				IsLocked:        false,
				IsOutWorkSpace:  false,
				LineColor:       "#00e000",
				Width:           uint16(width),
				Height:          uint16(width),
				TransformMatrix: [6]float32{1, 0, 0, 1, float32(xpos), float32(ypos)},
			})
			layers = append(layers, Layer{
				Id:              layerId,
				Name:            fmt.Sprintf("P:%d,S:%d", power, speedMmPerSec),
				Type:            "Fill",
				IsOutputEnabled: true,
				IsAirEnable:     false,
				Height:          uint16(width),
				Width:           uint16(width),
				Color:           "#00e000",
				LaserPower:      uint16(power),
				LaserSpeed:      uint32(speedMmPerMin),
				EngraveQuality:  float32(quality),
				EngraveCount:    1,
				Mode:            mode,
				Sort:            uint16(order),
				ProcessMethod:   processMethod,
				MaterialId:      "-1",
				IsCrossed:       false,
				IsLocked:        false,
				IsBidirectional: true,
			})
		}
	}
	allOfIt := LglStruct{
		Version: "1.2.1",
		LaserArea: LaserArea{
			Height: uint16(areaY),
			Width:  uint16(areaX),
		},
		Elements: elements,
		Layers:   layers,
	}
	b, err := json.MarshalIndent(allOfIt, "", "  ")
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	fileErr := os.WriteFile("matrix.lgl", []byte(b), 0740)
	if fileErr != nil {
		log.Fatal(fileErr)
	}
}
