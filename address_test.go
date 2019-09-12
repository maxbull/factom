// MIT License
//
// Copyright 2018 Canonical Ledgers, LLC
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package factom

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// Valid Test Addresses generated by factom-walletd
	// OBVIOUSLY NEVER USE THESE FOR ANY FUNDS!
	FAAddressStr = "FA2PdKfzGP5XwoSbeW1k9QunCHwC8DY6d8xgEdfm57qfR31nTueb"
	FsAddressStr = "Fs1ipNRjEXcWj8RUn1GRLMJYVoPFBL1yw9rn6sCxWGcxciC4HdPd"
	ECAddressStr = "EC2Pawhv7uAiKFQeLgaqfRhzk5o9uPVY8Ehjh8DnLXENosvYTT26"
	EsAddressStr = "Es2tFRhAqHnydaygVAR6zbpWTQXUDaXy1JHWJugQXnYavS8ssQQE"
)

type addressUnmarshalJSONTest struct {
	Name   string
	Adr    Address
	ExpAdr Address
	Data   string
	Err    string
}

var addressUnmarshalJSONTests = []addressUnmarshalJSONTest{{
	Name: "valid FA",
	Data: fmt.Sprintf("%q", FAAddressStr),
	Adr:  new(FAAddress),
	ExpAdr: func() *FAAddress {
		adr, _ := NewFsAddress(FsAddressStr)
		pub := adr.FAAddress()
		return &pub
	}(),
}, {
	Name: "valid Fs",
	Data: fmt.Sprintf("%q", FsAddressStr),
	Adr:  new(FsAddress),
	ExpAdr: func() *FsAddress {
		adr, _ := NewFsAddress(FsAddressStr)
		return &adr
	}(),
}, {
	Name: "valid EC",
	Data: fmt.Sprintf("%q", ECAddressStr),
	Adr:  new(ECAddress),
	ExpAdr: func() *ECAddress {
		adr, _ := NewEsAddress(EsAddressStr)
		pub := adr.ECAddress()
		return &pub
	}(),
}, {
	Name: "valid Es",
	Data: fmt.Sprintf("%q", EsAddressStr),
	Adr:  new(EsAddress),
	ExpAdr: func() *EsAddress {
		adr, _ := NewEsAddress(EsAddressStr)
		return &adr
	}(),
}, {
	Name: "invalid type",
	Data: `{}`,
	Err:  "json: cannot unmarshal object into Go value of type string",
}, {
	Name: "invalid type",
	Data: `5.5`,
	Err:  "json: cannot unmarshal number into Go value of type string",
}, {
	Name: "invalid type",
	Data: `["hello"]`,
	Err:  "json: cannot unmarshal array into Go value of type string",
}, {
	Name: "invalid length",
	Data: fmt.Sprintf("%q", FAAddressStr[0:len(FAAddressStr)-1]),
	Err:  "invalid length",
}, {
	Name: "invalid length",
	Data: fmt.Sprintf("%q", FAAddressStr+"Q"),
	Err:  "invalid length",
}, {
	Name: "invalid prefix",
	Data: fmt.Sprintf("%q", func() string {
		adr, _ := NewFAAddress(FAAddressStr)
		return adr.payload().StringPrefix([]byte{0x50, 0x50})
	}()),
	Err: "invalid prefix",
}, {
	Name: "invalid prefix",
	Data: fmt.Sprintf("%q", FsAddressStr),
	Err:  "invalid prefix",
}, {
	Name:   "invalid symbol/FA",
	Data:   fmt.Sprintf("%q", FAAddressStr[0:len(FAAddressStr)-1]+"0"),
	Err:    "invalid format: version and/or checksum bytes missing",
	Adr:    new(FAAddress),
	ExpAdr: new(FAAddress),
}, {
	Name:   "invalid checksum",
	Data:   fmt.Sprintf("%q", FAAddressStr[0:len(FAAddressStr)-1]+"e"),
	Err:    "checksum error",
	Adr:    new(FAAddress),
	ExpAdr: new(FAAddress),
}}

func testAddressUnmarshalJSON(t *testing.T, test addressUnmarshalJSONTest) {
	err := json.Unmarshal([]byte(test.Data), test.Adr)
	assert := assert.New(t)
	if len(test.Err) > 0 {
		assert.EqualError(err, test.Err)
		return
	}
	assert.NoError(err)
	assert.Equal(test.ExpAdr, test.Adr)
}

