# web-rtc-traps

使用WebRTC捕获真实IP，内网刺探

## 功能
- [x] 获取真实IP

## TODO
- [x] 内网刺探

## 编译


Mac 下编译 Windows 可执行文件
```bash
export GOOS=windows
export GOARCH=amd64
go build -ldflags="-s -w" -o web-rtc-traps.exe main.go
```