package botcommands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/foxbot/gavalink"
)

func Ready(s *discordgo.Session, event *discordgo.Ready) {
	idle := 0
	log.Println("discordgo ready!")
	s.UpdateWatchStatus(0, "gavalink")
	s.UpdateListeningStatus("gavalink")
	s.UpdateGameStatus(0, "gavalink")
	s.UpdateStatusComplex(discordgo.UpdateStatusData{Status: "gavalink", IdleSince: &idle})
	lavalink = gavalink.NewLavalink("1", event.User.ID)

	err := lavalink.AddNodes(gavalink.NodeConfig{
		REST:      "http://localhost:2333",
		WebSocket: "ws://localhost:2333",
		Password:  "youshallnotpass",
	})
	fmt.Println("SELAM DUNYALI")
	x, _ := lavalink.BestNode()
	if x == nil {
		fmt.Print("sdfs")
	}
	if err != nil {
		log.Println(err)
	}
}
func VoiceServerUpdate(s *discordgo.Session, event *discordgo.VoiceServerUpdate) {
	log.Println("received VSU")
	vsu := gavalink.VoiceServerUpdate{
		Endpoint: event.Endpoint,
		GuildID:  event.GuildID,
		Token:    event.Token,
	}

	if p, err := lavalink.GetPlayer(event.GuildID); err == nil {
		err = p.Forward(s.State.SessionID, vsu)
		if err != nil {
			log.Println(err)
		}
		return
	}

	node, err := lavalink.BestNode()
	if err != nil {
		log.Println(err)
		return
	}

	handler := new(gavalink.DummyEventHandler)
	player, err = node.CreatePlayer(event.GuildID, s.State.SessionID, vsu, handler)
	if err != nil {
		log.Println(err)
		return
	}
}
