package dns

import "encoding/binary"

func parseName(buf []byte, offset int) (string, int) {
	name := ""
	orig := offset
	jumped := false

	for {
		l := int(buf[offset])
		if l&0xC0 == 0xC0 {
			ptr := int(binary.BigEndian.Uint16(buf[offset:offset+2]) & 0x3FFF)
			if !jumped {
				orig = offset + 2
			}
			offset = ptr
			jumped = true
			continue
		}
		offset++
		if l == 0 {
			break
		}
		if name != "" {
			name += "."
		}
		name += string(buf[offset : offset+l])
		offset += l
	}
	if jumped {
		return name, orig
	}
	return name, offset
}

func encodeName(name string) []byte {
	out := []byte{}
	start := 0
	for i := 0; i <= len(name); i++ {
		if i == len(name) || name[i] == '.' {
			out = append(out, byte(i-start))
			out = append(out, name[start:i]...)
			start = i + 1
		}
	}
	return append(out, 0)
}
