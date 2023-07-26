package actions

import "postfix-spamdb/structs"
import "encoding/json"
import "fmt"
import "os"

func ReadJSON(filepath string) []structs.Entry {

	var entries []structs.Entry

	buffer, err1 := os.ReadFile(filepath)

	if err1 == nil && len(buffer) > 2 {

		err2 := json.Unmarshal(buffer, &entries)

		if err2 != nil {
			fmt.Println("Could not read \"" + filepath + "\"!")
			fmt.Println(err2)
		}

	}

	return entries

}
