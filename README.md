# snips

The program gathers github data about a specific user over a specific 
time period and prints it to `stdout` as markdown.

The time period is a start date and a day count inclusive of the start date.

### Usage

```
go install github.com/monopole/snips@latest
```

```
snips [--domain github.acmecorp.com] [--token {tokenForDomain}] \
    {githubUser} [{dateStart} [{dayCount}]] > snips.md
```

e.g., get data from [github.com] (the default domain) for one day by omitting `{dayCount}`:

```
snips monopole 2020-04-06 > snips.md
```

To render the output to a browser, try `pandoc`.
With no arguments, `pandoc` converts markdown to HTML.

```
sudo apt install pandoc
```

```
snips --token {token} monopole 2020-01-01 28 |\
    pandoc |\
    google-chrome "data:text/html;base64,$(base64 -w 0 <&0)"
```

### Authentication

[OAuth device flow]: https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps#device-flow

This program uses an [OAuth device flow] to get an API access token.

If that fails for some reason, you can use a "classic" (scope-free) token:

> To get a token, follow the instructions at:
> 
> https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token
> 
> Use _Generate new token (classic) for general use._
>
> Select the scopes:
> ```
> [x] repo
> [x] admin/read (to see membership)
> [x] user
> ```

Once you have a token, you can specify it as the value of the `--token` flag.


On usage, if the token owner and `{githubUser}` aren't the same, the program will fail
to read private repos associated with `{githubUser}`.
