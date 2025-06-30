package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

type BoostData struct {
	GuildID     string    `json:"guild_id"`
	UserID      string    `json:"user_id"`
	Username    string    `json:"username"`
	BoostTime   time.Time `json:"boost_time"`
	IsActive    bool      `json:"is_active"`
}

var (
	boosts = make(map[string][]BoostData)
	commandPrefix = "!"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN environment variable is required")
	}

	if prefix := os.Getenv("COMMAND_PREFIX"); prefix != "" {
		commandPrefix = prefix
	}

	// Create Discord session
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}

	// Register event handlers
	dg.AddHandler(ready)
	dg.AddHandler(messageCreate)
	dg.AddHandler(guildMemberUpdate)

	// Set intents
	dg.Identify.Intents = discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	// Open connection
	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening connection:", err)
	}

	fmt.Println("Discord Boost Bot is now running. Press CTRL+C to exit.")
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
	case "boostinfo":
		handleBoostInfo(s, m, args)
	case "ping":
		handlePing(s, m)
	}
}

func guildMemberUpdate(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
	// Check if premium since (boost time) changed
	before := m.BeforeUpdate
	if before == nil {
		return
	}

	// If user started boosting
	if m.Member.PremiumSince != nil && (before.PremiumSince == nil || before.PremiumSince.Before(*m.Member.PremiumSince)) {
		handleNewBoost(s, m.Member, m.GuildID)
	}

	// If user stopped boosting
	if m.Member.PremiumSince == nil && before.PremiumSince != nil {
		handleBoostRemoved(s, m.Member, m.GuildID)
	}
}

func handleNewBoost(s *discordgo.Session, member *discordgo.Member, guildID string) {
	// Add boost to our tracking
	boost := BoostData{
		GuildID:   guildID,
		UserID:    member.User.ID,
		Username:  member.User.Username,
		BoostTime: *member.PremiumSince,
		IsActive:  true,
	}

	if _, exists := boosts[guildID]; !exists {
		boosts[guildID] = []BoostData{}
	}
	boosts[guildID] = append(boosts[guildID], boost)

	// Send boost notification
	guild, err := s.Guild(guildID)
	if err != nil {
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🚀 Server Boosted!",
		Description: fmt.Sprintf("Thank you **%s** for boosting the server!", member.User.Username),
		Color:       0xFF69B4, // Pink color
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Current Boost Level",
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

	// Try to send in boost channel, fallback to system channel
	channelID := os.Getenv("BOOST_CHANNEL_ID")
	if channelID == "" {
		channelID = guild.SystemChannelID
	}

	if channelID != "" {
		s.ChannelMessageSendEmbed(channelID, embed)
	}
}

func handleBoostRemoved(s *discordgo.Session, member *discordgo.Member, guildID string) {
	// Mark boost as inactive
	if guildBoosts, exists := boosts[guildID]; exists {
		for i := range guildBoosts {
			if guildBoosts[i].UserID == member.User.ID && guildBoosts[i].IsActive {
				boosts[guildID][i].IsActive = false
				break
			}
		}
	}

	// Send notification
	channelID := os.Getenv("BOOST_CHANNEL_ID")
	guild, _ := s.Guild(guildID)
	if guild != nil && guild.SystemChannelID != "" && channelID == "" {
		channelID = guild.SystemChannelID
	}

	if channelID != "" {
		embed := &discordgo.MessageEmbed{
			Title:       "💔 Boost Removed",
			Description: fmt.Sprintf("**%s** is no longer boosting the server.", member.User.Username),
			Color:       0x808080, // Gray color
			Timestamp:   time.Now().Format(time.RFC3339),
		}
		s.ChannelMessageSendEmbed(channelID, embed)
	}
}

func handleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "🚀 Discord Boost Bot Commands",
		Description: "Here are all available commands:",
		Color:       0x7289DA,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   commandPrefix + "boosts",
				Value:  "Show current server boost information",
				Inline: false,
			},
			{
				Name:   commandPrefix + "boostcount",
				Value:  "Display total number of boosts",
				Inline: false,
			},
			{
				Name:   commandPrefix + "topboosters",
				Value:  "Show top server boosters",
				Inline: false,
			},
			{
				Name:   commandPrefix + "boostinfo @user",
				Value:  "Get boost information for a specific user",
				Inline: false,
			},
			{
				Name:   commandPrefix + "ping",
				Value:  "Check bot latency",
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Discord Boost Bot",
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func handleBoosts(s *discordgo.Session, m *discordgo.MessageCreate) {
	guild, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error fetching guild information.")
		return
	}

	var nextLevelBoosts int
	var nextLevelPerks string

	switch guild.PremiumTier {
	case 0:
		nextLevelBoosts = 2
		nextLevelPerks = "Animated Server Icon, 128 Kbps Audio, Custom Server Invite Background"
	case 1:
		nextLevelBoosts = 7
		nextLevelPerks = "Server Banner, 256 Kbps Audio, Animated Server Icon"
	case 2:
		nextLevelBoosts = 14
		nextLevelPerks = "Vanity URL, 384 Kbps Audio, Server Banner"
	default:
		nextLevelBoosts = 0
		nextLevelPerks = "Maximum level reached!"
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("🚀 %s Boost Status", guild.Name),
		Color:       0xFF69B4,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: guild.IconURL("256"),
		},
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
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if nextLevelBoosts > 0 {
		remaining := nextLevelBoosts - guild.PremiumSubscriptionCount
		if remaining > 0 {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Next Level",
				Value:  fmt.Sprintf("%d more boosts needed", remaining),
				Inline: true,
			})
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Next Level Perks",
				Value:  nextLevelPerks,
				Inline: false,
			})
		}
	} else {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Status",
			Value:  "Maximum boost level reached! 🎉",
			Inline: false,
		})
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func handleBoostCount(s *discordgo.Session, m *discordgo.MessageCreate) {
	guild, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error fetching guild information.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("**%s** currently has **%d** boosts! 🚀", guild.Name, guild.PremiumSubscriptionCount))
}

