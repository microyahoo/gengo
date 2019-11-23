// Package generators has the generators for the set-gen utility.
package generators

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	// "k8s.io/klog"
)

const (
	jsonFileType = "json"
	errCodeTag   = "ErrCode"
	beginTag     = "Begin"
	endTag       = "End"
)

var _ = os.Exit

var typesMap map[types.Name]constTypeInfo = make(map[types.Name]constTypeInfo)
var errCodeMap = make(map[string]string)

type constTypeInfo struct {
	tag  string
	name string
	val  interface{}
	t    *types.Type
}

// NameSystems returns the name system used by the generators in this package.
func NameSystems() namer.NameSystems {
	return namer.NameSystems{
		"public":  namer.NewPublicNamer(0),
		"private": namer.NewPrivateNamer(0),
		"raw":     namer.NewRawNamer("", nil),
	}
}

// DefaultNameSystem returns the default name system for ordering the types to be
// processed by the generators in this package.
func DefaultNameSystem() string {
	return "public"
}

// Packages makes the sets package definition.
func Packages(context *generator.Context, arguments *args.GeneratorArgs) generator.Packages {
	// boilerplate, err := arguments.LoadGoBoilerplate()
	// if err != nil {
	// 	klog.Fatalf("Failed loading boilerplate: %v", err)
	// }
	context.FileTypes = map[string]generator.FileType{
		jsonFileType: jsonFile{
			c:        context,
			Format:   jsonFormat,
			Assemble: jsonAssemble,
		},
	}

	return generator.Packages{&generator.DefaultPackage{
		PackageName: "codes",
		PackagePath: arguments.OutputPackagePath,
		// HeaderText:  boilerplate,
		// PackageDocumentation: []byte(
		// 	`// Package codes has auto-generated set types.
		// `),

		// GeneratorFunc returns a list of generators. Each generator makes a
		// single file.
		GeneratorFunc: func(c *generator.Context) (generators []generator.Generator) {
			// fmt.Printf("***generator.Context = %+v\n", c)

			// Since we want a file per type that we generate a set for, we
			// have to provide a function for this.
			generators = append(generators, &errCodeGen{
				DefaultGen:    generator.DefaultGen{},
				outputPackage: arguments.OutputPackagePath,
				imports:       generator.NewImportTracker(),
			})
			return generators
		},
		FilterFunc: func(c *generator.Context, t *types.Type) bool {
			// It would be reasonable to filter by the type's package here.
			// It might be necessary if your input directory has a big
			// import graph.
			// fmt.Printf("**t.Name = %+v, t.Kind = %+v, t.CommentLines = %v, t.SecondClosestCommentLines = %v\n",
			// t.Name, t.Kind, t.CommentLines, t.SecondClosestCommentLines)
			switch t.Kind {
			// we only care about const type
			case types.DeclarationOf:
				// Only some structs can be keys in a map. This is triggered by the line
				// // +ErrCode=Rbd,Common
				tags := extractTags(errCodeTag, t.CommentLines)
				// fmt.Printf("\t**name = %+v, **%+v, tags = %+v, len(tags)=%v, t.ConstValue = %#v, type(t.ConstValue) = %T\n", t.Name.Name, t.Name, tags, len(tags), t.ConstValue, t.ConstValue)
				if len(tags) == 0 {
					return false
				}
				typesMap[t.Name] = constTypeInfo{
					// there is only one value of one tag
					tag:  tags[0],
					name: t.Name.Name,
					val:  t.ConstValue,
					t:    t,
				}
				return true
			}
			return false
		},
	}}
}

// errCodeGen produces a file with a set for a single type.
type errCodeGen struct {
	generator.DefaultGen
	outputPackage string
	imports       namer.ImportTracker
}

func (g *errCodeGen) FileType() string {
	return jsonFileType
}

func (g *errCodeGen) Filename() string {
	return "error_code.json"
}

func (g *errCodeGen) Namers(c *generator.Context) namer.NameSystems {
	return namer.NameSystems{
		"raw": namer.NewRawNamer(g.outputPackage, g.imports),
	}
}

