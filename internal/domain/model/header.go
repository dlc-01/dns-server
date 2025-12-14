package model

type Header struct {
	ID     uint16
	QR     bool
	Opcode uint8
	RD     bool
	RCode  uint8
}
