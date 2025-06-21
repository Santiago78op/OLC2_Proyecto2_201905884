// backend/codegen/arm64/instructions.go
package arm64

import (
	"fmt"
	"strings"
)

// ARM64Instruction representa una instrucción ARM64
type ARM64Instruction struct {
	Opcode   string   // Mnemónico de la instrucción (mov, add, etc.)
	Operands []string // Operandos de la instrucción
	Comment  string   // Comentario opcional
	Label    string   // Etiqueta opcional
}

// String convierte la instrucción a su representación en assembly
func (instr *ARM64Instruction) String() string {
	var parts []string

	// Agregar etiqueta si existe
	if instr.Label != "" {
		return instr.Label + ":"
	}

	// Construir la instrucción
	if len(instr.Operands) > 0 {
		parts = append(parts, fmt.Sprintf("    %s %s", instr.Opcode, strings.Join(instr.Operands, ", ")))
	} else {
		parts = append(parts, fmt.Sprintf("    %s", instr.Opcode))
	}

	// Agregar comentario si existe
	if instr.Comment != "" {
		parts = append(parts, fmt.Sprintf(" // %s", instr.Comment))
	}

	return strings.Join(parts, "")
}

// InstructionBuilder ayuda a construir instrucciones ARM64
type InstructionBuilder struct {
	instructions []*ARM64Instruction
}

// NewInstructionBuilder crea un nuevo builder
func NewInstructionBuilder() *InstructionBuilder {
	return &InstructionBuilder{
		instructions: make([]*ARM64Instruction, 0),
	}
}

// Add agrega una instrucción al builder
func (ib *InstructionBuilder) Add(opcode string, operands []string, comment string) *InstructionBuilder {
	instr := &ARM64Instruction{
		Opcode:   opcode,
		Operands: operands,
		Comment:  comment,
	}
	ib.instructions = append(ib.instructions, instr)
	return ib
}

// Label agrega una etiqueta
func (ib *InstructionBuilder) Label(label string) *InstructionBuilder {
	instr := &ARM64Instruction{
		Label: label,
	}
	ib.instructions = append(ib.instructions, instr)
	return ib
}

// GetInstructions retorna todas las instrucciones generadas
func (ib *InstructionBuilder) GetInstructions() []*ARM64Instruction {
	return ib.instructions
}

// GetAssembly retorna el código assembly como string
func (ib *InstructionBuilder) GetAssembly() string {
	var lines []string
	for _, instr := range ib.instructions {
		lines = append(lines, instr.String())
	}
	return strings.Join(lines, "\n")
}

// ============ MÉTODOS DE CONVENIENCIA PARA INSTRUCCIONES COMUNES ============

// MOV - Move register or immediate
func (ib *InstructionBuilder) MOV(dest, src string, comment string) *InstructionBuilder {
	return ib.Add("mov", []string{dest, src}, comment)
}

// ADD - Add registers or immediate
func (ib *InstructionBuilder) ADD(dest, src1, src2 string, comment string) *InstructionBuilder {
	return ib.Add("add", []string{dest, src1, src2}, comment)
}

// SUB - Subtract registers or immediate
func (ib *InstructionBuilder) SUB(dest, src1, src2 string, comment string) *InstructionBuilder {
	return ib.Add("sub", []string{dest, src1, src2}, comment)
}

// MUL - Multiply registers
func (ib *InstructionBuilder) MUL(dest, src1, src2 string, comment string) *InstructionBuilder {
	return ib.Add("mul", []string{dest, src1, src2}, comment)
}

// SDIV - Signed divide
func (ib *InstructionBuilder) SDIV(dest, src1, src2 string, comment string) *InstructionBuilder {
	return ib.Add("sdiv", []string{dest, src1, src2}, comment)
}

// LDR - Load register from memory
func (ib *InstructionBuilder) LDR(dest, src string, comment string) *InstructionBuilder {
	return ib.Add("ldr", []string{dest, src}, comment)
}

// STR - Store register to memory
func (ib *InstructionBuilder) STR(src, dest string, comment string) *InstructionBuilder {
	return ib.Add("str", []string{src, dest}, comment)
}

