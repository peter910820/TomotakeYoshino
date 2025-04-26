package utils

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

func GetOptions(i *discordgo.InteractionCreate, name string) (string, error) {
	for _, v := range i.ApplicationCommandData().Options {
		if v.Name == name {
			value, ok := v.Value.(string)
			if !ok {
				return "", errors.New("value translate fail")
			} else {
				return value, nil
			}
		}
	}
	return "", errors.New("option not found")
}
