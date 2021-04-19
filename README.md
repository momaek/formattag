## Formattag

The tool is used to align golang struct's tags.

eg.:

Before
```golang
// TestStruct this is a test struct
type TestStruct struct {
	ID            string `json:"id" xml:"id"`
	IfNotModified string `json:"if_not_modified" xml:"if_not_modified"`
	Name          string `json:"name" xml:"name"`

	ThisIsAStructWodeTianNa struct {
		FieldFromThisIsAStructWodeTianNa string `json:"field_from_this_is_a_struct_wode_tian_na" xml:"field_from_this_is_a_struct_wode_tian_na"`
		TianName                         string `json:"tian_name" xml:"tian_name"`
	} `json:"this_is_a_struct_wode_tian_na" xml:"this_is_a_struct_wode_tian_na"`

	T    time.Time `json:"t" xml:"t"`
	Fset Fset      `json:"fset" xml:"fset"`
}

type Fset struct{}

```

After
```golang
// TestStruct this is a test struct
type TestStruct struct {
	ID            string `json:"id"              xml:"id"`
	IfNotModified string `json:"if_not_modified" xml:"if_not_modified"`
	Name          string `json:"name"            xml:"name"`

	ThisIsAStructWodeTianNa struct {
		FieldFromThisIsAStructWodeTianNa string `json:"field_from_this_is_a_struct_wode_tian_na" xml:"field_from_this_is_a_struct_wode_tian_na"`
		TianName                         string `json:"tian_name"                                xml:"tian_name"`
	} `json:"this_is_a_struct_wode_tian_na" xml:"this_is_a_struct_wode_tian_na"`

	T    time.Time `json:"t"    xml:"t"`
	Fset Fset      `json:"fset" xml:"fset"`
}

type Fset struct{}
```

### Installation
Compile from source, which requires [Go 1.15 or newer](https://golang.org/doc/install):
```
go get github.com/momaek/formattag
```

### Usage 

```
formattag -file /path/to/your/golang/file
```

This command will change your go file.

#### Vim
Add this to your ~/.vimrc:
```
set rtp+=$GOPATH/src/github.com/momaek/formattag/vim
```
If you have multiple entries in your GOPATH, replace $GOPATH with the right value.

Running `:PrettyTag` will run formattag on the current file.

Optionally, add this to your `~/.vimrc` to automatically run `formattag` on :w
```
autocmd BufWritePost,FileWritePost *.go execute 'PrettyTag' | checktime
```

#### VSCode
Please Install [Run On Save](https://marketplace.visualstudio.com/items?itemName=emeraldwalk.RunOnSave)
Add this to `settings.json` commands:
```json
{
    "match": "\\.go$",
    "isAsync": false,
    "cmd": "/path/to/formattag -file ${file}"
}
```

