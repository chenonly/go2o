package promodel

import (
	"database/sql"
	"go2o/core/domain/interface/pro_model"
)

var _ promodel.IAttrService = new(attrServiceImpl)

type attrServiceImpl struct {
	rep promodel.IProModelRepo
}

func NewAttrService(rep promodel.IProModelRepo) *attrServiceImpl {
	return &attrServiceImpl{
		rep: rep,
	}
}

// 获取属性
func (a *attrServiceImpl) GetAttr(attrId int32) *promodel.Attr {
	return a.rep.GetAttr(attrId)
}

// 保存属性
func (a *attrServiceImpl) SaveAttr(v *promodel.Attr) (id int32, err error) {
	var i int
	// 如不存在，则新增
	if v.Id <= 0 {
		i, err = a.rep.SaveAttr(v)
		v.Id = int32(i)
		if v == nil {
			return v.Id, err
		}
	}
	// 保存项
	if v.Items != nil {
		v.ItemValues = ""
		for i, iv := range v.Items {
			iv.ProModel = v.ProModel
			iv.AttrId = v.Id
			if i > 0 {
				v.ItemValues += ","
			}
			v.ItemValues += iv.Value
		}
		err = a.saveAttrItems(v.Id, v.Items)
	}
	// 再次保存
	if err == nil {
		_, err = a.rep.SaveAttr(v)
	}
	return v.Id, err
}

// 保存属性项
func (a *attrServiceImpl) saveAttrItems(attrId int32, items []*promodel.AttrItem) (err error) {
	var i int
	pk := attrId
	// 获取存在的项
	old := a.rep.SelectAttrItem("attr_id = ?", pk)
	// 分析当前项目并加入到MAP中
	delList := []int32{}
	currMap := make(map[int32]*promodel.AttrItem, len(items))
	for _, v := range items {
		currMap[v.Id] = v
	}
	// 筛选出要删除的项
	for _, v := range old {
		if currMap[v.Id] == nil {
			delList = append(delList, v.Id)
		}
	}
	// 删除项
	for _, v := range delList {
		a.rep.DeleteAttrItem(v)
	}
	// 保存项
	for _, v := range items {
		if v.AttrId == 0 {
			v.AttrId = pk
		}
		if v.AttrId == pk {
			if i, err = a.rep.SaveAttrItem(v); err == nil {
				v.Id = int32(i)
			}
		}
	}
	return err
}

// 保存属性项
func (a *attrServiceImpl) SaveItem(v *promodel.AttrItem) (int32, error) {
	id, err := a.rep.SaveAttrItem(v)
	return int32(id), err
}

// 删除属性
func (a *attrServiceImpl) DeleteAttr(attrId int32) error {
	_, err := a.rep.BatchDeleteAttrItem("attr_id=?", attrId)
	if err == nil || err == sql.ErrNoRows {
		err = a.rep.DeleteAttr(attrId)
	}
	return err
}

// 删除属性项
func (a *attrServiceImpl) DeleteItem(itemId int32) error {
	return a.rep.DeleteAttrItem(itemId)
}

// 获取属性的属性项
func (a *attrServiceImpl) GetItems(attrId int32) []*promodel.AttrItem {
	return a.rep.SelectAttrItem("attr_id=?", attrId)
}

// 获取产品模型的属性
func (a *attrServiceImpl) GetModelAttrs(proModel int32) []*promodel.Attr {
	arr := a.rep.SelectAttr("pro_model=?", proModel)
	for _, v := range arr {
		v.Items = a.GetItems(v.Id)
	}
	return arr
}

// 获取产品的属性
//func (a *attrServiceImpl)GetGoodsAttrs(proId int32) []*ProAttr