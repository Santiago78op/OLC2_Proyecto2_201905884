// backend/codegen/arm64/optimizer.go
package arm64

import (
	"fmt"
	"strconv"
	"strings"
)

// ARM64Optimizer realiza optimizaciones básicas en código ARM64
type ARM64Optimizer struct {
	instructions []*ARM64Instruction
	changes      int
}

// NewARM64Optimizer crea un nuevo optimizador ARM64
func NewARM64Optimizer() *ARM64Optimizer {
	return &ARM64Optimizer{
		changes: 0,
	}
}

// Optimize aplica optimizaciones al código ARM64
func (opt *ARM64Optimizer) Optimize(instructions []*ARM64Instruction) []*ARM64Instruction {
	opt.instructions = instructions
	opt.changes = 0

	fmt.Printf("🔧 Aplicando optimizaciones ARM64...\n")

	totalChanges := 0
	iteration := 0

	// Aplicar optimizaciones iterativamente
	for {
		iteration++
		opt.changes = 0

		fmt.Printf("  📊 Iteración de optimización ARM64 %d\n", iteration)

		// Aplicar diferentes tipos de optimizaciones
		opt.instructions = opt.removeRedundantMoves()
		opt.instructions = opt.combineLoadStore()
		opt.instructions = opt.optimizeImmediates()
		opt.instructions = opt.removeDeadCode()
		opt.instructions = opt.peepholeOptimizations()

		totalChanges += opt.changes
		fmt.Printf("    ✅ %d optimizaciones ARM64 aplicadas en esta iteración\n", opt.changes)

		// Si no hubo cambios, terminar
		if opt.changes == 0 {
			break
		}

		// Límite de seguridad
		if iteration > 5 {
			fmt.Printf("  ⚠️ Límite de iteraciones ARM64 alcanzado\n")
			break
		}
	}

	fmt.Printf("✅ Optimización ARM64 completada: %d optimizaciones totales\n", totalChanges)

	return opt.instructions
}

// removeRedundantMoves elimina movimientos redundantes
func (opt *ARM64Optimizer) removeRedundantMoves() []*ARM64Instruction {
	var optimized []*ARM64Instruction

	for i, instr := range opt.instructions {
		keep := true

		// Eliminar mov x, x (movimiento a sí mismo)
		if instr.Opcode == "mov" && len(instr.Operands) == 2 {
			if instr.Operands[0] == instr.Operands[1] {
				keep = false
				opt.changes++
				fmt.Printf("    🗑️ Eliminado mov redundante: %s\n", instr.String())
			}
		}

		// Eliminar movimientos redundantes consecutivos
		if keep && i > 0 && instr.Opcode == "mov" && len(instr.Operands) == 2 {
			prevInstr := opt.instructions[i-1]
			if prevInstr.Opcode == "mov" && len(prevInstr.Operands) == 2 {
				// mov a, b seguido de mov b, a
				if instr.Operands[0] == prevInstr.Operands[1] &&
					instr.Operands[1] == prevInstr.Operands[0] {
					keep = false
					opt.changes++
					fmt.Printf("    🗑️ Eliminado mov redundante consecutivo: %s\n", instr.String())
				}
			}
		}

		if keep {
			optimized = append(optimized, instr)
		}
	}

	return optimized
}

// combineLoadStore combina operaciones de load/store consecutivas
func (opt *ARM64Optimizer) combineLoadStore() []*ARM64Instruction {
	var optimized []*ARM64Instruction

	for i := 0; i < len(opt.instructions); i++ {
		instr := opt.instructions[i]
		combined := false

		// Buscar oportunidades de combinación con la siguiente instrucción
		if i+1 < len(opt.instructions) {
			nextInstr := opt.instructions[i+1]

			// Combinar str + str consecutivos en stp
			if opt.canCombineStores(instr, nextInstr) {
				stpInstr := opt.createStoreParInstr(instr, nextInstr)
				if stpInstr != nil {
					optimized = append(optimized, stpInstr)
					i++ // Saltar la siguiente instrucción
					combined = true
					opt.changes++
					fmt.Printf("    🔗 Combinado str+str en stp: %s\n", stpInstr.String())
				}
			}

			// Combinar ldr + ldr consecutivos en ldp
			if !combined && opt.canCombineLoads(instr, nextInstr) {
				ldpInstr := opt.createLoadPairInstr(instr, nextInstr)
				if ldpInstr != nil {
					optimized = append(optimized, ldpInstr)
					i++ // Saltar la siguiente instrucción
					combined = true
					opt.changes++
					fmt.Printf("    🔗 Combinado ldr+ldr en ldp: %s\n", ldpInstr.String())
				}
			}
		}

		if !combined {
			optimized = append(optimized, instr)
		}
	}

	return optimized
}

