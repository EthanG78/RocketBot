package main

import(
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"io/ioutil"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func main(){

	//Open token.txt for reading
	file, err := os.Open("token.txt")
	if err != nil{
		log.Fatalf("Error opening file: %v\n", err)
	}

	defer file.Close()

	//Reads secret bot token from token.txt
	data, err := ioutil.ReadAll(file)
	if err != nil{
		log.Fatalf("Error reading file: %s\n", err)
	}

	//Remove newline character from token string
	token := string(data)
	token = strings.TrimSuffix(token, "\n")

	//Initialize bot with discordgo package
	discord, err := discordgo.New("Bot " + token)
	if err != nil{
		log.Fatalf("Error initializing discord bot: %v\n", err)
	}

	//Open a websocket connection to discord and start listening
	err = discord.Open()
	if err != nil{
		log.Fatalf("Error creating websocket connection to discord: %v\n", err)
	}

	//Add handlers
	discord.AddHandler(onMessageCreate)

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}

//A function that is called everytime a new message is created 
func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate){

	//The bot will ignore all messages that are created by itself
	if m.Author.ID == s.State.User.ID{
		return
	}

	//If the message is !ping, reply with pong!
	if m.Content == "!ping"{
		s.ChannelMessageSend(m.ChannelID, "pong!")
	}
}