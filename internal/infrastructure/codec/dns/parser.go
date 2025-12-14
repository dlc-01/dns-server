package dns

import (
	"dns-server/internal/domain/model"
	"encoding/binary"
)

func Parse(buf []byte) model.Message {
	flags := binary.BigEndian.Uint16(buf[2:4])

	msg := model.Message{
		Header: model.Header{
			ID:     binary.BigEndian.Uint16(buf[0:2]),
			QR:     flags&(1<<15) != 0,
			Opcode: uint8(flags>>11) & 0xF,
			RD:     flags&(1<<8) != 0,
			RCode:  uint8(flags & 0xF),
		},
	}

	qd := binary.BigEndian.Uint16(buf[4:6])
	an := binary.BigEndian.Uint16(buf[6:8])

	offset := 12

	for i := 0; i < int(qd); i++ {
		name, off := parseName(buf, offset)
		offset = off

		qtype := binary.BigEndian.Uint16(buf[offset : offset+2])
		qclass := binary.BigEndian.Uint16(buf[offset+2 : offset+4])
		offset += 4

		msg.Questions = append(msg.Questions, model.Question{
			Name:  name,
			Type:  qtype,
			Class: qclass,
		})
	}

	for i := 0; i < int(an); i++ {
		name, off := parseName(buf, offset)
		offset = off

		typ := binary.BigEndian.Uint16(buf[offset : offset+2])
		class := binary.BigEndian.Uint16(buf[offset+2 : offset+4])
		ttl := binary.BigEndian.Uint32(buf[offset+4 : offset+8])
		rdlen := binary.BigEndian.Uint16(buf[offset+8 : offset+10])
		offset += 10

		data := buf[offset : offset+int(rdlen)]
		offset += int(rdlen)

		msg.Answers = append(msg.Answers, model.Answer{
			Name:  name,
			Type:  typ,
			Class: class,
			TTL:   ttl,
			Data:  data,
		})
	}

	return msg
}
