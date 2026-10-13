# gh-auth

A small CLI that keeps your local **git** identity in sync with the GitHub
account you are currently authenticated as through the [GitHub CLI](https://cli.github.com/)
(`gh`), **including GPG commit signing**.

When you work with more than one GitHub account, it is easy to switch `gh`
accounts but forget to update `git config` — so commits end up authored by (or
signed with the key of) the wrong identity. `gh-auth` reads the active `gh`
account and writes the matching `user.name`, `user.email`, `user.signingkey`
and `commit.gpgsign` into your git configuration for you.

## How it works

`gh-auth` shells out to the tools you already have:

- **`gh`** to read the active account (`gh api user`) and to switch accounts
  (`gh auth switch`).
- **`gpg`** (optional) to find a secret key whose UID matches the account email,
  so commit signing can be enabled automatically.
- **`git`** to read and write the configuration.

The account email is normalized by stripping any `+tag` subaddress (so
`user+work@example.com` becomes `user@example.com`), and GPG key ids are
reduced to the 16-character long id. When no matching key is found, you are
prompted to enter one or to skip signing.

## Requirements

- [`gh`](https://cli.github.com/) — required.
- `gpg` — optional, only needed for GPG commit signing.
- `git` — required.

## Installation

Install into your `$GOPATH/bin`:

```bash
make install
```

Or build a binary into `./bin`:

```bash
make build
```

Or install directly with Go:

```bash
go install github.com/jcchavezs/gh-auth/cmd/gh-auth@latest
```

## Usage

```text
gh-auth <command>
```

### `switch`

Runs `gh auth switch` to let you pick a GitHub account, then configures your
**global** git identity (name, email, and GPG signing) to match it.

```bash
gh-auth switch
```

### `sign-status`

Shows the active GitHub account alongside your current global git signing
configuration (name, email, GPG key, and whether signing is enabled).

```bash
gh-auth sign-status
```

### `status`

Passthrough to `gh auth status`.

```bash
gh-auth status
```

## Development

```bash
make test       # unit tests
make test-e2e   # end-to-end tests (require git and gh)
make lint       # golangci-lint
make vulncheck  # govulncheck
```

Run `make install-tools` first to install the linting and vulnerability
scanning tools.

## License

See [LICENSE](LICENSE).