func TestAddress(t *testing.T) {
	for _, test := range addressUnmarshalJSONTests {
		if test.Adr != nil {
			t.Run("UnmarshalJSON/"+test.Name, func(t *testing.T) {
				testAddressUnmarshalJSON(t, test)
			})
			continue
		}
		test.ExpAdr, test.Adr = &FAAddress{}, &FAAddress{}
		t.Run("UnmarshalJSON/FA", func(t *testing.T) {
			testAddressUnmarshalJSON(t, test)
		})
		test.ExpAdr, test.Adr = &ECAddress{}, &ECAddress{}
		t.Run("UnmarshalJSON/EC", func(t *testing.T) {
			testAddressUnmarshalJSON(t, test)
		})
	}

	fa, _ := NewFAAddress(FAAddressStr)
	fs, _ := NewFsAddress(FsAddressStr)
	ec, _ := NewECAddress(ECAddressStr)
	es, _ := NewEsAddress(EsAddressStr)
	strToAdr := map[string]Address{FAAddressStr: fa, FsAddressStr: fs,
		ECAddressStr: ec, EsAddressStr: es}
	for adrStr, adr := range strToAdr {
		t.Run("MarshalJSON/"+adr.PrefixString(), func(t *testing.T) {
			data, err := json.Marshal(adr)
			assert := assert.New(t)
			assert.NoError(err)
			assert.Equal(fmt.Sprintf("%q", adrStr), string(data))
		})
		t.Run("Payload/"+adr.PrefixString(), func(t *testing.T) {
			assert.EqualValues(t, adr, adr.Payload())
		})
	}

	t.Run("FsAddress", func(t *testing.T) {
		pub, _ := NewFAAddress(FAAddressStr)
		priv, _ := NewFsAddress(FsAddressStr)
		assert := assert.New(t)
		assert.Equal(pub, priv.FAAddress())
		assert.Equal(pub.PublicAddress(), priv.PublicAddress())
		assert.Equal(pub.RCDHash(), priv.RCDHash(), "RCDHash")
	})
	t.Run("EsAddress", func(t *testing.T) {
		pub, _ := NewECAddress(ECAddressStr)
		priv, _ := NewEsAddress(EsAddressStr)
		assert := assert.New(t)
		assert.Equal(pub, priv.ECAddress())
		assert.Equal(pub.PublicAddress(), priv.PublicAddress())
	})

	t.Run("New", func(t *testing.T) {
		for _, adrStr := range []string{FAAddressStr, FsAddressStr,
			ECAddressStr, EsAddressStr} {
			t.Run(adrStr, func(t *testing.T) {
				assert := assert.New(t)
				adr, err := NewAddress(adrStr)
				assert.NoError(err)
				assert.Equal(adrStr, fmt.Sprintf("%v", adr))
			})
			t.Run("Public/"+adrStr, func(t *testing.T) {
				assert := assert.New(t)
				adr, err := NewPublicAddress(adrStr)
				if adrStr[1] == 's' {
					assert.EqualError(err, "invalid prefix")
					return
				}
				assert.NoError(err)
				assert.Equal(adrStr, fmt.Sprintf("%v", adr))
			})
			t.Run("Private/"+adrStr, func(t *testing.T) {
				assert := assert.New(t)
				adr, err := NewPrivateAddress(adrStr)
				if adrStr[1] != 's' {
					assert.EqualError(err, "invalid prefix")
					return
				}
				assert.NoError(err)
				assert.Equal(adrStr, fmt.Sprintf("%v", adr))
			})
		}

		t.Run("invalid length", func(t *testing.T) {
			assert := assert.New(t)

			_, err := NewAddress("too short")
			assert.EqualError(err, "invalid length")

			_, err = NewPrivateAddress("too short")
			assert.EqualError(err, "invalid length")

			_, err = NewPublicAddress("too short")
			assert.EqualError(err, "invalid length")
		})

		t.Run("unrecognized prefix", func(t *testing.T) {
			adr, _ := NewFAAddress(FAAddressStr)
			adrStr := adr.payload().StringPrefix([]byte{0x50, 0x50})
			assert := assert.New(t)

			_, err := NewAddress(adrStr)
			assert.EqualError(err, "unrecognized prefix")

			_, err = NewPrivateAddress(adrStr)
			assert.EqualError(err, "unrecognized prefix")

			_, err = NewPublicAddress(adrStr)
			assert.EqualError(err, "unrecognized prefix")
		})
	})

	t.Run("Generate/Fs", func(t *testing.T) {
		var err error
		fs, err = GenerateFsAddress()
		assert.NoError(t, err)
	})
	t.Run("Generate/Es", func(t *testing.T) {
		var err error
		es, err = GenerateEsAddress()
		assert.NoError(t, err)
	})

	c := NewClient(nil, nil)
	t.Run("Save/Fs", func(t *testing.T) {
		err := fs.Save(c)
		assert.NoError(t, err)
	})
	t.Run("Save/Es", func(t *testing.T) {
		err := es.Save(c)
		assert.NoError(t, err)
	})

	t.Run("GetPrivateAddress/Fs", func(t *testing.T) {
		assert := assert.New(t)
		_, err := fs.GetPrivateAddress(nil)
		assert.NoError(err)

		fa = fs.FAAddress()
		newFs, err := fa.GetPrivateAddress(c)
		assert.NoError(err)
		assert.Equal(fs, newFs)
	})
	t.Run("GetPrivateAddress/Es", func(t *testing.T) {
		assert := assert.New(t)
		_, err := es.GetPrivateAddress(nil)
		assert.NoError(err)

		ec = es.ECAddress()
		newEs, err := ec.GetPrivateAddress(c)
		assert.NoError(err)
		assert.Equal(es, newEs)
		assert.Equal(ec.PublicKey(), es.PublicKey())
	})

	t.Run("GetAddresses", func(t *testing.T) {
		adrs, err := c.GetAddresses()
		assert := assert.New(t)
		assert.NoError(err)
		assert.NotEmpty(adrs)
	})
	t.Run("GetPrivateAddresses", func(t *testing.T) {
		adrs, err := c.GetPrivateAddresses()
		assert := assert.New(t)
		assert.NoError(err)
		assert.NotEmpty(adrs)
	})
	t.Run("GetFAAddresses", func(t *testing.T) {
		adrs, err := c.GetFAAddresses()
		assert := assert.New(t)
		assert.NoError(err)
		assert.NotEmpty(adrs)
	})
	t.Run("GetFsAddresses", func(t *testing.T) {
		adrs, err := c.GetFsAddresses()
		assert := assert.New(t)
		assert.NoError(err)
		assert.NotEmpty(adrs)
	})
	t.Run("GetECAddresses", func(t *testing.T) {
		adrs, err := c.GetECAddresses()
		assert := assert.New(t)
		assert.NoError(err)
		assert.NotEmpty(adrs)
	})
	t.Run("GetEsAddresses", func(t *testing.T) {
		adrs, err := c.GetEsAddresses()
		assert := assert.New(t)
		assert.NoError(err)
		assert.NotEmpty(adrs)
	})

	for _, adr := range strToAdr {
		t.Run("GetBalance/"+adr.PrefixString(), func(t *testing.T) {
			balance, err := adr.GetBalance(c)
			assert := assert.New(t)
			assert.NoError(err)
			assert.Equal(uint64(0), balance)
		})
	}
	fundedEC, _ := NewECAddress("EC1zANmWuEMYoH6VizJg6uFaEdi8Excn1VbLN99KRuxh3GSvB7YQ")
	t.Run("GetBalance/"+fundedEC.String(), func(t *testing.T) {
		balance, err := fundedEC.GetBalance(c)
		assert := assert.New(t)
		assert.NoError(err)
		assert.NotEqual(uint64(0), balance)
	})

	t.Run("Remove/Fs", func(t *testing.T) {
		err := fs.Remove(c)
		assert.NoError(t, err)
	})
	t.Run("Remove/Es", func(t *testing.T) {
		err := es.Remove(c)
		assert.NoError(t, err)
	})

	t.Run("Scan", func(t *testing.T) {
		var adr FAAddress
		err := adr.Scan(5)
		assert := assert.New(t)
		assert.EqualError(err, "invalid type")

		in := make([]byte, 32)
		in[0] = 0xff
		err = adr.Scan(in[:10])
		assert.EqualError(err, "invalid length")

		err = adr.Scan(in)
		assert.NoError(err)
		assert.EqualValues(in, adr[:])
	})

	t.Run("Value", func(t *testing.T) {
		var adr FAAddress
		adr[0] = 0xff
		val, err := adr.Value()
		assert := assert.New(t)
		assert.NoError(err)
		assert.Equal(adr[:], val)
	})

}
