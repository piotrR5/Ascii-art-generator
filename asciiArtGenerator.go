package main

import (
	_ "fmt"
	"image"
	_ "image/jpeg"
	"os"
	"strconv"
	"strings"
)

func useVar(numbers ...any) {
	for _, n := range numbers {
		_ = n
	}
}

func main() {
	asciiArtLogo := "[ AsciiArt ]"
	asciiSamples := "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/|()1{}[]?-_+~<>i!lI;:,^`'.."
	highContrast := false
	useVar(asciiSamples)

	userDefinedHeight := 0
	userDefinedWidth := 0
	userDefinedInput := ""
	userDefinedOutput := ""
	reversed := false
	_, _ = userDefinedOutput, userDefinedInput
	useVar(userDefinedHeight, userDefinedWidth)
	//////////////////////////////////getting the parameters///////////////////////////////
	if len(os.Args) > 1 {
		//fmt.Println(os.Args)
		var arguments []string
		for _, n := range os.Args {
			arguments = append(arguments, strings.Fields(n)...)
		}

		for i := 0; i < len(arguments); i++ {
			if arguments[i] == "-h" {
				userDefinedHeight, _ = strconv.Atoi(arguments[i+1])
			} else if arguments[i] == "-w" {
				userDefinedWidth, _ = strconv.Atoi(arguments[i+1])
			} else if arguments[i] == "-i" {
				userDefinedInput = arguments[i+1]
			} else if arguments[i] == "-o" {
				userDefinedOutput = arguments[i+1]
			} else if arguments[i] == "-r" {
				reversed = true
			} else if arguments[i] == "-c" {
				highContrast = true
			} else {
				continue
			}
		}
	} else if len(os.Args) <= 1 {
		println("Specify arguments: -h for desired height of the image, -w for desired width of the image, -i for input image path and -o for output image path.")
		os.Exit(0)
	}

	if userDefinedHeight == 0 || userDefinedWidth == 0 || userDefinedInput == "" || userDefinedOutput == "" {
		println("Specify arguments: -h for desired height of the image, -w for desired width of the image, -i for input image path and -o for output image path.")
		os.Exit(0)
	}

	println(asciiArtLogo+" desired width: ", userDefinedWidth, " desired height: ", userDefinedHeight, " input path: ", userDefinedInput, " output path: ", userDefinedOutput)

	//////////////////////////////////getting the parameters///////////////////////////////

	//////////////////////////////////getting the input file//////////////////////////////

	inputFile, err := os.Open(userDefinedInput)
	if err != nil {
		println(asciiArtLogo, "Could not open: ", userDefinedInput, " ", err)
		panic(err)
	}
	defer inputFile.Close()

	//////////////////////////////////getting the input file//////////////////////////////

	//////////////////////////////////getting the input image/////////////////////////////
	//inputDataType = strings.Trim(inputDataType, "image/")

	inputImage, dataType, err := image.Decode(inputFile)

	if err != nil {
		panic(err)
	}

	inputHeight, inputWidth := inputImage.Bounds().Max.Y, inputImage.Bounds().Max.X
	useVar(inputImage, dataType, inputHeight, inputWidth)

	println(asciiArtLogo, "detected image type: ", dataType)
	println(asciiArtLogo, "input image height: ", inputHeight, " input image width: ", inputWidth)

	scaleX := inputWidth / userDefinedWidth
	scaleY := inputHeight / userDefinedHeight

	if scaleX == 0 {scaleX++}
	if scaleY == 0 {scaleY++}

	println(asciiArtLogo, "scale Y:  ", scaleY, " scale X: ", scaleX)

	////////////////////////////////////////////////////////////////////////////////////

	//////////////////////////////////////making image grayscale////////////////////////
	grayImg := image.NewGray(inputImage.Bounds())
	for y := 0; y < inputHeight; y++ {
		for x := 0; x < inputWidth; x++ {
			grayImg.Set(x, y, inputImage.At(x, y))
		}
	}

	//////debug write grayscale image////////
	/*
		debugFile, err := os.Create("grayscale.png")
		if err != nil {
			// handle error
			log.Fatal(err)
		}
		defer debugFile.Close()

		if err := png.Encode(debugFile, grayImg); err != nil {
			log.Fatal(err)
		}
	*/
	//////debug write grayscale image////////

	asciiMatrix := make([][]uint8, userDefinedWidth)
	for i := range asciiMatrix {
		asciiMatrix[i] = make([]uint8, userDefinedHeight)
	}
	avarages := make([][]float32, userDefinedWidth)
	for i := range avarages {
		avarages[i] = make([]float32, userDefinedHeight)
	}
	useVar(asciiMatrix, avarages)
	/////////////////////////////////scale image down///////////////////////////////

	for x := 0; x/scaleX < userDefinedWidth; x = x + scaleX {
		for y := 0; y/scaleY < userDefinedHeight; y = y + scaleY {
			avarages[x/scaleX][y/scaleY] += float32(grayImg.Pix[y*grayImg.Stride+x])
		}
	}

	minAvarage, maxAvarage := float32(255.0), float32(0.0)

	if highContrast {
		for i := range avarages {
			for j := range avarages[i] {
				if avarages[i][j] < minAvarage {
					minAvarage = avarages[i][j]
				}
				if avarages[i][j] > maxAvarage {
					maxAvarage = avarages[i][j]
				}
			}
		}

		mean := maxAvarage + minAvarage/2

		for i := range avarages {
			for j := range avarages[i] {
				if avarages[i][j] < mean {
					avarages[i][j] = (avarages[i][j] - minAvarage) / 2
				} else if avarages[i][j] > mean {
					avarages[i][j] = (avarages[i][j] + maxAvarage) / 2
				}
			}
		}
	}

	//fmt.Print(grayImg)
	//fmt.Print(avarages)
	/////////////////////////////////scale image down///////////////////////////////

	//debug create scaled down image//
	/*

		scaledImage := image.NewGray(image.Rect(0, 0, userDefinedWidth, userDefinedHeight))
		for x := 0; x < userDefinedWidth; x++ {
			for y := 0; y < userDefinedHeight; y++ {
				scaledImage.SetGray(x, y, color.Gray{uint8(avarages[x][y])})
			}
		}

		scaledFile, err := os.Create("scaledDown.png")
		if err != nil {
			// handle error
			panic(err)
		}
		defer scaledFile.Close()

		if err := png.Encode(scaledFile, scaledImage); err != nil {
			panic(err)
		}

	*/
	//debug create scaled down image//

	/////////////////////////////////convert to ascii matrix////////////////////////

	for n := range avarages {
		for m := range avarages[n] {
			var index int
			if reversed {
				index = (len(asciiSamples) * int(avarages[n][m]) / 256)
			} else {
				index = 67 - (len(asciiSamples) * int(avarages[n][m]) / 256)
			}
			asciiMatrix[n][m] = asciiSamples[index]
		}
	}

	print(asciiMatrix)

	outputFile, err := os.Create(userDefinedOutput)

	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	for x := 0; x < len(asciiMatrix[0]); x++ {
		for y := 0; y < len(asciiMatrix); y++ {
			outputFile.WriteString(string(asciiMatrix[y][x]))
		}
		outputFile.WriteString("\n")
	}

	////////////////////////////////////////////////////////////////////////////////
}
