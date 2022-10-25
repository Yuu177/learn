# HTTP & HTTPS

- HTTPS

HTTPS （全称：Hyper Text Transfer Protocol over SecureSocket Layer），是以安全为目标的 HTTP 通道，在  HTTP 的基础上通过传输加密和身份认证保证了传输过程的安全性。

HTTPS 在 HTTP 的基础下加入 SSL，HTTPS  的安全基础是 SSL，因此加密的详细内容就需要 SSL。

HTTPS 主要由两部分组成：HTTP + SSL / TLS，也就是在 HTTP 上又加了一层处理加密信息的模块。服务端和客户端的信息传输都会通过 TLS 进行加密，所以传输的数据都是加密后的数据。

假如没有 https，访问网站时候传输数据都是未被加密的，在这种情况下数据就存在安全隐患。

- 为什么https 不支持支持 IP 和端口的形式访问

因为 https 申请的证书是和域名进行绑定的，直接使用 https + ip:port 地址证书就会校验失败，会造成用户无法访问。

