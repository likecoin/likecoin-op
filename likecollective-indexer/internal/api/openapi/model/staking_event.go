package model

import (
	"likecollective-indexer/ent"
	"likecollective-indexer/ent/stakingevent"
	"likecollective-indexer/internal/database"
	"likecollective-indexer/openapi/api"

	"github.com/holiman/uint256"
)

func MakeStakingEvent(stakingEvent *ent.StakingEvent) api.StakingEvent {
	switch stakingEvent.EventType {
	case stakingevent.EventTypeStaked:
		return api.StakingEvent{
			Type: api.StakingEventStakedStakingEvent,
			StakingEventStaked: api.StakingEventStaked{
				EventType: api.StakingEventStakedEventTypeStaked,
				BookNft:   api.EvmAddress(stakingEvent.NftClassAddress),
				Account:   api.EvmAddress(stakingEvent.AccountEvmAddress),
				Amount:    api.Uint256((*uint256.Int)(stakingEvent.StakedAmountAdded).String()),
				Datetime:  stakingEvent.Datetime,
			},
		}
	case stakingevent.EventTypeUnstaked:
		return api.StakingEvent{
			Type: api.StakingEventUnstakedStakingEvent,
			StakingEventUnstaked: api.StakingEventUnstaked{
				EventType: api.StakingEventUnstakedEventTypeUnstaked,
				BookNft:   api.EvmAddress(stakingEvent.NftClassAddress),
				Account:   api.EvmAddress(stakingEvent.AccountEvmAddress),
				Amount:    api.Uint256((*uint256.Int)(stakingEvent.StakedAmountRemoved).String()),
				Datetime:  stakingEvent.Datetime,
			},
		}
	case stakingevent.EventTypeRewardClaimed:
		return api.StakingEvent{
			Type: api.StakingEventRewardClaimedStakingEvent,
			StakingEventRewardClaimed: api.StakingEventRewardClaimed{
				EventType: api.StakingEventRewardClaimedEventTypeRewardClaimed,
				BookNft:   api.EvmAddress(stakingEvent.NftClassAddress),
				Account:   api.EvmAddress(stakingEvent.AccountEvmAddress),
				Amount:    api.Uint256((*uint256.Int)(stakingEvent.PendingRewardAmountRemoved).String()),
				Datetime:  stakingEvent.Datetime,
			},
		}
	case stakingevent.EventTypeRewardDeposited:
		return api.StakingEvent{
			Type: api.StakingEventRewardDepositedStakingEvent,
			StakingEventRewardDeposited: api.StakingEventRewardDeposited{
				EventType: api.StakingEventRewardDepositedEventTypeRewardDeposited,
				BookNft:   api.EvmAddress(stakingEvent.NftClassAddress),
				Account:   api.EvmAddress(stakingEvent.AccountEvmAddress),
				Amount:    api.Uint256((*uint256.Int)(stakingEvent.PendingRewardAmountAdded).String()),
				Datetime:  stakingEvent.Datetime,
			},
		}
	case stakingevent.EventTypeAllRewardsClaimed:
		return api.StakingEvent{
			Type: api.StakingEventAllRewardsClaimedStakingEvent,
			StakingEventAllRewardsClaimed: api.StakingEventAllRewardsClaimed{
				EventType: api.StakingEventAllRewardsClaimedEventTypeAllRewardsClaimed,
				Account:   api.EvmAddress(stakingEvent.AccountEvmAddress),
				ClaimedAmountList: []api.StakingEventAllRewardsClaimedClaimedAmountListItem{
					{
						BookNft: api.EvmAddress(stakingEvent.NftClassAddress),
						Amount:  api.Uint256((*uint256.Int)(stakingEvent.PendingRewardAmountRemoved).String()),
					},
				},
				Datetime: stakingEvent.Datetime,
			},
		}
	}

	panic("unknown staking event type: " + stakingEvent.EventType)
}

type StakingEventFilterParams struct {
	FilterNFTClassIn []api.EvmAddress
	FilterAccountIn  []api.EvmAddress
	FilterEventType  *api.StakingEventType
}

func (p *StakingEventFilterParams) ToEntFilter() database.QueryStakingEventsFilter {
	var filterNFTClassIn *[]string
	var filterAccountIn *[]string
	var filterEventType *string
	if len(p.FilterNFTClassIn) > 0 {
		_filterNFTClassIn := make([]string, len(p.FilterNFTClassIn))
		for i, nftClass := range p.FilterNFTClassIn {
			_filterNFTClassIn[i] = string(nftClass)
		}
		filterNFTClassIn = &_filterNFTClassIn
	}
	if len(p.FilterAccountIn) > 0 {
		_filterAccountIn := make([]string, len(p.FilterAccountIn))
		for i, account := range p.FilterAccountIn {
			_filterAccountIn[i] = string(account)
		}
		filterAccountIn = &_filterAccountIn
	}
	if p.FilterEventType != nil {
		_filterEventType := string(*p.FilterEventType)
		filterEventType = &_filterEventType
	}
	return database.NewQueryStakingEventsFilter(filterNFTClassIn, filterAccountIn, filterEventType)
}

type StakingEventPagination struct {
	// Pagination.limit.
	PaginationLimit api.OptInt
	// Pagination.key.
	PaginationKey api.OptInt
	// Reverse.
	Reverse api.OptBool
}

func (p *StakingEventPagination) ToEntPagination() database.StakingEventPagination {
	limit := FromOpt(p.PaginationLimit)
	if limit != nil && *limit == 0 {
		limit = nil
	}

	key := FromOpt(p.PaginationKey)
	if key != nil && *key == 0 {
		key = nil
	}

	reverse := FromOpt(p.Reverse)

	return database.StakingEventPagination{
		Limit:   limit,
		Key:     key,
		Reverse: reverse,
	}
}
