
&ensp;& ensp;&ensp; This program is a blockchain public chain demo designed to mimic the functionality of Bitcoin. It mainly applies knowledge related to cryptography, consensus algorithms, peer-to-peer networks, and blockchain tamper proof structures. It combines various knowledge points together to create a simple and complete executable public chain demo

<hr>

###Program features:

-Based on the proof of work consensus algorithm, data is stored in the structure of blockchain
-Decentralization, utilizing P2P technology to make each node relatively independent of each other
-Proactively search for peer nodes in the network, automatically connect and store them in the local node pool
-When a node exits, it will broadcast to the entire network, and the remaining nodes will dynamically update the current pool of connectable nodes
-Successful mining nodes obtain accounting rights and broadcast the latest synchronized blocks to the entire network. After verification by other nodes, they are stored in the local blockchain
-The transaction transfer uses the UTXO transaction model, which supports multiple transfers in one transaction
-Support importing Chinese mnemonic words and generating public-private key pairs from mnemonic words (using elliptic curve algorithm)
-Transaction transfers use private keys for digital signatures, public key verification, and the UTXO structure avoids replay attacks on signatures
-Establish a separate data table for unused UTXO and optimize transfer transaction speed
-Use Merkle tree to generate the root hash of transactions (the current demo does not distinguish between block headers and blocks, only wants to practice using this data structure)
-Persistent blockchain and public-private key information, stored in the local database of each node (each node has its own independent database)
-Customize mining difficulty value and absenteeism mining reward value
-Customize the size of the trading pool, mining will only begin after a specified number of transactions are completed

<hr>

###Main modules:
-Command scheduling module
-UTXO transaction generation module
-Cryptography encryption and decryption module
-Block generation and verification module
-Data persistence module
-P2P network communication module
-Log output module
<br>
 
####Command scheduling module
&ensp;& ensp;&ensp;   After starting the program, the console captures user input information and parses commands and values following the commands based on the user's input. Perform relevant operations on the program according to different commands
<br>
####UTXO transaction generation module
&ensp;& ensp;&ensp; The transaction transfer module is based on the UTXO model, but does not introduce Bitcoin scripts. Instead, a digitally signed byte array is directly used as a substitute in the script. When user A transfers money to user B, user A needs to use their private key to digitally sign the "* * input * *" (which includes the transaction hash, index, and other information owned by user A), generate the transaction, and send it to other nodes. Other nodes then use user A's public key to verify the signature</ br>&ensp;&ensp;&ensp; Due to the special structure of UTXO, it naturally avoids replay attacks and does not require adding nonce values like the Ethereum account system. However, in order to avoid the problem of repeated calculation of UTXO, it is not supported to transfer the same address again before the previous transfer is packaged into the block& ensp;&ensp; Supports multiple transfers for one transaction, and has created a UTXO data table specifically for storing all unused * * outputs * * in the blockchain to optimize transfer query speed</ br>&ensp;&ensp;&ensp; To learn more about UTXO, it is recommended to refer to [this article]（ https://draveness.me/utxo-account-models )

>Because the selected elliptical curve is ECDSA, there may be scalability attacks in the digital signature section, which requires isolation verification. Students with energy can do it themselves

<br>

####Cryptography encryption and decryption module
1. Unidirectional hash function: sha256 ripemd-160
   
Mainly used to convert the entire block into a fixed length string through computation, facilitating data verification
2. Encoding and decoding algorithm: base58
   
Due to the excessively long original length of the private key, which is not conducive to memory, base58 encoding is used to visually encode the private key and address
3. Asymmetric encryption: elliptic curve algorithm (crypto/elliptic p256)
   
Extract 7 pairs of Chinese words as seeds through mnemonic text, generate public-private key pairs using elliptic curve algorithm, use the private key for digital signature of transaction data, and verify the signature with the public key to ensure the identity of the initiator.
The public key generates an address through a series of operations, which is used to query the balance and receive transfer tokens
&ensp;& ensp;&ensp; The address generation rules are as follows:
    
