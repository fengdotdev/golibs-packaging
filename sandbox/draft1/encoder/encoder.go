package encoder

import (
	"errors"
)

var (
	marker = "@" // Fixed 1-byte marker
)

func UpdateMarker(newMarker string) error {
	if len(newMarker) != 1 {
		return errors.New("marker must be exactly 1 character")
	}

	// Check if the new marker is a valid ASCII character
	if newMarker[0] < 32 || newMarker[0] > 126 {
		return errors.New("marker must be a valid ASCII character (32-126)")
	}

	marker = newMarker
	return nil
}
func GetMarker() string {
	return marker
}

// search for marker and Header in the data, and reemplace with
func Escaper(in []byte, search string, replace string) []byte {
	if len(search) == 0 || len(replace) == 0 {
		return in // No replacement needed
	}

	out := make([]byte, 0, len(in))
	for i := 0; i < len(in); i++ {
		if i+len(search) <= len(in) && string(in[i:i+len(search)]) == search {
			out = append(out, []byte(replace)...)
			i += len(search) - 1 // Skip the length of the search string
		} else {
			out = append(out, in[i])
		}
	}
	return out
}
func EscaperReverse(in []byte, search string, replace string) []byte {
	if len(search) == 0 || len(replace) == 0 {
		return in // No replacement needed
	}

	out := make([]byte, 0, len(in))
	for i := 0; i < len(in); i++ {
		if i+len(replace) <= len(in) && string(in[i:i+len(replace)]) == replace {
			out = append(out, []byte(search)...)
			i += len(replace) - 1 // Skip the length of the replace string
		} else {
			out = append(out, in[i])
		}
	}
	return out
}





func NewEncodeObj(header string, version uint8, data []byte) (*EncodeObj, error) {
	if len(header) != 10 {
		return nil, errors.New("header must be exactly 10 characters")
	}
	if len(data) < 0 || len(data) > 0xFFFFFFFF {
		return nil, errors.New("data length must be between 0 and 4294967295")
	}

	obj := &EncodeObj{
		Mark:    marker, // Fixed 1-byte mark
		Header:  header,
		Version: version,
		Length:  len(data),
		Data:    data,
	}

	// Calculate CRC32 checksum
	obj.CRC32 = crc32Checksum(data)

	return obj, nil
}

type EncodeObj struct {

	// The structure of the encoded object is as follows:
	// MARK: 1byte + HEADER: 10bytes + VERSION : 1byte + CRC32: 4bytes + LENGTH: 8bytes + DATA: variable-length + HEADER: 10bytes + MARK: 1byte

	// example: // "@" + "abcdefghij" + version +  CRC32 +  length + data +  "abcdefghij" + "@"

	Mark    string // Fixed for 1 ascii character eg "@" (1 byte)
	Header  string // Fixed for 10 ascii characters "abcdefghij" (10 bytes)
	Version uint8  // 1-byte version 0 - 255 start from 1
	CRC32   uint32 // 4-byte CRC32 checksum
	Length  int    // fixed length to the data, 0 - 4294967295 (8 bytes)

	Data []byte // Variable-length data
}

func (e *EncodeObj) Encode() ([]byte, error) {
	panic("not implemented yet")
}

func (e *EncodeObj) Decode(data []byte) error {
	panic("not implemented yet")

}

func (e *EncodeObj) Validate() error {
	// check mark
	if e.Mark != marker {
		return errors.New("mark must be '@'")
	}

	// check header
	if len(e.Header) != 10 {
		return errors.New("header must be exactly 10 characters")
	}

	// check version
	if e.Version < 1 || e.Version > 255 {
		return errors.New("version must be between 1 and 255")
	}

	// Calculate CRC32 checksum and compare
	calculatedCRC32 := crc32Checksum(e.Data)
	if calculatedCRC32 != e.CRC32 {
		return errors.New("CRC32 checksum does not match")
	}

	// check length
	if len(e.Data) != e.Length {
		return errors.New("data length does not match specified length")
	}

	return nil
}
