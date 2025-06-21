// backend/codegen/intermediate/ir.go
package intermediate

import (
	"fmt"
	"strings"
)

// IROpcode representa las operaciones de la representación intermedia
type IROpcode int

const (
	// Operaciones de carga y almacenamiento
	IR_LOAD IROpcode = iota
	IR_STORE
	IR_LOAD_IMMEDIATE
	IR_LOAD_ADDRESS

	// Operaciones aritméticas
	IR_ADD
	IR_SUB
	IR_MULT
	IR_DIV
	IR_MOD
	IR_NEG

	// Operaciones lógicas
	IR_AND
	IR_OR
	IR_NOT

	// Operaciones de comparación
	IR_CMP_EQ
	IR_CMP_NE
	IR_CMP_LT
	IR_CMP_LE
	IR_CMP_GT
	IR_CMP_GE

	// Control de flujo
	IR_LABEL
	IR_BRANCH
	IR_BRANCH_IF_TRUE
	IR_BRANCH_IF_FALSE
	IR_CALL
	IR_RETURN
	IR_ENTER_FUNCTION
	IR_EXIT_FUNCTION

	// Operaciones de memoria
	IR_ALLOC_LOCAL
	IR_ALLOC_PARAM
	IR_PUSH
	IR_POP

	// Operaciones de conversión
	IR_INT_TO_FLOAT
	IR_FLOAT_TO_INT
	IR_TO_STRING

	// Operaciones de vectores/arrays
	IR_LOAD_ARRAY_ELEMENT
	IR_STORE_ARRAY_ELEMENT
	IR_ARRAY_LENGTH

	// I/O básicas
	IR_PRINT
	IR_PRINT_LN

	// Noop
	IR_NOP
)

// String convierte IROpcode a string para debugging
func (op IROpcode) String() string {
	switch op {
	case IR_LOAD:
		return "LOAD"
	case IR_STORE:
		return "STORE"
	case IR_LOAD_IMMEDIATE:
		return "LOAD_IMM"
	case IR_LOAD_ADDRESS:
		return "LOAD_ADDR"
	case IR_ADD:
		return "ADD"
	case IR_SUB:
		return "SUB"
	case IR_MULT:
		return "MULT"
	case IR_DIV:
		return "DIV"
	case IR_MOD:
		return "MOD"
	case IR_NEG:
		return "NEG"
	case IR_AND:
		return "AND"
	case IR_OR:
		return "OR"
	case IR_NOT:
		return "NOT"
	case IR_CMP_EQ:
		return "CMP_EQ"
	case IR_CMP_NE:
		return "CMP_NE"
	case IR_CMP_LT:
		return "CMP_LT"
	case IR_CMP_LE:
		return "CMP_LE"
	case IR_CMP_GT:
		return "CMP_GT"
	case IR_CMP_GE:
		return "CMP_GE"
	case IR_LABEL:
		return "LABEL"
	case IR_BRANCH:
		return "BRANCH"
	case IR_BRANCH_IF_TRUE:
		return "BRANCH_IF_TRUE"
	case IR_BRANCH_IF_FALSE:
		return "BRANCH_IF_FALSE"
	case IR_CALL:
		return "CALL"
	case IR_RETURN:
		return "RETURN"
	case IR_ENTER_FUNCTION:
		return "ENTER_FUNC"
	case IR_EXIT_FUNCTION:
		return "EXIT_FUNC"
	case IR_ALLOC_LOCAL:
		return "ALLOC_LOCAL"
	case IR_ALLOC_PARAM:
		return "ALLOC_PARAM"
	case IR_PUSH:
		return "PUSH"
	case IR_POP:
		return "POP"
	case IR_INT_TO_FLOAT:
		return "INT_TO_FLOAT"
	case IR_FLOAT_TO_INT:
		return "FLOAT_TO_INT"
	case IR_TO_STRING:
		return "TO_STRING"
	case IR_LOAD_ARRAY_ELEMENT:
		return "LOAD_ARRAY_ELEM"
	case IR_STORE_ARRAY_ELEMENT:
		return "STORE_ARRAY_ELEM"
	case IR_ARRAY_LENGTH:
		return "ARRAY_LENGTH"
	case IR_PRINT:
		return "PRINT"
	case IR_PRINT_LN:
		return "PRINT_LN"
	case IR_NOP:
		return "NOP"
	default:
		return "UNKNOWN"
	}
}

// IROperandType representa el tipo de operando
type IROperandType int

