package command

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

type Metadata struct {
	Cmd   string
	Place string
}

type Command interface {
	ToHash() string
	ToProtoHash() string
	Debug() string
	GetMetadata() *Metadata
	GetCommand() string
	GetCaption() string
}

type ForceCommand interface {
	Command
	GetExecutionDate() string
}

type UserInputCommand struct {
	Text     string
	Metadata *Metadata
}

func (c *UserInputCommand) ToHash() string {
	return ToHash(c.Metadata)
}

func (c *UserInputCommand) ToProtoHash() string {
	return ToProtoHash(c.Metadata)
}

func (c *UserInputCommand) Debug() string {
	return fmt.Sprintf("cmd: %s, state name: %s, text: %s, hash: %s, data: %s", c.Metadata.Cmd, c.Metadata.Place, c.Text, c.ToHash(), c.GetCommand())
}

func (c *UserInputCommand) GetMetadata() *Metadata {
	return c.Metadata
}

func (c *UserInputCommand) GetCommand() string {
	return c.Text
}

func (c *UserInputCommand) GetCaption() string {
	return "user_text_command"
}

type ButtonPressedCommand struct {
	ButtonCommand string
	ButtonText    string
	Metadata      *Metadata
}

func (c *ButtonPressedCommand) ToHash() string {
	return ToHash(c.Metadata)
}

func (c *ButtonPressedCommand) ToProtoHash() string {
	return ToProtoHash(c.Metadata)
}

func (c *ButtonPressedCommand) Debug() string {
	return fmt.Sprintf("cmd: %s, state name: %s, button cmd: %s, button text: %s, hash: %s", c.Metadata.Cmd, c.Metadata.Place, c.ButtonCommand, c.ButtonText, c.ToHash())
}

func (c *ButtonPressedCommand) GetMetadata() *Metadata {
	return c.Metadata
}

func (c *ButtonPressedCommand) GetCommand() string {
	return c.ButtonCommand
}

func (c *ButtonPressedCommand) GetCaption() string {
	return c.ButtonText
}

func CreateCommand(cmd string, place string, arguments []interface{}) Command {
	if cmd == "text_input" {
		text := ""

		if len(arguments) > 0 {
			text = arguments[0].(string)
		}

		return &UserInputCommand{
			Text:     text,
			Metadata: &Metadata{Cmd: cmd, Place: place},
		}
	}

	if cmd == "button" {
		buttonCommand := ""
		buttonText := ""

		if len(arguments) > 0 {
			buttonCommand = arguments[0].(string)
		}

		if len(arguments) > 1 {
			buttonText = arguments[1].(string)
		}

		return &ButtonPressedCommand{
			ButtonCommand: buttonCommand,
			ButtonText:    buttonText,
			Metadata:      &Metadata{Cmd: cmd, Place: place},
		}
	}

	return nil
}

func ToHash(metadata *Metadata) string {
	if metadata == nil {
		panic("METADATA = NIL")
	}
	hash := md5.Sum([]byte(metadata.Cmd + metadata.Place))
	return hex.EncodeToString(hash[:])
}

func ToProtoHash(metadata *Metadata) string {
	hash := md5.Sum([]byte(metadata.Cmd))
	return hex.EncodeToString(hash[:])
}
