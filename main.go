// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"os"
	"github.com/Top-Ranger/gogitui/helper"
	"path"
	"os/exec"
)

func main() {
	defer func() {
		r := recover()
		if r != nil {
			helper.ShowError(fmt.Sprint(r))
			panic(r)
			os.Exit(1)
		}
	}()

	config, err := helper.LoadConfig()
	if err != nil {
		panic(err)
	}

	exit := false
	for !exit {
		option, _ := helper.Menu("Select action", []string{"git pull", "git push", "Add repository", "Remove repository", "Exit"})
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
					helper.ShowError(fmt.Sprintln("Error while git pull at", targets[i], ":\n", string(output), "\n\nError:", err))
				}
			}
			helper.CloseProgressbar(handle)
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