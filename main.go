package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

func main() {
	fmt.Println("Batch Img Conv GO")

	//Hardcoded input directory
	const inp_dir = "C:\\Users\\DMVil\\Work\\PCAP\\anyimgtojpg\\inp"
	const out_dir = "C:\\Users\\DMVil\\Work\\PCAP\\anyimgtojpg\\out"

	var inp_paths []string
	var undecoded_paths []string

	err := filepath.WalkDir(inp_dir, func(fp string, file fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if file.IsDir() {
			undecoded_paths = append(undecoded_paths, file.Name())
			return nil
		}
		fmt.Println(fp)
		inp_paths = append(inp_paths, fp)
		file_ext := path.Ext(fp)

		err = convertToJpeg(fp, file_ext, out_dir)

		if err != nil {
			undecoded_paths = append(undecoded_paths, file.Name())
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Undecoded Paths: ")
	for _, file_name := range undecoded_paths {
		fmt.Println(file_name)
	}
}

func decodeFromFileExt(img_file *os.File, file_ext string) (image.Image, error) {
	if file_ext == ".jpeg" || file_ext == ".jpg" {
		jpg_img, err := jpeg.Decode(img_file)

		if err == nil {
			return jpg_img, nil
		}
		img_file.Seek(0, 0)
	}

	return nil, fmt.Errorf("File extension does not match image encoding")
}

func convertToJpeg(fp string, file_extension string, out_dir string) error {

	// Open the file
	img_file, err := os.Open(fp)

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
		// Try all decoders available
	}

	// concat the filename from fp with out_dir
	out_filepath := out_dir + filepath.Base(fp)

	// Create file to be written
	out_file, err := os.Create(out_filepath)
	if err != nil {
		fmt.Println("Failed to create output file")
		return err
	}

	jpeg.Encode(out_file, decoded_img, &jpeg.Options{Quality: 100})

	return fmt.Errorf("Failed to convert the following file to jpeg: ", fp)
}
