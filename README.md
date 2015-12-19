# ifviva
Golang Web Framework For http://www.ifviva.com

该框架正在开发中，尚未正式投入使用。

# 快速开始
app.go:
```go
package app

import "ifviva"

type ArticleCtrl struct {
    ifviva.Controller
}

func (ctrl *ArticleCtrl) GetArticle(id String) {
    ctrl.Status(200)
    ctrl.Json(map[string]string{
        "id": id
    })
}

func main() {
    app := ifviva.Application{}
    app.Get("/articles/:id", func(ctx ifviva.Context) {
        articleCtrl := ArticleCtrl{}
        articleCtrl.Init(ctx)
        articleCtrl.GetArticle(ctx.Params["id"])
    })
    app.Run(":3000")
}
```

# 主要设计
简约的路由功能，支持Restful API，但一个路由只对应一个处理函数，不在路由层作复杂的中间件处理，而把这复杂的工作交给控制器。

这样做除了要多写几行代码之外，有以下的好处：
## 解耦合
在不需要的时候，可以不使用框架提供的控制器，只使用路由匹配功能

## 灵活高效
将复杂的处理过程交给控制器，可以有效避免中间件之间交换数据带来的性能和可靠性的损失。如果使用中间件，代码大概会是这样子：
```go
type session struct {
    id string
}

func main() {
    app := ifviva.Application{}
    app.Get("/", func(ctx ifviva.Context) {
        // ctx.Data是一个map[string]interface{}类型
        ctx.Data["session"] = session{"123"}
    }, func(ctx ifviva.Context) {
        session := ctx.Data["session"].(session)
        ctrl.Res.WriteHeader(200)
        ctrl.Res.Write([]byte(session.id))
    })
}
```

这样把数据以空接口的形式保存，会带来性能和可靠性的损失。而如果用控制器，既能保持灵活性，又能保持高效可靠，避免类型转换。
```go
package app

import "ifviva"

type UserCtrl struct {
    ifviva.Controller
    session session // 可以根据需要定义中间件的数据
}

func (ctrl *UserCtrl) Init(ctx ifviva.Context) {
    // 在这里覆盖初始化方法
    ctrl.Controller.Init(ctx)
    ctrl.session = session{ctx.Params["id"]}
}

func (ctrl *UserCtrl) Login() {
    ctrl.Status(200)
    ctrl.Text(ctrl.session.id) // 可以获取session的值，而不用作类型转换之类的工作
}

func main() {
    app := ifviva.Application{}
    app.Post("/users/:id/login", func(ctx ifviva.Context) {
        userCtrl := UserCtrl{}
        userCtrl.Init(ctx)
        userCtrl.Login()
    })
    app.Run(":3000")
}
```

# API文档
(建设中)

# License
[MIT](./LICENSE)
