# Twitter like feed

[![Build status][github-actions-image]][github-actions-url]

[github-actions-image]: https://github.com/kenchan0130/twitter-like-feed/workflows/CI/badge.svg
[github-actions-url]: https://github.com/kenchan0130/twitter-like-feed/actions?query=workflow%3A%22CI%22

A web application that provides the like of the target account of Twitter as rss feed.

Sample is [here](https://twitter-like-feed-19rl.onrender.com).

## Environment Variables

| variable              | description                                                         |
|-----------------------|---------------------------------------------------------------------|
| BEARER_TOKEN          | Bearer Token authenticates requests on behalf of your developer App |
| CACHE_EXPIRES_SECONDS | Cache expires seconds for RSS Item, default 7200 (2 hr)             |

## Endpoints

| endpoint          | content                        |
|-------------------|--------------------------------|
| `/`               | redirect to `/health`.         |
| `/health`         | return "ok" as text.           |
| `/feed/:username` | return a rss feed of username. |

## Development

```sh
go get -u github.com/kenchan0130/twitter-like-feed
```

You may also clone this project instead.
And, please run the program.

```sh
go run main.go
```

## Deploy

Any changes to the mastar branch are automatically deployed to heroku by GitHub Action.
