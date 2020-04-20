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
func (this *ItemComponent) GetItemCount(tid int32) uint64 {
	container := this.GetContainerByTID(tid)
	if container != nil {
		return container.GetItemCount(tid)
	}

	return 0
}

//>> 检查是否能加道具
func (this *ItemComponent) TryAddItem(tid int32, count uint64) ItemError {
	container := this.GetContainerByTID(tid)
	if container != nil {
		return container.TryAddItem(tid, count)
	}
	return NewItemError(ErrContainerNotExist)
}

//>> 检查是否能加道具
func (this *ItemComponent) TryAddItems(items []ItemTidDesc) ItemError {
	for _, item:= range items {
		container := this.GetContainerByTID(item.TID)
		if container != nil {
			return container.TryAddItems(items)
		}
		return NewItemError(ErrContainerNotExist)
	}
	return  nil
}
//
////>> 增加道具，成功返回增加后的道具, 因为堆叠原因也可能有多个
//func (this *ItemComponent) AddItem(tid int32, count uint64, reason ItemChangeReason) ([]ItemInterface, error) {
//	container := this.GetContainerByTID(tid)
//	if container != nil {
//		return container.AddItem(tid, count, reason)
//	}
//	return  []ItemInterface{}, NewItemError(ErrContainerNotExist)
//}
//
////>> 批量加道具, 成功返回增加后的道具切片
//func (this *ItemComponent) AddItems(typ ContainerType, items []ItemTidDesc, reason ItemChangeReason) ([]ItemInterface, error) {
//	container := this.GetContainerByType(typ)
//	if container == nil {
//		return  []ItemInterface{}, NewItemError(ErrContainerNotExist)
//	}
//
//	//return container.AddItems(items, reason)
//}
//
////>> 检查是否能扣道具，成功返回nil
//func (this *ItemComponent) TryReduceItemByUID(uid, count uint64) ItemError {
//
//}
//
////>> 检查是否能扣道具，成功返回nil
//func (this *ItemComponent) TryReduceItemByTID(tid int32, count uint64) ItemError {
//
//}
//
////>> 检查是否能扣道具，成功返回nil
//func (this *ItemComponent) TryReduceItems([]ItemTidDesc) ItemError {
//
//}
//
////>> 扣道具，成功返回nil
//func (this *ItemComponent) ReduceItemByUID(uid, count uint64, reason ItemChangeReason) ItemError {
//
//}
//
////>> 扣道具，成功返回nil
//func (this *ItemComponent) ReduceItemByTID(tid int32, count uint64, reason ItemChangeReason) ItemError {
//
//}
//
////>> 扣道具，成功返回nil
//func (this *ItemComponent) ReduceItems(items []ItemTidDesc, reason ItemChangeReason) ItemError {
//
//}
//
////>> 扣并给道具，保证事务性
//func (this *ItemComponent) ReduceAndAddItems(delItems, giveItems []ItemTidDesc, reason ItemChangeReason) ItemError {
//
//}
//
////>> 扣并给道具，保证事务性
//func (this *ItemComponent) ReduceAndAddItemByUID(delUIDs []ItemUidDesc, giveItems []ItemTidDesc, reason ItemChangeReason) ItemError {
//
//}
