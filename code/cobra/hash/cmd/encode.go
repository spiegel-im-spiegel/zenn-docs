package cmd

import (
	"crypto"
	"fmt"

	"sample/hash/encode"

	"github.com/spf13/cobra"
	"github.com/spiegel-im-spiegel/gocli/rwi"
)

func newEncodeCmd(ui *rwi.RWI) *cobra.Command {
	encodeCmd := &cobra.Command{
		Use:     "encode",
		Aliases: []string{"enc", "e"},
		Short:   "hash input data",
		Long:    "hash input data (detail)",
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := encode.Value(ui.Reader(), crypto.SHA256)
			if err != nil {
				return err
			}
			fmt.Fprintf(ui.Writer(), "%x\n", v)
			return nil
		},
	}
	return encodeCmd
}
