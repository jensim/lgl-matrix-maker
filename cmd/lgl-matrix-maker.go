package main

import "fmt"

struct LaserArea {
  int height,
  int width
}

struct Layer {
    string "id", // : "64ac3715-4f49-4f8b-a1d6-9ca169d7d080",
    string "name", // : "03",
    string "type", // : "Line",
    bool "isOutputEnabled", // : true,
    bool "isAirEnable", // : false,
    int "height", // : 100,
    int "width", // : 100,
    string "color", // : "#00e000",
    int "laserPower", // : 10,
    int "laserSpeed", // : 1000,
    float "engraveQuality", // : 0.091,
    int "engraveCount", // : 1,
    string "mode", // : "Fill",
    int "sort", // : 0,
    string "processMethod", // : "Engrave",
    string "laserType", // : "",
    int "materialId", // : "-1",
    string "material", // : "",
    int "thickness", // : 0,
    bool "isCrossed", // : false,
    bool "isOverscan", // : false,
    float "overScanPer", // : 2.5,
    string "imageMode", // : "Stucki",
    float "scanAngle", // : 0,
    bool "isLocked", // : false,
    int "kerfOffset", // : 0,
    bool "isBidirectional", // : true
}

struct Element {
    string "id", //uuid
    string "type", // : "shape-rect",
    int "angle", //: 0,
    bool "selectable", // : true,
    bool "visible", // : true,
    string "layerId", // : "64ac3715-4f49-4f8b-a1d6-9ca169d7d080",
    bool "isLocked", // : false,
    bool "isOutWorkSpace", // : false,
    string "lineColor", // : "#00e000",
    int "width", // : 5,
    int "height", // : 6,
    int[] transformMatrix, // [1,0,0,1,x,y]
}

struct LglStruct {
    string version, //"1.2.1",
    LaserArea LaserArea,
    Layer[] Layers,
    Element[] Elements,
}

func main() {
    fmt.Println("Hello, World!")
}
