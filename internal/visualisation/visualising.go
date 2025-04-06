package visualisation

import (
	"image"
	"sync"

	fyne "fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	fynecanvas "fyne.io/fyne/v2/canvas"
	fynecontainer "fyne.io/fyne/v2/container"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

// An example of how to use the visualiser can be seen in README.md contained in this package.

type Visualiser interface {
	// Create window with static plot
	CreatePlotWindow(p *plot.Plot) Figure

	// Create empty plot, used in e.g. real-time plots
	CreateEmptyPlotWindow() Figure

	// Start the GUI application, generally is blocking and
	// needs to run on main thread
	Run()
}

type Figure interface{
	// A figure belongs to a specific visualiser
	// e.g. A FayneFigure is used to update plots in Fayne windows

	UpdatePlotWindow(p *plot.Plot)
}

type fayneVisualiser struct {
	app        fyne.App
	mainWindow fyne.Window
}

type fayneFigure struct {
	imageCanvas *fynecanvas.Image
}

var (
	globalVis Visualiser
	globalVisOnce sync.Once
)

/*
	Creates singleton visualiser that needs to run from main thread

	Return another Visualiser to change the backend GUI
*/
func GetVisualiser() Visualiser {
	globalVisOnce.Do(func() {
		visApp := fyneapp.New()
		FayneVis := &fayneVisualiser{app: visApp}
		globalVis = FayneVis

		// Create a persistent invisible window
		// One window needs to be open or main thread will die
		FayneVis.mainWindow = visApp.NewWindow("Persistent Window")
		FayneVis.mainWindow.Resize(fyne.NewSize(0, 0))
		FayneVis.mainWindow.SetFixedSize(true)
		FayneVis.mainWindow.SetCloseIntercept(func() {
			FayneVis.mainWindow.Hide()
		})
		FayneVis.mainWindow.Show() // Keep this window always open

		hide := func(){
			// Hack to hide the main window on startup to keep app running....
			// time.Sleep(100 * time.Millisecond)
			FayneVis.mainWindow.Hide()
		}
		visApp.Lifecycle().SetOnStarted(hide)
		visApp.Lifecycle().SetOnEnteredForeground(hide)
		visApp.Lifecycle().SetOnExitedForeground(hide)
	})
	return globalVis
}

func (vis *fayneVisualiser) Run() {
	vis.app.Run()
}

func (fig *fayneFigure) UpdatePlotWindow(p *plot.Plot) {
	img := getPlotImg(p)

	fig.imageCanvas.Image = img
	fynecanvas.Refresh(fig.imageCanvas)

}

// Creates an empty plotting window
func (vis *fayneVisualiser) CreateEmptyPlotWindow() Figure {
	w := vis.app.NewWindow("")

	imageCanvas := fynecanvas.NewImageFromImage(image.NewAlpha(image.Rect(0, 0, 600, 400)))
	imageCanvas.FillMode = fynecanvas.ImageFillOriginal

	content := fynecontainer.NewCenter(imageCanvas)
	w.SetContent(content)
	w.Resize(fyne.NewSize(700, 500))
	w.Show()

	return &fayneFigure{imageCanvas}
}

// Create a plotting window with initial plot
func (vis *fayneVisualiser) CreatePlotWindow(p *plot.Plot) Figure {
	w := vis.app.NewWindow(p.Title.Text)

	imageCanvas := fynecanvas.NewImageFromImage(image.NewAlpha(image.Rect(0, 0, 600, 400)))
	imageCanvas.FillMode = fynecanvas.ImageFillOriginal

	content := fynecontainer.NewCenter(imageCanvas)
	w.SetContent(content)
	w.Resize(fyne.NewSize(700, 500))
	w.Show()

	// Get image of current plot
	img := getPlotImg(p)

	imageCanvas.Image = img
	fynecanvas.Refresh(imageCanvas)
	return &fayneFigure{imageCanvas}
}

func getPlotImg(p *plot.Plot) image.Image{

	imgCanvas := vgimg.New(6*vg.Inch, 4*vg.Inch)
	p.Draw(draw.New(imgCanvas))
	return imgCanvas.Image()
}

func ScatterPlt(points plotter.XYs) *plot.Plot {
	p := plot.New()
	p.Title.Text = "Real-Time Plot"
	p.X.Label.Text = "X-axis"
	p.Y.Label.Text = "Y-axis"

	scatter, err := plotter.NewScatter(points)
	if err != nil {
		panic(err)
	}
	p.Add(scatter)

	return p
}

func LinePlt(points plotter.XYs) *plot.Plot {
	p := plot.New()
	p.Title.Text = "Real-Time Plot"
	p.X.Label.Text = "X-axis"
	p.Y.Label.Text = "Y-axis"

	line, err := plotter.NewLine(points)
	if err != nil {
		panic(err)
	}
	p.Add(line)

	return p
}
