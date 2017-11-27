package utils

import (
	"github.com/qor/admin"
	//"unsafe"
	//"fmt"
)

type HookTableAdmin struct {
	*admin.Admin
}

func AddResource(qor_admin *admin.Admin, value interface{}, tbName string) {
	//configuration := &admin.Config{}

	//var res *admin.Resource = new(admin.Resource)
	//res.Config = configuration
	//var adm *admin.Admin = (*admin.Admin)(unsafe.Pointer(res))
	//adm = qor_admin

	//res = &admin.Resource{Value: value, Name: name}
	//res.Permission = configuration.Permission

}
