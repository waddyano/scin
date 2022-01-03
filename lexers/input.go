package lexers

import "os"

type Input struct {
	filename      string
	file          *os.File
	data          []byte
	cursor        int
	marker        int
	token         int
	limit         int
	line          int
	state         int
	eof           bool
	bol           bool
	raw_str_delim []byte
}

func fill(in *Input) int {
	return 1
}

func peek(in *Input) byte {
	if in.cursor >= in.limit {
		return 0
	}
	return in.data[in.cursor]
}
