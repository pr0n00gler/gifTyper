package typer

import (
	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"strings"
)

var palette = color.Palette{
	image.Transparent,
	image.Black,
	image.White,
	color.RGBA{G: 255, A: 255},
	color.RGBA{G: 100, A: 255},
}

type Typer struct {
	maxLineSize    int
	maxLinesCount  int
	frameW, frameH int
	font           font.Face
}

func InitGenerator(maxLineSize, maxLinesCount, frameW, frameH int, fontFile string) (*Typer, error) {
	var (
		err error
	)

	generator := &Typer{
		maxLinesCount: maxLinesCount,
		maxLineSize:   maxLineSize,
		frameW:        frameW,
		frameH:        frameH,
	}
	if generator.font, err = gg.LoadFontFace(fontFile, 32); err != nil {
		return nil, err
	}
	return generator, nil
}

func (t *Typer) drawFrame(line string, x, y float64) image.Image {
	dc := gg.NewContext(t.frameW, t.frameH)
	dc.SetRGBA(1, 1, 1, 0)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	dc.SetFontFace(t.font)
	dc.DrawString(line, x, y)
	dc.Clip()
	return dc.Image()
}

func (t *Typer) drawBackground() image.Image {
	dc := gg.NewContext(t.frameW, t.frameH)
	dc.SetRGBA(1, 1, 1, 1)
	dc.Clear()
	return dc.Image()
}

func (t *Typer) drawFrames(lines []string, framesCount int) []image.Image {
	frames := make([]image.Image, 0, framesCount)
	frames = append(frames, t.drawBackground())

	var typedLine strings.Builder
	for i, line := range lines {
		typedLine.Grow(len(line))
		for _, symbol := range line {
			typedLine.WriteRune(symbol)
			frame := t.drawFrame(typedLine.String(), 0, float64(i+1)*32)
			frames = append(frames, frame)
		}
		typedLine.Reset()
	}
	return frames
}

func (t *Typer) GenerateGIF(line string) (*gif.GIF, error) {
	lines, framesCount, err := t.divideTextOnLines(line)
	if err != nil {
		return nil, err
	}
	frames := t.drawFrames(lines, framesCount)
	outGif := &gif.GIF{}
	for _, frame := range frames {
		bounds := frame.Bounds()
		palettedImage := image.NewPaletted(bounds, palette)
		draw.Draw(palettedImage, palettedImage.Rect, frame, bounds.Min, draw.Src)

		outGif.Image = append(outGif.Image, palettedImage)
		outGif.Delay = append(outGif.Delay, 30)
	}
	return outGif, nil
}

func (t *Typer) divideTextOnLines(text string) ([]string, int, error) {
	framesCount := 0

	lines := make([]string, 0)
	var line strings.Builder
	line.Grow(t.maxLineSize)
	for _, character := range text {
		if _, err := line.WriteRune(character); err != nil {
			return nil, framesCount, err
		}
		if line.Len() >= t.maxLineSize {
			lines = append(lines, line.String())
			line.Reset()
			line.Grow(t.maxLineSize)
		}
		framesCount++
	}
	if line.Len() > 0 {
		lines = append(lines, line.String())
	}
	return lines, framesCount, nil
}
