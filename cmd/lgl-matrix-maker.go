package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

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

	elements = append(elements, Element{
		Id:             uuid.NewString(),
		Type:           "text",
		Angle:          -90,
		Selectable:     true,
		Visible:        true,
		LayerId:        textLayer.Id,
		IsLocked:       false,
		IsOutWorkSpace: true,
		LineColor:      "#000000",
		Path:           "M 60.013 51.543 C 48.618 54.736 42.274 61.286 40.551 71.637 C 39.588 77.423 41.469 83.406 45.736 88.129 C 50.705 93.63 54.315 95.256 72.038 99.978 C 80.466 102.223 88.906 105.015 90.795 106.182 C 97.349 110.233 98.571 118.791 93.464 124.86 C 89.78 129.238 83.829 131.242 74.418 131.273 C 60.741 131.319 51.831 125.69 49.054 115.25 C 47.846 110.71 47.444 110.535 40.75 111.642 C 36.451 112.352 36.256 113.133 38.495 120.672 C 40.359 126.948 46.666 134.29 52.91 137.454 C 67.426 144.81 90.149 142.77 100.376 133.193 C 107.55 126.476 110.058 113.882 105.922 105.339 C 101.849 96.926 95.852 93.608 75.41 88.459 C 55.578 83.465 52 81.282 52 74.182 C 52 65.568 59.568 60.568 72.479 60.651 C 84.678 60.729 91.558 65.157 93.546 74.21 C 94.204 77.205 95.359 77.498 102.119 76.379 L 105.738 75.781 L 104.844 71.802 C 101.763 58.09 91.784 51.016 74.587 50.355 C 68.48 50.12 63.701 50.509 60.013 51.543 M 124 95.91 L 124 141 L 130 141 L 136 141 L 136 122.5 L 136 104 L 152.818 104 C 165.128 104 171.042 103.583 174.878 102.444 C 185.605 99.258 192.176 89.575 192.176 76.953 C 192.176 65.782 187.452 58.006 178 53.617 C 173.897 51.712 171.317 51.496 148.75 51.174 L 124 50.821 L 124 95.91 M 208 96 L 208 141 L 241.5 141 L 275 141 L 275 136 L 275 131 L 247.5 131 L 220 131 L 220 115 L 220 99 L 244.5 99 L 269 99 L 269 94 L 269 89 L 244.5 89 L 220 89 L 220 75 L 220 61 L 246.5 61 L 273 61 L 273 56 L 273 51 L 240.5 51 L 208 51 L 208 96 M 291 96 L 291 141 L 324.5 141 L 358 141 L 358 136 L 358 131 L 330.5 131 L 303 131 L 303 115 L 303 99 L 328 99 L 353 99 L 353 94 L 353 89 L 328 89 L 303 89 L 303 75 L 303 61 L 329.5 61 L 356 61 L 356 56 L 356 51 L 323.5 51 L 291 51 L 291 96 M 374 95.919 L 374 141 L 393.951 141 C 404.924 141 416.286 140.494 419.201 139.876 C 438.222 135.84 448.349 120.183 448.261 94.948 C 448.205 79.243 443.857 67.653 435.05 59.734 C 426.802 52.319 424.024 51.695 397.25 51.237 L 374 50.839 L 374 95.919 M 136 77.362 L 136 94 L 150.25 93.994 C 165.653 93.988 171.191 93.011 174.845 89.654 C 180.386 84.564 181.691 75.781 177.906 69.046 C 174.23 62.506 171.838 61.758 152.838 61.21 L 136 60.723 L 136 77.362 M 386.513 62.106 C 386.216 62.881 386.091 78.586 386.236 97.007 L 386.5 130.5 L 400.5 130.419 C 421.271 130.299 426.827 127.805 432.369 116.109 C 435.413 109.683 435.5 109.112 435.5 95.5 C 435.5 82.19 435.362 81.219 432.699 75.796 C 429.692 69.671 425.548 65.555 420.087 63.267 C 415.035 61.15 387.248 60.19 386.513 62.106",
		TransformMatrix: [6]float32{
			0.019361970033389339,
			0,
			0,
			0.021734919178404595,
			-4,
			25,
		},
		Text:       "SPEED",
		FontFamily: "Arial",
		FontSize:   8.69,
		FontWeight: false,
		FontStyle:  false,
		RowSpace:   "1.0",
		WordSpace:  "0.0",
		TextAlign:  "center",
		Curve:      0,
	})
	elements = append(elements, Element{
		Id:             uuid.NewString(),
		Type:           "text",
		Angle:          0,
		Selectable:     true,
		Visible:        true,
		LayerId:        textLayer.Id,
		IsLocked:       false,
		IsOutWorkSpace: true,
		LineColor:      "#000000",
		Path:           "M 151.024 51.603 C 132.334 56.626 120.935 73.275 120.768 95.794 C 120.664 109.936 124.502 120.421 133.119 129.533 C 148.38 145.668 176.23 146.089 192.963 130.438 C 206.54 117.738 210.065 92.008 200.898 72.506 C 192.803 55.284 171.129 46.2 151.024 51.603 M 41 95.922 L 41 141 L 47 141 L 53 141 L 53 122.5 L 53 104 L 68.951 104 C 87.341 104 93.328 102.994 98.882 98.97 C 105.155 94.426 108.262 88.186 108.79 79.073 C 109.331 69.727 106.839 62.895 100.951 57.575 C 95.052 52.245 92.039 51.712 65.25 51.256 L 41 50.844 L 41 95.922 M 225.147 95.75 C 231.651 120.362 236.98 140.634 236.987 140.798 C 236.994 140.961 239.786 140.961 243.19 140.798 C 249.125 140.512 249.414 140.376 250.192 137.5 C 250.639 135.85 255.456 118.314 260.898 98.53 L 270.792 62.56 L 272.236 67.53 C 273.03 70.264 277.89 87.912 283.036 106.75 L 292.392 141 L 298.128 141 C 302.73 141 303.967 140.654 304.377 139.25 C 305.392 135.782 326.895 56.435 327.548 53.75 L 328.217 51 L 322.287 51 L 316.357 51 L 313.701 62.25 C 304.413 101.585 298.12 126.296 297.721 125 C 297.466 124.175 296.3 119 295.13 113.5 C 293.959 108 289.718 91.8 285.706 77.5 L 278.411 51.5 L 271.26 51.207 L 264.11 50.915 L 262.112 57.707 C 257.394 73.747 245.059 118.815 244.446 122.25 C 243.832 125.696 241.953 127.644 242.036 124.75 C 242.055 124.063 238.366 107.188 233.837 87.25 L 225.603 51 L 219.461 51 L 213.32 51 L 225.147 95.75 M 340 96 L 340 141 L 373.5 141 L 407 141 L 407 136 L 407 131 L 379.5 131 L 352 131 L 352 115 L 352 99 L 376.5 99 L 401 99 L 401 94 L 401 89 L 376.5 89 L 352 89 L 352 75 L 352 61 L 378 61 L 404 61 L 404 56 L 404 51 L 372 51 L 340 51 L 340 96 M 423 95.927 L 423 141 L 429 141 L 435 141 L 435 121 L 435 101 L 444.75 101.022 C 460.945 101.058 462.253 102.109 479.384 128.819 L 487.196 141 L 494.629 141 L 502.061 141 L 498.704 135.75 C 496.857 132.863 493.376 127.35 490.968 123.5 C 483.767 111.985 478.741 105.798 474.294 102.974 L 470.087 100.304 L 475.81 98.713 C 490.274 94.691 497.568 82.344 493.581 68.631 C 491.601 61.82 487.06 56.389 481 53.583 C 476.86 51.667 474.358 51.474 449.75 51.177 L 423 50.855 L 423 95.927 M 153.332 61.851 C 143.26 65.332 137.283 71.888 134.544 82.46 C 132.55 90.158 132.571 102.84 134.59 110.037 C 139.188 126.422 155.795 135.399 172.517 130.538 C 180.065 128.344 188.803 119.984 191.337 112.533 C 195.154 101.31 194.361 85.593 189.494 75.988 C 182.892 62.96 167.517 56.949 153.332 61.851 M 53 77.448 L 53 94.172 L 69.25 93.836 C 87.062 93.468 89.945 92.647 94.395 86.677 C 97.358 82.702 97.382 72.328 94.438 68.385 C 89.731 62.08 88.614 61.752 70.067 61.216 L 53 60.723 L 53 77.448 M 435 76 L 435 91 L 451.25 90.956 C 463.868 90.922 468.612 90.531 472.477 89.206 C 479.141 86.922 482 82.732 482 75.252 C 482 70.418 481.626 69.426 478.703 66.503 C 474.072 61.872 469.777 61.012 451.25 61.006 L 435 61 L 435 76",
		TransformMatrix: [6]float32{
			0.019361970033389339,
			0,
			0,
			0.021734919178404595,
			25,
			-4,
		},
		Text:       "POWER",
		FontFamily: "Arial",
		FontSize:   8.69,
		RowSpace:   "1.0",
		WordSpace:  "0.0",
		TextAlign:  "center",
		Curve:      0,
	})

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
