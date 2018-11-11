package core

type HLNode struct {
	hlTarget string
	permiss int //0 - SOMENTE LEITURA 1 - LEITURA E ESCRITA
}

func HLNodeCreate(hlTarget string, permiss int) (HLNode) {
		hlNode := HLNode{}
		hlNode = HLNode{hlTarget, permiss}
		return hlNode
}

func (hlNode HLNode) GetHLTarget() (string) {
		return hlNode.hlTarget
}

func (hlNode HLNode) GetPermiss() (int) {
		return hlNode.permiss
}
