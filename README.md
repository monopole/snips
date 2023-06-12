[github.com]: github.com
[OAuth device flow]: https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps#device-flow
[token settings]: https://github.com/settings/tokens
[obtaining a classic token]: https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token#creating-a-personal-access-token-classic
[pandoc]: https://pandoc.org/
[`go`]: https://go.dev

# snips

Reports GitHub activity about specific users over a common time period.

### Installation

Install the [`go`] tool and assure that `$(go env GOPATH)/bin` is on your `PATH`.

Then:
```
go install github.com/monopole/snips@latest
```

## Usage

To get data from a GitHub enterprise instance at _Acme Corporation_ 
for several users during September 2020:

```
snips \
    --domain github.acmecorp.com \
    --day-start 2020-Sep-01 \
    --day-count 30 \
     alice bob charlie > /tmp/snips.html
```

The report is emitted as HTML to `stdout`.

To render directly to a browser, try:

```
snips {args} |\
    google-chrome "data:text/html;base64,$(base64 -w 0 <&0)"
```

Use `--markdown` to get markdown instead of HTML.

#### Taking a recent snapshot

To get recent data for user `thockin` from [github.com]:

```
snips thockin > /tmp/snips.html
```

The time period is a start date and a day count,
or a start and an end date, inclusive.

The default value for `--day-count` is 14 (two weeks).

If `--day-start` is omitted, a value of _today_ minus `day-count` is used.

The default `--domain` is `github.com`.


### Authentication Token

An [OAuth device flow] for the given `--domain` is triggered when either

 * the `--get-gh-token` flag is present,
 * or both the shell variable `GH_TOKEN` and the `--gh-token` override flag are empty.

Use this
```
export GH_TOKEN=$(snips --domain github.acmecorp.com --get-gh-token)
```
to login once, then use `snips` multiple times.

This program cannot return activity conducted in private repos,
unless the person being looked up matches the person who got the token.

#### Fallback to classic flow

If the OAuth flow fails for some reason (e.g.
this program has no clientId for the `--domain` being used),
then try the instructions for [obtaining a classic token].

In this flow, select the scopes:
```
 [x] admin/read (to see organization membership)
 [x] repo (to see pull requests)
 [x] user (to see public info about the user)
```

A classic token may be used with the `--gh-token` flag
as if it had been created via the OAuth device flow.

Protect this classic token like a password. During creation,
give it an expiration period, and/or delete it after
use at the [token settings] page.
