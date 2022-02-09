package dal

import (
	"errors"

	"github.com/yzx9/otodo/entity"
)

var files = make(map[string]entity.File)

func InsertFile(file entity.File) error {
	files[file.ID] = file
	return nil
}

func GetFile(id string) (entity.File, error) {
	file, ok := files[id]
	if !ok {
		return entity.File{}, errors.New("file not found")
	}

	return file, nil
}