const (
	IR_OPERAND_TEMP         IROperandType = iota // Temporal (t1, t2, etc.)
	IR_OPERAND_VAR                               // Variable
	IR_OPERAND_IMMEDIATE                         // Valor inmediato
	IR_OPERAND_LABEL                             // Etiqueta
	IR_OPERAND_PARAM                             // Parámetro de función
	IR_OPERAND_GLOBAL                            // Variable global
	IR_OPERAND_ARRAY_ACCESS                      // Acceso a array
)

// IROperand representa un operando en la representación intermedia
type IROperand struct {
	Type     IROperandType
	Name     string      // Nombre del operando
	Value    interface{} // Valor (para inmediatos)
	DataType string      // Tipo de dato (int, float, string, etc.)
	Offset   int         // Offset para variables locales/parámetros
	Index    *IROperand  // Para acceso a arrays
}

// String convierte IROperand a string
func (op *IROperand) String() string {
	if op == nil {
		return "nil"
	}

	switch op.Type {
	case IR_OPERAND_TEMP:
		return fmt.Sprintf("%%%s", op.Name)
	case IR_OPERAND_VAR:
		return fmt.Sprintf("$%s", op.Name)
	case IR_OPERAND_IMMEDIATE:
		return fmt.Sprintf("#%v", op.Value)
	case IR_OPERAND_LABEL:
		return fmt.Sprintf("@%s", op.Name)
	case IR_OPERAND_PARAM:
		return fmt.Sprintf("param_%s", op.Name)
	case IR_OPERAND_GLOBAL:
		return fmt.Sprintf("global_%s", op.Name)
	case IR_OPERAND_ARRAY_ACCESS:
		if op.Index != nil {
			return fmt.Sprintf("$%s[%s]", op.Name, op.Index.String())
		}
		return fmt.Sprintf("$%s[]", op.Name)
	default:
		return fmt.Sprintf("unknown_%s", op.Name)
	}
}

// IsImmediate verifica si el operando es un valor inmediato
func (op *IROperand) IsImmediate() bool {
	return op != nil && op.Type == IR_OPERAND_IMMEDIATE
}

// IsTemp verifica si el operando es un temporal
func (op *IROperand) IsTemp() bool {
	return op != nil && op.Type == IR_OPERAND_TEMP
}

// IsVariable verifica si el operando es una variable
func (op *IROperand) IsVariable() bool {
	return op != nil && (op.Type == IR_OPERAND_VAR || op.Type == IR_OPERAND_PARAM || op.Type == IR_OPERAND_GLOBAL)
}

// IRInstruction representa una instrucción en la representación intermedia
type IRInstruction struct {
	Op      IROpcode   // Operación
	Dest    *IROperand // Operando destino
	Src1    *IROperand // Primer operando fuente
	Src2    *IROperand // Segundo operando fuente
	Label   string     // Etiqueta (para saltos)
	Comment string     // Comentario para debugging
	LineNo  int        // Número de línea en el código fuente
}

// String convierte IRInstruction a string para debugging
func (instr *IRInstruction) String() string {
	var parts []string

	// Agregar etiqueta si existe
	if instr.Label != "" {
		parts = append(parts, fmt.Sprintf("%s:", instr.Label))
	}

	// Instrucción base
	instrStr := instr.Op.String()

	// Agregar operandos según el tipo de instrucción
	switch instr.Op {
	case IR_LABEL:
		instrStr = instr.Label + ":"
	case IR_BRANCH:
		if instr.Src1 != nil {
			instrStr = fmt.Sprintf("%s %s", instrStr, instr.Src1.String())
		}
	case IR_BRANCH_IF_TRUE, IR_BRANCH_IF_FALSE:
		if instr.Src1 != nil && instr.Src2 != nil {
			instrStr = fmt.Sprintf("%s %s, %s", instrStr, instr.Src1.String(), instr.Src2.String())
		}
	case IR_CALL:
		if instr.Dest != nil && instr.Src1 != nil {
			instrStr = fmt.Sprintf("%s %s = %s", instrStr, instr.Dest.String(), instr.Src1.String())
		} else if instr.Src1 != nil {
			instrStr = fmt.Sprintf("%s %s", instrStr, instr.Src1.String())
		}
	case IR_RETURN:
		if instr.Src1 != nil {
			instrStr = fmt.Sprintf("%s %s", instrStr, instr.Src1.String())
		}
	case IR_LOAD_IMMEDIATE:
		if instr.Dest != nil && instr.Src1 != nil {
			instrStr = fmt.Sprintf("%s %s, %s", instrStr, instr.Dest.String(), instr.Src1.String())
		}
	case IR_ALLOC_LOCAL:
		if instr.Dest != nil {
			instrStr = fmt.Sprintf("%s %s [%d bytes]", instrStr, instr.Dest.String(), instr.Dest.Offset)
		}
	default:
		// Formato general: OP dest, src1, src2
		if instr.Dest != nil {
			if instr.Src1 != nil && instr.Src2 != nil {
				instrStr = fmt.Sprintf("%s %s, %s, %s", instrStr, instr.Dest.String(), instr.Src1.String(), instr.Src2.String())
			} else if instr.Src1 != nil {
				instrStr = fmt.Sprintf("%s %s, %s", instrStr, instr.Dest.String(), instr.Src1.String())
			} else {
				instrStr = fmt.Sprintf("%s %s", instrStr, instr.Dest.String())
			}
		} else if instr.Src1 != nil {
			instrStr = fmt.Sprintf("%s %s", instrStr, instr.Src1.String())
		}
	}

	parts = append(parts, instrStr)

	// Agregar comentario si existe
	if instr.Comment != "" {
		parts = append(parts, fmt.Sprintf("; %s", instr.Comment))
	}

	return strings.Join(parts, " ")
}

