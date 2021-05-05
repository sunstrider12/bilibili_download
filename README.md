# 哔哩哔哩直播录制软件


## 0x01 使用方法说明
首先chrome打开live.bilibili.com,打开`开发者工具`->`NETWORK` 随便点开一个连接,找到request,然后找到`cookie` 把`cookie`复制到`yml`配置文件内.
如果不想登录账号,可使用默认的cookie
```
LIVE_BUVID=AUTO5016202281239115; _uuid=24E05F69-88C2-5566-9D3D-6891D9F6005223821infoc; buvid3=F0089EBB-EB4F-43C2-9412-66523928791E13424infoc; bfe_id=5112800f2e3d3cf17a473918472e345c; PVID=2
```
## 0x02 配置文件说明

| Dir        | 下载目录                    |
|------------|-------------------------|
| RoomNum    | 房间号                     |
| CheckTime  | 如未开播,多长时间检查一次           |
| SaveSpace  | 文件超过多少MB之后储存新的文件,0为不分割  |
| NeedTicker | 是否需要定时                  |
| BeginTime  | 开始检测时间24小时制,如下午1点就是1300 |
| EndTime    | 结束检测时间24小时制             |

可以多个直播间同时录制

## 0x03 TODO:
增加flv文件解析,不再重新请求录制
增加录制清晰度选项
