package Group

import (
	cachepb_pb "GeeCache/cachepb"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"net/url"
)

type NodeGetter struct {
	baseURL string
}

// http://127.0.0.1:9090/geecache/<group name>/<key>
func (ng *NodeGetter) Get(gogo *cachepb_pb.Request) (*cachepb_pb.Response, error){
	url := fmt.Sprintf("%v%v/%v", ng.baseURL, url.QueryEscape(gogo.GetGroupname()), url.QueryEscape(gogo.GetKey()))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", resp.Status)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	come:=&cachepb_pb.Response{}
	err = proto.Unmarshal(bytes, come)
	if err != nil {
		return nil, fmt.Errorf("proto.Unmarshal: %v", err)
	}
	return come, nil
}
