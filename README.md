# Twitter like feed

A web application that provides the like of the target account of Twitter as rss feed.

Sample is [here](https://twitter-like-feed.herokuapp.com/).

## Endpoints

| endpoint          | content                                        |
|-------------------|------------------------------------------------|
| `/`               | return "pong" as text                          |
| `/feed/:username` | return a rss feed of username (without atmark) |

## Development

```sh
go get -u github.com/kenchan0130/twitter-like-feed
go run main.go
```

## Deploy

```sh
git push heroku master
```
