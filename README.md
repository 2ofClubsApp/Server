<h1 align="center">
  <p align="center">2ofClubs Server</p>
  <h4 align="center">Finding the perfect club for you!</h4>
  <a href="https://2ofClubs.app><img src="https://avatars3.githubusercontent.com/u/64863952?s=400&u=293c427becbc89d1388ece6182462f14ad81d3a5&v=4" alt="2ofClubs"></a>
</h1>
<p align="center">
  <a href="https://goreportcard.com/report/github.com/2ofClubsApp/2ofclubs-server"><img src="https://goreportcard.com/badge/github.com/2ofClubsApp/2ofclubs-server" alt="Go Report Card"/> </a>
  <a href="https://hub.docker.com/r/2ofclubsapp/server"><img src="https://img.shields.io/docker/cloud/build/2ofclubsapp/server" alt="Docker Build Status"></a>
  <a href="https://hub.docker.com/r/2ofclubsapp/server"><img src="https://img.shields.io/docker/pulls/2ofclubsapp/server" alt="Docker Pulls"></a>
  <a href="https://pkg.go.dev/mod/github.com/2-of-clubs/2ofclubs-server"><img src="https://pkg.go.dev/badge/github.com/2ofClubsApp/2ofclubs-server?status.svg" alt="Go Doc"/></a>
  <a href="#License" alt="License"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License"/></a>
</p>

<p align="center">
  <a href="#introduction">Introduction</a> •
  <a href="#installation">Installation</a> •
  <a href="#documentation">Documentation</a> •
  <a href="#blog">Blog</a> •
  <a href="#contact">Contact</a> •
  <a href="#honorable-mentions">Honorable Mentions</a> •
  <a href="#license">License</a>
</p>

## Introduction
2ofClubs is a web app for helping students find clubs that suit their preferences.

- **Easy to Use!**
> Find the perfect club for you in just a couple of swipes!

- **Explore anywhere**
> You can find clubs and their hosted events wherever you go!

## Installation
2ofClubs-Server is available as a [Docker Container](https://hub.docker.com/r/2ofclubsapp/server) on [Docker Hub](https://hub.docker.com)

### Deployment using Docker

1. Pull the latest image from Docker Hub:

```
docker pull 2ofclubsapp/server
```

2. Run `docker-compose` 
Make sure your version of Docker supports docker-compose v3.3 or later

```
docker-compose up --build -d
```

3. The server should be listening and serving on port `8080`

You can check out our full [installation guide](https://2ofclubs.app/docs/installation), app requirements and more on our website.

### Configuration
* By default, the 2ofClubs-Server is listening and serving on port `8080`. This can be changed in the `docker-compose.yaml` file.
* App environemnt variables can be set in `app.env`
* Database environment variables can be set in `db.env`

## Documentation
Our documentation can be found [here](https://2ofclubs.app/docs)

## Blog
Checkout our [blog](https://2ofclubs.app/blog) for updates and changes!

## Contact
You can contact us through these channels:
- [Email](mailto:hello@2ofclubs.app)
- [Github Issues](https://github.com/2ofClubsApp/2ofclubs-server/issues)

## Contributing
- Properly [document](https://blog.golang.org/godoc) and [format](https://golang.org/doc/effective_go.html#formatting) your Go Code.
- If you're submitting a feature request or PR, please follow their associated templates

For any code that can be reworked/improved upon, feel free to submit a PR

## Honorable mentions
Thank you so much for your help and guidance! We couldn't have done this without you!

[Nicholas Hawkes](https://github.com/hawkesnc) and [Michael Liut](https://github.com/MichaelLiut)

## License
The 2ofClubs-Server is [MIT licensed](./LICENSE)

<a href="https://2ofClubs.app"><img src="https://user-images.githubusercontent.com/41246112/83603397-5d4d6800-a542-11ea-9dcd-3916bc86474d.png" alt="2ofClubsServer"/>
