---
name: switching-gh-auth-profile
description: Run gh, git or other commands under a specific gh-auth profile (GitHub account), so the right credentials, git author/committer identity and GPG signing are used. Use when the user asks to run something "as" or "with" a given gh profile/account, e.g. "commit with the work profile" or "open a PR using jcchavezs".
---

# Run commands with a gh-auth profile

`gh-auth shell <profile> -- <command> [args...]` runs a single command with
`GH_CONFIG_DIR`, git author/committer identity and (if configured) commit
signing set for that profile. Nothing is changed globally.

## Usage

```bash
gh-auth shell <profile> -- git commit -m "message"
gh-auth shell <profile> -- gh pr create --fill
gh-auth shell <profile> -- git push
```

- The `--` is required; everything after it is the command, run directly (not via a shell).
  For pipes or `&&`, wrap it: `gh-auth shell <profile> -- sh -c 'git add . && git commit -m "msg"'`.
- The exit status of the command is propagated.
- Nested use fails with "profile is already enabled" if `GHAUTH_PROFILE` is set.
- If the profile has never been set up, `gh auth login` runs interactively,
  which cannot work non-interactively: ask the user to run `gh-auth shell <profile>` once first.
- Add `--refresh-profile` (before `--`) to re-resolve the git identity from GitHub.

## Steps

0. Check that the CLI is installed with `command -v gh-auth`. If it is missing, ask the user whether to install it
   (do not install without their confirmation). If they agree, run
   `gh release download --repo jcchavezs/gh-auth --pattern 'install.sh' -O - | sh` (requires `gh`; honors `INSTALL_DIR`
   and `VERSION`). Then re-check `command -v gh-auth`, and if the install dir is not in `PATH`, tell the user how to add it.
1. Get the profile name from the user (profiles live in `~/.config/gh-auth/<profile>`; `ls` that directory to list them).
2. Prefix each gh/git command that needs that identity with `gh-auth shell <profile> --`.
3. Never commit or push unless the user asked for it.
