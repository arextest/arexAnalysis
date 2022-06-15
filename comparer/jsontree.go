package comparer

import (
	"errors"
	"fmt"
	"hash/fnv"

	json "github.com/json-iterator/go"
)

// PathError : an error that includes path information
type PathError struct {
	Path string
	Err  error
}

func newPathError(path string, v ...interface{}) *PathError {
	return &PathError{
		Path: path,
		Err:  errors.New(fmt.Sprint(v...)),
	}
}

func newPathErrorf(path, format string, v ...interface{}) *PathError {
	return &PathError{
		Path: path,
		Err:  fmt.Errorf(format, v...),
	}
}

func (err *PathError) Error() string {
	return fmt.Sprintf("%v; %s", err.Err, err.Path)
}

// JSONType json type
type JSONType uint8

const (
	// Object 0
	Object JSONType = iota
	Array
	String
	Number
	Boolean
	Null
	Error
)

var jsonTypeStrings = []string{
	Object:  "object",
	Array:   "array",
	String:  "string",
	Number:  "number",
	Boolean: "boolean",
	Null:    "null",
	Error:   "error",
}

func (t JSONType) String() string {
	if int(t) > len(jsonTypeStrings) {
		return fmt.Sprintf("Unknown (%d)", t)
	}
	return jsonTypeStrings[t]
}

// Fingerprint 计算
func Fingerprint(b []byte) uint64 {
	hash := fnv.New64a()
	hash.Write(b)
	return hash.Sum64()
}

//JSONNode JSON TREE
type JSONNode struct {
	parent      *JSONNode
	children    []*JSONNode
	class       JSONType
	key         string
	val         interface{}
	init        bool
	index       int
	len         int
	fingerprint uint64
	err         *error
}

// NewJSONTree : creates an empty *JsonTree. initialize the tree with json.Unmarshal().
func NewJSONTree() *JSONNode {
	return newNode(nil)
}

func newNode(val interface{}) *JSONNode {
	return &JSONNode{
		index: -1,
		val:   val,
	}
}

// NewNodeByInterface : CREATE NULL Node
func NewNodeByInterface(val interface{}) *JSONNode {
	tree := newNode(val)
	tree.getType()
	return tree
}

// NewNull : CREATE NULL Node
func NewNull() *JSONNode {
	tree := NewJSONTree()
	tree.getType()
	return tree
}

// NewString : ??
func NewString(s string) *JSONNode {
	tree := newNode(s)
	tree.getType()
	return tree
}

// NewNumber node
func NewNumber(x float64) *JSONNode {
	tree := newNode(x)
	tree.getType()
	return tree
}

// NewBoolean bool
func NewBoolean(b bool) *JSONNode {
	tree := newNode(b)
	tree.getType()
	return tree
}

// NewObject add object
func NewObject(o map[string]interface{}) *JSONNode {
	if o == nil {
		o = make(map[string]interface{}, 0)
	}
	tree := newNode(o)
	tree.getType()
	return tree
}

// NewArray interface{} array
func NewArray(a []interface{}) *JSONNode {
	if a == nil {
		a = make([]interface{}, 0, 1)
	}
	tree := newNode(a)
	tree.getType()
	return tree
}

// Type return type info
func (tree *JSONNode) Type() JSONType {
	return tree.class
}

// Len 长度
func (tree *JSONNode) Len() (int, error) {
	switch tree.Type() {
	case Object:
		return len(tree.val.(map[string]interface{})), nil
	case Array:
		return len(tree.val.([]interface{})), nil
	default:
		return 0, fmt.Errorf("not an array or an object (%v); %s", tree.Type(), tree.path())
	}
}

// Err any error encountered due to non-existent keys, out of range indices, etc.
func (tree *JSONNode) Err() error {
	if tree.err != nil {
		return *tree.err
	}
	return nil
}

// Root the root of tree. this value can not be nil.
func (tree *JSONNode) Root() *JSONNode {
	if tree.parent == nil {
		return tree
	}
	iterTree := tree.parent
	for iterTree.parent != nil {
		iterTree = iterTree.parent
	}
	return iterTree
}

// Parent tree's parent. this value can be nil
func (tree *JSONNode) Parent() *JSONNode {
	return tree.parent
}

// GetIndex a *JsonTree representing the i-th element in the array tree. if tree.Err()
// is not nil then the Err() method of returned *JsonTree returns the same error.
// if tree is not an array then the Err() method of the returned *JsonTree's
// returns a *PathError.
func (tree *JSONNode) GetIndex(i int) *JSONNode {
	child := NewJSONTree()
	defer child.getType()
	child.index = i
	child.parent = tree
	child.err = tree.err
	if child.err != nil {
		return child
	}
	switch {
	case !tree.init:
		child.errUninitialized()
	case tree.class == Array:
		a := tree.val.([]interface{})
		if 0 <= i && i < len(a) {
			child.val = a[i]
		} else {
			child.errIndexOutOfRange()
		}
	default:
		child.errTypeError(Array)
	}
	return child
}

