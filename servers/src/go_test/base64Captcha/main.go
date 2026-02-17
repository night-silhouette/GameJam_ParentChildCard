package main

import (
	"fmt"
	"os"

	"github.com/mojocn/base64Captcha"
)

func main() {
	// 1. 定义驱动（数字验证码）
	driver := base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80)

	// 2. 使用默认的内存存储
	store := base64Captcha.DefaultMemStore

	// 3. 创建验证码实例并生成
	c := base64Captcha.NewCaptcha(driver, store)
	id, b64s, answer, err := c.Generate()

	if err != nil {
		fmt.Println("生成失败:", err)
		return
	}

	// 4. 打印结果
	fmt.Println("验证码ID:", id)
	fmt.Println("验证码答案:", answer)
	// 因为 Base64 字符串太长，这里只打印前 50 位
	fmt.Println("图片Base64(前50位):", b64s[:50])

	// 5. 【重点】为了方便你直接看图，我们将它存成一个 HTML 文件
	// 这样你在 Mac 上直接双击这个文件就能看到生成的验证码图片了
	htmlContent := fmt.Sprintf("<html><body><img src=\"%s\"></body></html>", b64s)
	err = os.WriteFile("index.html", []byte(htmlContent), 0644)
	if err == nil {
		fmt.Println("\n✅ 生成成功！请在当前目录下查找 'index.html'，用浏览器打开即可查看图片。")
	}
}
