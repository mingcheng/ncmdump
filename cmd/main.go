package main

import (
	"fmt"
	"github.com/mingcheng/ncmdump"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func processFile(name string) error {
	name, err := filepath.Abs(name)
	if err != nil {
		return err
	}

	fp, err := os.Open(name)
	if err != nil {
		log.Println(err)
		return err
	}
	defer fp.Close()

	meta, err := ncmdump.DumpMeta(fp)
	if err != nil {
		return err
	}

	data, err := ncmdump.Dump(fp)
	if err != nil {
		return err
	}

	outputFilePath := fmt.Sprintf("%s.%s",
		strings.TrimSuffix(filepath.Base(name), filepath.Ext(name)), meta.Format)
	outputFilePath = filepath.Join(filepath.Dir(name), outputFilePath)

	if err := ioutil.WriteFile(outputFilePath, data, 0644); err != nil {
		return err
	}

	if err := addMeta(outputFilePath, meta); err != nil {
		return err
	}

	return nil
}

// addMeta to update music file meta from dumped data
func addMeta(musicFile string, meta ncmdump.Meta) error {
	switch strings.ToLower(meta.Format) {
	case "mp3":
		modifier := MP3{
			FilePath: musicFile,
			Meta:     meta,
		}

		return modifier.Update()

	default:
		return fmt.Errorf("unknown format %s", meta.Format)
	}

	return nil
}

func main() {
	argc := len(os.Args)
	if argc <= 1 {
		log.Println("please input file path!")
		return
	}
	files := make([]string, 0)

	for i := 0; i < argc-1; i++ {
		path := os.Args[i+1]
		if info, err := os.Stat(path); err != nil {
			log.Fatalf("Path %s does not exist.", info)
		} else if info.IsDir() {
			filelist, err := ioutil.ReadDir(path)
			if err != nil {
				log.Fatalf("Error while reading %s: %s", path, err.Error())
			}
			for _, f := range filelist {
				files = append(files, filepath.Join(path, "./", f.Name()))
			}
		} else {
			files = append(files, path)
		}
	}

	for _, filename := range files {
		if filepath.Ext(filename) == ".ncm" {
			processFile(filename)
		}
	}
}
