package wav

import (
	"bufio"
	bin "encoding/binary"
	"fmt"
	"io"
	"os"
)

type WavData struct {
	ChunkID   [4]byte // B
	ChunkSize uint32  // L
	Format    [4]byte // B

	AudioFormat   uint16 // L
	NumChannels   uint16 // L
	SampleRate    uint32 // L
	ByteRate      uint32 // L
	BlockAlign    uint16 // L
	BitsPerSample uint16 // L

	Data []byte // L
}

func ReadWavData(file io.ReadSeeker) (wav WavData, err error) {
	var SubchunkID [4]byte
	var SubchunkSize uint32
	var SubchunkStart int64

	var file_buf *bufio.Reader

	if err = bin.Read(file, bin.BigEndian, &wav.ChunkID); err != nil {
		return
	}
	if err = bin.Read(file, bin.LittleEndian, &wav.ChunkSize); err != nil {
		return
	}
	if err = bin.Read(file, bin.BigEndian, &wav.Format); err != nil {
		return
	}

	for {
		if err = bin.Read(file, bin.BigEndian, &SubchunkID); err != nil {
			return
		}
		if err = bin.Read(file, bin.LittleEndian, &SubchunkSize); err != nil {
			return
		}

		SubchunkStart, err = file.Seek(0, 1)
		if err != nil {
			return
		}
		file_buf = bufio.NewReader(file)

		switch string(SubchunkID[:4]) {
		case "fmt ":
			if err = bin.Read(file_buf, bin.LittleEndian, &wav.AudioFormat); err != nil {
				return
			}
			if err = bin.Read(file_buf, bin.LittleEndian, &wav.NumChannels); err != nil {
				return
			}
			if err = bin.Read(file_buf, bin.LittleEndian, &wav.SampleRate); err != nil {
				return
			}
			if err = bin.Read(file_buf, bin.LittleEndian, &wav.ByteRate); err != nil {
				return
			}
			if err = bin.Read(file_buf, bin.LittleEndian, &wav.BlockAlign); err != nil {
				return
			}
			if err = bin.Read(file_buf, bin.LittleEndian, &wav.BitsPerSample); err != nil {
				return
			}

		case "data":
			wav.Data = make([]byte, SubchunkSize)
			err = bin.Read(file_buf, bin.LittleEndian, &wav.Data)
			return
		}
		if err = file.Seek(SubchunkStart+int64(SubchunkSize), 0); err != nil {
			return
		}
	}

	/*
	 *   fmt.Printf( "\n" )
	 *   fmt.Printf( "ChunkID*: %s\n", ChunkID )
	 *   fmt.Printf( "ChunkSize: %d\n", ChunkSize )
	 *   fmt.Printf( "Format: %s\n", Format )
	 *   fmt.Printf( "\n" )
	 *   fmt.Printf( "Subchunk1ID: %s\n", Subchunk1ID )
	 *   fmt.Printf( "Subchunk1Size: %d\n", Subchunk1Size )
	 *   fmt.Printf( "AudioFormat: %d\n", AudioFormat )
	 *   fmt.Printf( "NumChannels: %d\n", NumChannels )
	 *   fmt.Printf( "SampleRate: %d\n", SampleRate )
	 *   fmt.Printf( "ByteRate: %d\n", ByteRate )
	 *   fmt.Printf( "BlockAlign: %d\n", BlockAlign )
	 *   fmt.Printf( "BitsPerSample: %d\n", BitsPerSample )
	 *   fmt.Printf( "\n" )
	 *   fmt.Printf( "Subchunk2ID: %s\n", Subchunk2ID )
	 *   fmt.Printf( "Subchunk2Size: %d\n", Subchunk2Size )
	 *   fmt.Printf( "NumSamples: %d\n", Subchunk2Size / uint32(NumChannels) / uint32(BitsPerSample/8) )
	 *   fmt.Printf( "\ndata: %v\n", len(data) )
	 *   fmt.Printf( "\n\n" )
	 */
	return
}

const (
	mid16 uint16 = 1 >> 2
	big16 uint16 = 1 >> 1
	big32 uint32 = 65535
)

func btou(b []byte) (u []uint16) {
	u = make([]uint16, len(b)/2)
	for i, _ := range u {
		val := uint16(b[i*2])
		val += uint16(b[i*2+1]) << 8
		u[i] = val
	}
	return
}

func btoi16(b []byte) (u []int16) {
	u = make([]int16, len(b)/2)
	for i, _ := range u {
		val := int16(b[i*2])
		val += int16(b[i*2+1]) << 8
		u[i] = val
	}
	return
}

func btof32(b []byte) (f []float32) {
	u := btoi16(b)
	f = make([]float32, len(u))
	for i, v := range u {
		f[i] = float32(v) / float32(32768)
	}
	return
}

func utob(u []uint16) (b []byte) {
	b = make([]byte, len(u)*2)
	for i, val := range u {
		lo := byte(val)
		hi := byte(val >> 8)
		b[i*2] = lo
		b[i*2+1] = hi
	}
	return
}
