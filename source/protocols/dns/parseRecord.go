package dns

import "encoding/binary"
import "errors"
import "strings"

func parseRecord(buffer []byte, pointer *Pointer) (*Record, error) {

	record := Record{}
	offset := pointer.Offset
	length := uint(len(buffer))

	if offset < length && offset >= 12 {

		labels, length_labels := parseLabels(buffer, pointer)

		if length_labels > 0 && len(labels) > 0 {

			record.Name = strings.Join(labels, ".")

			typ := binary.BigEndian.Uint16(buffer[offset+length_labels : offset+length_labels+2])
			class := binary.BigEndian.Uint16(buffer[offset+length_labels+2 : offset+length_labels+4])
			ttl := binary.BigEndian.Uint32(buffer[offset+length_labels+4 : offset+length_labels+8])
			data_length := binary.BigEndian.Uint16(buffer[offset+length_labels+8 : offset+length_labels+10])

			record.Type = Type(typ)
			record.Class = Class(class)

			if ttl > 0 {
				record.TTL = ttl
			}

			if data_length > 0 {

				if record.Type == TypeA {

					record.Data = []byte(buffer[offset+length_labels+10 : offset+length_labels+10+uint(data_length)])

				} else if record.Type == TypeNS {

					record_pointer := Pointer{}
					record_pointer.Offset = offset + length_labels + 10

					labels, _ := parseLabels(buffer, &record_pointer)

					if len(labels) > 0 {
						record.Data = renderLabels(labels)
					}

				} else if record.Type == TypeCNAME {

					record_pointer := Pointer{}
					record_pointer.Offset = offset + length_labels + 10

					labels, _ := parseLabels(buffer, &record_pointer)

					if len(labels) > 0 {
						record.Data = renderLabels(labels)
					}

				} else if record.Type == TypePTR {

					record_pointer := Pointer{}
					record_pointer.Offset = offset + length_labels + 10

					labels, _ := parseLabels(buffer, &record_pointer)

					if len(labels) > 0 {
						record.Data = renderLabels(labels)
					}

				} else if record.Type == TypeMX {

					record_pointer := Pointer{}
					preference := buffer[offset+length_labels+10 : offset+length_labels+12]

					record_pointer.Offset = offset + length_labels + 12
					labels, _ := parseLabels(buffer, &record_pointer)

					if len(labels) > 0 {

						data := make([]byte, 0)
						data = append(data, preference...)
						data = append(data, renderLabels(labels)...)
						record.Data = data

					}

				} else if record.Type == TypeAAAA {

					record.Data = []byte(buffer[offset+length_labels+10 : offset+length_labels+10+uint(data_length)])

				} else if record.Type == TypeSRV {

					priority := buffer[offset+length_labels+10 : offset+length_labels+12]
					weight := buffer[offset+length_labels+12 : offset+length_labels+14]
					port := buffer[offset+length_labels+14 : offset+length_labels+16]

					record_pointer := Pointer{}
					record_pointer.Offset = offset + length_labels + 16

					labels, _ := parseLabels(buffer, &record_pointer)

					if len(labels) > 0 {

						data := make([]byte, 0)
						data = append(data, priority...)
						data = append(data, weight...)
						data = append(data, port...)
						data = append(data, renderLabels(labels)...)
						record.Data = data

					}

				} else {

					record.Data = []byte(buffer[offset+length_labels+10 : offset+length_labels+10+uint(data_length)])

				}

			}

			pointer.Offset = offset + length_labels + 10 + uint(data_length)

		}

		return &record, nil

	} else {

		return nil, errors.New("Empty Record")

	}

}