// canCombineStores verifica si dos stores se pueden combinar
func (opt *ARM64Optimizer) canCombineStores(instr1, instr2 *ARM64Instruction) bool {
	if instr1.Opcode != "str" || instr2.Opcode != "str" {
		return false
	}

	if len(instr1.Operands) != 2 || len(instr2.Operands) != 2 {
		return false
	}

	// Verificar que las direcciones sean consecutivas
	return opt.areConsecutiveAddresses(instr1.Operands[1], instr2.Operands[1])
}

// canCombineLoads verifica si dos loads se pueden combinar
func (opt *ARM64Optimizer) canCombineLoads(instr1, instr2 *ARM64Instruction) bool {
	if instr1.Opcode != "ldr" || instr2.Opcode != "ldr" {
		return false
	}

	if len(instr1.Operands) != 2 || len(instr2.Operands) != 2 {
		return false
	}

	// Verificar que las direcciones sean consecutivas
	return opt.areConsecutiveAddresses(instr1.Operands[1], instr2.Operands[1])
}

// areConsecutiveAddresses verifica si dos direcciones son consecutivas
func (opt *ARM64Optimizer) areConsecutiveAddresses(addr1, addr2 string) bool {
	// Implementación simplificada: buscar patrones como [fp, #8] y [fp, #16]
	if strings.Contains(addr1, "[") && strings.Contains(addr2, "[") {
		// Extraer base y offset de ambas direcciones
		base1, offset1 := opt.parseAddress(addr1)
		base2, offset2 := opt.parseAddress(addr2)

		// Mismo registro base y offsets consecutivos (diferencia de 8)
		return base1 == base2 && (offset2-offset1 == 8 || offset1-offset2 == 8)
	}

	return false
}

// parseAddress extrae la base y offset de una dirección
func (opt *ARM64Optimizer) parseAddress(addr string) (string, int) {
	// Parsear direcciones como [fp, #8] o [sp, #-16]
	addr = strings.Trim(addr, "[]")
	parts := strings.Split(addr, ",")

	if len(parts) != 2 {
		return "", 0
	}

	base := strings.TrimSpace(parts[0])
	offsetStr := strings.TrimSpace(parts[1])
	offsetStr = strings.TrimPrefix(offsetStr, "#")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return base, 0
	}

	return base, offset
}

// createStoreParInstr crea una instrucción stp a partir de dos str
func (opt *ARM64Optimizer) createStoreParInstr(instr1, instr2 *ARM64Instruction) *ARM64Instruction {
	// Determinar el orden correcto basado en los offsets
	addr1 := instr1.Operands[1]
	addr2 := instr2.Operands[1]

	_, offset1 := opt.parseAddress(addr1)
	_, offset2 := opt.parseAddress(addr2)

	var reg1, reg2, addr string
	if offset1 < offset2 {
		reg1 = instr1.Operands[0]
		reg2 = instr2.Operands[0]
		addr = opt.adjustAddressForPair(addr1)
	} else {
		reg1 = instr2.Operands[0]
		reg2 = instr1.Operands[0]
		addr = opt.adjustAddressForPair(addr2)
	}

	return &ARM64Instruction{
		Opcode:   "stp",
		Operands: []string{reg1, reg2, addr},
		Comment:  "combined store pair",
	}
}

// createLoadPairInstr crea una instrucción ldp a partir de dos ldr
func (opt *ARM64Optimizer) createLoadPairInstr(instr1, instr2 *ARM64Instruction) *ARM64Instruction {
	// Similar a createStoreParInstr pero para loads
	addr1 := instr1.Operands[1]
	addr2 := instr2.Operands[1]

	_, offset1 := opt.parseAddress(addr1)
	_, offset2 := opt.parseAddress(addr2)

	var reg1, reg2, addr string
	if offset1 < offset2 {
		reg1 = instr1.Operands[0]
		reg2 = instr2.Operands[0]
		addr = opt.adjustAddressForPair(addr1)
	} else {
		reg1 = instr2.Operands[0]
		reg2 = instr1.Operands[0]
		addr = opt.adjustAddressForPair(addr2)
	}

	return &ARM64Instruction{
		Opcode:   "ldp",
		Operands: []string{reg1, reg2, addr},
		Comment:  "combined load pair",
	}
}

// adjustAddressForPair ajusta una dirección para usar con instrucciones pair
func (opt *ARM64Optimizer) adjustAddressForPair(addr string) string {
	// Para simplificar, retornar la dirección base
	// En una implementación real, ajustaríamos el offset apropiadamente
	return addr
}

