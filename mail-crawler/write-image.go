package main

import (
	"fmt"
	"os"
)

func writeImageToFs(image []byte, filename string) {
	f, err := os.Create(getConfigInstance().Fs.PhotosStorePath + string(os.PathSeparator) + filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = f.Write(image)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
