package file

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/restaurant/internal/pkg/config"
	"io"
	"mime/multipart"
	"os"
)

func UploadMultiple(files []*multipart.FileHeader, directory string) (string, []string, error) {
	var (
		created       []string
		pqStringArray = "{"
	)
	if len(files) < 0 {
		return "", nil, fmt.Errorf("empty parts")
	}

	for k, v := range files {

		link, name, err := touch(v, directory)
		if err != nil {
			return "", created, err
		}
		created = append(created, name)
		if k == 0 {
			pqStringArray += fmt.Sprintf(`"%s"`, link)
		} else {
			pqStringArray += fmt.Sprintf(` ,"%s"`, link)
		}
	}
	pqStringArray += "}"

	return pqStringArray, created, nil
}

func UploadSingle(part *multipart.FileHeader, directory string) (string, string, error) {
	return touch(part, directory)
}

func touch(part *multipart.FileHeader, directory string) (string, string, error) {
	var (
		contentType  = part.Header.Get("Content-Type")
		relativePath = fmt.Sprintf("media/%s/", directory)
	)

	cfg := config.NewConfig()
	base := cfg.BaseDestination

	reader, err := part.Open()
	if err != nil {
		return "", "", err
	}

	ext, ok := extensions[contentType]
	if !ok {
		return "", "", fmt.Errorf("content-type not supported, src:[%s]", contentType)
	}

	name := fmt.Sprintf("%s.%s", uuid.NewString(), ext)
	if _, err := os.Stat(base + relativePath); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(base+relativePath, os.ModePerm); err != nil {
			return "", "", fmt.Errorf("%v dest:[%s]", err, relativePath)
		}
	}
	destination := relativePath + name

	writer, err := os.Create(base + destination)
	if err != nil {
		return "", "", err
	}

	_, err = io.Copy(writer, reader)
	if err != nil {
		return "", "", err
	}

	return destination, base + destination, nil
}
