package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/tdewolff/canvas"
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
	LaserSpeed      float64 `json:"laserSpeed"`
	LaserType       string  `json:"laserType"`
	EngraveQuality  float32 `json:"engraveQuality"`
	EngraveCount    uint16  `json:"engraveCount"`
	Mode            string  `json:"mode"`
	Sort            uint16  `json:"sort"`
	ProcessMethod   string  `json:"processMethod"`
	MaterialId      string  `json:"materialId"`
	IsCrossed       bool    `json:"isCrossed"`
	IsLocked        bool    `json:"isLocked"`
	IsBidirectional bool    `json:"isBidirectional"`
	IsOverscan      bool    `json:"isOverscan"`
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
	Angle           float32    `json:"angle,omitempty"`      // Optional field, may not be present
	Path            string     `json:"path,omitempty"`       // Optional field, may not be present
	Text            string     `json:"text,omitempty"`       // Optional field, may not be present
	FontFamily      string     `json:"fontFamily,omitempty"` // Optional field, may not be present
	FontSize        float32    `json:"fontSize,omitempty"`   // Optional field, may not be present
	FontWeight      bool       `json:"fontWeight,omitempty"` // Optional field, may not be present
	FontStyle       bool       `json:"fontStyle,omitempty"`  // Optional field, may not be present
	RowSpace        string     `json:"rowSpace,omitempty"`   // Optional field, may not be present
	WordSpace       string     `json:"wordSpace,omitempty"`  // Optional field, may not be present
	TextAlign       string     `json:"textAlign,omitempty"`  // Optional field, may not be present
	Curve           float32    `json:"curve,omitempty"`      // Optional field, may not be present
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
	var speedIncrements float64
	flag.Float64Var(&speedIncrements, "si", 200, "increments in speed per tile, measured in mm/s")
	var speedMin float64
	flag.Float64Var(&speedMin, "sm", 3000, "first tile, minimum speed")
	var quality float64
	flag.Float64Var(&quality, "q", 0.01, "EngraveQuality")
	var outputFileName string
	flag.StringVar(&outputFileName, "o", "matrix.lgl", "output file name")

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
	var laserType string
	flag.StringVar(&laserType, "laserType", "BlueLight", "Laser type, e.g., \"BlueLight\", \"RedLight\", etc.")

	flag.Parse()

	if mode == "Fill" && processMethod == "Cut" {
		mode = "Line"
		log.Println("Cut only supports Line mode, mode set to Line")
	}
	if laserType != "BlueLight" && laserType != "RedLight" {
		log.Println("Unknown laser type used!")
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
			speedMmPerSec := (float64(yIdxInverted) * speedIncrements) + speedMin
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
			speedStr := strconv.FormatFloat(speedMmPerSec, 'f', 1, 64)
			layers = append(layers, Layer{
				Id:              layerId,
				Name:            fmt.Sprintf("P:%d,S:%s", power, speedStr),
				Type:            "Fill",
				IsOutputEnabled: true,
				IsAirEnable:     false,
				IsCrossed:       false,
				IsOverscan:      false,
				Height:          uint16(width),
				Width:           uint16(width),
				Color:           "#00e000",
				LaserPower:      uint16(power),
				LaserSpeed:      speedMmPerMin,
				LaserType:       laserType,
				EngraveQuality:  float32(quality),
				EngraveCount:    1,
				Mode:            mode,
				Sort:            uint16(order),
				ProcessMethod:   processMethod,
				MaterialId:      "-1",
				IsLocked:        false,
				IsBidirectional: true,
			})
		}
	}
	textLayer := layers[len(layers)/2]
	textLayer.Name = "Text Layer"
	textLayer.Mode = "Line"
	textLayer.ProcessMethod = "Engrave"
	textLayer.Id = uuid.NewString()
	textLayer.Color = "#000000"
	layers = append(layers, textLayer)

	powerElement := Element{
		Id:              uuid.NewString(),
		Type:            "text",
		Angle:           0,
		Selectable:      true,
		Visible:         true,
		LayerId:         textLayer.Id,
		IsLocked:        false,
		IsOutWorkSpace:  true,
		LineColor:       "#000000",
		Path:            textToPath("POWER%", 8),
		TransformMatrix: [6]float32{1, 0, 0, -1, 25, -4},
		Text:            "POWER%",
		FontFamily:      "Arial",
		FontSize:        8.69,
		RowSpace:        "1.0",
		WordSpace:       "0.0",
		TextAlign:       "center",
		Curve:           0,
	}
	elements = append(elements, powerElement)

	speedElement := powerElement
	speedElement.Id = uuid.NewString()
	speedElement.Angle = -90
	speedElement.Path = textToPath("SPEED (mm/s)", 8)
	speedElement.TransformMatrix = [6]float32{1, 0, 0, -1, -4, 25}
	speedElement.Text = "SPEED (mm/s)"
	elements = append(elements, speedElement)

	for xidx := 0; xidx < x; xidx++ {
		power := (xidx * powerIncrements) + powerMin
		xpos := (width / 2) + (space * float64(xidx+1)) + (width * float64(xidx))
		newElement := powerElement
		newElement.Id = uuid.NewString()
		newElement.Text = fmt.Sprintf("%d%%", power)
		newElement.Path = textToPath(newElement.Text, 2)
		newElement.TransformMatrix = [6]float32{1, 0, 0, -1, float32(xpos), -1.3}
		elements = append(elements, newElement)
	}
	for yidx := 0; yidx < y; yidx++ {
		ypos := (width / 2) + (space * float64(yidx+1)) + (float64(yidx) * width)
		yIdxInverted := (y - 1 - yidx)
		speedMmPerSec := (float64(yIdxInverted) * speedIncrements) + speedMin
		speedStr := strconv.FormatFloat(speedMmPerSec, 'f', 1, 64)

		newElement := powerElement
		newElement.Id = uuid.NewString()
		newElement.Text = speedStr
		newElement.Path = textToPath(newElement.Text, 2)
		newElement.TransformMatrix = [6]float32{1, 0, 0, -1, -1.3, float32(ypos)}
		elements = append(elements, newElement)
	}

	allOfIt := LglStruct{
		Version: "2.0.0",
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

	fileErr := os.WriteFile(outputFileName, []byte(b), 0740)
	if fileErr != nil {
		log.Fatal(fileErr)
	}
}

func textToPath(text string, fontSize float64) string {
	fontFamily := canvas.NewFontFamily("Arial")
	if err := fontFamily.LoadSystemFont("Arial Regular, sans-serif", canvas.FontRegular); err != nil {
		panic(err)
	}
	face := fontFamily.Face(fontSize, canvas.Black, canvas.FontBold, canvas.FontNormal)
	path, _, err := face.ToPath(text)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s", path)
}