// STP - Store pair of registers
func (ib *InstructionBuilder) STP(reg1, reg2, addr string, comment string) *InstructionBuilder {
	return ib.Add("stp", []string{reg1, reg2, addr}, comment)
}

// LDP - Load pair of registers
func (ib *InstructionBuilder) LDP(reg1, reg2, addr string, comment string) *InstructionBuilder {
	return ib.Add("ldp", []string{reg1, reg2, addr}, comment)
}

// CMP - Compare registers or immediate
func (ib *InstructionBuilder) CMP(reg1, reg2 string, comment string) *InstructionBuilder {
	return ib.Add("cmp", []string{reg1, reg2}, comment)
}

// B - Unconditional branch
func (ib *InstructionBuilder) B(label string, comment string) *InstructionBuilder {
	return ib.Add("b", []string{label}, comment)
}

// BL - Branch with link (function call)
func (ib *InstructionBuilder) BL(label string, comment string) *InstructionBuilder {
	return ib.Add("bl", []string{label}, comment)
}

// BEQ - Branch if equal
func (ib *InstructionBuilder) BEQ(label string, comment string) *InstructionBuilder {
	return ib.Add("b.eq", []string{label}, comment)
}

// BNE - Branch if not equal
func (ib *InstructionBuilder) BNE(label string, comment string) *InstructionBuilder {
	return ib.Add("b.ne", []string{label}, comment)
}

// BLT - Branch if less than
func (ib *InstructionBuilder) BLT(label string, comment string) *InstructionBuilder {
	return ib.Add("b.lt", []string{label}, comment)
}

// BLE - Branch if less than or equal
func (ib *InstructionBuilder) BLE(label string, comment string) *InstructionBuilder {
	return ib.Add("b.le", []string{label}, comment)
}

// BGT - Branch if greater than
func (ib *InstructionBuilder) BGT(label string, comment string) *InstructionBuilder {
	return ib.Add("b.gt", []string{label}, comment)
}

// BGE - Branch if greater than or equal
func (ib *InstructionBuilder) BGE(label string, comment string) *InstructionBuilder {
	return ib.Add("b.ge", []string{label}, comment)
}

// RET - Return from function
func (ib *InstructionBuilder) RET(comment string) *InstructionBuilder {
	return ib.Add("ret", []string{}, comment)
}

// ============ UTILIDADES PARA DIRECCIONAMIENTO ============

// FormatImmediate formatea un valor inmediato
func FormatImmediate(value interface{}) string {
	return fmt.Sprintf("#%v", value)
}

// FormatMemory formatea una dirección de memoria con base + offset
func FormatMemory(base string, offset int) string {
	if offset == 0 {
		return fmt.Sprintf("[%s]", base)
	}
	return fmt.Sprintf("[%s, #%d]", base, offset)
}

// FormatPreIndex formatea direccionamiento pre-index (actualiza base)
func FormatPreIndex(base string, offset int) string {
	return fmt.Sprintf("[%s, #%d]!", base, offset)
}

// FormatPostIndex formatea direccionamiento post-index (actualiza base después)
func FormatPostIndex(base string, offset int) string {
	return fmt.Sprintf("[%s], #%d", base, offset)
}

// ============ CONSTANTES PARA CONDICIONES ============

const (
	COND_EQ = "eq" // Equal
	COND_NE = "ne" // Not equal
	COND_LT = "lt" // Less than (signed)
	COND_LE = "le" // Less than or equal (signed)
	COND_GT = "gt" // Greater than (signed)
	COND_GE = "ge" // Greater than or equal (signed)
	COND_LO = "lo" // Lower (unsigned)
	COND_LS = "ls" // Lower or same (unsigned)
	COND_HI = "hi" // Higher (unsigned)
	COND_HS = "hs" // Higher or same (unsigned)
)

// ============ HELPERS PARA TIPOS DE DATOS ============

// GetRegisterForType retorna el sufijo de registro apropiado para un tipo
func GetRegisterForType(dataType string) string {
	switch dataType {
	case "int", "bool", "string", "pointer":
		return "x" // 64-bit
	case "float":
		return "d" // Double precision float
	case "char":
		return "w" // 32-bit
	default:
		return "x" // Default a 64-bit
	}
}

