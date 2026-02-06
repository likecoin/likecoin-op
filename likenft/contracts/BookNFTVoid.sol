// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {Base64} from "@openzeppelin/contracts/utils/Base64.sol";

error ErrVoidContract();

/**
 * @title BookNFTVoid
 * @notice A void implementation for BookNFT that disables all ERC721 functionality.
 * @dev This contract is designed to replace BookNFT via beacon proxy upgrade to:
 *      - Hide NFTs from marketplaces like OpenSea by not reporting ERC721 interface
 *      - Freeze all transfers, approvals, and burns
 *      - Provide legacy query functions for apps that need to verify historical data
 *
 *      Storage layout must match BookNFT exactly for proxy compatibility.
 *      Ownership data is stored in OpenZeppelin's internal ERC721 storage slots
 *      and can be verified via blockchain Transfer event logs.
 */
contract BookNFTVoid is Initializable {
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    // Must match BookNFT storage layout exactly for proxy compatibility
    struct BookNFTStorage {
        string name;
        string symbol;
        string metadata;
        uint64 max_supply;
        uint256 _currentIndex;
        mapping(uint256 => string) tokenURIMap;
        uint96 royaltyFraction;
        address protocolBeacon;
    }

    // keccak256(abi.encode(uint256(keccak256("likeprotocol.booknft.storage")) - 1)) & ~bytes32(uint256(0xff))
    bytes32 private constant CLASS_DATA_STORAGE =
        0x8303e9d27d04c843c8d4a08966b1e1be0214fc0b3375d79db0a8252068c41f00;

    function _getClassStorage()
        private
        pure
        returns (BookNFTStorage storage $)
    {
        // solhint-disable-next-line no-inline-assembly
        assembly {
            $.slot := CLASS_DATA_STORAGE
        }
    }

    // ============ ERC165 Interface - Returns false for all ============

    /**
     * @notice Always returns false - this contract does not support any standard interfaces
     * @dev This causes marketplaces to not recognize this as an NFT contract
     */
    function supportsInterface(bytes4) external pure returns (bool) {
        return false;
    }

    // ============ ERC721 Interface - All Revert ============

    function balanceOf(address) external pure returns (uint256) {
        revert ErrVoidContract();
    }

    function ownerOf(uint256) external pure returns (address) {
        revert ErrVoidContract();
    }

    function tokenURI(uint256) external pure returns (string memory) {
        revert ErrVoidContract();
    }

    function name() external pure returns (string memory) {
        revert ErrVoidContract();
    }

    function symbol() external pure returns (string memory) {
        revert ErrVoidContract();
    }

    function transferFrom(address, address, uint256) external pure {
        revert ErrVoidContract();
    }

    function safeTransferFrom(address, address, uint256) external pure {
        revert ErrVoidContract();
    }

    function safeTransferFrom(
        address,
        address,
        uint256,
        bytes calldata
    ) external pure {
        revert ErrVoidContract();
    }

    function approve(address, uint256) external pure {
        revert ErrVoidContract();
    }

    function setApprovalForAll(address, bool) external pure {
        revert ErrVoidContract();
    }

    function getApproved(uint256) external pure returns (address) {
        revert ErrVoidContract();
    }

    function isApprovedForAll(address, address) external pure returns (bool) {
        revert ErrVoidContract();
    }

    // ============ ERC721Enumerable - All Revert ============

    function totalSupply() external pure returns (uint256) {
        revert ErrVoidContract();
    }

    function tokenByIndex(uint256) external pure returns (uint256) {
        revert ErrVoidContract();
    }

    function tokenOfOwnerByIndex(
        address,
        uint256
    ) external pure returns (uint256) {
        revert ErrVoidContract();
    }

    // ============ ERC721Burnable - Revert ============

    function burn(uint256) external pure {
        revert ErrVoidContract();
    }

    // ============ BookNFT Custom Functions - All Revert ============

    function contractURI() external pure returns (string memory) {
        revert ErrVoidContract();
    }

    function mint(address, string[] calldata, string[] calldata) external pure {
        revert ErrVoidContract();
    }

    function batchMint(
        address[] calldata,
        string[] calldata,
        string[] calldata
    ) external pure {
        revert ErrVoidContract();
    }

    function transferWithMemo(
        address,
        address,
        uint256,
        string calldata
    ) external pure {
        revert ErrVoidContract();
    }

    function batchTransferWithMemo(
        address,
        address[] calldata,
        uint256[] calldata,
        string[] calldata
    ) external pure {
        revert ErrVoidContract();
    }

    // ============ Legacy Query Functions ============
    // These allow apps to verify historical data without ERC721 interface

    /**
     * @notice Returns the total number of tokens that were minted
     * @return The current token index (total minted count)
     */
    function legacyTotalSupply() external view returns (uint256) {
        return _getClassStorage()._currentIndex;
    }

    /**
     * @notice Returns the collection name
     * @return The name of the collection
     */
    function legacyName() external view returns (string memory) {
        return _getClassStorage().name;
    }

    /**
     * @notice Returns the collection symbol
     * @return The symbol of the collection
     */
    function legacySymbol() external view returns (string memory) {
        return _getClassStorage().symbol;
    }

    /**
     * @notice Returns the raw metadata for a specific token
     * @param tokenId The token ID to query
     * @return The raw metadata string (not base64 encoded)
     */
    function legacyTokenURI(
        uint256 tokenId
    ) external view returns (string memory) {
        BookNFTStorage storage $ = _getClassStorage();
        return $.tokenURIMap[tokenId];
    }

    /**
     * @notice Returns the token metadata as a data URI (matching original format)
     * @param tokenId The token ID to query
     * @return The metadata as a base64-encoded data URI
     */
    function legacyTokenURIEncoded(
        uint256 tokenId
    ) external view returns (string memory) {
        BookNFTStorage storage $ = _getClassStorage();
        return
            string(
                abi.encodePacked(
                    "data:application/json;base64,",
                    Base64.encode(abi.encodePacked($.tokenURIMap[tokenId]))
                )
            );
    }

    /**
     * @notice Returns the raw contract-level metadata
     * @return The raw metadata string (not base64 encoded)
     */
    function legacyContractMetadata() external view returns (string memory) {
        return _getClassStorage().metadata;
    }

    /**
     * @notice Returns the contract metadata as a data URI (matching original format)
     * @return The metadata as a base64-encoded data URI
     */
    function legacyContractURIEncoded() external view returns (string memory) {
        BookNFTStorage storage $ = _getClassStorage();
        return
            string(
                abi.encodePacked(
                    "data:application/json;base64,",
                    Base64.encode(abi.encodePacked($.metadata))
                )
            );
    }

    /**
     * @notice Returns the maximum supply configured for this collection
     * @return The maximum supply
     */
    function legacyMaxSupply() external view returns (uint64) {
        return _getClassStorage().max_supply;
    }

    /**
     * @notice Returns the protocol beacon address
     * @return The address of the LikeProtocol beacon
     */
    function legacyProtocolBeacon() external view returns (address) {
        return _getClassStorage().protocolBeacon;
    }

    /**
     * @notice Returns the royalty fraction (basis points, denominator is 10000)
     * @return The royalty fraction
     */
    function legacyRoyaltyFraction() external view returns (uint96) {
        return _getClassStorage().royaltyFraction;
    }

    // Note: Cannot provide legacyOwnerOf because ERC721 ownership is stored
    // in OpenZeppelin's internal mappings (ERC721Storage), not in BookNFTStorage.
    // Historical ownership can be verified via blockchain Transfer event logs.
}
