# [![CircleCI Build Status](https://circleci.com/gh/raravena80/ya.svg?style=shield)](https://circleci.com/gh/raravena80/ya) [![Coverage Status](https://coveralls.io/repos/github/raravena80/ya/badge.svg?branch=master)](https://coveralls.io/github/raravena80/ya?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/raravena80/ya)](https://goreportcard.com/report/github.com/raravena80/ya) [![Apache Licensed](https://img.shields.io/badge/license-Apache2.0-blue.svg)](https://raw.githubusercontent.com/raravena80/ya/master/LICENSE)

![ya](https://user-images.githubusercontent.com/7659560/32466351-7f0fec64-c2fb-11e7-8299-1aad4fbdd28e.png)


Ya runs or copies items across multiple servers using SSH or SCP

## Usage
```
Ya runs commands or copies files or directories,
across multiple servers, using SSH or SCP

Usage:
  ya [command]

Available Commands:
  help        Help about any command
  scp         Copy files to multiple servers
  ssh         Run command acrosss multiple servers

Flags:
      --config string      config file (default is $HOME/.ya.yaml)
  -h, --help               help for ya
  -k, --key string         Ssh key to use for authentication, full path (default "/Users/raravena/.ssh/id_rsa")
  -m, --machines strings   Hosts to run command on
  -p, --port int           Ssh port to connect to (default 22)
  -a, --useagent           Use agent for authentication
  -u, --user string        User to run the command as (default "raravena")

Use "ya [command] --help" for more information about a command.
```

## Config

Sample `~/.ya.yaml`

```
ya:
  user: ubuntu
  key: /Users/username/.ssh/id_rsa
  useagent: true
  machines:
    - 172.1.1.1
    - 172.1.1.2
    - 172.1.1.3
    - 172.1.1.4
    - 172.1.1.5
  command: sudo rm -f /var/log/syslog.*
```
