.PHONY: all build test vet fmt run clean cross testdata uitest layout layout-sweep touch-sweep touch-modals cdp run-auth settings-check users-check users-deploy-check

BINARY := gofs
VERSION ?= $(shell git describe --tags --always 2>/dev/null || echo dev)

# 前端 UI 测试（jsdom）用的路径，可用环境变量覆盖
NODE_MODULES ?= $(HOME)/.workbuddy/binaries/node/workspace/node_modules
NODE_BIN     ?= $(shell ls $(HOME)/.workbuddy/binaries/node/versions/*/bin/node 2>/dev/null | tail -1 || echo node)
GOFS_URL     ?= http://127.0.0.1:5000
CHROME       ?= /Applications/Google Chrome.app/Contents/MacOS/Google Chrome
CDP_PORT     ?= 9333
MEASURE_URL  ?= http://127.0.0.1:5001
export NODE_MODULES NODE_BIN GOFS_URL CDP_PORT

all: fmt vet test build

## build: 构建当前平台的二进制
build:
	go build -trimpath -ldflags "-s -w" -o $(BINARY) .
	@ls -lh $(BINARY)

## test: 运行全部测试
test:
	go test ./...

## testv: 运行全部测试并查看每个用例
testv:
	go test ./internal/server/ -v

## uitest: 用 jsdom 跑真实页面交互（需先启动服务，且开启鉴权）
##   用法：另开一个终端跑 `make run-auth`，然后 `make uitest`
uitest:
	NODE_PATH=$(NODE_MODULES) $(NODE_BIN) scripts/test-ui.mjs $(GOFS_URL)

## run-auth: 以鉴权模式启动，供 uitest / settings-check / users-check 使用
##   admin 可写根目录；test 只读，用来验证「改设置需要根目录写权限」
##   ⚠️ 显式指定 --user-file 指向临时文件：否则测试会写进你真实的
##      <用户配置目录>/gofs/users.json，那里面可能有你在用的账号。
run-auth: build
	./$(BINARY) -A -b 127.0.0.1 -p 5000 -a 'admin:secret@/:rw' -a 'test:test@/' \
	  --user-file /tmp/gofs-run-auth-users.json ./data

## settings-check: 端到端校验「服务设置」（权限 / 范围约束 / CSRF / 上限生效）
##   用法：另开一个终端跑 `make run-auth`，然后 `make settings-check`
settings-check:
	./scripts/check-settings.sh

## users-check: 端到端校验用户管理（CRUD / 权限生效 / 改密后旧密码立刻失效 / 落盘）
##   用法：另开一个终端跑 `make run-auth`，然后 `make users-check`
##   （脚本会创建/删除测试账号并自行收尾，不会残留）
users-check:
	./scripts/check-users.sh

## users-deploy-check: 用户管理的部署形态（重启后仍在 / 损坏文件拒绝启动 / 功能开关）
##   会自己起停多个实例，不需要预先启动服务
users-deploy-check:
	./scripts/check-users-deploy.sh

## cdp: 启动带调试端口的 headless Chrome，供 layout / layout-sweep 使用
##   （前台运行，另开终端跑 make layout）
cdp:
	$(CHROME) --headless=new --no-sandbox --disable-gpu --disable-dev-shm-usage \
	  --no-proxy-server --user-data-dir=/tmp/chrome-cdp \
	  --remote-debugging-port=$(CDP_PORT) --window-size=1440,900 about:blank

## layout: 量一个视口宽度下的表格布局（需先 make cdp + 起一个无鉴权实例）
##   用法：make layout W=1280；目标开了鉴权时加 MEASURE_AUTH=user:pass
layout:
	$(NODE_BIN) scripts/measure-layout.mjs $(MEASURE_URL) $(or $(W),1440) 900

## layout-sweep: 扫一遍所有响应式断点，每个宽度都要「列数一致 + 不溢出」
layout-sweep:
	@for w in 1600 1440 1281 1280 1024 900 820 760 560 500 414 375; do \
	  printf '  %-6spx  ' $$w; \
	  if $(NODE_BIN) scripts/measure-layout.mjs $(MEASURE_URL) $$w 900 >/tmp/gofs-layout.log 2>&1; then \
	    grep -oE '列数一致（[0-9]+ 列）' /tmp/gofs-layout.log | tr -d '\n'; \
	    printf '  \033[32mOK\033[0m\n'; \
	  else \
	    printf '\033[31mFAIL\033[0m\n'; sed -n '2p;14,20p' /tmp/gofs-layout.log; \
	  fi; \
	done

## touch-sweep: 同样是断点扫描，但模拟触摸设备
##   额外校验「无横向溢出 / 可点目标 ≥ 40px / 可点元素不能是半透明」
touch-sweep:
	@for w in 320 375 414 480 560 561 700 768 900 1024 1280 1440; do \
	  printf '  %-6spx  ' $$w; \
	  if TOUCH=1 $(NODE_BIN) scripts/measure-layout.mjs $(MEASURE_URL) $$w 812 >/tmp/gofs-touch.log 2>&1; then \
	    printf '\033[32mOK\033[0m  '; \
	    grep -oE '最小边 < 40px 的有 [0-9]+ 个' /tmp/gofs-touch.log; \
	  else \
	    printf '\033[31mFAIL\033[0m\n'; sed -n '2p;/横向溢出/p;/触摸目标过小/,$p' /tmp/gofs-touch.log | head -8; \
	  fi; \
	done

## touch-modals: 单独校验弹窗类界面（编辑器 / 行操作菜单 / 服务设置）的触摸可用性
touch-modals:
	@for o in sheet editor settings; do \
	  printf '  %-6s  ' $$o; \
	  if TOUCH=1 MEASURE_OPEN=$$o $(NODE_BIN) scripts/measure-layout.mjs $(MEASURE_URL) 375 812 >/tmp/gofs-modal.log 2>&1; then \
	    printf '\033[32mOK\033[0m  '; \
	    grep -oE '最小边 < 40px 的有 [0-9]+ 个' /tmp/gofs-modal.log; \
	  else \
	    printf '\033[31mFAIL\033[0m\n'; sed -n '/触摸目标过小/,$p' /tmp/gofs-modal.log | head -8; \
	  fi; \
	done

## vet: 静态检查
vet:
	go vet ./...

## fmt: 格式化
fmt:
	gofmt -w .
	gofmt -l .

## run: 以全权限模式服务 ./data 目录
run: build
	./$(BINARY) -A -b 127.0.0.1 -p 5000 ./data

## testdata: 生成端到端手工测试用的目录与压缩包（/tmp/gofs-test）
testdata:
	python3 scripts/gen_testdata.py

## cross: 交叉编译常用平台
cross:
	GOOS=linux   GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o dist/gofs-linux-amd64 .
	GOOS=linux   GOARCH=arm64 go build -trimpath -ldflags "-s -w" -o dist/gofs-linux-arm64 .
	GOOS=darwin  GOARCH=arm64 go build -trimpath -ldflags "-s -w" -o dist/gofs-darwin-arm64 .
	GOOS=darwin  GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o dist/gofs-darwin-amd64 .
	GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o dist/gofs-windows-amd64.exe .
	@ls -lh dist/

## clean: 清理构建产物
clean:
	rm -rf $(BINARY) dist
