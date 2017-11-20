package hook

import (
	"github.com/qor/qor"
	"github.com/qor/roles"
	"github.com/jinzhu/gorm"
	"github.com/qor/qor/resource"
	"fmt"
	"errors"
	"github.com/qor/admin"
)

var ErrProcessorSkipLeft = errors.New("resource: skip left")

func (f *Hook) findManyHandler(result interface{}, context *qor.Context) error {
	db := context.GetDB()
	tableName := getTableName(f.MTables, result)
	res := f.Admin.GetResource(tableName)
	if res.HasPermission(roles.Read, context) {
		if _, ok := db.Table(tableName).Get("qor:getting_total_count"); ok {
			return context.GetDB().Table(tableName).Count(result).Error
		}
		return context.GetDB().Table(tableName).Set("gorm:order_by_primary_key", "DESC").Find(result).Error
	}
	return roles.ErrPermissionDenied
}

func (f *Hook) saveHandler(result interface{}, context *qor.Context) error {
	tableName := getTableName(f.MTables, result)
	res := f.Admin.GetResource(tableName)
	if (context.GetDB().Table(tableName).NewScope(result).PrimaryKeyZero() &&
		res.HasPermission(roles.Create, context)) || // has create permission
		res.HasPermission(roles.Update, context) { // has update permission
		return context.GetDB().Table(tableName).Save(result).Error
	}
	return roles.ErrPermissionDenied
}

func (f *Hook) deleteHandler(result interface{}, context *qor.Context) error {
	tableName := getTableName(f.MTables, result)
	res := f.Admin.GetResource(tableName)
	if res.HasPermission(roles.Delete, context) {
		if !context.GetDB().Table(tableName).Where("id=?", context.ResourceID).First(result).RecordNotFound() {
			return context.GetDB().Table(tableName).Where("id=?", context.ResourceID).Delete(result).Error
		}

		return gorm.ErrRecordNotFound
	}
	return roles.ErrPermissionDenied
}

func (f *Hook) findOneHandler(result interface{}, metaValues *resource.MetaValues, context *qor.Context) error {
	tableName := getTableName(f.MTables, result)
	res := f.Admin.GetResource(tableName)
	if res.HasPermission(roles.Read, context) {
		var (
			primaryQuerySQL string
			primaryParams   []interface{}
		)

		if metaValues == nil {
			primaryQuerySQL, primaryParams = res.ToPrimaryQueryParams(context.ResourceID, context)
		} else {
			primaryQuerySQL, primaryParams = res.ToPrimaryQueryParamsFromMetaValue(metaValues, context)
		}

		if primaryQuerySQL != "" {
			if metaValues != nil {
				if destroy := metaValues.Get("_destroy"); destroy != nil {
					if fmt.Sprint(destroy.Value) != "0" && res.HasPermission(roles.Delete, context) {
						context.GetDB().Table(tableName).Where(append([]interface{}{primaryQuerySQL}, primaryParams...)).Delete(result)
						return ErrProcessorSkipLeft
					}
				}
			}
			return context.GetDB().Table(tableName).Where(append([]interface{}{primaryQuerySQL}, primaryParams...)).First(result).Error
		}
		return errors.New("failed to find")
	}
	return roles.ErrPermissionDenied
}

func (f *Hook) replaceResource(re *admin.Resource) {
	re.FindManyHandler = f.findManyHandler
	re.FindOneHandler = f.findOneHandler
	re.SaveHandler = f.saveHandler
	re.DeleteHandler = f.deleteHandler
}
