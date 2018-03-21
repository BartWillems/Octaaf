package main

import (
	"bytes"
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	humanize "github.com/dustin/go-humanize"

	"github.com/PuerkitoBio/goquery"
	"github.com/o1egl/govatar"
	"github.com/tidwall/gjson"
	"gopkg.in/telegram-bot-api.v4"
)

func sendRoll(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	rand.Seed(time.Now().UnixNano())
	roll := strconv.Itoa(rand.Intn(9999999999-1000000000) + 1000000000)
	points := [9]string{"👌 Dubs", "🙈 Trips", "😱 Quads", "🤣😂 Penta", "👌👌🤔🤔😂😂 Hexa", "🙊🙉🙈🐵 Septa", "🅱️Octa", "💯💯💯 El Niño"}
	var dubscount int8 = -1

	for i := len(roll) - 1; i > 0; i-- {
		if roll[i] == roll[i-1] {
			dubscount++
		} else {
			break
		}
	}

	if dubscount > -1 {
		roll = points[dubscount] + " " + roll
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, roll)
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}

func count(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("%v", message.MessageID))
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}

func m8Ball(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {

	if len(message.CommandArguments()) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Oi! You have to ask question hé 🖕")
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
		return
	}

	answers := [20]string{"👌 It is certain",
		"👌 It is decidedly so",
		"👌 Without a doubt",
		"👌 Yes definitely",
		"👌 You may rely on it",
		"👌 As I see it, yes",
		"👌 Most likely",
		"👌 Outlook good",
		"👌 Yes",
		"👌 Signs point to yes",
		"☝ Reply hazy try again",
		"☝ Ask again later",
		"☝ Better not tell you now",
		"☝ Cannot predict now",
		"☝ Concentrate and ask again",
		"🖕 Don't count on it",
		"🖕 My reply is no",
		"🖕 My sources say no",
		"🖕 Outlook not so good",
		"🖕 Very doubtful"}
	rand.Seed(time.Now().UnixNano())
	roll := rand.Intn(19)
	msg := tgbotapi.NewMessage(message.Chat.ID, answers[roll])
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}

func sendBodegem(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewLocation(message.Chat.ID, 50.8614773, 4.211304)
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}

func where(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	argument := strings.Replace(message.CommandArguments(), " ", "+", -1)

	location, found := getLocation(argument)

	if !found {
		msg := getMessageConfig(message, "This place does not exist 🙈🙈🙈🤔🤔🤔")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewLocation(message.Chat.ID, location.lat, location.lng)
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}

func what(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	query := message.CommandArguments()
	resp, _ := http.Get(fmt.Sprintf("https://api.duckduckgo.com/?q=%v&format=json&no_html=1&skip_disambig=1", query))
	body, _ := ioutil.ReadAll(resp.Body)

	result := gjson.Get(string(body), "AbstractText").String()

	if len(result) == 0 {
		msg := getMessageConfig(message, fmt.Sprintf("What is this *%v* you speak of? 🤔", query))
		bot.Send(msg)
		return
	}

	msg := getMessageConfig(message, fmt.Sprintf("*%v:* %v", query, result))
	bot.Send(msg)
}

func weather(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	argument := strings.Replace(message.CommandArguments(), " ", "+", -1)

	location, found := getLocation(argument)

	if !found {
		msg := getMessageConfig(message, "No data found 🙈🙈🙈🤔🤔🤔")
		bot.Send(msg)
		return
	}

	resp, _ := http.Get(fmt.Sprintf("https://graphdata.buienradar.nl/forecast/json/?lat=%v&lon=%v", location.lat, location.lng))
	body, _ := ioutil.ReadAll(resp.Body)
	weatherJSON := string(body)

	reply := "No weather data found."

	forecasts := gjson.Get(weatherJSON, "forecasts").Array()
	raining := false

	if len(forecasts) > 0 {
		reply = "☀️☀️☀️ It's not going to rain in " + message.CommandArguments()
		if forecasts[0].Get("precipation").Num > 0 {
			reply = "🌧🌧🌧 It's now raining in " + message.CommandArguments()
			raining = true
		}
	}

	for _, forecast := range forecasts {
		if raining && forecast.Get("precipation").Num == 0 {
			reply += ", but it's expected to stop "
			rain, err := dateparse.ParseAny(forecast.Get("datetime").String())
			if err != nil {
				reply += " in " + forecast.Get("datetime").String()
			} else {
				reply += humanize.Time(rain)
			}
			break
		} else if forecast.Get("precipation").Num > 0 {
			rain, err := dateparse.ParseAny(forecast.Get("datetime").String())
			if err != nil {
				reply = "🌦🌦🌦 Expected rain from " + forecast.Get("datetime").String()
			} else {
				reply = "🌦🌦🌦 Expected rain " + humanize.Time(rain)
			}
			break
		}
	}

	msg := getMessageConfig(message, "*Weather:* "+reply)
	bot.Send(msg)
}

func sendAvatar(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	img, err := govatar.GenerateFromUsername(govatar.MALE, message.From.UserName)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	buf := new(bytes.Buffer)
	png.Encode(buf, img)

	bytes := tgbotapi.FileBytes{Name: "avatar.png", Bytes: buf.Bytes()}
	msg := tgbotapi.NewPhotoUpload(message.Chat.ID, bytes)
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}

func bol(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	bolURL := "https://www.bol.com/nl/nieuwsbrieven.html?country=BE"
	doc, _ := goquery.NewDocument(bolURL)
	token := "bogusTokenValue"

	doc.Find(".newsletter_subscriptions input").Each(func(i int, node *goquery.Selection) {
		name, found := node.Attr("name")
		if found && name == "token" {
			token, _ = node.Attr("value")
		}
	})

	http.PostForm(bolURL,
		url.Values{
			"emailAddress":             {message.CommandArguments()},
			"subscribedNewsLetters[0]": {"DAGAANBIEDINGEN"},
			"subscribedNewsLetters[1]": {"SOFT_OPTIN"},
			"subscribedNewsLetters[2]": {"HARD_OPTIN"},
			"subscribedNewsLetters[3]": {"B2B"},
			"token":                    {token},
			"updateNewsletters":        {"Voorkeuren+opslaan"}})

	msg := getMessageConfig(message, fmt.Sprintf("Succesfully subscribed *%v* to the bol.com mailing lists!", message.CommandArguments()))
	bot.Send(msg)
}

func search(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	if len(message.CommandArguments()) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "What do you expect me to do? 🤔🤔🤔🤔")
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
		return
	}

	// Basic url that disables ads
	url := "https://duckduckgo.com/lite?k1=-1&q=" + message.CommandArguments()

	if message.Command() == "search_nsfw" {
		url += "&kp=-2"
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		msg := getMessageConfig(message, "Uh oh, server error 🤔")
		bot.Send(msg)
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.89 Safari/537.36")
	resp, _ := client.Do(req)

	doc, _ := goquery.NewDocumentFromReader(resp.Body)

	url, found := doc.Find(".result-link").First().Attr("href")

	if found {
		msg := getMessageConfig(message, url)
		bot.Send(msg)
		return
	}

	msg := getMessageConfig(message, "I found nothing 😱😱😱")
	bot.Send(msg)
}