-Generate public key through elliptic curve algorithm
-Hash the public key with sha256 and ripemd160 to obtain the publickeyHash
-Add a version byte array before publickeyHash to obtain version PublickeyHash
-Perform two sha256 hashes on version PublickeyHash and take the first 4 bytes to obtain tailfHash
-Concatenate tailfHash after version PublickeyHash to obtain the final hash of the public key, which is the finalHash
-Finally, perform Base58 encoding on the finalHash to obtain the Bitcoin address

  
>There used to be a question about why generating addresses in Bitcoin is so complicated. Since asymmetric encryption only has a public key and cannot deduce a private key, why not directly use the public key as the address, but hash the public key multiple times to obtain the address? It wasn't until I recently read an article that I realized that quantum computers can crack elliptic curve encryption and quickly find private key information through the public key. However, quantum computers are difficult to reverse Hash algorithms (or require 2-80 steps to crack Hash), so placing your Bitcoin in an unpaid address (according to the UTXO transaction model, the output stores the public key Hash instead of the public key, which also explains why the UTXO input stores the public key and the output stores the public key Hash) is quite secure. That is to say, addresses that have already been spent are not secure in front of quantum computers, while addresses that have not been spent have strong quantum resistance.
####Block generation and verification
&ensp;& ensp;&ensp;   Based on the POW consensus algorithm to generate blocks, the mining difficulty (a string of large numbers) is first defined according to the difficulty value (which can be defined in the configuration file). By calling GO's own random number packet crypto/rand, the random number nonce is continuously transformed (the previous version used the method of accumulating nonce values, but the probability of branching is too high), and the block itself is continuously hashed to make the final calculated block hash value smaller than the currently defined mining difficulty, in order to obtain the right to extract the block</ br>&ensp;&ensp;&ensp;   Block nodes can receive reward tokens and have accounting rights, which will be broadcasted across the entire network after block generation. After receiving the block, the remaining P2P nodes first verify the block's own hash, then check whether the pre hash in the block is consistent with the local previous block hash, and finally store it in the local database.
<br>
####Data persistence module
&ensp;& ensp;&ensp;  The persistence layer is based on the KV type database blot and includes an additional layer of encapsulation, with the main interfaces being put, view, and delete. Each call to the interface will open and close the database handle separately, so there will be no situation where it is occupied by other threads. The database has established three tables: BlockBucket (for storing detailed information about blocks), AddrBucket (for storing local wallet data), and UTXOBucket (for storing unconsumed UTXO data)
<br>
####P2P network communication module
&ensp;& ensp;  Using MDNS technology suitable for local area network addressing, due to the bug that the package used cannot find the network in Windows, it is recommended to run this program on Linux/MAC</br>&ensp& ensp;  After the node is started, it will automatically search for other peer nodes in the local area network. Once discovered, it will be stored in the node pool (stored in memory). The first twelve bytes of data communicated between nodes are defaulted as commands, and feedback on local blockchain related information will be provided based on different commands& ensp;&ensp; The main operating principle is to distribute blocks and mine after receiving transactions:</br>
</br>&ensp;& ensp; Process of obtaining blocks:
1. Compare block heights with each other
2. Obtain the missing block hash
3. Receive missing entire blocks through block hashing
4. Block verification, stored in the database

&ensp;& ensp; Mining process:
1. Send transaction data to all network nodes through a certain node
2. The node receives the transaction and performs signature verification and balance verification on the transaction
3. After verification, deposit into the trading pool and start mining once the size of the trading pool is met
4. Successful mining, broadcasting block height across the entire network
5. Send blocks to other nodes
6. Other nodes perform block verification and store it in the database
####Log output module
&ensp;& ensp;&ensp;  Using a self-made log package, the program will generate log files with log and port numbers by default in the current directory (which can be set in the configuration file) after startup. All debug information generated by the program will be printed in this log file. It is recommended to open a window for real-time monitoring to facilitate the interaction between nodes and the detailed steps of block generation

[Characteristics of Log Package]:

-Support directed output of logs to specified files
-Support one click hiding of debugging information
-Supports color printing (both Windows/Linux/Mac support it)
-Display the class name, function/method name of the output log
 
<br>
<hr>

