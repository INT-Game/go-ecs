# go-ecs

一个轻量级、高性能的 Go 语言 Entity-Component-System（ECS）框架，适用于游戏开发和高性能数据驱动应用。

## 特性
- 简单易用的 API
- 支持组件、实体、系统的灵活组合
- 高效的查询与批量操作
- 支持系统生命周期管理

## 安装

```shell
go get github.com/INT-Game/go-ecs
```

## 快速开始

```go
package main

import (
	"fmt"
	ecs "github.com/INT-Game/go-ecs/ecs"
)

type NameComponent struct {
	ecs.Component
	Name string
}

type NameSystem struct {
	ecs.System
}

func NewNameSystem(commands *ecs.Commands, query *ecs.Query) *NameSystem {
	return &NameSystem{
		System: *ecs.NewSystem(commands, query),
	}
}

func (s *NameSystem) Update() {
	entities := s.Query.Query(&NameComponent{})
	for _, entity := range entities {
		comp, ok := s.Query.Get(entity, &NameComponent{})
		if ok {
			fmt.Println(comp.(*NameComponent).Name)
		}
	}
}

func main() {
	w := ecs.NewWorld()
	commands := ecs.NewCommands(w)
	query := ecs.NewQuery(w)

	w.AddUpdateSystem(NewNameSystem(commands, query))

	nameComponent := ecs.CreateComponent[*NameComponent](w)
	nameComponent.Name = "TestNameComponent"

	commands.Spawn(nameComponent)
	w.Update()
}
```

## 目录结构
- `ecs/`：核心 ECS 实现
- `array/`、`sparse_set/`：底层数据结构
- `main.go`：示例入口

## 许可证

MIT License，详见 LICENSE 文件。
