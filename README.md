# EduPage2 - server
[![codecov](https://codecov.io/gh/DislikesSchool/EduPage2-server/graph/badge.svg?token=BXliJAEhQz)](https://codecov.io/gh/DislikesSchool/EduPage2-server)

The backend API server for the EduPage2 project.

## Installation
There are 3 ways to install and run the server:
1. [Bare-metal](#bare-metal)
2. [Docker](#docker)
3. [Docker-compose](#docker-compose)

### Bare-metal
1. Install [Golang](https://go.dev/doc/install)
2. Clone the repository
```bash
git clone https://github.com/DislikesSchool/EduPage2-server.git
```
3. Install dependencies
```bash
go mod download
```
4. Build the project
```bash
go build ./cmd/server
```
5. Run the server
```bash
./server
```

### Docker
1. Clone the repository
```bash
git clone https://github.com/DislikesSchool/EduPage2-server.git
```
2. Build the Docker image
```bash
docker build -t edupage2-server .
```
3. Run the Docker container
```bash
docker run --mount type=bind,src=./config.yaml,dst=/config.yaml,ro -p 8080:8080 edupage2-server
```

### Docker-compose
1. Clone the repository
```bash
git clone https://github.com/DislikesSchool/EduPage2-server.git
```
2. Start the container
```bash
docker compose up -d
```

## Configuration
The server can be configured using the `config.yaml` file. The default configuration is as follows:
```yaml
# General server configuration
server:
  # The host address to bind to (do not change this unless you know what you are doing)
  host: "0.0.0.0"
  # The port to bind to
  port: "8080"
  # The mode to run the server in (either "development" or "production")
  mode: "production"

# School configuration
schools:
  # Whitelisted schools (users will only be allowed to login if they are a member of this school)
  whitelist:
    - "schoolname" # This is your school's unique name (https://schoolname.edupage.org)
  # Whether to use the whitelist as a blacklist instead
  blacklist: false # Set this to true with a empty whitelist to allow users from any school

# JWT configuration
jwt:
  # The secret key to use for signing JWT tokens (change this to a secure random value)
  secret: "YourDefaultSecretKey"
```

If you want to host a public instance, remove all schools from the whitelist and set `blacklist` to `true`.

## API
The server provides a RESTful API for the frontend to interact with. The API documentation can be found [here](https://ep2.ypal.me/docs/index.html), or at the same endpoint on your own instance.
