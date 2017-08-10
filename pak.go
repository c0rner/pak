package pak

import (
	"errors"
)

const (
	hdrPeekLen = 12    // Bytes needed to peek into PAK header
	hdrMagic   = "PAK" // Magic PAK file identifier
	hdrData    = "DATA"
)

// File flags
const (
	TypeDir = 1 << iota // Directory
)

// Error messages
var (
	ErrBadPAK      = errors.New("PAK file is damaged")
	ErrNotValidPAK = errors.New("Not a valid PAK file")
	ErrUnsupported = errors.New("Unsupported PAK version")
)
