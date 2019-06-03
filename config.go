package typhon4g

import (
	"strings"
	"time"

	"github.com/bingoohuang/gou"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ConfigService defines the structure used for config refresh/write service
type ConfigService struct {
	C        *TyphonContext
	Setting  gou.UrlHttpSettings
	UpdateFn func([]FileContent)
}

// Start starts the refreshing loop of config service.
func (c ConfigService) Start() {
	c.Try("")

	d := SecondsDuration(c.C.ConfigRefreshIntervalSeconds)
	timer := time.NewTimer(d)
	defer timer.Stop()

	for range timer.C {
		c.Try("")
		timer.Reset(d)
	}
}

// Try tries to refresh conf defined by confFile or all (confFile is empty).
func (c ConfigService) Try(confFile string) (bool, ConfFile) {
	hit, cf := gou.RandomIterateSlice(c.C.ConfigServers, func(url string) (bool, interface{}) {
		return c.TryURL(url, confFile, &c.Setting)
	})

	if cf != nil {
		return hit, cf.(ConfFile)
	}

	return hit, nil
}

// TryURL tries to refresh conf defined by confFile or all (confFile is empty) in specified URL.
func (c ConfigService) TryURL(url, confFile string, setting *gou.UrlHttpSettings) (bool, interface{}) {
	confFileCrc := ""
	if confFile != "" {
		confFileCrc = confFile + ":0"
	} else {
		confFileCrc = c.CreateConfFileCrcs()
	}
	if confFileCrc == "" {
		return true, nil
	}

	clientURL := url + "?confFileCrc=" + confFileCrc
	var rsp struct {
		Status  int           `json:"status"`
		Message string        `json:"message"`
		Data    []FileContent `json:"data"`
	}

	err := gou.RestGetV2(clientURL, &rsp, setting)
	if err != nil {
		logrus.Warnf("fail to RefreshConfig %s, error %v", clientURL, err)
		return false, nil
	}

	if len(rsp.Data) == 0 {
		return false, nil
	}

	c.UpdateFn(rsp.Data)
	report := c.C.SaveFileContents(rsp.Data, true)
	if report != nil && len(report.Items) != 0 {
		c.UploadReport(report)
	}

	if confFile == "" {
		return true, nil
	}

	return true, c.C.LoadConfFile(confFile)
}

// CreateConfFileCrcs creates conf files and their crcs.
func (c ConfigService) CreateConfFileCrcs() string {
	confFileCrc := make([]string, 0)
	c.C.WalkFileContents(func(cf string, fc *FileContent) {
		confFileCrc = append(confFileCrc, fc.ConfFile+":"+fc.Crc)
	})

	return strings.Join(confFileCrc, ",")
}

// UploadReport uploads the listeners reports.
func (c ConfigService) UploadReport(report *ClientReport) {
	_, _ = gou.RandomIterateSlice(c.C.ConfigServers, func(url string) bool {
		return c.TryUploadReport(url, report)
	})
}

// SetBasicAuth adds basic auth the the http request.
func (c ConfigService) SetBasicAuth(r *gou.UrlHttpRequest) error {
	if c.C.postAuth != "" {
		usr, pas := gou.Split2(c.C.postAuth, ":", true, false)
		r.SetBasicAuth(usr, pas)
	}
	return nil
}

// TryUploadReport try to upload report.
func (c ConfigService) TryUploadReport(url string, report *ClientReport) bool {
	var rsp RspHead
	reportURL := strings.Replace(url, "/client/config/", "/admin/report/", 1)
	rspBody, err := gou.RestPostFn(reportURL, report, &rsp, c.SetBasicAuth)
	logrus.Infof("report response %s, error %v", string(rspBody), err)

	return rsp.Status == 200
}

// PostConf posts the conf to the server, returns crc and error info.
func (c ConfigService) PostConf(confFile, raw, clientIps string) (string, error) {
	ok, res := gou.RandomIterateSlice(c.C.ConfigServers, func(url string) (bool, interface{}) {
		return c.TryPost(url, confFile, raw, clientIps)
	})

	if ok && res != nil {
		return res.(string), nil
	}

	return "", errors.New("failed to post")
}

// TryPost try to post conf the server in specified url.
func (c ConfigService) TryPost(url, confFile, raw, clientIps string) (bool, interface{}) {
	postURL := strings.Replace(url, "/client/config/", "/admin/release/", 1) +
		"/" + confFile + "?clientIps=" + clientIps
	var rsp PostRsp

	_, _ = gou.RestPostFn(postURL, ReqBody{Data: raw}, &rsp, c.SetBasicAuth)
	ok := rsp.Status == 200

	if ok {
		fileContent := []FileContent{{AppID: c.C.AppID, ConfFile: confFile, Content: raw, Crc: rsp.Crc}}
		fileContent[0].init()
		c.UpdateFn(fileContent)
		c.C.SaveFileContents(fileContent, false)
	}
	return ok, rsp.Crc
}

// ListenerResults gets the listener results from the server.
func (c ConfigService) ListenerResults(confFile, crc string) ([]ClientReportRspItem, error) {
	ok, res := gou.RandomIterateSlice(c.C.ConfigServers, func(url string) (bool, interface{}) {
		return c.TryListenerResults(url, confFile, crc)
	})

	if ok {
		return res.([]ClientReportRspItem), nil
	}

	return nil, errors.New("failed to ListenerResults")
}

// TryListenerResults tries to get listeners results from the server specified url.
func (c ConfigService) TryListenerResults(url, confFile, crc string) (bool, interface{}) {
	reportURL := strings.Replace(url, "/config/", "/report/", 1) + "?confFile=" + confFile + "&crc=" + crc

	var rsp ClientReportRsp
	err := gou.RestGetV2(reportURL, &rsp, &c.Setting)
	if err != nil || rsp.Status != 200 {
		logrus.Warnf("fail to TryListenerResults %s, error %v", reportURL, err)
		return false, nil
	}

	return true, rsp.Data
}
