package insights

import "strconv"
import "strings"

func splitDiffLines(buffer []byte) []string {

	value := string(buffer)

	if value == "" {
		return make([]string, 0)
	}

	lines := strings.Split(value, "\n")

	if strings.HasSuffix(value, "\n") == true {
		lines = lines[:len(lines)-1]
	}

	return lines

}

func lineRange(start int, count int) string {

	if count == 0 {
		return strconv.Itoa(start) + ",0"
	}

	return strconv.Itoa(start+1) + "," + strconv.Itoa(count)

}

// UnifiedDiff returns a small unified diff with three lines of context. The
// dataset edits made by the helper affect one contiguous region per file, so
// a single hunk is sufficient.
func UnifiedDiff(path string, old_buffer []byte, new_buffer []byte) string {

	old_lines := splitDiffLines(old_buffer)
	new_lines := splitDiffLines(new_buffer)

	prefix := 0

	for prefix < len(old_lines) && prefix < len(new_lines) {

		if old_lines[prefix] != new_lines[prefix] {
			break
		}

		prefix++

	}

	suffix := 0

	for suffix < (len(old_lines)-prefix) && suffix < (len(new_lines)-prefix) {

		if old_lines[len(old_lines)-1-suffix] != new_lines[len(new_lines)-1-suffix] {
			break
		}

		suffix++

	}

	if prefix == len(old_lines) && prefix == len(new_lines) {
		return ""
	}

	context := 3
	start := prefix - context

	if start < 0 {
		start = 0
	}

	old_end := len(old_lines) - suffix

	if old_end < prefix {
		old_end = prefix
	}

	new_end := len(new_lines) - suffix

	if new_end < prefix {
		new_end = prefix
	}

	old_context_end := old_end + context

	if old_context_end > len(old_lines) {
		old_context_end = len(old_lines)
	}

	new_context_end := new_end + context

	if new_context_end > len(new_lines) {
		new_context_end = len(new_lines)
	}

	if old_context_end < old_end {
		old_context_end = old_end
	}

	if new_context_end < new_end {
		new_context_end = new_end
	}

	builder := strings.Builder{}

	builder.WriteString("--- a/" + path + "\n")
	builder.WriteString("+++ b/" + path + "\n")

	old_count := old_context_end - start
	new_count := new_context_end - start

	builder.WriteString("@@ -" + lineRange(start, old_count) + " +" + lineRange(start, new_count) + " @@\n")

	for line := start; line < prefix; line++ {
		builder.WriteString(" " + old_lines[line] + "\n")
	}

	for line := prefix; line < old_end; line++ {
		builder.WriteString("-" + old_lines[line] + "\n")
	}

	for line := prefix; line < new_end; line++ {
		builder.WriteString("+" + new_lines[line] + "\n")
	}

	for line := old_end; line < old_context_end; line++ {
		builder.WriteString(" " + old_lines[line] + "\n")
	}

	return builder.String()

}