// GetLoadStoreInstruction retorna la instrucción de load/store apropiada para un tipo
func GetLoadStoreInstruction(dataType string, isLoad bool) string {
	switch dataType {
	case "int", "bool", "string", "pointer":
		if isLoad {
			return "ldr"
		}
		return "str"
	case "float":
		if isLoad {
			return "ldr" // Para doubles también se usa ldr con registros D
		}
		return "str"
	case "char":
		if isLoad {
			return "ldrb" // Load byte
		}
		return "strb" // Store byte
	default:
		if isLoad {
			return "ldr"
		}
		return "str"
	}
}

// ============ HELPERS PARA CONSTANTES ============

// IsValidImmediate verifica si un valor puede ser usado como inmediato
func IsValidImmediate(value interface{}, instruction string) bool {
	switch v := value.(type) {
	case int:
		// ARM64 tiene diferentes limitaciones según la instrucción
		switch instruction {
		case "add", "sub":
			// ADD/SUB pueden usar inmediatos de 12 bits (0-4095) o 12 bits shifted
			return (v >= 0 && v <= 4095) || (v >= 0 && v <= 0xFFF000 && v%0x1000 == 0)
		case "mov":
			// MOV puede usar inmediatos más complejos, simplificamos a 16 bits
			return v >= -32768 && v <= 65535
		case "cmp":
			// CMP similar a ADD/SUB
			return (v >= 0 && v <= 4095) || (v >= 0 && v <= 0xFFF000 && v%0x1000 == 0)
		default:
			// Por defecto, asumir 12 bits
			return v >= 0 && v <= 4095
		}
	case float64:
		// Los flotantes requieren carga desde memoria generalmente
		return false
	case bool:
		// Los booleanos se pueden representar como 0 o 1
		return true
	default:
		return false
	}
}

// ============ TEMPLATE DE FUNCIÓN ============

// GenerateFunctionProlog genera el prólogo estándar de una función
func GenerateFunctionProlog(ib *InstructionBuilder, functionName string, localStackSize int, preservedRegs []Register) {
	ib.Label(functionName)

	// Guardar frame pointer y link register
	ib.STP("fp", "lr", FormatPreIndex("sp", -16), "save fp and lr")

	// Establecer nuevo frame pointer
	ib.MOV("fp", "sp", "set frame pointer")

	// Reservar espacio para variables locales si es necesario
	if localStackSize > 0 {
		// Redondear a múltiplo de 16 para mantener alineación de stack
		alignedSize := (localStackSize + 15) & ^15
		ib.SUB("sp", "sp", FormatImmediate(alignedSize), fmt.Sprintf("allocate %d bytes for locals", alignedSize))
	}

	// Guardar registros preservados que se van a usar
	if len(preservedRegs) > 0 {
		// Guardar registros de a pares para mantener alineación
		for i := 0; i < len(preservedRegs); i += 2 {
			if i+1 < len(preservedRegs) {
				ib.STP(preservedRegs[i].String(), preservedRegs[i+1].String(),
					FormatPreIndex("sp", -16), "save callee-saved registers")
			} else {
				// Registro impar, guardar solo
				ib.STR(preservedRegs[i].String(), FormatPreIndex("sp", -8), "save callee-saved register")
			}
		}
	}
}

// GenerateFunctionEpilog genera el epílogo estándar de una función
func GenerateFunctionEpilog(ib *InstructionBuilder, preservedRegs []Register) {
	// Restaurar registros preservados
	if len(preservedRegs) > 0 {
		// Restaurar en orden inverso
		for i := len(preservedRegs) - 1; i >= 0; i -= 2 {
			if i > 0 {
				ib.LDP(preservedRegs[i-1].String(), preservedRegs[i].String(),
					FormatPostIndex("sp", 16), "restore callee-saved registers")
				i-- // Saltar el registro ya procesado
			} else {
				// Registro impar
				ib.LDR(preservedRegs[i].String(), FormatPostIndex("sp", 8), "restore callee-saved register")
			}
		}
	}

	// Restaurar stack pointer
	ib.MOV("sp", "fp", "restore stack pointer")

	// Restaurar frame pointer y link register
	ib.LDP("fp", "lr", FormatPostIndex("sp", 16), "restore fp and lr")

	// Retornar
	ib.RET("return to caller")
}
