# experiment-p2p

# Build/Install
`make`

# Run
## server
`./gop2p server`

## client A 
`./gop2p client -t A -a <ip of server>`

## client A 
`./gop2p client -t B -a <ip of server>`

After the p2p connection is built, the server will exit. The two client will send a random int [0,100) to each other

> **note**
> the default <ip of serrver> is my vps' ip

![explain](./explain.png)
ref: https://blog.nowcoder.net/n/4cd4fd71452c40c18cdfeb92383117e5
