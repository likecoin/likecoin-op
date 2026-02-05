import { expect } from "chai";
import { BaseContract } from "ethers";
import { ethers, upgrades } from "hardhat";

import { BookConfigLoader, BookTokenConfigLoader } from "./BookConfigLoader";
import { createProtocol } from "./ProtocolFactory";

describe("BookNFTVoid", () => {
  before(async function () {
    this.LikeProtocol = await ethers.getContractFactory("LikeProtocol");
    this.BookNFTVoid = await ethers.getContractFactory("BookNFTVoid");
    const [protocolOwner, classOwner, likerLand, randomSigner] =
      await ethers.getSigners();

    this.protocolOwner = protocolOwner;
    this.classOwner = classOwner;
    this.likerLand = likerLand;
    this.randomSigner = randomSigner;
  });

  let protocolContract: BaseContract;
  let nftClassId: string;
  let nftClassContract: BaseContract;
  let voidContract: BaseContract;
  let voidClassContract: BaseContract;
  let originalName: string;
  let originalSymbol: string;
  let originalMaxSupply: bigint;
  let originalContractURI: string;
  let mintedTokenURI: string;

  beforeEach(async function () {
    const { likeProtocolContract } = await createProtocol(this.protocolOwner);

    protocolContract = likeProtocolContract;

    const likeProtocolOwnerSigner = protocolContract.connect(
      this.protocolOwner,
    );

    const NewClassEvent = new Promise<{ id: string }>((resolve, reject) => {
      likeProtocolOwnerSigner.on(
        "NewBookNFT",
        (id: string, params: any, event: any) => {
          event.removeListener();
          resolve({ id });
        },
      );
      setTimeout(() => {
        reject(new Error("timeout"));
      }, 20000);
    });

    const bookConfig = BookConfigLoader.load(
      "./test/fixtures/BookConfig0.json",
    );

    await likeProtocolOwnerSigner.newBookNFT({
      creator: this.classOwner,
      updaters: [this.classOwner, this.likerLand],
      minters: [this.classOwner, this.likerLand],
      config: bookConfig,
    });

    const newClassEvent = await NewClassEvent;
    nftClassId = newClassEvent.id;
    nftClassContract = await ethers.getContractAt("BookNFT", nftClassId);
    expect(await nftClassContract.owner()).to.equal(this.classOwner.address);

    // Mint a token before upgrading to void
    const likeClassOwnerSigner = nftClassContract.connect(this.classOwner);
    const tokenMetadata = JSON.stringify({
      image: "ipfs://QmUEV41Hbi7qkxeYSVUtoE5xkfRFnqSd62fa5v8Naya5Ys",
      description: "Test NFT Description",
      name: "Test NFT #1",
    });
    mintedTokenURI = tokenMetadata;

    await likeClassOwnerSigner
      .mint(this.classOwner.address, ["test mint"], [tokenMetadata])
      .then((tx: any) => tx.wait());

    // Store original values before upgrade
    originalName = await nftClassContract.name();
    originalSymbol = await nftClassContract.symbol();
    originalMaxSupply = await nftClassContract.maxSupply();
    originalContractURI = await nftClassContract.contractURI();

    // Deploy BookNFTVoid implementation
    const bookNFTVoidOwnerSigner = this.BookNFTVoid.connect(this.protocolOwner);
    voidContract = await bookNFTVoidOwnerSigner.deploy();

    // Upgrade to void implementation
    await likeProtocolOwnerSigner.upgradeTo(voidContract.getAddress());

    // Get the void contract at the proxy address
    voidClassContract = await ethers.getContractAt("BookNFTVoid", nftClassId);
  });

  describe("Interface Detection", () => {
    it("should return false for ERC721 interface (0x80ac58cd)", async function () {
      expect(await voidClassContract.supportsInterface("0x80ac58cd")).to.equal(
        false,
      );
    });

    it("should return false for ERC721Metadata interface (0x5b5e139f)", async function () {
      expect(await voidClassContract.supportsInterface("0x5b5e139f")).to.equal(
        false,
      );
    });

    it("should return false for ERC721Enumerable interface (0x780e9d63)", async function () {
      expect(await voidClassContract.supportsInterface("0x780e9d63")).to.equal(
        false,
      );
    });

    it("should return false for ERC2981 interface (0x2a55205a)", async function () {
      expect(await voidClassContract.supportsInterface("0x2a55205a")).to.equal(
        false,
      );
    });

    it("should return false for ERC165 interface (0x01ffc9a7)", async function () {
      expect(await voidClassContract.supportsInterface("0x01ffc9a7")).to.equal(
        false,
      );
    });

    it("should return false for random interface", async function () {
      expect(await voidClassContract.supportsInterface("0x12345678")).to.equal(
        false,
      );
    });
  });

  describe("ERC721 Functions - All Should Revert", () => {
    it("balanceOf should revert", async function () {
      await expect(
        voidClassContract.balanceOf(this.classOwner.address),
      ).to.be.rejectedWith("ErrVoidContract()");
    });

    it("ownerOf should revert", async function () {
      await expect(voidClassContract.ownerOf(0)).to.be.rejectedWith(
        "ErrVoidContract()",
      );
    });

    it("tokenURI should revert", async function () {
      await expect(voidClassContract.tokenURI(0)).to.be.rejectedWith(
        "ErrVoidContract()",
      );
    });

    it("name should revert", async function () {
      await expect(voidClassContract.name()).to.be.rejectedWith(
        "ErrVoidContract()",
      );
    });

    it("symbol should revert", async function () {
      await expect(voidClassContract.symbol()).to.be.rejectedWith(
        "ErrVoidContract()",
      );
    });

    it("transferFrom should revert", async function () {
      await expect(
        voidClassContract.transferFrom(
          this.classOwner.address,
          this.randomSigner.address,
          0,
        ),
      ).to.be.rejectedWith("ErrVoidContract()");
    });

    it("safeTransferFrom (3 args) should revert", async function () {
      await expect(
        voidClassContract["safeTransferFrom(address,address,uint256)"](
          this.classOwner.address,
          this.randomSigner.address,
          0,
        ),
      ).to.be.rejectedWith("ErrVoidContract()");
    });

    it("safeTransferFrom (4 args) should revert", async function () {
      await expect(
        voidClassContract["safeTransferFrom(address,address,uint256,bytes)"](
          this.classOwner.address,
          this.randomSigner.address,
          0,
          "0x",
        ),
      ).to.be.rejectedWith("ErrVoidContract()");
    });

    it("approve should revert", async function () {
      await expect(
        voidClassContract.approve(this.randomSigner.address, 0),
      ).to.be.rejectedWith("ErrVoidContract()");
    });

    it("setApprovalForAll should revert", async function () {
      await expect(
        voidClassContract.setApprovalForAll(this.randomSigner.address, true),
      ).to.be.rejectedWith("ErrVoidContract()");
    });

    it("getApproved should revert", async function () {
      await expect(voidClassContract.getApproved(0)).to.be.rejectedWith(
        "ErrVoidContract()",
      );
    });

    it("isApprovedForAll should revert", async function () {
      await expect(
        voidClassContract.isApprovedForAll(
          this.classOwner.address,
          this.randomSigner.address,
        ),
      ).to.be.rejectedWith("ErrVoidContract()");
    });
  });

  describe("ERC721Enumerable Functions - All Should Revert", () => {
    it("totalSupply should revert", async function () {
      await expect(voidClassContract.totalSupply()).to.be.rejectedWith(
        "ErrVoidContract()",
      );
    });

    it("tokenByIndex should revert", async function () {
      await expect(voidClassContract.tokenByIndex(0)).to.be.rejectedWith(
        "ErrVoidContract()",
      );
    });

    it("tokenOfOwnerByIndex should revert", async function () {
      await expect(
        voidClassContract.tokenOfOwnerByIndex(this.classOwner.address, 0),
      ).to.be.rejectedWith("ErrVoidContract()");
    });
  });

  describe("ERC721Burnable Functions - Should Revert", () => {
    it("burn should revert", async function () {
      await expect(voidClassContract.burn(0)).to.be.rejectedWith(
        "ErrVoidContract()",
      );
    });
  });

  describe("BookNFT Custom Functions - All Should Revert", () => {
    it("contractURI should revert", async function () {
      await expect(voidClassContract.contractURI()).to.be.rejectedWith(
        "ErrVoidContract()",
      );
    });

    it("mint should revert", async function () {
      await expect(
        voidClassContract.mint(this.classOwner.address, ["memo"], ["metadata"]),
      ).to.be.rejectedWith("ErrVoidContract()");
    });

    it("batchMint should revert", async function () {
      await expect(
        voidClassContract.batchMint(
          [this.classOwner.address],
          ["memo"],
          ["metadata"],
        ),
      ).to.be.rejectedWith("ErrVoidContract()");
    });

    it("transferWithMemo should revert", async function () {
      await expect(
        voidClassContract.transferWithMemo(
          this.classOwner.address,
          this.randomSigner.address,
          0,
          "memo",
        ),
      ).to.be.rejectedWith("ErrVoidContract()");
    });

    it("batchTransferWithMemo should revert", async function () {
      await expect(
        voidClassContract.batchTransferWithMemo(
          this.classOwner.address,
          [this.randomSigner.address],
          [0n],
          ["memo"],
        ),
      ).to.be.rejectedWith("ErrVoidContract()");
    });
  });

  describe("Legacy Query Functions - Should Preserve Data", () => {
    it("legacyTotalSupply should return the minted count", async function () {
      expect(await voidClassContract.legacyTotalSupply()).to.equal(1n);
    });

    it("legacyName should return the original name", async function () {
      expect(await voidClassContract.legacyName()).to.equal(originalName);
    });

    it("legacySymbol should return the original symbol", async function () {
      expect(await voidClassContract.legacySymbol()).to.equal(originalSymbol);
    });

    it("legacyMaxSupply should return the original max supply", async function () {
      expect(await voidClassContract.legacyMaxSupply()).to.equal(
        originalMaxSupply,
      );
    });

    it("legacyTokenURI should return the raw token metadata", async function () {
      expect(await voidClassContract.legacyTokenURI(0)).to.equal(
        mintedTokenURI,
      );
    });

    it("legacyTokenURIEncoded should return base64 encoded data URI", async function () {
      const encoded = await voidClassContract.legacyTokenURIEncoded(0);
      expect(encoded).to.include("data:application/json;base64,");

      // Decode and verify
      const base64Part = encoded.replace("data:application/json;base64,", "");
      const decoded = Buffer.from(base64Part, "base64").toString("utf-8");
      expect(decoded).to.equal(mintedTokenURI);
    });

    it("legacyContractMetadata should return the raw contract metadata", async function () {
      const bookConfig = BookConfigLoader.load(
        "./test/fixtures/BookConfig0.json",
      );
      expect(await voidClassContract.legacyContractMetadata()).to.equal(
        bookConfig.metadata,
      );
    });

    it("legacyContractURIEncoded should return base64 encoded data URI", async function () {
      const encoded = await voidClassContract.legacyContractURIEncoded();
      expect(encoded).to.equal(originalContractURI);
    });

    it("legacyProtocolBeacon should return the protocol address", async function () {
      const protocolAddress = await protocolContract.getAddress();
      expect(await voidClassContract.legacyProtocolBeacon()).to.equal(
        protocolAddress,
      );
    });

    it("legacyRoyaltyFraction should return the royalty fraction", async function () {
      // Default royalty fraction is 0 unless set
      expect(await voidClassContract.legacyRoyaltyFraction()).to.equal(0n);
    });
  });

  describe("Storage Preservation After Upgrade", () => {
    it("should preserve all data after upgrade from BookNFT to BookNFTVoid", async function () {
      // Verify all legacy data is accessible
      expect(await voidClassContract.legacyName()).to.equal(originalName);
      expect(await voidClassContract.legacySymbol()).to.equal(originalSymbol);
      expect(await voidClassContract.legacyMaxSupply()).to.equal(
        originalMaxSupply,
      );
      expect(await voidClassContract.legacyTotalSupply()).to.equal(1n);
      expect(await voidClassContract.legacyTokenURI(0)).to.equal(
        mintedTokenURI,
      );
    });

    it("should be able to upgrade back to BookNFT and restore functionality", async function () {
      const likeProtocolOwnerSigner = protocolContract.connect(
        this.protocolOwner,
      );

      // Deploy a new BookNFT implementation
      const BookNFT = await ethers.getContractFactory("BookNFT");
      const newBookNFT = await BookNFT.deploy();
      await newBookNFT.initialize({
        creator: this.protocolOwner.address,
        updaters: [],
        minters: [],
        config: {
          name: "BookNFT Implementation",
          symbol: "BOOKNFTV0",
          metadata: "{}",
          max_supply: 1n,
        },
      });

      // Upgrade back to BookNFT
      await likeProtocolOwnerSigner.upgradeTo(newBookNFT.getAddress());

      // Get BookNFT contract at proxy address
      const restoredContract = await ethers.getContractAt(
        "BookNFT",
        nftClassId,
      );

      // Verify ERC721 functions work again
      expect(await restoredContract.supportsInterface("0x80ac58cd")).to.equal(
        true,
      );
      expect(await restoredContract.name()).to.equal(originalName);
      expect(await restoredContract.symbol()).to.equal(originalSymbol);
      expect(
        await restoredContract.balanceOf(this.classOwner.address),
      ).to.equal(1n);
      expect(await restoredContract.ownerOf(0)).to.equal(
        this.classOwner.address,
      );
    });
  });
});
