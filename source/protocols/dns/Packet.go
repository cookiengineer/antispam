package dns

import "antispam/types"
import "math/rand"

type Pointer struct {
	Offset uint `json:"offset"`
}

type Packet struct {
	Identifier  uint16     `json:"identifier"`
	Type        string     `json:"type"`
	Questions   []Question `json:"questions"`
	Answers     []Record   `json:"answers"`
	Authorities []Record   `json:"authorities"`
	Additionals []Record   `json:"additionals"`
	Codes       struct {
		Operation OperationCode `json:"opcode"`
		Response  ResponseCode  `json:"rcode"`
	} `json:"codes"`
	Flags struct {
		AuthorativeAnswer  bool `json:"AA"`
		Truncated          bool `json:"TC"`
		RecursionAvailable bool `json:"RA"`
		RecursionDesired   bool `json:"RD"`
	} `json:"flags"`
	Server *types.Server `json:"server"`
}

func NewPacket() Packet {

	var packet Packet

	packet.Identifier  = uint16(rand.Uint64())
	packet.Questions   = make([]Question, 0)
	packet.Answers     = make([]Record, 0)
	packet.Authorities = make([]Record, 0)
	packet.Additionals = make([]Record, 0)
	packet.Server      = nil

	return packet

}

func (packet *Packet) SetIdentifier(value uint16) {
	packet.Identifier = value
}

func (packet *Packet) SetOperationCode(value OperationCode) {
	packet.Codes.Operation = value
}

func (packet *Packet) SetResponseCode(value ResponseCode) {
	packet.Codes.Response = value
}

func (packet *Packet) SetType(value string) {

	if value == "query" {
		packet.Type = "query"
	} else if value == "response" {
		packet.Type = "response"
	}

}

func (packet *Packet) AddQuestion(value Question) {
	packet.Questions = append(packet.Questions, value)
}

func (packet *Packet) AddAnswer(value Record) {
	packet.Answers = append(packet.Answers, value)
}

func (packet *Packet) AddAuthority(value Record) {
	packet.Authorities = append(packet.Authorities, value)
}

func (packet *Packet) AddAdditional(value Record) {
	packet.Additionals = append(packet.Additionals, value)
}

func (packet *Packet) Bytes() []byte {

	bytes := make([]byte, 12)

	bytes[0] = byte(packet.Identifier >> 8)
	bytes[1] = byte(packet.Identifier & 0xff)

	if packet.Type == "query" {

		packet.Flags.RecursionDesired = true

		bytes[2] = byte(0b00000001)
		bytes[3] = byte(0b00000000)

	} else if packet.Type == "response" {

		packet.Flags.RecursionAvailable = true

		bytes[2] = byte(0b10000001)
		bytes[3] = byte(0b10000000)

	}

	if packet.Codes.Operation > 0 && packet.Codes.Operation < 16 {
		bytes[2] |= byte(packet.Codes.Operation << 3)
	}

	if packet.Flags.AuthorativeAnswer == true {
		bytes[2] |= 0b00000100
	}

	if packet.Flags.Truncated == true {
		bytes[2] |= 0b00000010
	}

	if packet.Flags.RecursionDesired == true {
		bytes[2] |= 0b00000001
	}

	if packet.Flags.RecursionAvailable == true {
		bytes[3] |= 0b10000000
	}

	if packet.Codes.Response > 0 && packet.Codes.Response < 16 {
		bytes[3] |= byte(packet.Codes.Response)
	}

	bytes[4] = byte(len(packet.Questions) >> 8)
	bytes[5] = byte(len(packet.Questions) & 0xff)
	bytes[6] = byte(len(packet.Answers) >> 8)
	bytes[7] = byte(len(packet.Answers) & 0xff)
	bytes[8] = byte(len(packet.Authorities) >> 8)
	bytes[9] = byte(len(packet.Authorities) & 0xff)
	bytes[10] = byte(len(packet.Additionals) >> 8)
	bytes[11] = byte(len(packet.Additionals) & 0xff)

	for q := 0; q < len(packet.Questions); q++ {

		bytes_question := packet.Questions[q].Bytes()
		bytes = append(bytes, bytes_question...)

	}

	for a := 0; a < len(packet.Answers); a++ {

		bytes_record := packet.Answers[a].Bytes()
		bytes = append(bytes, bytes_record...)

	}

	for a := 0; a < len(packet.Authorities); a++ {

		bytes_record := packet.Authorities[a].Bytes()
		bytes = append(bytes, bytes_record...)

	}

	for a := 0; a < len(packet.Additionals); a++ {

		bytes_record := packet.Additionals[a].Bytes()
		bytes = append(bytes, bytes_record...)

	}

	return bytes

}

func (packet *Packet) SetServer(value types.Server) {
	packet.Server = &value
}
