package bots

import (
	"fmt"
	"github.com/andatoshiki/toshiki-e5subot/config"
	"github.com/andatoshiki/toshiki-e5subot/service/srv_client"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
	"time"
)

var (
	UserStatus       map[int64]int
	UserClientId     map[int64]string
	UserClientSecret map[int64]string
)

const (
	StatusNone = iota
	StatusBind1
	StatusBind2
)

func init() {
	UserStatus = make(map[int64]int)
	UserClientId = make(map[int64]string)
	UserClientSecret = make(map[int64]string)
}

func bMy(m *tb.Message) {
	clients := srv_client.GetClients(m.Chat.ID)
	var inlineKeys [][]tb.InlineButton
	for _, client := range clients {
		inlineBtn := tb.InlineButton{
			Unique: "my" + strconv.Itoa(client.ID),
			Text:   client.Alias,
			Data:   strconv.Itoa(client.ID),
		}
		bot.Handle(&inlineBtn, bMyInlineBtn)
		inlineKeys = append(inlineKeys, []tb.InlineButton{inlineBtn})
	}

	bot.Send(m.Chat,
		fmt.Sprintf("👉🏻 Please selet an account profile to view information details\n\nOwned account counts: %d/%d", len(srv_client.GetClients(m.Chat.ID)), config.BindMaxNum),
		&tb.ReplyMarkup{InlineKeyboard: inlineKeys})
}
func bMyInlineBtn(c *tb.Callback) {
	id, _ := strconv.Atoi(c.Data)
	client, err := srv_client.GetClient(id)
	if err != nil {
		bot.Send(c.Message.Chat, "❎ Failed to fetch account information details")
		return
	}
	bot.Send(c.Message.Chat,
		fmt.Sprintf("🔎 Account Details\nAlias:%s\nms_id: %s\nclient_id: %s\nclient_secret: %s\nLast updated: %s",
			client.Alias,
			client.MsId,
			client.ClientId,
			client.ClientSecret,
			time.Unix(client.Uptime, 0).Format("2006-01-02 15:04:05"),
		),
	)
	bot.Respond(c)
}

func bOnText(m *tb.Message) {
	switch UserStatus[m.Chat.ID] {
	case StatusBind1:
		bBind1(m)
	case StatusBind2:
		bBind2(m)
	default:
		bot.Send(m.Chat, "❓ Send /help to view list of commands and bot descriptions.")
	}
}
