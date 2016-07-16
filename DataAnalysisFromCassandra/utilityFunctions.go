package DataAnalysisFromCassandra

import (
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"image/color"
	"math/rand"
	"os"
	"strconv"
)

func PlotToFile(name string, data ...[]float64) error {
	file, err := os.Create("/tmp/" + name + ".jpg")
	defer file.Close()
	if err != nil {
		return err
	}
	dataXYs := make([]plotter.XYs, 0, len(data))

	for i, tab := range data {
		dataXYs = append(dataXYs, make(plotter.XYs, len(tab)))
		for j := 0; j < len(dataXYs[i]); j++ {
			dataXYs[i][j].X = float64(j)
			dataXYs[i][j].Y = data[i][j]
		}
	}

	p, err := plot.New()
	if err != nil {
		return err
	}

	p.Title.Text = name
	p.X.Label.Text = "t"
	p.Y.Label.Text = "Data"
	p.Add(plotter.NewGrid())

	for i, tab := range dataXYs {
		line, err := plotter.NewLine(tab)
		if err != nil {
			return err
		}

		line.LineStyle.Width = vg.Points(1)
		line.LineStyle.Color = color.RGBA{R: uint8(rand.Int() % 200), B: uint8(rand.Int() % 200), A: 255}

		p.Add(line)
		p.Legend.Add(strconv.Itoa(i), line)
	}

	wt, err := p.WriterTo(vg.Inch*16, vg.Inch*16, "jpg")
	if err != nil {
		return err
	}

	_, err = wt.WriteTo(file)

	return err
}
