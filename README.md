## RPi Menagerie agent
Simple agent to register Raspberry Pi with Menagerie service

### To build on a RPi
Follow the instructions here to install and compile `go` on your RPi. In addition to the steps shown there, you have to also add the following to `.bashrc`:

```
#golang setup
export GOPATH=/root/CODE/GO/
export PATH=/usr/local/go/bin:$PATH
```

Now you can build your binaries on a RPi.

---

Menagerie RPi Agent:
It should create an entry on Menagerie for a RPi, storing:
    - IPs √
    - serial number √
    - MAC address(es) √

It should not create an UUID. The board has a unique serial number. We can send that and let menagerie create a UUID. We can then send the same UUID.

It should boot with the RPi, @onboot.
It should create $HOSTIP env var.
It should ensure $HOSTNAME env var exists.

It should provide a status and a health endpoint.