// optimizeImmediates optimiza el uso de valores inmediatos
func (opt *ARM64Optimizer) optimizeImmediates() []*ARM64Instruction {
	var optimized []*ARM64Instruction

	for _, instr := range opt.instructions {
		optimizedInstr := opt.optimizeInstructionImmediates(instr)
		optimized = append(optimized, optimizedInstr)
	}

	return optimized
}

// optimizeInstructionImmediates optimiza los inmediatos en una instrucción
func (opt *ARM64Optimizer) optimizeInstructionImmediates(instr *ARM64Instruction) *ARM64Instruction {
	// Clonar la instrucción para no modificar la original
	newInstr := &ARM64Instruction{
		Opcode:   instr.Opcode,
		Operands: make([]string, len(instr.Operands)),
		Comment:  instr.Comment,
		Label:    instr.Label,
	}
	copy(newInstr.Operands, instr.Operands)

	// Optimizar operaciones con 0
	if instr.Opcode == "add" && len(instr.Operands) == 3 {
		if instr.Operands[2] == "#0" {
			// add x, y, #0 -> mov x, y
			newInstr.Opcode = "mov"
			newInstr.Operands = []string{instr.Operands[0], instr.Operands[1]}
			newInstr.Comment = "optimized: add with 0"
			opt.changes++
			fmt.Printf("    ⚡ Optimizado add con 0: %s\n", newInstr.String())
		}
	}

	// Optimizar operaciones con 1
	if instr.Opcode == "mul" && len(instr.Operands) == 3 {
		if instr.Operands[2] == "#1" {
			// mul x, y, #1 -> mov x, y
			newInstr.Opcode = "mov"
			newInstr.Operands = []string{instr.Operands[0], instr.Operands[1]}
			newInstr.Comment = "optimized: mul with 1"
			opt.changes++
			fmt.Printf("    ⚡ Optimizado mul con 1: %s\n", newInstr.String())
		}
	}

	// Optimizar sub con 0
	if instr.Opcode == "sub" && len(instr.Operands) == 3 {
		if instr.Operands[2] == "#0" {
			// sub x, y, #0 -> mov x, y
			newInstr.Opcode = "mov"
			newInstr.Operands = []string{instr.Operands[0], instr.Operands[1]}
			newInstr.Comment = "optimized: sub with 0"
			opt.changes++
			fmt.Printf("    ⚡ Optimizado sub con 0: %s\n", newInstr.String())
		}
	}

	return newInstr
}

// removeDeadCode elimina código muerto
func (opt *ARM64Optimizer) removeDeadCode() []*ARM64Instruction {
	var optimized []*ARM64Instruction

	// Mapear todas las etiquetas referenciadas
	referencedLabels := make(map[string]bool)
	for _, instr := range opt.instructions {
		if instr.Opcode == "b" || instr.Opcode == "bl" ||
			strings.HasPrefix(instr.Opcode, "b.") {
			if len(instr.Operands) > 0 {
				referencedLabels[instr.Operands[0]] = true
			}
		}
	}

	// Eliminar etiquetas no referenciadas y código inalcanzable
	for i, instr := range opt.instructions {
		keep := true

		// Eliminar etiquetas no referenciadas
		if instr.Label != "" && !referencedLabels[instr.Label] {
			keep = false
			opt.changes++
			fmt.Printf("    🗑️ Eliminada etiqueta no referenciada: %s\n", instr.Label)
		}

		// Eliminar código después de saltos incondicionales
		if keep && instr.Opcode == "b" && i+1 < len(opt.instructions) {
			nextInstr := opt.instructions[i+1]
			// Si la siguiente instrucción no es una etiqueta, podría ser código muerto
			if nextInstr.Label == "" {
				// Verificar si hay alguna etiqueta antes de la siguiente instrucción
				hasLabel := false
				for j := i + 1; j < len(opt.instructions); j++ {
					if opt.instructions[j].Label != "" {
						hasLabel = true
						break
					}
					if opt.instructions[j].Opcode != "" {
						break
					}
				}
				if !hasLabel {
					// Marcar instrucciones siguientes como muertas hasta encontrar etiqueta
					// (Esta es una implementación simplificada)
				}
			}
		}

		if keep {
			optimized = append(optimized, instr)
		}
	}

	return optimized
}

