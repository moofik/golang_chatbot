package runtime

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

const TYPE_TEXT_INPUT = "text_input"
const TYPE_PENDING = "pending"
const TYPE_BUTTON = "button"
const TYPE_INSTANT = "instant"
const TYPE_VALIDATE = "validate"

const MARKER_PENDING_PREFIX = "PENDING_"

type CommandMetadata struct {
	Cmd        string
	Place      string
	Uniqueness string
}

type Command interface {
	ToHash() string
	ToProtoHash() string
	ToUniquenessHash() string
	Debug() string
	GetMetadata() *CommandMetadata
	GetInput() string
	GetCaption() string
	Pass(p ChatProvider, initCmd Command, t TokenProxy) (bool, error)
	GetType() string
}

type ForceCommand interface {
	Command
	GetExecutionDate() string
}

//
type InstantTransitionCommand struct {
	Metadata *CommandMetadata
}

func (c *InstantTransitionCommand) ToHash() string {
	return ToHash(c.Metadata)
}

func (c *InstantTransitionCommand) ToProtoHash() string {
	return ToProtoHash(c.Metadata)
}

func (c *InstantTransitionCommand) Debug() string {
	return fmt.Sprintf("cmd: %s, state name: %s, hash: %s, CMD + PLACE + UNIQ: %s", c.Metadata.Cmd, c.Metadata.Place, c.ToHash(), c.Metadata.Cmd+c.Metadata.Place+c.Metadata.Uniqueness)
}

