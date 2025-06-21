// backend/codegen/arm64/registers.go
package arm64

import (
	"fmt"
)

// Register representa un registro ARM64
type Register int

const (
	// Registros de propósito general (64-bit)
	X0 Register = iota
	X1
	X2
	X3
	X4
	X5
	X6
	X7
	X8  // Registro de resultado indirecto
	X9  // Temporary
	X10 // Temporary
	X11 // Temporary
	X12 // Temporary
	X13 // Temporary
	X14 // Temporary
	X15 // Temporary
	X16 // IP0 - Intra-procedure call temporary
	X17 // IP1 - Intra-procedure call temporary
	X18 // Platform register
	X19 // Callee-saved
	X20 // Callee-saved
	X21 // Callee-saved
	X22 // Callee-saved
	X23 // Callee-saved
	X24 // Callee-saved
	X25 // Callee-saved
	X26 // Callee-saved
	X27 // Callee-saved
	X28 // Callee-saved
	X29 // Frame pointer (FP)
	X30 // Link register (LR)
	SP  // Stack pointer

	// Registros de 32-bit (parte baja de los X registers)
	W0
	W1
	W2
	W3
	W4
	W5
	W6
	W7
	W8
	W9
	W10
	W11
	W12
	W13
	W14
	W15
	W16
	W17
	W18
	W19
	W20
	W21
	W22
	W23
	W24
	W25
	W26
	W27
	W28
	W29
	W30

	// Registros especiales
	FP = X29 // Frame pointer alias
	LR = X30 // Link register alias

	// Registros de punto flotante
	D0
	D1
	D2
	D3
	D4
	D5
	D6
	D7
	D8
	D9
	D10
	D11
	D12
	D13
	D14
	D15
	D16
	D17
	D18
	D19
	D20
	D21
	D22
	D23
	D24
	D25
	D26
	D27
	D28
	D29
	D30
	D31

	INVALID_REGISTER = -1
)

// String convierte un registro a su representación en assembly
func (r Register) String() string {
	switch {
	case r >= X0 && r <= X28:
		return fmt.Sprintf("x%d", int(r-X0))
	case r == X29:
		return "fp"
	case r == X30:
		return "lr"
	case r == SP:
		return "sp"
	case r >= W0 && r <= W30:
		return fmt.Sprintf("w%d", int(r-W0))
	case r >= D0 && r <= D31:
		return fmt.Sprintf("d%d", int(r-D0))
	default:
		return "INVALID"
	}
}

// Is64Bit verifica si es un registro de 64 bits
func (r Register) Is64Bit() bool {
	return (r >= X0 && r <= SP) || (r >= D0 && r <= D31)
}

// Is32Bit verifica si es un registro de 32 bits
func (r Register) Is32Bit() bool {
	return r >= W0 && r <= W30
}

// IsFloat verifica si es un registro de punto flotante
func (r Register) IsFloat() bool {
	return r >= D0 && r <= D31
}

// ToX convierte un registro W a su equivalente X
func (r Register) ToX() Register {
	if r >= W0 && r <= W30 {
		return Register(int(r-W0) + int(X0))
	}
	return r
}

// ToW convierte un registro X a su equivalente W
func (r Register) ToW() Register {
	if r >= X0 && r <= X30 {
		return Register(int(r-X0) + int(W0))
	}
	return r
}

// RegisterAllocator gestiona la asignación de registros
type RegisterAllocator struct {
	// Registros disponibles para asignación temporal
	availableTemps []Register

	// Registros disponibles para preservar valores (callee-saved)
	availableSaved []Register

	// Mapeo de variables/temporales IR a registros físicos
	allocation map[string]Register

	// Registros actualmente en uso
	used map[Register]bool

	// Variables que fueron spilled al stack
	spilled map[string]int // variable -> offset en stack

	// Offset actual del stack para spilling
	stackOffset int
}

// NewRegisterAllocator crea un nuevo allocator de registros
func NewRegisterAllocator() *RegisterAllocator {
	return &RegisterAllocator{
		// Registros temporales (caller-saved)
		availableTemps: []Register{X0, X1, X2, X3, X4, X5, X6, X7, X9, X10, X11, X12, X13, X14, X15},

		// Registros preservados (callee-saved)
		availableSaved: []Register{X19, X20, X21, X22, X23, X24, X25, X26, X27, X28},

		allocation:  make(map[string]Register),
		used:        make(map[Register]bool),
		spilled:     make(map[string]int),
		stackOffset: 0,
	}
}

