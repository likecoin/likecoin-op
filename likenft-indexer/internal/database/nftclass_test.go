package database_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"likenft-indexer/ent"
	"likenft-indexer/ent/schema/typeutil"
	"likenft-indexer/internal/database"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/smartystreets/goconvey/convey"
)

type nftClassUtil struct {
}

func (nftClassUtil) makeNFTClass(
	nftClassRepository database.NFTClassRepository,
	owner *ent.Account,
	address string,
	latestEventBlockHeight uint64,
) error {
	err := nftClassRepository.InsertNFTClass(
		context.Background(),
		address,
		"Test NFT Class",
		"TEST",
		nil,
		[]string{},
		big.NewInt(1000),
		typeutil.Uint64(10000),
		nil,
		"banner.jpg",
		"featured.jpg",
		common.HexToAddress("0x789").String(),
		typeutil.Uint64(latestEventBlockHeight),
		typeutil.Uint64(latestEventBlockHeight),
		time.Now().UTC(),
		owner,
	)
	return err
}

func TestNFTClass(t *testing.T) {
	util := &nftClassUtil{}

	Convey("Test UpdateNFTClassesLatestEventBlockNumber", t, func() {
		dbService := newTestService(t)
		defer dbService.Close()

		accountRepository := database.MakeAccountRepository(dbService)
		nftClassRepository := database.MakeNFTClassRepository(dbService)

		account, err := accountRepository.GetOrCreateAccount(context.Background(), &ent.Account{
			EvmAddress: common.HexToAddress("0x2333").String(),
		})
		if err != nil {
			t.Fatal(err)
		}

		err = util.makeNFTClass(
			nftClassRepository,
			account,
			common.HexToAddress("0x8Ed946DF65a6dD97fcD071c995710b9bf14361c3").String(),
			100,
		)
		if err != nil {
			t.Fatal(err)
		}

		err = util.makeNFTClass(
			nftClassRepository,
			account,
			common.HexToAddress("0x6eF99D096D32b4E5E243ABAD5f2284e8de81CE9d").String(),
			100,
		)
		if err != nil {
			t.Fatal(err)
		}

		err = util.makeNFTClass(
			nftClassRepository,
			account,
			common.HexToAddress("0xCc9A758392d4FD72Def70F4cB4d3D74002353e87").String(),
			100,
		)
		if err != nil {
			t.Fatal(err)
		}

		err = nftClassRepository.UpdateNFTClassesLatestEventBlockNumber(
			context.Background(),
			[]string{
				"0x8ed946df65a6dd97fcd071c995710b9bf14361c3",
				"0x6ef99d096d32b4e5e243abad5f2284e8de81ce9d",
			},
			200,
		)
		So(err, ShouldBeNil)
		nftClass, err := nftClassRepository.QueryNFTClassByAddress(
			context.Background(),
			common.HexToAddress("0x8ed946df65a6dd97fcd071c995710b9bf14361c3").String(),
		)
		So(err, ShouldBeNil)
		So(nftClass.LatestEventBlockNumber, ShouldEqual, 200)

		nftClass, err = nftClassRepository.QueryNFTClassByAddress(
			context.Background(),
			common.HexToAddress("0x6ef99d096d32b4e5e243abad5f2284e8de81ce9d").String(),
		)
		So(err, ShouldBeNil)
		So(nftClass.LatestEventBlockNumber, ShouldEqual, 200)

		nftClass, err = nftClassRepository.QueryNFTClassByAddress(
			context.Background(),
			common.HexToAddress("0xcc9a758392d4fd72def70f4cb4d3d74002353e87").String(),
		)
		So(err, ShouldBeNil)
		So(nftClass.LatestEventBlockNumber, ShouldEqual, 100)

	})

}
