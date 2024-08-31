# tfvpn

tfvpn is a tool that instantly connects you to a vpn server located on the threefold grid and tunnels all of your traffic through the server.

## Requirements

- [python 3.8+](https://www.python.org/downloads/)
- [sshuttle](https://sshuttle.readthedocs.io/en/stable/installation.html)
- sudo or root access to install the required packages
  - you will be prompted to enter your password to install the required packages

Actually, the tool will install the requirements for you if they are not already installed on your system; however, you are more than welcomed to to install them beforehand.

## Usage

`REQUIRED` EnvVars:

- `MNEMONICS`: your secret words
- `NETWORK`: one of (dev, qa, test, main), default is `dev`

```bash
tfvpn [OPTIONS] [COMMAND]

COMMANDS:
    - up:      deploy the vpn server on the grid and connect to it
    OPTIONS:
        - --region: the region of the vpn server [optional]
        - --country: the country of the vpn server [optional]
        - --city: the city of the vpn server [optional]
    - down:    disconnect from the vpn server and cancel the deployment
```

Export env vars using:

```bash
export MNEMONICS=your_mnemonics
export NETWORK=working_network
```

Run:

```bash
make build
```

To use any of the commands, run:

```bash
./bin/tfvpn [COMMAND]
```

For example:

```bash
./bin/tfvpn up --country egypt --city cairo --region africa
```

Or just run:

```bash
./bin/tfvpn up
```

if you are not interested in specifying the region, country, or city.

## Usage For Each Command

### up

The up command deploys the vpn server on the grid and connects the host to the server.

Specify the region, country, and city of the vpn server using the flags if a specific location for the server, if desired, and the tool will try to find a server in that location.

The tool uses the public/private key pair located in the `~/.ssh` directory and named `id_rsa.pub` and `id_rsa` to connect to and authenticate with the server.

```bash
./bin/tfvpn up [OPTIONS]
```

OPTIONS:

- `--region`: the region of the vpn server `optional`
- `--country`: the country of the vpn server `optional`
- `--city`: the city of the vpn server `optional`

### Example

```bash
./bin/tfvpn up
```

output:

```bash
10:38AM INF starting peer session=tf-507370 twin=8658
10:38AM INF checking if sshuttle, and python3 are installed on the system
10:38AM INF all requirements are installed/present successfully!
10:38AM INF deploying vpn server... name=vpn node_id=11
10:38AM INF vpn server deployed successfully! public_ip=185.206.122.34
10:38AM INF using ssh agent for authentication
10:38AM INF SSH connection attempt failed, retrying...
10:38AM INF SSH connection attempt failed, retrying...
10:38AM INF SSH connection attempt failed, retrying...
10:38AM INF SSH connection attempt failed, retrying...
10:38AM INF SSH connection attempt failed, retrying...
10:38AM INF SSH connection attempt failed, retrying...
10:38AM INF SSH connection attempt failed, retrying...
10:38AM INF SSH connection attempt failed, retrying...
10:38AM INF SSH connection attempt failed, retrying...
10:38AM INF connecting to the vpn server...
10:38AM INF connection established successfully!
```

### down

The down command cancels the vpn server deployment, stops the vpn connection, and remove the server's host from the known hosts in your machine.

```bash
./bin/tfvpn down
```

### Example

```bash
./bin/grid-compose down
```

output:

```bash
10:39AM INF starting peer session=tf-508547 twin=8658
10:39AM INF disconnecting from vpn server
10:39AM INF killing sshuttle...
10:39AM INF killed the shuttle process
10:39AM INF canceling deployment...
10:39AM INF canceling contracts project name=8658/vpn
10:39AM INF project is canceled project name=8658/vpn
10:39AM INF deployment canceled successfully
10:39AM INF removing host 185.206.122.34 from known hosts...
10:39AM INF host 185.206.122.34 is removed
10:39AM INF disconnected from vpn server
```

## Testing

To run the tests, run:

```bash
make test
```

## Linting

To lint the code, run:

```bash
make lint
```

## Clean

To clean the project, run:

```bash
make clean
```
