package pak

import "encoding/binary"

type tocBuf []byte

func (b *tocBuf) byte() byte {
	res := (*b)[0]
	*b = (*b)[1:]
	return res
}

func (b *tocBuf) uint32() uint32 {
	res := binary.LittleEndian.Uint32(*b)
	*b = (*b)[4:]
	return res
}

func (b *tocBuf) string() string {
	length := b.byte()
	res := make([]byte, length)
	copy(res, *b)
	*b = (*b)[length:]
	return string(res)
}
