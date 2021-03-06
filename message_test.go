// Copyright (c) 2019 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.package encoding

package iso8583_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/rvflash/iso8583/field"

	"github.com/rvflash/iso8583/encoding"

	"github.com/matryer/is"
	"github.com/rvflash/iso8583"
)

func TestUnmarshal(t *testing.T) {
	var (
		msg = []string{
			"ascii_network_management_request",
			"ascii_headed_network_management_request",
			"ascii_network_management_response",
			//"ascii_financial_transaction_request",
		}
		are = is.New(t)
	)
	for _, name := range msg {
		name := name
		t.Run(name, func(t *testing.T) {
			src, err := message(name)
			are.NoErr(err)
			dst := new(iso8583.Message)
			dst.Format, err = encoding.Parse(src.Format)
			are.NoErr(err)
			dst.Header = src.Header
			err = iso8583.Unmarshal([]byte(src.Message), dst)
			are.NoErr(err)
			are.Equal(dst.MTI.String(), src.MTI)
			are.Equal(dst.Format, encoding.ASCII)
			are.Equal(dst.Header, src.Header)
			are.Equal(len(dst.Data), len(src.Fields))

			for k, v := range src.Fields {
				are.Equal(dst.Data[field.ID(k)].String(), v)
			}
		})
	}
}

type iso struct {
	Header  bool             `json:"header,omitempty"`
	Format  string           `json:"encoding,omitempty"`
	Message string           `json:"message"`
	MTI     string           `json:"mti,omitempty"`
	Fields  map[uint8]string `json:"fields,omitempty"`
}

func message(name string) (*iso, error) {
	b, err := ioutil.ReadFile("testdata/" + name + ".json")
	if err != nil {
		return nil, err
	}
	msg := new(iso)
	err = json.Unmarshal(b, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
