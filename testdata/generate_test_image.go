package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"

	"github.com/chai2010/webp"
)

func main() {
	// 画像サイズを定義
	width := 800
	height := 600

	// 新しいRGBA画像を作成
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 背景色を青に設定
	blue := color.RGBA{0, 0, 255, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{blue}, image.Point{}, draw.Src)

	// 赤い長方形を描画
	red := color.RGBA{255, 0, 0, 255}
	rect := image.Rect(width/4, height/4, width*3/4, height*3/4)
	draw.Draw(img, rect, &image.Uniform{red}, image.Point{}, draw.Src)

	// 画像をJPEGファイルとして保存
	jpegFile, err := os.Create("testdata/test_image.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer jpegFile.Close()

	// 高品質の画像として保存（圧縮テスト用にあえて高品質に）
	err = jpeg.Encode(jpegFile, img, &jpeg.Options{Quality: 90})
	if err != nil {
		log.Fatal(err)
	}

	// 画像をPNGファイルとしても保存
	pngFile, err := os.Create("testdata/test_image.png")
	if err != nil {
		log.Fatal(err)
	}
	defer pngFile.Close()

	err = png.Encode(pngFile, img)
	if err != nil {
		log.Fatal(err)
	}

	// 画像をWebPファイルとしても保存
	webpFile, err := os.Create("testdata/test_image.webp")
	if err != nil {
		log.Fatal(err)
	}
	defer webpFile.Close()

	err = webp.Encode(webpFile, img, &webp.Options{Quality: 90})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("テスト画像が生成されました: testdata/test_image.jpg, testdata/test_image.png, testdata/test_image.webp")
}
