package main

import (
	"flag"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/unixpickle/essentials"
	"github.com/unixpickle/ffmpego"
	"github.com/unixpickle/model3d/model2d"
	"github.com/unixpickle/model3d/model3d"
)

func main() {
	var inPath string
	var outPath string
	var fps float64
	var duration float64
	var pauseTime float64
	var timeExponent float64
	var average bool
	flag.StringVar(&inPath, "in", "", "input file path")
	flag.StringVar(&outPath, "out", "", "output file path")
	flag.Float64Var(&fps, "fps", 5.0, "frame rate")
	flag.Float64Var(&duration, "duration", 7.0, "duration of animation")
	flag.Float64Var(&pauseTime, "pause", 1.0, "pause at end of animation")
	flag.Float64Var(&timeExponent, "exponent", 0.8, "exponent to control rate of added points")
	flag.BoolVar(&average, "average", false, "average all colors in each Voronoi cell")
	flag.Parse()

	if inPath == "" || outPath == "" {
		essentials.Die("missing required flags: -in and -out. See -help.")
	}

	f, err := os.Open(inPath)
	essentials.Must(err)
	img, _, err := image.Decode(f)
	f.Close()
	essentials.Must(err)

	b := img.Bounds()
	var coords []model2d.Coord
	for i := 0; i < b.Dx(); i++ {
		for j := 0; j < b.Dy(); j++ {
			coords = append(coords, model2d.XY(float64(i), float64(j)))
		}
	}

	for i := 0; i < len(coords)-1; i++ {
		j := i + rand.Intn(len(coords)-i)
		coords[i], coords[j] = coords[j], coords[i]
	}

	writer, err := ffmpego.NewVideoWriter(outPath, b.Dx(), b.Dy(), fps)
	essentials.Must(err)
	defer writer.Close()

	interpTarget := math.Log(float64(len(coords)))
	step := 1 / (fps * duration)
	for t := step; true; t += step {
		log.Printf("timestep %f", t)
		count := int(math.Round(math.Exp(interpTarget * math.Pow(t, timeExponent))))
		count = essentials.MinInt(count, len(coords))

		var frame image.Image
		if average {
			frame = RenderFrameAverage(img, coords[:count])
		} else {
			frame = RenderFrame(img, coords[:count])
		}
		essentials.Must(writer.WriteFrame(frame))

		if count == len(coords) {
			break
		}
	}
	for i := 0.0; i < fps*pauseTime; i++ {
		essentials.Must(writer.WriteFrame(img))
	}
}

func RenderFrame(img image.Image, coords []model2d.Coord) image.Image {
	tree := model2d.NewCoordTree(coords)
	b := img.Bounds()
	res := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	for i := 0; i < b.Dy(); i++ {
		for j := 0; j < b.Dx(); j++ {
			source := tree.NearestNeighbor(model2d.XY(float64(j), float64(i)))
			c := img.At(int(source.X)+b.Min.X, int(source.Y)+b.Min.Y)
			res.Set(j, i, c)
		}
	}
	return res
}

func RenderFrameAverage(img image.Image, coords []model2d.Coord) image.Image {
	tree := model2d.NewCoordTree(coords)
	sums := map[model2d.Coord]model3d.Coord3D{}
	counts := map[model2d.Coord]float64{}

	b := img.Bounds()
	for i := 0; i < b.Dy(); i++ {
		for j := 0; j < b.Dx(); j++ {
			source := tree.NearestNeighbor(model2d.XY(float64(j), float64(i)))
			r, g, b, _ := img.At(j+b.Min.X, i+b.Min.Y).RGBA()
			rgb := model3d.XYZ(float64(r)/256, float64(g)/256, float64(b)/256)
			sums[source] = sums[source].Add(rgb)
			counts[source]++
		}
	}

	res := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	for i := 0; i < b.Dy(); i++ {
		for j := 0; j < b.Dx(); j++ {
			source := tree.NearestNeighbor(model2d.XY(float64(j), float64(i)))
			avg := sums[source].Scale(1 / counts[source])
			res.Set(j, i, color.RGBA{
				R: uint8(avg.X),
				G: uint8(avg.Y),
				B: uint8(avg.Z),
				A: 0xff,
			})
		}
	}
	return res
}
