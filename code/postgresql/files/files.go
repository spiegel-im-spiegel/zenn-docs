package files

import (
	"io"
	"os"

	"github.com/spiegel-im-spiegel/errs"
)

func GetBinary(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("path", path))
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	return b, errs.Wrap(err, errs.WithContext("path", path))
}
