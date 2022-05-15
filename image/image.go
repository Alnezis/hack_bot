package image

import (
	"encoding/base64"
	"errors"
	"github.com/fogleman/gg"
	"image/color"
	"io"
	"net/http"
	"os"
)

func DrawCircle(img *gg.Context, x0, y0, r int, c color.Color) {
	x, y, dx, dy := r-1, 0, 1, 1
	err := dx - (r * 2)
	img.SetColor(c)
	for x > y {
		img.SetPixel(x0+x, y0+y)
		img.SetPixel(x0+y, y0+x)
		img.SetPixel(x0-y, y0+x)
		img.SetPixel(x0-x, y0+y)
		img.SetPixel(x0-x, y0-y)
		img.SetPixel(x0-y, y0-x)
		img.SetPixel(x0+y, y0-x)
		img.SetPixel(x0+x, y0-y)

		if err <= 0 {
			y++
			err += dy
			dy += 2
		}
		if err > 0 {
			x--
			dx += 2
			err += dx - (r * 2)
		}
	}
}

func DownloadFile(URL, fileName string) error {
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func ToBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
