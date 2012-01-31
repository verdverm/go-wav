package wav

import (
  bin "encoding/binary"
  "os"
  "bufio"
  "fmt"
)

type WavData struct {
  bChunkID [4]byte         // B
  ChunkSize uint32       // L
  bFormat [4]byte          // B

  bSubchunk1ID [4]byte     // B
  Subchunk1Size uint32   // L
  AudioFormat uint16     // L
  NumChannels uint16     // L
  SampleRate uint32      // L
  ByteRate uint32        // L
  BlockAlign uint16      // L
  BitsPerSample uint16   // L

  bSubchunk2ID [4]byte     // B
  Subchunk2Size uint32   // L
  data []byte             // L
}

func ReadWavData( fn string ) (wav WavData) {
  ftotal, err := os.OpenFile(fn, os.O_RDONLY, 0)
  if err != nil {
    fmt.Printf( "Error opening\n" )
  }
  file := bufio.NewReader(ftotal)

  bin.Read( file, bin.BigEndian, &wav.bChunkID )
  bin.Read( file, bin.LittleEndian, &wav.ChunkSize )
  bin.Read( file, bin.BigEndian, &wav.bFormat )

  bin.Read( file, bin.BigEndian, &wav.bSubchunk1ID )
  bin.Read( file, bin.LittleEndian, &wav.Subchunk1Size )
  bin.Read( file, bin.LittleEndian, &wav.AudioFormat )
  bin.Read( file, bin.LittleEndian, &wav.NumChannels )
  bin.Read( file, bin.LittleEndian, &wav.SampleRate )
  bin.Read( file, bin.LittleEndian, &wav.ByteRate )
  bin.Read( file, bin.LittleEndian, &wav.BlockAlign )
  bin.Read( file, bin.LittleEndian, &wav.BitsPerSample )


  bin.Read( file, bin.BigEndian, &wav.bSubchunk2ID )
  bin.Read( file, bin.LittleEndian, &wav.Subchunk2Size )

  wav.data = make( []byte, wav.Subchunk2Size )
  bin.Read( file, bin.LittleEndian, &wav.data )

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
  mid16 uint16 = 1>>2
  big16 uint16 = 1>>1
  big32 uint32 = 65535
)

func btou( b []byte ) (u []uint16) {
  u = make( []uint16, len(b)/2 )
  for i,_ := range u {
    val := uint16(b[i*2])
    val += uint16(b[i*2+1])<<8
    u[i] = val
  }
  return
}

func btoi16( b []byte ) (u []int16) {
  u = make( []int16, len(b)/2 )
  for i,_ := range u {
    val := int16(b[i*2])
    val += int16(b[i*2+1])<<8
    u[i] = val
  }
  return
}

func btof32( b []byte ) (f []float32) {
  u := btoi16(b)
  f = make([]float32, len(u))
  for i,v := range u {
    f[i] = float32(v)/float32(32768)
  }
  return
}

func utob( u []uint16 ) (b []byte) {
  b = make( []byte, len(u)*2 )
  for i,val := range u {
    lo := byte(val)
    hi := byte(val>>8)
    b[i*2] = lo
    b[i*2+1] = hi
  }
  return
}



