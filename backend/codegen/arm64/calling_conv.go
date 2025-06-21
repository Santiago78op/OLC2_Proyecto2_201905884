// backend/codegen/arm64/calling_conv.go
package arm64

import (
	"fmt"

	"main.go/codegen/intermediate"
)

// CallingConvention implementa la convención de llamada AAPCS64 para ARM64
type CallingConvention struct {
	// Registros para parámetros enteros/punteros
	intParamRegisters []Register

	// Registros para parámetros de punto flotante
	floatParamRegisters []Register

	// Registro para valor de retorno
	returnRegister Register

	// Registros temporales (caller-saved)
	callerSavedRegs []Register

	// Registros preservados (callee-saved)
	calleeSavedRegs []Register
}

// NewCallingConvention crea una nueva instancia de la convención de llamada
func NewCallingConvention() *CallingConvention {
	return &CallingConvention{
		intParamRegisters:   []Register{X0, X1, X2, X3, X4, X5, X6, X7},
		floatParamRegisters: []Register{D0, D1, D2, D3, D4, D5, D6, D7},
		returnRegister:      X0,
		callerSavedRegs:     []Register{X0, X1, X2, X3, X4, X5, X6, X7, X8, X9, X10, X11, X12, X13, X14, X15, X16, X17, X18},
		calleeSavedRegs:     []Register{X19, X20, X21, X22, X23, X24, X25, X26, X27, X28},
	}
}

// ParameterAssignment representa la asignación de un parámetro
type ParameterAssignment struct {
	Operand     *intermediate.IROperand
	Register    Register // INVALID_REGISTER si va en stack
	StackOffset int      // Offset en stack si no cabe en registro
	IsFloat     bool     // True si es parámetro de punto flotante
}

// FunctionCallInfo contiene información sobre una llamada a función
type FunctionCallInfo struct {
	FunctionName string
	Parameters   []*ParameterAssignment
	ReturnType   string
	StackSpace   int // Espacio total necesario en stack para parámetros
}

// AnalyzeParameters analiza los parámetros de una función y determina su asignación
func (cc *CallingConvention) AnalyzeParameters(params []*intermediate.IROperand) *FunctionCallInfo {
	info := &FunctionCallInfo{
		Parameters: make([]*ParameterAssignment, 0),
		StackSpace: 0,
	}

	intRegIndex := 0
	floatRegIndex := 0
	stackOffset := 0

	for _, param := range params {
		assignment := &ParameterAssignment{
			Operand:     param,
			Register:    INVALID_REGISTER,
			StackOffset: -1,
			IsFloat:     cc.isFloatType(param.DataType),
		}

		if assignment.IsFloat {
			// Parámetro de punto flotante
			if floatRegIndex < len(cc.floatParamRegisters) {
				assignment.Register = cc.floatParamRegisters[floatRegIndex]
				floatRegIndex++
			} else {
				// No hay más registros float, usar stack
				assignment.StackOffset = stackOffset
				stackOffset += 8
			}
		} else {
			// Parámetro entero/puntero
			if intRegIndex < len(cc.intParamRegisters) {
				assignment.Register = cc.intParamRegisters[intRegIndex]
				intRegIndex++
			} else {
				// No hay más registros int, usar stack
				assignment.StackOffset = stackOffset
				stackOffset += 8
			}
		}

		info.Parameters = append(info.Parameters, assignment)
	}

	// Redondear stack space a múltiplo de 16 para alineación
	info.StackSpace = (stackOffset + 15) & ^15

	return info
}

// GenerateCall genera el código ARM64 para una llamada a función
func (cc *CallingConvention) GenerateCall(ib *InstructionBuilder, callInfo *FunctionCallInfo, allocator *RegisterAllocator) {
	// 1. Guardar registros caller-saved que están en uso
	cc.saveCallerSavedRegisters(ib, allocator)

	// 2. Reservar espacio en stack para parámetros si es necesario
	if callInfo.StackSpace > 0 {
		ib.SUB("sp", "sp", FormatImmediate(callInfo.StackSpace), "allocate space for parameters")
	}

	// 3. Cargar parámetros en registros o stack
	for _, param := range callInfo.Parameters {
		cc.loadParameter(ib, param, allocator)
	}

	// 4. Realizar la llamada
	ib.BL(callInfo.FunctionName, fmt.Sprintf("call %s", callInfo.FunctionName))

	// 5. Limpiar stack de parámetros si es necesario
	if callInfo.StackSpace > 0 {
		ib.ADD("sp", "sp", FormatImmediate(callInfo.StackSpace), "clean up parameter space")
	}

	// 6. Restaurar registros caller-saved
	cc.restoreCallerSavedRegisters(ib, allocator)
}

