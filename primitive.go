package handler

import (
	"bytes"
	"fmt"
	"image"
	_ "image-to-primitive/packrd"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fogleman/primitive/primitive"
	"github.com/gobuffalo/packr/v2"
	"github.com/nfnt/resize"
)

var (
	Input      string
	Outputs    flagArray
	Background string
	Configs    shapeConfigArray
	Alpha      int
	InputSize  int
	OutputSize int
	Mode       int
	Workers    int
	Nth        int
	Repeat     int
	V, VV      bool
)

type flagArray []string

func (i *flagArray) string() string {
	return strings.Join(*i, ", ")
}

func (i *flagArray) set(value string) error {
	*i = append(*i, value)
	return nil
}

type shapeConfig struct {
	Count  int
	Mode   int
	Alpha  int
	Repeat int
}
type shapeConfigArray []shapeConfig

func (i *shapeConfigArray) string() string {
	return ""
}

func (i *shapeConfigArray) set(value string) error {
	n, _ := strconv.ParseInt(value, 0, 0)
	*i = append(*i, shapeConfig{int(n), Mode, Alpha, Repeat})
	return nil
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// H handler
func H(w http.ResponseWriter, r *http.Request) {

	startParams := time.Now()

	imgURL := r.URL.Query().Get("img")
	m := r.URL.Query().Get("mode")
	n := r.URL.Query().Get("shape")
	o := r.URL.Query().Get("output")
	box := packr.New("assets", "assets")
	invalidURL, _ := box.FindString("invalid-url.jpg")
	maxShape, _ := box.FindString("max-shape.jpg")
	somethingWrong, _ := box.FindString("something-wrong.jpg")
	cannotDecode, _ := box.FindString("cannot-decode.jpg")

	endParams := time.Since(startParams)

	if len(imgURL) > 0 {
		nInt, err := strconv.Atoi(n)
		if err != nil {
			w.Header().Set("Content-Type", "image/jpeg")
			w.Header().Set("Content-Length", strconv.Itoa(len([]byte(somethingWrong))))
			w.Write([]byte(somethingWrong))
			return
		}

		if nInt > 32 {
			w.Header().Set("Content-Type", "image/jpeg")
			w.Header().Set("Content-Length", strconv.Itoa(len([]byte(maxShape))))
			w.Write([]byte(maxShape))
			return
		}

		startHTTPGet := time.Now()
		reqImg, err := http.Get(imgURL)
		if err != nil {
			w.Header().Set("Content-Type", "image/jpeg")
			w.Header().Set("Content-Length", strconv.Itoa(len([]byte(invalidURL))))
			w.Write([]byte(invalidURL))
			return
		}
		defer reqImg.Body.Close()
		endHTTPGet := time.Since(startHTTPGet)

		startReadReqBody := time.Now()
		body, err := ioutil.ReadAll(reqImg.Body)
		if err != nil {
			w.Header().Set("Content-Type", "image/jpeg")
			w.Header().Set("Content-Length", strconv.Itoa(len([]byte(cannotDecode))))
			w.Write([]byte(cannotDecode))
		}
		endReadReqBody := time.Since(startReadReqBody)

		startImgDecode := time.Now()
		img, _, _ := image.Decode(bytes.NewReader(body))
		endImgDecode := time.Since(startImgDecode)

		// run primitive
		startResize := time.Now()
		input := resize.Thumbnail(128, 128, img, resize.Bilinear)
		endResize := time.Since(startResize)

		startPrimitive := time.Now()
		rand.Seed(time.Now().UTC().UnixNano())
		var bg primitive.Color
		bg = primitive.MakeColor(primitive.AverageImageColor(input))

		Alpha = 128
		InputSize = 128
		OutputSize = 1024

		mInt, err := strconv.Atoi(m)
		if err != nil {
			panic(err)
		}

		// for some weird reason, need this reset
		Configs = []shapeConfig{}

		Configs.set(n)

		if len(Configs) == 1 {
			Configs[0].Mode = mInt
			Configs[0].Alpha = Alpha
			Configs[0].Repeat = 0
		}

		model := primitive.NewModel(img, bg, 1024, 1)
		primitive.Log(1, "%d: t=%.3f, score=%.6f\n", 0, 0.0, model.Score)
		start := time.Now()
		frame := 0
		for j, config := range Configs {
			primitive.Log(1, "count=%d, mode=%d, alpha=%d, repeat=%d\n",
				config.Count, config.Mode, config.Alpha, config.Repeat)

			for i := 0; i < config.Count; i++ {
				frame++

				// find optimal shape and add it to the model
				t := time.Now()
				n := model.Step(primitive.ShapeType(config.Mode), config.Alpha, config.Repeat)
				nps := primitive.NumberString(float64(n) / time.Since(t).Seconds())
				elapsed := time.Since(start).Seconds()
				primitive.Log(1, "%d: t=%.3f, score=%.6f, n=%d, n/s=%s\n", frame, elapsed, model.Score, n, nps)

				last := j == len(Configs)-1 && i == config.Count-1
				if last {
					buffer := new(bytes.Buffer)
					jpeg.Encode(buffer, model.Context.Image(), nil)
					endPrimitive := time.Since(startPrimitive)

					w.Header().Set("Timing-Parsed-Params", fmt.Sprintf("%v", endParams))
					w.Header().Set("Timing-HTTP-Get", fmt.Sprintf("%v", endHTTPGet))
					w.Header().Set("Timing-Read-Req-Body", fmt.Sprintf("%v", endReadReqBody))
					w.Header().Set("Timing-Img-Decode", fmt.Sprintf("%v", endImgDecode))
					w.Header().Set("Timing-Primitive", fmt.Sprintf("%v", endPrimitive))
					w.Header().Set("Timing-Resize", fmt.Sprintf("%v", endResize))

					switch o {
					case "jpg":
						w.Header().Set("Content-Type", "image/jpeg")
						jpeg.Encode(buffer, model.Context.Image(), nil)
						w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
						w.Write(buffer.Bytes())
					case "png":
						w.Header().Set("Content-Type", "image/png")
						png.Encode(buffer, model.Context.Image())
						w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
						w.Write(buffer.Bytes())
					case "svg":
						w.Header().Set("Content-Type", "image/svg+xml")
						w.Header().Set("Content-Length", strconv.Itoa(len(model.SVG())))
						w.Write([]byte(model.SVG()))
					}
				}
			}
		}
	}

}
