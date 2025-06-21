// backend/codegen/output/assembler.go
package output

import (
	"fmt"
	"strings"

	"main.go/codegen/arm64"
	"main.go/codegen/intermediate"
)

// ARM64Assembler convierte IR a código ARM64 assembly
type ARM64Assembler struct {
	allocator      *arm64.RegisterAllocator
	callingConv    *arm64.CallingConvention
	ib             *arm64.InstructionBuilder
	labelCounter   int
	stringLiterals map[int]string // ID -> string literal
}

// NewARM64Assembler crea un nuevo ensamblador ARM64
func NewARM64Assembler() *ARM64Assembler {
	return &ARM64Assembler{
		allocator:      arm64.NewRegisterAllocator(),
		callingConv:    arm64.NewCallingConvention(),
		ib:             arm64.NewInstructionBuilder(),
		labelCounter:   0,
		stringLiterals: make(map[int]string),
	}
}

// AssembleProgram convierte un programa IR completo a assembly ARM64
func (asm *ARM64Assembler) AssembleProgram(program *intermediate.IRProgram) (string, error) {
	fmt.Printf("🔧 Generando código ARM64...\n")

	// Reset del estado
	asm.ib = arm64.NewInstructionBuilder()
	asm.stringLiterals = program.StringTable

	// Generar header del assembly
	asm.generateHeader()

	// Generar sección de datos (strings, constantes)
	asm.generateDataSection(program)

	// Generar sección de código
	asm.generateCodeSection(program)

	// Obtener el assembly final
	assembly := asm.ib.GetAssembly()

	fmt.Printf("✅ Código ARM64 generado: %d líneas\n", len(strings.Split(assembly, "\n")))

	return assembly, nil
}

// generateHeader genera el header del archivo assembly
func (asm *ARM64Assembler) generateHeader() {
	asm.ib.Add(".arch", []string{"armv8-a"}, "specify ARM architecture")
	asm.ib.Add(".file", []string{"\"vlancherry_program.s\""}, "source file name")
	asm.ib.Add("", []string{}, "") // línea vacía
}

// generateDataSection genera la sección de datos
func (asm *ARM64Assembler) generateDataSection(program *intermediate.IRProgram) {
	if len(asm.stringLiterals) == 0 {
		return
	}

	asm.ib.Add(".section", []string{".rodata"}, "read-only data section")

	// Generar string literals
	for id, str := range asm.stringLiterals {
		labelName := fmt.Sprintf(".LC%d", id)
		asm.ib.Label(labelName)
		asm.ib.Add(".string", []string{fmt.Sprintf("\"%s\"", str)}, "string literal")
	}

	asm.ib.Add("", []string{}, "") // línea vacía
}

// generateCodeSection genera la sección de código
func (asm *ARM64Assembler) generateCodeSection(program *intermediate.IRProgram) {
	asm.ib.Add(".text", []string{}, "code section")

	// Generar código para cada función
	for _, function := range program.Functions {
		asm.generateFunction(function)
		asm.ib.Add("", []string{}, "") // línea vacía entre funciones
	}

	// Generar punto de entrada si hay función main
	if asm.hasMainFunction(program) {
		asm.generateEntryPoint()
	}
}

// generateFunction genera el código ARM64 para una función IR
func (asm *ARM64Assembler) generateFunction(function *intermediate.IRFunction) {
	fmt.Printf("  🔧 Generando función: %s\n", function.Name)

	// Reset del allocador para cada función
	asm.allocator.Reset()

	// Hacer la función globalmente visible
	asm.ib.Add(".global", []string{function.Name}, fmt.Sprintf("make %s global", function.Name))
	asm.ib.Add(".type", []string{function.Name, "@function"}, "function type")

	// Analizar la función y preparar registros
	asm.analyzeFunction(function)

	// Generar prólogo
	asm.callingConv.GenerateFunctionProlog(asm.ib, function, asm.allocator)

	// Generar código para cada instrucción IR
	for _, instr := range function.Instructions {
		asm.generateInstruction(instr)
	}

	// Generar epílogo si no hay return explícito al final
	lastInstr := function.Instructions[len(function.Instructions)-1]
	if lastInstr.Op != intermediate.IR_RETURN && lastInstr.Op != intermediate.IR_EXIT_FUNCTION {
		asm.callingConv.GenerateFunctionEpilog(asm.ib, asm.allocator)
	}

	// Marcar el final de la función
	asm.ib.Add(".size", []string{function.Name, fmt.Sprintf(". - %s", function.Name)}, "function size")
}

