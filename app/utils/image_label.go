package utils

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
)

type Options struct {
	Dpi       float64
	Size      float64
	Fontfile  string
	TextColor color.RGBA
}

func imageToRGBA(src image.Image) *image.RGBA {
	// No conversion needed if image is an *image.RGBA.
	if dst, ok := src.(*image.RGBA); ok {
		return dst
	}

	// Use the image/draw package to convert to *image.RGBA.
	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
	return dst
}

func writePaymentThrough(d *font.Drawer, str string) {
	d.Dot = fixed.Point26_6{
		X: fixed.I(830),
		Y: fixed.I(525),
	}

	d.DrawString(str)
}

func writeBuyAmount(d *font.Drawer, str string) {
	d.Dot = fixed.Point26_6{
		X: fixed.I(830),
		Y: fixed.I(605),
	}

	d.DrawString(str)
}

func writePaymentSum(d *font.Drawer, str string) {
	d.Dot = fixed.Point26_6{
		X: fixed.I(830),
		Y: fixed.I(685),
	}

	d.DrawString(str)
}

func writePaymentAddress(d *font.Drawer, str string) {
	d.Dot = fixed.Point26_6{
		X: fixed.I(830),
		Y: fixed.I(765),
	}

	d.DrawString(str)
}

func LabelImage(openFileName string, outputFileName string, through string, buyAmt string, paymentSum string, paymentAddr string, options *Options) {
	flag.Parse()

	fontBytes, err := ioutil.ReadFile(options.Fontfile)
	if err != nil {
		log.Println(err)
		return
	}
	f, err := truetype.Parse(fontBytes)

	img, err := os.Open(openFileName)

	if err != nil {
		log.Fatalf(err.Error())
	}
	fg := image.NewUniform(options.TextColor)
	jpgPic, _ := jpeg.Decode(img)
	rgba := imageToRGBA(jpgPic)
	d := &font.Drawer{
		Dst: rgba,
		Src: fg,
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    options.Size,
			DPI:     options.Dpi,
			Hinting: font.HintingNone,
		}),
	}

	writePaymentThrough(d, through)
	writeBuyAmount(d, buyAmt)
	writePaymentSum(d, paymentSum)
	writePaymentAddress(d, paymentAddr)
	outFile, err := os.Create(outputFileName)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)

	err = jpeg.Encode(b, rgba, &jpeg.Options{100})
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("Wrote %s OK.", outputFileName)
}
