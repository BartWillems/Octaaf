package trump

import (
	"bytes"
	"image"
	"image/png"

	"github.com/fogleman/gg"
	log "github.com/sirupsen/logrus"
)

func LoadOrder(img image.Image, fontPath string, fontSize float64) gg.Context {
	dc := gg.NewContextForImage(img)

	dc.SetRGB(0, 0, 0)

	if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
		log.Warningf("Unable to load fontface: %v", err)
	}

	return *dc
}

func Order(img image.Image, fontPath string, fontSize float64, message string) ([]byte, error) {
	trump := LoadOrder(img, fontPath, fontSize)

	trump.DrawStringWrapped(message, 500, 170, 0.5, 0.5, 160, 1.5, gg.AlignLeft)

	buf := new(bytes.Buffer)
	err := png.Encode(buf, trump.Image())

	return buf.Bytes(), err
}
