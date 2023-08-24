package botcommands

import (
	"errors"
	"fmt"
	"log"
	"math/rand"

	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/foxbot/gavalink"
)

type User = discordgo.User

var lavalink *gavalink.Lavalink
var player *gavalink.Player

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
	for _, r := range fields {
		fmt.Println(r)
	}
	go command.Callback(s, m, fields[1:]...)

}

type Command struct {
	Prefix   string
	Callback func(s *discordgo.Session, m *discordgo.MessageCreate, options ...string) error
}

func commandPlay(s *discordgo.Session, m *discordgo.MessageCreate, options ...string) error {
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Println("fail find channel")
		return nil
	}

	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		log.Println("fail find guild")
		return nil
	}

	for _, vs := range g.VoiceStates {
		if vs.UserID == m.Author.ID {
			log.Println("trying to connect to channel")
			err = s.ChannelVoiceJoinManual(c.GuildID, vs.ChannelID, false, false)
			if err != nil {
				log.Println(err)
			} else {
				log.Println("channel voice join succeeded")
			}
		}
	}
	qs := strings.Join(options, "%20")
	query := fmt.Sprintf("ytsearch:%v", qs)
	fmt.Println(query)
	node, err := lavalink.BestNode()
	if err != nil {
		log.Println(err)
	}
	tracks, err := node.LoadTracks(query)
	if err != nil {
		log.Println(err)
	}
	if tracks.Type != gavalink.TrackLoaded {
		log.Println("weird tracks type: ", tracks.Type)
	}
	track := tracks.Tracks[0].Data

	err = player.Play(track)
	if err != nil {
		log.Println(err)
	}
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
	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(len(members))
	mem := ""
	for _, me := range members {

		go printUser(s, c, me, &mem, &wg, &mu)

	}
	wg.Wait()
	mem = mem[:1999]
	s.ChannelMessageSend(m.ChannelID, mem)

	return nil
}
func commandStop(s *discordgo.Session, m *discordgo.MessageCreate, options ...string) error {
	err := player.Stop()
	return err
}
func printUser(s *discordgo.Session, c *discordgo.Channel, me *discordgo.Member, mem *string, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()

	now := time.Now()
	permissions, _ := s.State.UserChannelPermissions(me.User.ID, c.ID)
	passed := time.Until(now)
	if permissions&discordgo.PermissionViewChannel != 0 {
		mu.Lock()
		fmt.Println(me.User.Username)
		*mem += fmt.Sprintf("%v Time Passed:%v\n", me.User.Username, passed)
		mu.Unlock()
	}

}
func printPermUser(s *discordgo.Session, c *discordgo.Channel, me *discordgo.Member, mem *string, wg *sync.WaitGroup) {
	defer wg.Done()

	permissions, err := s.State.UserChannelPermissions(me.User.ID, c.ID)
	if err != nil {
		fmt.Println(err)
	}

	*mem += fmt.Sprintf("%v Permission:%v\n", me.User.Username, permissions)

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
		sonuc = "Yendimm knk ağla bağır"
	default:
		sonuc = "Sakinn, kimse kimseye hiçbi şey yapamadı"
	}
	stringToShow := fmt.Sprintf("Ben Sectim: %v\n %v Secti: %v\n Sonuc:%v", botChoose, m.Author.Username, userChoose, sonuc)
	s.ChannelMessageSend(m.ChannelID, stringToShow)

	return nil
}
func commandKoyluler(s *discordgo.Session, m *discordgo.MessageCreate, options ...string) error {
	c, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Fatal(err)
	}

	if c.Type != discordgo.ChannelTypeGuildText {
		log.Default()
	}
	poem := "Köylüleri niçin öldürmeliyiz?\n" +
		"Çünkü onlar ağır kanlı adamlardır\n" +
		"Değişen bir dünyaya karşı\n" +
		"Kerpiç duvarlar gibi katı\n" +
		"Çakır dikenleri gibi susuz\n" +
		"Kayıtsızca direnerek yaşarlar.\n" +
		"Aptal, kaba ve kurnazdırlar.\n" +
		"İnanarak ve kolayca yalan söylerler.\n" +
		"Paraları olsa da\n" +
		"Yoksul görünmek gibi bir hünerleri vardır.\n" +
		"Her şeyi hafife alır ve herkese söverler.\n" +
		"Yağmuru, rüzgarı ve güneşi\n" +
		"Bir gün olsun ekinleri akıllarına gelmeden\n" +
		"Düşünemezler…\n" +
		"Ve birbirlerinin sınırlarını sürerek\n" +
		"Topraklarını büyütmeye çalışırlar.\n" +
		"\n" +
		"Köylüleri niçin öldürmeliyiz?\n" +
		"Çünkü onlar karılarını döverler\n" +
		"Seslerinin tonu yumuşak değildir\n" +
		// Diğer satırlar burada devam eder...
		"Yarı gecelerde yıldızlara bakarak\n" +
		"Başka dünyaları düşünmek gibi bir tutkuları yoktur.\n" +
		"Gökyüzünü baharda yağmur yağarsa\n" +
		"Ve yaz güneşleri ekinlerini yetirirse severler.\n" +
		"Hayal güçleri kıttır ve hiçbir yeniliğe\n" +
		"-Bu verimi yüksek bir tohum bile olsa-\n" +
		"Sonuçlarını görmeden inanmazlar.\n" +
		"Dünyanın gelişimine bir katkıları yoktur.\n" +
		"Mülk düşkünüdürler amansız derecede\n" +
		"Bir ülkenin geleceği\n" +
		"Küçücük topraklarının ipoteği altındadır.\n" +
		"Ve birer kaya parçası gibi dururlar su geçirmeden\n" +
		"Zamanın derin ırmakları önünde…\n" +
		"\n" +
		"KÖYLÜLERİ, SÖYLEYİN NASIL\n" +
		"NASIL KURTARALIM?"
	s.ChannelMessageSend(m.ChannelID, poem)
	return nil
}
func getCommands() map[string]Command {
	return map[string]Command{
		"c/play": {
			Prefix:   "!play",
			Callback: commandPlay,
		},
		"c/agalar": {
			Prefix:   "!agalar",
			Callback: commandGetUserNames,
		},
		"c/tkm": {
			Prefix:   "!tkm",
			Callback: commandPRS,
		},
		"c/koyluler": {
			Prefix:   "!koyluler",
			Callback: commandKoyluler,
		},
		"c/stop": {
			Prefix:   "!stop",
			Callback: commandStop,
		},
	}
}