func (c *InstantTransitionCommand) GetMetadata() *CommandMetadata {
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

func (c *InstantTransitionCommand) Pass(p ChatProvider, initCmd Command, t TokenProxy) (bool, error) {
	return true, nil
}

func (c *InstantTransitionCommand) GetType() string {
	return "instant"
}

//
type UserInputCommand struct {
	Text     string
	Metadata *CommandMetadata
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
	return fmt.Sprintf("cmd: %s, state name: %s, hash: %s, CMD + PLACE + UNIQ: %s", c.Metadata.Cmd, c.Metadata.Place, c.ToHash(), c.Metadata.Cmd+c.Metadata.Place+c.Metadata.Uniqueness)
}

func (c *UserInputCommand) GetMetadata() *CommandMetadata {
	return c.Metadata
}

func (c *UserInputCommand) GetInput() string {
	return c.Text
}

func (c *UserInputCommand) GetCaption() string {
	return "user_text_command"
}

func (c *UserInputCommand) Pass(p ChatProvider, initCmd Command, t TokenProxy) (bool, error) {
	return true, nil
}

func (c *UserInputCommand) GetType() string {
	return "text_input"
}

//
type ButtonPressedCommand struct {
	ButtonCommand string
	ButtonText    string
	ButtonUrl     string
	Metadata      *CommandMetadata
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
	return fmt.Sprintf("cmd: %s, state name: %s, hash: %s, CMD + PLACE + UNIQ: %s", c.Metadata.Cmd, c.Metadata.Place, c.ToHash(), c.Metadata.Cmd+c.Metadata.Place+c.Metadata.Uniqueness)
}

func (c *ButtonPressedCommand) GetMetadata() *CommandMetadata {
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

func (c *ButtonPressedCommand) GetUrl() string {
	return c.ButtonUrl
}

func (c *ButtonPressedCommand) Pass(p ChatProvider, initCmd Command, t TokenProxy) (bool, error) {
	return true, nil
}

type RecognizeInputCommand struct {
	Text     string
	Metadata *CommandMetadata
	Marker   string
}

func (c *RecognizeInputCommand) ToUniquenessHash() string {
	return ToUniquenessHash(c.Metadata)
}

func (c *RecognizeInputCommand) ToHash() string {
	return ToHash(c.Metadata)
}

func (c *RecognizeInputCommand) ToProtoHash() string {
	return ToProtoHash(c.Metadata)
}

func (c *RecognizeInputCommand) Debug() string {
	return fmt.Sprintf("cmd: %s, state name: %s, hash: %s, CMD + PLACE + UNIQ: %s", c.Metadata.Cmd, c.Metadata.Place, c.ToHash(), c.Metadata.Cmd+c.Metadata.Place+c.Metadata.Uniqueness)
}

func (c *RecognizeInputCommand) GetMetadata() *CommandMetadata {
	return c.Metadata
}

func (c *RecognizeInputCommand) GetInput() string {
	return c.Text
}

func (c *RecognizeInputCommand) GetCaption() string {
	return "recognize_input"
}

func (c *RecognizeInputCommand) Pass(p ChatProvider, initCmd Command, t TokenProxy) (bool, error) {
	return c.Marker == initCmd.GetInput(), nil
}

func (c *RecognizeInputCommand) GetType() string {
	return "recognize_input"
}

type PendingCommand struct {
	Text     string
	Metadata *CommandMetadata
	Marker   string
}

func (c *PendingCommand) ToUniquenessHash() string {
	return ToUniquenessHash(c.Metadata)
}

func (c *PendingCommand) ToHash() string {
	return ToHash(c.Metadata)
}

func (c *PendingCommand) ToProtoHash() string {
	return ToProtoHash(c.Metadata)
}

func (c *PendingCommand) Debug() string {
	return fmt.Sprintf("cmd: %s, state name: %s, hash: %s, CMD + PLACE + UNIQ: %s", c.Metadata.Cmd, c.Metadata.Place, c.ToHash(), c.Metadata.Cmd+c.Metadata.Place+c.Metadata.Uniqueness)
}

func (c *PendingCommand) GetMetadata() *CommandMetadata {
	return c.Metadata
}

func (c *PendingCommand) GetInput() string {
	return c.Marker
}

func (c *PendingCommand) GetCaption() string {
	return "pending"
}

func (c *PendingCommand) Pass(p ChatProvider, initCmd Command, t TokenProxy) (bool, error) {
	return c.Marker == initCmd.GetInput(), nil
}

func (c *PendingCommand) GetType() string {
	return TYPE_PENDING
}

func CreatePendingCommand(text string, marker string) Command {
	marker = MARKER_PENDING_PREFIX + marker

	return &PendingCommand{
		Marker:   marker,
		Metadata: &CommandMetadata{Cmd: "text_input", Uniqueness: marker},
	}
}

//
func CreateCommand(cmd string, place string, arguments []interface{}, commandRegistry func(string, string, []interface{}) Command) Command {
	if cmd == "text_input" {
		text := ""

		if len(arguments) > 0 {
			text = arguments[0].(string)
		}

		return &UserInputCommand{
			Text:     text,
			Metadata: &CommandMetadata{Cmd: cmd, Place: place},
		}
	}

	if cmd == "button" {
		buttonCommand := ""
		buttonText := ""
		buttonUrl := ""

		if len(arguments) > 0 {
			buttonCommand = arguments[0].(string)
		}

		if len(arguments) > 1 {
			buttonText = arguments[1].(string)
		}

		if len(arguments) > 2 {
			buttonUrl = arguments[2].(string)
		}

		return &ButtonPressedCommand{
			ButtonCommand: buttonCommand,
			ButtonText:    buttonText,
			ButtonUrl:     buttonUrl,
			Metadata:      &CommandMetadata{Cmd: cmd, Place: place, Uniqueness: buttonCommand},
		}
	}

	if cmd == "instant" {
		return &InstantTransitionCommand{Metadata: &CommandMetadata{
			Cmd:        cmd,
			Place:      place,
			Uniqueness: "",
		}}
	}

	if cmd == "recognize_input" {
		marker := ""

		if len(arguments) > 0 {
			marker = arguments[0].(string)
		}

		return &RecognizeInputCommand{
			Marker:   marker,
			Metadata: &CommandMetadata{Cmd: "text_input", Place: place, Uniqueness: marker},
		}
	}

	if cmd == "pending" {
		marker := ""

		if len(arguments) > 0 {
			marker = arguments[0].(string)
		}

		return CreatePendingCommand("", marker)
	}

	if commandRegistry != nil {
		return commandRegistry(cmd, place, arguments)
	}

	return nil
}

func ToHash(metadata *CommandMetadata) string {
	hash := md5.Sum([]byte(metadata.Cmd + metadata.Place + metadata.Uniqueness))
	return hex.EncodeToString(hash[:])
}

func ToProtoHash(metadata *CommandMetadata) string {
	hash := md5.Sum([]byte(metadata.Cmd))
	return hex.EncodeToString(hash[:])
}

func ToUniquenessHash(metadata *CommandMetadata) string {
	hash := md5.Sum([]byte(metadata.Uniqueness))
	return hex.EncodeToString(hash[:])
}
