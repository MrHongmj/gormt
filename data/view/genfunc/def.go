package genfunc

const (
	genTnf = `
// TableName get sql table name.获取数据库表名
func (m *{{.StructName}}) TableName() string {
	return "{{.TableName}}"
}
`
	genBase = `
package {{.PackageName}}
import (
	"context"

	"github.com/jinzhu/gorm"
)

var globalIsRelated bool = true  // 全局预加载

// prepare for other
type _BaseMgr struct {
	*model.BaseModel
	ctx       *context.Context
	isRelated bool
}

// SetCtx set context
func (obj *_BaseMgr) SetCtx(c *context.Context) {
	obj.ctx = c
}

// GetDB get gorm.DB info
func (obj *_BaseMgr) GetDB() *gorm.DB {
	return obj.DB
}

// UpdateDB update gorm.DB info
func (obj *_BaseMgr) UpdateDB(db *gorm.DB) {
	obj.DB = db
}

// GetIsRelated Query foreign key Association.获取是否查询外键关联(gorm.Related)
func (obj *_BaseMgr) GetIsRelated() bool {
	return obj.isRelated
}

// SetIsRelated Query foreign key Association.设置是否查询外键关联(gorm.Related)
func (obj *_BaseMgr) SetIsRelated(b bool) {
	obj.isRelated = b
}

type options struct {
	query map[string]interface{}
}

// Option overrides behavior of Connect.
type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}


// OpenRelated 打开全局预加载
func OpenRelated() {
	globalIsRelated = true
}

// CloseRelated 关闭全局预加载
func CloseRelated() {
	globalIsRelated = true
}

	`

	genlogic = `{{$obj := .}}{{$list := $obj.Em}}
type {{$obj.StructName}}Model struct {
	*model.BaseModel
}


func New{{$obj.StructName}}Model(db *gorm.DB) *{{$obj.StructName}}Model {
	if db == nil {
		panic(fmt.Errorf("{{$obj.StructName}} need init by db"))
	}
	return &{{$obj.StructName}}Model{BaseModel: &model.BaseModel{DB: db}}
}

// GetTableName get sql table name.获取数据库名字
func (obj *{{$obj.StructName}}Model) GetTableName() string {
	return "{{$obj.TableName}}"
}

// Get 获取
func (obj *{{$obj.StructName}}Model) Get(fields string) (result {{$obj.StructName}}, err error) {
	if	fields == "" {
		fields = "*"
	}
	err = obj.DB.Table(obj.GetTableName()).Select(fields).Find(&result).Error
	{{GenPreloadList $obj.PreloadList false}}
	return
}

// Gets 获取批量结果
func (obj *{{$obj.StructName}}Model) Gets(fields string) (results []*{{$obj.StructName}}, err error) {
	if	fields == "" {
		fields = "*"
	}
	err = obj.DB.Table(obj.GetTableName()).Select(fields).Find(&results).Error
	{{GenPreloadList $obj.PreloadList true}}
	return
}

 //////////////////////////primary index case ////////////////////////////////////////////
 {{range $ofm := $obj.Primay}}
 // {{GenFListIndex $ofm 1}} primay or index 获取唯一内容
 func (obj *{{$obj.StructName}}Model) {{GenFListIndex $ofm 1}}({{GenFListIndex $ofm 2}},fields string) (result {{$obj.StructName}}, err error) {
	if	fields == "" {
		fields = "*"
	}	
	err = obj.DB.Table(obj.GetTableName()).Select(fields).Where("{{GenFListIndex $ofm 3}}", {{GenFListIndex $ofm 4}}).Find(&result).Error
	{{GenPreloadList $obj.PreloadList false}}
	return
}
 {{end}}

 {{range $ofm := $obj.Index}}
 // {{GenFListIndex $ofm 1}}  获取多个内容
 func (obj *{{$obj.StructName}}Model) {{GenFListIndex $ofm 1}}({{GenFListIndex $ofm 2}},fields string) (results []*{{$obj.StructName}}, err error) {
	if	fields == "" {
		fields = "*"
	}	
	err = obj.DB.Table(obj.GetTableName()).Select(fields).Where("{{GenFListIndex $ofm 3}}", {{GenFListIndex $ofm 4}}).Find(&results).Error
	{{GenPreloadList $obj.PreloadList true}}
	return
}
 {{end}}

`
	genPreload = `if err == nil && obj.isRelated { {{range $obj := .}}{{if $obj.IsMulti}}
		err = obj.DB.New().Table("{{$obj.ForeignkeyTableName}}").Where("{{$obj.ForeignkeyCol}} = ?", result.{{$obj.ColStructName}}).Find(&result.{{$obj.ForeignkeyStructName}}List).Error // {{$obj.Notes}}
		{{else}} 
		err = obj.DB.New().Table("{{$obj.ForeignkeyTableName}}").Where("{{$obj.ForeignkeyCol}} = ?", result.{{$obj.ColStructName}}).Find(&result.{{$obj.ForeignkeyStructName}}).Error // {{$obj.Notes}} 
		{{end}} {{end}}}
`
	genPreloadMulti = `if err == nil && obj.isRelated {
		for i := 0; i < len(results); i++ { {{range $obj := .}}{{if $obj.IsMulti}}
		if err = obj.DB.New().Table("{{$obj.ForeignkeyTableName}}").Where("{{$obj.ForeignkeyCol}} = ?", results[i].{{$obj.ColStructName}}).Find(&results[i].{{$obj.ForeignkeyStructName}}List).Error;err != nil { // {{$obj.Notes}}
				return	
			} {{else}} 
		if err = obj.DB.New().Table("{{$obj.ForeignkeyTableName}}").Where("{{$obj.ForeignkeyCol}} = ?", results[i].{{$obj.ColStructName}}).Find(&results[i].{{$obj.ForeignkeyStructName}}).Error; err != nil { // {{$obj.Notes}} 
				return
			} {{end}} {{end}}
	}
}`
)
