package haar

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/image/draw"
)

const (
	prefix            = "iqdb_"
	imgSize           = 128
	numberPixels      = imgSize * imgSize
	colorChannels     = 3
	haarCoefficients  = 40
	maxColorIntensity = 256 // Maximum value for an 8-bit color channel (0–255) for Red, Green, or Blue
)

type Signature struct {
	Avgl [colorChannels]float64
	Sig  [colorChannels][haarCoefficients]int16
}

func scaleImage(source image.Image) *image.RGBA {
	scaled := image.NewRGBA(image.Rect(0, 0, imgSize, imgSize))
	draw.CatmullRom.Scale(scaled, scaled.Bounds(), source, source.Bounds(), draw.Over, nil)
	return scaled
}

func imageToRGBChannels(img *image.RGBA) (r, g, b []float64) {
	r = make([]float64, numberPixels)
	g = make([]float64, numberPixels)
	b = make([]float64, numberPixels)

	for y := range imgSize {
		for x := range imgSize {
			color := img.RGBAAt(x, y)
			index := x + y*imgSize
			r[index], g[index], b[index] = float64(color.R), float64(color.G), float64(color.B)
		}
	}

	return r, g, b
}

func convertRGBToYIQ(r, g, b []float64) (y, i, q []float64) {
	y = make([]float64, numberPixels)
	i = make([]float64, numberPixels)
	q = make([]float64, numberPixels)

	for px := range numberPixels {
		Y := 0.299*r[px] + 0.587*g[px] + 0.114*b[px]
		I := 0.596*r[px] - 0.275*g[px] - 0.321*b[px]
		Q := 0.212*r[px] - 0.523*g[px] + 0.311*b[px]

		y[px], i[px], q[px] = Y, I, Q
	}

	return y, i, q
}

func haarTransform(src1, src2, src3 []float64) (out1, out2, out3 []float64) {
	out1 = haar2D(src1)
	out2 = haar2D(src2)
	out3 = haar2D(src3)

	// Average intensity of the three color channels
	out1[0] /= (maxColorIntensity * imgSize)
	out2[0] /= (maxColorIntensity * imgSize)
	out3[0] /= (maxColorIntensity * imgSize)

	return out1, out2, out3
}

func haar2D(src []float64) []float64 {
	out := make([]float64, numberPixels)
	copy(out, src)

	for y := range imgSize {
		row := make([]float64, imgSize)
		for x := range imgSize {
			row[x] = out[x+y*imgSize]
		}

		transformed := haar1D(row)
		for x := range imgSize {
			out[x+y*imgSize] = transformed[x]
		}
	}

	for x := range imgSize {
		column := make([]float64, imgSize)
		for y := range imgSize {
			column[y] = out[x+y*imgSize]
		}

		transformed := haar1D(column)
		for y := range imgSize {
			out[x+y*imgSize] = transformed[y]
		}
	}

	return out
}

func haar1D(src []float64) []float64 {
	x := make([]float64, imgSize)
	copy(x, src)

	N := imgSize
	for N > 1 {
		half := N / 2

		a := make([]float64, half)
		d := make([]float64, half)

		for i := range half {
			a[i] = (x[2*i] + x[2*i+1]) / math.Sqrt2
			d[i] = (x[2*i] - x[2*i+1]) / math.Sqrt2
		}

		for i := range half {
			x[i] = a[i]
			x[half+i] = d[i]
		}

		N = half
	}

	return x
}

func getLargestIndices(src []float64) [haarCoefficients]int16 {
	// Keep the original index while swapping during selection
	type indexValue struct {
		index int16
		value float64
	}

	signature := [haarCoefficients]int16{}
	sorted := make([]*indexValue, 0, numberPixels-1)

	for i := 1; i < numberPixels; i++ {
		sorted = append(sorted, &indexValue{
			index: int16(i),
			value: math.Abs(src[i]),
		})
	}

	for i := range haarCoefficients {
		largest := i

		for j := i; j < numberPixels-1; j++ {
			if sorted[j].value > sorted[largest].value {
				largest = j
			}
		}

		sorted[i], sorted[largest] = sorted[largest], sorted[i]
		index := sorted[i].index

		signature[i] = index
		if src[index] < 0 {
			signature[i] = -index
		}
	}

	return signature
}

func SignatureFromBlob(blob []byte) (*Signature, error) {
	sourceImg, _, err := image.Decode(bytes.NewReader(blob))
	if err != nil {
		return nil, wrapErr(err)
	}

	scaledImg := scaleImage(sourceImg)

	r, g, b := imageToRGBChannels(scaledImg)
	y, i, q := convertRGBToYIQ(r, g, b)
	ch1, ch2, ch3 := haarTransform(y, i, q)

	signatures := [colorChannels][haarCoefficients]int16{
		getLargestIndices(ch1),
		getLargestIndices(ch2),
		getLargestIndices(ch3),
	}

	for ch := range colorChannels {
		sort.Slice(signatures[ch][:], func(i, j int) bool {
			return signatures[ch][i] < signatures[ch][j]
		})
	}

	return &Signature{
		Avgl: [colorChannels]float64{
			ch1[0], ch2[0], ch3[0],
		},
		Sig: signatures,
	}, nil
}

func SignatureFromHash(hash string) (*Signature, error) {
	hash = strings.TrimPrefix(hash, prefix)
	reader := strings.NewReader(hash)

	avgl := [colorChannels]float64{}
	for i := range avgl {
		buf := make([]byte, 16)

		if _, err := reader.Read(buf); err != nil {
			return nil, wrapErr(err)
		}

		bits, err := strconv.ParseUint(string(buf), 16, 64)
		if err != nil {
			return nil, wrapErr(err)
		}

		avgl[i] = math.Float64frombits(bits)
	}

	signatures := [colorChannels][haarCoefficients]int16{}
	for ch := range colorChannels {
		for i := range haarCoefficients {
			buf := make([]byte, 4)

			if _, err := reader.Read(buf); err != nil {
				return nil, wrapErr(err)
			}

			val, err := strconv.ParseUint(string(buf), 16, 16)
			if err != nil {
				return nil, wrapErr(err)
			}

			signatures[ch][i] = int16(val)
		}
	}

	return &Signature{
		Avgl: avgl,
		Sig:  signatures,
	}, nil
}

func (s Signature) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString(prefix)

	for _, v := range s.Avgl {
		bits := math.Float64bits(v)
		hex := fmt.Sprintf("%016x", bits)
		buffer.WriteString(hex)
	}

	for ch := range colorChannels {
		for i := range haarCoefficients {
			hex := fmt.Sprintf("%04x", uint16(s.Sig[ch][i]))
			buffer.WriteString(hex)
		}
	}

	return buffer.String()
}

func wrapErr(err error) error {
	return fmt.Errorf("haar: %v", err)
}
