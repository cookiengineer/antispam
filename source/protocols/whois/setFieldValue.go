package whois

import "errors"
import "reflect"
import "strings"

func setFieldValue(reflect_value *reflect.Value, key string, value string) error {

	reflect_name := reflect_value.String()
	reflect_type := reflect_value.Type()

	for f := 0; f < reflect_type.NumField(); f++ {

		field := reflect_value.Field(f)
		field_name := reflect_type.Field(f).Name
		field_type := reflect_value.Field(f).Type()

		tags := strings.Split(reflect_type.Field(f).Tag.Get("whois"), ",")

		for t := 0; t < len(tags); t++ {

			if tags[t] == key {

				value_type := field_type.String()

				if value_type == "bool" {

					tmp, err := parseBool(value, value_type)

					if err == nil {
						field.SetBool(tmp)
					} else {
						return errors.New("whois: Cannot unmarshal string into Go struct field " + reflect_name + "." + field_name + " of type " + value_type)
					}

				} else if value_type == "[]bool" {

					if strings.Contains(value, "\n") {

						lines := strings.Split(value, "\n")
						slice := reflect.MakeSlice(field_type, len(lines), len(lines))

						for l := 0; l < len(lines); l++ {

							tmp, err := parseBool(lines[l], value_type[2:])

							if err == nil {
								slice.Index(l).SetBool(tmp)
							}

						}

						field.Set(slice)

					} else {

						tmp, err := parseBool(value, value_type[2:])

						if err == nil {

							slice := reflect.MakeSlice(field_type, 1, 1)
							slice.Index(0).SetBool(tmp)
							field.Set(slice)

						}

					}

				} else if value_type == "string" {

					field.SetString(value)

				} else if value_type == "[]string" {

					if strings.Contains(value, "\n") {

						lines := strings.Split(value, "\n")
						slice := reflect.MakeSlice(field_type, len(lines), len(lines))

						for l := 0; l < len(lines); l++ {
							slice.Index(l).SetString(lines[l])
						}

						field.Set(slice)

					} else {

						slice := reflect.MakeSlice(field_type, 1, 1)
						slice.Index(0).SetString(value)
						field.Set(slice)

					}

				} else if value_type == "int" ||
					value_type == "int8" ||
					value_type == "int16" ||
					value_type == "int32" ||
					value_type == "int64" {

					tmp, err := parseInt(value, value_type)

					if err == nil {
						field.SetInt(tmp)
					} else {
						return errors.New("whois: Cannot unmarshal string into Go struct field " + reflect_name + "." + field_name + " of type " + value_type)
					}

				} else if value_type == "[]int" ||
					value_type == "[]int8" ||
					value_type == "[]int16" ||
					value_type == "[]int32" ||
					value_type == "[]int64" {

					if strings.Contains(value, "\n") {

						lines := strings.Split(value, "\n")
						slice := reflect.MakeSlice(field_type, len(lines), len(lines))

						for l := 0; l < len(lines); l++ {

							tmp, err := parseInt(lines[l], value_type[2:])

							if err == nil {
								slice.Index(l).SetInt(tmp)
							}

						}

						field.Set(slice)

					} else {

						tmp, err := parseInt(value, value_type[2:])

						if err == nil {

							slice := reflect.MakeSlice(field_type, 1, 1)
							slice.Index(0).SetInt(tmp)
							field.Set(slice)

						}

					}

				} else if value_type == "uint" ||
					value_type == "uint8" ||
					value_type == "uint16" ||
					value_type == "uint32" ||
					value_type == "uint64" {

					tmp, err := parseUint(value, value_type)

					if err == nil {
						field.SetUint(tmp)
					} else {
						return errors.New("whois: Cannot unmarshal string into Go struct field " + reflect_name + "." + field_name + " of type " + value_type)
					}

				} else if value_type == "[]uint" ||
					value_type == "[]uint8" ||
					value_type == "[]uint16" ||
					value_type == "[]uint32" ||
					value_type == "[]uint64" {

					if strings.Contains(value, "\n") {

						lines := strings.Split(value, "\n")
						slice := reflect.MakeSlice(field_type, len(lines), len(lines))

						for l := 0; l < len(lines); l++ {

							tmp, err := parseUint(lines[l], value_type[2:])

							if err == nil {
								slice.Index(l).SetUint(tmp)
							}

						}

						field.Set(slice)

					} else {

						tmp, err := parseUint(value, value_type[2:])

						if err == nil {

							slice := reflect.MakeSlice(field_type, 1, 1)
							slice.Index(0).SetUint(tmp)
							field.Set(slice)

						} else if value_type == "[]uint8" && strings.HasSuffix(err.Error(), "invalid syntax") {

							// value_type could have been []byte as it's an alias to []uint8
							bytes := []byte(value)
							slice := reflect.MakeSlice(field_type, len(bytes), len(bytes))

							for b := 0; b < len(bytes); b++ {
								slice.Index(b).SetUint(uint64(bytes[b]))
							}

							field.Set(slice)

						}

					}

				} else if value_type == "float32" || value_type == "float64" {

					tmp, err := parseFloat(value, value_type)

					if err == nil {
						field.SetFloat(tmp)
					} else {
						return errors.New("whois: Cannot unmarshal string into Go struct field " + reflect_name + "." + field_name + " of type " + value_type)
					}

				} else if value_type == "[]float32" || value_type == "[]float64" {

					if strings.Contains(value, "\n") {

						lines := strings.Split(value, "\n")
						slice := reflect.MakeSlice(field_type, len(lines), len(lines))

						for l := 0; l < len(lines); l++ {

							tmp, err := parseFloat(lines[l], value_type[2:])

							if err == nil {
								slice.Index(l).SetFloat(tmp)
							}

						}

						field.Set(slice)

					} else {

						tmp, err := parseFloat(value, value_type[2:])

						if err == nil {

							slice := reflect.MakeSlice(field_type, 1, 1)
							slice.Index(0).SetFloat(tmp)
							field.Set(slice)

						}

					}

				}

			}

		}

	}

	return nil

}
