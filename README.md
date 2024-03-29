kssh
====

[![Actions Status](https://github.com/kaginawa/kssh/workflows/Go/badge.svg)](https://github.com/kaginawa/kssh/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/kaginawa/kssh)](https://goreportcard.com/report/github.com/kaginawa/kssh)

[Kaginawa](https://github.com/kaginawa/kaginawa)-powered SSH client.

## Download

See [Releases](https://github.com/kaginawa/kssh/releases) page.

## Usage

Login using custom ID as the current user:

```
kssh <CUSTOM_ID>
# Example: kssh debug1
```

Login using custom ID as a different user:

```
kssh <USER>@<CUSTOM_ID>
# Example: alice@debug1
```

Login using MAC address as the current user:

```
kssh <MAC>
# Example: kssh f0:18:98:eb:c7:27
```

Login using MAC address as a different user:

```
kssh <USER>@<MAC>
# Example: kssh alice@f0:18:98:eb:c7:27
```

Login, run and exit:

```
kssh [<USER>@]<CUSTOM_ID|MAC> <COMMAND>
# Example: kssh debug1 uname -a
# Do not specify interactive commands (e.g. vi)
```

## Options

- `-k <API_KEY>` - specify admin API key
- `-c <CONFIG>` - specify config file path
- `-s <SERVER>` - specify [kaginawa server](https://github.com/kaginawa/kaginawa-server) address
- `-f <PROCESURE_FILE>` - specify procedure (line-separated list of commands) file
- `-m <MINUTES>` - specify freshness threshold by minutes (default = 15)
- `-l` - listen a local port for transferring non-SSH TCP connections trough the SSH tunnel

## Configuration

Default file name of the configuration file is `kssh.conf` and location is [platform-dependent](https://golang.org/pkg/os/#UserConfigDir).

- Linux: `~/config/kssh.conf`
- macOS: `~/Library/Application Support/kssh.conf`
- Windows: `%AppData%\kssh.conf`

Format:

```
AdminKey <API_KEY>
Server <SERVER>
```

Supported parameters:

| Key             | Default | Description |
| --------------- | ------- | ----------- |
| Server          |         | Host name or IP address of the kaginawa-server (e.g. http://10.128.1.100) |
| AdminKey        |         | API key issued at kaginawa-server |
| DefaultUser     | $USER   | Default login user |
| DefaultPassword |         | Default password for login user (WARNING: understand security risks) |

## License

kssh licenced under [BSD 3-Clause](LICENSE).

## Author

[mikan](https://github.com/mikan)
