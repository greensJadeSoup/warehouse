package cp_middleware

import (
	"bytes"
	"warehouse/v5-go-component/cp_log"
	"encoding/json"
	"net/http"
)

// DingTalk 钉钉机器人
type DingTalk struct {
	WebHoop string
	Data    []byte
}

// NewDingTalk 钉钉机器人
func NewDingTalk() *DingTalk {
	return &DingTalk{}
}

// GetWebHoop 机器人地址
func (*DingTalk) GetWebHoop() string {
	return `https://oapi.dingtalk.com/robot/send?access_token=2f40779609a201996b365d20c689818bc10dec0dac1f0effa2876a7d00146a0d`
}

// Send 发送
func (d *DingTalk) Send() error {
	client := &http.Client{}

	postBytesReader := bytes.NewReader(d.Data)
	d.WebHoop = d.GetWebHoop()
	req, err := http.NewRequest("POST", d.WebHoop, postBytesReader)

	if err != nil {
		cp_log.Info(err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	// if response.StatusCode == 200 {
	// 	body, _ := ioutil.ReadAll(response.Body)
	// 	bodystr := string(body)
	// 	cp_log.Info(bodystr)
	// }

	return nil
}

// DingTalkText 消息
type DingTalkText struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

// NewDingTalkText 新的文本消息
func NewDingTalkText(text string) ([]byte, error) {
	s := &DingTalkText{}
	s.MsgType = "text"
	s.Text.Content = text

	return json.Marshal(s)
}

// DingTalkMarkDown MarkDown 消息
type DingTalkMarkDown struct {
	MsgType  string `json:"msgtype"`
	MarkDown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
	At struct {
		AtMobiles []string `json:"atMobiles"`
		IsAtAll   bool     `json:"isAtAll"`
	}
}

// NewDingTalkMarkDown markdown消息
func NewDingTalkMarkDown(title, text string) ([]byte, error) {
	s := &DingTalkMarkDown{}
	s.MsgType = "markdown"
	s.MarkDown.Title = title
	s.MarkDown.Text = text

	return json.Marshal(s)
}
