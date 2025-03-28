import { expect } from "chai";
import { BaseContract, EventLog } from "ethers";
import { ethers, upgrades } from "hardhat";

describe("LikeProtocol", () => {
  before(async function () {
    this.LikeProtocol = await ethers.getContractFactory("LikeProtocol");
    this.LikeProtocolMock = await ethers.getContractFactory("LikeProtocolMock");
    const [ownerSigner, randomSigner] = await ethers.getSigners();

    this.ownerSigner = ownerSigner;
    this.randomSigner = randomSigner;
  });

  let deployment: BaseContract;
  let contractAddress: string;
  let contract: any;
  beforeEach(async function () {
    const likeProtocol = await upgrades.deployProxy(
      this.LikeProtocol,
      [this.ownerSigner.address],
      {
        initializer: "initialize",
      },
    );
    deployment = await likeProtocol.waitForDeployment();
    contractAddress = await deployment.getAddress();
    contract = await ethers.getContractAt("LikeProtocol", contractAddress);
  });

  it("should be upgradable", async function () {
    const likeProtocolMockOwnerSigner = this.LikeProtocolMock.connect(
      this.ownerSigner,
    );
    const newLikeProtocol = await upgrades.upgradeProxy(
      contractAddress,
      likeProtocolMockOwnerSigner,
    );
    expect(await newLikeProtocol.getAddress()).to.equal(contractAddress);
    expect(await newLikeProtocol.owner()).to.equal(this.ownerSigner.address);
    expect(await newLikeProtocol.version()).to.equal(2n);
  });

  it("should not be upgradable by random address", async function () {
    const likeProtocolMockRandomSigner = this.LikeProtocolMock.connect(
      this.randomSigner,
    );
    await expect(
      upgrades.upgradeProxy(contractAddress, likeProtocolMockRandomSigner),
    ).to.be.rejected;
    expect(await contract.owner()).to.equal(this.ownerSigner.address);
  });

  it("should set the right owner", async function () {
    expect(await contract.owner()).to.equal(this.ownerSigner.address);
  });

  it("should allow ownership transfer", async function () {
    await contract.transferOwnership(this.randomSigner.address);
    expect(await contract.owner()).to.equal(this.randomSigner.address);
  });

  it("should not allow random ownership transfer", async function () {
    const likeProtocolSigner = contract.connect(this.randomSigner);
    try {
      await likeProtocolSigner.transferOwnership(this.randomSigner.address);
    } catch (error) {
      expect(error).to.be.instanceOf(Error);
    }
    expect(await contract.owner()).to.equal(this.ownerSigner.address);
  });

  it("should not paused by random address", async function () {
    const likeProtocolSigner = contract.connect(this.randomSigner);
    await expect(likeProtocolSigner.pause()).to.be.rejected;
  });

  it("should be paused by owner address", async function () {
    const likeProtocolSigner = contract.connect(this.ownerSigner);
    await expect(likeProtocolSigner.pause()).to.be.not.rejected;
  });

  it("should not unpaused by random address", async function () {
    const likeProtocolOwnerSigner = contract.connect(this.ownerSigner);
    await expect(likeProtocolOwnerSigner.pause()).to.be.not.rejected;
    const likeProtocolRandomSigner = contract.connect(this.randomSigner);
    await expect(likeProtocolRandomSigner.unpause()).to.be.rejected;
  });

  it("should be able to create new BookNFT", async function () {
    const likeProtocolOwnerSigner = contract.connect(this.ownerSigner);

    const newClass = async () => {
      await likeProtocolOwnerSigner
        .newBookNFT({
          creator: this.ownerSigner,
          updaters: [this.ownerSigner],
          minters: [this.ownerSigner],
          config: {
            name: "My Book",
            symbol: "KOOB",
            metadata: JSON.stringify({
              name: "Collection Name",
              symbol: "Collection SYMB",
              description: "Collection Description",
              image:
                "ipfs://bafybeiezq4yqosc2u4saanove5bsa3yciufwhfduemy5z6vvf6q3c5lnbi",
              banner_image: "",
              featured_image: "",
              external_link: "https://www.example.com",
              collaborators: [],
            }),
            max_supply: 10,
          },
        })
        .then((tx) => tx.wait());
    };

    const NewClassEvent = new Promise<{ id: string }>((resolve, reject) => {
      likeProtocolOwnerSigner.on("NewBookNFT", (id, params, event) => {
        event.removeListener();
        resolve({ id });
      });

      setTimeout(() => {
        reject(new Error("timeout"));
      }, 20000);
    });

    await expect(newClass()).to.be.not.rejected;
    const newClassEvent = await NewClassEvent;
    const classId = newClassEvent.id;

    const _newNFTClass = await ethers.getContractAt("BookNFT", classId);
    expect(await _newNFTClass.name()).to.equal("My Book");
    expect(await _newNFTClass.symbol()).to.equal("KOOB");
  });

  it("should not allow to create new BookNFT when paused", async function () {
    const likeProtocolSigner = contract.connect(this.ownerSigner);

    const classOperation = async () => {
      await likeProtocolSigner
        .newBookNFT({
          creator: this.ownerSigner,
          updaters: [this.ownerSigner],
          minters: [this.ownerSigner],
          config: {
            name: "My Book",
            symbol: "KOOB",
            metadata: JSON.stringify({
              name: "Collection Name",
              symbol: "Collection SYMB",
              description: "Collection Description",
              image:
                "ipfs://bafybeiezq4yqosc2u4saanove5bsa3yciufwhfduemy5z6vvf6q3c5lnbi",
              banner_image: "",
              featured_image: "",
              external_link: "https://www.example.com",
              collaborators: [],
            }),
            max_supply: 10,
          },
        })
        .then((tx) => tx.wait());
    };

    await expect(classOperation()).to.be.not.rejected;
    await expect(likeProtocolSigner.pause()).to.be.not.rejected;
    await expect(classOperation()).to.be.rejectedWith(
      "VM Exception while processing transaction: reverted with custom error 'EnforcedPause()'",
    );
    await expect(likeProtocolSigner.unpause()).to.be.not.rejected;
    await expect(classOperation()).to.be.not.rejected;
  });

  it("should be allow everyone to create new BookNFT", async function () {
    const likeNFTSigner = contract.connect(this.randomSigner);

    const newClass = async () => {
      await likeNFTSigner
        .newBookNFT({
          creator: this.randomSigner,
          updaters: [this.randomSigner],
          minters: [this.randomSigner],
          config: {
            name: "My Book",
            symbol: "KOOB",
            metadata: JSON.stringify({
              name: "Random by somone",
              symbol: "No data",
            }),
            max_supply: 10,
          },
        })
        .then((tx) => tx.wait());
    };

    await expect(newClass()).to.be.not.rejected;
  });

  it("should not allow everyone to create new BookNFT when paused", async function () {
    const likeProtocolOwnerSigner = contract.connect(this.ownerSigner);
    const likeProtocolRandomSigner = contract.connect(this.randomSigner);

    const classOperation = async () => {
      await likeProtocolRandomSigner
        .newBookNFT({
          creator: this.randomSigner,
          updaters: [this.randomSigner],
          minters: [this.randomSigner],
          config: {
            name: "My Book",
            symbol: "KOOB",
            metadata: JSON.stringify({
              name: "Random by somone",
              symbol: "No data",
            }),
            max_supply: 10,
          },
        })
        .then((tx) => tx.wait());
    };

    await expect(classOperation()).to.be.not.rejected;
    await expect(likeProtocolOwnerSigner.pause()).to.be.not.rejected;
    await expect(classOperation()).to.be.rejectedWith(
      "VM Exception while processing transaction: reverted with custom error 'EnforcedPause()'",
    );
    await expect(likeProtocolOwnerSigner.unpause()).to.be.not.rejected;
    await expect(classOperation()).to.be.not.rejected;
  });

  it("should only allow BookNFT owner to update the BookNFT", async function () {
    const likeProtocolOwnerSigner = contract.connect(this.ownerSigner);
    const likeProtocolRandomSigner = contract.connect(this.randomSigner);

    const newClass = async () => {
      await likeProtocolOwnerSigner
        .newBookNFT({
          creator: this.ownerSigner,
          updaters: [this.ownerSigner],
          minters: [this.ownerSigner],
          config: {
            name: "My Book",
            symbol: "KOOB",
            metadata: JSON.stringify({
              name: "Collection Name",
              symbol: "Collection SYMB",
              description: "Collection Description",
              image:
                "ipfs://bafybeiezq4yqosc2u4saanove5bsa3yciufwhfduemy5z6vvf6q3c5lnbi",
              banner_image: "",
              featured_image: "",
              external_link: "https://www.example.com",
              collaborators: [],
            }),
            max_supply: 10,
          },
        })
        .then((tx) => tx.wait());
    };

    const NewClassEvent = new Promise<{ id: string }>((resolve, reject) => {
      likeProtocolOwnerSigner.on("NewBookNFT", (id, params, event) => {
        event.removeListener();
        resolve({ id });
      });

      setTimeout(() => {
        reject(new Error("timeout"));
      }, 20000);
    });

    await expect(newClass()).to.be.not.rejected;
    const newClassEvent = await NewClassEvent;
    const classId = newClassEvent.id;

    const _newNFTClass = await ethers.getContractAt("BookNFT", classId);
    await expect(await _newNFTClass.owner()).to.equal(this.ownerSigner.address);
    await expect(await _newNFTClass.symbol()).to.equal("KOOB");

    await expect(
      likeProtocolRandomSigner.updateBookNFT({
        classId: classId,
        config: {
          name: "Hi Jack",
          symbol: "HIJACK",
          metadata: JSON.stringify({}),
          max_supply: 0,
        },
      }),
    ).to.be.rejected;
    await expect(await _newNFTClass.owner()).to.equal(this.ownerSigner.address);
    await expect(await _newNFTClass.symbol()).to.equal("KOOB");
  });

  it("should retain the BookNFT paused state after upgrade", async function () {
    const likeProtocolOwnerSigner = contract.connect(this.ownerSigner);

    expect(await likeProtocolOwnerSigner.paused()).to.be.false;
    await likeProtocolOwnerSigner.pause();
    expect(await likeProtocolOwnerSigner.paused()).to.be.true;

    const likeProtocolMockOwnerSigner = this.LikeProtocolMock.connect(
      this.ownerSigner,
    );
    const newLikeProtocol = await upgrades.upgradeProxy(
      contractAddress,
      likeProtocolMockOwnerSigner,
    );
    expect(await newLikeProtocol.owner()).to.equal(this.ownerSigner.address);
    expect(await newLikeProtocol.version()).to.equal(2n);

    const proxyContract = await ethers.getContractAt(
      "LikeProtocolMock",
      contractAddress,
    );
    expect(await proxyContract.version()).to.equal(2n);

    expect(await likeProtocolOwnerSigner.paused()).to.be.true;
    await likeProtocolOwnerSigner.unpause();
    expect(await likeProtocolOwnerSigner.paused()).to.be.false;
  });

  it("should retain the BookNFT mapping after upgrade", async function () {
    const likeProtocolOwnerSigner = contract.connect(this.ownerSigner);
    const newClass = async () => {
      await likeProtocolOwnerSigner
        .newBookNFT({
          creator: this.ownerSigner,
          updaters: [this.ownerSigner],
          minters: [this.ownerSigner],
          config: {
            name: "My Book",
            symbol: "KOOB",
            metadata: JSON.stringify({
              name: "Collection Name",
              symbol: "Collection SYMB",
              description: "Collection Description",
              image:
                "ipfs://bafybeiezq4yqosc2u4saanove5bsa3yciufwhfduemy5z6vvf6q3c5lnbi",
              banner_image: "",
              featured_image: "",
              external_link: "https://www.example.com",
              collaborators: [],
            }),
            max_supply: 10,
          },
        })
        .then((tx) => tx.wait());
    };

    const NewClassEvent = new Promise<{ id: string }>((resolve, reject) => {
      likeProtocolOwnerSigner.on("NewBookNFT", (id, params, event) => {
        event.removeListener();
        resolve({ id });
      });

      setTimeout(() => {
        reject(new Error("timeout"));
      }, 20000);
    });

    await expect(newClass()).to.be.not.rejected;
    const newClassEvent = await NewClassEvent;
    const classId = newClassEvent.id;
    expect(await likeProtocolOwnerSigner.isBookNFT(classId)).to.be.true;

    const likeProtocolMockOwnerSigner = this.LikeProtocolMock.connect(
      this.ownerSigner,
    );
    const newLikeProtocol = await upgrades.upgradeProxy(
      contractAddress,
      likeProtocolMockOwnerSigner,
    );
    expect(await newLikeProtocol.owner()).to.equal(this.ownerSigner.address);
    expect(await newLikeProtocol.version()).to.equal(2n);
    expect(await newLikeProtocol.isBookNFT(classId)).to.be.true;

    const proxyContract = await ethers.getContractAt(
      "LikeProtocolMock",
      contractAddress,
    );
    expect(await proxyContract.isBookNFT(classId)).to.be.true;
  });
});
