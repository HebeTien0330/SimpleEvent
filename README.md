## SimpleEvent：一个简单的事件系统（多语言实现）

> 本项目实现了一个简单的事件管理系统，用于在项目中处理事件订阅与发布。
> 系统支持普通事件监听以及一次性事件监听，并提供了对事件的添加、移除和触发功能。

## 目录结构

* manager: 事件管理系统的主文件，包括事件管理器类的定义和相关接口的实现。
* event: 定义了 Event 类。

## 使用说明

### TypeScript版本

#### 初始化事件系统

```typescript
import { installEventSystem } from "./manager";

// 将事件系统安装到两个目标对象上
installEventSystem(Target1);
installEventSystem(Target2);
```

#### 订阅事件

```typescript
Target1.listen('myEvent', (args) => {
    console.log('Received:', args);
}, { once: true }); // 添加一个只监听一次的事件
```

#### 发布事件

```typescript
Target2.onEvent('myEvent', { message: 'Hello World!' });
```

#### 移除事件监听

```typescript
const id = Target1.listen('anotherEvent', (args) => {
    console.log('Another event received:', args);
});

// 可通过事件id取消执行事件的监听
Target1.cancel('anotherEvent', id, false);
```

#### 多个事件监听

```typescript
const eventList = [
    { eventName: 'event1', callback: () => console.log('Event 1 triggered') },
    { eventName: 'event2', callback: () => console.log('Event 2 triggered'), once: true }
];

const ids = window.listenMulti(eventList);
```

#### API文档

**EventManager 类**

* getInstance(): 获取 EventManager 的单例实例。
* on(eventName: string, callback: Function, filter?: Function, once?: boolean): number: 注册事件监听器。
* off(eventName: string, evtId?: number, once?: boolean): 移除事件监听器。
* call(eventName: string, args: Params, evtId?: number): void: 触发事件。

**Event 类**

* execute(args: Params): void: 执行事件回调。

**handler 对象**

* listen(eventName: string, callback: Function, filter?: Function, once?: boolean): number
* listenMulti(eventList: Array `<any>`): Array `<number>`
* cancel(eventName: string, evtId?: number, once?: boolean)
* onEvent(eventName: string, args: any, evtId?: number)

#### 注意事项

* 确保在调用 installEventSystem 之前导入相关模块。
* 在取消监听时，如果未指定 eventId，则会取消所有同名事件的监听。
* 事件过滤器可选，如果提供，则只有满足条件的事件才会被触发。

---

### Go版本

#### 初始化事件系统

```go
package main

import (
    "fmt"
    "github.com/your_repo/event"
)

type Target1 struct {
    event.EventHandler
}

type Target2 struct {
    event.EventHandler
}

func main() {
    Target1 := &Target1{}
    Target1.OnInit()
    Target2 := &Target2{}
    Target1.OnInit()
}
```

### 注册事件

```go
eventId := Target1.Register("testEvent", func(args ...interface{}) {
    fmt.Println("Event triggered with args:", args)
}, event.DefaultFilter, false)
```

### 取消注册事件

```go
Target1.Deregister("testEvent", eventId)
```

### 触发事件

```go
Target2.OnEvent("testEvent", "Hello, World!")
```

#### 触发单个事件

```go
Target2.OnSEvent("testEvent", eventId, "Hello, World!")
```

#### API文档

**GetEventManager**

* 获取全局唯一的事件管理器实例。

```go
// 函数签名
func GetEventManager() *EventManager
```

**RegisterEvent**
注册一个事件监听器。

**参数：**
* eventName (string): 事件名称。
* callback (func(...any)): 回调函数。
* filter (func(...any) bool): 过滤器函数，默认为 DefaultFilter。
* once (bool): 是否只执行一次。

**返回值：**
* int: 事件ID。

```go
func (slf *EventManager) RegisterEvent(eventName string, callback func(...any), filter func(...any) bool, once bool) int
```

**RemoveEvent**
移除一个事件监听器。

**参数：**
* eventName (string): 事件名称。
* evtId (interface{}): 事件ID。

```go
func (slf *EventManager) RemoveEvent(eventName string, evtId interface{})
```

**TriggerEvent**
触发一个事件。

**参数：**
* eventName (string): 事件名称。
* evtId (int): 事件ID。
* args (...any): 传递给回调函数的参数。

```go
func (slf *EventManager) TriggerEvent(eventName string, evtId int, args ...any)
```

#### 注意事项

* 事件管理器支持并发安全的操作。
* 如果没有指定事件ID，则会触发所有注册的事件监听器。
* 通过 RegisterEvent 注册的事件监听器，必须通过 RemoveEvent 来取消注册以避免内存泄漏。

---

### Python版本

#### 引入相关函数

```python
from EventSystem import Listen, ListenMulti, Cancel, OnEvent
```

#### 监听事件
```python
# 监听事件
EvtId = Listen("UserLogin", callback_function, filter_function)

# 一次性监听事件
EvtId = Listen("UserLogin", callback_function, filter_function, Once=True)
```

#### 取消事件监听
```python
# 取消事件监听
Cancel("UserLogin", EvtId)

# 取消一次性事件监听
Cancel("UserLogin", EvtId, True)

# 取消所有事件监听
Cancel("UserLogin")
```

#### 触发事件
```python
# 触发事件
OnEvent("UserLogin", arg1, arg2)

# 触发特定事件
OnEvent("UserLogin", EvtId, arg1, arg2)
```

#### 监听多个事件
```python
# 监听多个事件
EvtIds = ListenMulti([
    ("UserLogin", callback_function, filter_function, False),
    ("UserLogout", callback_function, filter_function, True)
])
```

#### 注意事项

* 确保在调用 Listen 或 ListenOnce 之前已经正确初始化了 EventManager。
* 在调用 Cancel 时，需要提供正确的事件名称和事件 ID。
* 如果没有提供事件 ID，则默认会执行所有监听该事件的回调函数。


## TODO
使用其他语言实现

