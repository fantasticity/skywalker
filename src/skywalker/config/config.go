/*
 * Copyright (C) 2015 - 2016 Wiky L
 *
 * This program is free software: you can redistribute it and/or modify it
 * under the terms of the GNU General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 * See the GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.";
 */

package config

import (
	"flag"
	"github.com/hitoshii/golib/src/log"
	"os"
	"skywalker/agent"
	"skywalker/plugin"
	"skywalker/util"
)

type SkyWalkerExtraConfig SkyWalkerConfig

/* 服务配置 */
type SkyWalkerConfig struct {
	Name     string `json:name`
	BindAddr string `json:"bindAddr"`
	BindPort uint16 `json:"bindPort"`

	ClientProtocol string                 `json:"clientProtocol"`
	ClientConfig   map[string]interface{} `json:"clientConfig"`

	ServerProtocol string                 `json:"serverProtocol"`
	ServerConfig   map[string]interface{} `json:"serverConfig"`

	Log *log.LogConfig `json:"log"`

	DNSTimeout int64 `json:"dnsTimeout"`

	Daemon  bool                    `json:"daemon"`
	Plugins []*plugin.PluginConfig  `json:"plugin"`
	Extras  []*SkyWalkerExtraConfig `json:"extra"`
}

/* 将c的内容合并到cfg中 */
func (cfg *SkyWalkerConfig) Merge(c *SkyWalkerConfig) {
	if len(cfg.Name) == 0 {
		cfg.Name = c.Name
	}
	if len(cfg.BindAddr) == 0 {
		cfg.BindAddr = c.BindAddr
	}
	if len(cfg.ClientProtocol) == 0 {
		cfg.ClientProtocol = c.ClientProtocol
	}
	if cfg.ClientConfig == nil {
		cfg.ClientConfig = c.ClientConfig
	}
	if len(cfg.ServerProtocol) == 0 {
		cfg.ServerProtocol = c.ServerProtocol
	}
	if cfg.ServerConfig == nil {
		cfg.ServerConfig = c.ServerConfig
	}
	if cfg.Log == nil {
		cfg.Log = &log.LogConfig{
			ShowNamespace: c.Log.ShowNamespace,
			Loggers:       c.Log.Loggers,
		}
	} else if cfg.Log.Loggers == nil {
		cfg.Log.Loggers = defaultLoggers
	}
	cfg.Log.Namespace = cfg.Name
}

/*
 * 初始化配置
 * 设置日志、插件并检查CA和SA
 */
func (cfg *SkyWalkerConfig) Init() error {
	log.Init(cfg.Log)
	ca := cfg.ClientProtocol
	sa := cfg.ServerProtocol
	plugin.Init(cfg.Plugins, cfg.Name)
	if err := agent.CAInit(ca, cfg.Name, cfg.ClientConfig); err != nil {
		return err
	} else if err := agent.SAInit(sa, cfg.Name, cfg.ServerConfig); err != nil {
		return err
	}
	return nil
}

var (
	/* 默认配置 */
	defaultLoggers = []log.LoggerConfig{
		log.LoggerConfig{"DEBUG", "STDOUT"},
		log.LoggerConfig{"INFO", "STDOUT"},
		log.LoggerConfig{"WARNING", "STDERR"},
		log.LoggerConfig{"ERROR", "STDERR"},
	}
	defaultLogConfig = &log.LogConfig{
		Loggers: defaultLoggers,
	}
	gConfig = SkyWalkerConfig{
		Name:       "default",
		BindAddr:   "127.0.0.1",
		BindPort:   12345,
		DNSTimeout: 3600,
		/* 默认的日志输出 */
		Log:    defaultLogConfig,
		Daemon: false,
	}
)

const (
	DEFAULT_USER_CONFIG = "~/.config/skywalker.json"
	DEFAULT_GLOBAL_CONFIG = "/etc/skywalker.json"
)

/* 获取所有配置列表 */
func GetConfigs() []*SkyWalkerConfig {
	var configs []*SkyWalkerConfig

	gConfig.Log.Namespace = gConfig.Name
	configs = append(configs, &gConfig)
	for _, e := range gConfig.Extras {
		cfg := (*SkyWalkerConfig)(e)
		cfg.Merge(&gConfig)
		configs = append(configs, cfg)
	}
	return configs
}

/*
 * 查找配置文件，如果命令行参数-c指定了配置文件，则使用
 * 否则使用~/.config/skywalker.json
 * 否则使用/etc/skywalker.json
 */
func findConfigFile() string {
	file := flag.String("c", "", "the config file")
	flag.Parse()
	if len(*file) > 0 {
		return *file
	}
	checkRegularFile := func(filepath string) string {
		path := util.ResolveHomePath(filepath)
		info, err := os.Stat(path)
		if err == nil && info.Mode().IsRegular() {
			return path
		}
		return ""
	}
	if path := checkRegularFile(DEFAULT_USER_CONFIG); len(path) > 0 {
		return path
	} else if path := checkRegularFile(DEFAULT_GLOBAL_CONFIG); len(path) > 0 {
		return path
	}
	return ""
}

func init() {
	cfile := findConfigFile()
	if len(cfile) == 0 {
		util.FatalError("No Config Found!")
	} else if !util.LoadJsonFile(cfile, &gConfig) { /* 读取配置文件 */
		util.FatalError("Fail To Load Config From %s", cfile)
	}
	if gConfig.Log == nil {
		gConfig.Log = defaultLogConfig
	}

	/* 初始化DNS超时时间 */
	util.Init(gConfig.DNSTimeout)
}
