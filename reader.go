package pak

import (
	"fmt"
	"io"
	"os"
)

type ReadCloser struct {
	f *os.File
	Reader
}

func (r *ReadCloser) Close() error {
	return r.f.Close()
}

type Reader struct {
	File       []*File
	version    int
	dataOffset int64
	dataSize   int64
}

func (p *Reader) init(r io.ReaderAt, size int64) error {
	if size < hdrPeekLen {
		return ErrNotValidPAK
	}

	// Peek into the PAK header and perform
	// some sanity checks on the contents
	buf := make(tocBuf, hdrPeekLen)
	_, err := r.ReadAt(buf, 0)
	if err != nil {
		return err
	}

	// Verify existence of the magic PAK identifier
	if hdrMagic != string(buf[0:3]) {
		return ErrNotValidPAK
	}
	buf.skip(len(hdrMagic))

	p.version = int(buf.byte())
	p.dataOffset = int64(buf.uint32())
	p.dataSize = int64(buf.uint32())
	tocSize := p.dataOffset - int64(hdrPeekLen+len(hdrData))

	// Actual file size should equal the sum of
	// PAK header+toc length + PAK data length
	if size != p.dataOffset+p.dataSize {
		return ErrBadPAK
	}

	fmt.Printf("Version %d\nTOC size: %d\nData size: %d\n", p.version, tocSize, p.dataSize)

	toc := make([]byte, tocSize)
	_, err = r.ReadAt(toc, hdrPeekLen)
	if err != nil {
		return err
	}

	buf = tocBuf(toc)
	p.parseTOC(nil, &buf)

	/*
		// Verify we have a data identifier after TOC
		peekBuf = peekBuf[0:4]
		r.ReadAt(peekBuf, tocSize)
		if hdrData != string(peekBuf) {
			return ErrBadPAK
		}
	*/

	return nil
}

func (r *Reader) parseTOC(p *File, b *tocBuf) {
	f := new(File)
	r.File = append(r.File, f)
	f.name = b.string()
	f.parent = p
	f.flags = int(b.byte())

	if f.IsDirectory() {
		f.length = int64(b.uint32())
		for i := 0; i < int(f.length); i++ {
			r.parseTOC(f, b)
		}
	} else {
		f.offset = int64(b.uint32())
		f.length = int64(b.uint32())
		f.cksum = b.uint32()
	}
}

func OpenReader(path string) (*ReadCloser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	fstat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	r := new(ReadCloser)
	if err := r.init(f, fstat.Size()); err != nil {
		return nil, err
	}

	r.f = f

	return r, nil
}

type File struct {
	name   string
	parent *File
	offset int64
	length int64
	cksum  uint32
	flags  int
	pakr   io.ReaderAt
}

func (f *File) IsDirectory() bool {
	return int(f.flags)&TypeDir != 0
}

func (f *File) Size() int64 {
	return f.length
}

func (f *File) Cksum() uint32 {
	return f.cksum
}

func (f *File) Path() string {
	if f.parent == nil {
		return f.name
	}
	return f.parent.Path() + "/" + f.name
}

func (f *File) Name() string {
	return f.name
}