Main toolkits used
---------------------------
Package | Purpose
-------- | -----
[github.com/boltdb/bolt]( https://github.com/boltdb/bolt ）|K, v type database
[github.com/spf13/viper]( https://github.com/spf13/viper ）| Configuration file reading tool
[github.com/golang/crypto]( https://github.com/golang/crypto ）|Cryptography related tools
[github.com/libp2p/go-libp2p]( https://github.com/libp2p/go-libp2p ）P2P communication tool under IPFS
[github.com/corgi-kx/logcustom]( https://github.com/corgi-kx/logcustom ）|Log output tool

<hr>
 
###Program running tutorial:

**1. Compile after downloading**

This demo is recommended to be run on Linux/Mac, otherwise there may be issues with mnemonic garbled characters and inability to find peer-to-peer networks

```shell
git clone  https://github.com/corgi-kx/blockchain_golang.git
```
```shell
go build -mod=vendor -o chain main.go
```
<br>

**2. Open multiple windows**

To simplify the operation, start different ports on the same computer to simulate P2P nodes (three windows for program startup and three windows for real-time log viewing)
>When operating on a real machine, if no other nodes can be found, it may be a firewall issue. Please turn off the firewall and try again

! [Insert image description here]（ https://img-blog.csdnimg.cn/20191118103707708.png )

<br>

**3. Modify the configuration file**
  
Mainly modify the local listening IP and local listening port. Other defaults are sufficient</br>
It is not recommended to lower the difficulty threshold to avoid block forks. The demo has not yet processed block forks
```shell
vi config.yaml
```
```yaml
blockchain:
#Mining difficulty value, the higher the difficulty, the harder it is to mine
mine_difficulty_value: 24
#Number of mining reward tokens
token_reward_num: 25
#Transaction pool size (how many transactions must be met before mining begins)
trade_pool_length: 2
#Log storage path
log_path: "./ "
#Chinese mnemonic word seed path
chinese_mnemonic_path: "./ chinese_mnemonic_world.txt"
network:
#Local monitoring IP
listen_host: "192.168.0.164"
#Local listening port
listen_port: "9000"
#Unique identifier name for node group (if the names between nodes are different, the network cannot be found)
rendezvous_string: "meetme"
#The protocol ID of the network transport stream (if the IDs between nodes are different, data cannot be sent)
protocol_id: "/chain/1.1.0"

```

<br>

**4. Start the node, create a wallet, and generate a genesis block**

Start Node 1
```shell
./chain
```
! [Insert image description here]（ https://img-blog.csdnimg.cn/20191118101305498.png?x -oss-process=image/watermark, type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzM1OTExMTg0,size_16,color_FFFFFF,t_70)

By command, sir, create three wallet addresses

```
> generateWallet
Mnemonic words: ["Lung segment", "habitat", "tooth groove", "several dimensions", "Chinese Portuguese", "Mangyu", "Guanghua"]
Private key: 6HrLjHE4Qm31dZFGjemwNLZM3iqnxoSUqKb5VtEKbWzh
Address: 12BwtcVWimms9rrKxxoDev68woGyMYS4sk
> generateWallet
Mnemonic words: ["sprain", "cut wound", "myopathy", "sinking", "generalized", "voiced", "hernia"]
Private key: 7yBRSB46q8ZeEiYbwZDSvKzzsh1MYAygeo2i689uEMAf
Address: 1B6KYdABXZDwq8xGTbdKnpHBo11CkihxS
> generateWallet
Mnemonic words: ["ventricle", "deficiency vessel", "flap stomach", "black tea", "share", "Zhang copper", "wandering"]
Private key: 872CCeLS8bDrC7bdSoFrgUSWm57eqTdypEhKbErYC9xi
Address: 1E6aRBxfncAsypUnjGxPJYBR4J3gZ6hHD
```
Generate genesis block (assign the first address 100 Tokens)
```
> genesis -a 12BwtcVWimms9rrKxxoCev68woGyMYS4sk -v 100
Already generated genesis block
```

Log 1: Real time viewing of logs (showing the mining process)
```shell
tail -f log9000.txt 
```
! [Insert image description here]（ https://img-blog.csdnimg.cn/20191118144251486.png?x -oss-process=image/watermark, type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzM1OTExMTg0,size_16,color_FFFFFF,t_70)

<br>

**5. Synchronize blocks**

Node 2 and Node 3 sequentially modify the port number of the configuration file to 90019002, and start these two nodes to synchronize the genesis block</br>
At this point, the log of node 1 detects the presence of other nodes in the network
! [Insert image description here]（ https://img-blog.csdnimg.cn/20191118145703154.png ）Node 2 and Node 3 will automatically synchronize the genesis block after startup
! [Insert image description here]（ https://img-blog.csdnimg.cn/20191118145752942.png?x -oss-process=image/watermark, type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzM1OTExMTg0,size_16,color_FFFFFF,t_70)

<br>

**6. Perform transfer operation**

Each node is assigned a mining reward address (which can also be left unspecified, as no reward will be generated after the node mines)</br>
Node 1 sets mining reward address:
```
> setRewardAddr -a 12BwtcVWimms9rrKxxoCev68woGyMYS4sk
The address 12BwtcVWimms9rrKxxoDev68woGyMYS4sk has been set as the mining reward address!
```
Node 2 sets mining reward address:
```
> setRewardAddr -a 1B6KYdABXZDwq8xGTbdDknpHBo11CkihxS
The address 1B6KYdABXZDwq8xGTbdknpHBo11CkihxS has been set as the mining reward address!
```
Node 3 sets mining reward address:
```
> setRewardAddr -a 1E6aRBxfncAsypUnjGxPJYbR4JQ3gZ6hHD
The address 1E6aRBxfncAsypUnjGxPJYbR4JQ3gZ6hHD has been set as the mining reward address!
```
Node 1 performs a transfer operation (the genesis address transfers 10 tokens each like the other two addresses)
```
> transfer -from ["12BwtcVWimms9rrKxxoCev68woGyMYS4sk","12BwtcVWimms9rrKxxoCev68woGyMYS4sk"] -to ["1B6KYdABXZDwq8xGTbdDknpHBo11CkihxS","1E6aRBxfncAsypUnjGxPJYbR4JQ3gZ6hHD"] -amount [10,10]
The transfer command has been executed
```
! [Insert image description here]（ https://img-blog.csdnimg.cn/2019111815314125.png?x -oss-process=image/watermark, type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzM1OTExMTg0,size_16,color_FFFFFF,t_70)

<br>

**7. Check balance**

Among the three nodes, if node 2 mines a block, node 2 should receive a mining reward of 25 Tokens
! [Insert image description here]（ https://img-blog.csdnimg.cn/20191118153547470.png ）At this point, type the 'getBalance' command at any node to view the balance information of three addresses
```
> getBalance -a 12BwtcVWimms9rrKxxoCev68woGyMYS4sk
Address: 12BwtcVWimms9rrKxxoCev68woMYS4sk Balance: 80
> getBalance -a 1B6KYdABXZDwq8xGTbdDknpHBo11CkihxS
Address: 1B6KYdABXZDwq8xGTbdknpHBo11CkihxS Balance: 35
> getBalance -a 1E6aRBxfncAsypUnjGxPJYbR4JQ3gZ6hHD
Address: 1E6aRBxfncAsypUnjGxPJYBR4J3gZ6hHD Balance: 10
```

<br>

**8. View detailed block information**

Enter the ` printAllBlock ` command at any node to view block information

Block 1 is the genesis block, with only a 100UTXO output assigned to '12BwtcVWimms9rrKxxoDev68woGyMYS4sk'

You can see block 2:</br>
The first transaction, address' 12BwtcVWimms9rrKxxoCev68woGyMYS4sk ', first spends the UTXO with a genesis block quota of 100, generates a 90UTXO for itself, and generates a 10UTXO for address' 1B6KYdABXZDwq8xGTbdKnpHBo11CkihxS'</br>
The second transaction address' 12BwtcVWimms9rrKxxoCev68woGyMYS4sk 'uses the 90 limit UTXO output from the first transaction to generate an 80UTXO for itself and a 10UTXO for address' 1E6aRBxfncAsypUnjGxPJJybR4JQ3gZ6hHD'</br>
The third transaction is a mining reward transaction, so there is only output and no input. Generate 25UTXO for the address' 1B6KYdABXZDwq8xGTbdKnpHBo11CkihxS' (the 25 reward limit set in the configuration file)

```
> printAllBlock                                   
========================================================================================================
This block hash is 00000008acfb9a8dcf3b923f4eb6f2ddfc27dcaff861ea6848a9074ca46d85b
------------------------------Transaction data------------------------------
This transaction ID: 988cbe7f374855aa94addb873f22960cf43646bdaeb56253f3e683478270db
tx_input：
Transaction ID: bb717bd6717c8cae3829875187b97f256859277ad4a52ac57cdbc132895ca154
Index: 0
Signature information: 8c8b0628headebbc9e97b490a40a23494d3f8286f1af045f1e18d529c49a90afa194799182c264ee15871b5dd35c773e5dd46427fc8e2C268356ce09f6b60b
Public key: 8e0f1fe7d6177f11027818663048392ce8952ceff1ceec8edc84e176f46cedd338575f709b412eeab904d7027056354038f8aef7a1940f45264f7116ba793
Address: 12BwtcVWimms9rrKxxoDev68woGyMYS4sk
tx_output：
Amount: 90
Public key Hash: 0d0a1eb1baf838828a54ac97b09524f0b0c3210
Address: 12BwtcVWimms9rrKxxoDev68woGyMYS4sk
---------------
Amount: 10
Public key Hash: 6eb2d1846217aa089dfa26e3147b767e1 de0b08d
Address: 1B6KYdABXZDwq8xGTbdKnpHBo11CkihxS
This transaction ID: 443b4a4f04204bd8ed2bdcc096642a27457c27aa47c2ee81486d7440b05959521
tx_input：
Transaction ID: 988eccbe7f374855aa94addb873f22960cf43646bdaeb56253f3e683478270db
Index: 0
Signature information: 2a064297227ba07c7ea92eebb1d43f3fe4dfbd6c7e78be8ec2d30e20fa51500c8bdb591c11908a877aeef61b4c64f9a851c44af441cbe6893e1b80e42032c
Public key: 8e0f1fe7d6177f11027818663048392ce8952ceff1ceec8edc84e176f46cedd338575f709b412eeab904d7027056354038f8aef7a1940f45264f7116ba793
Address: 12BwtcVWimms9rrKxxoDev68woGyMYS4sk
tx_output：
Amount: 80
Public key Hash: 0d0a1eb1baf838828a54ac97b09524f0b0c3210
Address: 12BwtcVWimms9rrKxxoDev68woGyMYS4sk
---------------
Amount: 10
Public key Hash: 8fa79c32a067830be3b16ade637d370e1d1e6e0d
Address: 1E6aRBxfncAsypUnjGxPJYBR4J3gZ6hHD
This transaction ID: 2420c67272ab7832d6148a36b38166862d12e265f184439e6ab2e606b01245
tx_input：
tx_output：
Amount: 25
Public key Hash: 6eb2d1846217aa089dfa26e3147b767e1 de0b08d
Address: 1B6KYdABXZDwq8xGTbdKnpHBo11CkihxS
--------------------------------------------------------------------
Timestamp 2019-11-18 03:23:57 PM
Block height 2
Random number 2808567053068705071
Previous block hash 00000 7d7B7c7B54d9d1b0d1d06b6936e1bc613f6ab7de1ae0275cdaef4e4a4
========================================================================================================
This block has a hash of 00000 7d7b7c5b540d9d1b0d1d06b6936e1bc613f6ab7de1ae0275cdaef4e4a4
------------------------------Transaction data------------------------------
This transaction ID: bb717bd6717c8cae3829875187b97f256859277ad4a52ac57cdbc132895ca154
tx_input：
Transaction ID:
Index: -1
Signature information:
Public key:
Address:
tx_output：
Amount: 100
Public key Hash: 0d0a1eb1baf838828a54ac97b09524f0b0c3210
Address: 12BwtcVWimms9rrKxxoDev68woGyMYS4sk
--------------------------------------------------------------------
Timestamp 2019-11-18 10:43:41 AM
Block height 1
Random number 86040767999888393002
Previous block hash: 0000000 0000000
========================================================================================================

```

<br>

**9. Other**

You can also initiate a transfer at nodes 2 and 3, but first you need to import wallet information through mnemonic words, as shown in the following example:
```
>ImportMnword-m ["sprain", "cut wound", "myopathy", "sagging", "generalized", "voiced", "hernia"]
```
<br>
Please explore more features on your own:)

<br>
 


