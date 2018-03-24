// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"os"
	"github.com/Top-Ranger/gogitui/helper"
	"path"
	"os/exec"
	"log"
)

func main() {
	defer func() {
		r := recover()
		if r != nil {
			helper.ShowError(fmt.Sprint(r))
			log.Println(r)
			os.Exit(1)
		}
	}()

	fmt.Println("== gogitui ==")
	fmt.Println("A simple git ui managing multiple repositories")
	fmt.Println("Copyright 2018 Marcus Soll")
	fmt.Println("License: MIT")
	fmt.Println()

	config, err := helper.LoadConfig()
	if err != nil {
		panic(err)
	}

	exit := false
	for !exit {
		option, _ := helper.Menu("Select action", []string{"git pull", "git push", "git status", "git difftool", "git commit", "Add repository", "Remove repository", "Exit"})
		fmt.Println(option)
		switch option {
		case "git pull":
			targets, _ :=  helper.Checklist("Select repositories for pull:", config.Repositories, "on")
			if len(targets) == 0 {
				break
			}
			handle, _ := helper.CreateProgressbar("git pull", len(targets))
			for i := range targets {
				helper.SetProgressbarValue(handle, i)
				helper.SetProgressbarHeader(handle, targets[i])
				gitpull := exec.Command("/usr/bin/git", "pull")
				gitpull.Dir = targets[i]
				output, err := gitpull.CombinedOutput()
				fmt.Println(string(output))
				if err != nil {
					helper.ShowError(fmt.Sprintln("Error while git pull at", targets[i], ":\n", string(output), "\n\nError:", err))
				}
			}
			helper.CloseProgressbar(handle)
		case "git push":
			targets, _ :=  helper.Checklist("Select repositories for push:", config.Repositories, "on")
			if len(targets) == 0 {
				break
			}
			handle, _ := helper.CreateProgressbar("git push", len(targets))
			for i := range targets {
				helper.SetProgressbarValue(handle, i)
				helper.SetProgressbarHeader(handle, targets[i])
				gitpush := exec.Command("/usr/bin/git", "push")
				gitpush.Dir = targets[i]
				output, err := gitpush.CombinedOutput()
				fmt.Println(string(output))
				if err != nil {
					helper.ShowError(fmt.Sprintln("Error while git push at", targets[i], ":\n", string(output), "\n\nError:", err))
				}
			}
			helper.CloseProgressbar(handle)
		case "git status":
			targets, _ :=  helper.Checklist("Select repositories for push:", config.Repositories, "off")
			if len(targets) == 0 {
				break
			}
			for i := range targets {
				gitstatus := exec.Command("/usr/bin/git", "status")
				gitstatus.Dir = targets[i]
				output, err := gitstatus.CombinedOutput()
				fmt.Println(string(output))
				if err != nil {
					helper.ShowError(fmt.Sprintln("Error while git status at", targets[i], ":\n", string(output), "\n\nError:", err))
				}
				helper.ShowMessage(fmt.Sprint("Repository ", targets[i], ":\n\n", string(output)))
			}
		case "git difftool":
			targets, _ :=  helper.Checklist("Select repositories for push:", config.Repositories, "off")
			if len(targets) == 0 {
				break
			}
			for i := range targets {
				gitdifftool := exec.Command("/usr/bin/git", "difftool", "--dir-diff")
				gitdifftool.Dir = targets[i]
				output, err := gitdifftool.CombinedOutput()
				fmt.Println(string(output))
				if err != nil {
					helper.ShowError(fmt.Sprintln("Error while git difftool at", targets[i], ":\n", string(output), "\n\nError:", err))
				}
				if len(output) == 0 {
					helper.ShowMessage(fmt.Sprintln("No difference for", targets[i]))
				}
			}
		case "git commit":
			targets, _ :=  helper.Checklist("Select repositories for commit:", config.Repositories, "off")
			if len(targets) == 0 {
				break
			}
			for i := range targets {
				// Test if there are changes
				gitstatus := exec.Command("/usr/bin/git", "status", "--short")
				gitstatus.Dir = targets[i]
				output, err := gitstatus.CombinedOutput()
				fmt.Println(string(output))
				if err != nil {
					helper.ShowError(fmt.Sprintln("Error while git status at", targets[i], ":\n", string(output), "\n\nError:", err))
				}
				if len(output) == 0 {
					helper.ShowMessage(fmt.Sprintln("Nothing to commit for", targets[i]))
					continue
				}

				gitstatus = exec.Command("/usr/bin/git", "status")
				gitstatus.Dir = targets[i]
				output_status, err := gitstatus.CombinedOutput()
				if err != nil {
					helper.ShowError(fmt.Sprintln("Error while git status at", targets[i], ":\n", string(output_status), "\n\nError:", err))
				}

				commitExit := false
				for !commitExit {
					operator, _ := helper.Menu(fmt.Sprintln("Repository:", targets[i], "\n\n", string(output_status)), []string{"git commit -a", "git add -A && git commit", "git difftool", "Do nothing"})
					switch operator {
					case "git commit -a":
						message, _ := helper.TextInput(fmt.Sprintln("Commit message for", targets[i]))
						if message == "" {
							helper.ShowMessage("Aborting due to empty message")
							break
						}
						gitcommand := exec.Command("/usr/bin/git", "commit", "-a", "--message", message)
						gitcommand.Dir = targets[i]
						output, err := gitcommand.CombinedOutput()
						fmt.Println(string(output))
						if err != nil {
							helper.ShowError(fmt.Sprintln("Error while git commit -a at", targets[i], ":\n", string(output), "\n\nError:", err))
							break
						}
						commitExit = true
					case "git add -A && git commit":
						gitcommand := exec.Command("/usr/bin/git", "add", "-A")
						gitcommand.Dir = targets[i]
						output, err := gitcommand.CombinedOutput()
						fmt.Println(string(output))
						if err != nil {
							helper.ShowError(fmt.Sprintln("Error while git add -A at", targets[i], ":\n", string(output), "\n\nError:", err))
							break
						}
						message, _ := helper.TextInput(fmt.Sprintln("Commit message for", targets[i]))
						if message == "" {
							helper.ShowMessage("Aborting due to empty message")
							break
						}
						gitcommand = exec.Command("/usr/bin/git", "commit", "--message", message)
						gitcommand.Dir = targets[i]
						output, err = gitcommand.CombinedOutput()
						fmt.Println(string(output))
						if err != nil {
							helper.ShowError(fmt.Sprintln("Error while git commit at", targets[i], ":\n", string(output), "\n\nError:", err))
							break
						}
						commitExit = true
					case "git difftool":
						gitdifftool := exec.Command("/usr/bin/git", "difftool", "--dir-diff")
						gitdifftool.Dir = targets[i]
						output, err := gitdifftool.CombinedOutput()
						fmt.Println(string(output))
						if err != nil {
							helper.ShowError(fmt.Sprintln("Error while git difftool at", targets[i], ":\n", string(output), "\n\nError:", err))
						}
						if len(output) == 0 {
							helper.ShowMessage(fmt.Sprintln("No difference for", targets[i]))
						}
					case "Do nothing":
						fallthrough
					case "":
						commitExit = true
						break
					default:
						helper.ShowError(fmt.Sprintln("Unknown operator:", operator))
				}
			}
			}
		case "Add repository":
			repository, _ := helper.GetDir()
			if repository == "" {
				break
			}
			_, err := os.Stat(path.Join(repository, ".git"))
			if err != nil {
				helper.ShowError(fmt.Sprintln("Folder", repository, "does not contain a git repository."))
				break
			}
			alreadyAdded := false
			for i := range config.Repositories {
				if config.Repositories[i] == repository {
					alreadyAdded = true
					break
				}
			}
			if alreadyAdded {
				helper.ShowMessage(fmt.Sprintln("Repository", repository, "already added."))
			} else {
				config.Repositories = append(config.Repositories, repository)
				err := config.SaveConfig()
				if err != nil {
					panic(err)
				}
			}
		case "Remove repository":
			if len(config.Repositories) == 0 {
				helper.ShowMessage("No repositories are currently registered.")
				break
			}
			target, _ := helper.Menu("Select repository for removal:", config.Repositories)
			if target == "" {
				break
			}
			// Assume each repository is only once registered
			for i := range config.Repositories {
				if config.Repositories[i] == target {
					config.Repositories = append(config.Repositories[:i], config.Repositories[i+1:]...)
					config.SaveConfig()
					break
				}
			}
		case "Exit":
			fallthrough
		case "":
			exit = true
		default:
			panic(fmt.Sprintln("Unknown option:\n", option))
		}
	}
}