package error

import (
	"fmt"
	"strings"
	"tud/cmd/internal/file"
)

type Reporter struct {
	Errors []ErrorReport
	Source *[]byte
}

type ErrorReport struct {
	Pos        file.Position
	Message    string
	Where      string
	SourceLine string
}

func NewErrorReporter(source *[]byte) Reporter {
	return Reporter{
		Errors: make([]ErrorReport, 0),
		Source: source,
	}
}

// Copies struct so i doesnt hold state change
func (r *Reporter) NewError(pos file.Position, msg string) {
	src_line := r.findErrorLine(&pos)

	error_report := ErrorReport{
		Pos:        pos,
		Message:    msg,
		SourceLine: src_line,
		Where:      string(r.Source[pos.Start:pos.Offset]),
	}

	r.Errors = append(r.Errors, error_report)
}

func (r *Reporter) GetErrorsText() []string {
	for _, e := range r.Errors {
		fmt.Printf("Error: %s [%d:%d] at %s\n", e.Message, e.Pos.Line, e.Pos.LineCol+1, e.Where)
		fmt.Println(e.SourceLine)
		fmt.Printf("%s^\n", strings.Repeat(" ", e.Pos.LineCol))
	}
	return make([]string, 0)
}

func (r *Reporter) HasErrors() bool {
	return len(r.Errors) > 0
}

func (r *Reporter) AppendError(err ErrorReport) []ErrorReport {
	r.Errors = append(r.Errors, err)
	return r.Errors
}

func (r *Reporter) findErrorLine(pos *file.Position) string {
	start := pos.Offset
	source := r.Source

	for start > 0 && source[start-1] != '\n' {
		start--
	}

	end := pos.Offset
	for end < len(source) && source[end] != '\n' {
		end++
	}

	line_text := string(source[start:end])
	pos.LineCol = (pos.Offset - start) - 1
	pos.Start = start

	return line_text
}