func sendStallman(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	var url = "https://stallman.org/photos/rms-working/"

	doc, err := goquery.NewDocument(url)

	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Stallman went bork? 🤔🤔🤔🤔")
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
		return
	}

	var pages []string

	doc.Find("img").Each(func(i int, token *goquery.Selection) {
		url, exists := token.Parent().Attr("href")
		if exists {
			pages = append(pages, url)
		}
	})

	if len(pages) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "No stallman found... 🤔🤔🤔🤔")
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
		return
	}

	rand.Seed(time.Now().UnixNano())
	roll := rand.Intn(len(pages))

	log.Printf("Roll: %v", pages[roll])

	doc, err = goquery.NewDocument(url + pages[roll])

	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Stallman went bork? 🤔🤔🤔🤔")
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
		return
	}

	image, _ := doc.Find("img").First().Parent().Attr("href")

	log.Printf("Image: %v", image)
	log.Printf("Url: %v", url+path.Base(image))

	res, _ := http.Get(url + path.Base(image))

	content, err := ioutil.ReadAll(res.Body)

	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Stallman parser error... 🤔🤔🤔🤔")
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
		return
	}

	bytes := tgbotapi.FileBytes{Name: "stally.jpg", Bytes: content}
	msg := tgbotapi.NewPhotoUpload(message.Chat.ID, bytes)

	msg.Caption = message.CommandArguments()
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}

func sendImage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	argument := strings.Replace(message.CommandArguments(), " ", "+", -1)
	if len(argument) == 0 {
		msg := getMessageConfig(message, fmt.Sprintf("What am I to do, @%v? 🤔🤔🤔🤔", message.From.UserName))
		bot.Send(msg)
		return
	}

	query := "http://images.google.com/search?tbm=isch&q=" + argument

	if message.Command() == "img_sfw" {
		query += "&safe=on"
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", query, nil)

	if err != nil {
		msg := getMessageConfig(message, "Uh oh, server error 🤔")
		bot.Send(msg)
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.89 Safari/537.36")
	resp, err := client.Do(req)

	if err != nil {
		msg := getMessageConfig(message, fmt.Sprintf("Something went wrong while searching this query: `%v`", message.CommandArguments()))
		bot.Send(msg)
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		msg := getMessageConfig(message, fmt.Sprintf("Something went wrong while parsing this query response: `%v`", message.CommandArguments()))
		bot.Send(msg)
		return
	}

	var images []string

	doc.Find(".rg_di .rg_meta").Each(func(i int, token *goquery.Selection) {
		imageJSON := token.Text()
		imageURL := gjson.Get(imageJSON, "ou").String()

		if len(imageURL) > 0 {
			images = append(images, imageURL)
		}
	})

	timeout := time.Duration(2 * time.Second)
	client = &http.Client{
		Timeout: timeout,
	}

	for _, url := range images {

		res, err := client.Get(url)

		if err != nil {
			continue
		}

		content, err := ioutil.ReadAll(res.Body)

		if err != nil {
			continue
		}

		bytes := tgbotapi.FileBytes{Name: "image.jpg", Bytes: content}
		msg := tgbotapi.NewPhotoUpload(message.Chat.ID, bytes)

		msg.Caption = message.CommandArguments()
		msg.ReplyToMessageID = message.MessageID
		_, e := bot.Send(msg)

		if e == nil {
			return
		}
	}

	msg := getMessageConfig(message, "I did not find images for the query: `"+message.CommandArguments()+"`")
	bot.Send(msg)
}
