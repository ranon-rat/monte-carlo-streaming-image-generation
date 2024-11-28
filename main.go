package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/png"
	"io"
	"math/rand"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/fogleman/gg"
)

var (
	update    = make(map[string]chan EmptyChannel)
	images    = make(map[string]*gg.Context)
	nuevoChan = make(chan EmptyChannel)
	tActual   = int64(0)
)

const (
	boundaryWord = "MJPEGBOUNDARY"
	headerf      = "\r\n" +
		"--" + boundaryWord + "\r\n" +
		"Content-Type: image/jpeg\r\n" +
		"X-Timestamp: 0.000000\r\n" +
		"\r\n"
)

type photo struct {
	Photo string `json:"img"`
}
type EmptyChannel struct{}

func openThis(f io.Reader, id string, t int64) {
	//this open the image and print the pixels
	img, err := png.Decode(f)
	if err != nil {
		return
	}
	height, width := img.Bounds().Max.Y, img.Bounds().Max.X
	if _, e := images[id]; !e {
		k := gg.NewContext(width, height)
		k.SetRGB(0, 0, 0)
		k.DrawRectangle(0, 0, float64(width), float64(height))
		k.Fill()

		images[id] = k
		update[id] = make(chan EmptyChannel)
		nuevoChan <- EmptyChannel{}
	}
	dc := images[id]
	dc.SetRGBA255(0, 0, 0, 1)
	dc.DrawRectangle(0, 0, float64(width), float64(height))
	dc.Fill()
	for i := 0; i < 1000; i++ {
		if t != tActual {
			break
		}
		x, y := rand.Intn(width), rand.Intn(height)
		r, g, b, _ := img.At(x, y).RGBA()
		r8, g8, b8 := int(float64(r>>8)), int(float64(g>>8)), int(float64(b>>8))
		dc.SetRGB255(r8, g8, b8)
		dc.SetPixel(x, y)

	}

}
func imagePNG(input string) io.Reader {
	return base64.NewDecoder(base64.StdEncoding, strings.NewReader(input))
}
func imageMonteCarlo(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "multipart/x-mixed-replace;boundary="+boundaryWord)

	id := r.URL.Query().Get("id")
	cU := update[id]
	for cU == nil {
		<-nuevoChan
		cU = update[id]
	}
	<-cU

	for {
		image := images[id]
		if _, err := w.Write([]byte(headerf)); err != nil {
			delete(images, id)
			delete(update, id)
			return
		}
		if err := image.EncodePNG(w); err != nil {
			delete(images, id)
			delete(update, id)
			return
		}
		<-cU

	}
}
func picture(w http.ResponseWriter, r *http.Request) {
	// decode the bodyrequest
	var conf photo
	json.NewDecoder(r.Body).Decode(&conf)
	imageData := imagePNG(strings.Replace(conf.Photo, "data:image/octet-stream;base64,", "", 1))
	id := r.URL.Query().Get("id")

	t := time.Now().UnixNano()
	tActual = t
	go openThis(imageData, id, t)
	update[id] <- EmptyChannel{}

}

func main() {
	// clear the console
	out, _ := exec.Command("clear").Output()
	fmt.Println(string(out))
	// start the interface

	fmt.Println("\033[34mgo to http://localhost:8000 \033[0m")

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/image-result", imageMonteCarlo)
	http.HandleFunc("/picture", picture)
	http.ListenAndServe(":8000", nil)
	// this comment its for made the commit looks a little bit better
}
