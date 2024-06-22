<div align="center"><p>
    <h1>Relay Wizard</h1>
    <a href="https://github.com/nodetec/relaywizard/releases/latest">
      <img alt="Latest release" src="https://img.shields.io/github/v/release/nodetec/relaywizard?style=for-the-badge&logo=starship&color=C9CBFF&logoColor=D9E0EE&labelColor=302D41" />
    </a>
    <a href="https://github.com/nodetec/relaywizard/pulse">
      <img alt="Last commit" src="https://img.shields.io/github/last-commit/nodetec/relaywizard?style=for-the-badge&logo=starship&color=8bd5ca&logoColor=D9E0EE&labelColor=302D41"/>
    </a>
    <a href="https://github.com/nodetec/relaywizard/stargazers">
      <img alt="Stars" src="https://img.shields.io/github/stars/nodetec/relaywizard?style=for-the-badge&logo=starship&color=c69ff5&logoColor=D9E0EE&labelColor=302D41" />
    </a>
    <a href="https://github.com/nodetec/relaywizard/issues">
      <img alt="Issues" src="https://img.shields.io/github/issues/nodetec/relaywizard?style=for-the-badge&logo=bilibili&color=F5E0DC&logoColor=D9E0EE&labelColor=302D41" />
    </a>
    <a href="https://github.com/nodetec/relaywizard">
      <img alt="Repo Size" src="https://img.shields.io/github/repo-size/nodetec/relaywizard?color=%23DDB6F2&label=SIZE&logo=codesandbox&style=for-the-badge&logoColor=D9E0EE&labelColor=302D41" />
    </a>

  <p align="center">
    <img src="https://stars.medv.io/nodetec/relaywizard.svg", title="commits"/>
  </p>

</div>

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
