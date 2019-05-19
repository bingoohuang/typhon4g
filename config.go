package typhon4g

import (
	"github.com/bingoohuang/gou"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type ConfigService struct {
	C        *TyphonContext
	Setting  gou.UrlHttpSettings
	UpdateFn func([]FileContent)
}

type ConfigRsp struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    []FileContent `json:"data"`
}

func (c ConfigService) Start() {
	c.Try("")

	d := SecondsDuration(c.C.ConfigRefreshIntervalSeconds)
	timer := time.NewTimer(d)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			c.Try("")
			timer.Reset(d)
		}
	}
}

func (c ConfigService) Try(confFile string) (bool, ConfFile) {
	hit, cf := gou.RandomIterateSlice(c.C.ConfigServerUrls, func(url string) (bool, interface{}) {
		return c.TryUrl(url, confFile, &c.Setting)
	})

	if cf != nil {
		return hit, cf.(ConfFile)
	}

	return hit, nil
}

func (c ConfigService) TryUrl(url, confFile string, setting *gou.UrlHttpSettings) (bool, interface{}) {
	confFileCrc := ""
	if confFile != "" {
		confFileCrc = confFile + ":0"
	} else {
		confFileCrc = c.CreateConfFileCrcs()
	}
	if confFileCrc == "" {
		return true, nil
	}

	clientUrl := url + "?confFileCrc=" + confFileCrc
	var rsp ConfigRsp
	err := gou.RestGetV2(clientUrl, &rsp, setting)
	if err != nil {
		logrus.Warnf("fail to RefreshConfig %s, error %v", clientUrl, err)
		return false, nil
	}

	if len(rsp.Data) > 0 {
		c.UpdateFn(rsp.Data)
		report := c.C.SaveFileContents(rsp.Data)
		if report != nil && len(report.Items) != 0 {
			c.UploadReport(report)
		}

		if confFile == "" {
			return true, nil
		} else {
			return true, c.C.LoadConfFile(confFile)
		}
	}

	return false, nil
}

func (c ConfigService) CreateConfFileCrcs() string {
	confFileCrc := make([]string, 0)
	c.C.WalkFileContents(func(cf string, fc *FileContent) {
		confFileCrc = append(confFileCrc, fc.ConfFile+":"+fc.Crc)
	})

	return strings.Join(confFileCrc, ",")
}

func (c ConfigService) UploadReport(report *ClientReport) {
	_, _ = gou.RandomIterateSlice(c.C.ConfigServerUrls, func(url string) bool {
		return c.TryUploadReport(url, report)
	})
}

func (c ConfigService) SetBasicAuth(r *gou.UrlHttpRequest) error {
	if c.C.PostAuth != "" {
		usr, pas := gou.Split2(c.C.PostAuth, ":", true, false)
		r.SetBasicAuth(usr, pas)
	}
	return nil
}

func (c ConfigService) TryUploadReport(url string, report *ClientReport) bool {
	var rsp RspBase
	reportUrl := strings.Replace(url, "/client/config/", "/admin/report/", 1)
	rspBody, err := gou.RestPostFn(reportUrl, report, &rsp, c.SetBasicAuth)
	logrus.Infof("report response %s, error %v", string(rspBody), err)

	return rsp.Status == 200
}

func (c ConfigService) PostConf(confFile, raw, clientIps string) (string, error) {
	ok, res := gou.RandomIterateSlice(c.C.ConfigServerUrls, func(url string) (bool, interface{}) {
		return c.TryPost(url, confFile, raw, clientIps)
	})

	if ok && res != nil {
		return res.(string), nil
	}

	return "", errors.New("failed to post")
}

func (c ConfigService) TryPost(url, confFile, raw, clientIps string) (bool, interface{}) {
	postUrl := strings.Replace(url, "/client/config/", "/admin/release/", 1) +
		"/" + confFile + "?clientIps=" + clientIps
	req := ReqBody{Data: raw}
	var rsp Rsp

	_, _ = gou.RestPostFn(postUrl, req, &rsp, c.SetBasicAuth)
	return rsp.Status == 200, rsp.Crc
}

func (c ConfigService) ListenerResults(confFile, crc string) ([]ClientReportRspItem, error) {
	ok, res := gou.RandomIterateSlice(c.C.ConfigServerUrls, func(url string) (bool, interface{}) {
		return c.TryGetListenerResults(url, confFile, crc)
	})

	if ok {
		return res.([]ClientReportRspItem), nil
	}

	return nil, errors.New("failed to ListenerResults")
}

func (c ConfigService) TryGetListenerResults(url, confFile, crc string) (bool, interface{}) {
	reportUrl := strings.Replace(url, "/config/", "/report/", 1) + "?confFile=" + confFile + "&crc=" + crc

	var rsp ClientReportRsp
	err := gou.RestGetV2(reportUrl, &rsp, &c.Setting)
	if err != nil || rsp.Status != 200 {
		logrus.Warnf("fail to TryGetListenerResults %s, error %v", reportUrl, err)
		return false, nil
	}

	return true, rsp.Data
}
