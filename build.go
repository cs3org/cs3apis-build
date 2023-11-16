package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	gitMail = flag.String("git-author-mail", "cs3org-bot@hugo.labkode.com", "Git author mail")
	gitName = flag.String("git-author-name", "cs3org-bot", "Git author name")
	gitSSH  = flag.Bool("git-ssh", false, "Use git protocol instead of https for cloning repos")

	_pushGo     = flag.Bool("push-go", false, "Push Go library to github.com/cs3org/go-cs3apis")
	_pushPython = flag.Bool("push-python", false, "Push Python library to github.com/cs3org/python-cs3apis")
	_pushJs     = flag.Bool("push-js", false, "Push Js library to github.com/cs3org/js-cs3apis")
	_pushNode   = flag.Bool("push-node", false, "Push Node.js library to github.com/cs3org/node-cs3apis")
)

func init() {
	flag.Parse()
}

func getProtoOS() string {
	switch runtime.GOOS {
	case "darwin":
		return "osx"
	case "linux":
		return "linux"
	default:
		panic("no build procedure for " + runtime.GOOS)
	}
}

func clone(repo, dir string) {
	repo = getRepo(repo) // get git or https repo location
	cmd := exec.Command("git", "clone", "--quiet", repo)
	cmd.Dir = dir
	run(cmd)
}

func checkout(branch, dir string) {
	// See https://stackoverflow.com/questions/26961371/switch-on-another-branch-create-if-not-exists-without-checking-if-already-exi
	cmd := exec.Command("bash", "-c", fmt.Sprintf("git checkout %s || git checkout -b %s", branch, branch))
	cmd.Dir = dir
	run(cmd)
}

func update(dir string) error {
	cmd := exec.Command("git", "pull", "--quiet")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	return cmd.Run()
}

func isRepoDirty(repo string) bool {
	cmd := exec.Command("git", "status", "-s")
	cmd.Dir = repo
	changes := runAndGet(cmd)
	if changes != "" {
		fmt.Println("repo is dirty")
		fmt.Println(changes)
	}
	return changes != ""
}

func getCommitID(dir string) string {
	if os.Getenv("BUILD_GIT_COMMIT") != "" {
		return os.Getenv("BUILD_GIT_COMMIT")
	}

	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = dir
	commit := runAndGet(cmd)
	return commit
}

func getRepo(repo string) string {
	if *gitSSH {
		return fmt.Sprintf("git@github.com:%s", repo)
	}
	return fmt.Sprintf("https://github.com/%s", repo)
}

func commit(repo, msg string) {
	// set correct author name and mail
	cmd := exec.Command("git", "config", "user.email", *gitMail)
	cmd.Dir = repo
	run(cmd)

	cmd = exec.Command("git", "config", "user.name", *gitName)
	cmd.Dir = repo
	run(cmd)

	// check if repo is dirty
	if !isRepoDirty(repo) {
		// nothing to do
		return
	}

	cmd = exec.Command("git", "add", ".")
	cmd.Dir = repo
	run(cmd)

	cmd = exec.Command("git", "commit", "-m", msg)
	cmd.Dir = repo
	run(cmd)
}

func push(repo string) {
	protoBranch := getGitBranch(".")
	cmd := exec.Command("git", "push", "--set-upstream", "origin", protoBranch)
	cmd.Dir = repo
	run(cmd)
}

func getGitBranch(repo string) string {
	// check if branch is provided by env variable
	if os.Getenv("BUILD_GIT_BRANCH") != "" {
		return os.Getenv("BUILD_GIT_BRANCH")
	}

	// obtain branch from repo
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = repo
	branch := runAndGet(cmd)
	return branch
}

// getVersionFromGit returns a version string that identifies the currently
// checked out git commit.
func getVersionFromGit(repodir string) string {
	cmd := exec.Command("git", "describe",
		"--long", "--tags", "--dirty", "--always")
	cmd.Dir = repodir
	out, err := cmd.Output()
	if err != nil {
		panic(fmt.Sprintf("git describe returned error: %v\n", err))
	}

	version := strings.TrimSpace(string(out))
	return version
}

