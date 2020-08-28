// volu-img is a tool for generating resized product jpg's specifically for volusion store ftp upload

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path"
	"sync"

	"github.com/nfnt/resize"
)

var mux sync.Mutex

// image size structure
type size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// product data
type product struct {
	Sku     string   `json:"sku"`
	Img     string   `json:"img"`
	ImgAlts []string `json:"img_alts"`
	Large   size     `json:"large"`
	Medium  size     `json:"medium"`
	Small   size     `json:"small"`
	Tiny    size     `json:"tiny"`
	Thumb   size     `json:"thumb"`
}

func main() {

	// Need to provide json input and path to images, default image path is cwd
	var jsonFile = flag.String("i", "", "json input file path")
	var imgDir = flag.String("p", ".", "path to images for resizing")

	flag.Parse()

	// json file required
	if *jsonFile == "" {
		fmt.Println("need to specify json input file path using -i")
		os.Exit(1)
	}

	var products []product

	inFile, err := os.Open(*jsonFile)
	if err != nil {
		panic(err)
	}

	jsonParser := json.NewDecoder(inFile)
	err = jsonParser.Decode(&products)
	if err != nil {
		panic(err)
	}

	for _, p := range products {
		genImages(p, *imgDir)
	}

}

func genImages(p product, imgDir string) {

	imgPath := path.Join(imgDir, p.Img)

	file, err := os.Open(imgPath)
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	var wg sync.WaitGroup

	// see https://helpcenter.volusion.com/en/articles/1773795-product-image-file-names
	// Create the "main" image sizes
	wg.Add(5)
	createImage(img, fmt.Sprintf("%s-2.jpg", p.Sku), imgDir, p.Large, &wg)
	createImage(img, fmt.Sprintf("%s-1.jpg", p.Sku), imgDir, p.Small, &wg)
	createImage(img, fmt.Sprintf("%s-0.jpg", p.Sku), imgDir, p.Tiny, &wg)
	createImage(img, fmt.Sprintf("%s-2S.jpg", p.Sku), imgDir, p.Thumb, &wg)
	createImage(img, fmt.Sprintf("%s-2T.jpg", p.Sku), imgDir, p.Medium, &wg)

	// create the "alternative" image sizes, extension starts at 3+, e.g
	for i, imgAlt := range p.ImgAlts {

		altImgPath := path.Join(imgDir, imgAlt)

		file, err := os.Open(altImgPath)
		if err != nil {
			log.Fatal(err)
		}

		// decode jpeg into image.Image
		img, err := jpeg.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()

		// tiny and small sizes not created for alternative images
		wg.Add(3)
		createImage(img, fmt.Sprintf("%s-%d.jpg", p.Sku, i+3), imgDir, p.Large, &wg)
		createImage(img, fmt.Sprintf("%s-%dS.jpg", p.Sku, i+3), imgDir, p.Thumb, &wg)
		createImage(img, fmt.Sprintf("%s-%dT.jpg", p.Sku, i+3), imgDir, p.Medium, &wg)
	}
	wg.Wait()
}

// detSize determines the maximum dimension size to limit, if both height and width provided, only scales the largest dimension.
// this ensures the image will not be skewed and preserves aspect ratio
func detSize(s size, b image.Rectangle) size {
	switch {

	case s.Height > 0 && s.Width == 0:
		// max height
		return s

	case s.Width > 0 && s.Height == 0:
		// max width
		return s

	default:
		// both
		if b.Dy() >= b.Dx() {
			s.Width = 0
			return s
		}
		s.Height = 0
		return s

	}
}

func createImage(img image.Image, fileName, imgDir string, s size, wg *sync.WaitGroup) {

	defer wg.Done()

	bnds := img.Bounds()

	s = detSize(s, bnds)

	m := resize.Resize(uint(s.Width), uint(s.Height), img, resize.Lanczos3)

	// creates new directory "volu" in images directory, if does not exist
	mux.Lock()
	newDir := path.Join(imgDir, "volu")
	if _, err := os.Stat(newDir); os.IsNotExist(err) {
		os.Mkdir(newDir, 0755)
	}

	out, err := os.Create(path.Join(newDir, fileName))
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	mux.Unlock()

	// write new image to file
	jpeg.Encode(out, m, nil)
}
