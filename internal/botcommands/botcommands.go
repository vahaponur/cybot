package botcommands

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type User = discordgo.User

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	fields := strings.Fields(m.Content)
	command, ok := getCommands()[fields[0]]
	if !ok {
		return
	}

	command.Callback(s, m, fields[1:]...)

}

type Command struct {
	Prefix   string
	Callback func(s *discordgo.Session, m *discordgo.MessageCreate, options ...string) error
}

func commandPlay(s *discordgo.Session, m *discordgo.MessageCreate, options ...string) error {
	return nil
}
func commandGetUserNames(s *discordgo.Session, m *discordgo.MessageCreate, options ...string) error {
	c, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Fatal(err)
	}

	if c.Type != discordgo.ChannelTypeGuildText {
		log.Default()
	}
	members, err := s.GuildMembers(c.GuildID, "", 1000)
	mem := ""
	for _, me := range members {
		permissions, _ := s.UserChannelPermissions(me.User.ID, c.ID)
		if permissions&discordgo.PermissionReadMessages != 0 {
			mem += fmt.Sprintf("%v\n", me.User.Username)

		}
	}
	s.ChannelMessageSend(m.ChannelID, mem)

	return nil
}

const (
	Scissors = "Makas"
	Rock     = "Taş"
	Paper    = "Kağıt"
)

// Play Paper Rock Scissors
func commandPRS(s *discordgo.Session, m *discordgo.MessageCreate, options ...string) error {
	c, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Fatal(err)
	}

	if c.Type != discordgo.ChannelTypeGuildText {
		log.Default()
	}
	if len(options) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Bir secim yapcan knk")
		return errors.New("Not choosen a prs")
	}
	choose := options[0]
	if choose != Scissors && choose != Rock && choose != Paper {
		s.ChannelMessageSend(m.ChannelID, "Düzgün oyna lan")
		return errors.New("Given string is not a prs")
	}
	return nil
}
func getCommands() map[string]Command {
	return map[string]Command{
		"!play": {
			Prefix:   "!play",
			Callback: commandPlay,
		},
		"!agalar": {
			Prefix:   "!agalar",
			Callback: commandGetUserNames,
		},
		"!tkm": {
			Prefix:   "!tkm",
			Callback: commandPRS,
		},
	}
}
