package commands

import (
	"DiscordBot/util/rest"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math/rand"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	normal int = iota
	gray
	inverse
	// rainbow
	totalEffects
)

func getCacheDir(nestedPath string) (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	path := cacheDir + "/DiscordBot"

	if nestedPath != "" {
		path += "/" + nestedPath
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if e := os.MkdirAll(path, os.ModePerm); e != nil {
			return "", e
		}
	}

	return path, nil
}

func grayScaleImg(f *os.File, destPath string) error {
	img, format, err := image.Decode(f)
	if err != nil {
		return err
	}

	fmt.Println("format is " + format)

	size := img.Bounds().Size()
	rect := image.Rect(0, 0, size.X, size.Y)
	wImg := image.NewRGBA(rect)

	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			pixel := img.At(x, y)
			originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)

			r := float64(originalColor.R) * 0.92126
			g := float64(originalColor.G) * 0.97152
			b := float64(originalColor.B) * 0.90722

			grey := uint8((r + g + b) / 3)
			c := color.RGBA{
				R: grey,
				G: grey,
				B: grey,
				A: originalColor.A,
			}
			wImg.Set(x, y, c)
		}
	}

	fg, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer fg.Close()

	err = jpeg.Encode(fg, wImg, nil)
	if err != nil {
		return err
	}

	return nil
}

func rainbowImg(f *os.File, destPath string) error {
	img, format, err := image.Decode(f)
	if err != nil {
		return err
	}

	fmt.Println("format is " + format)

	size := img.Bounds().Size()
	rect := image.Rect(0, 0, size.X, size.Y)
	wImg := image.NewRGBA(rect)

	var commonAlpha uint8 = 100

	rainbowColors := []color.RGBA{
		color.RGBA{ // red
			R: 255, B: 0, G: 0, A: commonAlpha,
		},
		color.RGBA{ // orange
			R: 255, B: 127, G: 0, A: commonAlpha,
		},
		color.RGBA{ // yellow
			R: 255, B: 127, G: 0, A: commonAlpha,
		},
		color.RGBA{ // yellow
			R: 255, B: 255, G: 0, A: commonAlpha,
		},
		color.RGBA{ // green
			R: 0, B: 255, G: 0, A: commonAlpha,
		},
		color.RGBA{ // blue
			R: 0, B: 0, G: 255, A: commonAlpha,
		},
		color.RGBA{ // indigo
			R: 75, B: 0, G: 130, A: commonAlpha,
		},
		color.RGBA{ // violet
			R: 148, B: 0, G: 211, A: commonAlpha,
		},
	}

	padding := size.X / len(rainbowColors)

	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			padX := x
			for ; padX < x+padding && padX < size.X; padX++ {
				pixel := img.At(padX, y)
				originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)
				rainbowColor := rainbowColors[((x/padding)+y)%len(rainbowColors)]

				alpha := 1 - (1-rainbowColor.A)*(1-originalColor.A)

				c := color.RGBA{
					R: (rainbowColor.R * rainbowColor.A / alpha) + (originalColor.R * originalColor.A * (1 - rainbowColor.A) / alpha),
					G: (rainbowColor.G * rainbowColor.A / alpha) + (originalColor.G * originalColor.A * (1 - rainbowColor.A) / alpha),
					B: (rainbowColor.B * rainbowColor.A / alpha) + (originalColor.B * originalColor.A * (1 - rainbowColor.A) / alpha),
					A: alpha,
				}
				wImg.Set(padX, y, c)
			}
			x = padX
		}
	}

	fg, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer fg.Close()

	err = jpeg.Encode(fg, wImg, nil)
	if err != nil {
		return err
	}

	return nil
}

func inverseImg(f *os.File, destPath string) error {
	img, format, err := image.Decode(f)
	if err != nil {
		return err
	}

	fmt.Println("format is " + format)

	size := img.Bounds().Size()
	rect := image.Rect(0, 0, size.X, size.Y)
	wImg := image.NewRGBA(rect)

	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			pixel := img.At(x, y)
			originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)

			c := color.RGBA{
				R: 255 - originalColor.R,
				G: 255 - originalColor.G,
				B: 255 - originalColor.B,
				A: originalColor.A,
			}
			wImg.Set(x, y, c)
		}
	}

	fg, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer fg.Close()

	err = jpeg.Encode(fg, wImg, nil)
	if err != nil {
		return err
	}

	return nil
}

