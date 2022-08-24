package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type outgoing struct {
	Event     string `json:"event"`
	EventType string `json:"event_type"`
	AppId     string `json:"app_id"`
	UserId    string `json:"user_id"`
	MessageId string `json:"message_id"`

	PageTitle       string               `json:"page_title"`
	PageUrl         string               `json:"page_url"`
	BrowserLanguage string               `json:"browser_language"`
	ScreenSize      string               `json:"screen_size"`
	Attributes      map[string]typeValue `json:"attributes"`
	Traits          map[string]typeValue `json:"traits"`
}

type typeValue struct {
	Type  string
	Value string
}

func main() {
	app := fiber.New()
	app.Post("/", func(c *fiber.Ctx) error {
		inp := make(chan string)
		go worker(inp)

		inp <- string(c.Body())
		close(inp)
		return c.SendString("Done")
	})
	log.Fatal(app.Listen(":3000"))
}

func worker(inp <-chan string) {
	for i := range inp {
		response := outgoing{}
		var inputMap map[string]interface{}
		err := json.Unmarshal([]byte(i), &inputMap)
		if err != nil {
			panic(err)
		}
		i := 1
		att := "atrk"
		attributes := make(map[string]typeValue)
		for {
			if v, ok := inputMap[fmt.Sprintf("%v%v", att, i)]; ok {
				attributes[v.(string)] = typeValue{
					Type:  inputMap[fmt.Sprintf("atrt%v", i)].(string),
					Value: inputMap[fmt.Sprintf("atrv%v", i)].(string),
				}

			} else {
				break
			}
			i++
		}

		j := 1
		uatrk := "uatrk"
		traits := make(map[string]typeValue)
		for {
			if v, ok := inputMap[fmt.Sprintf("%v%v", uatrk, j)]; ok {
				traits[v.(string)] = typeValue{
					Type:  inputMap[fmt.Sprintf("uatrk%v", j)].(string),
					Value: inputMap[fmt.Sprintf("uatrv%v", j)].(string),
				}

			} else {
				break
			}
			j++
		}

		response.Event = inputMap["ev"].(string)
		response.EventType = inputMap["et"].(string)
		response.AppId = inputMap["id"].(string)
		response.UserId = inputMap["uid"].(string)
		response.MessageId = inputMap["mid"].(string)
		response.PageTitle = inputMap["t"].(string)
		response.PageUrl = inputMap["p"].(string)
		response.BrowserLanguage = inputMap["l"].(string)
		response.ScreenSize = inputMap["sc"].(string)
		response.Attributes = attributes
		response.Traits = traits
		json_data, err := json.Marshal(response)

		if err != nil {
			log.Fatal(err)
		}

		_, err = http.Post("https://webhook.site/ef213080-56d6-41d5-bafa-fba5afdfbab7", "application/json",
			bytes.NewBuffer(json_data))

		if err != nil {
			log.Fatal(err)
		}
	}
}

func ToJson(o interface{}) string {
	js, _ := json.Marshal(o)
	return string(js)
}
