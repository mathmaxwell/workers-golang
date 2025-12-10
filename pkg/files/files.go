package files

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"
)

// change fileName, here fileName is employees
func SaveFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	defer file.Close()
	os.MkdirAll("uploads/employees", 0755)
	filename := fmt.Sprintf("uploads/employees/%d-%s", time.Now().Unix(), header.Filename)
	out, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}
	return filename, nil
}
