// Code generated by running "go generate" in golang.org/x/text. DO NOT EDIT.

package translations

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

type dictionary struct {
	index []uint32
	data  string
}

func (d *dictionary) Lookup(key string) (data string, ok bool) {
	p, ok := messageKeyToIndex[key]
	if !ok {
		return "", false
	}
	start, end := d.index[p], d.index[p+1]
	if start == end {
		return "", false
	}
	return d.data[start:end], true
}

func init() {
	dict := map[string]catalog.Dictionary{
		"en": &dictionary{index: enIndex, data: enData},
	}
	fallback := language.MustParse("en")
	cat, err := catalog.NewFromMap(dict, catalog.Fallback(fallback))
	if err != nil {
		panic(err)
	}
	message.DefaultCatalog = cat
}

var messageKeyToIndex = map[string]int{
	"About ":            0,
	"App Registration":  2,
	"App Version: %s":   1,
	"Found Server":      6,
	"Manual Server":     7,
	"Token":             5,
	"Use Custom Server": 3,
	"Use TLS?":          4,
}

var enIndex = []uint32{ // 9 elements
	0x00000000, 0x0000000b, 0x0000001e, 0x0000002f,
	0x00000041, 0x0000004a, 0x00000050, 0x0000005d,
	0x0000006b,
} // Size: 60 bytes

const enData string = "" + // Size: 107 bytes
	"\x04\x00\x01 \x06\x02About\x02App Version: %[1]s\x02App Registration\x02" +
	"Use Custom Server\x02Use TLS?\x02Token\x02Found Server\x02Manual Server"

	// Total table size 167 bytes (0KiB); checksum: F89462C9
