package wplite

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/mattn/go-isatty"
	log "github.com/sirupsen/logrus"
)

type WPLiteEnv struct {
	File  string
	Theme string
	Title string
	User  string
	Pass  string
	Email string
	Port  int
}

type VCS struct {
	GitUrl string
	Branch string
}

type WPLite struct {
	ContainerName string
	OpenOnReady   bool
	ImageUrl      string
	Env           WPLiteEnv
	VCS           VCS
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (w *WPLite) WorkspaceContainerName() string {
	prefix := "wplite-"
	// create hash of the config and current directory
	// to ensure unique container names
	tw := *w
	tw.ImageUrl = ""
	jd, err := json.Marshal(tw)
	if err != nil {
		return prefix + randomString(8)
	}
	wd, err := os.Getwd()
	if err != nil {
		return prefix + randomString(8)
	}
	hashStr := fmt.Sprintf("%s-%s", string(jd), wd)
	h := sha1.New()
	h.Write([]byte(hashStr))
	hash := hex.EncodeToString(h.Sum(nil))
	return prefix + hash[:8]
}

func (w *WPLite) EnsureLocalPaths() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	wpContentDir := cwd + "/wp-content"
	if _, err := os.Stat(wpContentDir); os.IsNotExist(err) {
		os.Mkdir(wpContentDir, 0755)
	}

	staticDir := wpContentDir + "/static"
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		os.Mkdir(staticDir, 0755)
	}

	htaccessFile := cwd + "/.htaccess"
	if _, err := os.Stat(htaccessFile); os.IsNotExist(err) {
		file, err := os.Create(htaccessFile)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}

func (w *WPLite) EnsureEnvFile() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	envFile := path.Join(cwd, w.Env.File)
	reqVars := []string{
		"WP_TITLE",
		"WP_USER",
		"WP_PASS",
		"WP_EMAIL",
		//"WP_THEME",
		//"WP_PORT",
	}
	if _, err := os.Stat(envFile); !os.IsNotExist(err) {
		// set env vars from file
		err := godotenv.Load(w.Env.File)
		if err != nil {
			return err
		}
		w.Env.Title = os.Getenv("WP_TITLE")
		w.Env.User = os.Getenv("WP_USER")
		w.Env.Pass = os.Getenv("WP_PASS")
		w.Env.Email = os.Getenv("WP_EMAIL")
		w.Env.Theme = os.Getenv("WP_THEME")
		if os.Getenv("WP_PORT") != "" {
			w.Env.Port, err = strconv.Atoi(os.Getenv("WP_PORT"))
			if err != nil {
				return err
			}
		}
		return nil
	}
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		return nil
	}
	inVars := make(map[string]string)
	for _, v := range reqVars {
		fmt.Printf("%s: ", v)
		var val string
		fmt.Scanln(&val)
		inVars[v] = val
	}
	file, err := os.Create(envFile)
	if err != nil {
		return err
	}
	defer file.Close()
	for k, v := range inVars {
		file.WriteString(fmt.Sprintf("%s=%s\n", k, v))
		os.Setenv(k, v)
	}
	return nil
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch os := runtime.GOOS; os {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		fmt.Printf("unsupported platform")
	}
	cmd.Start()
}

func (w *WPLite) openOnReady(stdout io.Reader) {
	if !w.OpenOnReady {
		return
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
		if strings.Contains(m, "WordPress up and running") {
			openBrowser("http://localhost:" + strconv.Itoa(w.Env.Port))
			break
		}
	}
}