// analyzeFunction analiza una función para preparar la asignación de registros
func (asm *ARM64Assembler) analyzeFunction(function *intermediate.IRFunction) {
	// Reservar registros para parámetros
	for i, param := range function.Parameters {
		if i < 8 {
			paramReg := arm64.GetParameterRegister(i)
			asm.allocator.ReserveRegister(paramReg)
		}
	}

	// Pre-asignar registros para variables que aparecen frecuentemente
	varUsage := make(map[string]int)
	for _, instr := range function.Instructions {
		if instr.Dest != nil && instr.Dest.IsTemp() {
			varUsage[instr.Dest.Name]++
		}
		if instr.Src1 != nil && instr.Src1.IsTemp() {
			varUsage[instr.Src1.Name]++
		}
		if instr.Src2 != nil && instr.Src2.IsTemp() {
			varUsage[instr.Src2.Name]++
		}
	}

	// Asignar registros a variables más usadas primero
	for varName, usage := range varUsage {
		if usage > 2 { // Si se usa más de 2 veces, asignar registro
			asm.allocator.AllocateRegister(varName, "int", true)
		}
	}
}

// generateInstruction genera código ARM64 para una instrucción IR
func (asm *ARM64Assembler) generateInstruction(instr *intermediate.IRInstruction) {
	switch instr.Op {
	case intermediate.IR_LOAD:
		asm.generateLoad(instr)
	case intermediate.IR_STORE:
		asm.generateStore(instr)
	case intermediate.IR_LOAD_IMMEDIATE:
		asm.generateLoadImmediate(instr)
	case intermediate.IR_ADD:
		asm.generateAdd(instr)
	case intermediate.IR_SUB:
		asm.generateSub(instr)
	case intermediate.IR_MULT:
		asm.generateMult(instr)
	case intermediate.IR_DIV:
		asm.generateDiv(instr)
	case intermediate.IR_MOD:
		asm.generateMod(instr)
	case intermediate.IR_CMP_EQ:
		asm.generateCmpEq(instr)
	case intermediate.IR_CMP_NE:
		asm.generateCmpNe(instr)
	case intermediate.IR_CMP_LT:
		asm.generateCmpLt(instr)
	case intermediate.IR_CMP_LE:
		asm.generateCmpLe(instr)
	case intermediate.IR_CMP_GT:
		asm.generateCmpGt(instr)
	case intermediate.IR_CMP_GE:
		asm.generateCmpGe(instr)
	case intermediate.IR_BRANCH:
		asm.generateBranch(instr)
	case intermediate.IR_BRANCH_IF_TRUE:
		asm.generateBranchIfTrue(instr)
	case intermediate.IR_BRANCH_IF_FALSE:
		asm.generateBranchIfFalse(instr)
	case intermediate.IR_LABEL:
		asm.generateLabel(instr)
	case intermediate.IR_CALL:
		asm.generateCall(instr)
	case intermediate.IR_RETURN:
		asm.generateReturn(instr)
	case intermediate.IR_ALLOC_LOCAL:
		asm.generateAllocLocal(instr)
	case intermediate.IR_PRINT:
		asm.generatePrint(instr)
	case intermediate.IR_PRINT_LN:
		asm.generatePrintLn(instr)
	case intermediate.IR_ENTER_FUNCTION:
		// Ya manejado en el prólogo
	case intermediate.IR_EXIT_FUNCTION:
		asm.callingConv.GenerateFunctionEpilog(asm.ib, asm.allocator)
	case intermediate.IR_NOP:
		// No hacer nada
	default:
		asm.ib.Add("// TODO:", []string{fmt.Sprintf("implement %s", instr.Op.String())}, instr.Comment)
	}
}

// generateLoad genera código para cargar de memoria a registro
func (asm *ARM64Assembler) generateLoad(instr *intermediate.IRInstruction) {
	destReg := asm.getOrAllocateRegister(instr.Dest)

	if instr.Src1.IsVariable() {
		// Cargar desde variable local
		if srcReg, exists := asm.allocator.GetRegister(instr.Src1.Name); exists {
			asm.ib.MOV(destReg.String(), srcReg.String(), instr.Comment)
		} else if offset, isSpilled := asm.allocator.IsSpilled(instr.Src1.Name); isSpilled {
			asm.ib.LDR(destReg.String(), arm64.FormatMemory("fp", -offset), instr.Comment)
		}
	}
}

