package imageCompress

import (
	"testing"
	"bytes"
	"os"
)

func TestCompress(T *testing.T) {
	testText := "Hello World!"
	input := bytes.NewBufferString(testText)

	T.Log(string(input.Bytes()))
	swap, _ := compress(input)
	output, _ := extract(swap)
	T.Log(string(output))

	if testText != string(output) {
		T.Fail()
	}
}
func TestGenerateImage(T *testing.T) {
	in, err := os.Open("tests/inputData")
	out, err := os.Create("tests/outputData.png")

	// Create and save data to file
	GenerateImage(in, out)
	in.Close()
	out.Close()

	// Reopen and decode
	in, err = os.Open("tests/outputData.png")
	output := bytes.Buffer{}

	err = ExtractImage(in, &output)
	T.Log(string(output.Bytes()))

	if err != nil {
		T.Fatal(err)
	}
}
func TestEmbedColor(T *testing.T) {
	in := []byte{0x0F, 0xF0, 0x0B, 0xD0, 0x10, 0xAB}
	c := embedColor(in)
	out := extractColor(c)

	for i, v := range in {
		if v != out[i] {
			T.Log("In : ", in)
			T.Log("Out: ", out)
			T.Log("Col: ", c)
			T.Fatal()
		}
	}
}
