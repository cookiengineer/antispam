package types

import "encoding/binary"
import "encoding/hex"
import "hash/crc32"
import "strings"

func toHash(values ...string) string {

	checksum := crc32.ChecksumIEEE([]byte(strings.Join(values, "-")))

	tmp := make([]byte, 4)
	binary.LittleEndian.PutUint32(tmp, checksum)

	return hex.EncodeToString(tmp)

}
