package url

import (
	"fmt"
	"io"
	neturlpkg "net/url"
	"path"

	namepkg "github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/tliron/puccini/common/registry"
)

// TODO: support authentication

//
// DockerURL
//

type DockerURL struct {
	URL     *neturlpkg.URL
	String_ string `json:"string" yaml:"string"`
}

func NewDockerURL(neturl *neturlpkg.URL) *DockerURL {
	return &DockerURL{neturl, neturl.String()}
}

func NewValidDockerURL(neturl *neturlpkg.URL) (*DockerURL, error) {
	if (neturl.Scheme != "docker") && (neturl.Scheme != "") {
		return nil, fmt.Errorf("not a docker URL: %s", neturl.String())
	}

	// TODO

	return NewDockerURL(neturl), nil
}

// URL interface
// fmt.Stringer interface
func (self *DockerURL) String() string {
	return self.Key()
}

// URL interface
func (self *DockerURL) Format() string {
	format := self.URL.Query().Get("format")
	if format != "" {
		return format
	} else {
		return GetFormat(self.URL.Path)
	}
}

// URL interface
func (self *DockerURL) Origin() URL {
	url := *self
	url.URL.Path = path.Dir(url.URL.Path)
	return &url
}

// URL interface
func (self *DockerURL) Relative(path string) URL {
	if neturl, err := neturlpkg.Parse(path); err == nil {
		return NewDockerURL(self.URL.ResolveReference(neturl))
	} else {
		return nil
	}
}

// URL interface
func (self *DockerURL) Key() string {
	return self.String_
}

// URL interface
func (self *DockerURL) Open() (io.ReadCloser, error) {
	pipeReader, pipeWriter := io.Pipe()

	go func() {
		if err := self.WriteLayer(pipeWriter); err == nil {
			pipeWriter.Close()
		} else {
			pipeWriter.CloseWithError(err)
		}
	}()

	return pipeReader, nil
}

// URL interface
func (self *DockerURL) Release() error {
	return nil
}

func (self *DockerURL) WriteTarball(writer io.Writer) error {
	url := fmt.Sprintf("%s%s", self.URL.Host, self.URL.Path)
	if tag, err := namepkg.NewTag(url); err == nil {
		if image, err := remote.Image(tag); err == nil {
			return tarball.Write(tag, image, writer)
		} else {
			return err
		}
	} else {
		return err
	}
}

func (self *DockerURL) WriteLayer(writer io.Writer) error {
	pipeReader, pipeWriter := io.Pipe()

	go func() {
		if err := self.WriteTarball(pipeWriter); err != nil {
			pipeWriter.Close()
		} else {
			pipeWriter.CloseWithError(err)
		}
	}()

	decoder := registry.NewImageLayerDecoder(pipeReader)
	if _, err := io.Copy(writer, decoder.Decode()); err == nil {
		return nil
	} else {
		return err
	}
}
