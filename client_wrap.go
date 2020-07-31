package bbrpc

/**
 * 该文件对一些rpc进行封装的函数
 *
 */

import (
	"fmt"
)

// ListTransactionsSinceBlock 列出自某个区块以来的所有交易，不包含targetBlock的交易
// return: scan2BlockHash, tx, error
func (c *Client) ListTransactionsSinceBlock(targetBlockHash string, count int) (string, []TransactionDetail, error) {
	var fork *string = nil
	const defaultRecentHeight = 30 //如果提供的hash为空，则取最近的n个块的交易

	topForkHeight, err := c.Getforkheight(fork)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get fork height, %v", err)
	}

	if targetBlockHash == "" { //默认取最近30个块吧
		defaultTargetBlockHeight := topForkHeight - defaultRecentHeight
		if defaultTargetBlockHeight < 1 {
			defaultTargetBlockHeight = 1
		}
		defaultTargetBlockHash, err := c.Getblockhash(int(defaultTargetBlockHeight), nil)
		if err != nil || len(defaultTargetBlockHash) == 0 {
			return "", nil, fmt.Errorf("failed to get default target block hash, %v, len(hash): %d", err, len(defaultTargetBlockHash))
		}
		targetBlockHash = defaultTargetBlockHash[0]
	}

	topBlockHash, err := c.Getblockhash(int(topForkHeight), nil)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get top block hash, %v", err)
	}

	var all []TransactionDetail
	prevBlockHash := targetBlockHash

	block, err := c.Getblock(prevBlockHash)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get block [%s]", prevBlockHash)
	}
	if block.Height == uint(topForkHeight) { //已经是最新高度了
		return topBlockHash[0], []TransactionDetail{}, nil
	}
	prevBlockHeight := block.Height

	scanBlockCount := 0
	for {
		prevBlockHeight++
		blockHash, err := c.Getblockhash(int(prevBlockHeight), fork)
		if err != nil {
			return prevBlockHash, nil, fmt.Errorf("failed to get block hash @ [%d], %v", prevBlockHeight, err)
		}

		block, err := c.Getblock(blockHash[0])
		if err != nil {
			return prevBlockHash, nil, fmt.Errorf("failed to get block @ [%s], %v", blockHash, err)
		}

		if len(block.Tx) > 0 {
			for _, txid := range block.Tx {
				tx, err := c.Gettransaction(txid, pbool(false))
				if err != nil {
					return prevBlockHash, nil, fmt.Errorf("failed to get transaction [%s] at [%s(%d)]", txid, block.Hash, block.Height)
				}
				all = append(all, *tx)
			}
		}
		prevBlockHash = block.Hash

		scanBlockCount++
		if block.Height == uint(topForkHeight) || scanBlockCount >= count {
			break
		}
	}
	return prevBlockHash, all, nil
}

// ListBlockDetailsSince 列出自某个区块以来的所有区块详情，不包含targetBlock的交易
// return: topBlockHeight, blockDetails, error
func (c *Client) ListBlockDetailsSince(fork *string, targetBlockHash string, count int) (int, []BlockDetail, error) {
	// fmt.Println("[dbg]ListBlockDetailsSince", *fork, targetBlockHash, count)
	const defaultRecentHeight = 30 //如果提供的hash为空，则取最近的n个块的交易

	topForkHeight, err := c.Getforkheight(fork)
	if err != nil {
		return -1, nil, fmt.Errorf("failed to get fork height, %v", err)
	}

	if targetBlockHash == "" { //默认取最近30个块吧
		defaultTargetBlockHeight := topForkHeight - defaultRecentHeight
		if defaultTargetBlockHeight < 1 {
			defaultTargetBlockHeight = 1
		}
		defaultTargetBlockHash, err := c.Getblockhash(int(defaultTargetBlockHeight), fork)
		if err != nil || len(defaultTargetBlockHash) == 0 {
			return topForkHeight, nil, fmt.Errorf("failed to get default target block hash, %v, len(hash): %d", err, len(defaultTargetBlockHash))
		}
		targetBlockHash = defaultTargetBlockHash[0]
	}

	var all []BlockDetail

	var prevBlockHeight uint
	{
		block, err := c.Getblock(targetBlockHash)
		if err != nil {
			return topForkHeight, nil, fmt.Errorf("failed to get block [%s]", targetBlockHash)
		}
		if block.Height == uint(topForkHeight) { //已经是最新高度了
			return topForkHeight, nil, nil
		}
		prevBlockHeight = block.Height
	}

	scannedBlockCount := 0
	for {
		prevBlockHeight++
		blockHash, err := c.Getblockhash(int(prevBlockHeight), fork)
		if err != nil || len(blockHash) == 0 {
			return topForkHeight, nil, fmt.Errorf("failed to get block hash @ [%d], result len: [%d], %v", prevBlockHeight, len(blockHash), err)
		}

		var lastHeight uint
		for _, hash := range blockHash {
			detail, err := c.Getblockdetail(hash)
			if err != nil {
				return topForkHeight, nil, fmt.Errorf("failed to get block detail @ [%s], %v", blockHash, err)
			}
			all = append(all, *detail)
			lastHeight = detail.Height
		}
		if lastHeight == uint(topForkHeight) || scannedBlockCount >= count {
			break
		}
		scannedBlockCount++
	}
	return topForkHeight, all, nil
}

// ListBlockDetailsBetween 列出[from,to]内所有区块详情
// return: blockDetails, error
func (c *Client) ListBlockDetailsBetween(fork *string, fromHeight, toHeight int) ([]BlockDetail, error) {
	var all []BlockDetail
	cursorHeight := fromHeight
	for {
		blockHash, err := c.Getblockhash(cursorHeight, fork)
		if err != nil || len(blockHash) == 0 {
			return nil, fmt.Errorf("failed to get block hash @ [%d], result len: [%d], %v", cursorHeight, len(blockHash), err)
		}

		var lastHeight uint
		for _, hash := range blockHash {
			detail, err := c.Getblockdetail(hash)
			if err != nil {
				return nil, fmt.Errorf("failed to get block detail @ [%s], %v", blockHash, err)
			}
			all = append(all, *detail)
			lastHeight = detail.Height
		}
		if int(lastHeight) == toHeight {
			break
		}
		cursorHeight++
	}
	return all, nil
}
