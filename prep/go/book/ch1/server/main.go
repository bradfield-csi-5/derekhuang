// Server1 is a minimal "echo" server.
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lissajous(w, r.URL.Query())
	})
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

var palette = []color.Color{color.Black, color.RGBA{0x55, 0xAF, 0x55, 0xFF}}

const (
	blackIndex = 0
	greenIndex = 1
)

func lissajous(out io.Writer, params url.Values) {
	var cycles = 5.0 // number of complete x oscillator revolutions
	var res = 0.001  // angular resolution
	var size = 100.0 // image canvas covers [-size..+size]
	var nframes = 64 // number of animation frames
	var delay = 8    // delay between frames in 10ms units

	for k, v := range params {
		switch k {
		case "cycles":
			num, err := strconv.Atoi(v[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "lissajous: converting param %q: %q\n", k, v)
				return
			}
			cycles = float64(num)
		case "res":
			float, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "lissajous: converting param %q: %q\n", k, v)
				return
			}
			res = float64(float)
		case "size":
			num, err := strconv.Atoi(v[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "lissajous: converting param %q: %q\n", k, v)
				return
			}
			size = float64(num)
		case "nframes":
			num, err := strconv.Atoi(v[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "lissajous: converting param %q: %q\n", k, v)
				return
			}
			nframes = num
		case "delay":
			num, err := strconv.Atoi(v[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "lissajous: converting param %q: %q\n", k, v)
				return
			}
			delay = num
		}
	}

	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*int(size)+1, 2*int(size)+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(int(size)+int(x*size+0.5), int(size)+int(y*size+0.5),
				greenIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}
