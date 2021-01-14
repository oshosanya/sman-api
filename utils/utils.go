package utils

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/oshosanya/sman/definitions"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/jpeg"
	"os"
)

func CreateIDCard(passportFile *os.File, cardDetails definitions.IDCard) (*os.File, error) {
	imgFile1, err := os.Open("template.jpeg")
	if err != nil {
		return nil, err
	}
	img1, _, err := image.Decode(imgFile1)
	if err != nil {
		return nil, err
	}
	_, err = passportFile.Seek(0, 0)
	img2, _, err := image.Decode(passportFile)
	if err != nil {
		return nil, err
	}

	dstImage := imaging.Resize(img2, 300, 0, imaging.Lanczos)

	// Template image points
	//sp2 := image.Point{img1.Bounds().Dx(), 0}

	//r2 := image.Rectangle{sp2, sp2.Add(img2.Bounds().Size())}

	// Create base rectangle the size of template
	r2 := image.Rectangle{image.Point{0, 0}, img1.Bounds().Size()}
	r := image.Rectangle{image.Point{0, 0}, r2.Max}
	rgba := image.NewRGBA(r2)

	// Draw template on final image
	draw.Draw(rgba, img1.Bounds(), img1, image.Point{0, 0}, draw.Src)
	draw.Draw(rgba, r, dstImage, image.Point{-680, -185}, draw.Src)
	//draw.Dr
	//addLabel(rgba, -200, 0, "Hello")
	addLabel(rgba, 293, 275, cardDetails.Name)
	addLabel(rgba, 330, 315, cardDetails.Position)
	addLabel(rgba, 320, 355, cardDetails.Branch)
	addLabel(rgba, 293, 405, cardDetails.IDNumber)
	//addLabel(rgba, 0, 0, "Hel'poliutjyhtgrsfaexcvbnl;kjhtgfdwaw4e5tyguhin,bmnbvlo")

	out, err := os.Create(cardDetails.Name + "output.jpg")
	if err != nil {
		return nil, err
	}

	var opt jpeg.Options
	opt.Quality = 80

	jpeg.Encode(out, rgba, &opt)

	return out, nil
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{200, 100, 0, 255}
	var (
		fontFace *truetype.Font
		fontSize = 24.0
	)

	fontFace, _ = freetype.ParseFont(goregular.TTF)
	fontDrawer := &font.Drawer{
		Dst: img,
		Src: image.NewUniform(col),
		Face: truetype.NewFace(fontFace, &truetype.Options{
			Size:    fontSize,
			Hinting: font.HintingFull,
		}),
	}
	//textBounds, _ := fontDrawer.BoundString(label)
	xPosition := fixed.I(img.Rect.Min.X) + fixed.I(x)
	//textHeight := textBounds.Max.Y - textBounds.Min.Y
	yPosition := fixed.I(img.Rect.Min.Y) + fixed.I(y)
	fontDrawer.Dot = fixed.Point26_6{
		X: xPosition,
		Y: yPosition,
	}
	fontDrawer.DrawString(label)
}

func UploadFile(filename string) (string, error) {
	var bucket string
	//var bucket, key string
	bucket = os.Getenv("AWS_BUCKET")
	//key = os.Getenv("AWS_KEY")

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-west-3"),
	}))

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	f, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("failed to open file %q, %v", filename, err)
	}

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(f.Name()),
		Body:   f,
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file %q, %v", filename, err)
	}

	return result.Location, nil
}
