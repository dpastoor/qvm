package unarchive

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v4"
)

func trimTopDir(dir string) string {
	if pos := strings.Index(dir, "/"); pos >= 0 {
		return dir[pos+1:]
	}
	return dir
}

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

	// not sure if this should be a syncmap, or if a map is ok?
	// not sure if the handler itself is invoked serially or if it
	// is concurrent?
	dirMap := map[string]bool{}
	handler := func(ctx context.Context, f archiver.File) error {
		fileName := f.NameInArchive
		// currently on osx we get a top dir of ./bin and ./share
		// when in reality the
		if strings.HasPrefix(fileName, "quarto-") {
			fileName = trimTopDir(fileName)
		}
		newPath := filepath.Join(dir, fileName)
		fdir := filepath.Dir(newPath)
		if f.IsDir() {
			dirMap[fdir] = true
			return os.MkdirAll(newPath, f.Mode())
		} else {
			// check if we've seen the dir before, if not, we'll attemp to create
			// it in case its not there. This needs to be done as archive formats
			// do not necessarily always have the directory in order/present
			// eg zip dirs for quarto definitely are missing seemingly random dirs
			// when talking with charles about it, we were both unsure what might
			// be the reason, and assume its probably the powershell compress-archive
			// encantation, so rather than trying to go down that rabbit hole too far,
			// some additional checking here
			_, seenDir := dirMap[fdir]
			if !seenDir {
				dirMap[fdir] = true
				// linux default for new directories is 777 and let the umask handle
				// if should have other controls
				err := os.MkdirAll(newPath, 0777)
				if err != nil {
					return err
				}
			}
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