// Get a *JsonTree representing the value of key in the object tree. if tree.Err()
// is not nil then the Err() method of returned *JsonTree returns the same error.
// if tree is not an object then the Err() method of the returned *JsonTree's
// returns a *PathError.
func (tree *JSONNode) Get(key string) *JSONNode {
	child := NewJSONTree()
	defer child.getType()

	child.key = key
	child.parent = tree
	child.err = tree.err
	if child.err != nil {
		return child
	}
	switch {
	case !tree.init:
		child.errUninitialized()
	case tree.class == Object:
		val, ok := tree.val.(map[string]interface{})[key]
		if ok {
			child.val = val
		} else {
			child.errNoExist()
		}
	default:
		child.errTypeError(Object)
	}
	return child
}

// Interface val
func (tree *JSONNode) Interface() (interface{}, error) {
	return tree.val, tree.Err()
}

// converts tree to a string. returns a *PathError if tree is not a string.
func (tree *JSONNode) String() (string, error) {
	if !tree.init {
		return "", newPathErrorf(tree.path(), "uninitialized")
	}
	switch tree.class {
	case Error:
		return "", *tree.err
	case String:
		return tree.val.(string), nil
	default:
		return "", newPathErrorf(tree.path(), "not a string (%v)", tree.Type())
	}
}

// Number converts tree to a number. returns a *PathError if tree is not a number.
func (tree *JSONNode) Number() (float64, error) {
	if !tree.init {
		return 0, newPathErrorf(tree.path(), "uninitialized")
	}
	switch tree.class {
	case Error:
		return 0, *tree.err
	case Number:
		return tree.val.(float64), nil
	default:
		return 0, newPathErrorf(tree.path(), "not a number (%v)", tree.Type())
	}
}

// Boolean converts tree to a bool. returns a *PathError if tree is not a boolean.
func (tree *JSONNode) Boolean() (bool, error) {
	if !tree.init {
		return false, newPathErrorf(tree.path(), "uninitialized")
	}
	switch tree.class {
	case Error:
		return false, *tree.err
	case Boolean:
		return tree.val.(bool), nil
	default:
		return false, newPathErrorf(tree.path(), "not a bool (%v)", tree.Type())
	}
}

// Array converts tree to a slice. returns a *PathError if tree is not an array.
func (tree *JSONNode) Array() ([]interface{}, error) {
	if !tree.init {
		return nil, newPathErrorf(tree.path(), "uninitialized")
	}
	switch tree.class {
	case Error:
		return nil, *tree.err
	case Array:
		return tree.val.([]interface{}), nil
	default:
		return nil, newPathErrorf(tree.path(), "not an array (%v)", tree.Type())
	}
}

// Object converts tree to a map. returns a *PathError if tree is not an object.
func (tree *JSONNode) Object() (map[string]interface{}, error) {
	if !tree.init {
		return nil, newPathErrorf(tree.path(), "uninitialized")
	}
	switch tree.class {
	case Error:
		return nil, *tree.err
	case Object:
		return tree.val.(map[string]interface{}), nil
	default:
		return nil, newPathErrorf(tree.path(), "not an object (%v)", tree.Type())
	}
}

// IsNull returns true if tree is null. returns false in otherwise
// (other type, error, non existing keys, ...).
func (tree *JSONNode) IsNull() bool {
	return tree.class == Null
}

// UnmarshalJSON implements json.Unmarshaler
func (tree *JSONNode) UnmarshalJSON(p []byte) error {
	defer tree.getType()
	mapJSON := make(map[string]interface{})
	err := json.Unmarshal(p, &mapJSON)
	if err != nil {
		fmt.Println(err)
		return err
	}
	tree.buildTree(mapJSON)
	return nil
}

func (tree *JSONNode) buildTree(mapJSON map[string]interface{}) error {
	for k, v := range mapJSON {
		fmt.Printf("k %s v %s", k, v)
		oneNode := NewNodeByInterface(v)
		oneNode.key = k
	}
	return nil
}

// MarshalJSON implements json.Marshaler
func (tree *JSONNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(tree.val)
}

func (tree *JSONNode) newError(v ...interface{}) {
	err := error(newPathError(tree.path(), v...))
	tree.err = &err
}

func (tree *JSONNode) newErrorf(format string, v ...interface{}) {
	err := error(newPathErrorf(tree.path(), format, v...))
	tree.err = &err
}

func (tree *JSONNode) errUninitialized() {
	tree.newError("uninitialized")
}
func (tree *JSONNode) errNoExist() {
	tree.newError("key does not exist")
}
func (tree *JSONNode) errIndexOutOfRange() {
	tree.newErrorf("index out of range")
}
func (tree *JSONNode) errTypeError(expected JSONType) {
	if (expected == Object) || expected == Array {
		tree.newErrorf("not an %v (%v)", expected, tree.Type())
	} else {
		tree.newErrorf("not a %v (%v)", expected, tree.Type())
	}
}

func (tree *JSONNode) path() string {
	if tree.parent == nil {
		return "$"
	}
	pre := tree.parent.path()
	if tree.index >= 0 {
		return fmt.Sprintf("%s[%d]", pre, tree.index)
	} else {
		return fmt.Sprintf("%s.%s", pre, tree.key)
	}
}

func (tree *JSONNode) getType() {
	tree.init = true
	if tree.err != nil {
		tree.class = Error
		return
	}
	switch tree.val.(type) {
	case string:
		tree.class = String
	case float64:
		tree.class = Number
	case bool:
		tree.class = Boolean
	case nil:
		tree.class = Null
	case []interface{}:
		tree.class = Array
	case map[string]interface{}:
		tree.class = Object
	}
}
