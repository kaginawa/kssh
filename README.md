kssh
====

[![Actions Status](https://github.com/kaginawa/kssh/workflows/Go/badge.svg)](https://github.com/kaginawa/kssh/actions)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=kaginawa_kssh&metric=alert_status)](https://sonarcloud.io/dashboard?id=kaginawa_kssh)
[![Go Report Card](https://goreportcard.com/badge/github.com/kaginawa/kssh)](https://goreportcard.com/report/github.com/kaginawa/kssh)

[Kaginawa](https://github.com/kaginawa/kaginawa) powered SSH client.

## Usage

Login target using custom ID with current user name:

```
kssh <CUSTOM_ID>
# Example: kssh debug1
```

Login target using custom ID with specify user name:

```
kssh <USER>@<CUSTOM_ID>
# Example: alice@debug1
```

Login target using MAC address with current user name:

```
kssh <MAC>
# Example: kssh f0:18:98:eb:c7:27
```

Login target using MAC address with specify user name:

```
kssh <USER>@<MAC>
# Example: kssh alice@f0:18:98:eb:c7:27
```

## Options

- `-k <API_KEY>` - specify admin API key
- `-c <CONFIG>` - specify config file path
- `-s <SERVER>` - specify [kaginawa server](https://github.com/kaginawa/kaginawa-server) address

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
| Server          |         | Host name or IP address of the kaginawa-server (e.g. 10.128.1.100) |
| AdminKey        |         | API key issued at kaginawa-server |
| DefaultUser     | $USER   | Default login user |
| DefaultPassword |         | Default password for login user (WARNING: understand security risks) |

## Install

Via go get:

```
go get -u https://github.com/kaginawa/kssh
```

## License

kssh licenced under [BSD 3-Clause](LICENSE).

## Author

[mikan](https://github.com/mikan)
