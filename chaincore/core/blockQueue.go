package core

import  (
	"container/list"
)

type BlockQueue struct  {
	list.List
}


func NewBlockQueue() *BlockQueue {

	blockQueue := new(BlockQueue)
	return blockQueue

}

func (q *BlockQueue) Add(block *Block) {

	q.PushBack(block)

}

func (q *BlockQueue) GetTip() interface{} {

	for q.Len() > 0 {
		e := q.Front()
		return e.Value
	}
	return nil

}

