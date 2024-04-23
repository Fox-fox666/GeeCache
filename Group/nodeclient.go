package Group

import cachepb_pb "GeeCache/cachepb"

type NodePicker interface {
	PickNode(key string) (NodeClient, bool)
}

type NodeClient interface {
	Get(request *cachepb_pb.Request) (*cachepb_pb.Response, error)
}