func Me(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	avatarUrl := m.Author.AvatarURL("")
	var extension string
	{
		splitUrl := strings.Split(avatarUrl, ".")
		extension = splitUrl[len(splitUrl)-1]
	}

	cacheDir, err := getCacheDir("avatars/")
	if err != nil {
		fmt.Printf("Error getCacheDir() %s\n", err.Error())
		s.ChannelMessageSend(m.ChannelID, "Could not get image")
		return
	}

	path := cacheDir + m.Author.Avatar + "." + extension
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("Downloading file %s into file %s\n", avatarUrl, path)
		if e := rest.DownloadFile(avatarUrl, path); e != nil {
			fmt.Printf("Error Download file %s\n", e.Error())
			s.ChannelMessageSend(m.ChannelID, "Could not fetch profile image from discord servers")
			return
		}
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening file %s\n", err.Error())
		s.ChannelMessageSend(m.ChannelID, "Could not open cached profile image")
		return
	}
	defer f.Close()

	effect := rand.Intn(totalEffects)

	var message *discordgo.MessageSend

	switch effect {
	case normal:
		message = &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Image: &discordgo.MessageEmbedImage{
					URL: "attachment://" + path,
				},
			},
			Files: []*discordgo.File{
				{
					Name:   path,
					Reader: f,
				},
			},
		}
	case gray:
		grayPath := cacheDir + m.Author.Avatar + "-gray." + extension

		if _, err := os.Stat(grayPath); os.IsNotExist(err) {
			grayScaleImg(f, grayPath)
		}

		grayF, err := os.Open(grayPath)
		if err != nil {
			fmt.Printf("Error opening gray file %s\n", err.Error())
			s.ChannelMessageSend(m.ChannelID, "Could not open cached gray profile image")
			return
		}
		defer grayF.Close()

		message = &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Image: &discordgo.MessageEmbedImage{
					URL: "attachment://" + grayPath,
				},
			},
			Files: []*discordgo.File{
				{
					Name:   grayPath,
					Reader: grayF,
				},
			},
		}
	case inverse:
		inversePath := cacheDir + m.Author.Avatar + "-inverse." + extension

		if _, err := os.Stat(inversePath); os.IsNotExist(err) {
			inverseImg(f, inversePath)
		}

		inverseF, err := os.Open(inversePath)
		if err != nil {
			fmt.Printf("Error opening inverse file %s\n", err.Error())
			s.ChannelMessageSend(m.ChannelID, "Could not open cached inverse profile image")
			return
		}
		defer inverseF.Close()

		message = &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Image: &discordgo.MessageEmbedImage{
					URL: "attachment://" + inversePath,
				},
			},
			Files: []*discordgo.File{
				{
					Name:   inversePath,
					Reader: inverseF,
				},
			},
		}

	/*
		case rainbow:
			rainbowPath := cacheDir + m.Author.Avatar + "-rainbow." + extension

			if _, err := os.Stat(rainbowPath); os.IsNotExist(err) {
				rainbowImg(f, rainbowPath)
			}

			rainbowF, err := os.Open(rainbowPath)
			if err != nil {
				fmt.Printf("Error opening rainbow file %s\n", err.Error())
				s.ChannelMessageSend(m.ChannelID, "Could not open cached rainbow profile image")
				return
			}
			defer rainbowF.Close()

			message = &discordgo.MessageSend{
				Embed: &discordgo.MessageEmbed{
					Image: &discordgo.MessageEmbedImage{
						URL: "attachment://" + rainbowPath,
					},
				},
				Files: []*discordgo.File{
					{
						Name:   rainbowPath,
						Reader: rainbowF,
					},
				},
			}
	*/

	default:
		message = &discordgo.MessageSend{
			Content: "wtf",
		}

	}

	s.ChannelMessageSendComplex(m.ChannelID, message)
	// s.ChannelMessageSend(m.ChannelID, m.Author.AvatarURL(""))
}
