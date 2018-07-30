# snowflake
====

## 安装
```sh
go get github.com/anerg2046/snowflake
```

## 用法
默认有3个参数可以进行设置
- Epoch 设置一个小于当前时间的毫秒数，用于缩短时间位数
- MaxMachineBit 最大的机器标识bit位数，默认是8，取值范围1-255，如果为0会忽略机器标识这一段
- MaxSequenceBit 每毫秒内的流水号bit位数，默认是12，取值范围是0-4095

```go
package main

import (
	"fmt"

	"github.com/anerg2046/snowflake"
)

func main() {
	//以下3个参数为可选
	snowflake.Epoch = 1532620800000
	snowflake.MaxMachineBit = 10
	snowflake.MaxSequenceBit = 10
	//创建一个机器ID的节点，这里是1，如果为0则会忽略机器标识这一段
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	//获取一个全局唯一的ID
	id := node.NextID()
	//输出为int64类型
	fmt.Println(id.Int64())
	//输出为string类型
	fmt.Println(id.String())
	//输出为Base36类型，主要是缩短位数
	fmt.Println(id.Base36())

}
```