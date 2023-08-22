package botcommands

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

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
	Scissors = "makas"
	Rock     = "taş"
	Paper    = "kağıt"
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
	userChoose := options[0]
	if userChoose != Scissors && userChoose != Rock && userChoose != Paper {
		s.ChannelMessageSend(m.ChannelID, "Düzgün oyna lan")
		return errors.New("Given string is not a prs")
	}
	stringToSeed := func(input string) int64 {
		seed := int64(0)
		for _, char := range input {
			seed += int64(char)
		}
		return seed
	}
	source := rand.NewSource(stringToSeed(m.Author.ID) * int64(time.Now().Nanosecond()))
	random := rand.New(source)
	botChooseNum := random.Intn(3)
	botChoose := ""
	switch botChooseNum {
	case 0:
		botChoose = "kağıt"

	case 1:
		botChoose = "taş"

	case 2:
		botChoose = "makas"
	}
	getWinner := func(c1, c2 string) int {
		if c1 == Rock {
			switch c2 {
			case Rock:
				return 0
			case Scissors:
				return 1
			case Paper:
				return 2
			}
		}
		if c1 == Scissors {
			switch c2 {
			case Rock:
				return 2
			case Scissors:
				return 0
			case Paper:
				return 1
			}
		}
		if c1 == Paper {
			switch c2 {
			case Rock:
				return 1
			case Scissors:
				return 2
			case Paper:
				return 0
			}
		}
		return 0
	}
	winner := getWinner(userChoose, botChoose)
	sonuc := "Berabere"
	switch winner {
	case 1:
		sonuc = "Afferin Adam Oluyon"
	case 2:
		sonuc = "YENDİM PİÇ, ŞİMDİ SİKTİR GİT"
	default:
		sonuc = "Sakinn, kimse kimseye hiçbi şey yapamadı"
	}
	stringToShow := fmt.Sprintf("Ben Sectim: %v\n %v Secti: %v\n Sonuc:%v", botChoose, m.Author.Username, userChoose, sonuc)
	s.ChannelMessageSend(m.ChannelID, stringToShow)

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
