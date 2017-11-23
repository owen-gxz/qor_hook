package hook

import (
	"github.com/qor/admin"
	"github.com/jinzhu/gorm"
	"fmt"
	"errors"
	"github.com/qor/hook/utils"
	"encoding/json"

	"reflect"
	"strings"
	"github.com/qor/roles"
)

type Hook struct {
	Admin   *admin.Admin
	MTables map[string]interface{}
}

func (Hook) ResourceName() string {
	return "Hook"
}

func New(Admin *admin.Admin) *Hook {
	Admin.DB.AutoMigrate(&ResourceModel{}, &ResourceTableModel{})

	fb := &Hook{Admin: Admin}
	addHook(fb)
	flexs, err := loadHook(Admin.DB)
	if err != nil {
		panic(err)
		return nil
	}
	//表对应新增的字段
	tm := make(map[string]map[string]string)
	for _, r := range flexs {
		nv := make(map[string]string, 0)
		if tm[r.Name] == nil {
			nv[r.FieldName] = r.Type
			tm[r.Name] = nv
		} else {
			v := tm[r.Name]
			v[r.FieldName] = r.Type
		}
	}
	//设置新的values
	for k, v := range tm {
		fb.resourceLoadNew(k, v)
	}
	return fb
}

func (f *Hook) resourceLoadNew(k string, v map[string]string) {
	ar := f.Admin.GetResource(k)
	utils.SetNewValue(ar, v)
	vv := utils.GetSlices(ar)
	if f.MTables == nil {
		mm := make(map[string]interface{}, 0)
		mm[k] = vv
		f.MTables = mm
	} else {
		f.MTables[k] = vv
	}
	//替换原有crud
	f.replaceResource(ar)
	newFiled := make([]string, 0)
	for k, _ := range v {
		newFiled = append(newFiled, utils.Upper(k))
	}
	ar.OverrideEditAttrs(func() {
		ar.EditAttrs(ar.EditAttrs(), newFiled)
	})
	ar.OverrideIndexAttrs(func() {
		ar.IndexAttrs(ar.IndexAttrs(), newFiled)
	})
	ar.OverrideNewAttrs(func() {
		ar.NewAttrs(ar.NewAttrs(), newFiled)
	})
	ar.OverrideShowAttrs(func() {
		ar.ShowAttrs(ar.ShowAttrs(), newFiled)
	})
}

func getTableName(m map[string]interface{}, result interface{}) string {
	vv := reflect.TypeOf(result).Elem()
	for k, v := range m {
		////todo 两个相同的结构体会有问题
		if strings.Contains(fmt.Sprintf("%#v", v), vv.String()) {
			return k
		}
	}

	return ""
}

//增加Flexible resource
func addHook(f *Hook, config ...*admin.Config) {
	f.Admin.AddMenu(&admin.Menu{Name: "Hook"})
	hookTypeResource := f.Admin.AddResource(&ResourceModel{}, &admin.Config{
		Menu:       []string{"Hook"},
		Permission: roles.Allow(roles.Read, roles.Anyone).Allow(roles.Create, roles.Anyone).
			Allow(roles.Delete, roles.Anyone),
	})
	setType(hookTypeResource, f)
	hookTableResource := f.Admin.AddResource(&ResourceTableModel{}, &admin.Config{Menu: []string{"Hook"}})
	AddTable(hookTableResource, f)
}

func loadHook(db *gorm.DB) ([]ResourceModel, error) {
	list := make([]ResourceModel, 0)
	err := db.Model(&ResourceModel{}).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func getTables(db *gorm.DB) []string {

	dbType := db.Dialect().GetName()
	tables := make([]string, 0)
	var showTable string
	if dbType == Sqlite {
		showTable = "SELECT name FROM sqlite_master WHERE type='table' order by name"
	} else if dbType == Mysql {
		showTable = "show tables"
	}
	rows, err := db.Raw(showTable).Rows() // (*sql.Rows, error)
	if err != nil {
		panic(err)
		return nil
	}
	for rows.Next() {
		var name string
		rows.Scan(&name)

		tables = append(tables, name)
	}
	return tables
}

func addMeta(rm ResourceModel, Admin *admin.Admin) error {

	rs := Admin.GetResource(rm.Name)
	if rs == nil {
		return errors.New("resource error")
	}
	m := &admin.Meta{
		Name: utils.Upper(rm.FieldName),
		Type: rm.Type,
	}

	//更新value

	switch rm.Type {
	case "select_one":
		strs := make([]string, 0)
		err := json.Unmarshal([]byte(rm.FieldSelects), &strs)
		if err != nil {
			return err
		}
		m.Config = &admin.SelectOneConfig{Collection: strs}
	case "File":
		m.Type = "file"
	}
	rs.Meta(m)
	return nil
}

func delMeta(rm ResourceModel, Admin *admin.Admin) error {

	rs := Admin.GetResource(rm.Name)
	if rs == nil {
		return errors.New("resource error")
	}
	m := &admin.Meta{
		Name: utils.Upper(rm.FieldName),
		Type: "hidden",
	}
	rs.Meta(m)
	return nil
}
