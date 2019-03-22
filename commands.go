package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"math/rand"
	"net/http"
	"octaaf/cache"
	"octaaf/markdown"
	"octaaf/models"
	"octaaf/scrapers"
	"octaaf/trump"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	humanize "github.com/dustin/go-humanize"
	"github.com/gobuffalo/envy"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

func changelog(message *OctaafMessage) error {
	if Version == "" {
		return message.Reply("Current version not found, check the changelog here: " + GitURI + "/tags")
	}

	fetchSpan := message.Span.Tracer().StartSpan(
		"Fetching release info...",
		opentracing.ChildOf(message.Span.Context()),
	)
	releaseUrl := "https://gitlab.com/api/v4/projects/bartwillems%2Foctaaf/repository/tags/" + Version
	res, err := http.Get(releaseUrl)

	fetchSpan.Finish()

	if err != nil {
		fetchSpan.SetTag("error", true)
		fetchSpan.SetBaggageItem("error", err.Error())
		return message.Reply(fmt.Sprintf("Unable to fetch release info: %v", releaseUrl))
	}

	defer res.Body.Close()

	releaseJSON, err := ioutil.ReadAll(res.Body)

	if err != nil {
		message.LogError("Changelog failure: " + err.Error())
		return message.Reply("Unable to fetch release info")
	}

	release := gjson.ParseBytes(releaseJSON)

	msg := fmt.Sprintf("*%v*\n%v", release.Get("release.tag_name").String(), release.Get("release.description").String())
	return message.Reply(fmt.Sprintf("%v\n%v/tags/%v", msg, GitURI, Version))
}

func all(message *OctaafMessage) error {
	userIDSpan := message.Span.Tracer().StartSpan(
		"Fetch user ids from redis",
		opentracing.ChildOf(message.Span.Context()),
	)
	members := Redis.SMembers(fmt.Sprintf("members_%v", message.Chat.ID)).Val()

	userIDSpan.Finish()

	if len(members) == 0 {
		return message.Reply("I'm afraid I can't do that.")
	}

	usernamesSpan := message.Span.Tracer().StartSpan(
		"Load usernames",
		opentracing.ChildOf(message.Span.Context()),
	)

	// used to load the usernames in goroutines
	var wg sync.WaitGroup
	var response string
	// Get the members' usernames
	for _, member := range members {
		memberID, err := strconv.Atoi(member)

		if err != nil {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			username, err := getUserName(memberID, message.Chat.ID)
			if err == nil {
				response += fmt.Sprintf("@%v ", username)
			}
		}()
	}

	wg.Wait()
	usernamesSpan.Finish()
	return message.ReplyPlainText(fmt.Sprintf("%v %v", response, message.CommandArguments()))
}

func remind(message *OctaafMessage) error {
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)

	r, err := w.Parse(message.CommandArguments(), time.Now())

	if err != nil {
		message.LogError("Unable to parse reminder: " + err.Error())
		return message.Reply("Unable to parse")
	}

	if r == nil {
		message.LogError("No reminder found for message: " + message.CommandArguments())
		return message.Reply("No reminder found")
	}

	reminder := models.Reminder{
		ChatID:    message.Chat.ID,
		UserID:    message.From.ID,
		MessageID: message.MessageID,
		Message:   message.CommandArguments(),
		Deadline:  r.Time,
		Executed:  false}

	go startReminder(reminder)

	loc, err := time.LoadLocation(envy.Get("TZ", "Europe/Brussels"))
	var deadline string
	if err != nil {
		deadline = fmt.Sprintf("%s", reminder.Deadline)
	} else {
		deadline = reminder.Deadline.In(loc).Format("Monday January 2 15:04:05 2006")
	}

	return message.Reply(fmt.Sprintf("Reminder saved for %v!", deadline))
}

func sendRoll(message *OctaafMessage) error {
	rand.Seed(time.Now().UnixNano())
	roll := strconv.Itoa(rand.Intn(9999999999-1e9) + 1e9)
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
	return message.Reply(roll)
}