// generateStore genera código para almacenar de registro a memoria
func (asm *ARM64Assembler) generateStore(instr *intermediate.IRInstruction) {
	srcReg := asm.getRegisterForOperand(instr.Src1)

	if instr.Dest.IsVariable() {
		// Almacenar en variable local
		if destReg, exists := asm.allocator.GetRegister(instr.Dest.Name); exists {
			asm.ib.MOV(destReg.String(), srcReg.String(), instr.Comment)
		} else if offset, isSpilled := asm.allocator.IsSpilled(instr.Dest.Name); isSpilled {
			asm.ib.STR(srcReg.String(), arm64.FormatMemory("fp", -offset), instr.Comment)
		} else {
			// Variable no asignada, asignar ahora
			destReg := asm.allocator.AllocateRegister(instr.Dest.Name, instr.Dest.DataType, true)
			if destReg != arm64.INVALID_REGISTER {
				asm.ib.MOV(destReg.String(), srcReg.String(), instr.Comment)
			}
		}
	}
}

// generateLoadImmediate genera código para cargar un valor inmediato
func (asm *ARM64Assembler) generateLoadImmediate(instr *intermediate.IRInstruction) {
	destReg := asm.getOrAllocateRegister(instr.Dest)

	if instr.Src1.IsImmediate() {
		value := instr.Src1.Value
		if arm64.IsValidImmediate(value, "mov") {
			asm.ib.MOV(destReg.String(), arm64.FormatImmediate(value), instr.Comment)
		} else {
			// Valor inmediato muy grande, usar múltiples instrucciones
			asm.generateLargeImmediate(destReg, value, instr.Comment)
		}
	}
}

// generateAdd genera código para suma
func (asm *ARM64Assembler) generateAdd(instr *intermediate.IRInstruction) {
	destReg := asm.getOrAllocateRegister(instr.Dest)
	src1Reg := asm.getRegisterForOperand(instr.Src1)

	if instr.Src2.IsImmediate() && arm64.IsValidImmediate(instr.Src2.Value, "add") {
		asm.ib.ADD(destReg.String(), src1Reg.String(), arm64.FormatImmediate(instr.Src2.Value), instr.Comment)
	} else {
		src2Reg := asm.getRegisterForOperand(instr.Src2)
		asm.ib.ADD(destReg.String(), src1Reg.String(), src2Reg.String(), instr.Comment)
	}
}

// generateSub genera código para resta
func (asm *ARM64Assembler) generateSub(instr *intermediate.IRInstruction) {
	destReg := asm.getOrAllocateRegister(instr.Dest)
	src1Reg := asm.getRegisterForOperand(instr.Src1)

	if instr.Src2.IsImmediate() && arm64.IsValidImmediate(instr.Src2.Value, "sub") {
		asm.ib.SUB(destReg.String(), src1Reg.String(), arm64.FormatImmediate(instr.Src2.Value), instr.Comment)
	} else {
		src2Reg := asm.getRegisterForOperand(instr.Src2)
		asm.ib.SUB(destReg.String(), src1Reg.String(), src2Reg.String(), instr.Comment)
	}
}

// generateMult genera código para multiplicación
func (asm *ARM64Assembler) generateMult(instr *intermediate.IRInstruction) {
	destReg := asm.getOrAllocateRegister(instr.Dest)
	src1Reg := asm.getRegisterForOperand(instr.Src1)
	src2Reg := asm.getRegisterForOperand(instr.Src2)

	asm.ib.MUL(destReg.String(), src1Reg.String(), src2Reg.String(), instr.Comment)
}

// generateDiv genera código para división
func (asm *ARM64Assembler) generateDiv(instr *intermediate.IRInstruction) {
	destReg := asm.getOrAllocateRegister(instr.Dest)
	src1Reg := asm.getRegisterForOperand(instr.Src1)
	src2Reg := asm.getRegisterForOperand(instr.Src2)

	asm.ib.SDIV(destReg.String(), src1Reg.String(), src2Reg.String(), instr.Comment)
}

