package cmd

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func LaunchBuild(filename string) error {
	file, err := read(filename)
	if err != nil {
		return err
	}

	request, errChannel, err := toRequest(file)
	if err != nil {
		return err
	}

	client := http.Client{}
	res, err := client.Do(request)
	if err != nil {
		if errc := <-errChannel; errc != nil {
			return errors.New(fmt.Sprintf("multiple errors happend: %s %s", errc, err))
		} else {
			return err
		}
	}

	location := res.Header.Get("Location")
	fmt.Println(location)
	return nil
}

func toRequest(file *os.File) (*http.Request, chan error, error) {
	reader, writer := io.Pipe()
	newWriter := multipart.NewWriter(writer)
	errChannel := make(chan error, 1)
	go func() {
		defer file.Close()
		defer writer.Close()

		part, err := newWriter.CreateFormFile("file", filepath.Base(file.Name()))
		if err != nil {
			errChannel <- errors.New("unable to create multipart")
			return
		}

		if _, err := io.Copy(part, file); err != nil {
			errChannel <- err
			return
		}

		errChannel <- newWriter.Close()
	}()

	request, err := http.NewRequest("POST", "http://localhost:8088/files", reader)
	request.Header.Set("Content-Type", newWriter.FormDataContentType())
	return request, errChannel, err
}

func abs(filename string) (string, error) {
	if filepath.IsAbs(filename) {
		return filename, nil
	} else {
		if pwd, err := os.Getwd(); err == nil {
			return filepath.Join(pwd, filename), nil
		} else {
			return "", err
		}
	}
}

func read(filename string) (*os.File, error) {
	abs, err := abs(filename)

	if _, err := os.Stat(abs); os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("File %s not found", filename))
	}

	file, err := os.Open(abs)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Can not open file %s", abs))
	}
	return file, nil
}
