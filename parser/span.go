package jsonparser

import "io"

type span struct {
	s uint32
	l uint16
}

func (s *span) start() uint32 {
	return s.s
}

func (s *span) end() uint32 {
	return s.s + uint32(s.l)
}

func (s *span) len() uint16 {
	return s.l
}

func (s *span) textContent(r io.ReadSeeker) (string, error) {
	buf := make([]byte, s.l)
	_, err := r.Seek(int64(s.s), int(s.l))
	if err != nil {
		return "", err
	}
	_, err = r.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
