package main

import (
	"fmt"
	"image/color"
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// Function to draw uniform distribution graph
func drawUniformGraph(a, b float64) (*plot.Plot, error) {
	uniform := distuv.Uniform{Min: a, Max: b}

	p := plot.New()
	p.Title.Text = "Uniform Distribution"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "f(x), F(x)"

	n := 1000
	xValues := make(plotter.XYs, n)
	fValues := make(plotter.XYs, n)
	FValues := make(plotter.XYs, n)

	for i := 0; i < n; i++ {
		x := a - 5 + (float64(i)/float64(n))*(b+5-(a-5))
		xValues[i].X = x
		fValues[i].X = x
		FValues[i].X = x
		fValues[i].Y = uniform.Prob(x)
		FValues[i].Y = uniform.CDF(x)
	}

	// Create line plots
	fLine, err := plotter.NewLine(fValues)
	if err != nil {
		return nil, err
	}
	fLine.Color = color.RGBA{R: 255, A: 255}

	FLine, err := plotter.NewLine(FValues)
	if err != nil {
		return nil, err
	}
	FLine.Color = color.RGBA{G: 255, A: 255}

	p.Add(fLine, FLine)
	p.Legend.Add("PDF f(x)", fLine)
	p.Legend.Add("CDF F(x)", FLine)

	return p, nil
}

// Function to draw Normal distribution graph
func drawNormalGraph(mu, sigma float64) (*plot.Plot, error) {
	normal := distuv.Normal{Mu: mu, Sigma: sigma}

	p := plot.New()
	p.Title.Text = "Normal Distribution"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "f(x), F(x)"

	// Use 1000 points to plot the graph
	n := 1000
	xValues := make(plotter.XYs, n)
	fValues := make(plotter.XYs, n)
	FValues := make(plotter.XYs, n)

	// Generate points for the normal distribution
	for i := 0; i < n; i++ {
		x := mu - 4*sigma + (float64(i)/float64(n))*(mu+4*sigma-(mu-4*sigma))
		xValues[i].X = x
		fValues[i].X = x
		FValues[i].X = x
		fValues[i].Y = normal.Prob(x)
		FValues[i].Y = normal.CDF(x)
	}

	// Create line plots for the probability density function (PDF) and cumulative distribution function (CDF)
	fLine, err := plotter.NewLine(fValues)
	if err != nil {
		return nil, err
	}
	fLine.Color = color.RGBA{R: 255, A: 255}

	FLine, err := plotter.NewLine(FValues)
	if err != nil {
		return nil, err
	}
	FLine.Color = color.RGBA{G: 255, A: 255}

	p.Add(fLine, FLine)
	p.Legend.Add("PDF f(x)", fLine)
	p.Legend.Add("CDF F(x)", FLine)

	return p, nil
}

func main() {
	// Create new application
	myApp := app.New()
	myWindow := myApp.NewWindow("Graph Plotter")

	// Create input fields for uniform distribution
	aInput := widget.NewEntry()
	aInput.SetPlaceHolder("Enter a")
	bInput := widget.NewEntry()
	bInput.SetPlaceHolder("Enter b")

	// Create input fields for Normal distribution
	muInput := widget.NewEntry()
	muInput.SetPlaceHolder("Enter mu")
	sigmaInput := widget.NewEntry()
	sigmaInput.SetPlaceHolder("Enter sigma")

	// Create buttons
	drawUniformButton := widget.NewButton("Draw Uniform", func() {
		a, err := strconv.ParseFloat(aInput.Text, 64)
		if err != nil {
			fmt.Println("Invalid input for a")
			return
		}
		b, err := strconv.ParseFloat(bInput.Text, 64)
		if err != nil {
			fmt.Println("Invalid input for b")
			return
		}
		p, err := drawUniformGraph(a, b)
		if err != nil {
			fmt.Println("Error drawing uniform graph:", err)
			return
		}
		err = p.Save(4*vg.Inch, 4*vg.Inch, "uniform.svg")
		if err != nil {
			log.Fatal(err)
		}
	})

	// Create button for Normal distribution
	drawNormalButton := widget.NewButton("Draw Normal", func() {
		mu, err := strconv.ParseFloat(muInput.Text, 64)
		if err != nil {
			fmt.Println("Invalid input for mu")
			return
		}
		sigma, err := strconv.ParseFloat(sigmaInput.Text, 64)
		if err != nil {
			fmt.Println("Invalid input for sigma")
			return
		}
		p, err := drawNormalGraph(mu, sigma)
		if err != nil {
			fmt.Println("Error drawing normal graph:", err)
			return
		}
		err = p.Save(4*vg.Inch, 4*vg.Inch, "normal.svg")
		if err != nil {
			log.Fatal(err)
		}
	})

	// Layout
	formLayout := container.NewVBox(
		widget.NewLabel("Uniform Distribution"),
		aInput, bInput, drawUniformButton,
		widget.NewLabel("Normal Distribution"),
		muInput, sigmaInput, drawNormalButton,
	)

	// Set the content of the window
	myWindow.SetContent(formLayout)
	myWindow.Resize(fyne.NewSize(800, 400))
	myWindow.ShowAndRun()
}
