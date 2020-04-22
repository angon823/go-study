package bag

import (
	"fmt"
	"time"
)

//>> 容器类型
type ContainerType int16

const (
	KContainerTypeInvalid ContainerType = iota
	KContainerTypeBag                   // 背包
)

//>> 道具更改原因
type ItemChangeReason int

//>> 道具描述信息
type ItemTidDesc struct {
	TID   int32
	Count int64
}

type ItemUidDesc struct {
	UID   uint64
	Count int64
}

const (
	ErrSuccess = iota
	ErrItemNotExist
	ErrItemNotEnough
	ErrContainerNotExist
)

type itemError struct {
	Code   int    //>> 自定义错误码
	Detail string //>> 描述
	Param  []int  //>> 自定义错误参数，如道具数量不足错误，Param可以为[id,Count]
}

type ItemError *itemError

func (this *itemError) Error() string {
	return fmt.Sprintf("%d %s %v", this.Code, this.Detail, this.Param)
}

func NewItemError(code int) *itemError {
	err := itemError{Code: code, Param: make([]int, 1)}
	return &err
}

//>> 容器接口
type ContainerInterface interface {
	//>> 容器类型
	GetType() ContainerType
	//>> 占用格子数
	GetSize() int32
	//>> 负重
	GetLoad() int32
	//>> 最大格子数量
	GetMaxSize() int32
	//>> 负重上限
	GetMaxLoad() int32
	//>> 根据uid返回道具
	GetItemByUID(uid uint64) ItemInterface
	//>> 根据模板id返回道具列表
	GetItemsByTID(tid int32) []ItemInterface
	//>> 根据模板id返回一个“最佳删除”道具，可能优先绑定或快过期的
	GetItemForReduce(tid int32) ItemInterface
	//>> 获取道具数量
	GetItemCount(tid int32) int64
	//>> 检查是否能加道具
	TryAddItem(tid int32, count int64) ItemError
	//>> 检查是否能加道具
	TryAddItems(items []ItemTidDesc) ItemError
	//>> 增加道具，成功返回增加后的道具, 因为堆叠原因也可能有多个
	AddItem(tid int32, count int64, reason ItemChangeReason) ([]ItemInterface, ItemError)
	//>> 批量加道具, 成功返回增加后的道具切片
	AddItems(items []ItemTidDesc, reason ItemChangeReason) ([]ItemInterface, ItemError)
	//>> 检查是否能扣道具，成功返回nil
	TryReduceItemByUID(uid uint64, count int64) ItemError
	//>> 检查是否能扣道具，成功返回nil
	TryReduceItemByTID(tid int32, count int64) ItemError
	//>> 检查是否能扣道具，成功返回nil
	TryReduceItems([]ItemTidDesc) ItemError
	//>> 扣道具，成功返回nil
	ReduceItemByUID(uid uint64, count int64, reason ItemChangeReason) ItemError
	//>> 扣道具，成功返回nil
	ReduceItemByTID(tid int32, count int64, reason ItemChangeReason) ItemError
	//>> 扣道具，成功返回nil
	ReduceItems(items []ItemTidDesc, reason ItemChangeReason) ItemError
	//>> 扣并给道具，保证事务性
	ReduceAndAddItems(delItems, giveItems []ItemTidDesc, reason ItemChangeReason) ItemError
	//>> 扣并给道具，保证事务性
	ReduceAndAddItemByUID(delUIDs []ItemUidDesc, giveItems []ItemTidDesc, reason ItemChangeReason) ItemError
	//>> todo: 异步扣道具，结果调用回调, 有些一级货币可能无法同步扣除
	//AsyncReduceItem(items []ItemTidDesc, reason ItemChangeReason, cb func(err ItemError))
	//>> todo: 移动位置
	//MoveItem(uid uint64, targetPos int64) ItemError
	//>> todo：交换位置
	//SwapPosition(dstPos, scrPos int16) ItemError
	//>> todo: 如果有穿脱装备，会跨容器交换位置
	//SwapOut()
	//SwapIn()
}

