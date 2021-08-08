// +build run

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spiegel-im-spiegel/gpgpdump/parse"
	"github.com/spiegel-im-spiegel/gpgpdump/parse/context"
)

const openpgpStr = `
-----BEGIN PGP SIGNATURE-----
Version: GnuPG v2

iF4EARMIAAYFAlTDCN8ACgkQMfv9qV+7+hg2HwEA6h2iFFuCBv3VrsSf2BREQaT1
T1ZprZqwRPOjiLJg9AwA/ArTwCPz7c2vmxlv7sRlRLUI6CdsOqhuO1KfYXrq7idI
=ZOTN
-----END PGP SIGNATURE-----
`

func main() {
	p, err := parse.New(
		context.New(
			context.Set(context.ARMOR, true),
			context.Set(context.UTC, true),
		),
		strings.NewReader(openpgpStr),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		return
	}
	res, err := p.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		return
	}
	buf := &bytes.Buffer{}
	if err := toml.NewEncoder(buf).Encode(res); err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		return
	}
	if _, err = io.Copy(os.Stdout, buf); err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		return
	}
}
