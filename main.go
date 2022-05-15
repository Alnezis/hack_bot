package main

import (
	"encoding/json"
	"fmt"
	"github.com/fogleman/gg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"hack_bot/api"
	"hack_bot/app"
	"hack_bot/image"
	"hack_bot/validate"
	"image/color"
	"io/ioutil"
	"log"
	"net/http"
)

type Obj struct {
	Message     string `json:"message,omitempty"`
	Predictions []struct {
		Class      string  `json:"class"`
		Confidence float64 `json:"confidence"`
		Height     int     `json:"height"`
		Width      int     `json:"width"`
		X          float64 `json:"x"`
		Y          float64 `json:"y"`
	} `json:"predictions,omitempty"`
}

func main() {
	bot, err := tgbotapi.NewBotAPI(app.CFG.TokenTG)
	if err != nil {
		log.Panic(err)
	}

	//var numericKeyboard = tgbotapi.NewReplyKeyboard(
	//	tgbotapi.NewKeyboardButtonRow(
	//		tgbotapi.NewKeyboardButtonLocation("Поделиться местоположением"),
	//	),
	//)

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			if update.Message.Chat.ID != 388788625 {
				msg.Text = "Доступно только администраторам."
				bot.Send(msg)
			}

			if update.Message.Photo != nil {
				i := update.Message.Photo[len(update.Message.Photo)-1]
				file := tgbotapi.FileConfig{
					FileID: i.FileID,
				}

				if !validate.Validate(bot, i, msg) {
					continue
				}

				e, err := bot.GetFile(file)
				if err != nil {
					api.CheckErrInfo(err, "GetFile")
					continue
				}

				link := e.Link(bot.Token)
				msg.Text = link
				urlCheck := "https://detect.roboflow.com/hard-hat-sample-nh62g/2?api_key=GuzcZQ0lFTQHCzwDxEOq&image=" + link
				msg.Text = "Идет обработка..."

				bot.Send(msg)
				client := &http.Client{}
				req, err := http.NewRequest("POST", urlCheck, nil)
				if err != nil {
					log.Fatal(err)
				}

				resp, err := client.Do(req)
				if err != nil {
					log.Fatal(err)
				}
				defer resp.Body.Close()

				bodyText, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Fatal(err)
				}
				var obj Obj
				err = json.Unmarshal(bodyText, &obj)
				api.CheckErrInfo(err, "Unmarshal")

				if obj.Message != "" {
					msg.Text = "Error: " + obj.Message
					bot.Send(msg)
					continue
				}

				if obj.Predictions != nil {
					fmt.Println(len(obj.Predictions))
					if len(obj.Predictions) == 0 {
						msg.Text = "Поврежений не обнаружено."
						bot.Send(msg)
						continue
					}

					fileName := api.RandString(10) + ".jpg"
					err := image.DownloadFile(link, fileName)
					if err != nil {
						log.Fatal(err)
					}

					loadImage, err := gg.LoadImage(fileName)
					if err != nil {
						api.CheckErrInfo(err, "LoadImage")
						return
					}

					dc := gg.NewContextForImage(loadImage)

					caption := ""
					for _, v := range obj.Predictions {
						q := (dc.Width() + dc.Height()) / 11

						fmt.Println(q)
						image.DrawCircle(dc, int(v.X), int(v.Y), q, color.RGBA{255, 0, 0, 255})
						caption += fmt.Sprintf("x: %d, y: %d - %d", int(v.X), int(v.Y), int(v.Confidence*100)) + "%\n"
					}

					//x1 := obj.Predictions[0].X
					//y1 := obj.Predictions[0].Y + 10
					//x2 := obj.Predictions[0].X
					//y2 := obj.Predictions[0].Y + 10
					//r := 3.0
					//g := 3.0
					//b := 3.0
					//a := rand.Float64()*0.5 + 0.5
					//w := 60.0
					//dc.SetRGBA(r, g, b, a)
					//dc.SetLineWidth(w)
					//dc.DrawLine(x1, y1, x2, y2)
					//dc.Stroke()
					dc.SavePNG(fileName)

					msg := tgbotapi.NewPhoto(msg.ChatID, tgbotapi.FilePath(fileName))
					msg.Caption = caption
					if _, err = bot.Send(msg); err != nil {
						panic(err)
					}

				}
			}
		}

		//if update.Message.Location != nil {
		//	var url = "https://yandex.ru/maps?whatshere%5Bpoint%5D=" + api.ToString(update.Message.Location.Longitude) + "%2C" + api.ToString(update.Message.Location.Latitude) + "&whatshere%5Bzoom%5D=16.037464&ll=38.997331%2C45.03968899961266&z=16.037464"
		//	var numericKeyboard2 = tgbotapi.NewInlineKeyboardMarkup(
		//		tgbotapi.NewInlineKeyboardRow(
		//			tgbotapi.NewInlineKeyboardButtonURL("Открыть карту", url),
		//		),
		//	)
		//	msg.ReplyMarkup = numericKeyboard2
		//}
		//
		//msg.Text = "test"
		//// Send the message.
		//if _, err = bot.Send(msg); err != nil {
		//	panic(err)
		//}
		//} else if update.CallbackQuery != nil {
		//	// Respond to the callback query, telling Telegram to show the user
		//	// a message with the data received.
		//	fmt.Println(update.CallbackQuery.Data)
		//	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
		//	if _, err := bot.Request(callback); err != nil {
		//		panic(err)
		//	}
		//
		//	// And finally, send a message containing the data received.
		//	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
		//	if _, err := bot.Send(msg); err != nil {
		//		panic(err)
		//	}
		//}
	}
}
