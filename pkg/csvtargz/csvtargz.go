package csvtargz

import (
	"archive/tar"
	"compress/gzip"
	"encoding/csv"
	"io"
	"os"

	"github.com/jszwec/csvutil"
	"github.com/pkg/errors"
)

// ErrNoSuchFile is returned when no such file exists in archive.
var ErrNoSuchFile = errors.New("no such file")

// DecodeCSVFromTarGz decodes a CSV file from .tar.gz archive into dst.
func DecodeCSVFromTarGz(archivePath, csvFilename string, dst interface{}) error {
	return withCSVReaderFromTarGz(archivePath, csvFilename, func(csvReader *csv.Reader) error {
		csvDecoder, err := csvutil.NewDecoder(csvReader)
		if err != nil {
			return err
		}

		return csvDecoder.Decode(dst)
	})
}

func withCSVReaderFromTarGz(archivePath, csvFilename string, f func(csvReader *csv.Reader) error) error {
	gzFile, err := os.Open(archivePath)
	if err != nil {
		return err
	}

	defer func() {
		_ = gzFile.Close()
	}()

	gzReader, err := gzip.NewReader(gzFile)
	if err != nil {
		return err
	}

	defer func() {
		_ = gzReader.Close()
	}()

	csvReader, err := newCSVReaderFromTar(tar.NewReader(gzReader), csvFilename)
	if err != nil {
		return err
	}

	if err := f(csvReader); err != nil {
		return err
	}

	if err := gzReader.Close(); err != nil {
		return err
	}

	return gzFile.Close()
}

func newCSVReaderFromTar(tr *tar.Reader, csvFilename string) (*csv.Reader, error) {
	// iterate through the files in the archive and return csv reader of searched file
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return nil, ErrNoSuchFile
		}

		if err != nil {
			return nil, err
		}

		if hdr.Name == csvFilename {
			break
		}
	}

	return csv.NewReader(tr), nil
}
