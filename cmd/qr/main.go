package main

import (
	"fmt"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	qr, err := qrcode.New("https://gophercon25eu.glup3.dev")
	if err != nil {
		return fmt.Errorf("create qrcode failed: %v", err)
	}

	options := []standard.ImageOption{
		standard.WithLogoImageFilePNG("assets/gopher.png"),
		standard.WithLogoSizeMultiplier(2),
		standard.WithQRWidth(30),
	}
	writer, err := standard.New("public/qrcode.png", options...)
	if err != nil {
		return fmt.Errorf("create writer failed: %v", err)
	}
	defer writer.Close()

	if err = qr.Save(writer); err != nil {
		fmt.Printf("save qrcode failed: %v\n", err)
	}

	return nil
}
