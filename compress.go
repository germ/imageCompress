// This package (en|de)codes arbitrary data into a PNG with lzw compression
package imageCompress

import (
	"io"
	"io/ioutil"
	"image"
	"image/png"
	"image/color"
	"compress/lzw"
	"bytes"
	"math"
)


//Reads all data returning the data compressed with LZW
func compress(in io.Reader) ([]byte, error) {
	//Read file and write lzw compressed data to buffer
	inData, err := ioutil.ReadAll(in)

	outData := new(bytes.Buffer)
	outWriter:= lzw.NewWriter(outData, lzw.LSB, 8)
	outWriter.Write(inData)
	outWriter.Close()
	return outData.Bytes(), err
}

//Decompress LZW encoded data stored in raw
func extract(raw []byte) ([]byte, error) {
	buf := bytes.NewBuffer(raw)

	convWriter := lzw.NewReader(buf, lzw.LSB, 8)
	out, err := ioutil.ReadAll(convWriter)

	return out, err
}

//Creates an image from data and writes it to a file
func GenerateImage(in io.Reader, out io.Writer) error {
	raw, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	// Compress data
	buf := bytes.NewBuffer(raw)
	data, err := compress(buf)
	if err != nil{
		return err
	}

	// Write data to image, 3bytes/pixel
	pixels := int(math.Ceil(float64(len(data))/6.0))
	sq := int(math.Ceil(math.Sqrt(float64(pixels))))

	img := image.NewNRGBA64(image.Rect(0,0,sq,sq))

	// Pad the data buffer
	data = append(data, make([]byte, pixels*6 - len(data))...)

	// Embed data into image
	min, max := img.Bounds().Min, img.Bounds().Max
	i := 0
	for y := min.Y; y < max.Y; y++ {
		for x := min.X; x < max.X; x++ {
			img.SetNRGBA64(x, y, embedColor(data[i:i+6]))

			i += 6
			if i >= len(data) {
				goto L
			}
		}
	}
	L:

	err = png.Encode(out, img)
	return err
}

//Reads an image in and extracts the encoded data
func ExtractImage(in io.Reader, out io.Writer) error {
	var data []byte

	//Extract data from image colors
	img, err := png.Decode(in)
	if err != nil {
		return err
	}

	// Remove the data byte by byte
	min, max := img.Bounds().Min, img.Bounds().Max
	for y := min.Y; y < max.Y; y++ {
		for x := min.X; x < max.X; x++ {
			c := img.At(x,y)
			data = append(data, extractColor(c)...)
		}
	}

	// Decompress stream
	raw, err := extract(data)
	out.Write(raw)

	return err
}

// Convert a slice of bytes into a 64bit color
func embedColor(in []byte) color.NRGBA64 {
	// Cast array to uint16
	res := make([]uint16, 6)
	for i, v := range(in) {
		res[i] = uint16(v)
	}

	r := (res[0] << 8) | res[1]
	g := (res[2] << 8) | res[3]
	b := (res[4] << 8) | res[5]

	return color.NRGBA64{r, g, b, 0xFFFF}
}

// Extract bytes from a 64 bit color
func extractColor(in color.Color) []byte {
	ret := make([]byte, 6)
	r,g,b,_ := in.RGBA()

	ret[0] = byte(r >> 8)
	ret[1] = byte(r)
	ret[2] = byte(g >> 8)
	ret[3] = byte(g)
	ret[4] = byte(b >> 8)
	ret[5] = byte(b)

	return ret
}
