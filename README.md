[github.com]: github.com
[OAuth device flow]: https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps#device-flow
[token settings]: https://github.com/settings/tokens
[obtaining a classic token]: https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token#creating-a-personal-access-token-classic
[pandoc]: https://pandoc.org/

# snips

The program gathers GitHub activity data about specific users over a common
time period and prints it as markdown.

## Installation

Assuming you have the `go` tool installed and `$(go env GOPATH)/bin` on your `PATH`:
```
go install github.com/monopole/snips@latest
```

## Usage

To get data from a GitHub enterprise instance at _Acme Corporation_ 
for several users over a common period of September in 2020:

```
snips \
    --domain github.acmecorp.com \
    --token blah-blah-blah \
    --day-start 2020-sep-01 \
    --day-count 30 \
     torvalds thockin spf13 > /tmp/snips.md
```

To get recent data for user `thockin` from [github.com]:

```
snips thockin > /tmp/snips.md
```

The time period is a start date and a day count that includes the start date.

The default value for `--day-count` is 14 (two weeks).

If `--day-start` is omitted, a value of _today_ minus `day-count` is used.


### Authentication Token

The absence of a value for both the `--token` flag and the shell
variable `GH_TOKEN` triggers an [OAuth device flow].

The flow helps the user obtain an API access token from
the specified `--domain` (the default is `--domain github.com`).

A _newly_ obtained token will be echoed to `stderr`.

Placing the token value into the shell variable `GH_TOKEN`
allows one to omit the `--token` flag in subsequent usage.

#### Fallback to classic flow

If the [OAuth device flow] fails for some reason (which it will
if this app has no clientId for the `--domain` being used),
then try the instructions for [obtaining a classic token]
(being sure to do this from your desired `--domain`).

In this flow, select the scopes:
```
 [x] admin/read (to see organization membership)
 [x] repo (to see pull requests)
 [x] user (to see public info about the user)
```

A classic token may be used with the `--token` flag
as if it had been created via the OAuth flow.

Protect this classic token like a password. During creation,
give it an expiration period, and/or delete it after
use at the [token settings] page.

The tool will not return _private_ GitHub data for a username,
unless the username matches the user that authenticated to obtain the token.

### Rendered Output

To render the output to a browser, try [pandoc].

With no arguments, `pandoc` converts markdown to HTML.

```
sudo apt install pandoc
```

```
snips --token {token} monopole 2020-01-01 28 |\
    pandoc |\
    google-chrome "data:text/html;base64,$(base64 -w 0 <&0)"
```