// 道具更新类型
type ItemUpdateType int8

const (
	KItemUpdateTypeAdd    ItemUpdateType = iota // 增加
	KItemUpdateTypeDel                          // 删除
	KItemUpdateTypeUpdate                       // 更新
)

// 道具操作更新
type ItemOpRecord struct {
	UID       uint64
	Operation ItemUpdateType
}

type ContainerBase struct {
	typ      ContainerType
	curSize  int32
	maxSize  int32
	curLoad  int32
	maxLoad  int32
	items    map[uint64]ItemInterface
	tid2UIDs map[int32][]uint64

	// 更新队列
	updateQueue []ItemOpRecord
}

//>> 容器类型
func (this *ContainerBase) GetType() ContainerType {
	return this.typ
}

//>> 占用格子数
func (this *ContainerBase) GetSize() int32 {
	return this.curSize
}

//>> 负重
func (this *ContainerBase) GetLoad() int32 {
	return this.curLoad
}

//>> 最大格子数量
func (this *ContainerBase) GetMaxSize() int32 {
	return this.maxSize
}

//>> 负重上限
func (this *ContainerBase) GetMaxLoad() int32 {
	return this.maxLoad
}

//>> 根据uid返回道具
func (this *ContainerBase) GetItemByUID(uid uint64) ItemInterface {
	return this.items[uid]
}

//>> 根据模板id返回道具列表
func (this *ContainerBase) GetItemsByTID(tid int32) []ItemInterface {
	var ret []ItemInterface
	if tids, ok := this.tid2UIDs[tid]; ok {
		for uid := range tids {
			item, ok := this.items[uint64(uid)]
			if ok {
				ret = append(ret, item)
			}
		}
	}
	return ret
}

//>> 根据模板id返回一个“最佳删除”道具，可能优先绑定或快过期的
func (this *ContainerBase) GetItemForReduce(tid int32) ItemInterface {
	items := this.GetItemsByTID(tid)
	if len(items) > 0 {
		// todo: 根据规则找出先扣除的
		return items[0]
	}
	return nil
}

//>> 获取道具数量
func (this *ContainerBase) GetItemCount(tid int32) int64 {
	items := this.GetItemsByTID(tid)
	count := int64(0)
	for _, item := range items {
		count += item.GetCount()
	}
	return count
}

//>> 检查是否能加道具
func (this *ContainerBase) TryAddItem(tid int32, count int64) ItemError {
	//>> 计算负重
	//>> 计算格子
	//>> 计算堆叠
	return NewItemError(0)
}

//>> 检查是否能加道具
func (this *ContainerBase) TryAddItems(items []ItemTidDesc) ItemError {
	itemMap := ItemTidDesc{}.convertToMap(items)
	return this.tryAddItems(itemMap)
}

func (this *ContainerBase) tryAddItems(itemMap map[int32]int64) ItemError {
	for k, v := range itemMap {
		if err := this.TryAddItem(k, v); err != nil {
			return err
		}
	}
	return nil
}

//>> 增加道具，成功返回增加后的道具
func (this *ContainerBase) AddItem(tid int32, count int64, reason ItemChangeReason) ([]ItemInterface, ItemError) {
	if err := this.TryAddItem(tid, count); err != nil {
		return nil, err
	}

	var items []ItemInterface
	// todo: 计算堆叠
	for n := 1; n > 0; n-- {

		item := NewItem(tid, count)
		if item == nil {
			return nil, NewItemError(-1)
		}

		this.addItem(item, reason)
		items = append(items, item)
	}

	return items, nil
}

//>> 批量加道具, 成功返回增加后的道具切片
func (this *ContainerBase) AddItems(items []ItemTidDesc, reason ItemChangeReason) ([]ItemInterface, ItemError) {
	itemMap := ItemTidDesc{}.convertToMap(items)
	if err := this.tryAddItems(itemMap); err != nil {
		return nil, err
	}

	var ret []ItemInterface
	for tid, count := range itemMap {
		// todo: 计算堆叠
		total := int64(0)
		cur := int64(1)
		for {
			//maxOverlap := getItemMaxOverlap(itemDesc.TID)
			//curOverlap
			total += cur
			item := NewItem(tid, cur)
			if item == nil {
				return nil, NewItemError(-1)
			}
			ret = append(ret, item)
			this.addItem(item, reason)

			if total >= count {
				break
			}
		}
	}

	return ret, nil
}

