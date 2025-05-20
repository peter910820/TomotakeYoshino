package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// get slash command options
func GetOptions(i *discordgo.InteractionCreate, name string) (string, error) {
	for _, v := range i.ApplicationCommandData().Options {
		if v.Name == name {
			value, ok := v.Value.(string) // type assertion
			if !ok {
				return "", errors.New("value translate fail")
			} else {
				return value, nil
			}
		}
	}
	return "", errors.New("option not found")
}

// use interactionCreate to get user id
func GetUserID(i *discordgo.InteractionCreate) (string, error) {
	if i.Member != nil {
		return i.Member.User.ID, nil
	} else if i.User != nil {
		return i.User.ID, nil
	} else {
		logrus.Error("cannot found user id")
		return "", errors.New("cannot found user id")
	}
}

// use user id and session to get user name
func GetUserName(s *discordgo.Session, userID string) (string, error) {
	user, err := s.User(userID)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s#%s", user.Username, user.Discriminator), nil
}

// use json data to post request to specific url
func JsonRequest(url string, method string, data []byte) (*http.Response, []byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		logrus.Error(err)
		return nil, []byte{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error(err)
		return nil, []byte{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		return nil, []byte{}, err
	}
	return resp, body, nil
}

// handle slash command error return
func SlashCommandError(s *discordgo.Session, i *discordgo.InteractionCreate, err string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: err,
		},
	})
}
