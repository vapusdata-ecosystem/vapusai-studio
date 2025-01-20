package charts

import (
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
)

const (
	PieChart      = "pie"
	LineChart     = "line"
	BarChart      = "bar"
	DoughnutChart = "doughnut"
	Heatmap       = "heatmap"
)

type ChartCallTool struct {
	ChartType       string                   `json:"chartType"`
	GroupByFields   []string                 `json:"groupFields"`
	XAxisField      string                   `json:"xAxisField"`
	YAxisField      string                   `json:"yAxisField"`
	AggregateMethod string                   `json:"aggregateMethod"`
	Dataset         []map[string]interface{} `json:"dataset"`
	ChartTitle      string                   `json:"chartTitle"`
}

type ChartUtils struct {
	Filename   string
	FileType   string
	Path       string
	Height     int
	Width      int
	Logger     zerolog.Logger
	FullPath   string
	ChartBytes []byte
}

func (utils *ChartUtils) SetFullPath() {
	if utils.FileType == "" {
		utils.Filename = utils.Filename + ".png"
	} else {
		utils.Filename = utils.Filename + "." + strings.ToLower(utils.FileType)
	}
	utils.FullPath = filepath.Join(utils.Path, utils.Filename)
}

type Slices struct {
	Label string
	Value float64
}

type HeatmapDataset struct {
	XValues []float64
	YValues []float64
	ZValues []float64
	aggData map[interface{}]map[interface{}]float64
	xVals   []interface{}
	yVals   []interface{}
	hMap    [][]int
}
