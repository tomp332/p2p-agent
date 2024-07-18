package src

const (
	FileSystemNodeType NodeType = iota
	UnknownNodeType
)

type NodeType int

func (nt NodeType) String() string {
	switch nt {
	case FileSystemNodeType:
		return "FileSystemNode"
	default:
		return "Unknown"
	}
}
