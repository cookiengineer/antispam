package main

import "antispam/console"
import "bytes"
import "os/exec"
import "os"
import "path/filepath"
import "strings"

func main() {

	console.Group("Build")

	result := false
	go_compiler := "go"
	ld_flags := ""

	if len(os.Args) == 2 {

		if os.Args[1] == "--debug" {
			ld_flags = "-s -w"
		}

	}

	toolchain_folder, err0 := os.Getwd()

	if err0 == nil && go_compiler != "" && strings.HasSuffix(toolchain_folder, "/toolchain") {

		folder := filepath.Dir(toolchain_folder)

		entries, err1 := os.ReadDir(folder + "/source/cmds")

		if err1 == nil {

			for _, cmd_folder := range entries {

				go_cmd := cmd_folder.Name()
				go_os := "linux"
				go_arch := "amd64"
				go_output := folder + "/build/" + go_cmd + "_" + go_os + "_" + go_arch
				go_source := folder + "/source/cmds/" + go_cmd + "/main.go"

				var stdout bytes.Buffer
				var stderr bytes.Buffer

				cmd := exec.Command(
					"env",
					"CGO_ENABLED=0",
					"GOOS="+go_os,
					"GOARCH="+go_arch,
					go_compiler,
					"build",
					"-ldflags",
					ld_flags,
					"-o",
					go_output,
					go_source,
				)
				cmd.Dir = folder + "/source"

				console.Log(cmd.String())

				cmd.Stdout = &stdout
				cmd.Stderr = &stderr

				err1 := cmd.Run()

				if err1 == nil {

					console.Info("> " + go_os + " / " + go_arch)
					console.Log("> " + go_output)

					result = true

				} else {

					console.Error("> " + go_os + " / " + go_arch)

					result = false

					stdout_message := strings.TrimSpace(string(stdout.Bytes()))
					stderr_message := strings.TrimSpace(string(stderr.Bytes()))

					if stdout_message != "" {
						console.Error(stdout_message)
					}

					if stderr_message != "" {
						console.Error(stderr_message)
					}

				}

			}

		}

	}

	console.GroupEndResult(result, "Build")

}
