package botсont

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

var m = make(map[string] bool)
var mapINT = make(map[int64] bool)
func checkFor(token string, ID int64) bool {
	if m[token] {
		data, err := os.ReadFile("internal/overBotBuild/checkToken.txt")
		if err != nil {
			fmt.Println(err)
			return false
		}

		var useMap = make(map[string] int64)
		err = json.Unmarshal(data, &useMap)

		if err != nil {
			fmt.Println(err)
			return false
		}

		if useMap[token] != 0 {
			return false
		} else {
			useMap[token] = ID
			data, err := json.Marshal(useMap)
			if err != nil {
				log.Panic(err)
				return false
			}
		_ = os.WriteFile("internal/overBotBuild/checkToken.txt", data, 0644)
		return true

		}
	} else {
		return false
	}
}

func unparse() error{
	data, err := os.ReadFile("internal/overBotBuild/mapOfUsers.txt")
	if err != nil {
		return err
	}

	er := json.Unmarshal(data, &m)
	if er != nil {
		return er
	}

	return nil
}
func unparseTwo() error {
	data, err := os.ReadFile("internal/overBotBuild/mapOfAuth.txt")
	if err != nil {
		return err
	}

	er := json.Unmarshal(data, &mapINT)
	if er != nil {
		return er
	}
	return nil
}

func BotOper() {
	fmt.Println("Launching")

	er := unparse()
	if er != nil {
		fmt.Println(er)
		return
	}
	var isAuth = false
	er = unparseTwo()
	if er != nil {
	fmt.Println(er)
	return
	}

	var authTrue = false
	data, err := os.ReadFile("internal/overBotBuild/TelegramToken.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	token := string(data)
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		fmt.Println("JERE")
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		if !update.Message.IsCommand() {
			switch true {

			case authTrue:
				{
					fmt.Println("Ok, auth token awaits")
					if checkFor(update.Message.Text, int64(update.Message.From.ID)) {
						authTrue = false
						isAuth = true
						msg.Text = "Успешная авторизация"
					} else {
						msg.Text = "Провал("
					}
				}

			default:
				msg.Text = "Ознакомьтесь со списком команд"
			}
		} else {
			switch update.Message.Command() {
			case "authenticate":
				authTrue = true
				msg.Text = "Ok, authenticate"
			case "start":
				msg.Text = "Добро пожаловать! Войдите в систему, используя команду authenticate и пришлите логин."
			case "add":
				if isAuth  {
					msg.Text = "Хорошо, пришлите следующим сообщением номер аудитории!"
				} else {
					msg.Text = "Вы не авторизовались!"
				}
			default:
				msg.Text = "Please enter a valid command"
			}
		}




		_, err = bot.Send(msg)
		if err != nil {
			log.Panic(err)
		}

	}
}
