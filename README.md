# Nail Care

A vanity URL tool for Go packages. (Blame [Twitch](https://twitch.tv/hayden_dev) for the weird naming)

To use it, build the source and run the command as such:

```shell
$ ./nailcare [path],[type],[source] [path],[type],[source] [path],[type],[source] [path],[type],[source]
```

For example:

```shell
$ ./nailcare \
    zetman,git,https://github.com/hbjydev/zetman \
    centra/component-base,git,https://github.com/centra-oss/component-base
```

Now, if you browse to `localhost:8080/zetman` in your browser, you will be redirected to `pkg.go.dev/go.h4n.io/zetman`.