// loadParameter carga un parámetro en su registro o posición de stack asignada
func (cc *CallingConvention) loadParameter(ib *InstructionBuilder, param *ParameterAssignment, allocator *RegisterAllocator) {
	// Obtener el valor del parámetro
	var sourceOperand string

	if param.Operand.IsImmediate() {
		// Valor inmediato
		if param.Register != INVALID_REGISTER {
			// Cargar inmediato en registro
			if param.IsFloat {
				// Para flotantes, necesitamos cargar desde memoria generalmente
				// Por simplicidad, asumimos que ya está en memoria o usamos una aproximación
				ib.MOV(param.Register.String(), FormatImmediate(param.Operand.Value),
					fmt.Sprintf("load immediate %v", param.Operand.Value))
			} else {
				ib.MOV(param.Register.String(), FormatImmediate(param.Operand.Value),
					fmt.Sprintf("load immediate %v", param.Operand.Value))
			}
		} else {
			// Almacenar inmediato en stack
			tempReg := "x9" // Usar registro temporal
			ib.MOV(tempReg, FormatImmediate(param.Operand.Value), "load immediate to temp")
			ib.STR(tempReg, FormatMemory("sp", param.StackOffset), "store parameter to stack")
		}
	} else {
		// Variable o temporal
		if reg, exists := allocator.GetRegister(param.Operand.Name); exists {
			// Está en registro
			if param.Register != INVALID_REGISTER {
				// Mover de registro a registro de parámetro
				ib.MOV(param.Register.String(), reg.String(),
					fmt.Sprintf("move %s to parameter register", param.Operand.Name))
			} else {
				// Mover de registro a stack
				ib.STR(reg.String(), FormatMemory("sp", param.StackOffset),
					fmt.Sprintf("store %s to stack", param.Operand.Name))
			}
		} else if offset, isSpilled := allocator.IsSpilled(param.Operand.Name); isSpilled {
			// Está en stack (spilled)
			tempReg := "x9" // Usar registro temporal
			ib.LDR(tempReg, FormatMemory("fp", -offset), fmt.Sprintf("load spilled %s", param.Operand.Name))

			if param.Register != INVALID_REGISTER {
				ib.MOV(param.Register.String(), tempReg,
					fmt.Sprintf("move %s to parameter register", param.Operand.Name))
			} else {
				ib.STR(tempReg, FormatMemory("sp", param.StackOffset),
					fmt.Sprintf("store %s to stack", param.Operand.Name))
			}
		}
	}
}

// saveCallerSavedRegisters guarda los registros caller-saved que están en uso
func (cc *CallingConvention) saveCallerSavedRegisters(ib *InstructionBuilder, allocator *RegisterAllocator) {
	// Implementación simplificada: guardar todos los registros caller-saved
	// En una implementación real, solo guardaríamos los que están en uso

	ib.STP("x0", "x1", FormatPreIndex("sp", -16), "save caller-saved registers")
	ib.STP("x2", "x3", FormatPreIndex("sp", -16), "save caller-saved registers")
	ib.STP("x4", "x5", FormatPreIndex("sp", -16), "save caller-saved registers")
	ib.STP("x6", "x7", FormatPreIndex("sp", -16), "save caller-saved registers")
	ib.STP("x8", "x9", FormatPreIndex("sp", -16), "save caller-saved registers")
}

// restoreCallerSavedRegisters restaura los registros caller-saved
func (cc *CallingConvention) restoreCallerSavedRegisters(ib *InstructionBuilder, allocator *RegisterAllocator) {
	// Restaurar en orden inverso
	ib.LDP("x8", "x9", FormatPostIndex("sp", 16), "restore caller-saved registers")
	ib.LDP("x6", "x7", FormatPostIndex("sp", 16), "restore caller-saved registers")
	ib.LDP("x4", "x5", FormatPostIndex("sp", 16), "restore caller-saved registers")
	ib.LDP("x2", "x3", FormatPostIndex("sp", 16), "restore caller-saved registers")
	ib.LDP("x0", "x1", FormatPostIndex("sp", 16), "restore caller-saved registers")
}

// GenerateReturn genera el código para retornar un valor
func (cc *CallingConvention) GenerateReturn(ib *InstructionBuilder, returnValue *intermediate.IROperand, allocator *RegisterAllocator) {
	if returnValue == nil {
		// No hay valor de retorno
		return
	}

	targetReg := cc.returnRegister.String()

	if returnValue.IsImmediate() {
		// Valor inmediato
		ib.MOV(targetReg, FormatImmediate(returnValue.Value), "load return value")
	} else {
		// Variable o temporal
		if reg, exists := allocator.GetRegister(returnValue.Name); exists {
			// Está en registro
			if reg != cc.returnRegister {
				ib.MOV(targetReg, reg.String(), "move return value to x0")
			}
		} else if offset, isSpilled := allocator.IsSpilled(returnValue.Name); isSpilled {
			// Está en stack
			ib.LDR(targetReg, FormatMemory("fp", -offset), "load return value from stack")
		}
	}
}

// GetReturnRegister retorna el registro usado para valores de retorno
func (cc *CallingConvention) GetReturnRegister() Register {
	return cc.returnRegister
}

