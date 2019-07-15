package typhon4g

import (
	"time"
)

// FileContent structured the config file content detail.
type FileContent struct {
	AppID    string `json:"appID"`
	ConfFile string `json:"confFile"`
	Content  string `json:"content"`
	Crc      string `json:"crc"`

	conf ConfFile
}

func (f *FileContent) init() {
	f.conf = NewConfFile(f.ConfFile, f.Content)
}

func (f *FileContent) update(u FileContent) {
	f.Content = u.Content
	f.Crc = u.Crc
}

// ConfFileChangeEvent structured the change event content.
type ConfFileChangeEvent struct {
	ConfFile       string    `json:"confFile"`
	ConfFileFormat ConfFmt   `json:"confFileFormat"`
	Old            string    `json:"old"`
	Current        string    `json:"current"`
	ChangedTime    time.Time `json:"changedTime"` // 变更发生的时间(毫秒）
}

// ClientReportItem defines the structure of client listener report item.
type ClientReportItem struct {
	Time     string `json:"time"`
	Msg      string `json:"msg"`
	Ok       bool   `json:"ok"`
	ConfFile string `json:"confFile"`
	Crc      string `json:"crc"`
}

// ClientReportRspItem defines the structure of response to client report querying.
type ClientReportRspItem struct {
	ID    string `json:"id"`
	AppID string `json:"appID"`
	Host  string `json:"host"`
	Pid   string `json:"pid"`
	Bin   string `json:"bin"`
	ClientReportItem
}

// ClientReportRsp defines the top response structure of client report querying.
type ClientReportRsp struct {
	RspHead
	Data []ClientReportRspItem `json:"data"`
}

// ClientReport defines the structure of client report uploading.
type ClientReport struct {
	Host string `json:"host"`
	Pid  string `json:"pid"`
	Bin  string `json:"bin"`

	Items []ClientReportItem `json:"items"`
}

// RspHead defines the head of response.
type RspHead struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// ReqBody defines the request body.
type ReqBody struct {
	Data interface{} `json:"data"`
}

// PostRsp defines the response the post api.
type PostRsp struct {
	RspHead
	Crc string `json:"crc"`
}