func (g *errCodeGen) Finalize(c *generator.Context, w io.Writer) (err error) {
	// typesMap= map[types.Name]generators.constTypeInfo{
	//     types.Name{Package:"k8s.io/gengo/examples/error-code/inputs", Name:"ErrCodeRbdVolumeBegin", Path:""}:generators.constTypeInfo{tag:"Rbd,Volume", name:"ErrCodeRbdVolumeBegin", val:0, t:(*types.Type)(0xc0000e63c0)},
	//     types.Name{Package:"k8s.io/gengo/examples/error-code/inputs", Name:"ErrCodeRbdVolumeEnd", Path:""}:generators.constTypeInfo{tag:"Rbd,Volume", name:"ErrCodeRbdVolumeEnd", val:3, t:(*types.Type)(0xc0000e6780)},
	//     types.Name{Package:"k8s.io/gengo/examples/error-code/inputs", Name:"ErrCodeRbdVolumeUnknownParameter", Path:""}:generators.constTypeInfo{tag:"Rbd,Volume", name:"ErrCodeRbdVolumeUnknownParameter", val:1, t:(*types.Type)(0xc0000e6d80)}
	//     types.Name{Package:"k8s.io/gengo/examples/error-code/inputs", Name:"ErrCodeRbdVolume", Path:""}:generators.constTypeInfo{tag:"Rbd", name:"ErrCodeRbdVolume", val:2, t:(*types.Type)(0xc0017c7d40)},
	//     types.Name{Package:"k8s.io/gengo/examples/error-code/inputs", Name:"ErrCodeRbd", Path:""}:generators.constTypeInfo{tag:"", name:"ErrCodeRbd", val:2, t:(*types.Type)(0xc0017c7080)},
	// }
	fmt.Printf("\n%%$$$$typesMap= %#v\n", typesMap)
	rootNodeMap := make(map[string]*node)
	for _, info := range typesMap {
		// fmt.Printf("\n**type = %#v, info = %#v\n", t, info)
		// module
		if info.tag == "" {
			module := strings.TrimPrefix(info.name, errCodeTag)
			if _, ok := rootNodeMap[module]; !ok {
				rootNodeMap[module] = newRootNode(module, info.val.(int64))
			} else if !rootNodeMap[module].hasValue() {
				rootNodeMap[module].setValue(info.val.(int64))
			}
			// fmt.Printf("\t&&&Create a module node: %#v\n", rootNodeMap[module])
		} else if strings.Index(info.tag, ",") == -1 {
			if _, ok := rootNodeMap[info.tag]; !ok {
				rootNodeMap[info.tag] = newRootNode(info.tag, -1)
			}
			var subModuleNode *node
			var existed bool
			var modVal, subModVal int
			if _, ok := info.val.(int64); ok {
				var subModuleName = strings.TrimPrefix(info.name, catenate(errCodeTag, info.tag))
				if subModuleNode, existed = rootNodeMap[info.tag].child(subModuleName); !existed {
					subModuleNode = newInterNode(subModuleName, info.val.(int64),
						rootNodeMap[info.tag])
					// fmt.Printf("\t\t&&&Create a sub module node: %#v\n", subModuleNode)
				} else if !subModuleNode.hasValue() {
					subModuleNode.setValue(info.val.(int64))
					// fmt.Printf("\t\t&&&Set value of sub module node: %#v\n", subModuleNode)
				}
			} else if _, ok = info.val.(string); ok {
				// +ErrCode=Common
				// var ErrCodeCommonPrefix = "01-01"
				values := strings.Split(info.val.(string), "-")
				if len(values) != 2 {
					return fmt.Errorf("Invalid error code prefix: %s", info.val)
				}
				if modVal, err = strconv.Atoi(values[0]); err != nil {
					return fmt.Errorf("Failed to converse %s to integer", values[0])
				}
				rootNodeMap[info.tag].setValue(int64(modVal))
				if subModVal, err = strconv.Atoi(values[1]); err != nil {
					return fmt.Errorf("Failed to converse %s to integer", values[1])
				}
				if subModuleNode, existed = rootNodeMap[info.tag].child(info.tag); !existed {
					subModuleNode = newInterNode(info.tag, int64(subModVal),
						rootNodeMap[info.tag])
					// fmt.Printf("\t\t&&&Create a sub module node: %#v\n", subModuleNode)
				} else if !subModuleNode.hasValue() {
					subModuleNode.setValue(int64(subModVal))
					// fmt.Printf("\t\t&&&Set value of sub module node: %#v\n", subModuleNode)
				}
			} else {
				return fmt.Errorf("Invalid type of info value: %T", info.val)
			}
		} else {
			modules := strings.Split(info.tag, ",")
			if len(modules) != 2 {
				return fmt.Errorf("The length of module tags is larger than 2")
			}
			id := strings.TrimPrefix(info.name, catenate(catenate(errCodeTag, modules[0]), modules[1]))
			if id == beginTag || id == endTag {
				continue
			}
			// handle common module
			if modules[1] == "" {
				modules[1] = modules[0]
			}
			var parent = rootNodeMap[modules[0]]
			var subModuleNode *node
			var existed bool
			if parent == nil {
				parent = newRootNode(modules[0], -1)
				subModuleNode = newInterNode(modules[1], -1, parent)
				rootNodeMap[modules[0]] = parent
				// fmt.Printf("\t\t&&&Create a root node %#v and sub module node: %#v\n", parent, subModuleNode)
			}
			if subModuleNode, existed = parent.child(modules[1]); !existed {
				subModuleNode = newInterNode(modules[1], -1, parent)
				// fmt.Printf("\t\t&&&Create a sub module node: %#v\n", subModuleNode)
			}
			newLeafNode(info.name, info.val.(int64), subModuleNode)
			// leafNode := newLeafNode(info.name, info.val.(int64), subModuleNode)
			// fmt.Printf("\t\t\t&&&Create a leaf node: %#v\n", leafNode)
		}
	}
	for k, v := range rootNodeMap {
		// fmt.Printf("\n***rootNodeMap[key=%#v, value = %#v]\n", k, v)
		nodeStr, err := getString(v, "", true)
		if err != nil {
			fmt.Printf("\n***the error is %#v\n", err)
		}
		fmt.Printf("\n***rootNodeMap[key=%#v, value = \n%s]\n", k, nodeStr)

		fmt.Println("\n##################################")
		// for i, n := range v.getLeaves() {
		// 	fmt.Printf("---> %d, %#v\n", i+1, n)
		// }

	}
	errCodes := make(map[string]errCodeDesc)
	for _, v := range rootNodeMap {
		for _, n := range v.getLeaves() {
			// fmt.Printf("***The length of leaves: %d\n", len(v.getLeaves()))
			var mod, subMod int64
			subModuleNode := n.getParent()
			subMod = subModuleNode.getValue()
			moduleNode := subModuleNode.getParent()
			mod = moduleNode.getValue()
			errCodes[catenateErrCode(mod, subMod, n.getValue())] = errCodeDesc{
				Desc: strings.TrimPrefix(n.getName(), errCodeTag),
			}
		}
	}

	fmt.Println(errCodes)
	var codeBytes []byte
	if codeBytes, err = json.MarshalIndent(errCodes, "", "\t"); err != nil {
		return nil
	}
	w.Write(codeBytes)
	return nil
}

