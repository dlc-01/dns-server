package dns

import (
	"dns-server/internal/domain/model"
	"encoding/binary"
)

func Write(msg model.Message) []byte {
	buf := make([]byte, 12)

	binary.BigEndian.PutUint16(buf[0:2], msg.Header.ID)

	flags := uint16(0)
	if msg.Header.QR {
		flags |= 1 << 15
	}
	flags |= uint16(msg.Header.Opcode) << 11
	if msg.Header.RD {
		flags |= 1 << 8
	}
	flags |= uint16(msg.Header.RCode)

	binary.BigEndian.PutUint16(buf[2:4], flags)
	binary.BigEndian.PutUint16(buf[4:6], uint16(len(msg.Questions)))
	binary.BigEndian.PutUint16(buf[6:8], uint16(len(msg.Answers)))

	out := buf

	for _, q := range msg.Questions {
		out = append(out, encodeName(q.Name)...)
		tmp := make([]byte, 4)
		binary.BigEndian.PutUint16(tmp[0:2], q.Type)
		binary.BigEndian.PutUint16(tmp[2:4], q.Class)
		out = append(out, tmp...)
	}

	for _, a := range msg.Answers {
		out = append(out, encodeName(a.Name)...)
		tmp := make([]byte, 10)
		binary.BigEndian.PutUint16(tmp[0:2], a.Type)
		binary.BigEndian.PutUint16(tmp[2:4], a.Class)
		binary.BigEndian.PutUint32(tmp[4:8], a.TTL)
		binary.BigEndian.PutUint16(tmp[8:10], uint16(len(a.Data)))
		out = append(out, tmp...)
		out = append(out, a.Data...)
	}

	return out
}
