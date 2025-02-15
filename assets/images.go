package assets

import (
	"strings"

	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

type ImageID int16

var topID = 1
var stringToImageID = map[string]ImageID{}

var imageToEbiten = make([]*ebiten.Image, topID)

func GetImageID(s string) ImageID {
	if id, ok := stringToImageID[s]; ok {
		return id
	}
	return 0
}

func GetImage(id ImageID) *ebiten.Image {
	return imageToEbiten[id]
}

func init() {
	entries, err := fs.ReadDir(".")
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			if strings.HasSuffix(entry.Name(), ".png") {
				name := strings.TrimSuffix(entry.Name(), ".png")
				stringToImageID[name] = ImageID(topID)

				// Load it ebi-style
				img, err := fs.Open(entry.Name())
				if err != nil {
					panic(err)
				}
				defer img.Close()
				i, _, err := image.Decode(img)
				if err != nil {
					panic(err)
				}

				imageToEbiten = append(imageToEbiten, ebiten.NewImageFromImage(i))

				topID++
			}
		}
	}
}
