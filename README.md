# snips

The program gathers github data about a specific user over a specific 
time period, and prints it to `stdout` as markdown.

The period is specified by a start date and the number of days to
examine inclusive of the start date.

### Usage

```
go run . [--domain github.acmecorp.com] {githubAuthToken} {githubUser} [{dateStart} [{dayCount}]] > snips.md
```

e.g., show data for one day by omitting `{dayCount}`:

```
go run . $token monopole 2020-04-06 > snips.md
```

To render the output to a browser, try `pandoc`.
With no arguments, `pandoc` converts markdown to HTML.

```
sudo apt install pandoc
```

```
snips $token monopole 2020-01-01 28 |\
    pandoc |\
    google-chrome "data:text/html;base64,$(base64 -w 0 <&0)"
```

### Authentication

[githubApp]: https://docs.github.com/en/apps/creating-github-apps/setting-up-a-github-app/creating-a-github-app

This program is not an oauth-based App nor a [githubApp]; it requires a "classic"
auth token that should be protected as carefully as your github password.

> To get a token, read:
> 
> https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token
> 
> Use "Generate new token (classic) for general use."
>
> Select the scopes:
> ```
> [x] repo
> [x] admin/read (to see membership)
> [x] user
> ```

On usage, if the token owner and `{githubUser}` aren't the same, the program will fail
to read private repos associated with `{githubUser}`.

_TODO: convert this app to oauth or githubApp flow._
