package parse

import (
	"archive/zip"
	"bufio"

	"github.com/Sirupsen/logrus"
)

var log *logrus.Logger

// Entry is the entry point for this package
func Entry() {
	log = logrus.New()
	// Open a zip archive for reading.
	r, err := zip.OpenReader("./download/2013/01%20January/AUD_USD_Week2.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {
		log.Infoln("Contents of %s:\n", f.Name)
		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}
		scanner := bufio.NewScanner(rc)

		for scanner.Scan() {
			log.Infoln(scanner.Text())
			break
		}

		rc.Close()
		log.Println()
	}
}
