package unarchive

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v4"
)

func Unarchive(input io.Reader, dir string) error {
	// TODO: consider if should write to a more generic interface
	// like a writer, or if maybe if the function itself
	// should take the handler as an input so can be as generic
	// as you'd like in the handler
	format, input, err := archiver.Identify("", input)
	if err != nil {
		return err
	}
	// the list of files we want out of the archive; any
	// directories will include all their contents unless
	// we return fs.SkipDir from our handler
	// (leave this nil to walk ALL files from the archive)

	handler := func(ctx context.Context, f archiver.File) error {
		fileName := f.NameInArchive
		// currently on osx we get a top dir of ./bin and ./share
		// when in reality the
		if strings.HasPrefix(fileName, "quarto-") {
			fileName = strings.Join(strings.Split(fileName, "/")[1:], "/")
		}
		// do something with the file
		newPath := filepath.Join(dir, fileName)
		if f.IsDir() {
			return os.MkdirAll(newPath, f.Mode())
		}
		newFile, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY, f.Mode())
		if err != nil {
			return err
		}
		defer newFile.Close()
		// copy file data into tar writer
		af, err := f.Open()
		if err != nil {
			return err
		}
		defer af.Close()
		if _, err := io.Copy(newFile, af); err != nil {
			return err
		}
		return nil
	}
	// make sure the format is capable of extracting
	ex, ok := format.(archiver.Extractor)
	if !ok {
		return err
	}
	return ex.Extract(context.Background(), input, nil, handler)
}
