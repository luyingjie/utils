package utils

import (
	"io"
	"io/ioutil"
)

type ReadCloser struct {
	index      int
	content    []byte
	repeatable bool
}

func NewReadCloser(content []byte, repeatable bool) io.ReadCloser {
	return &ReadCloser{
		content:    content,
		repeatable: repeatable,
	}
}

func NewReadCloserWithReadCloser(r io.ReadCloser, repeatable bool) (io.ReadCloser, error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return &ReadCloser{
		content:    content,
		repeatable: repeatable,
	}, nil
}

func (b *ReadCloser) Read(p []byte) (n int, err error) {
	n = copy(p, b.content[b.index:])
	b.index += n
	if b.index >= len(b.content) {
		// Make it repeatable reading.
		if b.repeatable {
			b.index = 0
		}
		return n, io.EOF
	}
	return n, nil
}

func (b *ReadCloser) Close() error {
	return nil
}
