# lazy-clone

A command line tool for downloading files from a subdirectory in a Github repo.
For when individually downloading them is too tedious and cloning the entire repo is overkill :full_moon_with_face:

## Usage

- Clone this repo
- Run `go build`
- Run `./fetch [flags]`

For example, to download all the css files in [rabbitmq-website](https://github.com/rabbitmq/rabbitmq-website) from the [`site/css`](https://github.com/rabbitmq/rabbitmq-website/tree/live/site/css) subdirectory:

```sh
$ ./fetch -user rabbitmq -repo rabbitmq-website -dir site/css

amqp-0-9-1-quickref.css 1710
amqp-0-9-1-reference.css 3197
highlightjs_style.css 1211
rabbit-ie6.css 782
rabbit-next.css 737
rabbit.css 29069
search.css 894
tutorial.css 7256
```

Not overriding the `-dry-run` flag means by default this is a dry run. No files will be downloaded, and the list of files found and their size are simply printed.
To download the files listed above into a directory called `css`, the following command can be used:

```sh
$ ./fetch -user rabbitmq -repo rabbitmq-website -dir site/css -out css -dry-run false
```

## Options

| Flag       | Description                                   | Default |
| ---------- | --------------------------------------------- | ------- |
| `-dir`     | The name of the subdirectory to download from |
| `-dry-run` | Whether or not to execute the download        | true    |
| `-out`     | The output directory to create                | out     |
| `-repo`    | The Github repository name                    |
| `-user`    | The Github user                               |
