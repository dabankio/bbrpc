package bbrpc

/**
 * 该文件对一些rpc进行封装的函数
 *
 */

import (
	"fmt"
)

// ListTransactionsSinceBlock 列出自某个区块以来的所有交易，不包含targetBlock的交易
// return: topBlockHash, tx, error
func (c *Client) ListTransactionsSinceBlock(targetBlockHash string) (string, []TransactionDetail, error) {
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
	for {
		prevBlockHeight++
		blockHash, err := c.Getblockhash(int(prevBlockHeight), fork)
		if err != nil {
			return topBlockHash[0], nil, fmt.Errorf("failed to get block hash @ [%d], %v", prevBlockHeight, err)
		}

		block, err := c.Getblock(blockHash[0])
		if err != nil {
			return topBlockHash[0], nil, fmt.Errorf("failed to get block @ [%s], %v", blockHash, err)
		}

		prevBlockHash = block.Hash

		if len(block.Tx) > 0 {
			for _, txid := range block.Tx {
				tx, err := c.Gettransaction(txid, pbool(false))
				if err != nil {
					return topBlockHash[0], nil, fmt.Errorf("failed to get transaction [%s] at [%s(%d)]", txid, block.Hash, block.Height)
				}
				all = append(all, *tx)
			}
		}
		if block.Height == uint(topForkHeight) {
			break
		}
	}
	return topBlockHash[0], all, nil
}
