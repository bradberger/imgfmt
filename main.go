package main

import (
	"./optimizer"
	"bufio"
	"flag"
	"fmt"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/font"
	_ "golang.org/x/image/riff"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/vp8"
	_ "golang.org/x/image/vp8l"
	_ "golang.org/x/image/webp"
	_ "golang.org/x/image/webp/nycbcra"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime"
	"os"
	"path"
)

func main() {

	var img image.Image
	var err error
	var origFmt string

	// Process flags.
	mimeFlag := flag.String("format", "", "The mime type to output. If none specified, will default to the output file format or original format.")
	qualityFlag := flag.Int("quality", 0, "The quality level for the final image")
	widthFlag := flag.Uint("width", 0, "The width of the final image")
	heightFlag := flag.Uint("height", 0, "The height of the final image")
	dprFlag := flag.Float64("dpr", 1.0, "The Viewport DPR to optimize for")
	downlinkFlag := flag.Float64("downlink", 0.384, "The downlink speed to optimize for")
	saveDataFlag := flag.Bool("savedata", false, "Optimize to save data")
	flag.Parse()

	args := flag.Args()
	numArgs := len(args)
	stat, _ := os.Stdin.Stat()

	// @todo If using stdin, decode stdin into image, otherwise decode file into image.
	if (stat.Mode() & os.ModeCharDevice) == 0 {

		img, origFmt, err = image.Decode(bufio.NewReader(os.Stdin))

	} else {

		// If no stdin and no file, return error.
		if numArgs < 1 {
			fmt.Println("You must supply a file via stdin or the first argument")
			os.Exit(1)
		}

		// Try to open and read the file.
		f, e := os.Open(args[0])
		if e != nil {
			fmt.Printf("Error reading file: %s\n", e)
			os.Exit(1)
		}

		defer f.Close()
		img, origFmt, err = image.Decode(bufio.NewReader(f))

	}

	// Check to see decoding suceeded.
	if err != nil {
		fmt.Printf("Could not decode image from stdin: %s\n", err)
		os.Exit(1)
	}

	// Check for mime output flag.
	target := ""
	mimeType := *mimeFlag
	if mimeType == "" {

		// If using output file, get from extension.
		if numArgs > 1 {
			target = args[1]
		} else if numArgs > 1 {
			target = args[0]
		}

		if target != "" {
			mimeType = mime.TypeByExtension(path.Ext(target))
		} else {
			// Still no mime? Set to input image.
			mimeType = origFmt
		}

	}

	// Output. Check if target out, otherwise write to stdout.
	var w io.Writer

	if numArgs > 1 {

		// Output to file.
		f, err := os.Create(args[1])
		if err != nil {
			fmt.Printf("Could not write to %s\n", args[1])
			os.Exit(1)
		}

		defer f.Close()
		w = bufio.NewWriter(f)

	} else {
		// Output to stdout.
		w = bufio.NewWriter(os.Stdout)
	}

	// If here, have an io.Writer in w and are read to rock.
	o := optimizer.Options{
		Mime:     mimeType,
		Width:    *widthFlag,
		Height:   *heightFlag,
		Dpr:      *dprFlag,
		Quality:  *qualityFlag,
		SaveData: *saveDataFlag,
		Downlink: *downlinkFlag,
	}

	err = optimizer.Encode(w, img, o)
	if err != nil {
		fmt.Printf("Error encoding image: %s", err)
		os.Exit(1)
	}

}
