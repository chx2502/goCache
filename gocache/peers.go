package gocache

import pb "gocache/gocachepb"

// PeerPick 是定位某个特定 key 所在的节点所必须实现的接口
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter 是节点必须实现的接口
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}


