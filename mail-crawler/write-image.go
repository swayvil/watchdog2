package main

import (
	"log"
	"os"
)

// Call once, after the connection to the database is established
// To create, if they don't exist already, the folders hierarchy where the snapshots will be stored
func createFolders() {
	cameras := selectCameras()

	for i := 0; i < len(cameras); i++ {
		folderPath := getConfigInstance().Fs.PhotosStorePath + string(os.PathSeparator) + cameras[i] + string(os.PathSeparator)
		err1 := os.MkdirAll(folderPath+"small", os.ModePerm)
		if err1 != nil {
			log.Fatal(err1)
		}
		err2 := os.MkdirAll(folderPath+"big", os.ModePerm)
		if err2 != nil {
			log.Fatal(err2)
		}
		log.Printf("Initialize folder: %s\n", folderPath+"small")
		log.Printf("Initialize folder: %s\n", folderPath+"big")
	}
}

func writeImageToFs(image []byte, filename string) {
	f, err := os.Create(getConfigInstance().Fs.PhotosStorePath + string(os.PathSeparator) + filename)
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = f.Write(image)
	if err != nil {
		log.Fatal(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
}
