kssh
====

[![Actions Status](https://github.com/kaginawa/kssh/workflows/Go/badge.svg)](https://github.com/kaginawa/kssh/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/kaginawa/kssh)](https://goreportcard.com/report/github.com/kaginawa/kssh)
s://github.com/kagin
[Kaginawa](httpawa/kaginawa) powered SSH client.

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

Default configuration file is located at `~/.ssh/kconfig`.
It generates automatically when first use.

Format:

```
AdminKey <API_KEY>
Server <SERVER>
```

## Install

Via go get:

```
go get -u https://github.com/kaginawa/kssh
```

## License

kssh licenced under [BSD 3-Clause](LICENSE).

## Author

[mikan](https://github.com/mikan)
