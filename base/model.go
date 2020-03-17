package base

import (
	"time"
)

// FileRawWait structured the config file content detail.
type FileRawWait struct {
	Raw  FileRaw
	Wait chan bool
}

// FileRaw structured the config file content detail.
type FileRaw struct {
	AppID    string `json:"appID"`
	ConfFile string `json:"confFile"`
	Content  string `json:"content"`
	Crc      string `json:"crc"`

	TriggerChange bool
}

// FileContent structured the config file content detail.
type FileContent struct {
	FileRaw

	conf ConfFile
}

func (f *FileContent) init() {
	f.conf = NewConfFile(f.ConfFile, f.Content)
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
