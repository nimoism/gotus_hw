package progress

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
)

type Term struct{}

func (t Term) Size() (width, height int, err error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}
	return t.size(out)
}

func (t Term) size(text []byte) (width, height int, err error) {
	size := bytes.Split(bytes.Trim(text, "\n"), []byte(" "))
	height, err = strconv.Atoi(string(size[0]))
	if err != nil {
		return 0, 0, err
	}
	width, err = strconv.Atoi(string(size[1]))
	if err != nil {
		return 0, 0, err
	}
	return width, height, nil
}
