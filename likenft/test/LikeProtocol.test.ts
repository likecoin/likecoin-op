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

  it("should be able to create new class", async function () {
    const likeProtocolOwnerSigner = contract.connect(this.ownerSigner);

    const newClass = async () => {
      await likeProtocolOwnerSigner
        .newClass({
          creator: this.ownerSigner,
          updaters: [this.ownerSigner],
          minters: [this.ownerSigner],
          input: {
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
            config: {
              max_supply: 10,
            },
          },
        })
        .then((tx) => tx.wait());
    };

    const NewClassEvent = new Promise<{ id: string }>((resolve, reject) => {
      likeProtocolOwnerSigner.on("NewClass", (id, params, event) => {
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

    const _newNFTClass = await ethers.getContractAt("LikeNFTClass", classId);
    expect(await _newNFTClass.name()).to.equal("My Book");
    expect(await _newNFTClass.symbol()).to.equal("KOOB");
  });

  it("should not allow to create new class when paused", async function () {
    const likeProtocolSigner = contract.connect(this.ownerSigner);

    const classOperation = async () => {
      await likeProtocolSigner
        .newClass({
          creator: this.ownerSigner,
          updaters: [this.ownerSigner],
          minters: [this.ownerSigner],
          input: {
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
            config: {
              max_supply: 10,
            },
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

  it("should be allow everyone to create new class", async function () {
    const likeNFTSigner = contract.connect(this.randomSigner);

    const newClass = async () => {
      await likeNFTSigner
        .newClass({
          creator: this.randomSigner,
          updaters: [this.randomSigner],
          minters: [this.randomSigner],
          input: {
            name: "My Book",
            symbol: "KOOB",
            metadata: JSON.stringify({
              name: "Random by somone",
              symbol: "No data",
            }),
            config: {
              max_supply: 10,
            },
          },
        })
        .then((tx) => tx.wait());
    };

    await expect(newClass()).to.be.not.rejected;
  });

  it("should not allow everyone to create new class when paused", async function () {
    const likeProtocolOwnerSigner = contract.connect(this.ownerSigner);
    const likeProtocolRandomSigner = contract.connect(this.randomSigner);

    const classOperation = async () => {
      await likeProtocolRandomSigner
        .newClass({
          creator: this.randomSigner,
          updaters: [this.randomSigner],
          minters: [this.randomSigner],
          input: {
            name: "My Book",
            symbol: "KOOB",
            metadata: JSON.stringify({
              name: "Random by somone",
              symbol: "No data",
            }),
            config: {
              max_supply: 10,
            },
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

  it("should only allow LikeNFTClass owner to update the class", async function () {
    const likeProtocolOwnerSigner = contract.connect(this.ownerSigner);
    const likeProtocolRandomSigner = contract.connect(this.randomSigner);

    const newClass = async () => {
      await likeProtocolOwnerSigner
        .newClass({
          creator: this.ownerSigner,
          updaters: [this.ownerSigner],
          minters: [this.ownerSigner],
          input: {
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
            config: {
              max_supply: 10,
            },
          },
        })
        .then((tx) => tx.wait());
    };

    const NewClassEvent = new Promise<{ id: string }>((resolve, reject) => {
      likeProtocolOwnerSigner.on("NewClass", (id, params, event) => {
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

    const _newNFTClass = await ethers.getContractAt("LikeNFTClass", classId);
    await expect(await _newNFTClass.owner()).to.equal(this.ownerSigner.address);
    await expect(await _newNFTClass.symbol()).to.equal("KOOB");

    await expect(
      likeProtocolRandomSigner.updateClass({
        creator: this.randomSigner,
        updaters: [this.randomSigner],
        minters: [this.randomSigner],
        class_id: classId,
        input: {
          name: "Hi Jack",
          symbol: "HIJACK",
          metadata: JSON.stringify({}),
          config: {
            max_supply: 0,
          },
        },
      }),
    ).to.be.rejected;
    await expect(await _newNFTClass.owner()).to.equal(this.ownerSigner.address);
    await expect(await _newNFTClass.symbol()).to.equal("KOOB");
  });
});
