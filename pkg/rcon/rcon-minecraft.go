package rcon

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/seeruk/minecraft-rcon/rcon"
)

type Client struct {
	//*rcon.Client
	host     string
	port     int
	password string
}

type User struct {
	Health  float32
	XpLevel uint
	Position
}

type Position struct {
	X float32
	Y float32
	Z float32
}

func New(host string, port int, password string) (*Client, error) {
	return &Client{
		host:     host,
		port:     port,
		password: password,
	}, nil
}

type Command struct {
	command            string
	expression         string
	expressionNotFound string
}

func (c Client) command(command Command) ([]string, error) {
	// 毎回 client を作り直す
	client, err := rcon.NewClient(c.host, c.port, c.password)
	if err != nil {
		return nil, err
	}

	response, err := client.SendCommand(command.command)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(command.expression)
	extracted := re.FindStringSubmatch(response)
	if len(extracted) == 0 {
		re = regexp.MustCompile(command.expressionNotFound)
		extracted = re.FindStringSubmatch(response)
		if len(extracted) == 0 {
			return nil, fmt.Errorf(`"%s" is not match to "%s"`, command.expression, response)
		}
		return nil, nil
	}

	return extracted[1:], nil
}

func (c Client) List() ([]string, error) {
	result, err := c.command(Command{
		command:            `list`,
		expression:         `There are [0-9].* of a max [0-9].* players online: (.*)$`,
		expressionNotFound: `There are 0 of a max [0-9].* players online:`,
	})
	if err != nil {
		return nil, err
	} else if result == nil {
		return []string{}, nil
	}
	return strings.Split(result[0], ", "), nil
}

func (c Client) WhitelistAdd(username string) error {
	_, err := c.command(Command{
		command:            fmt.Sprintf(`whitelist add %s`, username),
		expression:         fmt.Sprintf(`Added %s to the whitelist`, username),
		expressionNotFound: `!!!not much!!!`,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c Client) WhitelistRemove(username string) error {
	_, err := c.command(Command{
		command:    fmt.Sprintf(`whitelist remove %s`, username),
		expression: fmt.Sprintf(`Removed %s from the whitelist`, username),

		expressionNotFound: `!!!not much!!!`,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c Client) WhitelistList() ([]string, error) {
	result, err := c.command(Command{
		command:            `whitelist list`,
		expression:         `There are [0-9].* whitelisted players: (.*)`,
		expressionNotFound: `There are no whitelisted players`,
	})
	if err != nil {
		return nil, err
	} else if result == nil {
		return []string{}, nil
	}
	return strings.Split(result[0], ", "), nil
}

// TODO: return nil,nil を辞める
func (c Client) DataGetEntity(username string) (*User, error) {
	array, err := c.command(Command{
		command:            fmt.Sprintf(`data get entity %s Health`, username),
		expression:         fmt.Sprintf(`^%s has the following entity data: (.*)f$`, username),
		expressionNotFound: `No entity was found`,
	})
	if err != nil {
		return nil, err
	} else if array == nil {
		return nil, nil
	}
	health, err := strconv.ParseFloat(array[0], 32)
	if err != nil {
		return nil, nil
	}

	array, err = c.command(Command{
		command:            fmt.Sprintf(`data get entity %s XpLevel`, username),
		expression:         fmt.Sprintf(`^%s has the following entity data: (.*)$`, username),
		expressionNotFound: `No entity was found`,
	})
	if err != nil {
		return nil, err
	} else if array == nil {
		return nil, nil
	}
	xpLevel, err := strconv.Atoi(array[0])
	if err != nil {
		return nil, nil
	}

	array, err = c.command(Command{
		command:            fmt.Sprintf(`data get entity %s Pos`, username),
		expression:         fmt.Sprintf(`^%s has the following entity data: \[(.*?)d, (.*?)d, (.*?)d\]`, username),
		expressionNotFound: `No entity was found`,
	})
	if err != nil {
		return nil, err
	} else if array == nil {
		return nil, nil
	}
	posX, err := strconv.ParseFloat(array[0], 32)
	if err != nil {
		return nil, err
	}
	posY, err := strconv.ParseFloat(array[1], 32)
	if err != nil {
		return nil, err
	}
	posZ, err := strconv.ParseFloat(array[2], 32)
	if err != nil {
		return nil, err
	}

	user := &User{
		Health:  float32(health),
		XpLevel: uint(xpLevel),
		Position: Position{
			X: float32(posX),
			Y: float32(posY),
			Z: float32(posZ),
		},
	}
	return user, nil
}

func (c Client) Title(msg string) ([]string, error) {
	result, err := c.command(Command{
		command:            fmt.Sprintf(`title @a title {"text": "%s"}`, msg),
		expression:         `Showing new title for (.*)$`,
		expressionNotFound: `No player was found`,
	})
	if err != nil {
		return nil, err
	} else if result == nil {
		return []string{}, nil
	}
	return strings.Split(result[0], ", "), nil
}
