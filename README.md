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

![reposync-workflow](https://mermaid.ink/img/pako:eNqNkl1vgjAUx7_KU546H4RQdXMPlOj25p6MezH6kMJeZAQSIKMFv3sj47bn8ZT_7u-_uwm7kZJnSHyyF7TOYRN-3Vagvu9JLKmQO3h4-Nb_bFCAaKsGBNa8OVdpD49djAWmEl4FP7IMxcUYH7VjwWTU_u4hSOYo0xyMhh9X--4TuKIKDP-DSt-BwRWEWRcUvEJYo8x5dmsWmpQRM10w2mxu8R7micGPjIIO7T5CcRz1sPiAqMANmOuSkbo-F2jmAFZByIW6LBfnG7YwmBGRFsskyDE9wPJN22B2Yo1sYMVjWhTvvqXub1I9PCXxgVXwMgereu4hIccGXrg0ZXp4TmYnTFuJ1yfSDrzxz7r7Konp0Qz8qeWTzq-TgJd1gZLxSj1i09D9e4GVAQyutfhMZlWm0sQiJYqSskytSHcltkTmWOKW-OqYUXHYkm11URxtJY_VahBfihYtIni7z4n_RotGqbbOqMSQUbVn5b9oTatfnN9p4nfkRHx3ZA-Gjus43nA0th3Ps8iZ-I498Dzbc4dfnIk7ccfTi0X-6ALDgQqMpp4zdcYjZ2KPxxbBjKk_tjZf2-1vfwHrVZGF?type=png)

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
