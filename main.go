package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"internal/botcommands"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("DISCORD_TOKEN")
	dg, err := discordgo.New("Bot " + token)
	dg.SyncEvents = false

	dg.AddHandler(botcommands.Ready)
	dg.AddHandler(botcommands.VoiceServerUpdate)
	dg.AddHandler(botcommands.MessageCreate)
	if err != nil {
		log.Fatal(err)
	}
	dg.Identify.Intents = discordgo.IntentsAll
	dg.Open()
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
	return

}
