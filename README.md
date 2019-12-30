# SSR Subscriber

A command-line tool for subscribe ssr config and check availability. 

基于 golang 的 SSR 订阅解析工具，并提供简单的端口和名称校验，保存输出为 ss-local 配置文件。

## 安装

```sh
# sudo pacman -S go
make
```

## 配置

ssr-subscriber 支持 文件、网络等订阅形式，可以参考代码库中的 config-example.yml 文件进行配置。示例配置文件如下：

```yaml
url:
  - <your-ssr-subscriber-url>

file:
  - <you-local-ssr-subscriber-file-base64-encoded>

check:
  timeout: 3s
  not: 免费|普通|回国|过期|剩余

output: <output-directory>

proxy: localhost:1088
```

其中，`check` 字段中 `timeout` 为检查对应 ssr 节点的主机和端口超时时间，`not` 字段为过滤和忽略 `Remarks` 中包含对应的字符串（支持正则表达式）。`proxy` 字段支持 socks5 的代理地址，~~强烈建议使用代理请求数据~~。

## 整合

建议配合 crontab 使用，定期抓取和更新订阅链接的节点信息，例如：

```
30 5 * * * rm -f $HOME/configs/*.json && ssr-subscriber -c $HOME/ssr-subscriber.yaml
```

## @TODO

- 代码优化和精简
- 支持其他格式的输出