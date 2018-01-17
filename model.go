package hook

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/qor"
	"github.com/qor/admin"
	"errors"
	"github.com/qor/roles"
	"fmt"
	"github.com/owengaozhen/qor-hook/utils"
)

//数据库默认使用mysql
type ResourceModel struct {
	gorm.Model

	Name         string `json:"name"`                   //表名称
	FieldName    string `json:"field_name"`             //字段名称
	FieldType    string `json:"field_size"`             //字段类型
	FieldSize    int    `json:"field_size"`             //字段大小
	FieldSelects string `json:"field_selects"`          //select 或多选按钮时使用 [{k:"",v:""}]
	Showed       bool   `json:"showed" sql:"default:1"` //是否显示
	Type         string `json:"type"`                   //展示的类型
}

func (ResourceModel) TableName() string {
	return "resource_models"
}

func setType(hookResource *admin.Resource, f *Hook) {
	hookResource.Meta(&admin.Meta{
		Name:   "Type",
		Label:  "模板类型",
		Type:   "select_one",
		Config: &admin.SelectOneConfig{Collection: TypeList},
	})
	hookResource.Meta(&admin.Meta{
		Name:   "FieldType",
		Label:  "字段类型",
		Type:   "select_one",
		Config: &admin.SelectOneConfig{Collection: SqlType},
	})
	hookResource.Meta(&admin.Meta{
		Name:  "FieldSize",
		Label: "字段长度",
		Type:  "number",
	})
	hookResource.Meta(&admin.Meta{
		Name:  "Showed",
		Type:  "checkbox",
		Label: "是否显示",
		Valuer: func(record interface{}, context *qor.Context) interface{} {
			return true
		},
	})
	hookResource.Meta(&admin.Meta{
		Name:  "FieldName",
		Label: "字段名称",
	})
	hookResource.Meta(&admin.Meta{
		Name:  "FieldSelects",
		Type:  "string",
		Label: "可选项",
		Valuer: func(record interface{}, context *qor.Context) interface{} {
			return `[""]`
		},
	})
	tables := getTables(f.Admin.DB)
	hookResource.Meta(&admin.Meta{
		Name:   "Name",
		Type:   "select_one",
		Label:  "数据库名称",
		Config: &admin.SelectOneConfig{Collection: tables},
	})
	hookResource.SaveHandler = f.saveHookHandler
	hookResource.DeleteHandler = f.deleteHookHandler
}

func (f *Hook) deleteHookHandler(result interface{}, context *qor.Context) error {
	rm, ok := result.(*ResourceModel)
	if !ok {
		return errors.New("model error")

	}
	tableName := ResourceModel{}.TableName()
	res := f.Admin.GetResource(tableName)
	//
	db := context.GetDB().Begin()
	if res.HasPermission(roles.Delete, context) {
		if !db.Table(tableName).Where("id=?", context.ResourceID).First(result).RecordNotFound() {
			delRes := f.Admin.GetResource(rm.Name)
			if delRes == nil {
				db.Rollback()
				return errors.New("get resource nil")
			}
			err := db.Table(tableName).Where("id=?", context.ResourceID).Delete(result).Error
			if err != nil {
				db.Rollback()
				return err
			}
			//err = db.Exec("ALTER TABLE ? DROP ?", rm.Name, rm.FieldName).Error
			err = db.Exec(fmt.Sprintf("ALTER TABLE %s DROP %s", rm.Name, rm.FieldName)).Error
			if err != nil {
				db.Rollback()
				return err
			}

			//删除更改对应的value
			utils.DelValueKey(delRes, rm.FieldName)
			//更改对应的mt
			//f.MTables[rm.Name] = utils.GetSlices(delRes)
			//覆盖字段
			delRes.OverrideEditAttrs(func() {
				delRes.EditAttrs(getDelFieldAttrs(delRes.EditAttrs(), rm.FieldName))
			})
			delRes.OverrideIndexAttrs(func() {
				delRes.IndexAttrs(getDelFieldAttrs(delRes.IndexAttrs(), rm.FieldName))
			})
			delRes.OverrideNewAttrs(func() {
				delRes.NewAttrs(getDelFieldAttrs(delRes.NewAttrs(), rm.FieldName))
			})
			delRes.OverrideShowAttrs(func() {
				delRes.ShowAttrs(getDelFieldAttrs(delRes.ShowAttrs(), rm.FieldName))
			})
			db.Commit()
			return nil
		}
		db.Rollback()
		return gorm.ErrRecordNotFound
	}
	db.Rollback()
	return roles.ErrPermissionDenied
}

func getDelFieldAttrs(rs []*admin.Section, delK string) []*admin.Section {
	n := make([]*admin.Section, 0)
	for _, item := range rs {
		if len(item.Rows) < 1 || utils.Upper(item.Rows[0][0]) == utils.Upper(delK) {
			continue
		}

		n = append(n, item)
		//rs.
	}
	return n
}

func (f *Hook) saveHookHandler(i interface{}, context *qor.Context) error {
	rm, ok := i.(*ResourceModel)
	if !ok {
		return errors.New("model error")

	}
	tx := context.GetDB().Begin()
	if rm.ID == 0 {
		err := tx.Where("name=? and field_name=? ", rm.Name, rm.FieldName).FirstOrCreate(rm).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		return nil
	}
	var v string
	switch rm.FieldType {
	case "DATETIME":
		fallthrough
	case "INT":
		v = rm.FieldType
	case "VARCHAR":
		v = fmt.Sprintf("%s(%d)", rm.FieldType, rm.FieldSize)
	}
	err := tx.Exec(fmt.Sprintf("ALTER TABLE %s ADD %s %s NULL", rm.Name, rm.FieldName, v)).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	//todo 重新加载value
	vm := make(map[string]string, 0)
	vm[rm.FieldName] = rm.Type
	f.resourceLoadNew(rm.Name, vm)
	err = addMeta(*rm, f.Admin)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

//  todo 加表暂时不行
type ResourceTableModel struct {
	gorm.Model

	Name string `json:"name"` //表名称
}

func (ResourceTableModel) TableName() string {
	return "resource_table_models"
}

func AddTable(hookTable *admin.Resource, f *Hook) {
	hookTable.SaveHandler = f.saveHookTableHandler
}

func (f *Hook) saveHookTableHandler(i interface{}, context *qor.Context) error {
	rm, ok := i.(*ResourceTableModel)
	if !ok {
		return errors.New("model error")
	}
	tx := context.GetDB().Begin()
	if rm.ID == 0 {
		err := tx.Exec(`CREATE TABLE ? (
			'id' int(10) unsigned NOT NULL AUTO_INCREMENT,
			'created_at' timestamp NULL DEFAULT NULL,
			'updated_at' timestamp NULL DEFAULT NULL,
			'deleted_at' timestamp NULL DEFAULT NULL,
			PRIMARY KEY ('id'),
			KEY 'idx_`+ rm.Name+ `_deleted_at' ('deleted_at')
		)`, rm.Name).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		tx.Commit()
		return nil
	}
	//todo 重写addResource方法



	tx.Commit()
	return nil
}