// generateMod genera código para módulo (ARM64 no tiene instrucción MOD directa)
func (asm *ARM64Assembler) generateMod(instr *intermediate.IRInstruction) {
	destReg := asm.getOrAllocateRegister(instr.Dest)
	src1Reg := asm.getRegisterForOperand(instr.Src1)
	src2Reg := asm.getRegisterForOperand(instr.Src2)

	// a % b = a - (a / b) * b
	tempReg := asm.allocator.AllocateRegister("temp_mod", "int", true)

	asm.ib.SDIV(tempReg.String(), src1Reg.String(), src2Reg.String(), "a / b")
	asm.ib.MUL(tempReg.String(), tempReg.String(), src2Reg.String(), "(a / b) * b")
	asm.ib.SUB(destReg.String(), src1Reg.String(), tempReg.String(), instr.Comment)

	asm.allocator.FreeRegister("temp_mod")
}

// generateCmpEq genera código para comparación de igualdad
func (asm *ARM64Assembler) generateCmpEq(instr *intermediate.IRInstruction) {
	destReg := asm.getOrAllocateRegister(instr.Dest)
	src1Reg := asm.getRegisterForOperand(instr.Src1)
	src2Reg := asm.getRegisterForOperand(instr.Src2)

	asm.ib.CMP(src1Reg.String(), src2Reg.String(), "compare for equality")
	asm.ib.Add("cset", []string{destReg.String(), "eq"}, instr.Comment)
}

// generateCmpNe genera código para comparación de desigualdad
func (asm *ARM64Assembler) generateCmpNe(instr *intermediate.IRInstruction) {
	destReg := asm.getOrAllocateRegister(instr.Dest)
	src1Reg := asm.getRegisterForOperand(instr.Src1)
	src2Reg := asm.getRegisterForOperand(instr.Src2)

	asm.ib.CMP(src1Reg.String(), src2Reg.String(), "compare for inequality")
	asm.ib.Add("cset", []string{destReg.String(), "ne"}, instr.Comment)
}

// generateCmpLt genera código para comparación menor que
func (asm *ARM64Assembler) generateCmpLt(instr *intermediate.IRInstruction) {
	destReg := asm.getOrAllocateRegister(instr.Dest)
	src1Reg := asm.getRegisterForOperand(instr.Src1)
	src2Reg := asm.getRegisterForOperand(instr.Src2)

	asm.ib.CMP(src1Reg.String(), src2Reg.String(), "compare for less than")
	asm.ib.Add("cset", []string{destReg.String(), "lt"}, instr.Comment)
}

// generateCmpLe genera código para comparación menor o igual
func (asm *ARM64Assembler) generateCmpLe(instr *intermediate.IRInstruction) {
	destReg := asm.getOrAllocateRegister(instr.Dest)
	src1Reg := asm.getRegisterForOperand(instr.Src1)
	src2Reg := asm.getRegisterForOperand(instr.Src2)

	asm.ib.CMP(src1Reg.String(), src2Reg.String(), "compare for less than or equal")
	asm.ib.Add("cset", []string{destReg.String(), "le"}, instr.Comment)
}

// generateCmpGt genera código para comparación mayor que
func (asm *ARM64Assembler) generateCmpGt(instr *intermediate.IRInstruction) {
	destReg := asm.getOrAllocateRegister(instr.Dest)
	src1Reg := asm.getRegisterForOperand(instr.Src1)
	src2Reg := asm.getRegisterForOperand(instr.Src2)

	asm.ib.CMP(src1Reg.String(), src2Reg.String(), "compare for greater than")
	asm.ib.Add("cset", []string{destReg.String(), "gt"}, instr.Comment)
}

// generateCmpGe genera código para comparación mayor o igual
func (asm *ARM64Assembler) generateCmpGe(instr *intermediate.IRInstruction) {
	destReg := asm.getOrAllocateRegister(instr.Dest)
	src1Reg := asm.getRegisterForOperand(instr.Src1)
	src2Reg := asm.getRegisterForOperand(instr.Src2)

	asm.ib.CMP(src1Reg.String(), src2Reg.String(), "compare for greater than or equal")
	asm.ib.Add("cset", []string{destReg.String(), "ge"}, instr.Comment)
}

// generateBranch genera código para salto incondicional
func (asm *ARM64Assembler) generateBranch(instr *intermediate.IRInstruction) {
	if instr.Src1 != nil && instr.Src1.Type == intermediate.IR_OPERAND_LABEL {
		asm.ib.B(instr.Src1.Name, instr.Comment)
	}
}

