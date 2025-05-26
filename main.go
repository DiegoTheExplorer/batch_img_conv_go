package main

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/jdeng/goheif"
	"golang.org/x/image/tiff"
	"golang.org/x/image/webp"
)

func main() {
	fmt.Println("Batch Img Conv GO")

	//Hardcoded input directory
	const inp_dir = "C:\\Users\\DMVil\\Work\\PCAP\\anyimgtojpg\\inp"
	const out_dir = "C:\\Users\\DMVil\\Work\\PCAP\\anyimgtojpg\\out"

	var inp_paths []string
	var undecoded_paths []string

	var total_file_count int = 0
	err := filepath.WalkDir(inp_dir, func(fp string, file fs.DirEntry, err error) error {
		total_file_count++
		if file.IsDir() {
			undecoded_paths = append(undecoded_paths, file.Name())
			return nil
		}
		inp_paths = append(inp_paths, fp)
		file_ext := path.Ext(fp)

		err = nil
		err = convertToJpeg(fp, file_ext, out_dir)

		if err != nil {
			undecoded_paths = append(undecoded_paths, file.Name())
			fmt.Println(err)
			return nil
		}

		fmt.Println("Jpeg encoding success for: ", filepath.Base(fp))
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Undecoded Paths: ")
	for _, file_name := range undecoded_paths {
		fmt.Println(file_name)
	}
	fmt.Println("Total undecoded paths: ", len(undecoded_paths))
	fmt.Println("Total number of files: ", total_file_count)
}

func decodeFromFileExt(img_file *os.File, file_ext string) (image.Image, error) {
	if file_ext == ".jpeg" || file_ext == ".jpg" {
		jpg_img, err := jpeg.Decode(img_file)

		if err == nil {
			return jpg_img, nil
		}
		img_file.Seek(0, 0)
	}
	if file_ext == ".png" {
		png_img, err := png.Decode(img_file)

		if err == nil {
			return png_img, nil
		}
		img_file.Seek(0, 0)
	}
	if file_ext == ".webp" {
		webp_img, err := webp.Decode(img_file)

		if err == nil {
			return webp_img, nil
		}
		img_file.Seek(0, 0)
	}
	if file_ext == "tiff" {
		tiff_img, err := tiff.Decode(img_file)

		if err == nil {
			return tiff_img, nil
		}
		img_file.Seek(0, 0)
	}
	if file_ext == "heif" {
		heif_img, err := goheif.Decode(img_file)

		if err == nil {
			return heif_img, nil
		}
		img_file.Seek(0, 0)
	}

	return nil, fmt.Errorf("File extension does not match image encoding")
}

func bruteForceDecode(img_file *os.File) (image.Image, error) {
	jpg_img, err := jpeg.Decode(img_file)

	if err == nil {
		return jpg_img, nil
	}
	img_file.Seek(0, 0)
	img, err := png.Decode(img_file)

	if err == nil {
		return img, nil
	}
	img_file.Seek(0, 0)
	img, err = webp.Decode(img_file)

	if err == nil {
		return img, nil
	}
	img_file.Seek(0, 0)
	img, err = tiff.Decode(img_file)

	if err == nil {
		return img, nil
	}
	img_file.Seek(0, 0)
	img, err = goheif.Decode(img_file)

	if err == nil {
		return img, nil
	}
	img_file.Seek(0, 0)

	return nil, fmt.Errorf("File extension does not match image encoding")
}

func convertToJpeg(fp string, file_extension string, out_dir string) error {
	// Open the file
	img_file, err := os.Open(fp)
	defer img_file.Close()

	if err != nil {
		fmt.Println("Error opening the file: ", fp)
		return err
	}

	var decoded_img image.Image
	file_ext_guess_fail := false

	// Try to use file extension to infer image decoding type
	if file_extension != "" {
		decoded_img, err = decodeFromFileExt(img_file, file_extension)

		if err != nil {
			file_ext_guess_fail = true
		}
	}

	if file_ext_guess_fail {
		decoded_img, err = bruteForceDecode(img_file)
	}

	if decoded_img == nil {
		err_out := errors.New("Decoded image was nil for: \n	" + fp)
		return err_out
	}

	// concat the filename from fp with out_dir
	new_filename := strings.Split(filepath.Base(fp), ".")[0] + ".jpg"
	out_filepath := out_dir + "/" + new_filename

	// Create file to be written
	out_file, err := os.Create(out_filepath)
	defer out_file.Close()

	if err != nil {
		fmt.Println("Failed to create output file")
		return err
	}

	err = jpeg.Encode(out_file, decoded_img, &jpeg.Options{Quality: 100})
	if err != nil {
		return err
	}

	return nil
}
