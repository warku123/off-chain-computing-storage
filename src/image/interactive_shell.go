package image

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Interactive_shell(v *ipfs_api) (err error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		cmdstring, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		out, exit_signal, err := Execute_command(v, cmdstring)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if len(out) > 0 {
			fmt.Println("log:" + out)
		}

		if exit_signal {
			break
		}
	}
	return nil
}

func Execute_command(v *ipfs_api, command string) (out_string string, exit_signal bool, err error) {
	command = strings.TrimSuffix(command, "\n")
	command = strings.TrimSpace(command)

	arr_command_args := strings.Fields(command)

	switch arr_command_args[0] {
	case "exit":
		return "", true, nil
	case "write":
		if len(arr_command_args) == 3 {
			err = v.WriteDB(arr_command_args[1], arr_command_args[2])
			if err != nil {
				return "", false, err
			}
		} else {
			return "", false, errors.New("Need 2 args to write DB")
		}
	case "read":
		if len(arr_command_args) == 3 {
			version, err := strconv.Atoi(arr_command_args[2])
			if err != nil {
				return "", false, err
			}

			out, err := v.ReadDB(arr_command_args[1], version)
			if err != nil {
				return "", false, err
			}
			return out, false, err
		} else {
			return "", false, errors.New("Need 2 args to read DB")
		}
	case "persist":
		if len(arr_command_args) == 3 {
			err = v.WriteDB(arr_command_args[1], arr_command_args[2])
			if err != nil {
				return "", false, err
			}
		} else {
			return "", false, errors.New("Need 2 args to persist data")
		}
	}
	return "", false, nil
}
