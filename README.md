### SSH Tester

#### Background 

Have you noticed that when you ssh into a server, you  have to enter the key passphrase _only_ for the key
that is actually configured on that server, but not other keys?

This is possible because the client sends over a selection of keys (pubkeys), and the server
informs the client which keys are acceptable. In consequence, it is possible to, given an ssh public key, figure
out if the key is configured for access to the server.

This repo is made as a little investigation into how this works. I did not find any existing
easy way to do it via the standard linux `ssh` binary, so I instead implemented it in golang. 

This repo is, essentially, a copy-pasta of the golang `x/crypto/ssh` package, with a lot of internals
gradually gutted out and removed. 

### What does it do

So, you can give it 

- a file with public keys (on the same format as a `authorized_keys`-file).
- a set of usernames,
- a host (and a port)

It will now,

- For each pubkey,
   - For each username,
     - Let you now if the username/pubkey combo is acceptable at the server.

### What does this mean

From the 'blue' perspective: within an organization, one can use this to check whether a user has been successfully disenrolled from
a set of servers, e.g. after a user has left the organization. All you need is the old `authorized_keys`-file containing the pubkey,
(and the username), and then without actually having any authorization on the server(s) in question, you can check if the user
still has access.

From a 'red' perspective, one can traverse all users in a github organization and download all the user ssh pubkeys.
A github user ssh key is available at e.g. [https://github.com/holiman.keys](https://github.com/holiman.keys).
And one can then visit all the organization's public servers, and see who has access where. 

### Example

This will probably bitrot, but here's how it looks right now

```
$ go run ./cmd/sshx -host 192.168.197.323 -keyfile ./pubkeys.txt --user root,foo
2025/03/28 09:12:04 INFO Testing user=root pubkey="ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAINUYPLDQQxX1GGXQZR3Pj/Xy7ckadVdiG9R61bVlq4oP"
2025/03/28 09:12:04 INFO Testing user=foo pubkey="ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAINUYPLDQQxX1GGXQZR3Pj/Xy7ckadVdiG9R61bVlq4oP"
2025/03/28 09:12:04 INFO Testing user=root pubkey="ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFKF5kDRXv4SWBYrk36i4iLRl2BZG3ESMQjMLsUpiHz5"
2025/03/28 09:12:05 INFO Server accepted user/key user=root pubkey="ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFKF5kDRXv4SWBYrk36i4iLRl2BZG3ESMQjMLsUpiHz5"
2025/03/28 09:12:05 INFO Testing user=foo pubkey="ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFKF5kDRXv4SWBYrk36i4iLRl2BZG3ESMQjMLsUpiHz5"
```