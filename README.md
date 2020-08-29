**Secure computation using Fully Homomorphic Encryption based on Lattice Based Cryptography.**

### Motivation

Often times, it is more convenient and cheaper to lease compute infrastructure than to build it. On the other hand, sensitive data
cannot always be handled by leased infrastructure due to legal requirements or other risk factors. In order
to get the best of both worlds, it would be ideal if the data could still be processed by a third party without being readable by the third party.
Fully Homomorphic Encryption is one way to satisfy these needs.


### Fully Homomorphic Encryption

Fully Homomorphic Encryption (FHE) allows for a server to perform computations (additions and multiplications) on encrypted data.
When the result of the calculations is sent back to the client after decryption, the client will have the same results as if the computations
had been performed on unencrypted data.
This feature allows for use cases where the client doesn't trust the server which performs the computations.

#### Lattigo

This project uses the Lattigo (https://github.com/ldsec/lattigo) Go library in order to implement FHE on a set of data sent
to the server for computation. In particular, it uses lattigo/bfv.


### Usage

This project consists of the server, performing computations on the encrypted data, and the client encrypting the data
and uploading it to the server. The client and the server communicate over gRPC.
Both, the client and the server can be compiled using `make` command. `make server` creates server binary in the build dir,
while `make client` creates client binary in the build dir. When the server binary is executed, the server will start listening
for incoming connections on the default port. When client binary is executed, it will connect the server on the default port and
send it instructions and data.

For demonstration purposes, the client can send a set of encrypted vectors to server. The server will store each vector in a file.
Upon request from the client, the server will perform an addition on stored vectors and return the encrypted result to the client.
The client will then decrypt the result and display it.
This shows that it is possible to use a server on untrusted infrastructure in order to execute computation without revealing the data
on which the computation is performed.

In order for the client to be able to successfully execute encryption/decryption, it will need to generate the keypair.
The client can be instructed to generate keypair by using the `g` switch like this `client -g`.
Once generated, the keypair will be used for the encryption/decryption so it is important not to lose it.

By using the `client -w` command followed by a set of integers, e.g. `client -w 324 54 66`, one can instruct the client to construct a new encrypted vector and
push it to the server for storage.

In order to receive the sum of all encrypted vectors stored on a server, one can execute `client -e`.

#### Configuration

It is possible to change port on which server listens to incoming connection as well as the token required for server to authorize operation.
Server reads it's settings from environment variables, while client reads it's settings from the `config.toml` file which should be
present in the same directory where the client binary resides.

##### Server configuration

Server accepts its configuration parameters via environment variables.

| variable        | description                                                                               | default   |
|-----------------|-------------------------------------------------------------------------------------------|-----------|
| SERVER_TOKEN    | Token string that needs to be supplied by client for operation to be authorized by server | 123       |
| SERVER_PORT     | IP and Port on which the server listens for incoming connections                          | :50051    |
| SERVER_FILE_DIR | The directory in which the incoming data will be stored                                   | /tmp      |

##### Client configuration

Client reads it's configuration from the `config.toml` file which should be present in the same dir from which the binary is run.
`config-example.toml` in examples directory shows the format of the configuration file.

| variable        | description                                                                               | default         |
|-----------------|-------------------------------------------------------------------------------------------|-----------------|
| token           | Token string that needs to be supplied by client for operation to be authorized by server | 123             |
| fhe_server      | IP and Port on which the server listens for incoming connections                          | 127.0.0.1:50051 |

##### Client command flags

`client -g` generetes new keypair used for encrytion/decryption on client side.

`client -w a b c ...` where a, b, c... are integers instructs client to create vector a,b,c... and send it to server.

`client -e` instructs server to perform evaluation on all stored data, fetches the result, decrypts it and shows it to user.



