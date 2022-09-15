# MiraiGo-module-read60s

ID: `com.aimerneige.read60s`

Module for [MiraiGo-Template](https://github.com/Logiase/MiraiGo-Template)

## 功能

- 接受到关键词后发送「每天 60 秒读懂世界」图片

## 鸣谢

本项目调用且依赖了 <https://api.2xb.cn/zaob> 的服务。

## 使用方法

在适当位置引用本包

```go
package example

imports (
    // ...

    _ "github.com/yukichan-bot-module/MiraiGo-module-read60s"

    // ...
)

// ...
```

如果需要自定义触发关键词，修改你的 `application.yaml` 配置文件：

```yaml
aimerneige:
  read60s:
    keywords:
      - "今日新闻"
      - "60s"
      - "早报"
```

## LICENSE

<a href="https://www.gnu.org/licenses/agpl-3.0.en.html">
<img src="https://www.gnu.org/graphics/agplv3-155x51.png">
</a>

本项目使用 `AGPLv3` 协议开源，您可以在 [GitHub](https://github.com/yukichan-bot-module/MiraiGo-module-read60s) 获取本项目源代码。为了整个社区的良性发展，我们强烈建议您做到以下几点：

- **间接接触（包括但不限于使用 `Http API` 或 跨进程技术）到本项目的软件使用 `AGPLv3` 开源**
- **不鼓励，不支持一切商业使用**
