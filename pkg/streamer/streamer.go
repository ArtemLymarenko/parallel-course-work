package streamer

import (
	"bytes"
)

type Connection interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
}

func readChunks(conn Connection) (chunkSize int, totalChunks int, err error) {
	chunkSizeBuff := make([]byte, 4)
	_, err = conn.Read(chunkSizeBuff)
	if err != nil {
		return 0, 0, err
	}
	chunkSize = int(ParseBufferedInt32(chunkSizeBuff))

	totalChunksBuff := make([]byte, 4)
	_, err = conn.Read(totalChunksBuff)
	if err != nil {
		return 0, 0, err
	}
	totalChunks = int(ParseBufferedInt32(totalChunksBuff))

	return chunkSize, totalChunks, nil
}

func ParseBufferedInt32(buff []byte) int32 {
	if len(buff) != 4 {
		return 0
	}
	return int32(buff[0])<<24 | int32(buff[1])<<16 | int32(buff[2])<<8 | int32(buff[3])
}

func ReadBuff(conn Connection) ([]byte, error) {
	chunkSize, totalChunks, err := readChunks(conn)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	for i := 0; i < totalChunks; i++ {
		chunk := make([]byte, chunkSize)
		n, err := conn.Read(chunk)
		if err != nil {
			return nil, err
		}

		buffer.Write(chunk[:n])
	}

	return buffer.Bytes(), nil
}

func WriteInt32ToBuffer(value int) []byte {
	return []byte{
		byte(value >> 24),
		byte(value >> 16),
		byte(value >> 8),
		byte(value),
	}
}

func WriteBuff(conn Connection, chunkSize int, requestBin []byte) error {
	requestLen := len(requestBin)
	totalChunks := requestLen / chunkSize
	if requestLen%chunkSize != 0 {
		totalChunks++
	}

	_, err := conn.Write(WriteInt32ToBuffer(chunkSize))
	if err != nil {
		return err
	}

	_, err = conn.Write(WriteInt32ToBuffer(totalChunks))
	if err != nil {
		return err
	}

	offset := 0
	for offset < requestLen {
		end := offset + chunkSize
		if end > requestLen {
			end = requestLen
		}

		n, err := conn.Write(requestBin[offset:end])
		if err != nil {
			return err
		}
		offset += n
	}

	return nil
}
