package files

import (
	"crypto/sha256"
	"encoding/gob"
	errors2 "errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"telegramBot/bot/errors"
	"telegramBot/storage"
)

type FileStorage struct {
	basePath string
}

func NewFileStorage(basePath string) *FileStorage {
	return &FileStorage{basePath: basePath}
}

const defaultPerm = 0774

func (fileStorage *FileStorage) SaveLocation(username string, location string) (err error) {
	filePath := filepath.Join(fileStorage.basePath, username)
	if err := os.MkdirAll(filePath, defaultPerm); err != nil {
		return err
	}

	name, err := hash(username, location)
	if err != nil {
		return err
	}

	filePath = filepath.Join(filePath, name)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(location); err != nil {
		return err
	}

	return nil
}

func (fileStorage *FileStorage) Location(username string) (location *string, err error) {
	defer func() {
		err = errors.WrapIfError("can't retrieve location", err)
	}()

	path := filepath.Join(fileStorage.basePath, username)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedLocation
	}

	return fileStorage.decodeLocation(filepath.Join(path, files[0].Name()))
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

func hash(username, location string) (string, error) {
	hash := sha256.New()

	if _, err := io.WriteString(hash, location); err != nil {
		return "", errors.Wrap("unable to calculate hash", err)
	}

	if _, err := io.WriteString(hash, username); err != nil {
		return "", errors.Wrap("unable to calculate hash", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func (fileStorage *FileStorage) Save(page *storage.Page) (err error) {
	defer func() {
		err = errors.WrapIfError("unable to save:", err)
	}()

	filePath := filepath.Join(fileStorage.basePath, page.Username)
	if err := os.MkdirAll(filePath, defaultPerm); err != nil {
		return err
	}

	name, err := fileName(page)
	if err != nil {
		return err
	}

	filePath = filepath.Join(filePath, name)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func (fileStorage *FileStorage) PickRandom(username string) (page *storage.Page, err error) {
	defer func() {
		err = errors.WrapIfError("can't pick random page", err)
	}()

	path := filepath.Join(fileStorage.basePath, username)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	n := rand.Intn(len(files))
	file := files[n]

	return fileStorage.decodePage(filepath.Join(path, file.Name()))
}

func (fileStorage *FileStorage) Remove(page *storage.Page) error {
	filename, err := fileName(page)
	if err != nil {
		return err
	}

	path := filepath.Join(fileStorage.basePath, page.Username, filename)

	if err := os.Remove(path); err != nil {
		message := fmt.Sprintf("can't remove file at %s", path)
		return errors.Wrap(message, err)
	}

	return nil
}

func (fileStorage *FileStorage) IsExists(page *storage.Page) (bool, error) {
	filename, err := fileName(page)
	if err != nil {
		return false, err
	}

	path := filepath.Join(fileStorage.basePath, page.Username, filename)

	switch _, err = os.Stat(path); {
	case errors2.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		message := fmt.Sprintf("can't check if file %s exists", path)
		return false, errors.Wrap(message, err)
	}

	return true, nil
}

func (fileStorage *FileStorage) decodePage(filepath string) (*storage.Page, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, errors.Wrap("can't decode page", err)
	}
	defer func() { _ = file.Close() }()

	var page storage.Page

	if err := gob.NewDecoder(file).Decode(&page); err != nil {
		return nil, errors.Wrap("can't decode page", err)
	}

	return &page, nil
}

func fileName(page *storage.Page) (string, error) {
	return page.Hash()
}
