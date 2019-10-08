package runtime

type Builder struct {
	labels     []*Label
	funcs      []Func
	strings    []string
	stringMap  map[string]int
	instr      []Instruction
	exprs      []Expr
	consts     []Value
	inputCount int
}

func NewBuilder() *Builder {
	return &Builder{
		stringMap: make(map[string]int),
	}
}

func (b *Builder) NewInput() int {
	inputIndex := b.inputCount
	b.inputCount++
	return inputIndex
}

func (b *Builder) RegisterConst(v Value) int {
	constIndex := len(b.consts)
	b.consts = append(b.consts, v)
	return constIndex
}

func (b *Builder) RegisterExternalFunc(fn Func) int {
	funcIndex := len(b.funcs)
	b.funcs = append(b.funcs, fn)
	return funcIndex
}

func (b *Builder) NewLabel() *Label {
	label := &Label{
		index: len(b.labels),
		addr:  -1,
	}
	b.labels = append(b.labels, label)
	return label
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

func (b *Builder) AssignLabel(label *Label) {
	label.addr = len(b.instr)
}

func (b *Builder) EmitOp(op Operation) {
	b.addInstr(Instruction{op: op})
}

func (b *Builder) EmitLoadConst(constIndex int) {
	b.addInstr(Instruction{op: LoadConst, extra: constIndex})
}

func (b *Builder) EmitLoadInput(inputIndex int) {
	b.addInstr(Instruction{op: LoadInput, extra: inputIndex})
}

func (b *Builder) EmitPushNumber(num float64) {
	b.addInstr(Instruction{op: PushNumber, vnum: num})
}

func (b *Builder) EmitPushString(str string) {
	strIndex := b.newString(str)
	b.addInstr(Instruction{op: PushString, extra: strIndex})
}

func (b *Builder) EmitPushBool(v bool) {
	extra := 0
	if v {
		extra = 1
	}
	b.addInstr(Instruction{op: PushBool, extra: extra})
}

func (b *Builder) EmitJump(op Operation, label *Label) {
	b.addInstr(Instruction{op: op, extra: label.index})
}

func (b *Builder) EmitCall(argCount int) {
	b.addInstr(Instruction{op: Call, extra: argCount})
}

func (b *Builder) addInstr(i Instruction) {
	b.instr = append(b.instr, i)
}

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

func (b *Builder) Build() *Program {
	return &Program{
		exprs:      b.exprs,
		funcs:      b.funcs,
		strings:    b.strings,
		consts:     b.consts,
		inputCount: b.inputCount,
	}
}
