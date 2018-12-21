# debased
## Overview  
debased aims to combine the power of blockchain and distributed systems with clean and robust development tools in order to provide a fast, reliable, and secure database solution to everyone.
## Philosophy
Data should be secure, accessible, and uncensored.
## The debased System
### Transactions
All communication between nodes on the debased system are transactions. Each transaction consists of instructions to be executed, a bounty, and a signature. Transactions range from moving funds to managing your database and adding/editing/deleting data.
***
All "edited" and "deleted" data is still stored, but moved to a smaller data footprint with longer look up times. These records will not be returned unless explicitly requested.
***
### Decentralized
debased has no governing authority. The system self-regulates, as long as half of the coins in the consensus process approve a transaction the transaction is ratified.
### Mutual Distrust
The debased system removes the need for trust. There is ZERO chance of down time, data loss, or data breach. Everything stored on the debased system is always available and secure.
### Blockchain
debased uses blockchain technology to store all transactions that occur on the system. This guarantees all transactions on debased are indestructible and immutable.
### Proof of Stake
The system comes to consensus using the computationally cheap proof of stake system instead of the extremely intensive proof of work algorithm. This minimizes energy usage and computational complexity without sacrificing security. This comes at the cost of higher network traffic but this trade off is definitively worthwhile.
Proof of stake works having every potential update to the blockchain sent to the other nodes on the network. The system then begins a vote, in which each node can place a bet with their position (pass/fail) and their confidence (how many coins they are willing to risk). This allows the system to come to consensus without requiring pure hashing power à la bitcoin and ethereum.
### Elliptic Curve Signatures
The core cryptographic strength of debased comes from  the elliptic curve digital signature algorithm. To use, a user generates a private key (from this private key the public key and account ID can be generated). This private key should **NEVER BE SHARED**. This key is used to sign transactions. The signing process creates a sha hash of all the data in the transaction returning two values (r and s). These two values are passed along with the transaction and the public key. These three values can then be used to verify the transaction is from the holder of the private key.
### Confidence
If a node repeatedly sends wrongly signed messages the receiving nodes have no hard evidence to present to the debased system to take their collateral (this is because an incorrectly signed message can be forged). In this case, the receiving node can loses confidence in the node. This causes communication from the node to be lowered in priority. If the offending node drops low enough the receiving node can start a which hunt.
### Which Hunt
This is a voting process used to kick a node from the network and take its collateral. The vote passes if half the coins in the vote process approve the measure. This means that a node will get kicked if it acts too poorly, even if a node remains below this threshold they will still lose value compared to the well behaved nodes given the confidence system.
### Double Spending
To avoid double spending, debased continues to break from bitcoin and ethereum. Instead of using a transaction counter for each account, debased uses liquid and illiquid balances. Whenever a transaction is signed by an account the collateral and max cost allowed is moved to the illiquid balance. After the transaction is completed the unused amount is returned from illiquid to liquid. If the transaction fails to conform to debased rules, the collateral is taken from the illiquid balance and given to the node that found the issue and all the nodes that were forced to use resources on an invalid transaction.
