package trump

import (
	"bytes"
	"image"
	"image/png"

	"github.com/fogleman/gg"
	log "github.com/sirupsen/logrus"
)

// TrumpConfig is the configuration for the fonts & text alignments
type TrumpConfig struct {
	FontPath   string  `toml:"font_path"`
	FontSize   float64 `toml:"font_size"`
	LineHeight float64 `toml:"line_height"`
}

func LoadOrder(img image.Image, cfg *TrumpConfig) gg.Context {
	dc := gg.NewContextForImage(img)

	dc.SetRGB(0, 0, 0)

	if err := dc.LoadFontFace(cfg.FontPath, cfg.FontSize); err != nil {
		log.Warningf("Unable to load fontface: %v", err)
	}

	return *dc
}

func Order(img image.Image, cfg *TrumpConfig, message string) ([]byte, error) {
	trump := LoadOrder(img, cfg)

	trump.DrawStringWrapped(message, 420, 150, 0, 0, 170, cfg.LineHeight, gg.AlignLeft)

	buf := new(bytes.Buffer)
	err := png.Encode(buf, trump.Image())

	return buf.Bytes(), err
}
