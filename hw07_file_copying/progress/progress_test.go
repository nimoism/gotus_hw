package progress

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBar(t *testing.T) {
	var total int64 = 100
	bar10 := NewBar(total, &bytes.Buffer{}, 10)
	bar20 := NewBar(total, &bytes.Buffer{}, 20)
	bar30 := NewBar(total, &bytes.Buffer{}, 30)

	for i, tst := range []struct {
		bar   *Bar
		add   int
		wrote int // zero as skip
		out   string
		reset bool
		err   error
	}{
		{bar: bar10, add: 0, out: "\r[...] 0%", reset: true},
		{bar: bar10, add: 15, out: "\r[...] 15%"},
		{bar: bar10, add: 30, out: "\r[#..] 45%"},
		{bar: bar10, add: 30, out: "\r[##.] 75%"},
		{bar: bar10, add: 25, out: "\r[###] 100%\n"},

		{bar: bar20, add: 50, out: "\r[######.......] 50%", reset: true},
		{bar: bar20, add: 10, out: "\r[#######......] 60%"},

		{bar: bar30, add: 50, out: "\r[###########............] 50%", reset: true},
		{bar: bar30, add: 50, out: "\r[#######################] 100%\n"},

		{bar: bar30, add: 50, out: "\r[###########............] 50%", reset: true},
		{bar: bar30, add: 60, wrote: 50, err: ErrWriteLimitExceed},
	} {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if tst.reset {
				tst.bar.Reset()
			}
			out := tst.bar.out.(*bytes.Buffer)
			out.Reset()
			n, err := tst.bar.Write(make([]byte, tst.add))
			if tst.err != nil {
				require.EqualError(t, err, tst.err.Error())
			} else {
				require.NoError(t, err)
			}
			var nExp int
			if tst.wrote == 0 {
				nExp = tst.add
			} else {
				nExp = tst.wrote
			}
			require.Equal(t, nExp, n)
			require.Equal(t, tst.out, out.String())
		})
	}
}
