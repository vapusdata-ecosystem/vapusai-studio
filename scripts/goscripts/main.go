package main

import (
	"log"

	vapuspublish "github.com/vapusdata-ecosystem/vapusdata/scripts/goscripts/publish"
)

func main() {
	chartVersion := vapuspublish.HelmChartOps()
	log.Println(chartVersion)
}
