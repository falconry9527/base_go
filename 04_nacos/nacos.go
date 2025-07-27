package main

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"log"
)

func main() {
	// 1. 创建配置客户端
	sc := []constant.ServerConfig{
		{
			IpAddr: "127.0.0.1",
			Port:   8848,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         "f2c4a3b4-a3ec-49dc-81f7-f78b2632eb6a",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}

	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		log.Fatalf("Failed to create config client: %v", err)
	}

	// 2. 获取配置
	dataId := "223423424243"
	group := "DEFAULT_GROUP"

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		log.Fatalf("Failed to get config: %v", err)
	}
	fmt.Printf("Initial config content: %s\n", content)

	// 3. 监听配置变化 ： 监听多个dataId，需要循环监听
	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Printf("\nConfig changed:\ngroup:%s, dataId:%s, data:%s\n", group, dataId, data)
		},
	})
	if err != nil {
		log.Fatalf("Failed to listen config: %v", err)
	}

	fmt.Println("Listening for config changes... Press Ctrl+C to exit.")

	// 保持程序运行
	select {}
}