// generateBranchIfTrue genera código para salto condicional si verdadero
func (asm *ARM64Assembler) generateBranchIfTrue(instr *intermediate.IRInstruction) {
	if instr.Src1 != nil && instr.Src2 != nil && instr.Src2.Type == intermediate.IR_OPERAND_LABEL {
		condReg := asm.getRegisterForOperand(instr.Src1)
		asm.ib.CMP(condReg.String(), "#0", "test condition")
		asm.ib.BNE(instr.Src2.Name, instr.Comment)
	}
}

// generateBranchIfFalse genera código para salto condicional si falso
func (asm *ARM64Assembler) generateBranchIfFalse(instr *intermediate.IRInstruction) {
	if instr.Src1 != nil && instr.Src2 != nil && instr.Src2.Type == intermediate.IR_OPERAND_LABEL {
		condReg := asm.getRegisterForOperand(instr.Src1)
		asm.ib.CMP(condReg.String(), "#0", "test condition")
		asm.ib.BEQ(instr.Src2.Name, instr.Comment)
	}
}

// generateLabel genera una etiqueta
func (asm *ARM64Assembler) generateLabel(instr *intermediate.IRInstruction) {
	asm.ib.Label(instr.Label)
}

// generateCall genera código para llamada a función
func (asm *ARM64Assembler) generateCall(instr *intermediate.IRInstruction) {
	if instr.Src1 != nil && instr.Src1.Type == intermediate.IR_OPERAND_LABEL {
		functionName := instr.Src1.Name

		// Para funciones builtin, generar código específico
		if asm.isBuiltinFunction(functionName) {
			asm.generateBuiltinCall(functionName, instr)
		} else {
			// Llamada a función regular
			asm.ib.BL(functionName, instr.Comment)

			// Si hay destino, mover resultado de X0
			if instr.Dest != nil {
				destReg := asm.getOrAllocateRegister(instr.Dest)
				if destReg.String() != "x0" {
					asm.ib.MOV(destReg.String(), "x0", "move return value")
				}
			}
		}
	}
}

// generateReturn genera código para retorno de función
func (asm *ARM64Assembler) generateReturn(instr *intermediate.IRInstruction) {
	// Mover valor de retorno a X0 si es necesario
	if instr.Src1 != nil {
		asm.callingConv.GenerateReturn(asm.ib, instr.Src1, asm.allocator)
	}

	// Generar epílogo y retorno
	asm.callingConv.GenerateFunctionEpilog(asm.ib, asm.allocator)
}

// generateAllocLocal genera código para alocar variable local
func (asm *ARM64Assembler) generateAllocLocal(instr *intermediate.IRInstruction) {
	if instr.Dest != nil {
		// Simplemente asignar un registro o espacio en stack
		asm.allocator.AllocateRegister(instr.Dest.Name, instr.Dest.DataType, false)
	}
}

// generatePrint genera código para función print
func (asm *ARM64Assembler) generatePrint(instr *intermediate.IRInstruction) {
	if instr.Src1 != nil {
		argReg := asm.getRegisterForOperand(instr.Src1)

		// Mover argumento a X0 para la llamada a printf
		if argReg.String() != "x0" {
			asm.ib.MOV("x0", argReg.String(), "move argument to x0")
		}

		// Llamar a función print del runtime
		asm.ib.BL("_vlc_print", instr.Comment)
	}
}

// generatePrintLn genera código para función println
func (asm *ARM64Assembler) generatePrintLn(instr *intermediate.IRInstruction) {
	if instr.Src1 != nil {
		argReg := asm.getRegisterForOperand(instr.Src1)

		// Mover argumento a X0 para la llamada
		if argReg.String() != "x0" {
			asm.ib.MOV("x0", argReg.String(), "move argument to x0")
		}

		// Llamar a función println del runtime
		asm.ib.BL("_vlc_println", instr.Comment)
	}
}

// ============ MÉTODOS AUXILIARES ============

// getOrAllocateRegister obtiene o asigna un registro para un operando
func (asm *ARM64Assembler) getOrAllocateRegister(operand *intermediate.IROperand) arm64.Register {
	if operand == nil {
		return arm64.INVALID_REGISTER
	}

	// Si ya tiene registro asignado, usarlo
	if reg, exists := asm.allocator.GetRegister(operand.Name); exists {
		return reg
	}

	// Asignar nuevo registro
	return asm.allocator.AllocateRegister(operand.Name, operand.DataType, operand.IsTemp())
}

