package main

import go_v8 "github.com/pefish/go-v8"

func main() {
	//  注册 worker 实例，并注册一个接受 JS 信息的回调
	worker := go_v8.New(func(msg []byte) []byte {
		println("In Go,", string(msg))
		return nil
	})

	utilCode := `...`
	jsCode := `
	// ArrayBuffer转为字符串，参数为ArrayBuffer对象
	function ab2str(buf) {
	  return String.fromCharCode.apply(null, new Uint16Array(buf));
	}
	
	// 字符串转为ArrayBuffer对象，参数为字符串
	function str2ab(str) {
	  var buf = new ArrayBuffer(str.length * 2); // 每个字符占用2个字节
	  var bufView = new Uint16Array(buf);
	  for (var i = 0, strLen = str.length; i < strLen; i++) {
		bufView[i] = str.charCodeAt(i);
	  }
	  return buf;
	}

    // 在 JS 中注册接收 Go 信息的函数，并打印接收到的信息
    V8Worker2.recv(msg => {
      V8Worker2.print("In js, " + ab2str(msg));
    });
    // 从 JS 向 Go 中发送『from js』信息
    V8Worker2.send(str2ab("from js"));
  `

	// 加载 JS 工具函数代码，这里按下不表
	worker.Load("utils.js", utilCode)
	// 加载业务代码
	worker.Load("foo.js", jsCode)
	// 从 Go 中向 JS 中发送『from Go』信息
	worker.SendBytes([]byte("from Go"))
}

// 运行结果：
// In Go, from js
// In js, from Go
