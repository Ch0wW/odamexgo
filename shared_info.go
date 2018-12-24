package odamexgo

import "encoding/binary"

// ReadString : Reads the string.
// Odamex ends every string with a 0 byte, so we can delimit them perfectly.
func (s *ServerQuery) ReadString() string {

	i := s.position
	end := s.position

	for s.buffer[end] != 0 {
		end = end + 1
	}

	// Advance the byte
	s.position = end + 1

	return string(s.buffer[i:end])

}

// ReadByte : Reads the byte.
// The second argument was done to prevent a warning from go-req
func (s *ServerQuery) ReadByte() (byte, error) {

	i := s.position

	// Advance the byte
	s.position = s.position + 1

	return s.buffer[i], nil
}

// ReadShort : Reads a short after shifting to Little Endian.
func (s *ServerQuery) ReadShort() int16 {

	i := s.position

	res := binary.LittleEndian.Uint16(s.buffer[i:])

	s.position = s.position + 2

	return int16(res)
}

// ReadLong : Reads a long after shifting to Little Endian.
func (s *ServerQuery) ReadLong() int32 {

	i := s.position
	res := binary.LittleEndian.Uint32(s.buffer[i:])

	s.position = s.position + 4

	return int32(res)
}

// ReadBool : Transform the byte into a boolean.
func (s *ServerQuery) ReadBool() bool {

	i := s.position

	s.position = s.position + 1
	if s.buffer[i] >= 1 {
		return true
	}
	return false
}
