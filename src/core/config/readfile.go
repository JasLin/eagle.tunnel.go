package config

import (
	"net"
	"path"
	"strings"
	"time"

	"github.com/eaglexiang/eagle.tunnel.go/src/core/protocols/et/comm"
	"github.com/eaglexiang/go-bytebuffer"
	"github.com/eaglexiang/go-logger"
	mynet "github.com/eaglexiang/go-net"
	"github.com/eaglexiang/go-settings"
	myuser "github.com/eaglexiang/go-user"
)

// ImportConfigFiles 导入配置
func ImportConfigFiles() {
	readConfigFile()

	initListen()
	initRelay()
	initProxyStatus()
	initLocalUser()
	initTimeout()
	initBufferSize()

	readConfigDir()
}

// readConfigFile 读取根据给定的配置文件
func readConfigFile() {
	if !settings.Exsit("config") {
		return
	}

	filePath := settings.Get("config")
	allConfLines := readLinesFromFile(filePath)
	err := settings.ImportLines(allConfLines)
	if err != nil {
		panic(err)
	}
}

func initTimeout() {
	timeout := settings.GetInt64("timeout")
	Timeout = time.Second * time.Duration(timeout)
	comm.Timeout = Timeout
}

func initBufferSize() {
	size := settings.GetInt64("buffer.size")
	bytebuffer.SetDefaultSize(int(size))
}

func initLocalUser() {
	// 读取本地用户
	if !settings.Exsit("user") {
		SetUser("null:null")
	} else {
		SetUser(settings.Get("user"))
	}
}

// initUserList 初始化用户列表
func initUserList() {
	if settings.Get("user-check") != "on" {
		return
	}

	usersPath := path.Join(settings.Get("config-dir"), "/users.list")
	importUsers(usersPath)
}

//SetUser 设置本地用户
func SetUser(user string) {
	var err error
	LocalUser, err = myuser.ParseValidUser(user)
	if err != nil {
		panic(err)
	}
}

func importUsers(usersPath string) {
	Users = make(map[string]*myuser.ValidUser)
	userLines := readLinesFromFile(usersPath)

	for _, line := range userLines {
		user, err := myuser.ParseValidUser(line)
		if err != nil {
			panic(err)
		}
		Users[user.ID] = user
	}
	logger.Info(len(Users), " users imported")
}

// finishPort 补全端口号
func finishIPEPort(ipe string) string {
	switch mynet.TypeOfAddr(ipe) {
	case mynet.IPv4Addr:
		if ip := net.ParseIP(ipe); ip != nil {
			// 不包含端口号
			ipe += ":8080"
		}
	case mynet.IPv6Addr:
		if strings.HasSuffix(ipe, "]") {
			// 不包含端口号
			ipe += ":8080"
		}
	}
	return ipe
}

// finishIPEs ipes的示例：192.168.0.1:8080,192.168.0.1:8081
func finishIPEs(ipes string) (newIPEs string) {
	_ipes := strings.Split(ipes, ",")
	for _, ipe := range _ipes {
		newIPEs += "," + finishIPEPort(ipe)
	}
	newIPEs = strings.TrimPrefix(newIPEs, ",") // 去掉头部多余的,符号
	return
}

func initProxyStatus() {
	var err error
	s := settings.Get("proxy-status")
	ProxyStatus, err = comm.ParseProxyStatus(s)
	if err != nil {
		panic(err)
	}
}

func initListen() {
	ipes := settings.Get("listen")
	ipes = finishIPEs(ipes)
	settings.Set("listen", ipes)
}

func initRelay() {
	ipes := settings.Get("relay")
	ipes = finishIPEs(ipes)
	settings.Set("relay", ipes)
}