# libclash

一个基于 Go 语言的 mihomo 内核启动示例。

## 目录结构

```
libclash/
├── main.go         // 程序入口，调用 libcore 启动/关闭内核
├── libcore.go      // 内核相关逻辑（需自行实现）
├── go.mod          // Go 模块文件
```

## 使用方法

1. **初始化 Go 模块（如未初始化）**
   ```sh
   go mod init libclash
   ```

2. **构建项目**
   ```sh
   go build
   ```

3. **运行程序**
   ```sh
   go run main.go
   ```

4. **自定义配置**
   - `cfgPath`：配置文件路径，留空使用默认路径
   - `cfgB64`：base64 编码的配置内容，留空不用
   - `geoMode`：是否启用地理模式

## 示例代码

```go
package main

import (
    "fmt"
    "os"
    "libclash/libcore"
)

func main() {
    cfgPath := "" // 留空则用默认路径
    cfgB64 := ""  // 不用 base64 配置
    geoMode := true

    if err := libcore.Run(cfgPath, cfgB64, geoMode); err != nil {
        fmt.Fprintf(os.Stderr, "启动失败: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("mihomo 内核已启动")
    libcore.Stop()
}
```

## 注意事项

- 请确保 `libcore.go` 的包名为 `libcore`，且与 `main.go` 在同一模块下。
- 如需自定义配置或扩展功能，请修改 `libcore.go