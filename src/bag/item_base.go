package bag

//>> 道具接口
type ItemInterface interface {
	//>> 模板ID
	GetTID() int32
	//>> 道具类型
	GetType() int32
	//>> UID
	GetUID() uint64
	//>> 数量
	GetCount() uint64
	//>> 创建时间
	GetCreateTime() int64
	//>> 在容器中所处位置
	GetPos() int16
	//>> 所处容器类型
	GetContainerType() ContainerType
	//>> 一些标记，比如是否是绑定的
	GetWeight() int32
	GetFlag() int

	SetTID(int32)
	SetUID(uint64)
	SetCount(uint64)
	SetCreateTime(int64)
	SetPos(int16)
	SetContainerType(ContainerType)
	SetFlag(int)
}

//>> 道具bit标记
const (
	IsBind = 1 << 0
)

//>> 道具基础信息
type ItemBase struct {
	tid        int32
	uid        uint64
	count      uint64
	createTime int64

	pos          int16
	containerTyp int16
	flag         int
}

func (this *ItemBase) GetTID() int32 {
	panic("implement me")
}

func (this *ItemBase) GetUID() uint64 {
	panic("implement me")
}

func (this *ItemBase) GetCount() uint64 {
	panic("implement me")
}

func (this *ItemBase) GetCreateTime() int64 {
	panic("implement me")
}

func (this *ItemBase) GetType() int32 {
	return getItemType(this.tid)
}

func (this *ItemBase) GetPos() int16 {
	panic("implement me")
}

func (this *ItemBase) GetContainerType() ContainerType {
	panic("implement me")
}

func (this* ItemBase) GetWeight() int32{
	return getItemWeight(this.tid)
}

func (this *ItemBase) GetFlag() int {
	panic("implement me")
}

func (this *ItemBase) SetTID(int32) {
	panic("implement me")
}

func (this *ItemBase) SetUID(uint64) {
	panic("implement me")
}

func (this *ItemBase) SetCount(uint64) {
	panic("implement me")
}

func (this *ItemBase) SetCreateTime(int64) {
	panic("implement me")
}

func (this *ItemBase) SetType(int32) {
	panic("implement me")
}

func (this *ItemBase) SetPos(int16) {
	panic("implement me")
}

func (this *ItemBase) SetContainerType(ContainerType) {
	panic("implement me")
}

func (this *ItemBase) SetFlag(int) {
	panic("implement me")
}


