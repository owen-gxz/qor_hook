```
git clone git@github.com:owengaozhen/qor-hook.git hook
mv -f hook/ $GOPATH/src/github.com/qor/
```

qor-hook no restart,table add column--only mysql


## Demo
```
    Admin := admin.New(&admin.AdminConfig{
        DB: DB,
    })
	hook.New(Admin)
```
