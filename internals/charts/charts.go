package charts

import (
	"fmt"
	"image/color"
	"math"
	"os"
	"sort"

	"github.com/rs/zerolog"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func BuildChartFromTool(tool *ChartCallTool, utils *ChartUtils, logger zerolog.Logger) error {
	var plot *plot.Plot
	var err error
	switch tool.ChartType {
	case BarChart:
		utils.Logger.Info().Msg("Building Bar Chart")
		plot, err = buildBarChart(tool, utils)
		if err != nil {
			utils.Logger.Error().Err(err).Msg("Error building bar chart")
			return err
		}
	case LineChart:
		utils.Logger.Info().Msg("Building Line Chart")
		plot, err = buildLineChart(tool, utils)
		if err != nil {
			utils.Logger.Error().Err(err).Msg("Error building line chart")
			return err
		}
	case Heatmap:
		utils.Logger.Info().Msg("Building Heatmap")
		return buildHeatMap(tool, utils, logger)
	default:
		return fmt.Errorf("unsupported chart type: %s", tool.ChartType)
	}
	if utils.ChartBytes == nil {
		utils.SetFullPath()
		err = plot.Save(font.Length(utils.Width)*vg.Inch, font.Length(utils.Height)*vg.Inch, utils.FullPath)
		if err != nil {
			utils.Logger.Err(err).Msg("Error saving bar chart")
			return err
		}
		utils.ChartBytes, err = os.ReadFile(utils.FullPath)
		if err != nil {
			utils.Logger.Err(err).Msg("Error reading chart bytes")
			return err
		}
	}
	return nil
}

func buildBarChart(tool *ChartCallTool, utils *ChartUtils) (*plot.Plot, error) {
	// If grouping is requested:
	if len(tool.GroupByFields) > 0 && tool.AggregateMethod != "" && tool.AggregateMethod != "none" {
		grouped, err := groupAndAggregate(tool.Dataset, tool.GroupByFields, tool.YAxisField, tool.AggregateMethod, utils)
		if err != nil {
			utils.Logger.Error().Err(err).Msg("Error grouping and aggregating data")
			return nil, err
		}
		// We'll treat "label" as category, "value" as numeric
		return buildPlainBar(grouped, "group", "value", tool.ChartTitle, utils)
	}
	xKey := tool.XAxisField
	if xKey == "" {
		xKey = "x"
	}
	yKey := tool.YAxisField
	if yKey == "" {
		yKey = "y"
	}
	return buildPlainBar(tool.Dataset, xKey, yKey, tool.ChartTitle, utils)
}

func buildPlainBar(
	data []map[string]interface{},
	xKey, yKey, chartTitle string, utils *ChartUtils) (*plot.Plot, error) {

	var categories []string
	vals := make(plotter.Values, 0, len(data))

	for i, row := range data {
		rawX, ok := row[xKey]
		if !ok {
			return nil, fmt.Errorf("missing '%s' in row %d", xKey, i)
		}
		category := fmt.Sprintf("%v", rawX)

		rawY, ok := row[yKey]
		if !ok {
			return nil, fmt.Errorf("missing '%s' in row %d", yKey, i)
		}
		yVal, okF := toFloat(rawY)
		if !okF {
			return nil, fmt.Errorf("value for '%s' in row %d not numeric", yKey, i)
		}

		categories = append(categories, category)
		vals = append(vals, yVal)
	}

	p := plot.New()
	p.Title.Text = chartTitle
	p.Y.Label.TextStyle.Rotation = 45
	p.X.Label.TextStyle.Color = color.RGBA{R: 100, G: 140, B: 240, A: 255}
	p.Y.Label.TextStyle.Color = color.RGBA{R: 100, G: 140, B: 240, A: 255}
	// Use nominal X for string categories
	p.NominalX(categories...)

	bars, err := plotter.NewBarChart(vals, vg.Points(14))
	if err != nil {
		return nil, err
	}
	bars.Color = color.RGBA{R: 0, G: 0, B: 0, A: 255}

	p.Add(bars)
	return p, nil
}

func buildHeatMap(tool *ChartCallTool, utils *ChartUtils, logger zerolog.Logger) error {
	chartBuilder := &HeatmapUtility{
		ChartUtils:        utils,
		ChartInstructions: tool,
	}
	return chartBuilder.Build(logger)
	// heatmapData := make(plotter.GridXYZ, len(heatmapDataset.XValues))
	// for i := range heatmapDataset.XValues {
	// 	heatmapData[i].X = heatmapDataset.XValues[i]
	// 	heatmapData[i].Y = heatmapDataset.YValues[i]
	// }
	// plotter.NewHeatMap(Z, nil, nil)
	// scatter, err := plotter.NewScatter(heatmapData)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// heatmap := plotter.NewHeatMap(Z, nil)
	// p.Add(heatmap)

	// scatter.GlyphStyle.Color = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	// scatter.GlyphStyle.Shape = draw.CircleGlyph{}

	// p.Add(scatter)
	// p.Add(plotter.NewGrid())
	// p.Legend.Add("Heatmap", scatter)
	// utils.SetFullPath()
	// heatmap := charts.NewHeatMap()

	// // Add the data to the heatmap
	// heatmap.AddSeries("heatmap", []opts.HeatMapData{})

	// // Set chart options
	// heatmap.SetGlobalOptions(
	// 	charts.WithTitleOpts(opts.Title{
	// 		Title: tool.ChartTitle,
	// 	}),
	// )
	// file, err := os.Create("heatmap.png")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()

	// // Render the chart as PNG to the file
	// err = heatmap.Render(file)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Heatmap saved as heatmap.png")
}

func buildLineChart(tool *ChartCallTool, utils *ChartUtils) (*plot.Plot, error) {
	// Typically, "line" ignores aggregator & uses (x,y) pairs directly
	rsDataset := []map[string]interface{}{}
	var err error
	if len(tool.GroupByFields) > 0 && tool.AggregateMethod != "" && tool.AggregateMethod != "none" {
		rsDataset, err = groupAndAggregate(tool.Dataset, tool.GroupByFields, tool.YAxisField, tool.AggregateMethod, utils)
		if err != nil {
			return nil, err
		}
		// We'll treat "label" as category, "value" as numeric
		// return buildPlainBar(grouped, "label", "value", tool.ChartTitle)
	} else {
		rsDataset = tool.Dataset
	}
	xKey := tool.XAxisField
	if xKey == "" {
		xKey = "x"
	}
	yKey := tool.YAxisField
	if yKey == "" {
		yKey = "y"
	}

	pts := make(plotter.XYs, 0, len(rsDataset))
	for i, row := range rsDataset {
		rawX, ok := row[xKey]
		if !ok {
			return nil, fmt.Errorf("missing '%s' in row %d", xKey, i)
		}
		xVal, okF := toFloat(rawX)
		if !okF {
			return nil, fmt.Errorf("x value '%s' in row %d not numeric", xKey, i)
		}

		rawY, ok := row[yKey]
		if !ok {
			return nil, fmt.Errorf("missing '%s' in row %d", yKey, i)
		}
		yVal, okF2 := toFloat(rawY)
		if !okF2 {
			return nil, fmt.Errorf("y value '%s' in row %d not numeric", yKey, i)
		}

		pts = append(pts, plotter.XY{X: xVal, Y: yVal})
	}

	// Sort by X to get a proper line
	sort.Slice(pts, func(i, j int) bool { return pts[i].X < pts[j].X })

	p := plot.New()
	p.Title.Text = "Line Chart"
	p.X.Label.Text = xKey
	p.Y.Label.Text = yKey

	line, err := plotter.NewLine(pts)
	if err != nil {
		return nil, err
	}
	line.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255} // red line

	p.Add(line)
	return p, nil
}

func toFloat(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		if val <= math.MaxUint64 {
			return float64(val), true
		}
	}
	// You could add more handling for strings, etc.
	return 0, false
}