// getRegisterForOperand obtiene el registro o carga el valor para un operando
func (asm *ARM64Assembler) getRegisterForOperand(operand *intermediate.IROperand) arm64.Register {
	if operand == nil {
		return arm64.INVALID_REGISTER
	}

	if operand.IsImmediate() {
		// Cargar inmediato en registro temporal
		tempReg := asm.allocator.AllocateRegister("temp_imm", operand.DataType, true)
		if tempReg != arm64.INVALID_REGISTER {
			if arm64.IsValidImmediate(operand.Value, "mov") {
				asm.ib.MOV(tempReg.String(), arm64.FormatImmediate(operand.Value), "load immediate")
			} else {
				asm.generateLargeImmediate(tempReg, operand.Value, "load large immediate")
			}
		}
		return tempReg
	}

	// Variable o temporal
	if reg, exists := asm.allocator.GetRegister(operand.Name); exists {
		return reg
	}

	// Si está spilled, cargar desde stack
	if offset, isSpilled := asm.allocator.IsSpilled(operand.Name); isSpilled {
		tempReg := asm.allocator.AllocateRegister("temp_load", operand.DataType, true)
		if tempReg != arm64.INVALID_REGISTER {
			asm.ib.LDR(tempReg.String(), arm64.FormatMemory("fp", -offset), "load from stack")
		}
		return tempReg
	}

	return arm64.INVALID_REGISTER
}

// generateLargeImmediate genera código para cargar un inmediato grande
func (asm *ARM64Assembler) generateLargeImmediate(reg arm64.Register, value interface{}, comment string) {
	switch v := value.(type) {
	case int:
		if v >= -65536 && v <= 65535 {
			// Usar MOVZ/MOVN para valores de 16 bits
			if v >= 0 {
				asm.ib.Add("movz", []string{reg.String(), fmt.Sprintf("#%d", v)}, comment)
			} else {
				asm.ib.Add("movn", []string{reg.String(), fmt.Sprintf("#%d", ^v)}, comment)
			}
		} else {
			// Usar múltiples instrucciones para valores más grandes
			low := v & 0xFFFF
			high := (v >> 16) & 0xFFFF
			asm.ib.Add("movz", []string{reg.String(), fmt.Sprintf("#%d", low)}, comment)
			if high != 0 {
				asm.ib.Add("movk", []string{reg.String(), fmt.Sprintf("#%d", high), "lsl #16"}, "load high bits")
			}
		}
	default:
		// Para otros tipos, usar MOV simple
		asm.ib.MOV(reg.String(), arm64.FormatImmediate(value), comment)
	}
}

// isBuiltinFunction verifica si es una función builtin
func (asm *ARM64Assembler) isBuiltinFunction(name string) bool {
	builtins := []string{"print", "println", "atoi", "parseFloat", "TypeOf"}
	for _, builtin := range builtins {
		if name == builtin {
			return true
		}
	}
	return false
}

// generateBuiltinCall genera código para llamadas a funciones builtin
func (asm *ARM64Assembler) generateBuiltinCall(functionName string, instr *intermediate.IRInstruction) {
	switch functionName {
	case "print", "println":
		// Ya manejado en generatePrint/generatePrintLn
	default:
		// Función builtin genérica
		asm.ib.BL("_vlc_"+functionName, instr.Comment)
		if instr.Dest != nil {
			destReg := asm.getOrAllocateRegister(instr.Dest)
			if destReg.String() != "x0" {
				asm.ib.MOV(destReg.String(), "x0", "move return value")
			}
		}
	}
}

// hasMainFunction verifica si el programa tiene función main
func (asm *ARM64Assembler) hasMainFunction(program *intermediate.IRProgram) bool {
	for _, function := range program.Functions {
		if function.Name == "main" {
			return true
		}
	}
	return false
}

// generateEntryPoint genera el punto de entrada del programa
func (asm *ARM64Assembler) generateEntryPoint() {
	// Punto de entrada estándar para programas ARM64
	asm.ib.Add(".global", []string{"_start"}, "program entry point")
	asm.ib.Label("_start")

	// Llamar a main
	asm.ib.BL("main", "call main function")

	// Terminar programa con exit syscall
	asm.ib.MOV("x8", "#93", "exit syscall number") // syscall __NR_exit
	asm.ib.MOV("x0", "x0", "exit code from main")  // código de salida
	asm.ib.Add("svc", []string{"#0"}, "system call")
}