//>> 检查是否能扣道具，成功返回nil
func (this *ContainerBase) TryReduceItemByUID(uid uint64, count int64) ItemError {
	item := this.GetItemByUID(uid)
	if item == nil {
		return NewItemError(ErrItemNotExist)
	}

	if item.GetCount() < count {
		err := NewItemError(ErrItemNotEnough)
		err.Param = append(err.Param, int(count-item.GetCount()))
		return err
	}

	return nil
}

//>> 检查是否能扣道具，成功返回nil
func (this *ContainerBase) TryReduceItemByTID(tid int32, count int64) ItemError {
	items := this.GetItemsByTID(tid)
	if len(items) == 0 {
		return NewItemError(ErrItemNotExist)
	}

	has := int64(0)
	for _, item := range items {
		has += item.GetCount()
		if count <= has {
			return nil
		}
	}

	err := NewItemError(ErrItemNotEnough)
	err.Param = append(err.Param, int(count-has))
	return err
}

//>> 检查是否能扣道具，成功返回nil
func (this *ContainerBase) TryReduceItems(items []ItemTidDesc) ItemError {
	itemMap := ItemTidDesc{}.convertToMap(items)
	if itemMap == nil {
		return nil
	}

	for k, v := range itemMap {
		if err := this.TryReduceItemByTID(k, v); err != nil {
			return err
		}
	}

	return nil
}

//>> 扣道具，成功返回nil
func (this *ContainerBase) ReduceItemByUID(uid uint64, count int64, reason ItemChangeReason) ItemError {
	if err := this.TryReduceItemByUID(uid, count); err != nil {
		return err
	}

	this.delItem(uid, count, reason)
	return nil
}

//>> 扣道具，成功返回nil
func (this *ContainerBase) ReduceItemByTID(tid int32, count int64, reason ItemChangeReason) ItemError {
	if err := this.TryReduceItemByTID(tid, count); err != nil {
		return err
	}

	for {
		item := this.GetItemForReduce(tid)
		if item == nil {
			panic("(this *ContainerBase) ReduceItemByTID()!")
		}
		if item.GetCount() > count {
			this.delItem(item.GetUID(), count, reason)
		} else {
			this.delItem(item.GetUID(), item.GetCount(), reason)
		}
		count -= item.GetCount()
		if count <= 0 {
			break
		}
	}

	return nil
}

//>> 扣道具，成功返回nil
func (this *ContainerBase) ReduceItems(items []ItemTidDesc, reason ItemChangeReason) ItemError {
	if err := this.TryReduceItems(items); err != nil {
		return nil
	}

	for _, v := range items {
		tid := v.TID
		count := v.Count
		for {
			item := this.GetItemForReduce(tid)
			if item == nil {
				panic("(this *ContainerBase) ReduceItemByTID()!")
			}
			if item.GetCount() > count {
				this.delItem(item.GetUID(), count, reason)
			} else {
				this.delItem(item.GetUID(), item.GetCount(), reason)
			}
			count -= item.GetCount()
			if count <= 0 {
				break
			}
		}
	}

	return nil
}

//>> 扣并给道具，保证事务性
func (this *ContainerBase) ReduceAndAddItems(delItems, giveItems []ItemTidDesc, reason ItemChangeReason) ItemError {
	if err := this.TryAddItems(giveItems); err != nil {
		return err
	}

	if err := this.ReduceItems(delItems, reason); err != nil {
		return err
	}

	this.AddItems(giveItems, reason)
	return nil
}

