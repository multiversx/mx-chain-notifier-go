package groups_test

import (
	"encoding/json"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-notifier-go/api/groups"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/stretchr/testify/require"
)

func TestGetHeader(t *testing.T) {
	t.Parallel()

	t.Run("header v2", func(t *testing.T) {
		t.Parallel()

		header := &block.HeaderV2{
			Header: &block.Header{
				Nonce:     1,
				ShardID:   1,
				TimeStamp: 2,
			},
			ScheduledGasProvided:  1,
			ScheduledGasPenalized: 1,
			ScheduledGasRefunded:  1,
		}
		saveBlockData := &data.ArgsSaveBlockData{
			Header: header,
		}
		headerBytes, _ := json.Marshal(saveBlockData)

		h, err := groups.GetHeader(headerBytes)
		require.Nil(t, err)
		require.Equal(t, header, h)
	})

	// t.Run("header v1", func(t *testing.T) {
	// 	t.Parallel()

	// 	header := &block.Header{
	// 		Nonce:              1,
	// 		ShardID:            1,
	// 		TimeStamp:          1,
	// 		EpochStartMetaHash: []byte("epoch start meta hash"),
	// 	}
	// 	saveBlockData := &data.ArgsSaveBlockData{
	// 		Header: header,
	// 	}
	// 	headerBytes, _ := json.Marshal(saveBlockData)

	// 	h, err := groups.GetHeader(headerBytes)
	// 	require.Nil(t, err)
	// 	require.Equal(t, header, h)
	// })

	t.Run("metablock", func(t *testing.T) {
		t.Parallel()

		header := &block.MetaBlock{
			Nonce:     2,
			TimeStamp: 2,
		}
		saveBlockData := &data.ArgsSaveBlockData{
			Header: header,
		}
		headerBytes, _ := json.Marshal(saveBlockData)

		h, err := groups.GetHeader(headerBytes)
		require.Nil(t, err)
		require.Equal(t, header, h)
	})
}
