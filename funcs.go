package bbrpc

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

func toJSONIndent(v interface{}) string { b, _ := json.MarshalIndent(v, "", "  "); return string(b) }

func ps(s string) *string      { return &s }
func pbool(b bool) *bool       { return &b }
func puint(i uint) *uint       { return &i }
func pstring(s string) *string { return &s }
func pint(i int) *int          { return &i }

// Pstring .
func Pstring(s string) *string { return &s }

// UtilDataEncoding 将tx data 进行编码
func UtilDataEncoding(data string) string {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(time.Now().Unix()))

	return strings.Join([]string{
		strings.Replace(uuid.New().String(), "-", "", -1),
		hex.EncodeToString(b),
		"00",
		hex.EncodeToString([]byte(data)),
	}, "")
}

// UtilDataDecoding .
func UtilDataDecoding(data string) (DataDetail, error) {
	var dd DataDetail
	if l := len(data); l < 32+8+2 {
		return dd, fmt.Errorf("invalid len: %d, should > 42", l)
	}
	dd.UUID = data[:32]

	timeBytes, err := hex.DecodeString(data[32 : 32+8])
	if err != nil {
		return dd, fmt.Errorf("unable to decode time, %v", err)
	}
	dd.UnixTime = binary.LittleEndian.Uint32(timeBytes)

	content, err := hex.DecodeString(data[42:])
	if err != nil {
		return dd, fmt.Errorf("unable to decode content, %v", err)
	}
	dd.Data = string(content)
	return dd, nil
}

// DataDetail .
type DataDetail struct {
	UUID     string
	UnixTime uint32
	Data     string
}