func run(cmd *exec.Cmd) {
	var b bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &b)
	cmd.Stdout = mw
	cmd.Stderr = mw
	err := cmd.Run()
	fmt.Println(cmd.Dir, cmd.Args)
	fmt.Println(b.String())
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		os.Exit(1)
	}
}

func runAndGet(cmd *exec.Cmd) string {
	var b bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &b)
	cmd.Stderr = mw
	out, err := cmd.Output()
	fmt.Println(cmd.Dir, cmd.Args)
	fmt.Println(b.String())
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		os.Exit(1)
	}
	return strings.TrimSpace(string(out))
}

// Works with Go 1.8+
// https://stackoverflow.com/questions/32649770/how-to-get-current-gopath-from-code
func getGoPath() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	return gopath
}

func sed(dir, suffix, old, new string) {
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, suffix) {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			newData := strings.ReplaceAll(string(data), old, new)
			err = ioutil.WriteFile(path, []byte(newData), 0)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func find(patterns ...string) []string {
	var files []string
	for _, p := range patterns {
		fs, err := filepath.Glob(p)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		files = append(files, fs...)
	}
	return files
}

func findProtos() []string {
	return find("cs3/*/*.proto", "cs3/*/*/*.proto", "cs3/*/*/*/*.proto")
}

func findFolders() []string {
	var folders []string
	err := filepath.Walk("cs3",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if info.IsDir() {
				folders = append(folders, path)
			}
			return nil
		})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return folders
}

func generate() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting generation of protobuf language bindings ...")
	fmt.Printf("current working directory: %s\n", cwd)

	cmd := exec.Command("git", "config", "--global", "--add", "safe.directory", cwd)
	run(cmd)

	// Remove build dir
	os.RemoveAll("build")
	os.MkdirAll("build", 0755)

	languages := []string{"go", "js", "node", "python"}

	// prepare language git repos
	for _, l := range languages {
		target := fmt.Sprintf("%s-cs3apis", l)
		fmt.Println("cloning repo for " + target)

		// Clone Go repo and set branch to current branch
		clone("cs3org/"+target, "build")
		protoBranch := getGitBranch(".")
		targetBranch := getGitBranch("build/" + target)
		fmt.Printf("Proto branch: %s\n%s branch: %s\n", l, protoBranch, targetBranch)

		if targetBranch != protoBranch {
			checkout(protoBranch, "build/"+target)
		}

		// remove leftovers (existing defs)
		os.RemoveAll(fmt.Sprintf("build/%s/cs3", target))

	}

	fmt.Println("Generating ...")
	cmd = exec.Command("buf", "generate")
	run(cmd)

	for _, l := range languages {
		target := fmt.Sprintf("%s-cs3apis", l)
		fmt.Println("Commiting changes for " + target)

		if !isRepoDirty("build/" + target) {
			fmt.Println("Repo is clean, nothing to do")
		}

		// get proto repo commit id
		hash := getCommitID(".")
		repo := "build/" + target
		msg := "Synced to https://github.com/cs3org/cs3apis/tree/" + hash
		commit(repo, msg)
	}
	fmt.Println("Generation done!")
}

func pushPython() {
	push("build/python-cs3apis")
}

func pushGo() {
	push("build/go-cs3apis")
}

func pushJS() {
	push("build/js-cs3apis")
}

func pushNode() {
	push("build/node-cs3apis")
}

func main() {
	generate()

	if *_pushGo {
		fmt.Println("Pushing Go ...")
		pushGo()
	}

	if *_pushPython {
		fmt.Println("Pushing Python ...")
		pushPython()
	}

	if *_pushJs {
		fmt.Println("Pushing Js ...")
		pushJS()
	}

	if *_pushNode {
		fmt.Println("Pushing Node.js ...")
		pushNode()
	}
}
