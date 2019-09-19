// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 148.

// Fetch saves the contents of a URL into a local file.
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)
//下面的代码是fetch（ 1.5节） 的改进版， 我们将http响应信息写入本地文件而不是从标准输出
//流输出。 我们通过path.Base提出url路径的最后一段作为文件名。

//对resp.Body.Close延迟调用我们已经见过了， 在此不做解释。 上例中， 通过os.Create打开文
//件进行写入， 在关闭文件时， 我们没有对f.close采用defer机制， 因为这会产生一些微妙的错
//误。 许多文件系统， 尤其是NFS， 写入文件时发生的错误会被延迟到文件关闭时反馈。 如果
//没有检查文件关闭时的反馈信息， 可能会导致数据丢失， 而我们还误以为写入操作成功。 如
//果io.Copy和f.close都失败了， 我们倾向于将io.Copy的错误信息反馈给调用者， 因为它先于
//f.close发生， 更有可能接近问题的本质。
//练习5.18： 不修改fetch的行为， 重写fetch函数， 要求使用defer机制关闭文
//!+
// Fetch downloads the URL and returns the
// name and length of the local file.
func fetch(url string) (filename string, n int64, err error) {
	resp, err := http.Get(url) //发起连接，发送一个http请求
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()  //再函数退出前关闭连接，释放资源

	local := path.Base(resp.Request.URL.Path) //path这是干啥用的？
	if local == "/" {
		local = "index.html"
	}
	f, err := os.Create(local) //创建文件
	if err != nil {
		return "", 0, err
	}
	//defer f.Close() //不推荐这种做法，因为“写入文件时发生的错误会被延迟到文件关闭时反馈”，所有要操作完文件旧的立刻关闭文件
	n, err = io.Copy(f, resp.Body)  //将http响应内容拷贝文件里（这种情况下，都不用打开文件操作，os.open？）
	// Close file, but prefer error from Copy, if any.？？
	//先关闭文件，在判断文件是否写入成功（这么做的目的是：写入文件时发生的错误会被延迟到文件关闭时反馈，所以写完就关闭文件，以及时获取写文件后的反馈信息）
	if closeErr := f.Close(); err == nil { //注意这种使用，先赋值运算，再if判断真假；
		err = closeErr
	}
	return local, n, err
}

//!-

func main() {
	for _, url := range os.Args[1:] {
		local, n, err := fetch(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch %s: %v\n", url, err)
			continue
		}
		fmt.Fprintf(os.Stderr, "%s => %s (%d bytes).\n", url, local, n)
	}
}






