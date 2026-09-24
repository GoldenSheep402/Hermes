
# Hermes

赫耳墨斯，来自希腊神话中的信使神。基于 jframe 编写的 PT 站，追求高性能与极致体验。

## 开发与 CI

- **测试**: `go test ./...`
- **依赖校验**: `go mod verify`
- **漏洞扫描**（需安装 [govulncheck](https://go.dev/blog/govulncheck)）: `govulncheck ./...`
- **一键脚本**: [`scripts/ci.sh`](scripts/ci.sh) — 依次执行 `verify`、全量测试、`govulncheck`（若已安装）

Tracker 行为、Redis Key、部署与安全说明见 [`mod/tracker/README.md`](mod/tracker/README.md)。

# Powered by jFrame
![jFrame](https://github.com/GoldenSheep402/Hermes/raw/main/docs/header.webp)
> 一个创意无限的 Golang 框架

- 更好的权限控制
- 加密算法
- 下周三前写好大纲