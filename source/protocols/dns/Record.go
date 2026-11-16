package dns

import "antispam/types"
import "encoding/binary"
import "encoding/hex"
import net_url "net/url"
import "strconv"
import "strings"

type Record struct {
	Name  string `json:"name"`
	Type  Type   `json:"type"`
	Class Class  `json:"class"`
	TTL   uint32 `json:"ttl"`
	Data  []byte `json:"data"`
}

func NewRecord(name string, typ Type) Record {

	var record Record

	record.SetName(name)
	record.SetType(typ)
	record.SetClass(ClassInternet)
	record.SetTTL(0)

	if record.Type == TypeA {
		record.SetIPv4("0.0.0.0")
	} else if record.Type == TypeNS {
		record.SetDomain("localhost")
	} else if record.Type == TypeCNAME {
		record.SetDomain("localhost")
	} else if record.Type == TypePTR {
		record.SetIPv4("0.0.0.0")
	} else if record.Type == TypeMX {
		record.SetDomain("localhost")
	} else if record.Type == TypeTXT {
		record.SetData([]byte{})
	} else if record.Type == TypeAAAA {
		record.SetIPv6("0000:0000:0000:0000:0000:0000:0000:0000")
	} else if record.Type == TypeSRV {

		data := []byte{
			0, 1, // priority
			0, 0, // weight
			0, 0, // port
		}

		data = append(data, renderLabels([]string{"localhost"})...)

		record.SetData(data)

	} else if record.Type == TypeURI {

		data := []byte{
			0, 1, // priority
			0, 0, // weight
		}

		data = append(data, []byte("\"http://localhost/index.html\"")...)

		record.SetData(data)

	}

	return record

}

func (record *Record) Bytes() []byte {

	buffer := make([]byte, 0)
	labels := strings.Split(record.Name, ".")

	if len(labels) > 0 {

		buffer = append(buffer, renderLabels(labels)...)
		buffer = append(buffer, byte(record.Type>>8))
		buffer = append(buffer, byte(record.Type&0xff))
		buffer = append(buffer, byte(record.Class>>8))
		buffer = append(buffer, byte(record.Class&0xff))
		buffer = append(buffer, byte(record.TTL>>24))
		buffer = append(buffer, byte(record.TTL>>16))
		buffer = append(buffer, byte(record.TTL>>8))
		buffer = append(buffer, byte(record.TTL&0xff))

		if len(record.Data) > 0 {

			data := record.Data
			data_length := uint16(len(data))
			buffer = append(buffer, byte(data_length>>8))
			buffer = append(buffer, byte(data_length&0xff))
			buffer = append(buffer, data...)

		} else {

			buffer = append(buffer, byte(0))
			buffer = append(buffer, byte(0))

		}

	}

	return buffer

}

func (record *Record) SetClass(value Class) {
	record.Class = value
}

func (record *Record) SetData(value []byte) {
	record.Data = value
}

func (record *Record) SetDomain(value string) {

	if record.Type == TypeA {

		if strings.HasSuffix(value, ".") == false {
			record.Name = value
		}

	} else if record.Type == TypeNS {

		if strings.HasSuffix(value, ".") == false {

			labels := strings.Split(value, ".")
			buffer := renderLabels(labels)

			if len(buffer) > 0 {
				record.Data = buffer
			}

		}

	} else if record.Type == TypeCNAME {

		if strings.HasSuffix(value, ".") == false {

			labels := strings.Split(value, ".")
			buffer := renderLabels(labels)

			if len(buffer) > 0 {
				record.Data = buffer
			}

		}

	} else if record.Type == TypeSOA {

		rname := make([]byte, 0)
		serial := make([]byte, 4)
		refresh := make([]byte, 4)
		retry := make([]byte, 4)
		expire := make([]byte, 4)

		if len(record.Data) > 0 {

			_, mname_length := parseLabels(record.Data[0:], nil)
			rname_labels, rname_length := parseLabels(record.Data[mname_length:], nil)

			offset := mname_length + rname_length
			serial = record.Data[offset : offset+4]
			refresh = record.Data[offset+4 : offset+8]
			retry = record.Data[offset+8 : offset+12]
			expire = record.Data[offset+12 : offset+16]

			rname = renderLabels(rname_labels)

		}

		labels := strings.Split(value, ".")
		buffer := renderLabels(labels)

		if len(buffer) > 0 {

			data := make([]byte, 0)
			data = append(data, buffer...)
			data = append(data, rname...)
			data = append(data, serial...)
			data = append(data, refresh...)
			data = append(data, retry...)
			data = append(data, expire...)

			record.Data = data

		}

	} else if record.Type == TypePTR {

		if strings.HasSuffix(value, ".") == false {

			labels := strings.Split(value, ".")
			buffer := renderLabels(labels)

			if len(buffer) > 0 {
				record.Data = buffer
			}

		}

	} else if record.Type == TypeMX {

		if strings.HasSuffix(value, ".") == false {

			preference := []byte{0, 0}

			if len(record.Data) > 2 {
				preference = record.Data[0:2]
			}

			labels := strings.Split(value, ".")
			buffer := renderLabels(labels)

			if len(buffer) > 0 {

				data := make([]byte, 0)
				data = append(data, preference...)
				data = append(data, buffer...)

				record.Data = data

			}

		}

	} else if record.Type == TypeTXT {

		if strings.HasSuffix(value, ".") == false {
			record.Name = value
		}

	} else if record.Type == TypeAAAA {

		if strings.HasSuffix(value, ".") == false {
			record.Name = value
		}

	} else if record.Type == TypeSRV {

		if strings.HasSuffix(value, ".") == false {

			priority := []byte{0, 1}
			weight := []byte{0, 0}
			port := []byte{0, 0}

			if len(record.Data) > 6 {
				priority = record.Data[0:2]
				weight = record.Data[2:4]
				port = record.Data[4:6]
			}

			labels := strings.Split(value, ".")
			buffer := renderLabels(labels)

			if len(buffer) > 0 {

				data := make([]byte, 0)
				data = append(data, priority...)
				data = append(data, weight...)
				data = append(data, port...)
				data = append(data, buffer...)

				record.Data = data

			}

		}

	} else if record.Type == TypeURI {

		if strings.HasSuffix(value, ".") == false {
			record.Name = value
		}

	}

}

