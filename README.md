# RepoSync

RepoSync is a powerful command-line tool written in Go for securely cloning all repositories from GitLab groups or GitHub organizations, including nested subgroups (GitLab) with zero-touch authentication configuration.

## Features

- **Secure token storage** - Credentials saved in user home directory with proper permissions
- **Complete repository cloning** - Clone all repositories from GitLab groups (with subgroups) or GitHub organizations
- **Flexible cloning methods** - Support for both HTTPS and SSH cloning
- **Directory structure preservation** - Automatic mirroring of organization hierarchy
- **Smart skipping** - Automatically skips already cloned repositories
- **Progress reporting** - Real-time progress indicators during cloning
- **Enterprise support** - Works with self-hosted GitLab and GitHub Enterprise
- **Input validation** - Comprehensive validation for all inputs
- **Error handling** - Robust error handling with retry mechanisms
- **Rate limiting** - Built-in rate limiting to prevent API throttling

## Installation

### Prerequisites

- **Git** [(Download Git)](https://github.com/git-guides/install-git)
- **Go 1.24+** [(Download Go)](https://go.dev/doc/install)

### Install RepoSync

```sh
go install github.com/itszeeshan/reposync@latest
```

### Make it globally accessible (Optional)

Move the binary to `/usr/local/bin` so it can be used system-wide:

```sh
sudo mv $(go env GOPATH)/bin/reposync /usr/local/bin/
```

Now you can run `reposync` from anywhere.

## Configuration

### Initial Setup

1. **Configure your credentials** (first-time setup):

```sh
reposync config
```

Follow the prompts to enter your GitLab and GitHub personal access tokens. Tokens are entered securely and hidden from terminal history.

2. **Verify configuration**:

```sh
cat ~/.reposync/config.json
```

### Token Requirements

- **GitHub**: Personal access token with `repo` scope
- **GitLab**: Personal access token with `read_api` scope

## Usage

### Basic Syntax

```sh
reposync -p <gitlab|github> -g <GROUP_ID|ORG_NAME> [-m <https|ssh>]
```

### Arguments

| Argument | Description                                     | Required |
| -------- | ----------------------------------------------- | -------- |
| `-p`     | Provider: `gitlab` or `github`                  | Yes      |
| `-g`     | Group ID (GitLab) or Organization name (GitHub) | Yes      |
| `-m`     | Clone method: `https` (default) or `ssh`        | No       |
| `-h`     | Show help message                               | No       |

### Examples

#### Clone GitHub organization with HTTPS

```sh
reposync -p github -g your-organization
```

#### Clone GitLab group with SSH

```sh
reposync -p gitlab -g 123456 -m ssh
```

#### Clone GitHub organization with SSH

```sh
reposync -p github -g your-organization -m ssh
```

## Directory Structure

### GitLab Group Structure

RepoSync preserves the exact GitLab group hierarchy:

```text
my-group/                   # Parent group
├── backend/                # Subgroup
│   ├── auth-service/       # Repository
│   └── payment-service/
├── frontend/
│   ├── web-app/
│   └── mobile-app/
└── tools/                  # Root group repositories
    ├── ci-cd-templates/
    └── monitoring-system/
```

### GitHub Organization Structure

GitHub repositories are cloned in a flat structure:

```text
my-organization/           # Organization root
├── docs-website/          # Repository
├── api-gateway/
├── user-management/
└── infrastructure-as-code/
```

## Use Cases

### Local Development Mirroring

Create exact replicas of your organization structure for local development:

```text
my-group/
└── frontend/
    └── web-app/  # Matches GitLab's "my-group > frontend > web-app"
```

**Benefits:**

- Instant context with paths identical to web interface
- Multi-repo workflows with predictable paths
- Easy navigation between related repositories

### Backup and Archiving

Maintain complete organizational structure for backups:

```text
backup-2024/
├── my-group/                  # Full hierarchy preserved
│   └── infrastructure/
└── my-organization/           # GitHub repos in flat structure
```

**Benefits:**

- Disaster recovery with exact group/repo relationships
- Version snapshots for compliance
- Audit trail maintenance

### CI/CD Pipeline Integration

Use consistent paths for automation scripts:

```bash
# Scan all backend services
for repo in my-group/backend/*; do
  trivy config "$repo"
done
```

**Benefits:**

- Predictable source paths for Docker/K8s builds
- Path-based automation scripts
- Consistent deployment workflows

### IDE Workspace Configuration

Configure multi-root workspaces with persistent paths:

```json
{
  "folders": [
    { "path": "my-group/frontend/web-app" },
    { "path": "my-group/backend/auth-service" }
  ]
}
```

**Benefits:**

- Persistent workspace configurations
- Cross-repo navigation in IDEs
- Single-window multi-repo development

## Troubleshooting

### Authentication Issues

If you receive permission errors:

```sh
reposync config  # Update stored credentials
```

### Token Management

- **Rotate tokens**: Re-run `reposync config`
- **Delete credentials**: `rm ~/.reposync/config.json`
- **Check token scopes**: Ensure tokens have required permissions

### Common Issues

1. **Permission denied errors**:

   - Verify token has correct scopes
   - Check if token is expired
   - Ensure you have access to the group/organization

2. **SSH cloning issues**:

   - Ensure SSH key is added to your account
   - Verify `ssh-agent` is running
   - Test SSH connection manually

3. **Network connectivity**:

   - Ensure you have network access to GitHub/GitLab
   - Check firewall/proxy settings
   - Verify DNS resolution

4. **Rate limiting**:
   - Tool includes built-in rate limiting
   - Wait and retry if you hit API limits
   - Consider using SSH for large organizations

### Error Messages

The tool provides specific error messages to help diagnose issues:

- **Invalid organization name**: Check format (alphanumeric and hyphens only)
- **Invalid group ID**: Ensure it's a valid integer
- **Token validation failed**: Token too short or empty
- **Repository cloning failed**: Network issues or permission problems

## Advanced Features

### Self-Hosted Instances

RepoSync supports self-hosted GitLab and GitHub Enterprise instances. Configuration can be extended to include custom URLs in the config file.

### Progress Reporting

Real-time progress indicators show:

- Number of repositories found
- Current cloning progress
- Success/failure status for each repository

### Retry Logic

Built-in retry mechanism for git clone operations:

- Automatic retry on network failures
- Exponential backoff between attempts
- Maximum retry limit to prevent infinite loops

### Rate Limiting

Automatic rate limiting to prevent API throttling:

- 100ms delay between API calls
- Respects GitHub/GitLab rate limits
- Prevents 429 (Too Many Requests) errors

## Contributing

Pull requests are welcome! If you encounter issues, feel free to open an issue on GitHub.

### Development Setup

1. Clone the repository
2. Install dependencies: `go mod download`
3. Run tests: `go test ./...`
4. Build: `go build -o reposync .`

## License

MIT License. See `LICENSE` for details.
