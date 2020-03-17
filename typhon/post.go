package typhon

import (
	"errors"
	"strings"

	"github.com/bingoohuang/gonet"
	"github.com/bingoohuang/gor"
	"github.com/bingoohuang/gou/enc"
	"github.com/bingoohuang/gou/str"
	"github.com/bingoohuang/typhon4g/base"
	"github.com/sirupsen/logrus"
)

// ClientReportRsp defines the top response structure of client report querying.
type ClientReportRsp struct {
	RspHead
	Data []base.ClientReportRspItem `json:"data"`
}

// ClientReport defines the structure of client report uploading.
type ClientReport struct {
	Host string `json:"host"`
	Pid  string `json:"pid"`
	Bin  string `json:"bin"`

	Items []base.ClientReportItem `json:"items"`
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

// UploadReport uploads the listeners reports.
func (c *Client) UploadReport(report *ClientReport) {
	_, _ = gor.IterateSlice(c.C.ConfigServers, -1, func(url string) bool {
		return c.TryUploadReport(url, report)
	})
}

// basicAuth adds basic auth the the http request.
func (c *Client) basicAuth(r *gonet.HTTPReq) {
	if c.C.PostAuth != "" {
		usr, pas := str.Split2(c.C.PostAuth, ":", true, false)
		r.BasicAuth(usr, pas)
	}
}

// TryUploadReport try to upload report.
func (c *Client) TryUploadReport(url string, report *ClientReport) bool {
	var rsp RspHead

	reportURL := strings.Replace(url, "/client/config/", "/admin/report/", 1)
	rspBody, err := c.C.Req.RestPostFn(reportURL, report, &rsp, c.basicAuth)
	logrus.Infof("report response %s, error %v", string(rspBody), err)

	return rsp.Status == 200
}

// PostConf posts the conf to the server with clientIps (blank/comma separated IP addresses or all)
// returns crc and error info.
func (c *Client) PostConf(confFile, raw, clientIps string) (string, error) {
	ok, res := gor.IterateSlice(c.C.ConfigServersParsed, -1, func(addr string) (bool, interface{}) {
		return c.TryPost(addr, confFile, raw, clientIps)
	})

	if ok && res != nil {
		c.fileRaw <- base.FileRawWait{
			FileRaw: base.FileRaw{
				AppID:    c.C.AppID,
				ConfFile: confFile,
				Content:  raw,
				Crc:      enc.Checksum([]byte(raw)),
			}}

		return res.(string), nil
	}

	return "", errors.New("failed to post")
}

// TryPost try to post conf the server in specified url.
func (c *Client) TryPost(addr, confFile, raw, clientIps string) (bool, interface{}) {
	addr = c.createConfigServer(addr)
	postURL := strings.Replace(addr, "/client/config/", "/admin/release/", 1) +
		"/" + confFile + "?clientIps=" + clientIps

	var rsp PostRsp

	_, _ = c.C.Req.RestPostFn(postURL, ReqBody{Data: raw}, &rsp, c.basicAuth)

	return rsp.Status == 200, rsp.Crc
}

// ListenerResults gets the listener results from the server.
func (c *Client) ListenerResults(confFile, crc string) ([]base.ClientReportRspItem, error) {
	if ok, res := gor.IterateSlice(c.C.ConfigServers, -1, func(url string) (bool, interface{}) {
		return c.tryListenerResults(url, confFile, crc)
	}); ok {
		return res.([]base.ClientReportRspItem), nil
	}

	return nil, errors.New("failed to ListenerResults")
}

// TryListenerResults tries to get listeners results from the server specified url.
func (c *Client) tryListenerResults(url, confFile, crc string) (bool, interface{}) {
	reportURL := strings.Replace(url, "/config/", "/report/", 1) + "?confFile=" + confFile + "&crc=" + crc

	var rsp ClientReportRsp
	err := c.C.Req.RestGet(reportURL, &rsp)

	if err != nil || rsp.Status != 200 {
		logrus.Warnf("fail to TryListenerResults %s, error %v", reportURL, err)
		return false, nil
	}

	return true, rsp.Data
}
