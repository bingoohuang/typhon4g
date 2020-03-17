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
	NamespaceName  string `json:"namespaceName,omitempty"`
	NotificationID int64  `json:"notificationId,omitempty"`
}

// Polling polling the specified addr.
func (c *Client) Polling(configServer string) error {
	notifications, filenameNamespaceMap := c.makeNotifications()
	if len(notifications) == 0 {
		logrus.Infof("nothing to polling")
		return errors.New("nothing to polling")
	}

	pollingAddr := c.pollingAddr(configServer, notifications)

	logrus.Infof("polling addr %s", pollingAddr)

	var rsp []notification

	if err := c.ReqPoll.RestGet(pollingAddr, &rsp); err != nil {
		if os.IsTimeout(err) {
			logrus.Infof("normal polling timeout %s", pollingAddr)
			return nil
		}

		return err
	}

	for _, u := range rsp {
		fileName := filenameNamespaceMap[u.NamespaceName]
		c.notifications.Store(fileName, u.NotificationID)
		c.readConfig(fileName, nil)
	}

	return nil
}

func (c *Client) makeNotifications() ([]notification, map[string]string) {
	notifications := make([]notification, 0)

	m := make(map[string]string)

	c.notifications.Range(func(k, v interface{}) bool {
		f := k.(string)
		namespace := strings.TrimSuffix(f, path.Ext(f))
		m[namespace] = f
		notifications = append(notifications, notification{
			NamespaceName:  namespace,
			NotificationID: v.(int64),
		})
		return true
	})

	return notifications, m
}

func (c *Client) pollingAddr(configServer string, notifications []notification) string {
	notificationBytes, _ := json.Marshal(notifications)
	json := string(notificationBytes)

	return fmt.Sprintf("%s/notifications/v2?appId=%s&cluster=%s&dataCenter=%s&notifications=%s&ip=%s",
		base.HTPPAddr(configServer),
		url.QueryEscape(c.AppID),
		url.QueryEscape(c.Cluster),
		url.QueryEscape(c.DataCenter),
		url.QueryEscape(json),
		url.QueryEscape(c.localIP),
	)
}
