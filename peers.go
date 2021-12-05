package geecache

// PeerPicker 实现获取相对应的节点
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter 通过组节点和 key 获取缓存值，相当与 HTTP 客户端
type PeerGetter interface {
	Get(group, key string) ([]byte, error)
}
