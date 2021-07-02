package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"github.com/fatih/structtag"
)

// structType contains a structType node and it's name. It's a convenient
// helper type, because *ast.StructType doesn't contain the name of the struct
type structType struct {
	name string
	node *ast.StructType
}
type Config struct {
	templateFile         string
	file                 string
	structName           string
	fset                 *token.FileSet
	offset               int
	line                 string
	start, end           int
	all                  bool
	skipUnexportedFields bool
	// max number of colums  we use to generate GetBy , DeleteBy, UpdateBy code
	maxColumns int8
	src        []byte
}

func parseConfig(args []string) (*Config, error) {
	var (
		flagFile     = flag.String("file", "", "Filename to be parsed")
		templateFile = flag.String("template", "gorm.template", "template to use generate code")
		flagStruct   = flag.String("struct", "", "Struct name to be processed")
	)

	if err := flag.CommandLine.Parse(args); err != nil {
		return nil, err
	}

	if flag.NFlag() == 0 {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return nil, flag.ErrHelp
	}
	cfg := &Config{
		file:         *flagFile,
		templateFile: *templateFile,
		structName:   *flagStruct,
	}

	return cfg, nil
}

func (c *Config) parse() (ast.Node, error) {
	c.fset = token.NewFileSet()
	var contents interface{}
	var err error
	c.src, err = ioutil.ReadFile(c.file)
	if err != nil {
		return &ast.File{}, err
	}

	return parser.ParseFile(c.fset, c.file, contents, parser.ParseComments)
}

// findSelection returns the start and end position of the fields that are
// suspect to change. It depends on the line, struct or offset selection.
func (c *Config) findSelection(node ast.Node) (int, int, error) {
	if c.line != "" {
		return c.lineSelection(node)
	} else if c.offset != 0 {
		return c.offsetSelection(node)
	} else if c.structName != "" {
		return c.structSelection(node)
	} else if c.all {
		return c.allSelection(node)
	} else {
		return 0, 0, errors.New("-line, -offset, -struct or -all is not passed")
	}
}
func (c *Config) lineSelection(file ast.Node) (int, int, error) {
	var err error
	splitted := strings.Split(c.line, ",")

	start, err := strconv.Atoi(splitted[0])
	if err != nil {
		return 0, 0, err
	}

	end := start
	if len(splitted) == 2 {
		end, err = strconv.Atoi(splitted[1])
		if err != nil {
			return 0, 0, err
		}
	}

	if start > end {
		return 0, 0, errors.New("wrong range. start line cannot be larger than end line")
	}

	return start, end, nil
}

func (c *Config) structSelection(file ast.Node) (int, int, error) {
	structs := collectStructs(file)

	var encStruct *ast.StructType
	for _, st := range structs {
		if st.name == c.structName {
			encStruct = st.node
		}
	}

	if encStruct == nil {
		return 0, 0, errors.New("struct name does not exist")
	}

	start := c.fset.Position(encStruct.Pos()).Line
	end := c.fset.Position(encStruct.End()).Line

	return start, end, nil
}

func (c *Config) offsetSelection(file ast.Node) (int, int, error) {
	structs := collectStructs(file)

	var encStruct *ast.StructType
	for _, st := range structs {
		structBegin := c.fset.Position(st.node.Pos()).Offset
		structEnd := c.fset.Position(st.node.End()).Offset

		if structBegin <= c.offset && c.offset <= structEnd {
			encStruct = st.node
			break
		}
	}

	if encStruct == nil {
		return 0, 0, errors.New("offset is not inside a struct")
	}

	// offset selects all fields
	start := c.fset.Position(encStruct.Pos()).Line
	end := c.fset.Position(encStruct.End()).Line

	return start, end, nil
}

// allSelection selects all structs inside a file
func (c *Config) allSelection(file ast.Node) (int, int, error) {
	start := 1
	end := c.fset.File(file.Pos()).LineCount()

	return start, end, nil
}

func isPublicName(name string) bool {
	for _, c := range name {
		return unicode.IsUpper(c)
	}
	return false
}

// collectStructs collects and maps structType nodes to their positions
func collectStructs(node ast.Node) map[token.Pos]*structType {
	structs := make(map[token.Pos]*structType, 0)

	collectStructs := func(n ast.Node) bool {
		var t ast.Expr
		var structName string

		switch x := n.(type) {
		case *ast.TypeSpec:
			if x.Type == nil {
				return true

			}

			structName = x.Name.Name
			t = x.Type
		case *ast.CompositeLit:
			t = x.Type
		case *ast.ValueSpec:
			structName = x.Names[0].Name
			t = x.Type
		case *ast.Field:
			// this case also catches struct fields and the structName
			// therefore might contain the field name (which is wrong)
			// because `x.Type` in this case is not a *ast.StructType.
			//
			// We're OK with it, because, in our case *ast.Field represents
			// a parameter declaration, i.e:
			//
			//   func test(arg struct {
			//   	Field int
			//   }) {
			//   }
			//
			// and hence the struct name will be `arg`.
			if len(x.Names) != 0 {
				structName = x.Names[0].Name
			}
			t = x.Type
		}

		// if expression is in form "*T" or "[]T", dereference to check if "T"
		// contains a struct expression
		t = deref(t)

		x, ok := t.(*ast.StructType)
		if !ok {
			return true
		}

		structs[x.Pos()] = &structType{
			name: structName,
			node: x,
		}
		return true
	}

	ast.Inspect(node, collectStructs)
	return structs
}

