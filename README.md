# RepoSync

**RepoSync** is a powerful command-line tool written in `Go` for securely cloning all repositories from a GitLab group or GitHub organization, including nested subgroups (GitLab) with zero-touch authentication configuration.

## Features

- **Secure token storage** - Credentials saved in user home directory
- Clone all repositories from **GitLab groups** (with subgroups) or **GitHub organizations**
- **HTTPS/SSH cloning** support
- Automatic directory structure mirroring
- Configuration wizard for easy setup
- Skipping of existing repositories

---

## Workflow

![reposync-workflow](https://mermaid.ink/img/pako:eNpdklGTmjAQx7_KTp45R4Qq0pnO9EDFO725Kfal6EMKe5IRCBOCowW_e2PitefxlP_u77-7CduRlGdIfLIXtM5hE37dVqC-70ksqZA7eHj41v9sUIBoqwYE1rw5V2kPj12MBaYSXgU_sgzFxRgftWPBZNT-7iFI5ijTHIyGH1f77hO4ogoM_4NK34HBFYRZFxS8QlijzHl2axaalBEzXTDabF7jHuaJwY-Mgg7tPkJxHPWw-ICowA2Y65KRuj4XaOYAVkHIhLosF-cbtjCYEZEWyyTIMT3A8k3bYHZijWxgxVNaFO--pe5vUj08JfGB1XCdg1X7OyTk2MALl6ZMD8_J7IRpK_H6RNqBN_5Zd18lMT2agT-1fNL5dRLwsi5QMl6pR2waun8vsDKAEWstXpJZlak0sUiJoqQsUyvSXYktkTmWuCW-OmZUHLZkW10UR1vJY7UaxJeiRYsI3u5z4r_RolGqrTMqMWRU7Vn5L1rT6hfnd5r4HTkR3x3Zg6HjOo43HI1tx_Mscia-Yw88z_bc4Rdn4k7c8fRikT-6wHCgAqOp50yd8ciZ2OOxRTBj6o-tzX7rNb_8BUdT5Og?type=png)

## Installation

### Prerequisites

- **Git** [(Download Git)](https://github.com/git-guides/install-git)
- **Go 1.20+** [(Download Go)](https://go.dev/doc/install)

### Install RepoSync

```sh
go install github.com/itszeeshan/reposync@latest
```

### (Optional) Make it globally accessible

Move the binary to `/usr/local/bin` so it can be used system-wide:

```sh
sudo mv $(go env GOPATH)/bin/reposync /usr/local/bin/
```

Now you can run `reposync` from anywhere. 🚀

### Configuration

1. **Configure your credentials** (first-time setup):

```sh
reposync config
```

Follow the prompts to enter your GitLab and GitHub personal access tokens

2. **Verify configuration**:

```sh
cat ~/.reposync/config.json
```

## Usage

```sh
# Basic syntax
reposync -p <gitlab|github> -g <GROUP_ID|ORG_NAME> [-m <https|ssh>]

# Examples
reposync -p gitlab -g 123456          # Clone GitLab group with HTTPS
reposync -p github -g your-org -m ssh # Clone GitHub org with SSH
```

### Arguments

| Argument | Description                                   |
| -------- | --------------------------------------------- |
| `-p`     | Provider: `gitlab` or `github` (required)     |
| `-g`     | Group ID (GitLab) or Organization ID (GitHub) |
| `-m`     | Clone method: `https` (default) or `ssh`      |

### Examples

#### Clone all repositories from a GitHub organization/Gitlab groups using HTTPS

```sh
reposync -p github -g your-org
```

#### Clone all repositories from a Github organization/GitLab group using SSH

```sh
reposync -p gitlab -t glpat_yourtoken -g 123456 -m ssh
```

---

## Troubleshooting

### Authentication Issues

```sh
# If receiving permission errors:
reposync config # Update stored credentials
```

### Token Management

- Rotate tokens: Re-run `reposync config`

- Delete credentials: rm `~/.reposync/config.json`

### Cloning Errors

- Ensure `git` is installed: `git --version`
- Ensure you have **network access** to GitHub/GitLab.

### Common Errors

1. Ensure tokens have required scopes:

- **GitHub**: `repo` scope

- **GitLab**: `read_api` scope

2. SSH cloning requires:

- SSH key added to your account

- `ssh-agent` running

---

## Contributing

Pull requests are welcome! If you encounter issues, feel free to open an issue on GitHub.

---

## License

MIT License. See `LICENSE` for details.
