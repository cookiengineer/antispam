package whois

import "errors"
import "reflect"
import "strings"

func Unmarshal(buffer []byte, object any) error {

	blocks := Parse(buffer)

	if len(blocks) > 0 {

		valueof := reflect.ValueOf(object)

		if valueof.Kind() != reflect.Ptr {
			return errors.New("ptr required")
		}

		// Indirect Pointer
		valueof = valueof.Elem()
		valueof_kind := valueof.Kind()
		// valueof_name := valueof.String()
		valueof_type := valueof.Type()

		if valueof_kind == reflect.Slice {

			// slice->pointer->struct
			slice_element := reflect.TypeOf(object).Elem().Elem()

			var length = 0
			var filtered []map[string]string

			for b := 0; b < len(blocks); b++ {

				first_keys := strings.Split(slice_element.Field(0).Tag.Get("whois"), ",")

				var can_parse bool = false

				for f := 0; f < len(first_keys); f++ {

					first_key := first_keys[f]

					_, ok := blocks[b][first_key]

					if ok == true {
						can_parse = true
						break
					}

				}

				if can_parse == true {
					filtered = append(filtered, blocks[b])
					length++
				}

			}

			if length > 0 && len(filtered) > 0 {

				slice := reflect.MakeSlice(valueof_type, 0, length)

				for f := 0; f < len(filtered); f++ {

					block := filtered[f]
					element := reflect.New(slice_element).Elem()

					for key, value := range block {

						err := setFieldValue(&element, key, value)

						if err != nil {
							return err
						}

					}

					slice = reflect.Append(slice, element)

				}

				valueof.Set(slice)

			} else {
				return errors.New("buffer does not contain any sections matching first field of structs")
			}

		} else if valueof_kind == reflect.Struct {

			for b := 0; b < len(blocks); b++ {

				block := blocks[b]

				for key, value := range block {

					err := setFieldValue(&valueof, key, value)

					if err != nil {
						return err
					}

				}

			}

		} else {

			return errors.New("ptr to slice of structs or ptr to struct required")

		}

	}

	return nil

}
