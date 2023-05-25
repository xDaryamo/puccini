package clout

import (
	contextpkg "context"

	"github.com/tliron/exturl"
	"github.com/tliron/kutil/util"
)

func Load(context contextpkg.Context, url string, format string, urlContext *exturl.Context) (*Clout, error) {
	var url_ exturl.URL

	var err error
	if url != "" {
		if url_, err = urlContext.NewValidURL(context, url, nil); err != nil {
			return nil, err
		}
	} else {
		if url_, err = urlContext.ReadToInternalURLFromStdin(context, format); err != nil {
			return nil, err
		}
	}

	if reader, err := url_.Open(context); err == nil {
		reader = util.NewContextualReadCloser(context, reader)
		defer reader.Close()
		return Read(reader, url_.Format())
	} else {
		return nil, err
	}
}
