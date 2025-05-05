package helpers

func PosToLineCol(input []byte, pos int) (line int, col int) {
	if pos > len(input) {
		pos = len(input)
	}

	line = 1
	lastNewline := -1

	for i := 0; i < pos; i++ {
		if input[i] == '\n' {
			line++
			lastNewline = i
		}
	}
	col = pos - lastNewline
	return line, col
}