// deref takes an expression, and removes all its leading "*" and "[]"
// operator. Uuse case : if found expression is a "*t" or "[]t", we need to
// check if "t" contains a struct expression.
func deref(x ast.Expr) ast.Expr {
	switch t := x.(type) {
	case *ast.StarExpr:
		return deref(t.X)
	case *ast.ArrayType:
		return deref(t.Elt)
	}
	return x
}

// collectStructModels rewrites the node for structs between the start and end
// positions
func (c *Config) collectStructModels(node ast.Node, start, end int) (models []ModelType, packageName string, err error) {
	errs := &rewriteErrors{errs: make([]error, 0)}

	rewriteFunc := func(n ast.Node) bool {

		var t ast.Expr
		var structName string

		switch x := n.(type) {

		case *ast.File:
			packageName = x.Name.Name
		case *ast.TypeSpec:
			if x.Type == nil {
				return true

			}

			structName = x.Name.Name
			t = x.Type
		case *ast.CompositeLit:
			t = x.Type
		case *ast.ValueSpec:
			structName = x.Names[0].Name
			t = x.Type
		case *ast.Field:
			// this case also catches struct fields and the structName
			// therefore might contain the field name (which is wrong)
			// because `x.Type` in this case is not a *ast.StructType.
			//
			// We're OK with it, because, in our case *ast.Field represents
			// a parameter declaration, i.e:
			//
			//   func test(arg struct {
			//   	Field int
			//   }) {
			//   }
			//
			// and hence the struct name will be `arg`.
			if len(x.Names) != 0 {
				structName = x.Names[0].Name
			}
			t = x.Type
		}

		// if expression is in form "*T" or "[]T", dereference to check if "T"
		// contains a struct expression
		t = deref(t)

		x, ok := t.(*ast.StructType)
		if !ok {
			return true
		}

		if c.structName != "" && structName != c.structName {
			return true
		}

		tmpModel := ModelType{
			ModelName:   structName,
			PackageName: packageName,
		}

		for _, f := range x.Fields.List {
			line := c.fset.Position(f.Pos()).Line

			if !(start <= line && line <= end) {
				continue
			}

			fieldName := ""
			if len(f.Names) != 0 {
				for _, field := range f.Names {
					if !c.skipUnexportedFields || isPublicName(field.Name) {
						fieldName = field.Name
						break
					}
				}
			}

			// anonymous field
			if f.Names == nil {
				ident, ok := f.Type.(*ast.Ident)
				if !ok {
					continue
				}

				if !c.skipUnexportedFields {
					fieldName = ident.Name
				}
			}

			// nothing to process, continue with next line
			if fieldName == "" {
				continue
			}

			if f.Tag == nil {
				continue
			}

			column, err := c.process(fieldName, f.Tag.Value)
			if err != nil {
				errs.Append(fmt.Errorf("%s:%d:%d:%s",
					c.fset.Position(f.Pos()).Filename,
					c.fset.Position(f.Pos()).Line,
					c.fset.Position(f.Pos()).Column,
					err))
				continue
			}

			start := f.Type.Pos() - 1
			end := f.Type.End() - 1

			// grab it in source
			column.GoType = string(c.src[start:end])

			if isSupportType(column.GoType) {
				column.VarName, _ = convertFieldName("camelcase", column.FieldName)
				tmpModel.Columns = append(tmpModel.Columns, column)
			}

		}

		models = append(models, tmpModel)

		return true
	}

	ast.Inspect(node, rewriteFunc)

	c.start = start
	c.end = end

	return
}

func (c *Config) process(fieldName, tagVal string) (ColumnType, error) {
	var tag string
	if tagVal != "" {
		var err error
		tag, err = strconv.Unquote(tagVal)
		if err != nil {
			return ColumnType{}, err
		}
	}

	tags, err := structtag.Parse(tag)
	if err != nil {
		return ColumnType{}, err
	}

	for _, tag := range tags.Tags() {
		if tag.Key == "gorm" {
			return extractColumn(fieldName, tag.Name), nil
		}
	}

	return ColumnType{}, nil

}

type rewriteErrors struct {
	errs []error
}

func (r *rewriteErrors) Error() string {
	var buf bytes.Buffer
	for _, e := range r.errs {
		buf.WriteString(fmt.Sprintf("%s\n", e.Error()))
	}
	return buf.String()
}

func (r *rewriteErrors) Append(err error) {
	if err == nil {
		return
	}

	r.errs = append(r.errs, err)
}

func isSupportType(typeName string) bool {
	switch typeName {
	case "string", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return true
	}

	return false

}

func (cfg *Config) GenerateGormCode() error {

	node, err := cfg.parse()
	if err != nil {
		return err
	}

	start, end, err := cfg.findSelection(node)
	if err != nil {
		return err
	}

	models, _, errs := cfg.collectStructModels(node, start, end)
	if errs != nil {
		if _, ok := errs.(*rewriteErrors); !ok {
			return errs
		}
	}

	extractUniqueIndex(models)

	// fmt.Printf("%+v", models)

	fileDir := filepath.Dir(cfg.file)

	return generateGormCode(models, fileDir, cfg.templateFile)
}