func (w *WPLite) DockerRun() error {
	l := log.WithFields(log.Fields{
		"fn": "DockerRun",
	})
	l.Debug("running wplite docker container...")
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	portMapping := strconv.Itoa(w.Env.Port) + ":80"
	cmds := []string{
		"docker", "run", "-d", "-p", portMapping,
		"--name", w.ContainerName,
		"-v", cwd + "/wp-content:/var/www/html/wp-content",
		"-v", cwd + "/.htaccess:/var/www/html/.htaccess",
		"--env-file", path.Join(cwd, w.Env.File),
		w.ImageUrl,
	}
	l.Debugf("running command: %s", strings.Join(cmds, " "))
	cmd := exec.Command(cmds[0], cmds[1:]...)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (w *WPLite) DockerStop() error {
	cmd := exec.Command("docker", "stop", w.ContainerName)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (w *WPLite) DockerRemove() error {
	cmd := exec.Command("docker", "rm", w.ContainerName)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (w *WPLite) Stop() error {
	l := log.WithFields(log.Fields{
		"fn": "Stop",
	})
	l.Info("stopping wplite...")
	if !w.DockerRunning() {
		return nil
	}
	if err := w.DockerStop(); err != nil {
		return err
	}
	if err := w.DockerRemove(); err != nil {
		return err
	}
	return nil
}

func (w *WPLite) WatchDockerLogs() error {
	cmd := exec.Command("docker", "logs", "-f", w.ContainerName)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	w.openOnReady(stdout)
	return nil
}

func (w *WPLite) Start() error {
	l := log.WithFields(log.Fields{
		"fn": "Start",
	})
	l.Info("starting wplite...")
	// if the docker container is already running, stop it first
	if w.DockerRunning() {
		if err := w.Stop(); err != nil {
			return err
		}
	}
	if err := w.EnsureEnvFile(); err != nil {
		return err
	}
	if err := w.EnsureLocalPaths(); err != nil {
		return err
	}
	if err := w.DockerRun(); err != nil {
		return err
	}
	w.WatchDockerLogs()
	return nil
}

func (w *WPLite) DockerRunning() bool {
	cmd := exec.Command("docker", "ps", "-q", "--filter", "name="+w.ContainerName)
	outdata, err := cmd.Output()
	if err != nil {
		return false
	}
	if len(outdata) == 0 {
		return false
	}
	return true
}

func (w *WPLite) streamLogsStdout() {
	cmd := exec.Command("docker", "logs", "-f", w.ContainerName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func (w *WPLite) Build(noStop, quiet bool) error {
	l := log.WithFields(log.Fields{
		"fn": "Build",
	})
	l.Info("building wplite static assets. depending on the size of your site, this may take a while...")
	// first, ensure the docker container is running
	// if not, throw an error
	if !w.DockerRunning() {
		w.OpenOnReady = false
		if err := w.Start(); err != nil {
			return err
		}
	}
	// then, trigger a static build
	var exitArg string
	if !noStop {
		exitArg = "--exit=true"
	}
	cmd := exec.Command("docker", "exec", w.ContainerName, "wp", "--allow-root", "wplite", "build", exitArg)
	if !quiet {
		go w.streamLogsStdout()
	}
	if err := cmd.Run(); err != nil {
		return err
	}
	w.Stop()
	l.Info("build complete. static assets written to wp-content/static")
	gitIgnoreStatic()
	return nil
}

func ignoreStatic() error {
	// add wp-content/static to .gitignore
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	gitIgnoreFile := path.Join(wd, ".gitignore")
	gitIgnore, err := os.OpenFile(gitIgnoreFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer gitIgnore.Close()
	newLine := "\nwp-content/static\n"
	// add the line to the file
	if _, err := gitIgnore.WriteString(newLine); err != nil {
		return err
	}
	return nil
}

func gitIgnoreStatic() {
	// if this is a git repo
	// and the wp-content/static directory is not gitignored
	// add it to the .gitignore file
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	gitDir := path.Join(wd, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return
	}
	staticDir := path.Join("wp-content", "static")
	gitIgnoreFile := path.Join(wd, ".gitignore")
	gitIgnore, err := os.Open(gitIgnoreFile)
	if err != nil {
		if err := ignoreStatic(); err != nil {
			return
		}
		return
	}
	defer gitIgnore.Close()
	scanner := bufio.NewScanner(gitIgnore)
	for scanner.Scan() {
		if scanner.Text() == staticDir {
			return
		}
	}
	if err := ignoreStatic(); err != nil {
		return
	}
}

func (w *WPLite) Pull() error {
	l := log.WithFields(log.Fields{
		"fn": "Pull",
	})
	l.Info("pulling wplite from remote git repository...")
	// if the .git directory does not exist, throw an error
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	gitDir := path.Join(wd, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) && w.VCS.GitUrl == "" {
		return fmt.Errorf("not a git repository")
	} else if os.IsNotExist(err) && w.VCS.GitUrl != "" {
		// if the .git directory does not exist, initialize a new git repository
		cmd := exec.Command("git", "init")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
		// then, add the remote origin
		cmd = exec.Command("git", "remote", "add", "origin", w.VCS.GitUrl)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	// if the wp-content/static directory is not gitignored, add it to the .gitignore file
	gitIgnoreStatic()
	// then, pull the latest changes
	cmd := exec.Command("git", "pull", "origin", w.VCS.Branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	l.Info("pull complete")
	return nil
}

func (w *WPLite) Push() error {
	l := log.WithFields(log.Fields{
		"fn": "Push",
	})
	l.Info("pushing wplite to remote git repository...")
	// if the docker container is running, stop it
	if w.DockerRunning() {
		if err := w.Stop(); err != nil {
			return err
		}
	}
	// if the .git directory does not exist, throw an error
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	gitDir := path.Join(wd, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) && w.VCS.GitUrl == "" {
		return fmt.Errorf("not a git repository")
	} else if os.IsNotExist(err) && w.VCS.GitUrl != "" {
		// if the .git directory does not exist, initialize a new git repository
		cmd := exec.Command("git", "init")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
		// then, add the remote origin
		cmd = exec.Command("git", "remote", "add", "origin", w.VCS.GitUrl)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	// if the wp-content/static directory is not gitignored, add it to the .gitignore file
	gitIgnoreStatic()
	// then, add, commit, and push the changes
	cmd := exec.Command("git", "add", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command("git", "commit", "-m", "wplite update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command("git", "push", "--set-upstream", "origin", w.VCS.Branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	l.Info("push complete")
	return nil
}
