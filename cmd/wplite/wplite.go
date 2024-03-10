package main

import (
	"flag"
	"os"

	"github.com/robertlestak/wplite/internal/wplite"
	log "github.com/sirupsen/logrus"
)

var (
	wpliteFlagset = flag.NewFlagSet("wplite", flag.ExitOnError)
)

func init() {
	ll, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		ll = log.InfoLevel
	}
	log.SetLevel(ll)
}

func cmdDockerRun(action string) {
	l := log.WithFields(log.Fields{
		"fn":     "cmdDockerRun",
		"action": action,
	})
	l.Debug("running wplite docker container")
	startFlagset := flag.NewFlagSet("start", flag.ExitOnError)
	title := startFlagset.String("title", "WPLite", "site title")
	user := startFlagset.String("user", "admin", "admin username")
	pass := startFlagset.String("pass", "admin", "admin password")
	email := startFlagset.String("email", "hello@example.com", "admin email")
	openOnReady := startFlagset.Bool("open", true, "open browser on ready")
	theme := startFlagset.String("theme", "twentytwentyfour", "theme")
	port := startFlagset.Int("port", 80, "port")
	envFile := startFlagset.String("env-file", ".wplite-env", "env file")
	wpliteImage := startFlagset.String("image", "robertlestak/wplite:latest", "wplite image")
	noStop := startFlagset.Bool("no-stop", false, "do not stop container after running")
	quietBuild := startFlagset.Bool("quiet", false, "do not output build logs")
	branch := startFlagset.String("git-branch", "main", "git branch")
	gitUrl := startFlagset.String("git-url", "", "git url")
	if err := startFlagset.Parse(os.Args[2:]); err != nil {
		l.WithError(err).Error("error parsing start flagset")
	}
	wpl := wplite.WPLite{
		ImageUrl:    *wpliteImage,
		OpenOnReady: *openOnReady,
		Env: wplite.WPLiteEnv{
			File:  *envFile,
			Theme: *theme,
			Title: *title,
			User:  *user,
			Pass:  *pass,
			Email: *email,
			Port:  *port,
		},
		VCS: wplite.VCS{
			Branch: *branch,
			GitUrl: *gitUrl,
		},
	}
	wpl.ContainerName = wpl.WorkspaceContainerName()
	switch action {
	case "build":
		if err := wpl.Build(*noStop, *quietBuild); err != nil {
			l.WithError(err).Error("error building wplite")
		}
	case "pull":
		if err := wpl.Pull(); err != nil {
			l.WithError(err).Error("error pulling wplite")
		}
	case "push":
		if err := wpl.Push(); err != nil {
			l.WithError(err).Error("error pushing wplite")
		}
	case "start":
		if err := wpl.Start(); err != nil {
			l.WithError(err).Error("error starting wplite")
		}
	case "stop":
		if err := wpl.Stop(); err != nil {
			l.WithError(err).Error("error stopping wplite")
		}
	}
}

func main() {
	l := log.WithFields(log.Fields{
		"fn": "main",
	})
	l.Debug("starting wplite")
	logLevel := wpliteFlagset.String("log-level", log.GetLevel().String(), "log level")
	wpliteFlagset.Parse(os.Args[1:])
	ll, err := log.ParseLevel(*logLevel)
	if err != nil {
		ll = log.InfoLevel
	}
	log.SetLevel(ll)
	if len(wpliteFlagset.Args()) == 0 {
		l.Debug("no command specified")
		return
	}
	switch wpliteFlagset.Args()[0] {
	case "build":
		cmdDockerRun("build")
	case "push":
		cmdDockerRun("push")
	case "start":
		cmdDockerRun("start")
	case "stop":
		cmdDockerRun("stop")
	default:
		l.Debug("unknown command")
	}
}
