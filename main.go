package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var commandPrefix = "!"

func main() {
	godotenv.Load()
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN required")
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}

	dg.AddHandler(ready)
	dg.AddHandler(messageCreate)
	dg.AddHandler(guildMemberUpdate)
	dg.Identify.Intents = discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening connection:", err)
	}

	fmt.Println("Discord Boost Bot running! Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	dg.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateGameStatus(0, "Server Boosts | "+commandPrefix+"help")
	log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || !strings.HasPrefix(m.Content, commandPrefix) {
		return
	}

	args := strings.Fields(m.Content)
	if len(args) == 0 {
		return
	}

	command := strings.ToLower(args[0][len(commandPrefix):])

	switch command {
	case "help":
		handleHelp(s, m)
	case "boosts":
		handleBoosts(s, m)
	case "boostcount":
		handleBoostCount(s, m)
	case "topboosters":
		handleTopBoosters(s, m)
	case "ping":
		handlePing(s, m)
	}
}

func guildMemberUpdate(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
	before := m.BeforeUpdate
	if before == nil {
		return
	}

	if m.Member.PremiumSince != nil && (before.PremiumSince == nil || before.PremiumSince.Before(*m.Member.PremiumSince)) {
		handleNewBoost(s, m.Member, m.GuildID)
	}

	if m.Member.PremiumSince == nil && before.PremiumSince != nil {
		handleBoostRemoved(s, m.Member, m.GuildID)
	}
}

func handleNewBoost(s *discordgo.Session, member *discordgo.Member, guildID string) {
	guild, err := s.Guild(guildID)
	if err != nil {
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🚀 Server Boosted!",
		Description: fmt.Sprintf("Thank you **%s** for boosting the server!", member.User.Username),
		Color:       0xFF69B4,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Current Level",
				Value:  fmt.Sprintf("Level %d", guild.PremiumTier),
				Inline: true,
			},
			{
				Name:   "Total Boosts",
				Value:  fmt.Sprintf("%d", guild.PremiumSubscriptionCount),
				Inline: true,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: member.User.AvatarURL("256"),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	channelID := os.Getenv("BOOST_CHANNEL_ID")
	if channelID == "" {
		channelID = guild.SystemChannelID
	}

	if channelID != "" {
		s.ChannelMessageSendEmbed(channelID, embed)
	}
}

func handleBoostRemoved(s *discordgo.Session, member *discordgo.Member, guildID string) {
	channelID := os.Getenv("BOOST_CHANNEL_ID")
	guild, _ := s.Guild(guildID)
	if guild != nil && guild.SystemChannelID != "" && channelID == "" {
		channelID = guild.SystemChannelID
	}

	if channelID != "" {
		embed := &discordgo.MessageEmbed{
			Title:       "💔 Boost Removed",
			Description: fmt.Sprintf("**%s** is no longer boosting the server.", member.User.Username),
			Color:       0x808080,
			Timestamp:   time.Now().Format(time.RFC3339),
		}
		s.ChannelMessageSendEmbed(channelID, embed)
	}
}

func handleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "🚀 Discord Boost Bot Commands",
		Description: "Available commands:",
		Color:       0x7289DA,
		Fields: []*discordgo.MessageEmbedField{
			{Name: commandPrefix + "boosts", Value: "Show server boost info", Inline: false},
			{Name: commandPrefix + "boostcount", Value: "Display total boosts", Inline: false},
			{Name: commandPrefix + "topboosters", Value: "Show top boosters", Inline: false},
			{Name: commandPrefix + "ping", Value: "Check bot latency", Inline: false},
		},
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func handleBoosts(s *discordgo.Session, m *discordgo.MessageCreate) {
	guild, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error fetching guild info.")
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:   fmt.Sprintf("🚀 %s Boost Status", guild.Name),
		Color:   0xFF69B4,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Current Level", Value: fmt.Sprintf("Level %d", guild.PremiumTier), Inline: true},
			{Name: "Total Boosts", Value: fmt.Sprintf("%d", guild.PremiumSubscriptionCount), Inline: true},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{URL: guild.IconURL("256")},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func handleBoostCount(s *discordgo.Session, m *discordgo.MessageCreate) {
	guild, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error fetching guild info.")
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("**%s** has **%d** boosts! 🚀", guild.Name, guild.PremiumSubscriptionCount))
}

func handleTopBoosters(s *discordgo.Session, m *discordgo.MessageCreate) {
	members, err := s.GuildMembers(m.GuildID, "", 1000)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error fetching members.")
		return
	}

	var boosters []string
	for _, member := range members {
		if member.PremiumSince != nil {
			boosters = append(boosters, fmt.Sprintf("• **%s** - <t:%d:R>", member.User.Username, member.PremiumSince.Unix()))
		}
	}

	if len(boosters) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No boosters found.")
		return
	}

	if len(boosters) > 10 {
		boosters = boosters[:10]
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🏆 Top Server Boosters",
		Description: strings.Join(boosters, "\n"),
		Color:       0xFFD700,
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func handlePing(s *discordgo.Session, m *discordgo.MessageCreate) {
	start := time.Now()
	msg, _ := s.ChannelMessageSend(m.ChannelID, "🏓 Pinging...")
	latency := time.Since(start)

	embed := &discordgo.MessageEmbed{
		Title: "🏓 Pong!",
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Latency", Value: fmt.Sprintf("%v", latency), Inline: true},
			{Name: "WebSocket", Value: fmt.Sprintf("%v", s.HeartbeatLatency()), Inline: true},
		},
		Color: 0x00FF00,
	}
	s.ChannelMessageEditEmbed(m.ChannelID, msg.ID, embed)
}