package Group

import (
	"GeeCache/ConsistentHash"
	cachepb_pb "GeeCache/cachepb"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/geecache/"
	defaultreplices = 50
)

type HttpNode struct {
	addr        string
	basePath    string
	mu          sync.Mutex
	nodes       *ConsistentHash.NodeMap
	nodeGetters map[string]*NodeGetter
}

// http://127.0.0.1:9898/geecache/<group name>/<key>
func (h *HttpNode) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if !strings.HasPrefix(path, defaultBasePath) {
		w.Write([]byte("Path error"))
		return
	}

	log.Printf("%s %s\n", r.Method, path)
	splitPath := strings.Split(path[len(defaultBasePath):], "/")
	groupname := splitPath[0]
	key := splitPath[1]
	group := GetGroup(groupname)
	if group == nil {
		w.Write([]byte("Path error,group no exist"))
		return
	}

	data, err := group.Get(key)
	if err != nil {
		log.Println(err)
		w.Write([]byte("net error"))
		return
	}

	response := &cachepb_pb.Response{Data: data.ToSilce()}
	proto, err := proto.Marshal(response)
	//w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(proto)
}

func NewhttpNode(addr string) *HttpNode {
	return &HttpNode{
		addr:     addr,
		basePath: defaultBasePath,
	}
}

func (h *HttpNode) PickNode(key string) (NodeClient, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if nodeaddr := h.nodes.Get(key); nodeaddr != "" && nodeaddr != h.addr {
		log.Printf("Peer pick %s\n", nodeaddr)
		return h.nodeGetters[nodeaddr], true
	}

	return nil, false
}

func (h *HttpNode) Set(addr ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.nodes = ConsistentHash.NewNodeMap(nil, defaultreplices)
	h.nodes.Add(addr...)
	h.nodeGetters = make(map[string]*NodeGetter)
	for _, s := range addr {
		h.nodeGetters[s] = &NodeGetter{baseURL: s + defaultBasePath}
	}
}