// AllocateRegister asigna un registro para una variable/temporal
func (ra *RegisterAllocator) AllocateRegister(name string, dataType string, preferTemp bool) Register {
	// Si ya está asignado, retornar el mismo registro
	if reg, exists := ra.allocation[name]; exists {
		return reg
	}

	// Determinar qué pool de registros usar
	var candidates []Register
	if preferTemp {
		candidates = ra.availableTemps
	} else {
		candidates = ra.availableSaved
	}

	// Buscar un registro libre
	for _, reg := range candidates {
		if !ra.used[reg] {
			ra.allocation[name] = reg
			ra.used[reg] = true
			return reg
		}
	}

	// Si no hay registros libres en el pool preferido, intentar el otro
	otherCandidates := ra.availableSaved
	if !preferTemp {
		otherCandidates = ra.availableTemps
	}

	for _, reg := range otherCandidates {
		if !ra.used[reg] {
			ra.allocation[name] = reg
			ra.used[reg] = true
			return reg
		}
	}

	// Si no hay registros disponibles, hacer spill
	return ra.spillVariable(name)
}

// spillVariable mueve una variable al stack cuando no hay registros disponibles
func (ra *RegisterAllocator) spillVariable(name string) Register {
	// Encontrar una variable para hacer spill (estrategia simple: la primera encontrada)
	var victimName string
	var victimReg Register

	for varName, reg := range ra.allocation {
		if varName != name { // No hacer spill de la variable que estamos asignando
			victimName = varName
			victimReg = reg
			break
		}
	}

	if victimName == "" {
		// No hay variables para hacer spill, usar stack directamente
		ra.stackOffset += 8 // 8 bytes por variable
		ra.spilled[name] = ra.stackOffset
		return INVALID_REGISTER
	}

	// Mover la víctima al stack
	ra.stackOffset += 8
	ra.spilled[victimName] = ra.stackOffset
	delete(ra.allocation, victimName)

	// Asignar el registro liberado a la nueva variable
	ra.allocation[name] = victimReg
	return victimReg
}

// GetRegister obtiene el registro asignado a una variable
func (ra *RegisterAllocator) GetRegister(name string) (Register, bool) {
	reg, exists := ra.allocation[name]
	return reg, exists
}

// IsSpilled verifica si una variable está en el stack
func (ra *RegisterAllocator) IsSpilled(name string) (int, bool) {
	offset, exists := ra.spilled[name]
	return offset, exists
}

// FreeRegister libera un registro para reutilización
func (ra *RegisterAllocator) FreeRegister(name string) {
	if reg, exists := ra.allocation[name]; exists {
		delete(ra.allocation, name)
		ra.used[reg] = false
	}
}

// GetStackSize retorna el tamaño total usado en el stack para spilling
func (ra *RegisterAllocator) GetStackSize() int {
	return ra.stackOffset
}

// PreserveCalleeSaved retorna la lista de registros callee-saved que necesitan preservarse
func (ra *RegisterAllocator) PreserveCalleeSaved() []Register {
	var toPreserve []Register

	for _, reg := range ra.availableSaved {
		if ra.used[reg] {
			toPreserve = append(toPreserve, reg)
		}
	}

	return toPreserve
}

// Reset reinicia el allocator para una nueva función
func (ra *RegisterAllocator) Reset() {
	ra.allocation = make(map[string]Register)
	ra.used = make(map[Register]bool)
	ra.spilled = make(map[string]int)
	ra.stackOffset = 0
}

// ReserveRegister marca un registro como usado (para parámetros de función, etc.)
func (ra *RegisterAllocator) ReserveRegister(reg Register) {
	ra.used[reg] = true
}

// GetParameterRegister retorna el registro usado para el parámetro n (siguiendo AAPCS64)
func GetParameterRegister(paramIndex int) Register {
	if paramIndex < 8 {
		return Register(int(X0) + paramIndex)
	}
	// Parámetros 8+ van en el stack
	return INVALID_REGISTER
}

// GetReturnRegister retorna el registro usado para valores de retorno
func GetReturnRegister(dataType string) Register {
	// Por simplicidad, siempre usar X0 para retornos
	return X0
}

// IsCallerSaved verifica si un registro es caller-saved (temporal)
func IsCallerSaved(reg Register) bool {
	callerSaved := []Register{X0, X1, X2, X3, X4, X5, X6, X7, X8, X9, X10, X11, X12, X13, X14, X15, X16, X17, X18}

	for _, r := range callerSaved {
		if r == reg {
			return true
		}
	}
	return false
}

// IsCalleeSaved verifica si un registro es callee-saved (preservado)
func IsCalleeSaved(reg Register) bool {
	calleeSaved := []Register{X19, X20, X21, X22, X23, X24, X25, X26, X27, X28, FP, LR}

	for _, r := range calleeSaved {
		if r == reg {
			return true
		}
	}
	return false
}
