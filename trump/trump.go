package trump

import (
	"bytes"
	"image/png"

	"github.com/fogleman/gg"
	log "github.com/sirupsen/logrus"
)

func LoadOrder(imgPath string, fontPath string, fontSize float64) (*gg.Context, error) {
	img, err := gg.LoadImage(imgPath)

	if err != nil {
		return nil, err
	}

	dc := gg.NewContext(600, 338)
	dc.SetRGB(0, 0, 0)
	dc.DrawImage(img, 0, 0)

	if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
		log.Warningf("Unable to load fontface: %v", err)
	}

	return dc, nil
}

func Order(trump *gg.Context, message string) ([]byte, error) {
	trump.DrawStringWrapped(message, 500, 170, 0.5, 0.5, 160, 1.5, gg.AlignLeft)

	buf := new(bytes.Buffer)
	err := png.Encode(buf, trump.Image())

	return buf.Bytes(), err
}
