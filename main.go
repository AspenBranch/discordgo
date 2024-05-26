package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create a new session using the DISCORD_TOKEN environment variable from Railway
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		fmt.Printf("Error while starting bot: %s", err)
		return
	}

	// Add the message handler
	dg.AddHandler(messageCreate)

	// Add the Guild Messages intent
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Connect to the gateway
	err = dg.Open()
	if err != nil {
		fmt.Printf("Error while connecting to gateway: %s", err)
		return
	}

	// Wait until Ctrl+C or another signal is received
	fmt.Println("The bot is now running. Press Ctrl+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Close the Discord session
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Don't proceed if the message author is a bot
	if m.Author.Bot {
		return
	}

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong")
		return
	}

	if m.Content == "hello" {
		s.ChannelMessageSend(m.ChannelID, "Hello")
		return
	}
	if m.Content == "Hello" {
		s.ChannelMessageSend(m.ChannelID, "Hello")
		return
	}
	if m.Content == "!echo" {
		content := m.ContentWithMentionsReplaced()
		s.ChannelMessageSend(m.ChannelID, content)
	}

	if m.Content == "!time" {
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		s.ChannelMessageSend(m.ChannelID, "Current server time is: "+currentTime)
	}
	if m.Content == "" {
		s.ChannelMessageSend(m.ChannelID, "")
		return
	}
	if m.Content == "" {
		s.ChannelMessageSend(m.ChannelID, "")
		return
	}

	if m.Content == "" {
		s.ChannelMessageSend(m.ChannelID, "")
		return
	}
	if strings.HasPrefix(m.Content, "!weather") {
		city := strings.TrimPrefix(m.Content, "!weather ")
		weather, err := getWeather(city)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Could not fetch weather data.")
			return
		}
		s.ChannelMessageSend(m.ChannelID, weather)
	}

	if strings.HasPrefix(m.Content, "!joke") {
		joke, err := getJoke()
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Could not fetch a joke.")
			return
		}
		s.ChannelMessageSend(m.ChannelID, joke)
	}

	if strings.HasPrefix(m.Content, "!quote") {
		quote, err := getQuote()
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Could not fetch a quote.")
			return
		}
		s.ChannelMessageSend(m.ChannelID, quote)
	}

	if strings.HasPrefix(m.Content, "!roll") {
		result := rand.Intn(6) + 1
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You rolled a %d", result))
	}

	if strings.HasPrefix(m.Content, "!avatar") {
		s.ChannelMessageSend(m.ChannelID, m.Author.AvatarURL(""))
	}

	if strings.HasPrefix(m.Content, "!ban") && len(m.Mentions) > 0 {
		user := m.Mentions[0]
		err := s.GuildBanCreateWithReason(m.GuildID, user.ID, "Banned by bot", 0)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Failed to ban user.")
		} else {
			s.ChannelMessageSend(m.ChannelID, "User banned successfully.")
		}
	}

	if strings.HasPrefix(m.Content, "!kick") && len(m.Mentions) > 0 {
		user := m.Mentions[0]
		err := s.GuildMemberDeleteWithReason(m.GuildID, user.ID, "Kicked by bot")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Failed to kick user.")
		} else {
			s.ChannelMessageSend(m.ChannelID, "User kicked successfully.")
		}
	}

	if strings.HasPrefix(m.Content, "!userinfo") {
		user := m.Author
		info := fmt.Sprintf("User Info:\nUsername: %s\nID: %s\nAvatar: %s", user.Username, user.ID, user.AvatarURL(""))
		s.ChannelMessageSend(m.ChannelID, info)
	}
}
	
func getWeather(city string) (string, error) {
	apiKey := "YourOpenWeatherMapAPIKey"
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	if result["cod"] != float64(200) {
		return "", fmt.Errorf("City not found")
	}

	main := result["main"].(map[string]interface{})
	weather := result["weather"].([]interface{})[0].(map[string]interface{})

	description := weather["description"].(string)
	temp := main["temp"].(float64)

	return fmt.Sprintf("Weather in %s: %s, %.2f°C", city, description, temp), nil
}

func getJoke() (string, error) {
	url := "https://official-joke-api.appspot.com/random_joke"
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var joke map[string]interface{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &joke)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s - %s", joke["setup"].(string), joke["punchline"].(string)), nil
}

func getQuote() (string, error) {
	url := "https://api.quotable.io/random"
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var quote map[string]interface{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &quote)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s - %s", quote["content"].(string), quote["author"].(string)), nil
}
