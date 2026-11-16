package main

import "antispam/console"
import "bytes"
import "os"
import "os/exec"
import "path/filepath"
import "strings"

func main() {

	console.Group("Build App")

	result := false
	go_compiler := "go"
	ld_flags := ""

	if len(os.Args) == 2 {

		if os.Args[1] == "--debug" {
			ld_flags = ""
		} else if os.Args[1] == "--release" {
			ld_flags = "-s -w"
		}

	}

	toolchain_folder, err0 := os.Getwd()

	if err0 == nil && strings.HasSuffix(toolchain_folder, "/toolchain") {

		folder := filepath.Dir(toolchain_folder)
		go_output := folder + "/build/antispam-gtk_linux_amd64"

		gtk_check := exec.Command("pkg-config", "--exists", "gtk4")
		gtk_err := gtk_check.Run()

		if gtk_err != nil {

			console.Error("GTK4 development files are required (pkg-config gtk4)")
			console.GroupEndResult(false, "Build App")

			return

		}

		go_args := []string{"build"}

		if ld_flags != "" {
			go_args = append(go_args, "-ldflags", ld_flags)
		}

		go_args = append(go_args, "-o", go_output, "./cmds/antispam-gtk")

		cmd := exec.Command(
			"env",
			"CGO_ENABLED=1",
			"GOOS=linux",
			"GOARCH=amd64",
			go_compiler,
		)

		cmd.Args = append(cmd.Args, go_args...)
		cmd.Dir = folder + "/source"

		console.Log(cmd.String())

		var stdout bytes.Buffer
		var stderr bytes.Buffer

		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err1 := cmd.Run()

		if err1 == nil {

			console.Info("> linux / amd64")
			console.Log("> " + go_output)

			result = true

		} else {

			console.Error("> linux / amd64")

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

	console.GroupEndResult(result, "Build App")

}
