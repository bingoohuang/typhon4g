package typhon4g

import (
	"time"
)

type FileContent struct {
	AppID    string `json:"appID"`
	ConfFile string `json:"confFile"`
	Content  string `json:"content"`
	Crc      string `json:"crc"`

	Conf ConfFile `json:"-"`
}

func (f *FileContent) init() {
	f.Conf = NewConfFile(f.ConfFile, f.Content)
}

type ConfFmt int

const (
	PropertiesFmt ConfFmt = iota
	TxtFmt
)

type ConfFileChangeEvent struct {
	ConfFile       string    `json:"confFile"`
	ConfFileFormat ConfFmt   `json:"confFileFormat"`
	Old            string    `json:"old"`
	Current        string    `json:"current"`
	ChangedTime    time.Time `json:"changedTime"` // 变更发生的时间(毫秒）
}

type ClientReportItem struct {
	Time     string `json:"time"`
	Msg      string `json:"msg"`
	Ok       bool   `json:"ok"`
	ConfFile string `json:"confFile"`
	Crc      string `json:"crc"`
}

type ClientReportRspItem struct {
	Id       string `json:"id"`
	AppID    string `json:"appID"`
	Host     string `json:"host"`
	Pid      string `json:"pid"`
	Bin      string `json:"bin"`
	ConfFile string `json:"confFile"`
	Crc      string `json:"crc"`
	ClientReportItem
}

type ClientReportRsp struct {
	RspBase
	Data []ClientReportRspItem `json:"data"`
}

type ClientReport struct {
	Host string `json:"host"`
	Pid  string `json:"pid"`
	Bin  string `json:"bin"`

	Items []ClientReportItem `json:"items"`
}

type RspBase struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ReqBody struct {
	Data string `json:"data"`
}

type Rsp struct {
	RspBase
	Crc string `json:"crc"`
}