func (record *Record) SetName(value string) {

	if strings.HasSuffix(value, ".") == false {

		var labels []string = strings.Split(value, ".")
		var valid bool = true
		var count int = 0

		for l := 0; l < len(labels); l++ {

			var label = labels[l]

			if len(label) < 64 {
				count += 1
				count += len(label)
			} else {
				valid = false
				break
			}

		}

		if valid == true && count < 255 {
			record.Name = value
		}

	}

}

func (record *Record) SetIPv4(value string) {

	if record.Type == TypeA {

		if types.IsIPv4(value) {

			ipv4 := types.ParseIPv4(value)

			if ipv4 != nil {
				record.Data = ipv4.Bytes(32)
			}

		}

	} else if record.Type == TypePTR {

		if types.IsIPv4(value) {

			ipv4 := types.ParseIPv4(value)

			if ipv4 != nil {

				tmp := strings.Split(ipv4.String(), ".")
				labels := make([]string, 0)

				for t := len(tmp) - 1; t >= 0; t-- {
					labels = append(labels, tmp[t])
				}

				labels = append(labels, "in-addr")
				labels = append(labels, "arpa")

				record.Name = strings.Join(labels, ".")

			}

		}

	}

}

func (record *Record) SetIPv6(value string) {

	if record.Type == TypeAAAA {

		if types.IsIPv6(value) {

			ipv6 := types.ParseIPv6(value)

			if ipv6 != nil {
				record.Data = ipv6.Bytes(128)
			}

		}

	} else if record.Type == TypePTR {

		if types.IsIPv6(value) {

			ipv6 := types.ParseIPv6(value)

			if ipv6 != nil {

				str := ipv6.String()
				tmp := strings.Split(strings.ReplaceAll(str[1:len(str)-1], ":", ""), "")
				labels := make([]string, 0)

				for t := len(tmp) - 1; t >= 0; t-- {
					labels = append(labels, tmp[t])
				}

				labels = append(labels, "ip6")
				labels = append(labels, "arpa")

				record.Name = strings.Join(labels, ".")

			}

		}

	}

}

func (record *Record) SetPort(value uint16) {

	if record.Type == TypeSRV {

		if len(record.Data) > 6 {
			record.Data[4] = byte(value >> 8)
			record.Data[5] = byte(value & 0xff)
		}

	} else if record.Type == TypeURI {

		// TODO: add port to URL

	}

}

func (record *Record) SetTTL(value uint32) {
	record.TTL = value
}

func (record *Record) SetType(value Type) {
	record.Type = value
}

func (record *Record) SetURL(value string) {

	if record.Type == TypeURI {

		url, err := net_url.Parse(value)

		if err == nil {

			priority := []byte{0, 1}
			weight := []byte{0, 0}

			if len(record.Data) > 4 {
				priority = record.Data[0:2]
				weight = record.Data[2:4]
			}

			str := url.String()

			if len(str) > 0 {

				data := make([]byte, 0)
				data = append(data, priority...)
				data = append(data, weight...)
				data = append(data, byte('"'))
				data = append(data, []byte(str)...)
				data = append(data, byte('"'))

				record.Data = data

			}

		}

	}

}

