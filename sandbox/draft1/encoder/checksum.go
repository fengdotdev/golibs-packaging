package encoder

import "hash/crc32"

func crc32Checksum(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}
