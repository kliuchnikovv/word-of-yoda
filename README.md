Design and implement “Word of Wisdom” tcp server.
 • TCP server should be protected from DDOS attacks with the Prof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
 • The choice of the POW algorithm should be explained.
 • After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
 • Docker file should be provided both for the server and for the client that solves the POW challenge.

## Why I Chose Hashcash

Hashcash was selected as the proof-of-work algorithm for this system due to several key advantages:

### 1. Simplicity and Proven Reliability

Hashcash offers an elegant and battle-tested algorithm that has been in use for a long time. Its simplicity makes it easy to implement while still providing robust security guarantees. The algorithm requires finding a nonce that, when hashed with the challenge data, produces a hash with a specified number of leading zeros - a computation that is deliberately resource-intensive.

### 2. Asymmetric Computational Cost

One of Hashcash's primary benefits is the asymmetry between generation and verification:
- **Generation**: Computationally expensive, requiring numerous hash attempts
- **Verification**: Extremely efficient, requiring only a single hash operation

This asymmetry makes it ideal for preventing denial-of-service attacks while maintaining fast verification on the server side.

### 3. Adjustable Difficulty

Hashcash allows us to dynamically adjust the difficulty level by simply changing the required number of leading zeros in the hash. This adaptability enables system to scale according to client capabilities and current load conditions.

### 4. No Cryptographic Secrets Required

Unlike many security mechanisms, Hashcash doesn't rely on shared secrets or key management. This eliminates complex key distribution issues and reduces potential security vulnerabilities.

### 5. Stateless Verification

The verification process is completely stateless, which simplifies server architecture and allows for horizontal scaling without complex synchronization between nodes.

### 6. Proven Track Record

Beyond its original anti-spam application, Hashcash's principles form the foundation of Bitcoin's mining algorithm, demonstrating its effectiveness at scale in high-stakes environments.