func (record *Record) ToDomain() string {

	var result string

	if record.Type == TypeA {

		result = record.Name

	} else if record.Type == TypeNS {

		labels, _ := parseLabels(record.Data, nil)

		if len(labels) > 0 {
			result = strings.Join(labels, ".")
		}

	} else if record.Type == TypeCNAME {

		labels, _ := parseLabels(record.Data, nil)

		if len(labels) > 0 {
			result = strings.Join(labels, ".")
		}

	} else if record.Type == TypeSOA {

		labels, _ := parseLabels(record.Data, nil)

		if len(labels) > 0 {
			result = strings.Join(labels, ".")
		}

	} else if record.Type == TypePTR {

		labels, _ := parseLabels(record.Data, nil)

		if len(labels) > 0 {
			result = strings.Join(labels, ".")
		}

	} else if record.Type == TypeMX {

		// preference := binary.BigEndian.Uint16(record.Data[0:2])
		labels, _ := parseLabels(record.Data[2:], nil)

		if len(labels) > 0 {
			result = strings.Join(labels, ".")
		}

	} else if record.Type == TypeTXT {

		result = record.Name

	} else if record.Type == TypeAAAA {

		result = record.Name

	} else if record.Type == TypeSRV {

		// priority := binary.BigEndian.Uint16(record.Data[0:2])
		// weight := binary.BigEndian.Uint16(record.Data[2:4])
		// port := binary.BigEndian.Uint16(record.Data[4:6])

		if len(record.Data) > 6 {

			labels, _ := parseLabels(record.Data[6:], nil)

			if len(labels) > 0 {
				result = strings.Join(labels, ".")
			}

		}

	} else if record.Type == TypeURI {

		result = record.Name

	}

	return result

}

func (record *Record) ToIPv4() string {

	var result string

	if record.Type == TypeA {

		data := record.Data

		if len(data) == 4 {

			tmp := ""
			tmp += strconv.FormatUint(uint64(data[0]), 10)
			tmp += "."
			tmp += strconv.FormatUint(uint64(data[1]), 10)
			tmp += "."
			tmp += strconv.FormatUint(uint64(data[2]), 10)
			tmp += "."
			tmp += strconv.FormatUint(uint64(data[3]), 10)

			result = tmp

		}

	} else if record.Type == TypePTR {

		if strings.HasSuffix(record.Name, ".in-addr.arpa") {

			tmp1 := strings.Split(record.Name[0:len(record.Name)-13], ".")
			tmp2 := ""

			for t := len(tmp1) - 1; t >= 0; t-- {

				tmp2 += tmp1[t]

				if t > 0 {
					tmp2 += "."
				}

			}

			result = tmp2

		}

	}

	return result

}

func (record *Record) ToIPv6() string {

	var result string

	if record.Type == TypeAAAA {

		data := record.Data

		if len(data) == 16 {

			tmp := ""
			tmp += "["
			tmp += hex.EncodeToString(data[0:2])
			tmp += ":"
			tmp += hex.EncodeToString(data[2:4])
			tmp += ":"
			tmp += hex.EncodeToString(data[4:6])
			tmp += ":"
			tmp += hex.EncodeToString(data[6:8])
			tmp += ":"
			tmp += hex.EncodeToString(data[8:10])
			tmp += ":"
			tmp += hex.EncodeToString(data[10:12])
			tmp += ":"
			tmp += hex.EncodeToString(data[12:14])
			tmp += ":"
			tmp += hex.EncodeToString(data[14:16])
			tmp += "]"

			result = tmp

		}

	} else if record.Type == TypePTR {

		if strings.HasSuffix(record.Name, ".ip6.arpa") {

			tmp1 := strings.Split(record.Name[0:len(record.Name)-9], ".")
			tmp2 := ""

			for t := len(tmp1) - 1; t >= 0; t-- {

				tmp2 += tmp1[t]

				if t%4 == 0 && t > 0 {
					tmp2 += ":"
				}

			}

			result = "[" + tmp2 + "]"

		}

	}

	return result

}

func (record *Record) ToPort() uint16 {

	var result uint16

	if record.Type == TypeSRV {

		if len(record.Data) > 6 {

			port := binary.BigEndian.Uint16(record.Data[4:6])

			if port != 0 {
				result = uint16(port)
			}

		}

	}

	return result

}

func (record *Record) ToURL() string {

	var result string

	if record.Type == TypeURI {

		if len(record.Data) > 4 {

			target := record.Data[4:]

			if target[0] == '"' && target[len(target)-1] == '"' {
				result = string(target[1:len(target)-1])
			}

		}

	}

	return result

}
