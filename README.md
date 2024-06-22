# Relay Wizard

![0622](https://github.com/nodetec/relaywizard/assets/29136904/eb226b30-9250-43c6-ba2a-0361446d790b)

Relay Wizard is a cli tool that helps you bootstrap a [nostr](https://nostr.com/) relay.

The program will automate the following steps:

1. install necessary dependencies
1. configuring nginx
1. setting up a firewall
1. obtaining a TLS certificate for HTTPS
1. installing the relay software
1. setting up a systemd service for your relay

## Installation

To install a relay, spin up a new Debian server, hook up a domain name, and run the following command:

```bash
curl -sL https://relayrunner.org/relaywizard.sh | bash
```

## Learn more

If you want to learn more about how to setup a relay from scratch, check out [relayrunner.org](https://relayrunner.org)

If you just want to know enough to get started, read the following sections to get a server, hook up a domain name and setup remote access:

- [Get a server](https://relayrunner.org/server/get-a-server)

- [Get a domain](https://relayrunner.org/server/domain-name)

- [Remote access](https://relayrunner.org/server/remote-access)

from here you should be able to run the installation command above and get started.

## Contributing

If you want to contribute consider adding a new package manager and test the script out on another Linux Distro, I also have plans to support multiple relay implementation options.
