package main

import (
	"fmt"
	"io/fs"
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

func convertToJpeg(fp string, file_extension string, out_dir string) error {

	return fmt.Errorf("error")
}
