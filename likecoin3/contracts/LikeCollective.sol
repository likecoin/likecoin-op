// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";

import {Likecoin} from "./Likecoin.sol";
import {LikeStakePosition} from "./LikeStakePosition.sol";

/// @custom:security-contact rickmak@oursky.com
contract LikeCollective is
    Initializable,
    OwnableUpgradeable,
    UUPSUpgradeable,
    PausableUpgradeable,
    ReentrancyGuardUpgradeable
{
    struct PoolData {
        uint256 totalStaked;
        uint256 totalRewarded;
        uint256 rewardPending;
        uint256 rewardIndex;
        mapping(uint256 tokenId => uint256 rewardIndex) rewardIndexes;
    }

    struct CollectiveData {
        Likecoin likecoin;
        LikeStakePosition likeStakePosition;
        mapping(address bookNFT => PoolData) pools;
    }

    uint256 public constant ACC_REWARD_PRECISION = 1e18; // Precision scalar for reward index

    // keccak256(abi.encode(uint256(keccak256("likecollective.storage")) - 1)) & ~bytes32(uint256(0xff))
    bytes32 private constant CLASS_DATA_STORAGE =
        0xe9c9d9e1df02920d747aa7516ca1d4362d70267096e6330bcfb24b265ac2ee00;

    function _getCollectiveData()
        private
        pure
        returns (CollectiveData storage $)
    {
        // solhint-disable-next-line no-inline-assembly
        assembly {
            $.slot := CLASS_DATA_STORAGE
        }
    }

    // solhint-disable gas-indexed-events
    event Staked(
        address indexed bookNFT,
        address indexed account,
        uint256 stakedAmount
    );
    event Unstaked(
        address indexed bookNFT,
        address indexed account,
        uint256 stakedAmount
    );
    event RewardClaimed(
        address indexed bookNFT,
        address indexed account,
        uint256 rewardedAmount
    );
    event RewardDeposited(
        address indexed bookNFT,
        address indexed account,
        uint256 rewardedAmount
    );
    event AllRewardClaimed(
        address indexed account,
        RewardData[] rewardedAmount
    );
    // solhint-enable gas-indexed-events

    struct RewardData {
        address bookNFT;
        uint256 rewardedAmount;
    }

    // Errors
    error ErrInvalidBookNFT(address bookNFT);
    error ErrInvalidOwner();
    error ErrInsufficientStake(uint256 required, uint256 available);
    error ErrNoRewardsToClaim();
    error ErrInvalidAmount();

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize(address initialOwner) public initializer {
        __UUPSUpgradeable_init();
        __Ownable_init(initialOwner);
        __Pausable_init();
        __ReentrancyGuard_init();
    }

    // Start Admin functions
    function pause() public onlyOwner {
        _pause();
    }

    function unpause() public onlyOwner {
        _unpause();
    }

    function setLikecoin(address likecoin) external onlyOwner {
        CollectiveData storage $ = _getCollectiveData();
        $.likecoin = Likecoin(likecoin);
    }

    function setLikeStakePosition(
        address likeStakePosition
    ) external onlyOwner {
        CollectiveData storage $ = _getCollectiveData();
        $.likeStakePosition = LikeStakePosition(likeStakePosition);
    }

    function _authorizeUpgrade(
        address _newImplementation
    ) internal override onlyOwner {
        // TODO: Add any additional authorization logic if needed
    }
    // End Admin functions

    // Private View functions
    function _getLikeStakePosition() internal view returns (LikeStakePosition) {
        return _getCollectiveData().likeStakePosition;
    }

    function _getPosition(
        address user
    ) internal view returns (uint256[] memory) {
        LikeStakePosition likeStakePosition = _getLikeStakePosition();
        return likeStakePosition.getUserPositions(user);
    }

    function _pendingRewardsOf(
        uint256 tokenId
    ) internal view returns (uint256) {
        CollectiveData storage $ = _getCollectiveData();
        LikeStakePosition.Position memory p = $.likeStakePosition.getPosition(
            tokenId
        );
        uint256 poolIndex = $.pools[p.bookNFT].rewardIndex;
        uint256 positionRewardIndex = $.pools[p.bookNFT].rewardIndexes[tokenId];
        return
            (p.stakedAmount * (poolIndex - positionRewardIndex)) /
            ACC_REWARD_PRECISION;
    }
    // End Private View functions

    function newStakePosition(
        address bookNFT,
        uint256 amount
    ) external whenNotPaused nonReentrant {
        CollectiveData storage $ = _getCollectiveData();
        PoolData storage pool = $.pools[bookNFT];
        $.likecoin.transferFrom(_msgSender(), address(this), amount);
        uint256 tokenId = $.likeStakePosition.mintPosition(
            _msgSender(),
            bookNFT,
            amount,
            pool.rewardIndex
        );
        pool.totalStaked = pool.totalStaked + amount;
        pool.rewardIndexes[tokenId] = pool.rewardIndex;
        emit Staked(bookNFT, _msgSender(), amount);
    }

    function increaseStakeToPosition(
        uint256 tokenID,
        uint256 amount
    ) external whenNotPaused nonReentrant {
        CollectiveData storage $ = _getCollectiveData();
        address owner = $.likeStakePosition.ownerOf(tokenID);
        if (owner != _msgSender()) revert ErrInvalidOwner();
        $.likecoin.transferFrom(_msgSender(), address(this), amount);

        LikeStakePosition.Position memory position = $
            .likeStakePosition
            .getPosition(tokenID);
        uint256 pendingRewards = _pendingRewardsOf(tokenID);
        uint256 newAmount = position.stakedAmount + amount + pendingRewards;
        address bookNFT = position.bookNFT;

        PoolData storage pool = $.pools[bookNFT];
        // Pending rewards are compounded into the position principal (newAmount),
        // so they must be added to the pool total as well, mirroring restakeRewardPosition.
        pool.totalStaked = pool.totalStaked + amount + pendingRewards;
        pool.rewardPending = pool.rewardPending - pendingRewards;
        pool.totalRewarded = pool.totalRewarded + pendingRewards;
        pool.rewardIndexes[tokenID] = pool.rewardIndex;
        $.likeStakePosition.updatePosition(
            tokenID,
            newAmount,
            pool.rewardIndex
        );

        emit Staked(bookNFT, _msgSender(), amount);
        emit RewardClaimed(bookNFT, _msgSender(), pendingRewards);
        if (pendingRewards > 0) {
            // Compounded rewards add to the staked principal too, so event consumers
            // (e.g. the indexer) that sum Staked amounts stay consistent with totalStaked.
            emit Staked(bookNFT, _msgSender(), pendingRewards);
        }
    }

    function decreaseStakePosition(
        uint256 tokenID,
        uint256 amount
    ) external whenNotPaused nonReentrant {
        CollectiveData storage $ = _getCollectiveData();
        address owner = $.likeStakePosition.ownerOf(tokenID);
        if (owner != _msgSender()) revert ErrInvalidOwner();

        LikeStakePosition.Position memory position = $
            .likeStakePosition
            .getPosition(tokenID);
        if (amount > position.stakedAmount) revert ErrInvalidAmount();

        uint256 pendingRewards = _pendingRewardsOf(tokenID);
        uint256 amountToTransfer = amount + pendingRewards;
        uint256 newAmount = position.stakedAmount - amount;
        address bookNFT = position.bookNFT;
        PoolData storage pool = $.pools[bookNFT];
        // Claim rewards
        pool.totalStaked = pool.totalStaked - amount;
        pool.rewardPending = pool.rewardPending - pendingRewards;
        pool.totalRewarded = pool.totalRewarded + pendingRewards;
        pool.rewardIndexes[tokenID] = pool.rewardIndex;
        $.likeStakePosition.updatePosition(
            tokenID,
            newAmount,
            pool.rewardIndex
        );

        $.likecoin.transfer(_msgSender(), amountToTransfer);
        emit Unstaked(bookNFT, _msgSender(), amount);
        emit RewardClaimed(bookNFT, _msgSender(), pendingRewards);
    }

    function removeStakePosition(
        uint256 tokenId
    ) external whenNotPaused nonReentrant {
        CollectiveData storage $ = _getCollectiveData();
        address owner = $.likeStakePosition.ownerOf(tokenId);
        if (owner != _msgSender()) revert ErrInvalidOwner();

        LikeStakePosition.Position memory position = $
            .likeStakePosition
            .getPosition(tokenId);
        uint256 amount = position.stakedAmount;
        uint256 pendingRewards = _pendingRewardsOf(tokenId);
        uint256 totalAmount = amount + pendingRewards;
        address bookNFT = position.bookNFT;
        PoolData storage pool = $.pools[bookNFT];
        // Claim rewards
        pool.totalStaked = pool.totalStaked - amount;
        pool.totalRewarded = pool.totalRewarded + pendingRewards;
        pool.rewardPending = pool.rewardPending - pendingRewards;
        delete pool.rewardIndexes[tokenId];

        $.likeStakePosition.burnPosition(tokenId);
        $.likecoin.transfer(_msgSender(), totalAmount);
        emit RewardClaimed(bookNFT, _msgSender(), pendingRewards);
        emit Unstaked(bookNFT, _msgSender(), amount);
    }

    function claimRewards(uint256 tokenId) external whenNotPaused nonReentrant {
        CollectiveData storage $ = _getCollectiveData();
        address owner = $.likeStakePosition.ownerOf(tokenId);
        if (owner != _msgSender()) revert ErrInvalidOwner();

        LikeStakePosition.Position memory position = $
            .likeStakePosition
            .getPosition(tokenId);
        address bookNFT = position.bookNFT;
        PoolData storage pool = $.pools[bookNFT];
        uint256 pendingRewards = _pendingRewardsOf(tokenId);
        pool.rewardPending -= pendingRewards;
        pool.totalRewarded += pendingRewards;
        pool.rewardIndexes[tokenId] = pool.rewardIndex;
        $.likeStakePosition.updatePositionRewardIndex(
            tokenId,
            pool.rewardIndex
        );
        $.likecoin.transfer(_msgSender(), pendingRewards);
        emit RewardClaimed(bookNFT, _msgSender(), pendingRewards);
    }

    function claimAllRewards(address user) external whenNotPaused nonReentrant {
        if (user != _msgSender()) revert ErrInvalidOwner();
        CollectiveData storage $ = _getCollectiveData();
        uint256[] memory positions = $.likeStakePosition.getUserPositions(user);
        uint256 totalRewards = 0;
        RewardData[] memory rewards = new RewardData[](positions.length);
        for (uint256 i = 0; i < positions.length; ++i) {
            uint256 p = positions[i];
            uint256 pendingRewards = _pendingRewardsOf(p);
            address bookNFT = $.likeStakePosition.getPosition(p).bookNFT;

            PoolData storage pool = $.pools[bookNFT];
            pool.rewardPending -= pendingRewards;
            pool.totalRewarded += pendingRewards;
            pool.rewardIndexes[p] = pool.rewardIndex;
            $.likeStakePosition.updatePositionRewardIndex(p, pool.rewardIndex);
            rewards[i] = RewardData({
                bookNFT: bookNFT,
                rewardedAmount: pendingRewards
            });
            totalRewards += pendingRewards;
        }
        $.likecoin.transfer(_msgSender(), totalRewards);
        emit AllRewardClaimed(_msgSender(), rewards);
    }

    function restakeRewardPosition(
        uint256 tokenId
    ) external whenNotPaused nonReentrant {
        CollectiveData storage $ = _getCollectiveData();
        address owner = $.likeStakePosition.ownerOf(tokenId);
        if (owner != _msgSender()) revert ErrInvalidOwner();

        LikeStakePosition.Position memory position = $
            .likeStakePosition
            .getPosition(tokenId);
        uint256 pendingRewards = _pendingRewardsOf(tokenId);
        address bookNFT = position.bookNFT;
        PoolData storage pool = $.pools[bookNFT];

        pool.totalStaked += pendingRewards;
        pool.rewardPending -= pendingRewards;
        pool.totalRewarded += pendingRewards;
        pool.rewardIndexes[tokenId] = pool.rewardIndex;
        $.likeStakePosition.updatePosition(
            tokenId,
            position.stakedAmount + pendingRewards,
            pool.rewardIndex
        );

        emit RewardClaimed(bookNFT, _msgSender(), pendingRewards);
        emit Staked(bookNFT, _msgSender(), pendingRewards);
    }

    function depositReward(
        address bookNFT,
        uint256 amount
    ) external whenNotPaused {
        CollectiveData storage $ = _getCollectiveData();
        $.likecoin.transferFrom(_msgSender(), address(this), amount);
        PoolData storage pool = $.pools[bookNFT];
        pool.rewardPending += amount;
        pool.rewardIndex += (amount * ACC_REWARD_PRECISION) / pool.totalStaked;
        emit RewardDeposited(bookNFT, _msgSender(), amount);
    }

    // View function for stake
    function getTotalStake(address bookNFT) external view returns (uint256) {
        CollectiveData storage $ = _getCollectiveData();
        return $.pools[bookNFT].totalStaked;
    }

    function getStakeForUser(
        address user,
        address bookNFT
    ) external view returns (uint256) {
        CollectiveData storage $ = _getCollectiveData();
        uint256[] memory positions = $
            .likeStakePosition
            .getUserPositionByBookNFT(user, bookNFT);
        uint256 totalStake = 0;
        for (uint256 i = 0; i < positions.length; ++i) {
            totalStake += $
                .likeStakePosition
                .getPosition(positions[i])
                .stakedAmount;
        }
        return totalStake;
    }

    // View function for pending rewards
    function getRewardsOfPosition(
        uint256 tokenId
    ) external view returns (uint256) {
        return _pendingRewardsOf(tokenId);
    }

    function getPendingRewardsPool(
        address bookNFT
    ) external view returns (uint256) {
        CollectiveData storage $ = _getCollectiveData();
        return $.pools[bookNFT].rewardPending;
    }

    function getPendingRewardsForUser(
        address user,
        address bookNFT
    ) external view returns (uint256) {
        CollectiveData storage $ = _getCollectiveData();
        // Projected index includes undistributed pending rewards
        uint256 poolIndex = $.pools[bookNFT].rewardIndex;
        uint256[] memory _positions = $
            .likeStakePosition
            .getUserPositionByBookNFT(user, bookNFT);

        uint256 totalPending = 0;
        for (uint256 i = 0; i < _positions.length; ++i) {
            LikeStakePosition.Position memory p = $
                .likeStakePosition
                .getPosition(_positions[i]);
            uint256 _pendingRewards = (p.stakedAmount *
                (poolIndex - p.rewardIndex)) / ACC_REWARD_PRECISION;
            totalPending += _pendingRewards;
        }
        return totalPending;
    }

    // solhint-enable no-unused-vars
}
