# ![ya](https://user-images.githubusercontent.com/7659560/32466351-7f0fec64-c2fb-11e7-8299-1aad4fbdd28e.png)

[![CircleCI Build Status](https://circleci.com/gh/raravena80/ya.svg?style=shield)](https://circleci.com/gh/raravena80/ya) [![Coverage Status](https://coveralls.io/repos/github/raravena80/ya/badge.svg?branch=master)](https://coveralls.io/github/raravena80/ya?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/raravena80/ya)](https://goreportcard.com/report/github.com/raravena80/ya) [![Documentation](https://godoc.org/github.com/raravena80/ya?status.svg)](http://godoc.org/github.com/raravena80/ya)  [![Apache Licensed](https://img.shields.io/badge/license-Apache2.0-blue.svg)](https://raw.githubusercontent.com/raravena80/ya/master/LICENSE)


Ya runs a command using SSH or copies items using SCP across multiple servers.

## Usage
```
Ya runs a command or copies files or directories,
across multiple servers, using SSH or SCP

Usage:
  ya [command]

Available Commands:
  help        Help about any command
  scp         Copy files to multiple servers
  ssh         Run command acrosss multiple servers

Flags:
  -s, --agentsock string   SSH agent socket file. If using SSH agent (default "/private/tmp/com.apple.launchd.67UG0GmO3V/Listeners")
      --config string      config file (default is $HOME/.ya.yaml)
  -h, --help               help for ya
  -k, --key string         Ssh key to use for authentication, full path (default "/Users/raravena/.ssh/id_rsa")
  -m, --machines strings   Hosts to run command on
  -p, --port int           Ssh port to connect to (default 22)
  -t, --timeout int        Timeout for connection (default 5)
  -a, --useagent           Use agent for authentication
  -u, --user string        User to run the command as (default "raravena")
  -v, --verbose            Set verbose output

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
```
SSH Specific
```
  ssh:
    command: sudo rm -f /var/log/syslog.*
```
SCP Specific
```
  scp:
    source: /Users/raravena/tmp2
    destination: /root/tmpboguss
    recursive: true
```

## SSH Examples

Makes /tmp/tmpdir in 17.2.2.2 and 17.2.3.2:
```
$ ya ssh -c "mkdir /tmp/tmpdir" -m 17.2.2.2,17.2.3.2
```

Creates /tmp/file file in host1 and host2
```
$ ya ssh -c "touch /tmp/file" -m host1,host2
```

Moves a file and creates another one in host1 and host2:
```
$ ya ssh -c "mv /tmp/file1 /tmp/file2; touch /tmp/file3" -m host1,host2
```

Runs with default in `~/.ya.yaml`
```
$ ya ssh
```

## SCP Examples

Copies from local /tmp/tmpfile to /tmp/tmpfile2 in 17.2.2.2 and 17.2.3.2:
```
$ ya scp --src /tmp/tmpfile --dst /tmp/tmpfile2  -m 17.2.2.2,17.2.3.2
```

Copies /tmp/dir directory under /tmp in host1 and host2 (recursively)
```
$ ya scp  -r --src /tmp/dir  --dst /tmp host1,host2
```

Runs with default in `~/.ya.yaml`
```
$ ya scp
```
