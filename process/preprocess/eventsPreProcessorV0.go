package preprocess

import (
	"encoding/json"

	"github.com/multiversx/mx-chain-core-go/core"
	coreData "github.com/multiversx/mx-chain-core-go/data"
	nodeData "github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/process"
)

// TODO: dismiss this implementation after http integration is fully deprecated
type eventsPreProcessorV0 struct {
	*baseEventsPreProcessor
}

// NewEventsPreProcessorV0 will create a new events data preprocessor instance
func NewEventsPreProcessorV0(args ArgsEventsPreProcessor) (*eventsPreProcessorV0, error) {
	baseEventsPreProcessor, err := newBaseEventsPreProcessor(args)
	if err != nil {
		return nil, err
	}

	return &eventsPreProcessorV0{
		baseEventsPreProcessor: baseEventsPreProcessor,
	}, nil
}

// SaveBlock will handle the block info data
func (d *eventsPreProcessorV0) SaveBlock(marshalledData []byte) error {
	blockData := &data.OutportBlockDataOld{}
	err := json.Unmarshal(marshalledData, blockData)
	if err != nil {
		return err
	}

	header, err := d.getHeader(marshalledData)
	if err != nil {
		return err
	}

	txsPool, err := d.parseTransactionsPool(blockData.TransactionsPool)
	if err != nil {
		return err
	}

	saveBlockData := &data.ArgsSaveBlockData{
		HeaderHash:             blockData.HeaderHash,
		Body:                   blockData.Body,
		SignersIndexes:         blockData.SignersIndexes,
		NotarizedHeadersHashes: blockData.NotarizedHeadersHashes,
		HeaderGasConsumption:   &blockData.HeaderGasConsumption,
		AlteredAccounts:        blockData.AlteredAccounts,
		NumberOfShards:         blockData.NumberOfShards,
		TransactionsPool:       txsPool,
		Header:                 header,
	}

	err = d.facade.HandlePushEvents(*saveBlockData)
	if err != nil {
		return err
	}

	return nil
}

func (d *eventsPreProcessorV0) parseTransactionsPool(txsPool *data.TransactionsPool) (*outport.TransactionPool, error) {
	if txsPool == nil {
		return nil, process.ErrNilTransactionsPool
	}

	return &outport.TransactionPool{
		Transactions:         d.parseTxs(txsPool.Txs),
		SmartContractResults: d.parseScrs(txsPool.Scrs),
		Logs:                 d.parseLogs(txsPool.Logs),
	}, nil
}

func (d *eventsPreProcessorV0) parseTxs(txs map[string]*data.NodeTransaction) map[string]*outport.TxInfo {
	newTxs := make(map[string]*outport.TxInfo)
	for hash, txHandler := range txs {
		if txHandler == nil {
			continue
		}

		newTxs[hash] = &outport.TxInfo{
			Transaction:    txHandler.TransactionHandler,
			FeeInfo:        &txHandler.FeeInfo,
			ExecutionOrder: uint32(txHandler.ExecutionOrder),
		}
	}

	return newTxs
}

func (d *eventsPreProcessorV0) parseScrs(scrs map[string]*data.NodeSmartContractResult) map[string]*outport.SCRInfo {
	newScrs := make(map[string]*outport.SCRInfo)
	for hash, scrHandler := range scrs {
		if scrHandler == nil {
			continue
		}

		newScrs[hash] = &outport.SCRInfo{
			SmartContractResult: scrHandler.TransactionHandler,
			FeeInfo:             &scrHandler.FeeInfo,
			ExecutionOrder:      uint32(scrHandler.ExecutionOrder),
		}
	}

	return newScrs
}

func (d *eventsPreProcessorV0) parseLogs(logs []*data.LogData) []*outport.LogData {
	newLogs := make([]*outport.LogData, len(logs))
	for _, logHandler := range logs {
		if logHandler == nil {
			continue
		}

		newLogs = append(newLogs, &outport.LogData{
			TxHash: logHandler.TxHash,
			Log:    logHandler.LogHandler,
		})
	}

	return newLogs
}

func (d *eventsPreProcessorV0) getHeader(marshaledData []byte) (nodeData.HeaderHandler, error) {
	headerTypeStruct := struct {
		HeaderType core.HeaderType
	}{}

	err := json.Unmarshal(marshaledData, &headerTypeStruct)
	if err != nil {
		return nil, err
	}

	header, err := d.getHeaderFromBytes(headerTypeStruct.HeaderType, marshaledData)
	if err != nil {
		return nil, err
	}

	return header, nil
}

func (d *eventsPreProcessorV0) getHeaderFromBytes(headerType core.HeaderType, headerBytes []byte) (header coreData.HeaderHandler, err error) {
	creator, err := d.emptyBlockCreator.Get(headerType)
	if err != nil {
		return nil, err
	}

	return block.GetHeaderFromBytes(d.marshaller, creator, headerBytes)
}

// RevertIndexedBlock will handle the revert block event
func (d *eventsPreProcessorV0) RevertIndexedBlock(marshalledData []byte) error {
	revertBlock := &data.RevertBlock{}
	err := d.marshaller.Unmarshal(revertBlock, marshalledData)
	if err != nil {
		return err
	}

	d.facade.HandleRevertEvents(*revertBlock)

	return nil
}

// FinalizedBlock will handle the finalized block event
func (d *eventsPreProcessorV0) FinalizedBlock(marshalledData []byte) error {
	finalizedBlock := &data.FinalizedBlock{}
	err := d.marshaller.Unmarshal(finalizedBlock, marshalledData)
	if err != nil {
		return err
	}

	d.facade.HandleFinalizedEvents(*finalizedBlock)

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (d *eventsPreProcessorV0) IsInterfaceNil() bool {
	return d == nil
}