type node struct {
	isRoot   bool
	parent   *node
	name     string
	val      int64
	children []*node
	isLeaf   bool
}

func (n *node) getName() string {
	return n.name
}

func (n *node) getParent() *node {
	if n.isRoot {
		return nil
	}
	return n.parent
}

func (n *node) getValue() int64 {
	return n.val
}

func (n *node) getLeaves() []*node {
	nodes := []*node{}
	if n.isLeaf {
		nodes = append(nodes, n)
	}
	for _, child := range n.children {
		nodes = append(nodes, child.getLeaves()...)
	}
	return nodes
}

func (n *node) hasValue() bool {
	// -1 stands for empty value
	return n.val != -1
}

func (n *node) setValue(val int64) {
	n.val = val
}

func (n *node) child(name string) (*node, bool) {
	if n.children == nil || len(n.children) == 0 {
		return nil, false
	}
	for _, node := range n.children {
		if node.name == name {
			return node, true
		}
	}
	return nil, false
}

func (n *node) addChild(child *node) bool {
	if n.isLeaf {
		return false
	}
	if n.children == nil {
		n.children = []*node{}
	}
	n.children = append(n.children, child)
	return true
}

func getString(n *node, prefix string, isTail bool) (string, error) {
	var buffer bytes.Buffer
	var err error
	buffer.WriteString(prefix)
	if isTail {
		_, err = buffer.WriteString("└── ")
		if err != nil {
			return "", err
		}
	} else {
		_, err = buffer.WriteString("├── ")
		if err != nil {
			return "", err
		}
	}
	buffer.WriteString(fmt.Sprintf("Node[name = %s, value = %d]\n", n.name, n.val))
	for i, child := range n.children {
		if i == len(n.children)-1 {
			s, err := getString(child, prefix+"   ", true)
			if err != nil {
				return "", err
			}
			_, err = buffer.WriteString(s)
			if err != nil {
				return "", err
			}
		} else {
			// buffer.WriteString(getString(child, prefix+"   ", false))
			s, err := getString(child, prefix+"   ", false)
			if err != nil {
				return "", err
			}
			_, err = buffer.WriteString(s)
			if err != nil {
				return "", err
			}
		}
	}
	return buffer.String(), nil
}

