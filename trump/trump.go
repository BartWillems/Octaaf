package trump

import (
	"bytes"
	"image"
	"image/png"

	"github.com/fogleman/gg"
	log "github.com/sirupsen/logrus"
)

// Config is the configuration for the fonts & text alignments
type Config struct {
	FontPath   string  `toml:"font_path" env:"TRUMP_FONT_PATH"`
	FontSize   float64 `toml:"font_size" env:"TRUMP_FONT_SIZE"`
	LineHeight float64 `toml:"line_height" env:"TRUMP_LINE_HEIGHT"`
}

// LoadOrder returns a gg context canvas with the presidential order template
func LoadOrder(img image.Image, cfg *Config) gg.Context {
	dc := gg.NewContextForImage(img)

	dc.SetRGB(0, 0, 0)

	if err := dc.LoadFontFace(cfg.FontPath, cfg.FontSize); err != nil {
		log.Warningf("Unable to load fontface: %v", err)
	}

	return *dc
}

// Order returns a byte array image of the presidential order
func Order(img image.Image, cfg *Config, message string) ([]byte, error) {
	trump := LoadOrder(img, cfg)

	trump.DrawStringWrapped(message, 420, 150, 0, 0, 170, cfg.LineHeight, gg.AlignLeft)

	buf := new(bytes.Buffer)
	err := png.Encode(buf, trump.Image())

	return buf.Bytes(), err
}
