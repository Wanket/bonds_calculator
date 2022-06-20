package moex

import (
	"encoding/csv"
	"io"
)

type CsvReader struct {
	*csv.Reader
}

func NewReader(r io.Reader) CsvReader {
	reader := CsvReader{csv.NewReader(r)}
	reader.Comma = ';'
	reader.FieldsPerRecord = -1
	reader.ReuseRecord = true

	return reader
}
