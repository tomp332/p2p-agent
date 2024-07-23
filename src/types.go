package src

const (
	FilesNodeType NodeType = iota
	UnknownNodeType
)

type NodeType int

func (nt NodeType) String() string {
	switch nt {
	case FilesNodeType:
		return "FilesNode"
	default:
		return "Unknown"
	}
}
