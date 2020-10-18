package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	_ "crypto/sha256"

	"github.com/spiegel-im-spiegel/gocli/exitcode"
	"github.com/spiegel-im-spiegel/gocli/rwi"
)

func TestEncode(t *testing.T) {
	testCases := []struct {
		inp  string
		outp string
		ext  exitcode.ExitCode
	}{
		{inp: "hello world\n", outp: "a948904f2f0f479b8f8197694b30184b0d2ed1c1cd2a1ec0fb85d299a192a447\n", ext: exitcode.Normal},
	}

	for _, tc := range testCases {
		r := strings.NewReader(tc.inp)
		wbuf := &bytes.Buffer{}
		ebuf := &bytes.Buffer{}
		ext := Execute(
			rwi.New(
				rwi.WithReader(r),
				rwi.WithWriter(wbuf),
				rwi.WithErrorWriter(ebuf),
			),
			[]string{"encode"},
		)
		if ext != tc.ext {
			t.Errorf("Execute() is \"%v\", want \"%v\".", ext, tc.ext)
			fmt.Println(ebuf.String())
		}
		str := wbuf.String()
		if str != tc.outp {
			t.Errorf("Execute() -> \"%v\", want \"%v\".", str, tc.outp)
		}
	}
}
