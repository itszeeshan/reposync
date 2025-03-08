# RepoSync

**RepoSync** is a powerfool command-line tool written in `golang` for cloning all repositories from a GitLab group or GitHub organization, including nested subgroups (GitLab) and all repositories (GitHub).

## Features

- Clone all repositories from a **GitLab** group or a **GitHub** organization.
- Recursively clone GitLab/Github subgroups.
- Supports **HTTPS** and **SSH** cloning.
- Automatically skips existing repositories.
- Simple and efficient with minimal dependencies.

---

## Workflow

![reposync-workflow](https://mermaid.ink/img/pako:eNpdklGTmjAQx7_KTp45R4Qq0pnO9EDFO725Kfal6EMKe5IRCBOCowW_e2PitefxlP_u77-7CduRlGdIfLIXtM5hE37dVqC-70ksqZA7eHj41v9sUIBoqwYE1rw5V2kPj12MBaYSXgU_sgzFxRgftWPBZNT-7iFI5ijTHIyGH1f77hO4ogoM_4NK34HBFYRZFxS8QlijzHl2axaalBEzXTDabF7jHuaJwY-Mgg7tPkJxHPWw-ICowA2Y65KRuj4XaOYAVkHIhLosF-cbtjCYEZEWyyTIMT3A8k3bYHZijWxgxVNaFO--pe5vUj08JfGB1XCdg1X7OyTk2MALl6ZMD8_J7IRpK_H6RNqBN_5Zd18lMT2agT-1fNL5dRLwsi5QMl6pR2waun8vsDKAEWstXpJZlak0sUiJoqQsUyvSXYktkTmWuCW-OmZUHLZkW10UR1vJY7UaxJeiRYsI3u5z4r_RolGqrTMqMWRU7Vn5L1rT6hfnd5r4HTkR3x3Zg6HjOo43HI1tx_Mscia-Yw88z_bc4Rdn4k7c8fRikT-6wHCgAqOp50yd8ciZ2OOxRTBj6o-tzX7rNb_8BUdT5Og?type=png)

## Installation

### Install Golang (If not already installed)

- Ensure you have **Git** [(Download Git)](https://github.com/git-guides/install-git):

```sh
# macOS (Homebrew)
brew install git

# Ubuntu/Debian
sudo apt-get update && sudo apt-get install git-all
```

Once Installed run the following command to check on your terminal `git version`.

- Ensure you have **Go** installed. If not, install it using[(Download Go)](https://go.dev/doc/install):

```sh
# macOS (Homebrew)
brew install go

# Ubuntu/Debian
sudo apt update && sudo apt install golang-go
```

Once Installed run the following command to check on your terminal `go version`.

### Install RepoSync

```sh
go install github.com/itszeeshan/reposync@latest
```

This will download, compile, and install `reposync` to `$GOPATH/bin`.

### (Optional) Make it globally accessible

Move the binary to `/usr/local/bin` so it can be used system-wide:

```sh
sudo mv $(go env GOPATH)/bin/reposync /usr/local/bin/
```

Now you can run `reposync` from anywhere. 🚀

## Usage

Run RepoSync from the terminal:

```sh
reposync -p <gitlab|github> -t <PERSONAL_ACCESS_TOKEN> -g <GROUP_ID|ORG_ID> [-m <https|ssh>]
```

### Arguments

| Argument | Description                                     |
| -------- | ----------------------------------------------- |
| `-p`     | Provider: `gitlab` or `github`                  |
| `-t`     | Personal Access Token (GitHub/GitLab)           |
| `-g`     | Group ID (GitLab) or Organization Name (GitHub) |
| `-m`     | Clone method: `https` (default) or `ssh`        |

### Examples

#### Clone all repositories from a GitHub organization/Gitlab groups using HTTPS

```sh
reposync -p github -t ghp_yourtoken -g your-org
```

#### Clone all repositories from a Github organization/GitLab group using SSH

```sh
reposync -p gitlab -t glpat_yourtoken -g 123456 -m ssh
```

---

## Troubleshooting

### Authentication Issues

- Ensure you have **correct permissions** for the provided access token.
- If using **SSH**, ensure your SSH key is added to your account.

### Cloning Errors

- Ensure `git` is installed: `git --version`
- Ensure you have **network access** to GitHub/GitLab.

### Permission Denied Issues

- Ensure your token has sufficient permissions:
- **GitHub**: Token must have `repo` scope.
- **GitLab**: Token must have `read_api` scope.

### Command Not Found

- Ensure `reposync` is in your `$PATH`. If not, run:

```sh
$ export PATH=$PATH:/usr/local/bin/
```

---

## Contributing

Pull requests are welcome! If you encounter issues, feel free to open an issue on GitHub.

---

## License

MIT License. See `LICENSE` for details.
