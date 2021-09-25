package files

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	"github.com/spiegel-im-spiegel/errs"
)

func Output(dst io.Writer, src interface{}) error {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	if err := enc.Encode(src); err != nil {
		return errs.Wrap(err)
	}
	if _, err := io.Copy(os.Stdout, buf); err != nil {
		return errs.Wrap(err)
	}
	return nil
}
