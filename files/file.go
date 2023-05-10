package files

import (
	"encoding/gob"
	errors2 "errors"
	"os"
	"path/filepath"
	"telegramBot/bot/errors"
)

var ErrNoSavedLocation = errors2.New("no saved location")

type FileStorage struct {
	basePath string
}

// NewFileStorage returns a new instance of the FileStorage
func NewFileStorage(basePath string) *FileStorage {
	return &FileStorage{basePath: basePath}
}

const defaultPerm = 0774

// SaveLocation saves user's location to the file /location_files/username
func (fileStorage *FileStorage) SaveLocation(username string, location string) (err error) {
	filePath := fileStorage.basePath
	if err := os.MkdirAll(filePath, defaultPerm); err != nil {
		return err
	}

	name := username

	filePath = filepath.Join(filePath, name)

	file, err := getFile(filePath)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(location); err != nil {
		return err
	}

	return nil
}

func getFile(filePath string) (*os.File, error) {
	if _, err := os.Stat(filePath); errors2.Is(err, os.ErrNotExist) {
		file, err := os.Create(filePath)
		if err != nil {
			return nil, err
		}
		return file, nil
	} else {
		file, err := os.OpenFile(filePath, os.O_WRONLY, defaultPerm)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
}

// Location returns a location from /location_files/username
func (fileStorage *FileStorage) Location(username string) (location *string, err error) {
	defer func() {
		err = errors.WrapIfError("can't retrieve location", err)
	}()

	path := filepath.Join(fileStorage.basePath, username)

	return fileStorage.decodeLocation(path)
}

func (fileStorage *FileStorage) decodeLocation(filepath string) (*string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, errors.Wrap("can't decode page", err)
	}
	defer func() { _ = file.Close() }()

	var location *string

	if err := gob.NewDecoder(file).Decode(&location); err != nil {
		return nil, errors.Wrap("can't decode location", err)
	}

	return location, nil
}
