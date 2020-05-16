# krypt
[![Software License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/gesquive/krypt/blob/master/LICENSE)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/gesquive/krypt)
[![Build Status](https://img.shields.io/circleci/build/github/gesquive/krypt?style=flat-square)](https://circleci.com/gh/gesquive/krypt)
[![Coverage Report](https://img.shields.io/codecov/c/gh/gesquive/krypt?style=flat-square)](https://codecov.io/gh/gesquive/krypt)

A command line file encrypter and decrypter.

This program started as a clone of ansible-vault. Support for additional ciphers and a different binary only file format was added for flexibility.

## Installing

### Compile
This project has only been tested with go1.14+. To compile just run `go get -u github.com/gesquive/krypt` and the executable should be built for you automatically in your `$GOPATH`. This project uses go mods, so you might need to set `GO111MODULE=on` in order for `go get` to complete properly.

Optionally you can run `make install` to build and copy the executable to `/usr/local/bin/` with correct permissions.

### Download
Alternately, you can download the latest release for your platform from [github](https://github.com/gesquive/krypt/releases).

Once you have an executable, make sure to copy it somewhere on your path like `/usr/local/bin` or `C:/Program Files/`.
If on a \*nix/mac system, make sure to run `chmod +x /path/to/krypt`.

### Homebrew
This app is also avalable from this [homebrew tap](https://github.com/gesquive/homebrew-tap). Just install the tap and then the app will be available.
```shell
$ brew tap gesquive/tap
$ brew install krypt
```

## Configuration

### Precedence Order
The application looks for variables in the following order:
 - command line flag
 - environment variable
 - config file variable
 - default

So any variable specified on the command line would override values set in the environment or config file.

### Config File
The application looks for a configuration file at the following locations in order:
 - `./config.yml`
 - `~/.config/krypt/config.yml`
 - `/etc/krypt/config.yml`

### Environment Variables
Optionally, instead of using a config file you can specify config entries as environment variables. Use the prefix "KRYPT_" in front of the uppercased variable name. For example, the config variable `password-file` would be the environment variable `KRYPT_PASSWORD_FILE`.

## Usage

```console
Encrypt or Decrypt files using different ciphers

Usage:
  krypt [flags] command

Available Commands:
  create      Create a new encrypted text file
  edit        Decrypt, edit and encrypt an encrypted file
  help        Help about any command
  list        List the available cipher methods
  reseal      Change the password/cipher on encrypted file(s)
  seal        Seal unencrypted file(s)
  unseal      Unseal encrypted file(s)
  view        Decrypt and view the contents of a sealed file without editing

Flags:
  -c, --config string   Path to a specific config file (default "./config.yml")
  -v, --version         Show the version info and exit
```

Optionally, a hidden debug flag is available in case you need additional output.
```console
Hidden Flags:
  -D, --debug                  Include debug statements in log output
```

## Documentation

This documentation can be found at github.com/gesquive/krypt

## License

This package is made available under an MIT-style license. See LICENSE.

## Contributing

PRs are always welcome!
