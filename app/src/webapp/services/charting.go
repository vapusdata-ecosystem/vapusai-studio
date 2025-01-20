package services

// import (
// 	"fmt"
// 	"image/color"
// 	"log"
// 	"os"

// 	"github.com/chenjiandongx/go-echarts/charts"
// )

// func mains() {
// 	// Example data
// 	categories := []string{"A", "B", "C", "D", "E"}
// 	values := []int{10, 20, 15, 30, 25}

// 	// Create a new bar chart
// 	bar := charts.NewBar()

// 	// Set the title
// 	bar.SetGlobalOptions(
// 		charts.WithTitleOpts(charts.Title{Title: "Bar Chart Example"}),
// 		charts.WithXAxisOpts(charts.XAxis{Data: categories}),
// 		charts.WithYAxisOpts(charts.YAxis{SplitLine: charts.SplitLine{Show: true}}),
// 	)

// 	// Add data to the chart
// 	bar.AddSeries("Category", generateBarItems(values))

// 	// Render the chart to a file
// 	f, err := os.Create("bar_chart.html")
// 	if err != nil {
// 		fmt.Println("Error creating file:", err)
// 		return
// 	}
// 	defer f.Close()

// 	// Render the chart to the file as HTML
// 	err = bar.Render(f)
// 	if err != nil {
// 		fmt.Println("Error rendering chart:", err)
// 	}
// 	fmt.Println("Bar chart saved as 'bar_chart.html'")
// }

// // Helper function to create bar chart items
// func generateBarItems(values []int) []charts.BarData {
// 	var items []charts.BarData
// 	for _, value := range values {
// 		items = append(items, charts.BarData{Value: value})
// 	}
// 	return items
// }

// func maind() {
// 	// Example data
// 	values := []float64{10, 20, 15, 30, 25}
// 	names := []string{"A", "B", "C", "D", "E"}

// 	// Create a new plot
// 	p, err := plot.New()
// 	if err != nil {
// 		log.Fatalf("could not create plot: %v", err)
// 	}

// 	// Create a bar chart
// 	barWidth := vg.Points(20)
// 	bars := make(plotter.Values, len(values))
// 	for i, v := range values {
// 		bars[i] = v
// 	}

// 	barChart, err := plotter.NewBarChart(bars, barWidth)
// 	if err != nil {
// 		log.Fatalf("could not create bar chart: %v", err)
// 	}

// 	barChart.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}

// 	// Set the title and labels
// 	p.Title.Text = "Bar Chart Example"
// 	p.X.Label.Text = "Categories"
// 	p.Y.Label.Text = "Values"

// 	// Set custom X axis labels
// 	p.X.Tick.Labels = names

// 	// Add the bar chart to the plot
// 	p.Add(barChart)

// 	// Save the plot as a PNG file
// 	file, err := os.Create("bar_chart.png")
// 	if err != nil {
// 		log.Fatalf("could not create file: %v", err)
// 	}
// 	defer file.Close()

// 	// Save the plot as an image
// 	err = p.Save(6*vg.Inch, 4*vg.Inch, file)
// 	if err != nil {
// 		log.Fatalf("could not save plot: %v", err)
// 	}

// 	fmt.Println("Bar chart saved as 'bar_chart.png'")
// }
