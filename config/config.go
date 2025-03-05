package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func NewConfig(path string) (*Config, error) {

	d := &Config{}
	err := d.initialize(path)
	if err != nil {
		return nil, err
	}

	err = d.checkConfig()
	if err != nil {
		return nil, err
	}

	err = d.parseChats()
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (conf *Config) initialize(path string) error {

	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	stats, err := os.Stat(abs)

	if os.IsNotExist(err) {
		return errors.New("given path does not exist")
	}

	if err != nil {
		return err
	}

	if stats.IsDir() {
		return errors.New("given path is a directiory")
	}

	conf.Path = abs

	bytes, err := ioutil.ReadFile(conf.Path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &conf.Config)
	if err != nil {
		return err
	}

	return nil
}

func (conf *Config) checkConfig() error {

	if conf.Config.APIKey == "" {
		return errors.New("APIKey cannot be null")
	}

	parent := filepath.Dir(conf.Path)

	chats := filepath.Join(parent, "chats")
	prompts := filepath.Join(parent, "prompts")

	stats, err := os.Stat(chats)
	if os.IsNotExist(err) {
		return errors.New("chats path does not exist")
	}
	if err != nil {
		return err
	}

	prompts_stats, err := os.Stat(prompts)
	if os.IsNotExist(err) {
		return errors.New("chats path does not exist")
	}
	if err != nil {
		return err
	}

	if !stats.IsDir() {
		return errors.New("chats path not is a directory")
	}

	if !prompts_stats.IsDir() {
		return errors.New("prompts path not is a directory")
	}

	conf.ChatsPath = chats
	conf.PromptsPath = prompts

	files, err := ioutil.ReadDir(conf.ChatsPath)
	if err != nil {
		return err
	}

	configs := make(map[string]string)
	for _, file := range files {
		fileName := file.Name()
		if !strings.HasSuffix(fileName, ".config.json") {
			continue
		}
		fileName = strings.TrimSuffix(fileName, ".config.json")
		configs[fileName] = file.Name()
	}

	for _, chat := range conf.Config.AvailableChats {
		if _, ok := configs[chat]; !ok {
			return errors.New(
				fmt.Sprintf("%v config not found in %v", chat, conf.ChatsPath))
		}
	}

	return nil
}

func (conf *Config) parseChats() error {

	conf.Config.Chats = make(ChatStringMap)
	conf.Config.IDChats = make(ChatIntMap)
	for _, chatName := range conf.Config.AvailableChats {
		filePath := filepath.Join(conf.ChatsPath,
			fmt.Sprintf("%v.config.json", chatName))

		c := &Chat{}

		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}
		err = json.Unmarshal(bytes, c)
		if err != nil {
			return err
		}

		c.ChatName = chatName

		for _, aiReply := range c.AiReplies {
			promptsPath := filepath.Join(conf.PromptsPath,
				fmt.Sprintf("%v.txt", aiReply.Prompt))

			bytes, err := ioutil.ReadFile(promptsPath)
			if err != nil {
				return err
			}
			aiReply.PromptText = string(bytes)
		}

		conf.Config.Chats[chatName] = c
		conf.Config.IDChats[c.ChatID] = c
	}

	return nil
}
