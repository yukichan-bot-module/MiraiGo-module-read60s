package read60s

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
)

var instance *read60s
var logger = utils.GetModuleLogger("com.aimerneige.read60s")
var keywordList []string

type read60s struct {
}

func init() {
	instance = &read60s{}
	bot.RegisterModule(instance)
}

func (r *read60s) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "com.aimerneige.read60s",
		Instance: instance,
	}
}

// Init 初始化过程
// 在此处可以进行 Module 的初始化配置
// 如配置读取
func (r *read60s) Init() {
	configKeywordList := config.GlobalConfig.GetStringSlice("aimerneige.read60s.keywords")
	if len(configKeywordList) == 0 {
		keywordList = []string{"今日新闻", "60s", "早报"}
	} else {
		for _, k := range configKeywordList {
			keywordList = append(keywordList, k)
		}
	}
}

// PostInit 第二次初始化
// 再次过程中可以进行跨 Module 的动作
// 如通用数据库等等
func (r *read60s) PostInit() {
}

// Serve 注册服务函数部分
func (r *read60s) Serve(b *bot.Bot) {
	b.GroupMessageEvent.Subscribe(func(c *client.QQClient, msg *message.GroupMessage) {
		solve60s(c, msg.ToString(), message.Source{
			PrimaryID:  msg.GroupCode,
			SourceType: message.SourceGroup,
		})
	})
	b.PrivateMessageEvent.Subscribe(func(c *client.QQClient, msg *message.PrivateMessage) {
		solve60s(c, msg.ToString(), message.Source{
			PrimaryID:  msg.Sender.Uin,
			SourceType: message.SourcePrivate,
		})
	})
}

// Start 此函数会新开携程进行调用
// ```go
//
//	go exampleModule.Start()
//
// ```
// 可以利用此部分进行后台操作
// 如 http 服务器等等
func (r *read60s) Start(b *bot.Bot) {
}

// Stop 结束部分
// 一般调用此函数时，程序接收到 os.Interrupt 信号
// 即将退出
// 在此处应该释放相应的资源或者对状态进行保存
func (r *read60s) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	// 别忘了解锁
	defer wg.Done()
}

func simpleText(msg string) *message.SendingMessage {
	return message.NewSendingMessage().Append(message.NewText(msg))
}

func isKeyword(s string) bool {
	for _, v := range keywordList {
		if s == v {
			return true
		}
	}
	return false
}

func solve60s(c *client.QQClient, msg string, target message.Source) {
	if !isKeyword(msg) {
		return
	}
	img, err := getImage()
	if err != nil {
		errMsg := simpleText("获取图片失败了呢，可能是服务器网络出问题了罢。")
		switch target.SourceType {
		case message.SourceGroup:
			c.SendGroupMessage(target.PrimaryID, errMsg)
		case message.SourcePrivate:
			c.SendPrivateMessage(target.PrimaryID, errMsg)
		}
		return
	}
	uploadedImage, err := c.UploadImage(target, img)
	imgMsg := message.NewSendingMessage().Append(uploadedImage)
	switch target.SourceType {
	case message.SourceGroup:
		c.SendGroupMessage(target.PrimaryID, imgMsg)
	case message.SourcePrivate:
		c.SendPrivateMessage(target.PrimaryID, imgMsg)
	}
}

func getImage() (io.ReadSeeker, error) {
	const apiURL = "https://api.2xb.cn/zaob"
	type apiResponse struct {
		Code     int    `json:"code"`
		Msg      string `json:"msg"`
		ImageURL string `json:"imageUrl"`
		DateTime string `json:"datetime"`
	}
	var response apiResponse
	resp, err := getRequest(apiURL, [][]string{})
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(resp, &response); err != nil {
		return nil, err
	}
	if response.Code != 200 {
		return nil, fmt.Errorf("API Error")
	}
	imgURL := response.ImageURL
	imgResp, err := getRequest(imgURL, [][]string{})
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(imgResp), nil
}

func getRequest(url string, queryList [][]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	for _, queryItem := range queryList {
		if len(queryItem) != 2 {
			return nil, fmt.Errorf("%v is not a valid query", queryItem)
		}
		q.Add(queryItem[0], queryItem[1])
	}
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
