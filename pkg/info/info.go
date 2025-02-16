package info

// Directory represents a directory information.
type Directory struct {
}

// File represents a file information.
type File struct {
	Size          int
	ChunkSize     int
	ChunksCount   int
	LastChunkSize int
}
