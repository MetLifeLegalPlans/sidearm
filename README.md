# sidearm

Demolish your infrastructure.

(spirit animal coming soon to a repo near you)

## Installing

### From source (recommended)

First make sure that you have the latest stable version of [the Go compiler](https://go.dev/dl/) installed and available on your `$PATH`. Then,

```sh
$ go install github.com/MetLifeLegalPlans/sidearm@latest
```

will place an executable called `sidearm` in your `$GOROOT/bin` directory. By default, `$GOROOT` is `$HOME/go`.

### With Docker

```sh
$ docker run -v $PWD/config.yml:/config.yml -t metlifelegalplans/sidearm:latest
```

You will need to mount your config file into the container so that it can be read.

## Running

`sidearm` can be run in any of four ways:

1. `server` mode (generates and distributes tasks to clients)
2. `client` mode (receives and performs tasks from clients)
3. `dashboard` mode (runs a dashboard that receives and stores statistics)
4. `standalone` mode (runs a client, server, and dashboard session all at once)

All of these modes require a config file to be passed in as the first argument, but the configuration is identical between all of them - meaning you can use the same config file for both server and client.

```sh
$ sidearm [client|server|standalone] <config.yml>
```

In addition, there will be a `dashboard` mode in later releases that collects information from your `client` instances and displays them in a nice format.

## Configuration

All of `sidearm`'s configuration is handled in a single YAML file. An annotated version of this showcasing all options is below.

```yaml
queue:
  bind: tcp://*:5557 # Where server mode will bind its outbound message queue

  # If you are running the server on a different machine, then this will be
  # the address of the server as your client machine sees it
  connect: tcp://example.com:5557
sink: # Configure the dashboard mode receiver
  bind: tcp://*:5558
  connect: tcp://example.com:5558
dbPath: ./out.db # Where the dashboard receiver stores its data

requests: 10000 # The number of requests to send per second at peak - REQUIRED
duration: 300 # How many seconds it takes to reach the max send rate, default 0

# scenarios is a list of possible request types that can be sent out
# Scenarios are chosen at random for every request, but can be assigned
# weights to influence the odds of them being picked
scenarios:
  # A URL is the only required field in a scenario
  - url: https://legalplans.com
    # HTTP verb, defaults to GET
    method: GET

    # Comparative likelihood that this scenario will be chosen
    # Defaults to 1, valid values are >= 1
    weight: 2

  - url: https://members.legalplans.com
    method: POST
    weight: 5

    # Dictionary that is mapped to the request body and sent
    # as JSON. Ignored for request types that have no body.
    body:
      myKey: 123
      otherKey:
        nestedKey: 'nestedValue'

  - url: https://login.legalplans.com/api/v1/loginOrSomething
    method: POST
    weight: 3
    body:
      username: horse
      password: radish
      email: horseradish@legalplans.com
```

## Minimal Configuration

That was an absurdly long YAML file. What's the minimum you can get away with?

```yaml
queue:
  bind: tcp://*:5557
requests: 100
scenarios:
  - url: https://legalplans.com
```

...well that wasn't so bad after all.

## Building from source

Outside of the Go compiler, you will need to have [ZeroMQ](https://zeromq.org/download) installed on your machine. Their website provides detailed instructions for a variety of different distributions. Version 4 of the library or above is required.

Once you have that, it is as simple as building any other Go program.

```sh
$ git clone https://github.com/MetLifeLegalPlans/sidearm.git
$ cd sidearm
$ go build
```

If ZeroMQ and its headers were installed correctly, you should now have an executable in the current directory named `sidearm`.
