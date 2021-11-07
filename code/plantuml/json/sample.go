//go:build run
// +build run

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	paapi5 "github.com/spiegel-im-spiegel/pa-api"
	"github.com/spiegel-im-spiegel/pa-api/query"
)

func main() {
	//Create client
	client := paapi5.New(
		paapi5.WithMarketplace(paapi5.LocaleJapan),
	).CreateClient(
		"mytag-20",
		"AKIAIOSFODNN7EXAMPLE",
		"1234567890",
	)

	//Make query
	q := query.NewGetItems(
		client.Marketplace(),
		client.PartnerTag(),
		client.PartnerType(),
	).
		ASINs([]string{"B09HK66P5X"}).
		EnableItemInfo().
		EnableImages().
		EnableParentASIN()

	//Requet and response
	body, err := client.Request(q)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	fmt.Println("@startjson \"book\"")
	if _, err := io.Copy(os.Stdout, bytes.NewReader(body)); err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	fmt.Println("\n@endjson")
}
