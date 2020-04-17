package bag

import (
	"fmt"
	"time"
)

//>> 容器类型
type ContainerType int16

const (
	KContainerTypeInvalid = iota
	KContainerTypeBag     // 背包
)

//>> 道具更改原因
type ItemChangeReason int

//>> 道具描述信息
type ItemDesc struct {
	tid   int32
	count uint64
}

type itemError struct {
	Code   int    //>> 自定义错误码
	Detail string //>> 描述
	Param  []int  //>> 自定义错误参数，如道具数量不足错误，Param可以为[id,count]
}

func (this *itemError) Error() string {
	return fmt.Sprintf("%d %s %v", this.Code, this.Detail, this.Param)
}

type ItemError *itemError

func NewItemError(code int) ItemError {
	err := itemError{Code: code}
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
	GetItemByTID(tid int32) []ItemInterface
	//>> 根据模板id返回一个“最佳删除”道具，可能优先绑定或快过期的
	GetItemForReduce(tid int32) ItemInterface
	//>> 获取道具数量
	GetItemCount(tid int32) uint64
	//>> 检查是否能加道具
	TryAddItem(tid int32, count uint64) ItemError
	//>> 检查是否能加道具
	TryAddItems(items []ItemDesc) ItemError
	//>> 增加道具，成功返回增加后的道具
	AddItem(tid int32, count uint64, reason ItemChangeReason) ItemInterface
	//>> 批量加道具, 成功返回增加后的道具切片
	AddItems(items []ItemDesc, reason ItemChangeReason) []ItemInterface
	//>> 检查是否能扣道具，成功返回nil
	TryReduceItemByUID(uid, count uint64) ItemError
	//>> 检查是否能扣道具，成功返回nil
	TryReduceItemByTID(tid int32, count uint64) ItemError
	//>> 检查是否能扣道具，成功返回nil
	TryReduceItems([]ItemDesc) ItemError
	//>> 扣道具，成功返回nil
	ReduceItemByUID(uid, count uint64, reason ItemChangeReason) ItemError
	//>> 扣道具，成功返回nil
	ReduceItem(tid int32, count uint64, reason ItemChangeReason) ItemError
	//>> 扣道具，成功返回nil
	ReduceItems(items []ItemDesc, reason ItemChangeReason) ItemError
	//>> 扣并给道具，保证事务性
	ReduceAndAddItems(delItems, giveItems []ItemDesc, reason ItemChangeReason) ItemError
	//>> 扣并给道具，保证事务性
	ReduceAndAddItemByUID(delUIDs []uint64, giveItems []ItemDesc, reason ItemChangeReason) ItemError
	//>> todo: 异步扣道具，结果调用回调, 有些一级货币可能无法同步扣除
	//AsyncReduceItem(items []ItemDesc, reason ItemChangeReason, cb func(err ItemError))
	//>> todo: 移动位置
	//MoveItem(uid uint64, targetPos int64) ItemError
	//>> todo：交换位置
	//SwapPosition(dstPos, scrPos int16) ItemError
	//>> todo: 如果有穿脱装备，会跨容器交换位置
	//SwapOut()
	//SwapIn()
}

type ContainerBase struct {
	typ      ContainerType
	curSize  int32
	maxSize  int32
	curLoad  int32
	maxLoad  int32
	items    map[uint64]ItemInterface
	tid2UIDs map[int32][]uint64
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
	return  this.items[uid]
}

//>> 根据模板id返回道具列表
func (this *ContainerBase) GetItemByTID(tid int32) []ItemInterface {
	var ret []ItemInterface
	if tids, ok := this.tid2UIDs[tid]; ok{
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
	items := this.GetItemByTID(tid)
	if len(items) > 0 {
		// todo: 根据规则找出先扣除的
		return items[0]
	}
	return nil
}

//>> 获取道具数量
func (this *ContainerBase) GetItemCount(tid int32) uint64 {
	items  := this.GetItemByTID(tid)
	count := uint64(0)
	for _, item := range items {
		count += item.GetCount()
	}
	return count
}

//>> 检查是否能加道具
func (this *ContainerBase) TryAddItem(tid int32, count uint64) ItemError {
	//>> 计算堆叠
	//>> 计算负重
	//>> 计算格子
	return NewItemError(0)
}

//>> 检查是否能加道具
func (this *ContainerBase) TryAddItems(items []ItemDesc) ItemError {
	for _, item := range items {
		if err := this.TryAddItem(item.tid, item.count); err != nil {
			//>> todo 合并错误
			return err
		}
	}
	return nil
}

//>> 增加道具，成功返回增加后的道具
func (this *ContainerBase) AddItem(tid int32, count uint64, reason ItemChangeReason) ItemInterface {
	if this.TryAddItem(tid, count) != nil {
		return  nil
	}

	item := NewItem(tid)
	if item == nil {
		return nil
	}

	this.addItem(item, reason)


	return item
}

//>> 批量加道具, 成功返回增加后的道具切片
func (this *ContainerBase) AddItems(items []ItemDesc, reason ItemChangeReason) []ItemInterface {

}

//>> 检查是否能扣道具，成功返回nil
func (this *ContainerBase) TryReduceItemByUID(uid, count uint64) ItemError {

}

//>> 检查是否能扣道具，成功返回nil
func (this *ContainerBase) TryReduceItemByTID(tid int32, count uint64) ItemError {

}

//>> 检查是否能扣道具，成功返回nil
func (this *ContainerBase) TryReduceItems([]ItemDesc) ItemError {

}

//>> 扣道具，成功返回nil
func (this *ContainerBase) ReduceItemByUID(uid, count uint64, reason ItemChangeReason) ItemError {

}

//>> 扣道具，成功返回nil
func (this *ContainerBase) ReduceItem(tid int32, count uint64, reason ItemChangeReason) ItemError {

}

//>> 扣道具，成功返回nil
func (this *ContainerBase) ReduceItems(items []ItemDesc, reason ItemChangeReason) ItemError {

}

//>> 扣并给道具，保证事务性
func (this *ContainerBase) ReduceAndAddItems(delItems, giveItems []ItemDesc, reason ItemChangeReason) ItemError {

}

//>> 扣并给道具，保证事务性
func (this *ContainerBase) ReduceAndAddItemByUID(delUIDs []uint64, giveItems []ItemDesc, reason ItemChangeReason) ItemError {

}

func (this *ContainerBase) addItem(item ItemInterface, reason ItemChangeReason) ItemError {
	//todo
	newGrid := int16(0)

	item.SetTID(tid)
	item.SetCount(count)
	item.SetType(getItemType(tid))
	item.SetCreateTime(time.Now().Unix())
	item.SetContainerType(this.GetType())
	item.SetPos(newGrid)

	//todo: 绑定信息等
	//item.SetFlag()

	return  nil
}

//>> 根据道具ID从配置表里找到道具类型
func getItemType(tid int32) int16 {
	return 0
}

//>> 返回道具负重
func getItemWeight(tid int32) int {
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

func NewItem(tid int32) ItemInterface{
	switch getItemType(tid) {
	case 0:
		// todo: uid
		return &ItemBase{uid:0, tid: tid}
	default:
	}
	return  nil
}

//>> 背包组件
type ItemComponent struct {
	containers map[ContainerType]ContainerInterface
}
