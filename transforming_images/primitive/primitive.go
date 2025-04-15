package primitive

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Mode int

const (
	ModeCombo Mode = iota
	ModeTriangle
	ModeRect
	ModeEllipse
	ModeCircle
	ModeRotatedRect
	ModeBeziers
	ModeRotatedEllipse
	ModePolygon
)

func WithMode(mode Mode) func() []string {
	return func() []string {
		return []string{"-m", fmt.Sprintf("%d", mode)}
	}
}

func Transform(image io.Reader, ext string, numShapes int, opts ...func() []string) (io.Reader, error) {
	var args []string
	for _, opt := range opts {
		args = append(args, opt()...)
	}

	in, err := tempFile("in_", ext)
	defer os.Remove(in.Name())

	out, err := tempFile("out_", ext)
	if err != nil {
		return nil, err
	}

	defer os.Remove(out.Name())

	_, err = io.Copy(in, image)
	if err != nil {
		return nil, err
	}

	_, err = primitive(in.Name(), out.Name(), numShapes, args...)
	if err != nil {
		return nil, err
	}

	b := bytes.NewBuffer(nil)
	_, err = io.Copy(b, out)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func primitive(inputFile, outputFile string, numShapes int, args ...string) (string, error) {
	argString := fmt.Sprintf("-i %s -o %s -n %d ", inputFile, outputFile, numShapes)
	modeArgs := append(strings.Fields(argString), args...)
	cmd := exec.Command("primitive", modeArgs...)
	b, err := cmd.CombinedOutput()
	if err != nil {
		return "", nil
	}

	return string(b), nil
}

func tempFile(prefix, extension string) (*os.File, error) {
	in, err := os.CreateTemp("", prefix)
	if err != nil {
		return nil, err
	}

	defer os.Remove(in.Name())

	return os.Create(fmt.Sprintf("%s.%s", in.Name(), extension))
}
