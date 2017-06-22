notifier
=============
## 功能
-  把smtp封装为一个简单http接口，配置到sender中用来发送报警邮件， 使用方法

```
curl http://$ip:4000/api/sender/mail -d "tos=a@a.com,b@b.com&subject=xx&content=yy"
```
- 利用第三方短信接口提供一个http接口，发送报警短信
```
curl http://$ip:4000/api/sender/sms -d "tos=187322432&content=yy"
```
