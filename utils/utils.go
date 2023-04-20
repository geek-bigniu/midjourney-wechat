package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func GetImageUrlData(imageUrl string) (bool, io.Reader) {
	resp, err := http.Get(imageUrl)
	if err != nil {
		fmt.Println(err)
		return false, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to download image. StatusCode: %d\n", resp.StatusCode)
		return false, nil
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return false, nil
	}

	imageReader := bytes.NewReader(imageData)
	return true, imageReader
}
