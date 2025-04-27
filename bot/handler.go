package bot

import "github.com/bwmarrin/discordgo"

func Ready(s *discordgo.Session, m *discordgo.Ready) {
	s.UpdateGameStatus(0, "サクラノ詩-櫻の森の上を舞う-")
	RegisterCommand(s)
}

func GuildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	welcomeMessage := "Welcome " + m.User.Username + "!"
	s.ChannelMessageSend(m.GuildID, welcomeMessage)
}
