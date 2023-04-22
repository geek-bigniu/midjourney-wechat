package errorhandler

import (
	"fmt"
	"os"
)

// 全局异常处理函数
func HandlePanic() {
	if err := recover(); err != nil {
		fmt.Println("程序运行失败:", err)
		fmt.Printf("按Enter键结束...")
		fmt.Scanln()
		os.Exit(0)
	}

}
