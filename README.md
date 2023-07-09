[github.com]: github.com
[OAuth device flow]: https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps#device-flow
[token settings]: https://github.com/settings/tokens
[obtaining a classic token]: https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token#creating-a-personal-access-token-classic
[pandoc]: https://pandoc.org/
[`go`]: https://go.dev

# snips

Reports GitHub activity about specific users over a common period of days.

## Usage

To get recent data for user _thockin_ from [github.com]:

```
snips thockin > /tmp/snips.html
```

The report is emitted as HTML to `stdout`.

> A one-liner to render to chrome on ubuntu:
>
> ```
> snips {args} | google-chrome "data:text/html;base64,$(base64 -w 0 <&0)"
> ```

Use `--md` to get markdown instead.

To get data from a GitHub enterprise instance at _Acme Corporation_
for several users during September 2020:

```
snips \
    --gh-domain github.acmecorp.com \
    --day-start 2020-Sep-01 \
    --day-count 30 \
     alice bob charlie > /tmp/snips.html
```

The time period is measured in days.
It can be specified using any two of the
following flags:
`--day-start`, `--day-end`, and `--day-count`.
The default _day-end_ is _today_.
The default _day-count_ is _14_.

## Installation

Install the [`go`] tool.

Assure that your `PATH` includes the value of `$(go env GOPATH)/bin`.

Enter:
```
go install github.com/monopole/snips@latest
```

## GitHub Authentication

An [OAuth device flow] for the given `--gh-domain` is triggered when either

 * the `--just-get-gh-token` flag is present,
 * or both the shell variable `GH_TOKEN` and the `--gh-token` override flag are empty.

This command
```
export GH_TOKEN=$(snips --gh-domain github.acmecorp.com --just-get-gh-token)
```
sets `GH_TOKEN`, allowing subsequent `snips` usage to skip prompting
for a password.

If the OAuth flow fails for some reason (e.g.
this program has no clientId for the `--gh-domain` being used),
then try the instructions for [obtaining a classic token].
In that flow, select these scopes:
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
