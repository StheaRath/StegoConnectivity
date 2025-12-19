package stego

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
)

// InjectMeta writes custom tEXt chunks into the PNG binary
func InjectMeta(pngData []byte, meta map[string]string) []byte {
	if len(meta) == 0 {
		return pngData
	}
	end := []byte("IEND")
	idx := bytes.LastIndex(pngData, end)
	if idx < 4 {
		return pngData
	}

	var buf bytes.Buffer
	buf.Write(pngData[:idx-4])

	for k, v := range meta {
		d := []byte(k + "\x00" + v)
		lenB := make([]byte, 4)
		binary.BigEndian.PutUint32(lenB, uint32(len(d)))

		crc := crc32.NewIEEE()
		crc.Write([]byte("tEXt"))
		crc.Write(d)

		crcB := make([]byte, 4)
		binary.BigEndian.PutUint32(crcB, crc.Sum32())

		buf.Write(lenB)
		buf.Write([]byte("tEXt"))
		buf.Write(d)
		buf.Write(crcB)
	}
	buf.Write(pngData[idx-4:])
	return buf.Bytes()
}

// ReadMeta extracts tEXt chunks from the PNG binary
func ReadMeta(pngData []byte) map[string]string {
	res := make(map[string]string)
	r := bytes.NewReader(pngData)

	// Skip PNG Signature (8 bytes)
	r.Seek(8, 0)

	buf := make([]byte, 4)
	for {
		if _, err := r.Read(buf); err != nil {
			break
		}
		lenVal := binary.BigEndian.Uint32(buf)

		r.Read(buf) // Chunk Type

		if string(buf) == "tEXt" {
			d := make([]byte, lenVal)
			r.Read(d)
			parts := bytes.SplitN(d, []byte("\x00"), 2)
			if len(parts) == 2 {
				res[string(parts[0])] = string(parts[1])
			}
			r.Seek(4, 1) // Skip CRC
		} else {
			r.Seek(int64(lenVal)+4, 1) // Skip Data + CRC
		}
	}
	return res
}