//>> 扣并给道具，保证事务性
func (this *ContainerBase) ReduceAndAddItemByUID(delItems []ItemUidDesc, giveItems []ItemTidDesc, reason ItemChangeReason) ItemError {
	if err := this.TryAddItems(giveItems); err != nil {
		return err
	}

	delItemMap := ItemUidDesc{}.convertToMap(delItems)
	if delItemMap != nil {
		for k, v := range delItemMap {
			if err := this.TryReduceItemByUID(k, v); err != nil {
				return err
			}
		}

		for k, v := range delItemMap {
			if err := this.ReduceItemByUID(k, v, reason); err != nil {
				panic("this.ReduceItemByUID")
			}
		}
	}

	this.AddItems(giveItems, reason)
	return nil
}

func (this *ContainerBase) addItem(item ItemInterface, reason ItemChangeReason) {

	item.SetCreateTime(time.Now().Unix())
	item.SetContainerType(this.GetType())

	//todo
	newGrid := int16(0)
	item.SetPos(newGrid)

	//todo: 绑定信息等
	//item.SetFlag()

	this.items[item.GetUID()] = item

	ii, ok := this.tid2UIDs[item.GetTID()]
	if !ok {
		ii = make([]uint64, 1)
	}
	ii = append(ii, item.GetUID())
	this.tid2UIDs[item.GetTID()] = ii

	this.updateQueue = append(this.updateQueue, ItemOpRecord{UID: item.GetUID(), Operation: KItemUpdateTypeAdd})
}

func (this *ContainerBase) delItem(uid uint64, count int64, reason ItemChangeReason) {
	item := this.GetItemByUID(uid)
	if item == nil {
		panic("(this *ContainerBase) delItem")
	}

	left := item.GetCount() - count
	if left > 0 {
		item.SetCount(left)
	} else if left == 0 {
		uids := this.tid2UIDs[item.GetTID()]
		if len(uids) == 0 {
			panic("(this *ContainerBase) delItem len(uids) == 0")
		} else if len(uids) == 1 {
			delete(this.tid2UIDs, item.GetTID())
		} else {
			idx := 0
			for i := 0; i < len(uids); i++ {
				if uids[i] == uid {
					idx = i
					break
				}
			}
			uids = append(uids[0:idx], uids[idx+1:]...)
			this.tid2UIDs[item.GetTID()] = uids
		}
		delete(this.items, uid)
	} else {
		panic("(this *ContainerBase) delItem left < 0")
	}

	op := ItemOpRecord{UID: uid}
	if left == 0 {
		op.Operation = KItemUpdateTypeDel
	} else {
		op.Operation = KItemUpdateTypeUpdate
	}

	this.updateQueue = append(this.updateQueue, op)
}

func (ItemUidDesc) convertToMap(items []ItemUidDesc) map[uint64]int64 {
	if len(items) == 0 {
		return nil
	}
	itemMap := make(map[uint64]int64, len(items))
	for _, item := range items {
		itemMap[item.UID] += item.Count
	}
	return itemMap
}

func (ItemTidDesc) convertToMap(items []ItemTidDesc) map[int32]int64 {
	if len(items) == 0 {
		return nil
	}
	itemMap := make(map[int32]int64, len(items))
	for _, item := range items {
		itemMap[item.TID] += item.Count
	}
	return itemMap
}

//>> 根据道具ID从配置表里找到道具类型
func getItemType(tid int32) int32 {
	return 0
}

//>> 返回道具负重
func getItemWeight(tid int32) int32 {
	return 0
}

//>> 返回道具最大堆叠
func getItemMaxOverlap(tid int32) int {
	return 0
}

//>> 根据类型判断道具默认在哪个的容器
func getItemContainerType(tid int32) ContainerType {
	retTyp := KContainerTypeInvalid
	switch getItemType(tid) {
	case 0:
		retTyp = KContainerTypeBag
	default:
		retTyp = KContainerTypeBag
	}
	return ContainerType(retTyp)
}

func NewItem(tid int32, count int64) ItemInterface {
	switch getItemType(tid) {
	case 0:
		// todo: uid
		return &ItemBase{uid: 0, tid: tid, count: count}
	default:
	}
	return nil
}