func handleTopBoosters(s *discordgo.Session, m *discordgo.MessageCreate) {
	members, err := s.GuildMembers(m.GuildID, "", 1000)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error fetching guild members.")
		return
	}

	var boosters []string
	for _, member := range members {
		if member.PremiumSince != nil {
			boosters = append(boosters, fmt.Sprintf("• **%s** - <t:%d:R>", member.User.Username, member.PremiumSince.Unix()))
		}
	}

	if len(boosters) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No current boosters found.")
		return
	}

	// Limit to top 10
	if len(boosters) > 10 {
		boosters = boosters[:10]
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🏆 Top Server Boosters",
		Description: strings.Join(boosters, "\n"),
		Color:       0xFFD700, // Gold color
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Showing %d boosters", len(boosters)),
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func handleBoostInfo(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Please mention a user: `"+commandPrefix+"boostinfo @user`")
		return
	}

	// Extract user ID from mention
	userID := strings.Trim(args[1], "<@!>")

	member, err := s.GuildMember(m.GuildID, userID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "User not found in this server.")
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("Boost Info - %s", member.User.Username),
		Color: 0x7289DA,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: member.User.AvatarURL("256"),
		},
	}

	if member.PremiumSince != nil {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Boost Status",
			Value:  "✅ Currently Boosting",
			Inline: true,
		})
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Boosting Since",
			Value:  fmt.Sprintf("<t:%d:F>", member.PremiumSince.Unix()),
			Inline: true,
		})
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Boost Duration",
			Value:  fmt.Sprintf("<t:%d:R>", member.PremiumSince.Unix()),
			Inline: true,
		})
		embed.Color = 0xFF69B4
	} else {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Boost Status",
			Value:  "❌ Not Currently Boosting",
			Inline: false,
		})
		embed.Color = 0x808080
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
			{
				Name:   "Message Latency",
				Value:  fmt.Sprintf("%v", latency),
				Inline: true,
			},
			{
				Name:   "WebSocket Latency", 
				Value:  fmt.Sprintf("%v", s.HeartbeatLatency()),
				Inline: true,
			},
		},
		Color: 0x00FF00,
	}
	
	s.ChannelMessageEditEmbed(m.ChannelID, msg.ID, embed)
}
