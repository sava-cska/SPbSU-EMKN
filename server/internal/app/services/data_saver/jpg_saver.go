package data_saver

import (
	"fmt"
	"github.com/google/uuid"
	"os"
)

type JpgSaver struct{}

func (saver *JpgSaver) HardSave(content []byte) (string, error) {
	return save(content)
}

func save(content []byte) (string, error) {
	filename := fmt.Sprintf("%s.jpg", uuid.New())
	fo, err := os.Create("data/" + filename)
	if err != nil {
		return "", err
	}
	// close fo on exit and check for its returned error
	defer func() {
		_ = fo.Close()
	}()
	if _, err := fo.Write(content); err != nil {
		return "", err
	}
	if err := fo.Sync(); err != nil {
		return "", err
	}
	return filename, nil
}
