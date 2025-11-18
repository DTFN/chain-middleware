package chaincode

type FileMeta struct {
	RawFileId      string   `json:"rawFileId"`
	Pk       	   string   `json:"pk"`
	Uploader       string   `json:"uploader"`
	FileName       string   `json:"fileName"`
	Author         string   `json:"author"`
	FileSize       uint64   `json:"fileSize"`
	TotalShards    uint64   `json:"totalShards"`
	DataShards     uint64   `json:"dataShards"`
	IpfsUrls       []string `json:"ipfsUrls"`
	Cids 		   []string `json:"cids"`
	Timestamp      int64    `json:"timestamp"`
}
