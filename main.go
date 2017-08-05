package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

var build = "0" // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "git plugin"
	app.Usage = "git plugin"
	app.Action = run
	app.Version = fmt.Sprintf("1.1.%s", build)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "remote",
			Usage:  "git remote url",
			EnvVar: "DRONE_REMOTE_URL",
		},
		cli.StringFlag{
			Name:   "path",
			Usage:  "git clone path",
			EnvVar: "DRONE_WORKSPACE",
		},
		cli.StringFlag{
			Name:   "sha",
			Usage:  "git commit sha",
			EnvVar: "DRONE_COMMIT_SHA",
		},
		cli.StringFlag{
			Name:   "ref",
			Value:  "refs/heads/master",
			Usage:  "git commit ref",
			EnvVar: "DRONE_COMMIT_REF",
		},
		cli.StringFlag{
			Name:   "event",
			Value:  "push",
			Usage:  "build event",
			EnvVar: "DRONE_BUILD_EVENT",
		},
		cli.StringFlag{
			Name:   "netrc.machine",
			Usage:  "netrc machine",
			EnvVar: "DRONE_NETRC_MACHINE",
		},
		cli.StringFlag{
			Name:   "netrc.username",
			Usage:  "netrc username",
			EnvVar: "DRONE_NETRC_USERNAME",
		},
		cli.StringFlag{
			Name:   "netrc.password",
			Usage:  "netrc password",
			EnvVar: "DRONE_NETRC_PASSWORD",
		},
		cli.IntFlag{
			Name:   "depth",
			Usage:  "clone depth",
			EnvVar: "PLUGIN_DEPTH",
		},
		cli.BoolTFlag{
			Name:   "recursive",
			Usage:  "clone submodules",
			EnvVar: "PLUGIN_RECURSIVE",
		},
		cli.BoolFlag{
			Name:   "tags",
			Usage:  "clone tags",
			EnvVar: "PLUGIN_TAGS",
		},
		cli.BoolFlag{
			Name:   "skip-verify",
			Usage:  "skip tls verification",
			EnvVar: "PLUGIN_SKIP_VERIFY",
		},
		cli.BoolFlag{
			Name:   "submodule-update-remote",
			Usage:  "update remote submodules",
			EnvVar: "PLUGIN_SUBMODULES_UPDATE_REMOTE,PLUGIN_SUBMODULE_UPDATE_REMOTE",
		},
		cli.GenericFlag{
			Name:   "submodule-override",
			Usage:  "json map of submodule overrides",
			EnvVar: "PLUGIN_SUBMODULE_OVERRIDE",
			Value:  &MapFlag{},
		},
		cli.StringFlag{
			Name:  "env-file",
			Usage: "source env file",
		},
		// ------------------------------------------------------
		// 自定义 drone-git 参数
		cli.StringFlag{
			Name:   "commit.message",
			Usage:  "commit message",
			EnvVar: "DRONE_COMMIT_MESSAGE",
		},
		cli.StringFlag{
			Name:   "commit.branch",
			Usage:  "commit branch",
			EnvVar: "DRONE_COMMIT_BRANCH",
		},
		cli.StringFlag{
			Name:   "pr.number",
			Usage:  "commit branch",
			EnvVar: "DRONE_PULL_REQUEST",
		},
		cli.StringFlag{
			Name:   "custom.remote.url",
			Usage:  "custom git remote url",
			EnvVar: "PLUGIN_REMOTE_URL",
		},
		cli.StringFlag{
			Name:   "branch",
			Usage:  "repo branch",
			EnvVar: "PLUGIN_BRANCH",
		},
		cli.StringFlag{
			Name:   "custom.path",
			Usage:  "custom git clone path",
			EnvVar: "PLUGIN_PATH",
		},
		cli.StringFlag{
			Name:   "debug",
			Usage:  "show debug info",
			EnvVar: "PLUGIN_DEBUG",
		},
		// ------------------------------------------------------
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	// 如果存在custom.remote.url, 则认为是要下载依赖包，比如base库等, 否则使用drone传入remote
	remote := c.String("remote")
	isDependencyRepo := false
	if c.String("custom.remote.url") != "" {
		remote = c.String("custom.remote.url")
		isDependencyRepo = true
	}

	// 如果存在branch, 则替换ref为refs/heads/$branch, 否则使用drone传入ref
	refs := c.String("ref")
	if c.String("branch") != "" {
		refs = fmt.Sprintf("refs/heads/%s", c.String("branch"))
	}

	// 如果存在custom.path, 则替换clone repo path, 否则使用drone传入workspace
	// drone默认为: /drone/src/github.com/octocat/hello-world
	path := c.String("path")
	if c.String("custom.path") != "" {
		path = c.String("custom.path")
	}

	plugin := Plugin{
		Repo: Repo{
			Clone: remote,
		},
		Build: Build{
			Commit:        c.String("sha"),
			Event:         c.String("event"),
			Path:          path,
			Ref:           refs,
			CommitMessage: c.String("commit.message"),
			Branch:        c.String("commit.branch"),
			PullReqNumber: c.String("pr.number"),
		},
		Netrc: Netrc{
			Login:    c.String("netrc.username"),
			Machine:  c.String("netrc.machine"),
			Password: c.String("netrc.password"),
		},
		Config: Config{
			Depth:            c.Int("depth"),
			Tags:             c.Bool("tags"),
			Recursive:        c.BoolT("recursive"),
			SkipVerify:       c.Bool("skip-verify"),
			SubmoduleRemote:  c.Bool("submodule-update-remote"),
			Submodules:       c.Generic("submodule-override").(*MapFlag).Get(),
			IsDependencyRepo: isDependencyRepo,
		},
	}

	// 显示环境变量和Plugin详细信息
	if c.Bool("debug") {
		logrus.Infof("plugin: %+v", plugin)
		for _, e := range os.Environ() {
			logrus.Info(e)
		}
	}

	return plugin.Exec()
}
