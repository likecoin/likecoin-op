import type { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox-viem";

import "@nomicfoundation/hardhat-ethers";
import "@openzeppelin/hardhat-upgrades";
import dotenv from "dotenv";
dotenv.config();

const config: HardhatUserConfig = {
  solidity: {
    version: "0.8.28",
    settings: {
      optimizer: {
        enabled: true,
        runs: 200,
        details: {
          yulDetails: {
            optimizerSteps: "u",
          },
        },
      },
      viaIR: true,
    },
  },

  networks: {
    ...(process.env.DEPLOY_WALLET_PRIVATE_KEY_LOCALHOST != null
      ? {
          localhost: {
            url: "http://127.0.0.1:8545",
            accounts: [`0x${process.env.DEPLOY_WALLET_PRIVATE_KEY_LOCALHOST}`],
          },
        }
      : {}),
    ...(process.env.DEPLOY_WALLET_PRIVATE_KEY_OP_SEPOLIA != null
      ? {
          sepolia: {
            url: "https://ethereum-sepolia-rpc.publicnode.com",
            chainId: 11155111,
            accounts: [`0x${process.env.DEPLOY_WALLET_PRIVATE_KEY_OP_SEPOLIA}`],
          },
        }
      : {}),
    ...(process.env.DEPLOY_WALLET_PRIVATE_KEY_OP_SEPOLIA != null
      ? {
          "optimism-sepolia": {
            url: "https://sepolia.optimism.io",
            chainId: 11155420,
            accounts: [`0x${process.env.DEPLOY_WALLET_PRIVATE_KEY_OP_SEPOLIA}`],
          },
        }
      : {}),
  },
};

export default config;