// peepholeOptimizations aplica optimizaciones peephole
func (opt *ARM64Optimizer) peepholeOptimizations() []*ARM64Instruction {
	var optimized []*ARM64Instruction

	for i := 0; i < len(opt.instructions); i++ {
		instr := opt.instructions[i]
		applied := false

		// Patrón: cmp + cset -> optimización directa
		if i+1 < len(opt.instructions) &&
			instr.Opcode == "cmp" &&
			opt.instructions[i+1].Opcode == "cset" {

			// Verificar si podemos optimizar la comparación
			if opt.canOptimizeComparison(instr, opt.instructions[i+1]) {
				optimizedInstr := opt.optimizeComparison(instr, opt.instructions[i+1])
				if optimizedInstr != nil {
					optimized = append(optimized, optimizedInstr)
					i++ // Saltar la siguiente instrucción
					applied = true
					opt.changes++
					fmt.Printf("    🎯 Optimizada comparación: %s\n", optimizedInstr.String())
				}
			}
		}

		// Patrón: mov + operación aritmética -> optimización con inmediato
		if !applied && i+1 < len(opt.instructions) &&
			instr.Opcode == "mov" &&
			opt.isArithmeticOp(opt.instructions[i+1].Opcode) {

			optimizedInstr := opt.tryInlineImmediate(instr, opt.instructions[i+1])
			if optimizedInstr != nil {
				optimized = append(optimized, optimizedInstr)
				i++ // Saltar la siguiente instrucción
				applied = true
				opt.changes++
				fmt.Printf("    🎯 Inlineado inmediato: %s\n", optimizedInstr.String())
			}
		}

		// Patrón: str + ldr de la misma ubicación -> optimización
		if !applied && i+1 < len(opt.instructions) &&
			instr.Opcode == "str" &&
			opt.instructions[i+1].Opcode == "ldr" {

			if opt.isStoreLoadFromSameLocation(instr, opt.instructions[i+1]) {
				// str x, [addr] seguido de ldr y, [addr] -> str x, [addr]; mov y, x
				movInstr := &ARM64Instruction{
					Opcode:   "mov",
					Operands: []string{opt.instructions[i+1].Operands[0], instr.Operands[0]},
					Comment:  "optimized: store-load bypass",
				}
				optimized = append(optimized, instr)
				optimized = append(optimized, movInstr)
				i++ // Saltar la siguiente instrucción
				applied = true
				opt.changes++
				fmt.Printf("    🎯 Optimizado store-load: %s\n", movInstr.String())
			}
		}

		if !applied {
			optimized = append(optimized, instr)
		}
	}

	return optimized
}

// canOptimizeComparison verifica si una comparación se puede optimizar
func (opt *ARM64Optimizer) canOptimizeComparison(cmpInstr, csetInstr *ARM64Instruction) bool {
	// Verificar que el patrón sea válido
	if len(cmpInstr.Operands) != 2 || len(csetInstr.Operands) != 2 {
		return false
	}

	// Verificar si la comparación es con 0
	return cmpInstr.Operands[1] == "#0"
}

// optimizeComparison optimiza una comparación + cset
func (opt *ARM64Optimizer) optimizeComparison(cmpInstr, csetInstr *ARM64Instruction) *ARM64Instruction {
	// Simplificar cmp x, #0 + cset y, ne -> cmp x, #0; csel y, #1, #0, ne
	// En este caso, simplemente retornamos las instrucciones como están
	// Una optimización real sería más compleja
	return csetInstr
}

// isArithmeticOp verifica si una operación es aritmética
func (opt *ARM64Optimizer) isArithmeticOp(opcode string) bool {
	arithmetic := []string{"add", "sub", "mul", "and", "orr", "eor"}
	for _, op := range arithmetic {
		if opcode == op {
			return true
		}
	}
	return false
}

// tryInlineImmediate intenta inlinear un inmediato en una operación aritmética
func (opt *ARM64Optimizer) tryInlineImmediate(movInstr, arithInstr *ARM64Instruction) *ARM64Instruction {
	// Verificar que mov carga un inmediato y la operación aritmética usa ese registro
	if len(movInstr.Operands) != 2 || len(arithInstr.Operands) != 3 {
		return nil
	}

	// Verificar que mov carga un inmediato
	if !strings.HasPrefix(movInstr.Operands[1], "#") {
		return nil
	}

	// Verificar que la operación aritmética usa el registro destino del mov
	movDest := movInstr.Operands[0]
	if arithInstr.Operands[1] != movDest && arithInstr.Operands[2] != movDest {
		return nil
	}

	// Crear instrucción optimizada
	newInstr := &ARM64Instruction{
		Opcode:   arithInstr.Opcode,
		Operands: make([]string, 3),
		Comment:  "inlined immediate",
	}

	newInstr.Operands[0] = arithInstr.Operands[0]

	if arithInstr.Operands[1] == movDest {
		newInstr.Operands[1] = arithInstr.Operands[2]
		newInstr.Operands[2] = movInstr.Operands[1]
	} else {
		newInstr.Operands[1] = arithInstr.Operands[1]
		newInstr.Operands[2] = movInstr.Operands[1]
	}

	// Verificar que el inmediato sea válido para la operación
	immediate := movInstr.Operands[1]
	if !opt.isValidImmediateForOp(immediate, arithInstr.Opcode) {
		return nil
	}

	return newInstr
}

