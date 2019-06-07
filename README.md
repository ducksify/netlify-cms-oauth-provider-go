# netlify-cms-oauth-provider-go

[![CircleCI Build Status](https://circleci.com/gh/CircleCI-Public/circleci-demo-go.svg?style=shield)](https://circleci.com/gh/CircleCI-Public/circleci-demo-go) [![MIT Licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/CircleCI-Public/circleci-demo-go/master/LICENSE.md)
[![Go Report Card](https://goreportcard.com/badge/github.com/maarek/netlify-cms-oauth-provider-go)](https://goreportcard.com/report/github.com/maarek/netlify-cms-oauth-provider-go)
[![Downloads](https://img.shields.io/github/downloads/maarek/netlify-cms-oauth-provider-go/latest/total.svg)](https://github.com/maarek/netlify-cms-oauth-provider-go/releases)
[![Latest release](https://img.shields.io/github/release/maarek/netlify-cms-oauth-provider-go.svg)](https://github.com/maarek/netlify-cms-oauth-provider-go/releases)

Netlify-CMS oauth client sending token in form as Netlify service itself, implementation in Go (golang)

inspired by [netlify-cms-github-oauth-provider](https://github.com/vencax/netlify-cms-github-oauth-provider) (node-js). Thanks Václav!

## 1) Install

```bash
# binary will be $GOPATH/bin/netlify-cms-oauth-provider-go
curl -sfL https://raw.githubusercontent.com/maarek/netlify-cms-oauth-provider-go/master/install.sh | sh -s -- -b $GOPATH/bin

# or install it into ./bin/
curl -sfL https://raw.githubusercontent.com/maarek/netlify-cms-oauth-provider-go/master/install.sh | sh -s

# In alpine linux (as it does not come with curl by default)
wget -O - -q https://raw.githubusercontent.com/maarek/netlify-cms-oauth-provider-go/master/install.sh | sh -s
```

## 2) Config

## Session Token Config

Session Secret token needs to be injected into the environment to run.

```
export SESSION_SECRET=secret-token-here
```

### Auth Provider Config

Configuration is done with environment variables, which can be supplied as command line arguments, added in your app hosting interface, or loaded from a .env ([dotenv](https://github.com/motdotla/dotenv)) file.

**Example .env file:**

```
HOST=localhost:3000
CALLBACK_HOST=localhost:3000
GITHUB_HOST=
GITHUB_KEY=
GITHUB_SECRET=
BITBUCKET_KEY=
BITBUCKET_SECRET=
GITLAB_KEY=
GITLAB_SECRET=
```

**Client ID & Client Secret:**
After registering your Oauth app, you will be able to get your client id and client secret on the next page.

### CMS Config

You also need to add `base_url` to the backend section of your netlify-cms's config file. `base_url` is the live URL of this repo with no trailing slashes.

```
backend:
  name: github
  repo: user/repo   # Path to your Github repository
  branch: master    # Branch to update
  base_url: https://your.server.com      # Path to ext auth provider
  api_root: http://ent-github.com/api/v3 # Path to enterprise github
```
