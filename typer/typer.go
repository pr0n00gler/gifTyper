package typer

import (
	"errors"
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

const defaultFontSize = 32
const defaultFontFile = "Roboto-Regular.ttf"
const defaultFrameWidth = 500
const defaultFrameHeight = 500
const defaultMaxLineSize = 20
const defaultMaxLineCount = 5
const defaultDelay = 30

type Typer struct {
	maxLineSize    int
	maxLinesCount  int
	frameW, frameH int
	font           font.Face
	fontPath       string
	fontSize       int
	delay          int
}

func InitGenerator() (*Typer, error) {
	var (
		err error
	)

	generator := &Typer{
		maxLinesCount: defaultMaxLineCount,
		maxLineSize:   defaultMaxLineSize,
		frameW:        defaultFrameWidth,
		frameH:        defaultFrameHeight,
		delay:         defaultDelay,
		fontPath:      defaultFontFile,
	}
	if generator.font, err = gg.LoadFontFace(defaultFontFile, defaultFontSize); err != nil {
		return nil, err
	}
	return generator, nil
}

func (t *Typer) SetFontSize(fontSize int) error {
	err := t.SetFont(t.fontPath, fontSize)
	if err != nil {
		return err
	}
	return nil
}

func (t *Typer) SetFont(fontFile string, fontSize int) error {
	var err error
	t.font, err = gg.LoadFontFace(fontFile, float64(fontSize))
	if err != nil {
		return err
	}
	return nil
}

func (t *Typer) SetDelay(delay int) error {
	if delay < 1 {
		return errors.New("Incorrect delay")
	}
	t.delay = delay
	return nil
}

func (t *Typer) SetMaxLinesCount(maxLinesCount int) error {
	if maxLinesCount < 1 {
		return errors.New("Incorrect lines count")
	}
	t.maxLinesCount = maxLinesCount
	return nil
}

func (t *Typer) setFrameWidth(width int) error {
	if width < 1 {
		return errors.New("Incorrect lines count")
	}
	t.frameW = width
	return nil
}

func (t *Typer) setFrameHeight(height int) error {
	if height < 1 {
		return errors.New("Incorrect lines count")
	}
	t.frameH = height
	return nil
}

func (t *Typer) countSpaceWidth() int {
	rect, _ := font.BoundString(t.font, "W")
	rectSizeWithoutSpace := rect.Max.X.Round()
	rect, _ = font.BoundString(t.font, " W")
	rectSizeWithSpace := rect.Max.X.Round()

	spaceLength := rectSizeWithSpace - rectSizeWithoutSpace
	println(spaceLength)
	return spaceLength
}

func (t *Typer) countStringWidth(s string) int {
	rect, _ := font.BoundString(t.font, s)
	width := rect.Max.X.Round()
	return width
}

func (t *Typer) countMaxLineSize(text string) int {
	var maxLineSize = 0
	widthLimit := t.frameW
	spaceLength := 1
	spaceWidth := t.countSpaceWidth()
	var currentSymbolsCount = 0
	var currentSymbolsWidth = 0
	words := strings.Split(text, " ")
	for _, word := range words {
		wordWidth := t.countStringWidth(word)
		if currentSymbolsWidth+wordWidth < widthLimit {
			currentSymbolsWidth += wordWidth + spaceWidth
			currentSymbolsCount += len(word) + spaceLength
		} else {
			maxLineSize = currentSymbolsCount
			currentSymbolsWidth = wordWidth + spaceWidth
			currentSymbolsCount = len(word) + spaceLength
		}
	}

	if maxLineSize == 0 {
		maxLineSize = currentSymbolsCount
	}

	return maxLineSize
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
	line = t.checkSpacesAfterPunctuationMarks(line)
	t.maxLineSize = t.countMaxLineSize(line)
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
		outGif.Delay = append(outGif.Delay, t.delay)
	}
	return outGif, nil
}

func (t *Typer) divideTextOnLines(text string) ([]string, int, error) {
	space := " "
	framesCount := 0

	words := strings.Split(text, space)
	lines := make([]string, 0)

	var line strings.Builder
	line.Grow(t.maxLineSize)
	for _, word := range words {
		if (line.Len() + len(word)) < t.maxLineSize {
			line.WriteString(word + space)
		} else {
			lines = append(lines, line.String())
			framesCount += line.Len()
			line.Reset()
			line.Grow(t.maxLineSize)
			line.WriteString(word + space)
		}
	}

	if line.Len() != 0 {
		framesCount += line.Len()
		lines = append(lines, line.String())
		line.Reset()
	}

	return lines, framesCount, nil
}

func (t *Typer) checkSpacesAfterPunctuationMarks(text string) string {
	var space byte = ' '
	punctuationMarks := []byte{',', '.', '!', '?', ':', ';', '-'}
	for index, _ := range text {
		if index == len(text)-1 {
			continue
		}
		if t.contains(punctuationMarks, text[index]) &&
			!t.contains(punctuationMarks, text[index+1]) &&
			text[index+1] != space {
			text = strings.ReplaceAll(text, string(text[index]), string(text[index])+string(space))
		}
	}
	return text
}

func (t *Typer) contains(s []byte, e byte) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
