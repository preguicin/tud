package file

import (
	"io"
	"os"
)

type Position struct {
	Line  int
	Col   int
	Start int
}

func NewPosition() Position {
	return Position{Line: 1, Col: 0, Start: 0}
}

type File interface {
	GetBytes() []byte
	Close() error
	Pos() *Position
}

type InMemFile struct {
	data []byte
	Position
}

func (inf *InMemFile) GetBytes() []byte {
	return inf.data
}
func (inf *InMemFile) Close() error {
	inf.data = nil
	return nil
}

func (inf *InMemFile) Pos() *Position { return &inf.Position }

func NewInMemFile(data []byte) *InMemFile {
	return &InMemFile{data: data, Position: NewPosition()}
}

type DiskFile struct {
	data []byte
	path string
	f    *os.File
	Position
}

func (df *DiskFile) GetBytes() []byte {
	return df.data
}

func (df *DiskFile) Close() error {
	err := df.f.Close()
	return err
}

func (df *DiskFile) Pos() *Position { return &df.Position }

func NewDiskFile(file *os.File) (*DiskFile, error) {

	data, err := io.ReadAll(file)
	if err != nil {
		file.Close()
		return nil, err
	}

	return &DiskFile{
		data:     data,
		f:        file,
		Position: NewPosition(),
	}, nil
}
