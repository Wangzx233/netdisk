# begonia源码分析

## 作用

``begonia``是一个``RPC``框架，可以实现远程函数调用。

之前已经有相关课程介绍了rpc，这里就不再赘述。

## 原理

函数远程调用的过程可以简化为：客户端把需要调用函数的**入参序列化**并传递给服务端，服务端**接受后反序列化**后**本地执行**获得结果，把**结果序列化**传递再给客户端，最后**客户端反序列化，获得结果**。

通俗的说，调用者就是把函数入参传给服务器，服务器执行完再把结果返回给调用者，只不过多了序列化等步骤而已。

## 为什么不用http请求

在刚接触``RPC``概念的时候，一直有一个疑问：``http``请求似乎也可以实现上述功能，那为什么还会出现``RPC``？

原因是``http``在传输的过程中**有用信息占比少**，请求头包含了大量臃肿的信息，使得它的**效率并不高**，其次，在部分场景下，确实可以用``http``请求代替``RPC``使用。

## ``begonia``节点

在集群模式下，``begonia``除了客户端和服务端外，还有一个服务中心节点。

服务端向服务中心注册服务，服务中心会保存服务相关信息，

当客户端发起调用后，服务中心就会根据保存的信息，把调用请求转发给服务端，并进行记录，

当服务端处理入参后，返回结果给服务中心，服务中心再根据记录把结果转发给服务端。

Begonia 文档有这么一句话，一切皆 RPC ，一切皆服务。甚至用**于注册的服务中心也是一个服务**，它注册的是**为其它服务端节点提供注册能力**的服务。

## 请求流程

<img src="C:\Users\wuhaoda\AppData\Roaming\Typora\typora-user-images\image-20210407220137402.png" alt="image-20210407220137402" style="zoom:80%;" />

当客户端获得服务的信息后，每次不需要重复黄色部分，只需要蓝色的调用部分就可以了。

## 架构

``begonia`` 从代码层面抽象成了三层，这三层分别提供不同的能力。

### App (Application) 层

这一层主要提供的是序列化、反序列化能力和正确的函数调用能力，用户调用的API就来自于这里。

### Logic 层

`Logic` 层收到 `App` 层的包已经经过了序列化，所以在这里仅需要将其封装为后发送给 `Dispatch` 层，然后注册该请求的回调，等待响应到来。

并且 `Logic` 层会在 `Dispatch` 层中注册事件，当有新的数据包到来时，`Logic` 会通过数据包的序列号来查找相对应的回调。

### Dispatch 层

`Dispatch` 层主要负责的是通讯，对端口的监听、或连接至对应的地址。

## 代码生成服务

``begonia``提供了反射和代码生成两种模式，在一般情况下， 使用反射即可实现 RPC 的功能，但反射的开销较大，用户可以使用代码生成来替换反射的实现，去除反射的开销。另外今后网校的项目预计会逐步改为代码生成模式运行，鉴于时间的原因，这里就只介绍一下begonia的代码生成模式。

在代码生成模式下，既然调用者要远程调用服务端的函数，那么他就要拥有函数的相关信息。

相信大家已经使用过grpc了，在grpc中，需要开发者编写包含服务出入参信息的``.proto``文件并编译成``.db.go``，以此来使客户端进行调用。而在begonia中更为方便，直接输入命令，服务端就会根据本地被注册的服务，直接生成它们的参数信息等数据，只需要把生成的``call``文件夹复制到客户端，客户端便可以获取这些信息。

### 服务端生成的文件

在服务端执行``$ begonia -s -r ./``后会生成``Service.begonia.go``文件

#### 出入参信息

```go
var (
	_EchoServiceFuncList []cRegister.FunInfo

	_EchoServiceSayHelloInSchema = `
{
			"namespace":"begonia.func.SayHello",  
			"type":"record",
			"name":"In",
			"fields":[ //字段
				{"name":"F1","type":"string","alias":"name"}

			]
		}`
	_EchoServiceSayHelloOutSchema = `
{
			"namespace":"begonia.func.SayHello",
			"type":"record",
			"name":"Out",
			"fields":[
				{"name":"F1","type":"string"}

			]
		}`
	_EchoServiceSayHelloInCoder  coding.Coder
	_EchoServiceSayHelloOutCoder coding.Coder
)
```
可以发现这里的``var``中``_EchoServiceSayHelloInSchema``与``_EchoServiceSayHelloOutSchema``字段分别保存了``SayHello``方法的入参与出参信息，解析器只需要获取这些信息就可以进行序列化等操作

