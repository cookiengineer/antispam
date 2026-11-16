package dns

import "encoding/binary"
import "errors"
import "strings"

func parseQuestion(buffer []byte, pointer *Pointer) (*Question, error) {

	question := Question{}
	offset := pointer.Offset
	length := uint(len(buffer))

	if offset < length && offset >= 12 {

		labels, length_labels := parseLabels(buffer, pointer)

		if length_labels > 0 && len(labels) > 0 {

			question.Name = strings.Join(labels, ".")

			typ := binary.BigEndian.Uint16(buffer[offset+length_labels : offset+length_labels+2])
			class := binary.BigEndian.Uint16(buffer[offset+length_labels+2 : offset+length_labels+4])

			question.Type = Type(typ)
			question.Class = Class(class)

			pointer.Offset = offset + length_labels + 4

		}

		return &question, nil

	} else {

		return nil, errors.New("Empty Question")

	}

}
