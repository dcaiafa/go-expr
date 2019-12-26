package runtime

import "github.com/dcaiafa/go-expr/expr/types"

// Builder builds a Program using low level primitives.
type Builder struct {
	labels    []*Label
	funcs     []Func
	strings   []string
	stringMap map[string]int
	instr     []Instruction
	exprs     []Expr
	consts    []Value
	inputs    []types.Type
}

// NewBuilder creates a new Builder.
func NewBuilder() *Builder {
	return &Builder{
		stringMap: make(map[string]int),
	}
}

// NewInput creates a new input that can be referenced in a LoadInput
// instruction.
func (b *Builder) NewInput(t types.Type) int {
	b.inputs = append(b.inputs, t)
	return len(b.inputs) - 1
}

// NewConst creates a new constant that can be referenced in a LoadConst
// instruction.
func (b *Builder) NewConst(v Value) int {
	constIndex := len(b.consts)
	b.consts = append(b.consts, v)
	return constIndex
}

// NewExternalFunc creates a new external function that can be referenced in a
// Call instruction.
func (b *Builder) NewExternalFunc(fn Func) int {
	funcIndex := len(b.funcs)
	b.funcs = append(b.funcs, fn)
	return funcIndex
}

// NewLabel creates a new label that can be used in EmitJump. The label is
// immediately ready to be used, but it must be assigned using AssignLabel
// before Build is called.
func (b *Builder) NewLabel() *Label {
	label := &Label{
		index: len(b.labels),
		addr:  -1,
	}
	b.labels = append(b.labels, label)
	return label
}

// AssignLabel assigns the label to the current address of the program being
// built.
func (b *Builder) AssignLabel(label *Label) {
	label.addr = len(b.instr)
}

// EmitOp emits an instruction with a simple operation.
func (b *Builder) EmitOp(op Operation) {
	b.addInstr(Instruction{op: op})
}

// EmitLoadConst emits a LoadConst instruction.
func (b *Builder) EmitLoadConst(constIndex int) {
	b.addInstr(Instruction{op: LoadConst, extra: constIndex})
}

// EmitLoadInput emits a LoadInput instruction.
func (b *Builder) EmitLoadInput(inputIndex int) {
	b.addInstr(Instruction{op: LoadInput, extra: inputIndex})
}

// EmitPushNumber emits a PushNumber instruction.
func (b *Builder) EmitPushNumber(num float64) {
	b.addInstr(Instruction{op: PushNumber, vnum: num})
}

// EmitPushString emits a PushString instruction.
func (b *Builder) EmitPushString(str string) {
	strIndex := b.newString(str)
	b.addInstr(Instruction{op: PushString, extra: strIndex})
}

// EmitPushBool emits a PushBool instruction.
func (b *Builder) EmitPushBool(v bool) {
	extra := 0
	if v {
		extra = 1
	}
	b.addInstr(Instruction{op: PushBool, extra: extra})
}

// EmitJump emits a Jump, JumpIfTrue or JumpIfFalse instruction.
func (b *Builder) EmitJump(op Operation, label *Label) {
	b.addInstr(Instruction{op: op, extra: label.index})
}

// EmitCall emits a Call instruction that consumes argCount arguments from the
// stack.
func (b *Builder) EmitCall(argCount int) {
	b.addInstr(Instruction{op: Call, extra: argCount})
}

func (b *Builder) addInstr(i Instruction) {
	b.instr = append(b.instr, i)
}

// FinishExpr finishes the current expression.
func (b *Builder) FinishExpr() {
	for i := 0; i < len(b.instr); i++ {
		if b.instr[i].op != Jump &&
			b.instr[i].op != JumpIfTrue &&
			b.instr[i].op != JumpIfFalse {
			continue
		}

		label := b.labels[b.instr[i].extra]
		if label.addr == -1 {
			panic("unassigned label")
		}

		b.instr[i].extra = label.addr
	}

	b.exprs = append(b.exprs, Expr(b.instr))
	b.instr = nil
}

// Build returns the Program.
func (b *Builder) Build() *Program {
	return &Program{
		exprs:   b.exprs,
		funcs:   b.funcs,
		strings: b.strings,
		consts:  b.consts,
		inputs:  b.inputs,
	}
}

func (b *Builder) newString(str string) int {
	index, ok := b.stringMap[str]
	if !ok {
		index = len(b.strings)
		b.strings = append(b.strings, str)
		b.stringMap[str] = index
	}
	return index
}
