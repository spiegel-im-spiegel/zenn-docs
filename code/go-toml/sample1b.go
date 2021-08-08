// +build run

package main

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/spiegel-im-spiegel/gpgpdump/parse/result"
)

const tomlStr = `[[Packet]]
  name = "Signature Packet (tag 2)"
  note = "94 bytes"

  [[Packet.Item]]
    name = "Version"
    note = "current"
    value = "4"

  [[Packet.Item]]
    name = "Signiture Type"
    value = "Signature of a canonical text document (0x01)"

  [[Packet.Item]]
    name = "Public-key Algorithm"
    value = "ECDSA public key algorithm (pub 19)"

  [[Packet.Item]]
    name = "Hash Algorithm"
    value = "SHA2-256 (hash 8)"

  [[Packet.Item]]
    name = "Hashed Subpacket"
    note = "6 bytes"

    [[Packet.Item.Item]]
      name = "Signature Creation Time (sub 2)"
      value = "2015-01-24T02:52:15Z"

  [[Packet.Item]]
    name = "Unhashed Subpacket"
    note = "10 bytes"

    [[Packet.Item.Item]]
      name = "Issuer (sub 16)"
      value = "0x31fbfda95fbbfa18"

  [[Packet.Item]]
    dump = "36 1f"
    name = "Hash left 2 bytes"

  [[Packet.Item]]
    name = "ECDSA value r"
    note = "256 bits"

  [[Packet.Item]]
    name = "ECDSA value s"
    note = "252 bits"
`

func main() {
	info := result.Info{}
	if err := toml.Unmarshal([]byte(tomlStr), &info); err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		return
	}
	fmt.Println(info.Packets[0].Name)
}
