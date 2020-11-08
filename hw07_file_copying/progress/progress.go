package progress

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

const (
	percentsWidth = 4
	minWidth      = 4 + percentsWidth // Minimum of `[#] 100%` string.
	drawErrorMsg  = "cant draw bar"
)

var ErrWriteLimitExceed = errors.New("write limit exceed")

type Bar struct {
	total int64
	read  int64
	out   io.Writer
	buf   bytes.Buffer
}

func NewBar(total int64, out io.Writer, width int) *Bar {
	if out == nil {
		out = os.Stdout
	}
	buf := bytes.NewBuffer(make([]byte, int(math.Max(float64(width), minWidth))))
	return &Bar{total: total, out: out, buf: *buf}
}

func (r *Bar) Write(b []byte) (int, error) {
	read := int64(len(b))
	if left := r.total - (r.read + read); left < 0 {
		return int(read + left), ErrWriteLimitExceed
	}
	r.read += read
	if err := r.draw(); err != nil {
		return int(read), err
	}
	return int(read), nil
}

// Draw to output like `[###..] 65% `.
func (r Bar) draw() error {
	r.buf.Reset()
	r.buf.WriteString("[")
	barWidth := r.buf.Cap() - percentsWidth - 3 // -`[] ` len
	frac := float32(r.read) / float32(r.total)
	barCopied := int(float32(barWidth) * frac)
	r.buf.WriteString(strings.Repeat("#", barCopied))
	barElapsed := barWidth - barCopied
	r.buf.WriteString(strings.Repeat(".", barElapsed))
	r.buf.WriteString("] ")
	r.buf.WriteString(fmt.Sprintf("%v%%", int(frac*100)))
	var err error
	if _, err = r.out.Write([]byte("\r")); err != nil {
		return fmt.Errorf("%v: %w", drawErrorMsg, err)
	}
	if _, err = r.buf.WriteTo(r.out); err != nil {
		return fmt.Errorf("%v: %w", drawErrorMsg, err)
	}
	if r.read == r.total {
		if _, err := r.out.Write([]byte("\n")); err != nil {
			return fmt.Errorf("%v: %w", drawErrorMsg, err)
		}
	}
	return nil
}

func (r *Bar) Reset() {
	r.read = 0
}
