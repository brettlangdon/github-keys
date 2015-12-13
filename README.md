github-keys
===========
`github-keys` is a web server written in golang which acts as a proxy to `https://api.github.com/users/<username>/keys`.

## Installing
`go get github.com/brettlangdon/github-keys`

## Running
`github-keys -username "<username>"`

### Options

`github-keys -username "<username>" [-ttl <seconds>] [-listen "[host]:<port>"]`