// func (n *node) String() string {
// 	if n.isLeaf {
// 		return fmt.Sprintf("\t---Leaf[name = %s, value = %d]", n.name, n.val)
// 	}
// 	var s string
// 	if n.isRoot {
// 		s = fmt.Sprintf("Root[name = %s, value = %d]\n", n.name, n.val)
// 	}
// 	for _, child := range n.children {
// 		s = fmt.Sprintf("%s\tIntermediate[name = %s, value = %d], --> %s\n", s, child.name, child.val, child.String())
// 		// s = fmt.Sprintf("%s\n\t%s", s, child.String())
// 	}
// 	return s
// }

func newRootNode(name string, val int64) *node {
	return &node{
		isRoot:   true,
		parent:   nil,
		name:     name,
		val:      val,
		children: []*node{},
		isLeaf:   false,
	}
}

func newInterNode(name string, val int64, parent *node, children ...*node) *node {
	n := &node{
		isRoot:   false,
		parent:   parent,
		name:     name,
		val:      val,
		children: children,
		isLeaf:   false,
	}
	parent.addChild(n)
	return n
}

func newLeafNode(name string, val int64, parent *node) *node {
	n := &node{
		isRoot:   false,
		parent:   parent,
		name:     name,
		val:      val,
		children: nil,
		isLeaf:   true,
	}
	parent.addChild(n)
	return n
}

func jsonFormat(b []byte) ([]byte, error) {
	return b, nil
}

func jsonAssemble(w io.Writer, f *generator.File) {
	// fmt.Fprintf(w, "{\n")
	w.Write(f.Body.Bytes())
	// fmt.Fprintf(w, "}")
}

type jsonFile struct {
	c        *generator.Context
	Format   func([]byte) ([]byte, error)
	Assemble func(io.Writer, *generator.File)
}

func (jf jsonFile) AssembleFile(f *generator.File, pathname string) error {
	destFile, err := os.Create(pathname)
	if err != nil {
		return err
	}
	defer destFile.Close()

	b := &bytes.Buffer{}
	et := generator.NewErrorTracker(b)
	jf.Assemble(et, f)
	if et.Error() != nil {
		return et.Error()
	}
	var formatted []byte
	if formatted, err = jf.Format(b.Bytes()); err != nil {
		err = fmt.Errorf("unable to format file %q (%v)", pathname, err)
		// Write the file anyway, so they can see what's going wrong and fix the generator.
		if _, err2 := destFile.Write(b.Bytes()); err2 != nil {
			return err2
		}
		return err
	}
	_, err = destFile.Write(formatted)
	return err
}

func (jf jsonFile) VerifyFile(f *generator.File, path string) error {
	return nil
}

func catenate(str1, str2 string) string {
	return fmt.Sprintf("%s%s", str1, str2)
}

func catenateErrCode(module, subModule, code int64) string {
	return fmt.Sprintf("%02X-%02X-%04X", module, subModule, code)
}

type errCodeDesc struct {
	Desc string `json:"desc"`
}

func (code errCodeDesc) String() string {
	return code.Desc
}
