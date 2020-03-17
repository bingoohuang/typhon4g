
```java
// com.ctrip.framework.apollo.core.MetaDomainConsts

* Select one available meta server from the comma separated meta server addresses, e.g.
* http://1.2.3.4:8080,http://2.3.4.5:8080

// random load balancing
Collections.shuffle(metaServers);

// check whether /services/config is accessible
if (NetUtil.pingUrl(metraAddress + "/services/config")) 

// com.ctrip.framework.apollo.internals.ConfigServiceLocator

private String assembleMetaServiceUrl() {
    String domainName = m_configUtil.getMetaServerDomainName();
    String appId = m_configUtil.getAppId();
    String localIp = m_configUtil.getLocalIp();

    Map<String, String> queryParams = Maps.newHashMap();
    queryParams.put("appId", queryParamEscaper.escape(appId));
    if (!Strings.isNullOrEmpty(localIp)) {
      queryParams.put("ip", queryParamEscaper.escape(localIp));
    }

    return domainName + "/services/config?" + MAP_JOINER.join(queryParams);
}

HttpResponse<List<ServiceDTO>> response = m_httpUtil.doGet(request, m_responseType);

public class ServiceDTO {
  private String appName;
  private String instanceId;
  private String homepageUrl;
}

// com.ctrip.framework.apollo.internals.RemoteConfigLongPollService

url =
    assembleLongPollRefreshUrl(lastServiceDto.getHomepageUrl(), appId, cluster, dataCenter,
        m_notifications);

// ...

final HttpResponse<List<ApolloConfigNotification>> response =
    m_httpUtil.doGet(request, m_responseType);

public class ApolloConfigNotification {
  private String namespaceName;
  private long notificationId;
  private volatile ApolloNotificationMessages messages;
}

public class ApolloNotificationMessages {
  private Map<String, Long> details;
}

// com.ctrip.framework.apollo.internals.RemoteConfigLongPollService
String assembleLongPollRefreshUrl(String uri, String appId, String cluster, String dataCenter,
                                    Map<String, Long> notificationsMap) {
    Map<String, String> queryParams = Maps.newHashMap();
    queryParams.put("appId", queryParamEscaper.escape(appId));
    queryParams.put("cluster", queryParamEscaper.escape(cluster));
    queryParams
        .put("notifications", queryParamEscaper.escape(assembleNotifications(notificationsMap)));

    if (!Strings.isNullOrEmpty(dataCenter)) {
      queryParams.put("dataCenter", queryParamEscaper.escape(dataCenter));
    }
    String localIp = m_configUtil.getLocalIp();
    if (!Strings.isNullOrEmpty(localIp)) {
      queryParams.put("ip", queryParamEscaper.escape(localIp));
    }

    String params = MAP_JOINER.join(queryParams);
    if (!uri.endsWith("/")) {
      uri += "/";
    }

    return uri + "notifications/v2?" + params;
}


// com.ctrip.framework.apollo.internals.RemoteConfigRepository

HttpResponse<ApolloConfig> response = m_httpUtil.doGet(request, ApolloConfig.class);

public class ApolloConfig {
  private String appId;
  private String cluster;
  private String namespaceName;
  private Map<String, String> configurations;
  private String releaseKey;
}

String assembleQueryConfigUrl(String uri, String appId, String cluster, String namespace,
                                String dataCenter, ApolloNotificationMessages remoteMessages, ApolloConfig previousConfig) {

    String path = "configs/%s/%s/%s";
    List<String> pathParams =
        Lists.newArrayList(pathEscaper.escape(appId), pathEscaper.escape(cluster),
            pathEscaper.escape(namespace));
    Map<String, String> queryParams = Maps.newHashMap();

    if (previousConfig != null) {
      queryParams.put("releaseKey", queryParamEscaper.escape(previousConfig.getReleaseKey()));
    }

    if (!Strings.isNullOrEmpty(dataCenter)) {
      queryParams.put("dataCenter", queryParamEscaper.escape(dataCenter));
    }

    String localIp = m_configUtil.getLocalIp();
    if (!Strings.isNullOrEmpty(localIp)) {
      queryParams.put("ip", queryParamEscaper.escape(localIp));
    }

    if (remoteMessages != null) {
      queryParams.put("messages", queryParamEscaper.escape(gson.toJson(remoteMessages)));
    }

    String pathExpanded = String.format(path, pathParams.toArray());

    if (!queryParams.isEmpty()) {
      pathExpanded += "?" + MAP_JOINER.join(queryParams);
    }
    if (!uri.endsWith("/")) {
      uri += "/";
    }
    return uri + pathExpanded;
}
```

