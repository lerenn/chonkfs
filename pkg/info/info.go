package info

type Directory struct {
}

type File struct {
	Size          int
	ChunkSize     int
	ChunksCount   int
	LastChunkSize int
}