#### ``Do``函数

```go
//代码生成模式下，一个rpc调用发送到服务端后，服务端会使用Do函数调用本地函数
func (d *EchoService) Do(ctx context.Context, fun string, param []byte) (result []byte, err error) {
   switch fun { //通过switch语句确定服务中的方法

   case "SayHello": 
      var in _EchoServiceSayHelloIn
      err = _EchoServiceSayHelloInCoder.DecodeIn(param, &in) //根据上述var中的信息反序列化param并绑定到in
      if err != nil {
         panic(err)
      }

      res1 := d.SayHello( //把in传入本地函数并执行，这里的SayHello就是我们编写的函数

         in.F1,
      )
      if err != nil {
         return nil, err
      }
      var out _EchoServiceSayHelloOut
      out.F1 = res1

      res, err := _EchoServiceSayHelloOutCoder.Encode(out) //把结果序列化
      if err != nil {
         panic(err)
      }
      return res, nil //返回序列化的结果

   default:
      err = errors.New("rpc call error: fun not exist")
   }
   return
}
```

**代码生成模式下，客户端发起的所有调用最后都会在服务端执行相应的Do函数**

### 中心节点``center``生成的代码

前面说过，``Center``节点其实是一个特殊的``server``节点，同时这个特殊的``server``节点还使用了代码生成服务。

既然是``server``节点，那就是说如同普通的``server``节点一样，``center``也注册了供其他节点调用的服务。

事实上，``center``**在本地**自己给自己注册了一个服务``register``，它的作用是为其它``server``提供注册能力。

``center``使用了代码生成服务，我们可以在``core/register/gencode.go``下看到这个服务生成的代码

```go
func (r *CoreRegister) Do(ctx context.Context, fun string, param []byte) (result []byte, err error) {
	switch fun {
	
    //服务端注册服务的时候使用的register其实就是来到了这里执行的
	case "Register":
		var in _CoreRegisterRegisterIn
		err = _CoreRegisterRegisterInCoder.DecodeIn(param, &in)
		if err != nil {
			panic(err)
		}

		err := r.Register(
			ctx,
			in.F1,
		)
		if err != nil {
			return nil, err
		}
		var out _CoreRegisterRegisterOut

		res, err := _CoreRegisterRegisterOutCoder.Encode(out)
		if err != nil {
			panic(err)
		}
		return res, nil

	case "ServiceInfo":
		var in _CoreRegisterServiceInfoIn
		err = _CoreRegisterServiceInfoInCoder.DecodeIn(param, &in)
		if err != nil {
			panic(err)
		}

		res1, err := r.ServiceInfo(

			in.F1,
		)
		if err != nil {
			return nil, err
		}
		var out _CoreRegisterServiceInfoOut
		out.F1 = res1

		res, err := _CoreRegisterServiceInfoOutCoder.Encode(out)
		if err != nil {
			panic(err)
		}
		return res, nil

	default:
		err = errors.New("rpc call error: fun not exist")
	}
	return
}
```

## 源码分析

下面来到了紧张刺激的扒源码环节，为了让叙述过程更有逻辑，这里会**按照流程顺序**逐步分析源码

同时底层的数据传输和tcp监听等过于细节的内容就不详细说明了，有兴趣的可以自行学习。

### center初始化

在终端里使用``bgacenter start``开启一个center节点，

center相关源码的入口在``cmd/bgacenter/main.go``

这里让我们打开``IDE``康康``center``部分的源码

### service初始化并注册服务

#### 生成一个Server

```go
s := begonia.NewServer(option.Addr(":12306"))
```

其实和center的初始化非常相似，不过也有一些区别

