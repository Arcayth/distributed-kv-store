package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type Command int

const (
	SetCommand Command = iota

	GetCommand

	ClearCommand

	UnknownCommand
)

func (c Command) String() string {
	switch c {
	case SetCommand:
		return "Set"
	case GetCommand:
		return "Get"
	case ClearCommand:
		return "Clear"
	default:
		return "Unknown"
	}
}

type CommandInfo struct {
	Cmd   Command
	Key   string
	Value string
}

func getU8(cur *bytes.Reader) (uint8, error) {
	var data [1]byte
	if _, err := cur.Read(data[:]); err != nil {
		if err == io.EOF {
			return 0, errors.New("unexpected end of input")
		}
		return 0, err
	}
	return data[0], nil
}

func getU32(cur *bytes.Reader) (uint32, error) {
	var data [4]byte
	if _, err := cur.Read(data[:]); err != nil {
		if err == io.EOF {
			return 0, errors.New("unexpected end of input")
		}
		return 0, err
	}
	return binary.LittleEndian.Uint32(data[:]), nil
}

func getString(cur *bytes.Reader) (string, error) {
	length, err := binary.ReadUvarint(cur)
	if err != nil {
		return "", err
	}

	strBuffer := make([]byte, length)
	_, err = io.ReadFull(cur, strBuffer)
	if err != nil {
		return "", nil
	}
	return string(strBuffer), nil
}

func parse(data *bytes.Reader) (CommandInfo, error) {
	cmdInfo := CommandInfo{}

	if _, err := getU8(data); err != nil {
		return cmdInfo, err
	}

	command, err := getU32(data)
	if err != nil {
		cmdInfo.Cmd = UnknownCommand
		return cmdInfo, err
	}

	switch command {
	case 0:
		cmdInfo, err = parse_set(data)
		return cmdInfo, nil
	case 1:
		cmdInfo, err = parse_get(data)
		return cmdInfo, nil
	case 2:
		cmdInfo, err = parse_clear(data)
		return cmdInfo, nil
	default:
		cmdInfo.Cmd = UnknownCommand
		return cmdInfo, errors.New("Unknown command")
	}
}

func parse_get(data *bytes.Reader) (CommandInfo, error) {
	cmdInfo := CommandInfo{Cmd: GetCommand}
	key, err := getString(data)
	if err != nil {
		cmdInfo.Key = ""
		return cmdInfo, err
	}

	cmdInfo.Key = key

	return cmdInfo, nil
}

func parse_set(data *bytes.Reader) (CommandInfo, error) {
	cmdInfo := CommandInfo{Cmd: SetCommand}
	key, err := getString(data)
	if err != nil {
		cmdInfo.Key = ""
		return cmdInfo, err
	}
	cmdInfo.Key = key

	value, err := getString(data)
	if err != nil {
		cmdInfo.Value = ""

		return cmdInfo, err
	}
	cmdInfo.Value = value

	return cmdInfo, nil
}

func parse_clear(data *bytes.Reader) (CommandInfo, error) {
	cmdInfo := CommandInfo{Cmd: ClearCommand}
	key, err := getString(data)
	if err != nil {
		cmdInfo.Key = ""
		return cmdInfo, err
	}

	cmdInfo.Key = key

	return cmdInfo, nil
}
func main() {
	data := []byte{1, 1, 0, 0, 0, 3, 102, 111, 111, 3, 98, 97, 114}
	reader := bytes.NewReader(data)
	cmd, err := parse(reader)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("parsed command : ", cmd)
}
