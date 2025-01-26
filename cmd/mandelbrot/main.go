package main

import (
	"image"
	"image/color"
	"image/png"
	"net/http"
	"strconv"
	"text/template"

	"github.com/go-chi/chi/v5"
)

type plane struct {
	minX, minY, step float64
}

var set = &plane{-2, -2, 0.005}

func main() {
	server := chi.NewRouter()
	server.Get("/", indexHandler)
	server.Post("/mandelbrot/reset", resetHandler)
	server.Post("/mandelbrot/zoom/{x}/{y}", zoomHandler)
	server.Get("/mandelbrot/section/{x}/{y}", mandelbrotSectionHandler)
	http.ListenAndServe(":8080", server)
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	set = &plane{-2, -2, 0.005}
	w.WriteHeader(http.StatusAccepted)
}

func zoomHandler(w http.ResponseWriter, r *http.Request) {
	x, err := strconv.Atoi(chi.URLParam(r, "x"))
	if err != nil {
		http.Error(w, "Invalid x parameter", http.StatusBadRequest)
		return
	}
	y, err := strconv.Atoi(chi.URLParam(r, "y"))
	if err != nil {
		http.Error(w, "Invalid y parameter", http.StatusBadRequest)
		return
	}

	centreX := set.minX + float64(x)*set.step
	centreY := set.minY + float64(y)*set.step
	set.step /= 10
	set.minX = centreX - set.step*400
	set.minY = centreY - set.step*400
	w.WriteHeader(http.StatusAccepted)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./cmd/mandelbrot/index.html"))
	tmpl.Execute(w, nil)
}

func mandelbrotSectionHandler(w http.ResponseWriter, r *http.Request) {
	x, err := strconv.Atoi(chi.URLParam(r, "x"))
	if err != nil {
		http.Error(w, "Invalid x parameter", http.StatusBadRequest)
		return
	}
	y, err := strconv.Atoi(chi.URLParam(r, "y"))
	if err != nil {
		http.Error(w, "Invalid y parameter", http.StatusBadRequest)
		return
	}

	img := image.NewRGBA(image.Rectangle{image.Point{x, y}, image.Point{x + 100, y + 100}})
	for i := x; i < x+100; i++ {
		for j := y; j < y+100; j++ {
			bounded, iterations := isBounded(set.minX+float64(i)*set.step, set.minY+float64(j)*set.step)
			if bounded {
				img.Set(i, j, color.RGBA{0, 0, 0, 255})
			} else {
				img.Set(i, j, color.RGBA{235 - iterations*2, 220 - iterations*2, 200 - iterations*2, 255})
			}
		}
	}
	png.Encode(w, img)
}

func isBounded(cx, cy float64) (bool, uint8) {
	var x, y, xx, yy float64 = 0.0, 0.0, 0.0, 0.0
	var i uint8 = 0
	for i = 0; i < 100; i++ {
		xy := x * y
		xx = x * x
		yy = y * y
		if xx+yy > 4 {
			return false, i
		}
		x = xx - yy + cx
		y = 2*xy + cy
	}
	return true, 100
}
