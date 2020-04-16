package apollo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/bingoohuang/typhon4g/base"
)

// notification defines the polling request and response structure of apollo polling service.
type notification struct {
	Namespace string `json:"namespaceName,omitempty"`
	NotifyID  int64  `json:"notificationId,omitempty"`
}

// Polling polling the specified addr.
func (c *Client) Polling(configServer string) error {
	notifications, namespaceMap := c.makeNotifications()
	if len(notifications) == 0 {
		logrus.Infof("nothing to polling")
		return errors.New("nothing to polling")
	}

	pollingAddr := c.pollingAddr(configServer, notifications)

	logrus.Infof("polling addr %s", pollingAddr)

	var rsp []notification
	if err := c.ReqPoll.RestGet(pollingAddr, &rsp); err != nil {
		isTimeout := os.IsTimeout(err)
		logrus.Warnf("normal polling %s, isTimeout %v, error %v", pollingAddr, isTimeout, err)

		if isTimeout {
			return nil
		}

		return err
	}

	for _, u := range rsp {
		fileName := namespaceMap[u.Namespace]

		c.notifications.Store(fileName, u.NotifyID)
		c.readConfig(fileName, nil)
	}

	return nil
}

func (c *Client) makeNotifications() ([]notification, map[string]string) {
	notifications := make([]notification, 0)

	m := make(map[string]string)

	c.notifications.Range(func(k, v interface{}) bool {
		f := k.(string)

		m[strings.TrimSuffix(f, path.Ext(f))] = f
		m[f] = f

		notifications = append(notifications, notification{Namespace: f, NotifyID: v.(int64)})

		return true
	})

	return notifications, m
}

func (c *Client) pollingAddr(configServer string, notifications []notification) string {
	notifyBytes, _ := json.Marshal(notifications)

	return fmt.Sprintf("%s/notifications/v2?appId=%s&cluster=%s&dataCenter=%s&notifications=%s&ip=%s",
		base.HTPPAddr(configServer),
		url.QueryEscape(c.AppID),
		url.QueryEscape(c.Cluster),
		url.QueryEscape(c.DataCenter),
		url.QueryEscape(string(notifyBytes)),
		url.QueryEscape(c.localIP),
	)
}
