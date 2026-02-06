import "@openzeppelin/hardhat-upgrades";
import hardhat, { ethers } from "hardhat";
import { deployWithGasEstimation } from "../utils/estimateGas/estimateGas";
import makeDefaultInterceptors from "../utils/estimateGas/interceptors";

async function main() {
  const BookNFTVoid = await ethers.getContractFactory("BookNFTVoid");
  const [deployer] = await ethers.getSigners();
  console.log("Deployer:", deployer.address);
  const gasBump = parseInt(process.env.GAS_BUMP) || 300000;
  console.log("Deploying with gas limit bump", gasBump);
  const feeIncrement = parseInt(process.env.FEE_BUMP) || 0;
  console.log("Deploying with fee per gas bump", feeIncrement);

  const bookNFTVoid = await deployWithGasEstimation(hardhat, BookNFTVoid, [], {
    deployer,
    interceptor: makeDefaultInterceptors("deploy", gasBump, feeIncrement),
  });

  const newImplementationAddress = await bookNFTVoid.getAddress();
  console.log(
    "New BookNFTVoid implementation is deployed to:",
    newImplementationAddress,
  );

  // Too many time the block-explorer not yet catch the contract, not calling here.
  console.log(
    "Run following to verify after block-explorer catch the deployment",
  );
  console.log(`
VERIFY=etherscan BOOKNFTVOID_ADDRESS=${newImplementationAddress} \\
    npm run script:${hardhat.network.name} scripts/verifyBookNFTVoid.ts`);
  console.log(`
VERIFY=blockscout BOOKNFTVOID_ADDRESS=${newImplementationAddress} \\
    npm run script:${hardhat.network.name} scripts/verifyBookNFTVoid.ts`);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
