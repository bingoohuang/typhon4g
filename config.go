package typhon4g

import (
	"github.com/bingoohuang/gou"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type ConfigService struct {
	c        *TyphonContext
	setting  gou.UrlHttpSettings
	updateFn func([]FileContent)
}

type ConfigRsp struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    []FileContent `json:"data"`
}

func (c ConfigService) start() {
	c.try("")

	d := time.Duration(c.c.ConfigRefreshIntervalSeconds) * time.Second
	timer := time.NewTimer(d)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			c.try("")
			timer.Reset(d)
		}
	}
}

func (c ConfigService) try(confFile string) bool {
	urls := c.c.ConfigServerUrls
	ok, _ := gou.IterateSlice(urls, gou.RandomIntN(uint64(len(urls))), func(url string) bool {
		return c.tryUrl(url, confFile, &c.setting)
	})

	return ok
}

func (c ConfigService) tryUrl(url, confFile string, setting *gou.UrlHttpSettings) bool {
	confFileCrc := ""
	if confFile != "" {
		confFileCrc = confFile + ":0"
	} else {
		confFileCrc = c.createConfFileCrcs()
	}
	if confFileCrc == "" {
		return true
	}

	clientUrl := url + "?confFileCrc=" + confFileCrc
	var rsp ConfigRsp
	err := gou.RestGetV2(clientUrl, &rsp, setting)
	if err != nil {
		logrus.Warnf("fail to RefreshConfig %s, error %v", clientUrl, err)
		return false
	}

	if len(rsp.Data) > 0 {
		c.updateFn(rsp.Data)
		report := c.c.saveCaches(rsp.Data)
		if report != nil {
			c.uploadReport(report)
		}
		return true
	}

	return false
}

func (c ConfigService) createConfFileCrcs() string {
	confFileCrc := make([]string, 0)
	c.c.iterateCache(func(cf string, fc *FileContent) {
		confFileCrc = append(confFileCrc, fc.ConfFile+":"+fc.Crc)
	})

	return strings.Join(confFileCrc, ",")
}

func (c ConfigService) uploadReport(report *ClientReport) {
	urls := c.c.ConfigServerUrls
	_, _ = gou.IterateSlice(urls, gou.RandomIntN(uint64(len(urls))), func(url string) bool {
		return c.tryUpload(url, report)
	})
}

func (c ConfigService) tryUpload(url string, report *ClientReport) bool {
	reportUrl := strings.Replace(url, "/config/", "/report/", 1)
	var rsp RspBase
	rspBody, err := gou.RestPost(reportUrl, report, &rsp)
	logrus.Infof("report response %s, error %v", string(rspBody), err)

	return rsp.Status == 200
}

func (c ConfigService) PostConf(confFile, raw, clientIps string) (string, error) {
	urls := c.c.ConfigServerUrls
	ok, res := gou.IterateSlice(urls, gou.RandomIntN(uint64(len(urls))), func(url string) (bool, interface{}) {
		return c.tryPost(url, confFile, raw, clientIps)
	})

	if ok && res != nil {
		return res.(string), nil
	}

	return "", errors.New("failed to post")
}

func (c ConfigService) tryPost(url, confFile, raw, clientIps string) (bool, interface{}) {
	postUrl := strings.Replace(url, "/client/config/", "/admin/release/", 1) + "/" + confFile + "?clientIps=" + clientIps
	req := ReqBody{Data: raw}
	var rsp Rsp

	_, _ = gou.RestPostFn(postUrl, req, &rsp, func(r *gou.UrlHttpRequest) error {
		if c.c.PostAuth != "" {
			usr, pas := gou.Split2(c.c.PostAuth, ":", true, false)
			r.SetBasicAuth(usr, pas)
		}
		return nil
	})

	return rsp.Status == 200, rsp.Crc
}

func (c ConfigService) GetListenerResults(confFile, crc string) ([]ClientReportRspItem, error) {
	urls := c.c.ConfigServerUrls
	ok, res := gou.IterateSlice(urls, gou.RandomIntN(uint64(len(urls))), func(url string) (bool, interface{}) {
		return c.tryGetListenerResults(url, confFile, crc)
	})

	if ok && res != nil {
		return res.([]ClientReportRspItem), nil
	}

	return nil, errors.New("failed to GetListenerResults")
}

func (c ConfigService) tryGetListenerResults(url, confFile, crc string) (bool, interface{}) {
	reportUrl := strings.Replace(url, "/config/", "/report/", 1) + "?confFile=" + confFile + "&crc=" + crc

	var rsp ClientReportRsp
	err := gou.RestGetV2(reportUrl, &rsp, &c.setting)
	if err != nil || rsp.Status != 200 {
		logrus.Warnf("fail to tryGetListenerResults %s, error %v", reportUrl, err)
		return false, nil
	}

	return true, rsp.Data
}
