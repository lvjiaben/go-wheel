package main

import (
	"fmt"
	"log"

	"github.com/lvjiaben/go-wheel/app/backend/gen"
	"github.com/lvjiaben/go-wheel/pkg/container"
)

func main() {
	// 初始化容器（生成器不需要嵌入文件）
	c := container.NewContainer(nil)
	defer c.Shutdown()

	if err := c.Initialize(); err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	// 运行 CLI
	cli := gen.NewCLI(c.GetDB())
	if err := cli.Run(); err != nil {
		fmt.Printf("错误：%v\n", err)
	}
}
