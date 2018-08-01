package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"octaaf/models"
	"octaaf/scrapers"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	humanize "github.com/dustin/go-humanize"
	"github.com/go-redis/cache"
	"github.com/tidwall/gjson"
	"gopkg.in/telegram-bot-api.v4"
)

func sendRoll(message *tgbotapi.Message) {
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
	reply(message, roll)
}

func count(message *tgbotapi.Message) {
	reply(message, fmt.Sprintf("%v", message.MessageID))
}

func whoami(message *tgbotapi.Message) {
	reply(message, fmt.Sprintf("%v", message.From.ID))
}

func m8Ball(message *tgbotapi.Message) {

	if len(message.CommandArguments()) == 0 {
		reply(message, "Oi! You have to ask question hé 🖕")
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
	Octaaf.Send(msg)
}

func sendBodegem(message *tgbotapi.Message) {
	msg := tgbotapi.NewLocation(message.Chat.ID, 50.8614773, 4.211304)
	msg.ReplyToMessageID = message.MessageID
	Octaaf.Send(msg)
}

func where(message *tgbotapi.Message) {
	argument := strings.Replace(message.CommandArguments(), " ", "+", -1)

	location, found := scrapers.GetLocation(argument)

	if !found {
		reply(message, "This place does not exist 🙈🙈🙈🤔🤔�")
		return
	}

	msg := tgbotapi.NewLocation(message.Chat.ID, location.Lat, location.Lng)
	msg.ReplyToMessageID = message.MessageID
	Octaaf.Send(msg)
}

func what(message *tgbotapi.Message) {
	query := message.CommandArguments()
	resp, err := http.Get(fmt.Sprintf("https://api.duckduckgo.com/?q=%v&format=json&no_html=1&skip_disambig=1", query))
	if err != nil {
		reply(message, "Just what is this? 🤔")
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		reply(message, "Just what is this? 🤔")
		return
	}

	result := gjson.Get(string(body), "AbstractText").String()

	if len(result) == 0 {
		reply(message, fmt.Sprintf("What is this %v you speak of? 🤔", Markdown(query, mdbold)))
		return
	}

	reply(message, fmt.Sprintf("%v: %v", Markdown(query, mdbold), result))
}

func weather(message *tgbotapi.Message) {
	weather, found := scrapers.GetWeatherStatus(message.CommandArguments())
	if !found {
		reply(message, "No data found 🙈🙈🙈🤔🤔🤔")
	} else {
		reply(message, "*Weather:* "+weather)
	}
}

func search(message *tgbotapi.Message) {
	if len(message.CommandArguments()) == 0 {
		reply(message, "What do you expect me to do? 🤔🤔🤔🤔")
		return
	}

	url, found := scrapers.Search(message.CommandArguments(), message.Command() == "search_nsfw")

	if found {
		reply(message, MDEscape(url))
		return
	}

	reply(message, "I found nothing 😱😱😱")
}

func sendStallman(message *tgbotapi.Message) {

	image, err := scrapers.GetStallman()

	if err != nil {
		reply(message, "Stallman went bork? 🤔🤔🤔🤔")
		return
	}

	bytes := tgbotapi.FileBytes{Name: "stally.jpg", Bytes: image}
	msg := tgbotapi.NewPhotoUpload(message.Chat.ID, bytes)

	msg.Caption = message.CommandArguments()
	msg.ReplyToMessageID = message.MessageID
	Octaaf.Send(msg)
}

func sendImage(message *tgbotapi.Message) {
	var images []string
	var err error
	key := fmt.Sprintf("images_%v", message.Chat.ID)
	if message.Command() != "more" {
		if len(message.CommandArguments()) == 0 {
			reply(message, fmt.Sprintf("What am I to do, @%v? 🤔🤔🤔🤔", message.From.UserName))
			return
		}

		images, err = scrapers.GetImages(message.CommandArguments(), message.Command() == "img_sfw")
		if err != nil {
			reply(message, "Something went wrong!")
			return
		}

		Codec.Set(&cache.Item{
			Key:        key,
			Object:     images,
			Expiration: 0,
		})
	} else {
		if err := Codec.Get(key, &images); err != nil {
			reply(message, "I can't fetch them for you right now.")
			return
		}

		// Randomly order images for a different /more
		for i := range images {
			j := rand.Intn(i + 1)
			images[i], images[j] = images[j], images[i]
		}
	}

	timeout := time.Duration(2 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}

	for _, url := range images {

		res, err := client.Get(url)

		if err != nil {
			continue
		}

		defer res.Body.Close()

		content, err := ioutil.ReadAll(res.Body)

		if err != nil {
			continue
		}

		bytes := tgbotapi.FileBytes{Name: "image.jpg", Bytes: content}
		msg := tgbotapi.NewPhotoUpload(message.Chat.ID, bytes)

		msg.Caption = message.CommandArguments()
		msg.ReplyToMessageID = message.MessageID
		_, e := Octaaf.Send(msg)

		if e == nil {
			return
		}
	}

	reply(message, "I did not find images for the query: `"+message.CommandArguments()+"`")
}

func xkcd(message *tgbotapi.Message) {
	image, err := scrapers.GetXKCD()

	if err != nil {
		reply(message, "Failed to parse XKCD image")
		return
	}

	bytes := tgbotapi.FileBytes{Name: "image.jpg", Bytes: image}
	msg := tgbotapi.NewPhotoUpload(message.Chat.ID, bytes)

	msg.Caption = message.CommandArguments()
	msg.ReplyToMessageID = message.MessageID
	Octaaf.Send(msg)
}

func doubt(message *tgbotapi.Message) {
	msg := tgbotapi.NewPhotoUpload(message.Chat.ID, "assets/doubt.jpg")
	msg.ReplyToMessageID = message.MessageID
	Octaaf.Send(msg)
}

func quote(message *tgbotapi.Message) {
	// Fetch a random quote
	if message.ReplyToMessage == nil {
		quote := models.Quote{}

		var err error

		if len(message.CommandArguments()) > 0 {
			query := DB.Where("chat_id = ? AND quote ilike '%' || ? || '%'", message.Chat.ID, message.CommandArguments())
			err = query.Order("random()").Limit(1).First(&quote)
		} else {
			err = DB.Where("chat_id = ?", message.Chat.ID).Order("random()").Limit(1).First(&quote)
		}

		log.Printf("ERROR %s", err)

		if err != nil {
			reply(message, "No quote found boi")
			return
		}

		config := tgbotapi.ChatConfigWithUser{
			ChatID:             message.Chat.ID,
			SuperGroupUsername: "",
			UserID:             quote.UserID}

		user, userErr := Octaaf.GetChatMember(config)

		if userErr != nil {
			reply(message, quote.Quote)
		} else {
			msg := fmt.Sprintf("\"%v\"", Markdown(quote.Quote, mdquote))
			msg += fmt.Sprintf(" \n    ~@%v", MDEscape(user.User.UserName))
			reply(message, msg)
		}

		return
	}

	// Unable to store this quote
	if message.ReplyToMessage.Text == "" {
		reply(message, "No text found in the comment. Not saving the quote!")
		return
	}

	err := DB.Save(&models.Quote{
		Quote:  message.ReplyToMessage.Text,
		UserID: message.ReplyToMessage.From.ID,
		ChatID: message.Chat.ID})

	if err != nil {
		reply(message, "Unable to save the quote...")
		return
	}

	reply(message, "Quote successfully saved!")
}

func nextLaunch(message *tgbotapi.Message) {
	res, err := http.Get("https://launchlibrary.net/1.3/launch?next=5&mode=verbose")

	if err != nil {
		reply(message, "Unable to fetch launch data")
		return
	}

	defer res.Body.Close()

	launchJSON, err := ioutil.ReadAll(res.Body)

	if err != nil {
		reply(message, "Unable to fetch launch data")
		return
	}

	launches := gjson.Get(string(launchJSON), "launches").Array()

	var msg = "*Next 5 launches:*"

	layout := "January 2, 2006 15:04:05 MST"

	for index, launch := range launches {
		whenStr := launch.Get("net").String()
		when, err := time.Parse(layout, whenStr)

		msg += fmt.Sprintf("\n*%v*: %v", index+1, MDEscape(launch.Get("name").String()))

		if err != nil {
			msg += fmt.Sprintf("\n	  %v", Markdown(whenStr, mdcursive))
		} else {
			msg += fmt.Sprintf("\n	  %v", Markdown(humanize.Time(when), mdcursive))
		}

		vods := launch.Get("vidURLs").Array()

		if len(vods) > 0 {
			msg += "\n    " + MDEscape(vods[0].String())
		}
	}

	reply(message, msg)
}

func issues(message *tgbotapi.Message) {
	res, err := http.Get("https://api.github.com/repos/bartwillems/Octaaf/issues?state=open")

	if err != nil {
		reply(message, "Unable to fetch open issues")
		return
	}

	defer res.Body.Close()

	issuesJSON, err := ioutil.ReadAll(res.Body)

	if err != nil {
		reply(message, "Unable to fetch open issues")
		return
	}

	issues := gjson.ParseBytes(issuesJSON)

	var msg = "*Octaaf issues:*"

	var count int

	issues.ForEach(func(key, value gjson.Result) bool {
		count++
		msg += fmt.Sprintf("\n*%v: %v*", count, MDEscape(value.Get("title").String()))
		msg += fmt.Sprintf("\n    *url:* %v", Markdown(value.Get("url").String(), mdcursive))
		msg += fmt.Sprintf("\n    *creator:* %v", Markdown(value.Get("user.login").String(), mdcursive))
		return true
	})

	reply(message, msg)
}

func kaliRank(message *tgbotapi.Message) {
	if message.Chat.ID != KaliID {
		reply(message, "You are not allowed!")
		return
	}

	kaliRank := []models.MessageCount{}
	err := DB.Order("diff DESC").Limit(5).All(&kaliRank)

	if err != nil {
		reply(message, "Unable to fetch the kali rankings")
		return
	}

	var msg = "*Kali rankings:*"
	for index, rank := range kaliRank {
		msg += fmt.Sprintf("\n`#%v:` *%v messages*   _~%v_", index+1, rank.Diff, rank.CreatedAt.Format("Monday, 2 January 2006"))
	}

	reply(message, msg)
}

func iasip(message *tgbotapi.Message) {
	server := "http://159.89.14.97:6969"

	res, err := http.Get(server)
	if err != nil {
		reply(message, "Unable to fetch iasip quote...you goddamn bitch you..")
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		reply(message, "Unable to fetch iasip quote...you goddamn bitch you..")
		return
	}

	reply(message, string(body))
}

func reported(message *tgbotapi.Message) {
	if message.Chat.ID != KaliID {
		reply(message, "Yeah well, you need to update to Strontbot Enterprise edition for Workgroups to use this feature.")
		return
	}

	reportCount, err := DB.Count(models.Report{})

	if err != nil {
		reply(message, "I can't seem to be able to count the reports.")
		return
	}

	config := tgbotapi.ChatConfigWithUser{
		ChatID:             message.Chat.ID,
		SuperGroupUsername: "",
		UserID:             ReporterID}

	reporter, err := Octaaf.GetChatMember(config)

	if err != nil {
		reply(message, fmt.Sprintf("So far, %v people have been reported by Dieter", reportCount))
	} else {
		reply(message, MDEscape(fmt.Sprintf("So far, %v people have been reported by: @%v", reportCount, reporter.User.UserName)))
	}
}