- server使用集群而p2p模式
- 因为注册服务实际上在center，所以server用``remoteRegister``封装``coreRegister``
- server用``link``连接到center而非``set``

这里就不去看源码了

#### 注册服务

```go
s.Register("Test", echoService)
```

源码见``app/server/server_ast.go``，

这里讲下源码

### center处理注册请求

service会向连接的center发送注册请求，那么center是如何处理请求的呢？

**之前忽略过center的``listen``部分，现在从这一部分的源码入手看一下center处理注册请求的逻辑**

### client初始化并发起调用

#### 客户端代码

客户端把call文件夹复制过来后，只需要使用``call.FuncName()``就可以实现远程调用

```go
package main

import (
	"begonia_test/call"
	"fmt"
)

func main() {
	res, err := call.SayHello("wuhaoda") //这里可以实现像调用本地函数一样，无感知的去调用远程函数
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
```

call文件夹下有``cli.begonia.go``和``NameService.begonia.go``两个文件

前者对客户端以代码生成模式进行初始化：

```go
var BegoniaCli begonia.Client = begonia.NewClientWithAst(option.Addr("wuhaoda.life:12306"))
```

后者在保存了服务出入参等信息的同时，还有一些其它内容

```go
func init() {
	app.ServiceAppMode = app.Ast

	bService, err := BegoniaCli.Service("Echo") //这里获取一个服务
	if err != nil {
		panic(err)
	}

	_EchoServiceServiceSayHello, err = bService.FuncSync("SayHello")
    //根据方法名拿到可以调用的远程方法

	_EchoServiceServiceSayHelloInCoder, err = coding.NewAvro(_EchoServiceServiceSayHelloInSchema)
	if err != nil {
		panic(err)
	}
	_EchoServiceServiceSayHelloOutCoder, err = coding.NewAvro(_EchoServiceServiceSayHelloOutSchema)
	if err != nil {
		panic(err)
	}

}
func SayHello(name string) (F1 string, err error) {
	var in _EchoServiceServiceSayHelloIn
	in.F1 = name //保存入参信息

	b, err := _EchoServiceServiceSayHelloInCoder.Encode(in) //序列化
	if err != nil {
		panic(err)
	}

	begoniaResTmp, err := _EchoServiceServiceSayHello(b) 
    //向拿到的方法中传入序列化的入参，获得被序列化的出参
	if err != nil {
		return
	}

	var out _EchoServiceServiceSayHelloOut
	err = _EchoServiceServiceSayHelloOutCoder.DecodeIn(begoniaResTmp.([]byte), &out)
	if err != nil {
		panic(err)
	}

	F1 = out.F1 //获得结果之后解码，反序列化

	return
}
```

#### 初始化client

现在详细看一下上面那些代码

```go
var BegoniaCli begonia.Client = begonia.NewClientWithAst(option.Addr("wuhaoda.life:12306"))
```

入口位于``/client.go``，继续看源码

#### 远程调用

```go
bService, err := BegoniaCli.Service("Echo")
```

首先在``init()``函数中，会获取所调用的服务

它位于``app/client/client.go``中，看一波源码

```go
_EchoServiceServiceSayHello, err = bService.FuncSync("SayHello")
```

然后用获取到的服务，根据方法名，去拿远程调用的方法

``FuncSync``位于``/app/client/client_service_ast.go``目录下，肝源码

### center处理client的``RPC``请求

具体操作是转发该请求给service，并监听service返回结果，收到结果后把其再转发给client

我们从``dispatch/dispatch_default_set``中的``Listen``开始看，

其实``Listen``后面的源码又是之前讲过了的，不过此时代理器``proxy``就起到了作用

### service处理center转发的请求

service在注册服务之后一直处于监听状态，监听到请求之后就会调用work函数

直接看``dispatch/dispatch_default_link.go``下的``Work()``函数

```go
// work 获得一个新的连接之后持续监听连接，然后把消息发送到msgCh里
func (d *linkDispatch) work(c conn.Conn) {

	id := ids.New()

	d.linkID = id
	log.Printf("link addr [%s] success, connID [%s]\n", c.Addr(), id)

	d.DoStartHook(id) // 变量初始化完成，这里去hook一些东西

	for {

		opcode, payload, err := c.Recv()
		if err != nil {
			c.Close()
			d.DoCloseHook(id, err)
			break
		}

		d.rt.Do(id, opcode, payload)
	}
}
```

