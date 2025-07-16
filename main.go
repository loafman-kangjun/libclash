package main

import (
    "fmt"
    "os"

    "github.com/loafman-kangjun/libclash/libcore"
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


    // 程序退出前关闭
    libcore.Stop()
}