// IRFunction representa una función en IR
type IRFunction struct {
	Name         string           // Nombre de la función
	Parameters   []*IROperand     // Parámetros
	ReturnType   string           // Tipo de retorno
	LocalVars    []*IROperand     // Variables locales
	Instructions []*IRInstruction // Instrucciones
	StackSize    int              // Tamaño del stack frame
}

// String convierte IRFunction a string
func (fn *IRFunction) String() string {
	var lines []string

	// Header de la función
	paramNames := make([]string, len(fn.Parameters))
	for i, param := range fn.Parameters {
		paramNames[i] = param.String()
	}

	header := fmt.Sprintf("function %s(%s) -> %s {",
		fn.Name, strings.Join(paramNames, ", "), fn.ReturnType)
	lines = append(lines, header)

	// Variables locales
	if len(fn.LocalVars) > 0 {
		lines = append(lines, "  ; Local variables:")
		for _, localVar := range fn.LocalVars {
			lines = append(lines, fmt.Sprintf("  ; %s [offset: %d]", localVar.String(), localVar.Offset))
		}
		lines = append(lines, "")
	}

	// Instrucciones
	for _, instr := range fn.Instructions {
		if instr.Op == IR_LABEL {
			lines = append(lines, fmt.Sprintf("%s:", instr.Label))
		} else {
			lines = append(lines, fmt.Sprintf("  %s", instr.String()))
		}
	}

	lines = append(lines, "}")
	return strings.Join(lines, "\n")
}

// IRProgram representa un programa completo en IR
type IRProgram struct {
	Functions   []*IRFunction  // Funciones del programa
	GlobalVars  []*IROperand   // Variables globales
	StringTable map[string]int // Tabla de strings literales
}

// String convierte IRProgram a string
func (prog *IRProgram) String() string {
	var lines []string

	// Variables globales
	if len(prog.GlobalVars) > 0 {
		lines = append(lines, "; Global variables:")
		for _, globalVar := range prog.GlobalVars {
			lines = append(lines, fmt.Sprintf("; %s", globalVar.String()))
		}
		lines = append(lines, "")
	}

	// Tabla de strings
	if len(prog.StringTable) > 0 {
		lines = append(lines, "; String table:")
		for str, id := range prog.StringTable {
			lines = append(lines, fmt.Sprintf("; str_%d: \"%s\"", id, str))
		}
		lines = append(lines, "")
	}

	// Funciones
	for i, function := range prog.Functions {
		if i > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, function.String())
	}

	return strings.Join(lines, "\n")
}

// AddFunction agrega una función al programa
func (prog *IRProgram) AddFunction(fn *IRFunction) {
	prog.Functions = append(prog.Functions, fn)
}

// AddGlobalVar agrega una variable global al programa
func (prog *IRProgram) AddGlobalVar(operand *IROperand) {
	prog.GlobalVars = append(prog.GlobalVars, operand)
}

// AddString agrega un string literal a la tabla y retorna su ID
func (prog *IRProgram) AddString(str string) int {
	if prog.StringTable == nil {
		prog.StringTable = make(map[string]int)
	}

	if id, exists := prog.StringTable[str]; exists {
		return id
	}

	id := len(prog.StringTable)
	prog.StringTable[str] = id
	return id
}

// NewIRProgram crea un nuevo programa IR
func NewIRProgram() *IRProgram {
	return &IRProgram{
		Functions:   make([]*IRFunction, 0),
		GlobalVars:  make([]*IROperand, 0),
		StringTable: make(map[string]int),
	}
}
