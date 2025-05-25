# INTRODUCTION
This is a command and control (C2) tool with communications over Tor network. Agent is made in Rust and C2 in Go.

C2 has two APIs; one for onion service listening on a unix socket, and another listening (by default) on port 8080 for connecting the user interface. C2's API for user interface integration is based on XENA's C2, offering a seemless integration with XENA's dark-themed, elegant user interface writen in Go, find out more at https://github.com/zarkones/XENA.

Agent does not need to "bundle" Tor service nor anything of a sort. Instead it relies on Arti, a rewrite of Tor in Rust. This allows easy compilation for many targets and way less trouble compared to trying to embed Tor writen in C into an agent. Arti does not offer all of the security features of legacy Tor implementation, however, that doesn't matter since you can (and should for the time being) run the onion service via the legacy battle tested Tor implementation.

Note that this is experimental and ongoing development effort.

# FEATURES
- Tor integration (allows for end to end encryption, hiding the C2's IP address)
- Execution of shell commands.
- Basic attempt of hiding C2 config in the agent's binary.

# SETUP
This guide assumes you have Go, Rust, and XENA ready.

- Run the C2 server via: cd api && go run . --api-key=your_secret_api_key_longer_than_16_chars
- Run Tor onion service based on the configuration made by the C2 program via: cd api && tor -f torrc
- When Tor onion service is ready you will see a .onion domain inside: api/onionservice/hostname, place that domain in function "get_address" in agent/src/config.rs
- Run the agent via: cd agent && cargo run

That's it. If you wish to know more about it execute the C2 server with "-h" command for help, or read the source code. It's less than 1,000 lines of code. A good starting point for further development.

![Promo Image 1](https://raw.githubusercontent.com/zarkones/OnionC2/production/promo/promo1.png)
