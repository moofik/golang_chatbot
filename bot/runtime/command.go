package runtime

import (
	"bot-daedalus/petrinet"
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

const TYPE_TEXT_INPUT = "text_input"
const TYPE_BUTTON = "button"
const TYPE_INSTANT = "instant"
const TYPE_VALIDATE = "validate"

type Metadata struct {
	Cmd        string
	Place      string
	Uniqueness string
}

type Command interface {
	ToHash() string
	ToProtoHash() string
	ToUniquenessHash() string
	Debug() string
	GetMetadata() *Metadata
	GetInput() string
	GetCaption() string
	Pass(p ChatProvider, t TokenProxy, tr *petrinet.Transition) bool
	GetType() string
}

type ForceCommand interface {
	Command
	GetExecutionDate() string
}

//
type InstantTransitionCommand struct {
	Metadata *Metadata
}

func (c *InstantTransitionCommand) ToHash() string {
	return ToHash(c.Metadata)
}

func (c *InstantTransitionCommand) ToProtoHash() string {
	return ToProtoHash(c.Metadata)
}

func (c *InstantTransitionCommand) Debug() string {
	return fmt.Sprintf("cmd: %s, state name: %s, hash: %s", c.Metadata.Cmd, c.Metadata.Place, c.ToHash())
}

func (c *InstantTransitionCommand) GetMetadata() *Metadata {
	return c.Metadata
}

func (c *InstantTransitionCommand) GetInput() string {
	return ""
}

func (c *InstantTransitionCommand) GetCaption() string {
	return "Моментальный переход"
}

func (c *InstantTransitionCommand) ToUniquenessHash() string {
	return ToUniquenessHash(c.Metadata)
}

func (c *InstantTransitionCommand) Pass(p ChatProvider, t TokenProxy, tr *petrinet.Transition) bool {
	return true
}

func (c *InstantTransitionCommand) GetType() string {
	return "instant"
}

//
type UserInputCommand struct {
	Text     string
	Metadata *Metadata
}

func (c *UserInputCommand) ToUniquenessHash() string {
	return ToUniquenessHash(c.Metadata)
}

func (c *UserInputCommand) ToHash() string {
	return ToHash(c.Metadata)
}

func (c *UserInputCommand) ToProtoHash() string {
	return ToProtoHash(c.Metadata)
}

func (c *UserInputCommand) Debug() string {
	return fmt.Sprintf("cmd: %s, state name: %s, text: %s, hash: %s, data: %s", c.Metadata.Cmd, c.Metadata.Place, c.Text, c.ToHash(), c.GetInput())
}

func (c *UserInputCommand) GetMetadata() *Metadata {
	return c.Metadata
}

func (c *UserInputCommand) GetInput() string {
	return c.Text
}

func (c *UserInputCommand) GetCaption() string {
	return "user_text_command"
}

func (c *UserInputCommand) Pass(p ChatProvider, t TokenProxy, tr *petrinet.Transition) bool {
	return true
}

func (c *UserInputCommand) GetType() string {
	return "text_input"
}

//
type ButtonPressedCommand struct {
	ButtonCommand string
	ButtonText    string
	Metadata      *Metadata
}

func (c *ButtonPressedCommand) ToUniquenessHash() string {
	return ToUniquenessHash(c.Metadata)
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

func (c *ButtonPressedCommand) GetInput() string {
	return c.ButtonCommand
}

func (c *ButtonPressedCommand) GetCaption() string {
	return c.ButtonText
}

func (c *ButtonPressedCommand) GetType() string {
	return "button"
}

func (c *ButtonPressedCommand) Pass(p ChatProvider, t TokenProxy, tr *petrinet.Transition) bool {
	return true
}

//
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
			Metadata:      &Metadata{Cmd: cmd, Place: place, Uniqueness: buttonCommand},
		}
	}

	if cmd == "instant" {
		return &InstantTransitionCommand{Metadata: &Metadata{
			Cmd:        cmd,
			Place:      place,
			Uniqueness: "",
		}}
	}

	return nil
}

func ToHash(metadata *Metadata) string {
	hash := md5.Sum([]byte(metadata.Cmd + metadata.Place + metadata.Uniqueness))
	return hex.EncodeToString(hash[:])
}

func ToProtoHash(metadata *Metadata) string {
	hash := md5.Sum([]byte(metadata.Cmd))
	return hex.EncodeToString(hash[:])
}

func ToUniquenessHash(metadata *Metadata) string {
	hash := md5.Sum([]byte(metadata.Uniqueness))
	return hex.EncodeToString(hash[:])
}
