<div align="center"><p>
    <h1>Relay Wizard 🧙</h1>
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
      <img alt="Repo size" src="https://img.shields.io/github/repo-size/nodetec/relaywizard?color=%23DDB6F2&label=SIZE&logo=codesandbox&style=for-the-badge&logoColor=D9E0EE&labelColor=302D41" />
    </a>
</div>

![0622](https://github.com/nodetec/relaywizard/assets/29136904/eb226b30-9250-43c6-ba2a-0361446d790b)

Relay Wizard is a CLI tool that helps you bootstrap a [Nostr](https://nostr.com/ "Nostr") relay.

The program will automate the following steps:

1. Install necessary dependencies
2. Set up a firewall
3. Configure nginx
4. Obtain a TLS certificate for HTTPS
5. Install the relay software
6. Set up a systemd service for your relay

## Installation

To install a relay, spin up a new Debian server, hook up a domain name, and run the following command:

```bash
curl -sL https://relaywizard.com/install.sh | bash
```

## Verification

If you prefer to manually verify the authenticity of the Relay Wizard binary before running it, then you can follow along with the verification process described here. This will minimize the possibility of the binary being compromised. To perform the verification you'll need to have `gnupg` and `curl` installed which are most likely already installed on your system, but if not here's how to install them on some operating systems:

### gnupg

#### Arch

```sh
sudo pacman -S gnupg
```

#### Debian/Ubuntu

```sh
sudo apt install -y gnupg
```

### curl

#### Arch

```sh
sudo pacman -S curl
```

#### Debian/Ubuntu

```sh
sudo apt install -y curl
```

Now you need to import the public key that signed the manifest file which you can do by running the following command:

```sh
curl https://keybase.io/nodetec/pgp_keys.asc | gpg --import
```

You're now ready to verify the manifest file. You will need to have the `rwz-x.x.x-manifest.sha512sum` and the `rwz-x.x.x-manifest.sha512sum.asc` files in the same directory as the Relay Wizard binary you downloaded where the `x.x.x` is replaced by whatever version of `rwz` you're verifying.

To verify the manifest file run the following command:

```sh
gpg --verify rwz-x.x.x-manifest.sha512sum.asc
```

Here's the command to run for the latest version of `rwz`:

```sh
gpg --verify rwz-0.3.0-alpha1-manifest.sha512sum.asc
```

You should see output similar to the following if the verification was successful:

```sh
gpg: assuming signed data in 'rwz-0.3.0-alpha1-manifest.sha512sum'
gpg: Signature made Thu Sep 26 21:04:47 2024 EDT
gpg:                using RSA key 252F57B9DCD920EBF14E6151A8841CC4D10CC288
gpg: Good signature from "NODE-TEC Devs <devs@node-tec.com>" [unknown]
gpg:                 aka "[jpeg image of size 5143]" [unknown]
Primary key fingerprint: 04BD 8C20 598F A5FD DE19  BECD 8F24 69F7 1314 FAD7
     Subkey fingerprint: 252F 57B9 DCD9 20EB F14E  6151 A884 1CC4 D10C C288
```

> Unless you tell GnuPG to trust the key, you'll see a warning similar to the following:

```sh
gpg: WARNING: This key is not certified with a trusted signature!
gpg:          There is no indication that the signature belongs to the owner.
```

This warning means that the key is not certified by another third party authority. If the downloaded file was a fake, then the signature verification process would fail and you would be warned that the fingerprints don't match.

When you get a warning like this it's also good practice to check the key against other sources, e.g., the [NODE-TEC Keybase](https://keybase.io/nodetec "NODE-TEC Keybase") or the [NODE-TEC GitHub](https://github.com/nodetec "NODE-TEC GitHub").

You have now verified the signature of the manifest file which ensures the integrity and authenticity of the file but not of the binary.

To verify the binary you'll need to recompute the SHA512 hash of the file, compare it with the corresponding hash in the manifest file, and ensure they match exactly which you can do by running the following command:

```sh
sha512sum --check rwz-x.x.x-manifest.sha512sum
```

Here's the command to run for the latest version of `rwz`:

```sh
sha512sum --check rwz-0.3.0-alpha1-manifest.sha512sum
```

If the verification was successful you should see the output similar to the following:

```sh
rwz-0.3.0-alpha1-x86_64-linux-gnu.tar.gz: OK
```

By completing the above steps you will have successfully verified the integrity of the binary.

## Learn more

If you want to learn more about how to setup a relay from scratch, check out [Relay Runner](https://relayrunner.org "Relay Runner").

If you just want to know enough to get started, read the following sections to get a server, hook up a domain name and setup remote access:

- [Get a server](https://relayrunner.org/server/get-a-server "Get a server")

- [Get a domain](https://relayrunner.org/server/domain-name "Get a domain")

- [Remote access](https://relayrunner.org/server/remote-access "Remote access")

from here you should be able to run the installation command above and get started.

## Contribute

If you want to contribute consider adding a new package manager and testing the script out on another Linux Distro. You can also look into adding support for more relay implementations.
