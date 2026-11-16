package dns

import "strings"

type Question struct {
	Name  string `json:"name"`
	Type  Type   `json:"type"`
	Class Class  `json:"class"`
}

func NewQuestion(name string, typ Type) Question {

	var question Question

	if typ == TypePTR {
		name = toPTRName(name)
	}

	question.SetName(name)
	question.SetType(typ)
	question.SetClass(ClassInternet)

	return question

}

func (question *Question) Bytes() []byte {

	buffer := make([]byte, 0)
	labels := renderLabels(strings.Split(question.Name, "."))

	if len(labels) > 0 {
		buffer = append(buffer, labels...)
	} else {
		buffer = append(buffer, 0)
		buffer = append(buffer, 0)
	}

	buffer = append(buffer, byte(question.Type>>8))
	buffer = append(buffer, byte(question.Type&0xff))
	buffer = append(buffer, byte(question.Class>>8))
	buffer = append(buffer, byte(question.Class&0xff))

	return buffer

}

func (question *Question) SetClass(value Class) {
	question.Class = value
}

func (question *Question) SetName(value string) {

	if strings.HasSuffix(value, ".") == false {

		labels := strings.Split(value, ".")
		valid := true
		count := int(0)

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
			question.Name = value
		}

	}

}

func (question *Question) SetType(value Type) {
	question.Type = value
}