// isFloatType verifica si un tipo de dato es de punto flotante
func (cc *CallingConvention) isFloatType(dataType string) bool {
	return dataType == "float" || dataType == "double"
}

// CalculateStackFrame calcula el tamaño del stack frame para una función
func (cc *CallingConvention) CalculateStackFrame(localVars []*intermediate.IROperand, spilledVars int) int {
	// Tamaño base para variables locales
	localSize := len(localVars) * 8

	// Espacio para variables spilled
	spillSize := spilledVars * 8

	// Espacio para registros preservados (estimación)
	preservedSize := len(cc.calleeSavedRegs) * 8

	// Total redondeado a múltiplo de 16
	total := localSize + spillSize + preservedSize
	return (total + 15) & ^15
}

// GenerateFunctionProlog genera el prólogo de función con convención de llamada
func (cc *CallingConvention) GenerateFunctionProlog(ib *InstructionBuilder, function *intermediate.IRFunction, allocator *RegisterAllocator) {
	// Calcular tamaño del stack frame
	stackSize := cc.CalculateStackFrame(function.LocalVars, allocator.GetStackSize())

	// Determinar qué registros preservados necesitamos guardar
	preservedRegs := allocator.PreserveCalleeSaved()

	// Generar prólogo estándar
	GenerateFunctionProlog(ib, function.Name, stackSize, preservedRegs)

	// Asignar parámetros a variables locales
	for i, param := range function.Parameters {
		if i < len(cc.intParamRegisters) {
			// Parámetro está en registro
			paramReg := cc.intParamRegisters[i]

			// Asignar variable local al mismo registro si es posible
			if allocatedReg := allocator.AllocateRegister(param.Name, param.DataType, true); allocatedReg != INVALID_REGISTER {
				if allocatedReg != paramReg {
					ib.MOV(allocatedReg.String(), paramReg.String(),
						fmt.Sprintf("copy parameter %s", param.Name))
				}
			} else {
				// Parámetro debe ir a stack
				if offset, _ := allocator.IsSpilled(param.Name); offset > 0 {
					ib.STR(paramReg.String(), FormatMemory("fp", -offset),
						fmt.Sprintf("store parameter %s to stack", param.Name))
				}
			}
		} else {
			// Parámetro está en stack del caller
			// Calculamos su offset: parámetros extra están después del link record
			paramOffset := 16 + (i-len(cc.intParamRegisters))*8

			if allocatedReg := allocator.AllocateRegister(param.Name, param.DataType, true); allocatedReg != INVALID_REGISTER {
				// Cargar de stack del caller a registro local
				ib.LDR(allocatedReg.String(), FormatMemory("fp", paramOffset),
					fmt.Sprintf("load parameter %s from caller stack", param.Name))
			} else {
				// Copiar de stack del caller a stack local
				tempReg := "x9"
				ib.LDR(tempReg, FormatMemory("fp", paramOffset),
					fmt.Sprintf("load parameter %s from caller stack", param.Name))
				if offset, _ := allocator.IsSpilled(param.Name); offset > 0 {
					ib.STR(tempReg, FormatMemory("fp", -offset),
						fmt.Sprintf("store parameter %s to local stack", param.Name))
				}
			}
		}
	}
}

// GenerateFunctionEpilog genera el epílogo de función con convención de llamada
func (cc *CallingConvention) GenerateFunctionEpilog(ib *InstructionBuilder, allocator *RegisterAllocator) {
	// Determinar qué registros preservados fueron guardados
	preservedRegs := allocator.PreserveCalleeSaved()

	// Generar epílogo estándar
	GenerateFunctionEpilog(ib, preservedRegs)
}

// IsCallerSaved verifica si un registro debe ser guardado por el caller
func (cc *CallingConvention) IsCallerSaved(reg Register) bool {
	for _, r := range cc.callerSavedRegs {
		if r == reg {
			return true
		}
	}
	return false
}

// IsCalleeSaved verifica si un registro debe ser guardado por el callee
func (cc *CallingConvention) IsCalleeSaved(reg Register) bool {
	for _, r := range cc.calleeSavedRegs {
		if r == reg {
			return true
		}
	}
	return false
}

// ValidateFunction valida que una función cumpla con la convención de llamada
func (cc *CallingConvention) ValidateFunction(function *intermediate.IRFunction) []string {
	var warnings []string

	// Verificar número de parámetros
	if len(function.Parameters) > len(cc.intParamRegisters)+len(cc.floatParamRegisters) {
		warnings = append(warnings, fmt.Sprintf("Function %s has %d parameters, some will be passed on stack",
			function.Name, len(function.Parameters)))
	}

	// Verificar tipos de parámetros
	for i, param := range function.Parameters {
		if cc.isFloatType(param.DataType) && i >= len(cc.floatParamRegisters) {
			warnings = append(warnings, fmt.Sprintf("Float parameter %s will be passed on stack", param.Name))
		}
	}

	return warnings
}
