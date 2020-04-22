package bag

import "fmt"

type Bag struct {
	ContainerBase
}

//>> 背包组件
type ItemComponent struct {
	//>> 容器
	containers map[ContainerType]ContainerInterface
}

func (this *ItemComponent) Init() {
	this.containers[KContainerTypeBag] = &Bag{}
}

func (this *ItemComponent) Update() {
	for _, container := range this.containers {
		//>> 更新至客户端
		fmt.Println(container.(*ContainerBase).updateQueue)
	}
}

func (this *ItemComponent) GetContainerByTID(tid int32) ContainerInterface {
	typ := getItemContainerType(tid)
	return this.GetContainerByType(typ)
}

func (this *ItemComponent) GetContainerByType(typ ContainerType) ContainerInterface {
	return this.containers[typ]
}

//>> 根据uid返回道具
func (this *ItemComponent) GetItemByUID(uid uint64) ItemInterface {
	for _, container := range this.containers {
		if item := container.GetItemByUID(uid); item != nil {
			return item
		}
	}
	return nil
}

//>> 根据模板id返回道具列表
func (this *ItemComponent) GetItemsByTID(tid int32) []ItemInterface {
	container := this.GetContainerByTID(tid)
	if container != nil {
		return container.GetItemsByTID(tid)
	}
	return []ItemInterface{}
}

//>> 获取道具数量
func (this *ItemComponent) GetItemCount(tid int32) int64 {
	container := this.GetContainerByTID(tid)
	if container != nil {
		return container.GetItemCount(tid)
	}
	return 0
}

//>> 检查是否能加道具
func (this *ItemComponent) TryAddItem(typ ContainerType, tid int32, count int64) ItemError {
	container := this.GetContainerByType(typ)
	if container != nil {
		return container.TryAddItem(tid, count)
	}
	return NewItemError(ErrContainerNotExist)
}

//>> 检查是否能加道具
func (this *ItemComponent) TryAddItems(typ ContainerType, items []ItemTidDesc) ItemError {
	container := this.GetContainerByType(typ)
	if container != nil {
		return container.TryAddItems(items)
	}
	return NewItemError(ErrContainerNotExist)
}

//>> 增加道具，成功返回增加后的道具, 因为堆叠原因也可能有多个
func (this *ItemComponent) AddItem(typ ContainerType, tid int32, count int64, reason ItemChangeReason) ([]ItemInterface, ItemError) {
	container := this.GetContainerByType(typ)
	if container != nil {
		return container.AddItem(tid, count, reason)
	}
	return nil, NewItemError(ErrContainerNotExist)
}

//>> 批量加道具, 成功返回增加后的道具切片
func (this *ItemComponent) AddItems(typ ContainerType, items []ItemTidDesc, reason ItemChangeReason) ([]ItemInterface, ItemError) {
	container := this.GetContainerByType(typ)
	if container != nil {
		return container.AddItems(items, reason)
	}

	return nil, NewItemError(ErrContainerNotExist)
}

//>> 检查是否能扣道具，成功返回nil
func (this *ItemComponent) TryReduceItemByUID(uid uint64, count int64) ItemError {
	for _, container := range this.containers {
		if container.TryReduceItemByUID(uid, count) == nil {
			return nil
		}
	}
	return NewItemError(ErrItemNotExist)
}

//>> 检查是否能扣道具，成功返回nil
func (this *ItemComponent) TryReduceItemByTID(typ ContainerType, tid int32, count int64) ItemError {
	container := this.GetContainerByType(typ)
	if container != nil {
		return container.TryReduceItemByTID(tid, count)
	}
	return NewItemError(ErrContainerNotExist)
}

//>> 检查是否能扣道具，成功返回nil
func (this *ItemComponent) TryReduceItems(typ ContainerType, items []ItemTidDesc) ItemError {
	container := this.GetContainerByType(typ)
	if container != nil {
		return container.TryReduceItems(items)
	}
	return NewItemError(ErrContainerNotExist)
}

//>> 扣道具，成功返回nil
func (this *ItemComponent) ReduceItemByUID(uid uint64, count int64, reason ItemChangeReason) ItemError {
	for _, container := range this.containers {
		if container.ReduceItemByUID(uid, count, reason) == nil {
			return nil
		}
	}
	return NewItemError(ErrItemNotExist)
}

//>> 扣道具，成功返回nil
func (this *ItemComponent) ReduceItemByTID(typ ContainerType, tid int32, count int64, reason ItemChangeReason) ItemError {
	container := this.GetContainerByType(typ)
	if container != nil {
		return container.ReduceItemByTID(tid, count, reason)
	}
	return NewItemError(ErrContainerNotExist)
}

//>> 扣道具，成功返回nil
func (this *ItemComponent) ReduceItems(typ ContainerType, items []ItemTidDesc, reason ItemChangeReason) ItemError {
	container := this.GetContainerByType(typ)
	if container != nil {
		return container.ReduceItems(items, reason)
	}
	return NewItemError(ErrContainerNotExist)
}

//>> 扣并给道具，保证事务性
func (this *ItemComponent) ReduceAndAddItems(typ ContainerType, delItems, giveItems []ItemTidDesc, reason ItemChangeReason) ItemError {
	container := this.GetContainerByType(typ)
	if container != nil {
		return container.ReduceAndAddItems(delItems, giveItems, reason)
	}
	return NewItemError(ErrContainerNotExist)
}

//>> 扣并给道具，保证事务性
func (this *ItemComponent) ReduceAndAddItemByUID(typ ContainerType, delUIDs []ItemUidDesc, giveItems []ItemTidDesc, reason ItemChangeReason) ItemError {
	container := this.GetContainerByType(typ)
	if container != nil {
		return container.ReduceAndAddItemByUID(delUIDs, giveItems, reason)
	}
	return NewItemError(ErrContainerNotExist)
}
