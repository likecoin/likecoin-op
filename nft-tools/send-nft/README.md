# send NFT

send NFT from csv

## Prerequisite

- Docker / Node.js 14
- env

## Installation

```bash
$ make npm install
```

## Configuration

Change the data in `./list.csv` (reference `./list_example.csv`)

Change the data in `./config/config.js`

DEFINE `MNEMONIC` or `PRIVATE_KEY`

Where `MNEMONIC` is your mnemonic phrase, `PRIVATE_KEY` is the hex of private key, `WAIT_TIME` is wait time after console (before send)

## Usage

```bash
$ make send-nft
```