func count(message *OctaafMessage) error {
	// Get the 100K count (eg 5 for >=500k <600K)
	major := int(message.MessageID / 1e5)

	// Eg: 5 * 100K + 100K = 600K - messageID (eg 599K) < 5000
	if (major*1e5+1e5)-message.MessageID < 5000 {
		return message.Reply("Not now kimosabi")
	}
	return message.Reply(fmt.Sprintf("%v", message.MessageID))
}

func whoami(message *OctaafMessage) error {
	return message.Reply(fmt.Sprintf("%v", message.From.ID))
}

func m8Ball(message *OctaafMessage) error {

	if len(message.CommandArguments()) == 0 {
		return message.Reply("Oi! You have to ask question hé 🖕")
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
	return message.Reply(answers[roll])
}

func sendBodegem(message *OctaafMessage) error {
	return message.Reply(scrapers.Location{Lat: 50.8614773, Lng: 4.211304})
}

func where(message *OctaafMessage) error {
	argument := strings.Replace(message.CommandArguments(), " ", "+", -1)

	span := message.Span.Tracer().StartSpan(
		"Fetch location",
		opentracing.ChildOf(message.Span.Context()),
	)
	location, found := scrapers.GetLocation(argument, settings.Google.APIKEY)
	span.Finish()

	if !found {
		return message.Reply("This place does not exist 🙈🙈🙈🤔🤔�")
	}

	go DB.Save(&models.LocationHistory{
		ChatID:    message.Chat.ID,
		MessageID: message.MessageID,
		UserID:    message.From.ID,
		Lat:       location.Lat,
		Lng:       location.Lng,
		Name:      message.CommandArguments(),
	})

	return message.Reply(location)
}

func what(message *OctaafMessage) error {
	span := message.Span.Tracer().StartSpan(
		"Trying to explain something...",
		opentracing.ChildOf(message.Span.Context()),
	)
	query := message.CommandArguments()
	result, found, err := scrapers.What(query)
	span.Finish()

	if err != nil {
		message.LogError(fmt.Sprintf("Unable to explain '%v'. Error: %v", query, err))
		return message.Reply("I do not know, something went bork...")
	}

	if !found {
		return message.Reply("That is forbidden knowledge.")
	}

	return message.Reply(fmt.Sprintf("%v: %v", markdown.Bold(query), result))
}

func weather(message *OctaafMessage) error {
	weatherSpan := message.Span.Tracer().StartSpan(
		"Fetching weather status...",
		opentracing.ChildOf(message.Span.Context()),
	)
	weather, found := scrapers.GetWeatherStatus(message.CommandArguments(), settings.Google.APIKEY)

	weatherSpan.SetTag("found", found == true)
	weatherSpan.Finish()
	if !found {
		return message.Reply("No data found 🙈🙈🙈")
	}
	return message.Reply("*Weather:* " + weather)
}

func search(message *OctaafMessage) error {
	if len(message.CommandArguments()) == 0 {
		return message.Reply("What do you expect me to do? 🤔🤔🤔🤔")
	}

	searchSpan := message.Span.Tracer().StartSpan(
		"Searching on duckduckgo...",
		opentracing.ChildOf(message.Span.Context()),
	)

	url, found := scrapers.Search(message.CommandArguments(), message.Command() == "search_nsfw")

	searchSpan.SetTag("found", found == true)
	searchSpan.Finish()

	if found {
		return message.Reply(markdown.Escape(url))
	}
	return message.Reply("I found nothing 😱😱😱")
}

func sendStallman(message *OctaafMessage) error {
	fetchSpan := message.Span.Tracer().StartSpan(
		"Fetch a stallman",
		opentracing.ChildOf(message.Span.Context()),
	)

	image, err := scrapers.GetStallman()

	fetchSpan.Finish()

	if err != nil {
		message.LogError("Unable to retrieve stallman: " + err.Error())
		return message.Reply("Stallman went bork?")
	}
	return message.Reply(image)
}

func sendImage(message *OctaafMessage) error {
	var images []string
	var err error
	more := message.Command() == "more"
	message.Span.SetTag("more", more)
	if !more {
		if len(message.CommandArguments()) == 0 {
			return message.Reply(fmt.Sprintf("What am I to do, @%v? 🤔🤔🤔🤔", message.From.UserName))
		}

		fetchSpan := message.Span.Tracer().StartSpan(
			"fetch query from google",
			opentracing.ChildOf(message.Span.Context()),
		)

		images, err = scrapers.GetImages(message.CommandArguments(), message.Command() == "img_sfw")

		fetchSpan.Finish()

		if err != nil {
			fetchSpan.SetTag("error", err)
			return message.Reply("Something went wrong!")
		}

		cache.Store(message.Chat.ID, "images", images)
	} else {
		if err := cache.Fetch(message.Chat.ID, "images", &images); err != nil {
			return message.Reply("I can't fetch them for you right now.")
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

	for attempt, url := range images {

		imgSpan := message.Span.Tracer().StartSpan(
			fmt.Sprintf("Download attempt %v", attempt),
			opentracing.ChildOf(message.Span.Context()),
		)

		imgSpan.SetTag("url", url)
		imgSpan.SetTag("attempt", attempt)

		res, err := client.Get(url)

		imgSpan.Finish()

		if err != nil {
			imgSpan.SetTag("error", err)
			continue
		}

		defer res.Body.Close()

		var img []byte

		if message.Command() == "img_censored" {
			log.Debug("Censoring image")
			src, _, err := image.Decode(res.Body)

			if err != nil {
				imgSpan.SetTag("error", true)
				imgSpan.SetBaggageItem("error", err.Error())
				continue
			}

			imgBlurred := imaging.Blur(src, 15)
			buf := new(bytes.Buffer)
			err = png.Encode(buf, imgBlurred.SubImage(imgBlurred.Bounds()))

			if err != nil {
				imgSpan.SetTag("error", true)
				imgSpan.SetBaggageItem("error", err.Error())
				continue
			}

			img = buf.Bytes()
		} else {
			img, err = ioutil.ReadAll(res.Body)

			if err != nil {
				imgSpan.SetTag("error", true)
				imgSpan.SetBaggageItem("error", err.Error())
				log.Errorf("Unable to load image %v; error: %v", url, err)
				continue
			}
		}

		err = message.Reply(img)

		if err == nil {
			return nil
		}
		imgSpan.SetTag("error", err)
	}

	return message.Reply("I did not find images for the query: `" + message.CommandArguments() + "`")
}

func xkcd(message *OctaafMessage) error {
	image, err := scrapers.GetXKCD()

	if err != nil {
		message.LogError("Failed to parse XKCD image: " + err.Error())
		return message.Reply("Failed to parse XKCD image")
	}

	return message.Reply(image)
}

func doubt(message *OctaafMessage) error {
	return message.Reply(assets.Doubt)
}

func quote(message *OctaafMessage) error {
	quoteSpan := message.Span.Tracer().StartSpan(
		"quote",
		opentracing.ChildOf(message.Span.Context()),
	)
	// Fetch a random quote
	if message.ReplyToMessage == nil {
		quoteSpan.SetOperationName("Loading quote...")
		quote := models.Quote{}

		err := quote.Search(DB, message.Chat.ID, message.CommandArguments())

		quoteSpan.Finish()

		if err != nil {
			log.Errorf("Quote fetch error: %v", err)
			quoteSpan.SetTag("error", true)
			quoteSpan.SetBaggageItem("error", err.Error())
			message.LogError("Quote retrieval failure: " + err.Error())
			return message.Reply("No quote found boi")
		}

		userSpan := message.Span.Tracer().StartSpan(
			"Loading username...",
			opentracing.ChildOf(message.Span.Context()),
		)

		user, userErr := getUser(quote.UserID, message.Chat.ID)

		userSpan.Finish()

		if userErr != nil {
			log.Errorf("Unable to find the user for id '%v' : %v", quote.UserID, userErr)
			userSpan.SetTag("error", true)
			userSpan.SetBaggageItem("error", userErr.Error())
			message.LogError(fmt.Sprintf("Unable to find the user for id '%v' : %v", quote.UserID, userErr.Error()))
			return message.Reply(quote.Quote)
		}

		if message.Command() == "presidential_quote" {
			msg := fmt.Sprintf(`"%v"`, quote.Quote)
			msg += fmt.Sprintf("\n    ~@%v", user.User.String())
			img, err := trump.Order(assets.Trump, &settings.Trump, msg)

			if err != nil {
				message.LogError("Presidential quote failure: " + err.Error())
				return message.Reply("Unable to call send a presidential quote.")
			}
			return message.Reply(img)
		}

		msg := fmt.Sprintf("\"%v\"", markdown.Quote(quote.Quote))
		msg += fmt.Sprintf(" \n    ~@%v", markdown.Escape(user.User.String()))
		return message.Reply(msg)
	}

	quoteSpan.SetOperationName("Saving quote...")

	// Unable to store this quote
	if message.ReplyToMessage.Text == "" {
		quoteSpan.SetTag("error", true)
		quoteSpan.SetBaggageItem("error", "No quote found")
		quoteSpan.Finish()
		return message.Reply("No text found in the comment. Not saving the quote!")
	}

	err := DB.Save(&models.Quote{
		Quote:  message.ReplyToMessage.Text,
		UserID: message.ReplyToMessage.From.ID,
		ChatID: message.Chat.ID})

	quoteSpan.Finish()

	if err != nil {
		log.Errorf("Unable to save quote '%v', error: %v", message.ReplyToMessage.Text, err)
		quoteSpan.SetTag("error", true)
		quoteSpan.SetBaggageItem("error", err.Error())
		return message.Reply("Unable to save the quote...")
	}

	return message.Reply("Quote successfully saved!")
}

func nextLaunch(message *OctaafMessage) error {
	fetchSpan := message.Span.Tracer().StartSpan(
		"Fetching nextlaunch data...",
		opentracing.ChildOf(message.Span.Context()),
	)
	res, err := http.Get("https://launchlibrary.net/1.3/launch?next=5&mode=verbose")

	fetchSpan.Finish()

	if err != nil {
		fetchSpan.SetTag("error", true)
		fetchSpan.SetBaggageItem("error", err.Error())
		return message.Reply("Unable to fetch launch data")
	}

	defer res.Body.Close()

	launchJSON, err := ioutil.ReadAll(res.Body)

	if err != nil {
		message.LogError("Launch failure: " + err.Error())
		return message.Reply("Unable to fetch launch data")
	}

	launches := gjson.GetBytes(launchJSON, "launches").Array()

	var msg = "*Next 5 launches:*"

	layout := "January 2, 2006 15:04:05 MST"

	for index, launch := range launches {
		whenStr := launch.Get("net").String()
		when, err := time.Parse(layout, whenStr)

		msg += fmt.Sprintf("\n*%v*: %v", index+1, markdown.Escape(launch.Get("name").String()))

		if err != nil {
			msg += fmt.Sprintf("\n	  %v", markdown.Cursive(whenStr))
		} else {
			msg += fmt.Sprintf("\n	  %v", markdown.Cursive(humanize.Time(when)))
		}

		vods := launch.Get("vidURLs").Array()

		if len(vods) > 0 {
			msg += "\n    " + markdown.Escape(vods[0].String())
		}
	}

	return message.Reply(msg)
}

func gitlabIssues(message *OctaafMessage) error {
	fetchSpan := message.Span.Tracer().StartSpan(
		"Fetching gitlab issues...",
		opentracing.ChildOf(message.Span.Context()),
	)
	res, err := http.Get("https://gitlab.com/api/v4/projects/bartwillems%2Foctaaf/issues?state=opened")

	fetchSpan.Finish()

	if err != nil {
		fetchSpan.SetTag("error", true)
		fetchSpan.SetBaggageItem("error", err.Error())
		return message.Reply("Unable to fetch gitlab issues")
	}

	defer res.Body.Close()

	issuesJSON, err := ioutil.ReadAll(res.Body)

	if err != nil {
		message.LogError("Issues failure: " + err.Error())
		return message.Reply("Unable to fetch gitlab issues")
	}

	issues := gjson.ParseBytes(issuesJSON).Array()

	var msg = "*Open Gitlab Issues:*"

	layout := "2006-01-02T15:04:05.000Z"

	for _, issue := range issues {
		idStr := issue.Get("iid").String()
		urlStr := issue.Get("web_url").String()
		titleStr := issue.Get("title").String()
		authorStr := issue.Get("author.name").String()
		whenStr := issue.Get("created_at").String()
		upvoteStr := issue.Get("upvotes").String()
		downvoteStr := issue.Get("downvotes").String()
		noteStr := issue.Get("user_notes_count").String()

		msg += fmt.Sprintf("\n*#%v*: [%v](%v)\n   👍%v  👎%v  💬%v", idStr, titleStr, urlStr, upvoteStr, downvoteStr, noteStr)
		msg += fmt.Sprintf("\n	  _Reported by %v_", authorStr)

		when, err := time.Parse(layout, whenStr)
		if err != nil {
			msg += fmt.Sprintf("\n	  _Opened %v_", whenStr)
		} else {
			msg += fmt.Sprintf("\n	  _Opened %v_", markdown.Cursive(humanize.Time(when)))
		}
	}

	return message.Reply(msg)
}

func kaliRank(message *OctaafMessage) error {
	if message.Chat.ID != settings.Telegram.KaliID {
		return message.Reply("You are not allowed!")
	}

	kaliRank := []models.MessageCount{}
	err := DB.Order("diff DESC").Limit(5).All(&kaliRank)

	if err != nil {
		message.LogError("Kalirank failure: " + err.Error())
		return message.Reply("Unable to fetch the kali rankings")
	}

	var msg = "*Kali rankings:*"
	for index, rank := range kaliRank {
		msg += fmt.Sprintf("\n`#%v:` *%v messages*   _~%v_", index+1, rank.Diff, rank.CreatedAt.Format("Monday, 2 January 2006"))
	}

	return message.Reply(msg)
}

func care(message *OctaafMessage) error {
	msg := `¯\_(ツ)_/¯`

	// Telegram is broken, so it will always print the wrong amount of backticks
	// This is why this specific command gets parsed as plaintext
	message.IsMarkdown = false

	log.Debug(msg)
	log.Debug(`TEST DINK \ TEST`)

	reply := message.ReplyToMessage
	if reply == nil {
		return message.Reply(msg)
	}

	return message.ReplyTo(msg, reply.MessageID)
}

func pollentiek(message *OctaafMessage) error {
	orientations := map[string][]string{
		"corrupte sos": []string{
			"Liever poen dan groen! 🤑🤑",
			"Zwijg bruine rakker!! 🤚🤚",
			"Wij staken voor uw toekomst 😴😴🍻🍻🍻",
			"Sommige mensen denken dat ze kost wat kost mogen gaan werken 🤓🤓🤓",
		},
		"karakterloze tsjeef": []string{
			"Eat, sleep, tsjeef, repeat 💅💅💅",
			"Is hier nog ergens een chassidische jood beschikbaar om op te komen voor mij? Aub ik smeek u Bartje maakt mij kapot.. 🕎🕎",
			"🆘🆘🆘 't Is al de schuld van de sossen! 🆘🆘🆘",
			"Ik heb geen probleem met moslims in de straat, maar ...💁💁💁",
		},
		"racistische marginale zot": []string{
			"🆘🆘🆘 't Is al de schuld van de sossen! 🆘🆘🆘",
			"Komt door al die vluchtelingen 🏃🏃🏃🏃",
			"Dit is fake nieuws. U kan die posts gewoon op internet vinden. Of zelf maken.\nIemand heeft mijn profielfoto en voornaam gestolen en post zo'n uitspraken in mijn naam.\nMaar die zijn niet van mij. 🤷‍♂️🤷‍♂️🤷‍♂️",
			"😤😤😤Het wordt hoog tijd dat de mensch terug zijn schild en zijn vriend draagtdt!!😤😤😤",
			"Moest Vlaams Belang meer zetels hebben zou dit niet gebeuren punt 🛋🛋🛋🛋",
			"😈Obama😈 and 😈Hillary😈 both smell like 🔥sulfur🔥.",
			"Goddamn liberals 😤😤😤",
			"Beter dood dan rood!🔴☠️🔴☠️🔴☠️",
			"Linkse ratten!! Rolt uw matten!!🐀🐀🐀",
			"Het is weer nen makaak ze 🙉🙉🙉",
			"'t Zijn altijd dezelfden!! 😒😒😒😒",
		},
		"gierige lafaard met geld": []string{
			"🆘🆘🆘 't Is al de schuld van de sossen! 🆘🆘🆘",
			"🇩🇪🇩🇪🇩🇪WIR SCHAFFEN DAS🇩🇪🇩🇪🇩🇪",
			"WIR HABEN DAS NICHT GEWURST🚿🚿🚿",
			"Gewoon doen, watermeloen 🍉🍉🍉🤤🤤",
			"Ge zijt ne flipflop! U en uw partij!🤔🤔🤔🤔",
			"🤤🤤🤤Here is how Bernie can still win..🤤🤤🤤",
		},
	}

	keys := reflect.ValueOf(orientations).MapKeys()

	rand.Seed(time.Now().UnixNano())
	orientation := keys[rand.Intn(len(keys))].String()

	msg := fmt.Sprintf("You are a fullblooded %v.\n", markdown.Bold(orientation))

	rand.Seed(time.Now().UnixNano())
	randomSayIndex := rand.Intn(len(orientations[orientation]))
	saying := orientations[orientation][randomSayIndex]

	msg += fmt.Sprintf("Don't forget to remind everyone around you by proclaiming at least once a day:\n\n%s", markdown.Bold(saying))

	return message.Reply(msg)
}

func presidentialOrder(message *OctaafMessage) error {
	if message.CommandArguments() == "" {
		return message.Reply("Please provide a presidential order.")
	}

	span := message.Span.Tracer().StartSpan(
		"Writing presidential order",
		opentracing.ChildOf(message.Span.Context()),
	)

	img, err := trump.Order(assets.Trump, &settings.Trump, message.CommandArguments())

	span.Finish()

	if err != nil {
		message.LogError("Presidential order failure: " + err.Error())
		return message.Reply("Unable to order dink")
	}

	return message.Reply(img)
}

// msgQuote stores a message as a certain type of quote, eg imgquote, vodquot, ...
func msgQuote(message *OctaafMessage) error {
	quoteType := message.Command()

	if quoteType != models.AudioQuote && quoteType != models.ImgQuote && quoteType != models.VodQuote {
		return message.Reply("Invalid quote type")
	}

	if message.ReplyToMessage == nil {
		quote := models.MsgQuote{}
		err := quote.Search(DB, message.Chat.ID, quoteType)

		if err != nil {
			message.LogError("Quote retrieval failure: " + err.Error())
			return message.Reply("Unable to fetch a quote.")
		}

		return message.Reply(quote)
	}

	if quoteType == models.ImgQuote {
		if message.ReplyToMessage.Photo == nil {
			return message.Reply("No image found in this message.")
		}
	}

	if quoteType == models.AudioQuote {
		if message.ReplyToMessage.Audio == nil && message.ReplyToMessage.Voice == nil {
			return message.Reply("No audio found in this message.")
		}
	}

	if quoteType == models.VodQuote {
		if message.ReplyToMessage.Video == nil && message.ReplyToMessage.VideoNote == nil && message.ReplyToMessage.Animation == nil {
			return message.Reply("No video found in this message.")
		}
	}

	quote := models.MsgQuote{
		MessageID: message.ReplyToMessage.MessageID,
		ChatID:    message.Chat.ID,
		UserID:    message.ReplyToMessage.From.ID,
		Type:      quoteType,
	}

	err := DB.Save(&quote)

	if err != nil {
		message.LogError("Quote save failure: " + err.Error())
		return message.Reply("Unable to store the quote.")
	}

	return message.Reply("Successfully saved the quote!")
}
