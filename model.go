package typhon4g

import "time"

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
	Properties ConfFmt = iota
	TXT
)

type ConfFileChangeEvent struct {
	ConfFile       string    `json:"confFile"`
	ConfFileFormat ConfFmt   `json:"confFileFormat"`
	Old            string    `json:"old"`
	Current        string    `json:"current"`
	ChangedTime    time.Time `json:"changedTime"` // 变更发生的时间(毫秒）
}
