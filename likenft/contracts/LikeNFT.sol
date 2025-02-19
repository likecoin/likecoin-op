// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";

import {MsgMintNFT} from "../types/msgs/MsgMintNFT.sol";
import {MsgMintNFTs} from "../types/msgs/MsgMintNFTs.sol";
import {MsgNewClass} from "../types/msgs/MsgNewClass.sol";
import {MsgUpdateClass} from "../types/msgs/MsgUpdateClass.sol";

import {Class} from "./Class.sol";

error ErrNftClassNotFound();

contract LikeNFT is
    Initializable,
    UUPSUpgradeable,
    OwnableUpgradeable,
    PausableUpgradeable
{
    struct LikeNFTStorage {
        address minter;
        mapping(address creator => mapping(address class_id => Class)) creatorClassIdClassMapping;
        Class[] classes;
    }
    // keccak256(abi.encode(uint256(keccak256("likenft.storage")) - 1)) & ~bytes32(uint256(0xff))
    bytes32 private constant DATA_STORAGE =
        0xf59cae2d8704429a88f4a038c4cff8d2643dc6b4647d519013fb42e0b4344200;
    function _getLikeNFTStorage()
        private
        pure
        returns (LikeNFTStorage storage $)
    {
        // solhint-disable-next-line no-inline-assembly
        assembly {
            $.slot := DATA_STORAGE
        }
    }

    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");

    event NewClass(address classId, MsgNewClass params);

    function initialize(address initialOwner) public initializer {
        __UUPSUpgradeable_init();
        __Ownable_init(initialOwner);
        __Pausable_init();
    }

    function pause() public {
        _pause();
    }

    function unpause() public {
        _unpause();
    }

    function newClass(MsgNewClass memory msgNewClass) public whenNotPaused {
        LikeNFTStorage storage $ = _getLikeNFTStorage();
        Class class = new Class(msgNewClass);
        address id = address(class);
        $.classes.push(class);
        $.creatorClassIdClassMapping[msgNewClass.creator][id] = class;
        emit NewClass(id, msgNewClass);
    }

    function updateClass(
        MsgUpdateClass memory msgUpdateClass
    ) public whenNotPaused {
        LikeNFTStorage storage $ = _getLikeNFTStorage();
        Class class = $.creatorClassIdClassMapping[msgUpdateClass.creator][
            msgUpdateClass.class_id
        ];
        if (address(class) == address(0)) {
            revert ErrNftClassNotFound();
        }
        class.update(msgUpdateClass.input);
    }

    function mintNFT(MsgMintNFT calldata msgMintNFT) external whenNotPaused {
        LikeNFTStorage storage $ = _getLikeNFTStorage();
        Class class = $.creatorClassIdClassMapping[msgMintNFT.creator][
            msgMintNFT.class_id
        ];
        if (address(class) == address(0)) {
            revert ErrNftClassNotFound();
        }
        string[] memory metadata_list = new string[](1);
        metadata_list[0] = msgMintNFT.input.metadata;
        class.mint(msgMintNFT.creator, metadata_list);
    }

    function mintNFTs(MsgMintNFTs calldata msgMintNFTs) external whenNotPaused {
        LikeNFTStorage storage $ = _getLikeNFTStorage();
        Class class = $.creatorClassIdClassMapping[msgMintNFTs.creator][
            msgMintNFTs.class_id
        ];
        if (address(class) == address(0)) {
            revert ErrNftClassNotFound();
        }
        string[] memory metadata_list = new string[](msgMintNFTs.inputs.length);
        for (uint i = 0; i < msgMintNFTs.inputs.length; i++) {
            metadata_list[i] = msgMintNFTs.inputs[i].metadata;
        }
        class.mint(msgMintNFTs.creator, metadata_list);
    }

    function _authorizeUpgrade(
        address _newImplementation // solhint-disable-next-line no-empty-blocks
    ) internal override {}
}