// isStoreLoadFromSameLocation verifica si store y load son de la misma ubicación
func (opt *ARM64Optimizer) isStoreLoadFromSameLocation(strInstr, ldrInstr *ARM64Instruction) bool {
	if len(strInstr.Operands) != 2 || len(ldrInstr.Operands) != 2 {
		return false
	}

	return strInstr.Operands[1] == ldrInstr.Operands[1]
}

// isValidImmediateForOp verifica si un inmediato es válido para una operación
func (opt *ARM64Optimizer) isValidImmediateForOp(immediate, opcode string) bool {
	// Extraer valor numérico
	immediateStr := strings.TrimPrefix(immediate, "#")
	value, err := strconv.Atoi(immediateStr)
	if err != nil {
		return false
	}

	// Verificar limitaciones según la operación
	switch opcode {
	case "add", "sub":
		// ADD/SUB pueden usar inmediatos de 12 bits
		return value >= 0 && value <= 4095
	case "and", "orr", "eor":
		// Operaciones lógicas tienen limitaciones más complejas
		return value >= 0 && value <= 65535
	case "mul":
		// MUL no acepta inmediatos en ARM64
		return false
	default:
		return false
	}
}

// OptimizeRegisterAllocation optimiza la asignación de registros
func (opt *ARM64Optimizer) OptimizeRegisterAllocation(instructions []*ARM64Instruction) []*ARM64Instruction {
	// Análisis de vida de registros simplificado
	registerUsage := make(map[string][]int) // registro -> lista de líneas donde se usa

	// Recopilar uso de registros
	for i, instr := range instructions {
		for _, operand := range instr.Operands {
			if opt.isRegister(operand) {
				registerUsage[operand] = append(registerUsage[operand], i)
			}
		}
	}

	// Identificar oportunidades de reutilización de registros
	// Esta es una implementación simplificada
	for reg, usage := range registerUsage {
		if len(usage) == 1 {
			fmt.Printf("    📊 Registro %s usado solo una vez en línea %d\n", reg, usage[0])
		}
	}

	// Por ahora, retornar las instrucciones sin cambios
	// Una implementación real haría reasignación de registros
	return instructions
}

// isRegister verifica si un operando es un registro
func (opt *ARM64Optimizer) isRegister(operand string) bool {
	// Detectar registros ARM64 (x0-x30, w0-w30, fp, lr, sp)
	if strings.HasPrefix(operand, "x") || strings.HasPrefix(operand, "w") {
		return true
	}

	return operand == "fp" || operand == "lr" || operand == "sp"
}

// PrintOptimizationStats imprime estadísticas de optimización
func (opt *ARM64Optimizer) PrintOptimizationStats(originalCount, optimizedCount int) {
	fmt.Printf("\n📊 Estadísticas de Optimización ARM64:\n")
	fmt.Printf("  📋 Instrucciones originales: %d\n", originalCount)
	fmt.Printf("  📋 Instrucciones optimizadas: %d\n", optimizedCount)
	fmt.Printf("  📈 Reducción: %d instrucciones (%.1f%%)\n",
		originalCount-optimizedCount,
		float64(originalCount-optimizedCount)/float64(originalCount)*100)
}

// AnalyzeCode analiza el código para identificar patrones de optimización
func (opt *ARM64Optimizer) AnalyzeCode(instructions []*ARM64Instruction) map[string]int {
	stats := make(map[string]int)

	for _, instr := range instructions {
		// Contar tipos de instrucciones
		stats[instr.Opcode]++

		// Detectar patrones específicos
		if instr.Opcode == "mov" && len(instr.Operands) == 2 {
			if strings.HasPrefix(instr.Operands[1], "#") {
				stats["immediate_loads"]++
			} else {
				stats["register_moves"]++
			}
		}

		if strings.HasPrefix(instr.Opcode, "b.") {
			stats["conditional_branches"]++
		}

		if instr.Opcode == "str" || instr.Opcode == "ldr" {
			stats["memory_ops"]++
		}
	}

	return stats
}