此时service在``for``循环中接收到之后获取到``opcode``和``payload``，执行``d.rt.Do``，这里的Do就是之前讲过的那个Do，流程再次重复，也不讲了。

## 总结

~~当我们在阅读源码时，我们在阅读什么~~

- 一些概念``hook``函数，``handle``函数，``context``控制并发，一些设计模式，心跳机制, 代码结构...
- 更熟悉``begonia``了，在写业务的过程中，阅读过源码更容易排查问题
- 了解到一些欠缺，官方包了解不深，``io``相关不熟悉...
- 明白自己的一些不足之处语法树，连接池...

ceSayHelloIn
      err = _EchoServiceSayHelloInCoder.DecodeIn(param, &in) //根据上述var中的信息反序列化param并绑定到in
      if err != nil {
         panic(err)
      }

      res1 := d.SayHello( //把in传入本地函数并执行，这里的SayHello就是我们编写的函数

         in.F1,
      )
      if err != nil {
         return nil, err
      }
      var out _EchoServiceSayHelloOut
      out.F1 = res1

      res, err := _EchoServiceSayHelloOutCoder.Encode(out) //把结果序列化
      if err != nil {
         panic(err)
      }
      return res, nil //返回序列化的结果

   default:
      err = errors.New("rpc call error: fun not exist")
   }
   return
}
```

**代码生成模式下，客户端发起的所有调用最后都会在服务端执行相应的Do函数**

### 中心节点``center``生成的代码

前面说过，``Center``节点其实是一个特殊的``server``节点，同时这个特殊的``server``节点还使用了代码生成服务。

既然是``server``节点，那就是说如同普通的``server``节点一样，``center``也注册了供其他节点调用的服务。

事实上，``center``**在本地**自己给自己注册了一个服务``register``，它的作用是为其它``server``提供注册能力。

``center``使用了代码生成服务，我们可以在``core/register/gencode.go``下看到这个服务生成的代码

```go
func (r *CoreRegister) Do(ctx context.Context, fun string, param []byte) (result []byte, err error) {
	switch fun {
	
    //服务端注册服务的时候使用的register其实就是来到了这里执行的
	case "Register":
		var in _CoreRegisterRegisterIn
		err = _CoreRegisterRegisterInCoder.DecodeIn(param, &in)
		if err != nil {
			panic(err)
		}

		err := r.Register(
			ctx,
			in.F1,
		)
		if err != nil {
			return nil, err
		}
		var out _CoreRegisterRegisterOut

		res, err := _CoreRegisterRegisterOutCoder.Encode(out)
		if err != nil {
			panic(err)
		}
		return res, nil

	case "ServiceInfo":
		var in _CoreRegisterServiceInfoIn
		err = _CoreRegisterServiceInfoInCoder.DecodeIn(param, &in)
		if err != nil {
			panic(err)
		}

		res1, err := r.ServiceInfo(

			in.F1,
		)
		if err != nil {
			return nil, err
		}
		var out _CoreRegisterServiceInfoOut
		out.F1 = res1

		res, err := _CoreRegisterServiceInfoOutCoder.Encode(out)
		if err != nil {
			panic(err)
		}
		return res, nil

	default:
		err = errors.New("rpc call error: fun not exist")
	}
	return
}
```

## 源码分析

下面来到了紧张刺激的扒源码环节，为了让叙述过程更有逻辑，这里会**按照流程顺序**逐步分析源码

同时底层的数据传输和tcp监听等过于细节的内容就不详细说明了，有兴趣的可以自行学习。

### center初始化

在终端里使用``bgacenter start``开启一个center节点，

center相关源码的入口在``cmd/bgacenter/main.go``

这里让我们打开``IDE``康康``center``部分的源码

### service初始化并注册服务

#### 生成一个Server

```go
s := begonia.NewServer(option.Addr(":12306"))
```
