# netlify-cms-oauth-provider-go

Netlify-CMS oauth client sending token in form as Netlify service itself, implementation in Go (golang)

inspired by [netlify-cms-github-oauth-provider](https://github.com/vencax/netlify-cms-github-oauth-provider) (node-js). Thanks Václav!

## 1) Install

Download the binary and run it.

## 2) Config

## Session Token Config

Session Secret token needs to be injected into the environment to run.

```bash
export SESSION_SECRET=secret-token-here
```

### Auth Provider Config

Configuration is done with environment variables, which can be supplied as command line arguments, added in your app hosting interface, or loaded from a .env ([dotenv](https://github.com/motdotla/dotenv)) file.

**Example .env file:**

```
PORT=3000
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

```yaml
backend:
  name: github
  repo: user/repo   # Path to your Github repository
  branch: master    # Branch to update
  base_url: https://your.server.com      # Path to ext auth provider
  api_root: http://ent-github.com
  site_domain: my.site.com
```
