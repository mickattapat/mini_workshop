package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/trace"
	"sync"
	"time"
)

type Photos []struct {
	AlbumID      int    `json:"albumId"`
	ID           int    `json:"id"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnailUrl"`
}

type Image struct {
	filePath string
	img      []byte
}

const MAX_DOWNLOAD = 1000 // photo

func main() {
	defer func() {
		fmt.Println("Main program exit successfully")
	}()
	log.SetFlags(log.Ltime)

	dir := "myDownLoadImages_" + time.Now().Format("14_02_2022") // set folder name
	if _, err := os.Stat(dir); err != nil { // make folder
		os.Mkdir(dir, 0777) // or ModeDir 
	}
	f, err := os.Create(dir + ".trace.log")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	trace.Start(f)
	defer trace.Stop()

	photos := Photos{}
	err = getJson("https://jsonplaceholder.typicode.com/photos", &photos)
	fmt.Println(err)
	fmt.Println(len(photos[:MAX_DOWNLOAD]))

	chImg := make(chan Image, len(photos[:MAX_DOWNLOAD])) // buffered channels
	token := make(chan struct{}, 20)                      // limited download
	counter := sync.WaitGroup{}
	for _, v := range photos[:MAX_DOWNLOAD] {
		photo := v
		counter.Add(1)
		go func() {
			defer counter.Done()
			if photo.ID > 2500 {
				photo.ThumbnailURL = "http://abc.jpg"
			}
			// allow maximum download 20
			// take a token
			token <- struct{}{}
			img, err := downloadImage(photo.ThumbnailURL) // use token to work
			// release token
			<-token
			if err != nil {
				log.Println(err)
				return
			}
			format, err := deCodeImages(img)
			if err != nil {
				log.Fatal(err)
			}
			absoluteFileName := filepath.Join(dir, fmt.Sprintf("%d.%s", photo.ID, format)) // TODO analyze file type
			chImg <- Image{filePath: absoluteFileName, img: img}
		}()
	}
	go func() {
		counter.Wait()
		close(chImg)
	}()
	for v := range chImg {
		err := saveImages(v.filePath, v.img)
		if err != nil {
			log.Println(err)
		}
	}
}

func saveImages(fileName string, img []byte) error {
	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("saveImage - cannot create file : %v", err)
	}
	defer f.Close()
	//f <-- img
	_, err = io.Copy(f, bytes.NewReader(img))
	if err != nil {
		return fmt.Errorf("saveImage - cannot save file : %v", err)
	}
	log.Printf("save : %v\n", fileName)
	return nil
}
func deCodeImages(img []byte) (string, error) {
	_, format, err := image.Decode(bytes.NewReader(img))
	return format, err
}

func downloadImage(url string) ([]byte, error) {
	errMsg := func(err error) error {
		return fmt.Errorf("downloadImg : %v", err)
	}
	res, err := http.Get(url)
	if err != nil {
		return nil, errMsg(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errMsg(err)
	}
	return body, nil
}

func getJson(url string, structType interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	switch v := structType.(type) {
	case *Photos:
		decoder := json.NewDecoder(res.Body)
		photos := structType.(*Photos)
		decoder.Decode(photos)
		return nil
	default:
		return fmt.Errorf("getJson : not support typr %v", v)
	}
}
