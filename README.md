# go-ecs

一个轻量级、高性能的 Go 语言 Entity-Component-System（ECS）框架，适用于游戏开发和高性能数据驱动应用。

## 特性

- 🚀 **高性能** - 基于稀疏集（Sparse Set）实现高效的组件存储与查询
- 🎯 **简单易用** - 提供简洁直观的 API 设计
- 🔧 **泛型支持** - 充分利用 Go 1.18+ 泛型特性，提供类型安全的操作
- 🧩 **灵活组合** - 支持组件、实体、系统的灵活组合
- 🔄 **对象池** - 内置组件对象池，减少 GC 压力
- 📦 **资源管理** - 支持全局资源（Resource）管理
- 📡 **事件系统** - 内置事件读写机制

## 安装

```shell
go get github.com/INT-Game/go-ecs
```

## 核心概念

### World（世界）

World 是 ECS 的核心容器，管理所有的实体、组件和系统。

```go
w := ecs.NewWorld()
```

### Component（组件）

组件是纯数据容器，不包含逻辑。通过嵌入 `ecs.Component` 来定义自定义组件：

```go
type PositionComponent struct {
    ecs.Component
    X, Y float64
}

type VelocityComponent struct {
    ecs.Component
    VX, VY float64
}
```

### Entity（实体）

实体是组件的容器，本身只是一个 ID 标识：

```go
// 创建组件
posComp := ecs.SpawnComponent[*PositionComponent](w)
posComp.X, posComp.Y = 100, 200

velComp := ecs.SpawnComponent[*VelocityComponent](w)
velComp.VX, velComp.VY = 1, 1

// 创建实体并附加组件
entity := ecs.SpawnEmptyEntity(w, posComp, velComp)
```

### System（系统）

系统包含处理组件的逻辑，通过嵌入 `ecs.System` 来定义：

```go
type MovementSystem struct {
    ecs.System
}

func NewMovementSystem(w *ecs.World) *MovementSystem {
    return &MovementSystem{
        System: *ecs.NewSystem(w),
    }
}

func (s *MovementSystem) Update() {
    // 查询所有同时拥有 Position 和 Velocity 组件的实体
    entities := s.Query.Query(&PositionComponent{}, &VelocityComponent{})
    for _, entity := range entities {
        pos, _ := s.Query.Get(entity, &PositionComponent{})
        vel, _ := s.Query.Get(entity, &VelocityComponent{})
        
        position := pos.(*PositionComponent)
        velocity := vel.(*VelocityComponent)
        
        position.X += velocity.VX
        position.Y += velocity.VY
    }
}
```

## 快速开始

```go
package main

import (
    "fmt"
    ecs "github.com/INT-Game/go-ecs/ecs"
)

// 定义组件
type NameComponent struct {
    ecs.Component
    Name string
}

// 定义系统
type NameSystem struct {
    ecs.System
}

func NewNameSystem(w *ecs.World) *NameSystem {
    return &NameSystem{
        System: *ecs.NewSystem(w),
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
    // 创建世界
    w := ecs.NewWorld()
    
    // 添加系统
    w.AddUpdateSystem(NewNameSystem(w))
    
    // 创建组件
    nameComponent := ecs.SpawnComponent[*NameComponent](w)
    nameComponent.Name = "Player1"
    
    // 创建实体
    ecs.SpawnEmptyEntity(w, nameComponent)
    
    // 运行更新循环
    w.Update()
}
```

## API 参考

### World

| 方法 | 说明 |
|------|------|
| `NewWorld()` | 创建新的 World 实例 |
| `AddStartUpSystem(system)` | 添加启动时执行一次的系统 |
| `AddUpdateSystem(system)` | 添加每帧更新的系统 |
| `Startup()` | 执行所有启动系统 |
| `Update()` | 执行所有更新系统 |
| `Shutdown()` | 清理世界中的所有资源 |
| `GetCommands()` | 获取命令对象 |
| `GetQuery()` | 获取查询对象 |

### Commands

| 方法 | 说明 |
|------|------|
| `DestroyEntity(entity)` | 标记实体待销毁 |
| `Execute()` | 执行所有待处理的命令 |
| `SetResource(component)` | 设置全局资源 |
| `RemoveResource(component)` | 移除全局资源 |

### Query

| 方法 | 说明 |
|------|------|
| `Query(components...)` | 查询包含指定组件的所有实体 |
| `Has(entity, component)` | 判断实体是否包含指定组件 |
| `Contains(entity, components...)` | 判断实体是否包含所有指定组件 |
| `Get(entity, component)` | 获取实体的指定组件 |

### Entity

| 方法 | 说明 |
|------|------|
| `SpawnEmptyEntity(world, components...)` | 创建实体并附加组件 |
| `SpawnEntity[T](world, components...)` | 创建自定义类型实体 |
| `SpawnComponent[T](world)` | 从对象池创建组件 |
| `GetComponent[T](entity)` | 泛型方式获取实体组件 |
| `AddComponents(components...)` | 向实体添加组件 |
| `RemoveComponents(components...)` | 从实体移除组件 |

### Resources

| 方法 | 说明 |
|------|------|
| `Has(resource)` | 判断是否存在指定资源 |
| `Get(resource)` | 获取指定资源 |
| `GetResource[T](resources)` | 泛型方式获取资源 |

### Events

| 方法 | 说明 |
|------|------|
| `NewEvents[T]()` | 创建事件实例 |
| `EventReader.Has()` | 判断是否有事件 |
| `EventReader.Get()` | 获取事件数据 |
| `EventWriter.Send(data)` | 发送事件 |

## 目录结构

```
go-ecs/
├── ecs/                 # 核心 ECS 实现
│   ├── world.go        # 世界管理
│   ├── entity.go       # 实体定义
│   ├── component.go    # 组件定义
│   ├── system.go       # 系统定义
│   ├── commands.go     # 命令模式实现
│   ├── query.go        # 查询系统
│   ├── spawner.go      # 实体/组件生成器
│   ├── resources.go    # 全局资源管理
│   ├── events.go       # 事件系统
│   └── pool.go         # 对象池
├── array/              # 动态数组实现
├── sparse_set/         # 稀疏集数据结构
├── main.go             # 示例入口
└── README.md
```

## 高级用法

### 自定义实体类型

```go
type PlayerEntity struct {
    *ecs.Entity
    PlayerID int
}

// 使用泛型创建
player := ecs.SpawnEntity[*PlayerEntity](w, posComp, velComp)
```

### 全局资源管理

```go
type GameConfig struct {
    ecs.Component
    Difficulty int
}

// 设置资源
config := &GameConfig{Difficulty: 1}
w.GetCommands().SetResource(config)

// 获取资源
resources := ecs.NewResources(w)
if cfg, ok := ecs.GetResource[*GameConfig](resources); ok {
    fmt.Println(cfg.Difficulty)
}
```

### 组件生命周期

```go
type MyComponent struct {
    ecs.Component
    Data []byte
}

func (c *MyComponent) Init() {
    // 组件初始化时调用
    c.Data = make([]byte, 1024)
}

func (c *MyComponent) Destroy() {
    // 组件销毁时调用
    c.Data = nil
}
```

## 性能提示

1. **使用组件查询** - 尽量使用 `Query.Query()` 批量查询，避免遍历所有实体
2. **对象池复用** - 使用 `SpawnComponent` 创建组件，框架会自动管理对象池
3. **延迟销毁** - 使用 `Commands.DestroyEntity()` 标记销毁，然后调用 `Commands.Execute()` 批量处理

## 许可证

MIT License，详见 [LICENSE](LICENSE) 文件。