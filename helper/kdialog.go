// SPDX-License-Identifier: MIT

package helper

import (
	"os/exec"
	"strings"
	"log"
	"regexp"
	"fmt"
	"strconv"
	"errors"
)

// Helper regular expression
var outputMatcher = regexp.MustCompile("\".*?\"")

// Opens a checklist which allows the selection of multiple items
func Checklist(header string, options []string, defaultMode string) ([]string, error) {
	cmd := exec.Command("/usr/bin/kdialog", "--title", "gogitui", "--checklist", header)
	for i := range options {
		if strings.HasPrefix(options[i],"-") {
			log.Println("Ignoring option", options[i])
			continue
		}
		cmd.Args = append(cmd.Args, fmt.Sprint("'", options[i], "'"), options[i], defaultMode)
	}
	out, err := cmd.Output()
	if err != nil {
		return []string{}, err
	}
	outputIndex := outputMatcher.FindAllIndex(out, -1)
	outputOptions := make([]string, 0, len(outputIndex))
	for i := range outputIndex {
		outputOptions = append(outputOptions, string(out[outputIndex[i][0]+2:outputIndex[i][1]-2]))
	}
	return outputOptions, nil
}

// Shows a menu which allows the selection of exactly one item or "" if none is selected
func Menu(header string, options []string) (string, error) {
	cmd := exec.Command("/usr/bin/kdialog", "--title", "gogitui", "--menu", header)
	for i := range options {
		if strings.HasPrefix(options[i],"-") {
			log.Println("Ignoring option", options[i])
			continue
		}
		cmd.Args = append(cmd.Args, options[i], options[i])
	}
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}

// Returns the path of an existing directory or "" if none is selected
func GetDir() (string, error) {
	cmd := exec.Command("/usr/bin/kdialog", "--getexistingdirectory", "--title", "gogitui")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// Shows a message
func ShowMessage(msg string) error {
	cmd := exec.Command("/usr/bin/kdialog", "--title", "gogitui", "--msgbox", msg)
	err := cmd.Run()
	return err
}

// Shows an error message
func ShowError(msg string) error {
	cmd := exec.Command("/usr/bin/kdialog", "--title", "gogitui", "--error", msg)
	err := cmd.Run()
	return err
}

// Creates a progress bar and returns a handle. The handle is used to control the progress bar later
func CreateProgressbar(header string, max int) (string, error) {
	cmd := exec.Command("/usr/bin/kdialog", "--title", "gogitui", "--progressbar", header, strconv.Itoa(max))
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}

// Updates the header text of the progress bar
func SetProgressbarHeader (handle, header string) error {
	handleSplit := strings.Split(handle, " ")
	if len(handleSplit) != 2 {
		return errors.New("Invalid handle")
	}
	cmd := exec.Command("/usr/bin/qdbus", handleSplit[0], handleSplit[1], "setLabelText", header)
	return cmd.Start()
}

// Sets the value of the progress bar
func SetProgressbarValue (handle string, value int) error {
	handleSplit := strings.Split(handle, " ")
	if len(handleSplit) != 2 {
		return errors.New("Invalid handle")
	}
	cmd := exec.Command("/usr/bin/qdbus", handleSplit[0], handleSplit[1], "Set", "", "value",  strconv.Itoa(value))
	return cmd.Start()
}

// Closes the progress bar
func CloseProgressbar (handle string) error {
	handleSplit := strings.Split(handle, " ")
	if len(handleSplit) != 2 {
		return errors.New("Invalid handle")
	}
	cmd := exec.Command("/usr/bin/qdbus", handleSplit[0], handleSplit[1], "close")
	return cmd.Start()
}

func TextInput(header string) (string, error) {
	cmd := exec.Command("/usr/bin/kdialog", "--title", "gogitui", "--textinputbox", header)
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}