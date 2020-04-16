package registry

import (
	"archive/tar"
	"errors"
	"io"
	"strings"
	"sync"

	gzip "github.com/klauspost/pgzip"
)

//
// ImageLayerDecoder
//
// Unzips the first tar entry with a ".tar.gz" extension
//

type ImageLayerDecoder struct {
	reader     io.Reader
	pipeReader *io.PipeReader
	pipeWriter *io.PipeWriter
	waitGroup  sync.WaitGroup
}

func NewImageLayerDecoder(reader io.Reader) *ImageLayerDecoder {
	pipeReader, pipeWriter := io.Pipe()
	return &ImageLayerDecoder{
		reader:     reader,
		pipeReader: pipeReader,
		pipeWriter: pipeWriter,
	}
}

func (self *ImageLayerDecoder) Decode() io.Reader {
	self.waitGroup.Add(1)
	go self.copy()
	return self.pipeReader
}

func (self *ImageLayerDecoder) Drain() {
	self.waitGroup.Wait()
}

func (self *ImageLayerDecoder) copy() {
	defer self.waitGroup.Done()

	tarReader := tar.NewReader(self.reader)

	for {
		if header, err := tarReader.Next(); err == nil {
			if (header.Typeflag == tar.TypeReg) && strings.HasSuffix(header.Name, ".tar.gz") {
				if gzipReader, err := gzip.NewReader(tarReader); err == nil {
					if _, err := io.Copy(self.pipeWriter, gzipReader); err == nil {
						self.pipeWriter.Close()
						break
					} else {
						self.pipeWriter.CloseWithError(err)
						break
					}
				} else {
					self.pipeWriter.CloseWithError(err)
					break
				}
			}
		} else if err == io.EOF {
			self.pipeWriter.CloseWithError(errors.New("\"*.tar.gz\" not found in tar"))
			break
		} else {
			self.pipeWriter.CloseWithError(err)
			break
		}
	}
}
