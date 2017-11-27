package hook

import (
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
	"errors"
	"github.com/qor/admin"
	"github.com/jinzhu/inflection"
	"strings"
)

var ErrProcessorSkipLeft = errors.New("resource: skip left")

func (f *Hook) newFindManyHandler(res *admin.Resource) func(result interface{}, context *qor.Context) error {
	return func(result interface{}, context *qor.Context) error {
		//// FIXME model name == resource.name(Plural)
		tableName := inflection.Plural(strings.ToLower(res.Name))
		re := f.ResourceMap[res.Name]
		context.DB = context.DB.Table(tableName)
		context.Config.DB = context.DB.Table(tableName)
		return re.FindManyHandler(result, context)
	}
}

func (f *Hook) newSaveHandler(res *admin.Resource) func(result interface{}, context *qor.Context) error {
	return func(result interface{}, context *qor.Context) error {
		// FIXME model name == resource.name(Plural)
		tableName := inflection.Plural(strings.ToLower(res.Name))
		re := f.ResourceMap[res.Name]
		context.DB = context.DB.Table(tableName)
		context.Config.DB = context.DB.Table(tableName)
		return re.SaveHandler(result, context)
	}
}

func (f *Hook) newDeleteHandler(res *admin.Resource) func(result interface{}, context *qor.Context) error {
	return func(result interface{}, context *qor.Context) error {
		// FIXME model name == resource.name(Plural)
		tableName := inflection.Plural(strings.ToLower(res.Name))
		re := f.ResourceMap[res.Name]
		context.DB = context.DB.Table(tableName)
		context.Config.DB = context.DB.Table(tableName)
		return re.DeleteHandler(result, context)
	}
}

func (f *Hook) newFindOneHandler(res *admin.Resource) func(result interface{}, metaValues *resource.MetaValues, context *qor.Context) error {
	return func(result interface{}, metaValues *resource.MetaValues, context *qor.Context) error {
		// FIXME model name == resource.name(Plural)
		tableName := inflection.Plural(strings.ToLower(res.Name))

		//
		//if res.HasPermission(roles.Read, context) {
		//	var (
		//		primaryQuerySQL string
		//		primaryParams   []interface{}
		//	)
		//
		//	if metaValues == nil {
		//		primaryQuerySQL, primaryParams = res.ToPrimaryQueryParams(context.ResourceID, context)
		//	} else {
		//		primaryQuerySQL, primaryParams = res.ToPrimaryQueryParamsFromMetaValue(metaValues, context)
		//	}
		//
		//	if primaryQuerySQL != "" {
		//		if metaValues != nil {
		//			if destroy := metaValues.Get("_destroy"); destroy != nil {
		//				if fmt.Sprint(destroy.Value) != "0" && res.HasPermission(roles.Delete, context) {
		//					context.GetDB().Table(tableName).Where(append([]interface{}{primaryQuerySQL}, primaryParams...)).Delete(result)
		//					return ErrProcessorSkipLeft
		//				}
		//			}
		//		}
		//		return context.GetDB().Table(tableName).Where(append([]interface{}{primaryQuerySQL}, primaryParams...)).First(result).Error
		//	}
		//	return errors.New("failed to find")
		//}
		//return roles.ErrPermissionDenied
		re := f.ResourceMap[res.Name]
		context.DB = context.DB.Table(tableName)
		context.Config.DB = context.DB.Table(tableName)
		return re.FindOneHandler(result, metaValues, context)
	}
}

func (f *Hook) replaceResource(re *admin.Resource) {
	if f.ResourceMap[re.Name] == nil {
		fm := FuncMap{
			FindManyHandler: re.FindManyHandler,
			FindOneHandler:  re.FindOneHandler,
			SaveHandler:     re.SaveHandler,
			DeleteHandler:   re.DeleteHandler,
		}
		f.ResourceMap[re.Name] = &fm
		re.FindManyHandler = f.newFindManyHandler(re)
		re.FindOneHandler = f.newFindOneHandler(re)
		re.SaveHandler = f.newSaveHandler(re)
		re.DeleteHandler = f.newDeleteHandler(re)
	}

}
