# Web3 SSO

## Supported/Tested wallets
- Metamask

## Authentication steps
1. Connect wallet
2. Issue challenge
3. Login

### Connect wallet
1. Client connects to wallet
2. Client gets wallet address

### Issue challenge
1. Client calls /challenge with wallet address in the params
2. Server generates a challenge and sends it to the client

### Login
1. Client signs challenge
2. Client calls /login with wallet address and signed challenge
3. Server verifies signed challenge was signed by wallet address
4. Server generates jwt and sends it to the client

## Languages
- JavaScript
- More in the future...

## SDKs and Examples
All the SDKs have a corresponding example. All examples use the same frontend with a JWT cookie for authentication.
To run each example, enter the corresponding example, install dependencies

### How to run example?
- Go to the example
- Install dependencies
- Run bin/start script

### Examples
- Node.js -> examples/node
- Express.js -> examples/express