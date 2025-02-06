package game

import (
    "bytes"
    _ "embed"
    "image"
    _ "image/png" // Required for PNG decoding
)

//go:embed assets/icon.png
var iconData []byte

func GetAppIcon() image.Image {
    img, _, err := image.Decode(bytes.NewReader(iconData))
    if err != nil {
        panic(err)
    }
    return img
}