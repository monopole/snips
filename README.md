# snips

This program acquires data about a github user over a specific
time period and prints it as markdown.


### Install

Something like:

```
go install github.com/monopole/snips
```

### Usage

```
snips {githubAuthToken} {githubUser} [{dateStart} [{dayCount}]] > snips.md
```

e.g.

```
go run . deadbeef0000deadbeef monopole 2020-04-06 > snips.md
```

To render the output to a browser, try


```
# With no arguments, pandoc converts markdown to html
sudo apt install pandoc
```

```
snips $token monopole 2020-01-01 28 |\
    pandoc |\
    google-chrome "data:text/html;base64,$(base64 -w 0 <&0)"
```

### Authentication

This program is not an oauth-App nor a github-App; it requires a "classic"
auth token that should be protected as carefully as your github password.

[classic github access token]: https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token#creating-a-personal-access-token-classic

> See these instructions on creating a [classic github access token].

```
go run . deadbeef0000deadbeef monopole 2020-04-06
```

The program gathers data for a time period specified by a start date and
number of days to examine inclusive of the start date.

This program is not an oauthApp nor a githubApp; it requires a "classic"
auth token that should be protected as carefully as your github password.

If the token owner and {githubUser} aren't the same, the program will fail
to read private repos associated with {githubUser}.

> To get a token, read:
> 
> https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token
> 
> Use "Generate new token (classic) for general use."
>
> Select scopes
>  [x] repo
>  [x] admin/read (to see membership)
>  [x] user
> ```

If the token owner and `{githubUser}` aren't the same, the program will fail
to read private repos associated with `{githubUser}`.


_TODO: convert this app to oauth or githubApp flow._

_TODO: upgrade go-github deps._
