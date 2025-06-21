package intermediate

import (
	"fmt"
)

// IROptimizer realiza optimizaciones básicas en la representación intermedia
type IROptimizer struct {
	program   *IRProgram
	changes   int                    // Contador de cambios realizados
	constants map[string]interface{} // Tabla de constantes conocidas
}

// NewIROptimizer crea un nuevo optimizador
func NewIROptimizer() *IROptimizer {
	return &IROptimizer{
		constants: make(map[string]interface{}),
	}
}

// Optimize realiza optimizaciones en el programa IR
func (opt *IROptimizer) Optimize(program *IRProgram) *IRProgram {
	opt.program = program

	fmt.Printf("🔧 Iniciando optimización del IR...\n")

	totalChanges := 0
	iteration := 0

	// Repetir optimizaciones hasta que no haya más cambios
	for {
		iteration++
		opt.changes = 0

		fmt.Printf("  📊 Iteración de optimización %d\n", iteration)

		// Aplicar diferentes tipos de optimizaciones
		for _, function := range opt.program.Functions {
			opt.optimizeFunction(function)
		}

		totalChanges += opt.changes
		fmt.Printf("    ✅ %d optimizaciones aplicadas en esta iteración\n", opt.changes)

		// Si no hubo cambios, terminar
		if opt.changes == 0 {
			break
		}

		// Límite de seguridad para evitar loops infinitos
		if iteration > 10 {
			fmt.Printf("  ⚠️ Límite de iteraciones alcanzado\n")
			break
		}
	}

	fmt.Printf("✅ Optimización completada: %d optimizaciones totales en %d iteraciones\n", totalChanges, iteration)

	return opt.program
}

// optimizeFunction aplica optimizaciones a una función específica
func (opt *IROptimizer) optimizeFunction(function *IRFunction) {
	fmt.Printf("    🔍 Optimizando función: %s\n", function.Name)

	// Limpiar tabla de constantes para cada función
	opt.constants = make(map[string]interface{})

	// Aplicar optimizaciones en orden
	opt.constantPropagation(function)
	opt.constantFolding(function)
	opt.deadCodeElimination(function)
	opt.removeRedundantLoads(function)
	opt.removeUnusedLabels(function)
	opt.peepholeOptimizations(function)
}

// ==================== PROPAGACIÓN DE CONSTANTES ====================

func (opt *IROptimizer) constantPropagation(function *IRFunction) {
	changed := false

	for _, instr := range function.Instructions {
		switch instr.Op {
		case IR_LOAD_IMMEDIATE:
			// Registrar constante conocida
			if instr.Dest != nil && instr.Src1 != nil && instr.Src1.IsImmediate() {
				opt.constants[instr.Dest.Name] = instr.Src1.Value
				instr.Comment += " [const registered]"
			}

		case IR_STORE:
			// Si almacenamos una constante, registrarla
			if instr.Dest != nil && instr.Src1 != nil && instr.Src1.IsImmediate() {
				opt.constants[instr.Dest.Name] = instr.Src1.Value
			}

		case IR_ADD, IR_SUB, IR_MULT, IR_DIV, IR_MOD:
			// Propagar constantes en operaciones binarias
			if instr.Src1 != nil && instr.Src1.IsTemp() {
				if value, exists := opt.constants[instr.Src1.Name]; exists {
					instr.Src1 = &IROperand{
						Type:     IR_OPERAND_IMMEDIATE,
						Value:    value,
						DataType: instr.Src1.DataType,
					}
					instr.Comment += " [const prop src1]"
					changed = true
				}
			}

			if instr.Src2 != nil && instr.Src2.IsTemp() {
				if value, exists := opt.constants[instr.Src2.Name]; exists {
					instr.Src2 = &IROperand{
						Type:     IR_OPERAND_IMMEDIATE,
						Value:    value,
						DataType: instr.Src2.DataType,
					}
					instr.Comment += " [const prop src2]"
					changed = true
				}
			}
		}
	}

	if changed {
		opt.changes++
	}
}

// ==================== FOLDING DE CONSTANTES ====================

func (opt *IROptimizer) constantFolding(function *IRFunction) {
	for i, instr := range function.Instructions {
		if opt.tryFoldConstant(instr) {
			opt.changes++

			// Reemplazar con una instrucción de carga inmediata
			newInstr := &IRInstruction{
				Op:      IR_LOAD_IMMEDIATE,
				Dest:    instr.Dest,
				Src1:    instr.Dest, // El destino ahora contiene el valor calculado
				Comment: fmt.Sprintf("folded: %s", instr.Comment),
			}

			function.Instructions[i] = newInstr
		}
	}
}

func (opt *IROptimizer) tryFoldConstant(instr *IRInstruction) bool {
	if instr.Src1 == nil || instr.Src2 == nil {
		return false
	}

	if !instr.Src1.IsImmediate() || !instr.Src2.IsImmediate() {
		return false
	}

	var result interface{}
	var success bool

	switch instr.Op {
	case IR_ADD:
		result, success = opt.foldAdd(instr.Src1.Value, instr.Src2.Value)
	case IR_SUB:
		result, success = opt.foldSub(instr.Src1.Value, instr.Src2.Value)
	case IR_MULT:
		result, success = opt.foldMult(instr.Src1.Value, instr.Src2.Value)
	case IR_DIV:
		result, success = opt.foldDiv(instr.Src1.Value, instr.Src2.Value)
	case IR_MOD:
		result, success = opt.foldMod(instr.Src1.Value, instr.Src2.Value)
	case IR_CMP_EQ:
		result, success = opt.foldCmpEq(instr.Src1.Value, instr.Src2.Value)
	case IR_CMP_NE:
		result, success = opt.foldCmpNe(instr.Src1.Value, instr.Src2.Value)
	case IR_CMP_LT:
		result, success = opt.foldCmpLt(instr.Src1.Value, instr.Src2.Value)
	case IR_CMP_LE:
		result, success = opt.foldCmpLe(instr.Src1.Value, instr.Src2.Value)
	case IR_CMP_GT:
		result, success = opt.foldCmpGt(instr.Src1.Value, instr.Src2.Value)
	case IR_CMP_GE:
		result, success = opt.foldCmpGe(instr.Src1.Value, instr.Src2.Value)
	default:
		return false
	}

	if success {
		// Actualizar el destino con el resultado calculado
		instr.Dest.Type = IR_OPERAND_IMMEDIATE
		instr.Dest.Value = result
		return true
	}

	return false
}

// Funciones auxiliares para folding aritmético
func (opt *IROptimizer) foldAdd(a, b interface{}) (interface{}, bool) {
	switch va := a.(type) {
	case int:
		if vb, ok := b.(int); ok {
			return va + vb, true
		}
	case float64:
		if vb, ok := b.(float64); ok {
			return va + vb, true
		}
		if vb, ok := b.(int); ok {
			return va + float64(vb), true
		}
	}
	return nil, false
}

func (opt *IROptimizer) foldSub(a, b interface{}) (interface{}, bool) {
	switch va := a.(type) {
	case int:
		if vb, ok := b.(int); ok {
			return va - vb, true
		}
	case float64:
		if vb, ok := b.(float64); ok {
			return va - vb, true
		}
		if vb, ok := b.(int); ok {
			return va - float64(vb), true
		}
	}
	return nil, false
}

func (opt *IROptimizer) foldMult(a, b interface{}) (interface{}, bool) {
	switch va := a.(type) {
	case int:
		if vb, ok := b.(int); ok {
			return va * vb, true
		}
	case float64:
		if vb, ok := b.(float64); ok {
			return va * vb, true
		}
		if vb, ok := b.(int); ok {
			return va * float64(vb), true
		}
	}
	return nil, false
}

func (opt *IROptimizer) foldDiv(a, b interface{}) (interface{}, bool) {
	switch va := a.(type) {
	case int:
		if vb, ok := b.(int); ok && vb != 0 {
			return va / vb, true
		}
	case float64:
		if vb, ok := b.(float64); ok && vb != 0.0 {
			return va / vb, true
		}
		if vb, ok := b.(int); ok && vb != 0 {
			return va / float64(vb), true
		}
	}
	return nil, false
}

func (opt *IROptimizer) foldMod(a, b interface{}) (interface{}, bool) {
	if va, ok := a.(int); ok {
		if vb, ok := b.(int); ok && vb != 0 {
			return va % vb, true
		}
	}
	return nil, false
}

// Funciones auxiliares para folding de comparaciones
func (opt *IROptimizer) foldCmpEq(a, b interface{}) (interface{}, bool) {
	return a == b, true
}

func (opt *IROptimizer) foldCmpNe(a, b interface{}) (interface{}, bool) {
	return a != b, true
}

func (opt *IROptimizer) foldCmpLt(a, b interface{}) (interface{}, bool) {
	switch va := a.(type) {
	case int:
		if vb, ok := b.(int); ok {
			return va < vb, true
		}
	case float64:
		if vb, ok := b.(float64); ok {
			return va < vb, true
		}
	}
	return nil, false
}

func (opt *IROptimizer) foldCmpLe(a, b interface{}) (interface{}, bool) {
	switch va := a.(type) {
	case int:
		if vb, ok := b.(int); ok {
			return va <= vb, true
		}
	case float64:
		if vb, ok := b.(float64); ok {
			return va <= vb, true
		}
	}
	return nil, false
}

func (opt *IROptimizer) foldCmpGt(a, b interface{}) (interface{}, bool) {
	switch va := a.(type) {
	case int:
		if vb, ok := b.(int); ok {
			return va > vb, true
		}
	case float64:
		if vb, ok := b.(float64); ok {
			return va > vb, true
		}
	}
	return nil, false
}

func (opt *IROptimizer) foldCmpGe(a, b interface{}) (interface{}, bool) {
	switch va := a.(type) {
	case int:
		if vb, ok := b.(int); ok {
			return va >= vb, true
		}
	case float64:
		if vb, ok := b.(float64); ok {
			return va >= vb, true
		}
	}
	return nil, false
}

// ==================== ELIMINACIÓN DE CÓDIGO MUERTO ====================

func (opt *IROptimizer) deadCodeElimination(function *IRFunction) {
	// Marcar temporales usados
	usedTemps := make(map[string]bool)

	// Primera pasada: marcar todos los temporales que son usados
	for _, instr := range function.Instructions {
		if instr.Src1 != nil && instr.Src1.IsTemp() {
			usedTemps[instr.Src1.Name] = true
		}
		if instr.Src2 != nil && instr.Src2.IsTemp() {
			usedTemps[instr.Src2.Name] = true
		}
	}

	// Segunda pasada: eliminar instrucciones que definen temporales no usados
	newInstructions := make([]*IRInstruction, 0)

	for _, instr := range function.Instructions {
		keep := true

		// Eliminar asignaciones a temporales no usados
		if instr.Dest != nil && instr.Dest.IsTemp() {
			if !usedTemps[instr.Dest.Name] {
				// Solo eliminar si la instrucción no tiene efectos secundarios
				if opt.hasNoSideEffects(instr) {
					keep = false
					opt.changes++
				}
			}
		}

		// Eliminar NOPs
		if instr.Op == IR_NOP {
			keep = false
			opt.changes++
		}

		if keep {
			newInstructions = append(newInstructions, instr)
		}
	}

	function.Instructions = newInstructions
}

func (opt *IROptimizer) hasNoSideEffects(instr *IRInstruction) bool {
	switch instr.Op {
	case IR_ADD, IR_SUB, IR_MULT, IR_DIV, IR_MOD:
		return true
	case IR_CMP_EQ, IR_CMP_NE, IR_CMP_LT, IR_CMP_LE, IR_CMP_GT, IR_CMP_GE:
		return true
	case IR_LOAD_IMMEDIATE, IR_LOAD:
		return true
	case IR_NEG, IR_NOT:
		return true
	default:
		return false // Las operaciones no listadas pueden tener efectos secundarios
	}
}

// ==================== ELIMINACIÓN DE CARGAS REDUND
// ==================== ELIMINACIÓN DE CARGAS REDUNDANTES ====================

func (opt *IROptimizer) removeRedundantLoads(function *IRFunction) {
	// Rastrear el último valor cargado para cada variable
	lastLoaded := make(map[string]*IROperand)

	for i, instr := range function.Instructions {
		switch instr.Op {
		case IR_LOAD:
			if instr.Src1 != nil && instr.Src1.IsVariable() {
				varName := instr.Src1.Name

				// Si ya cargamos esta variable recientemente y no ha sido modificada
				if lastTemp, exists := lastLoaded[varName]; exists {
					// Reemplazar esta carga con una copia del temporal anterior
					newInstr := &IRInstruction{
						Op:      IR_LOAD_IMMEDIATE,
						Dest:    instr.Dest,
						Src1:    lastTemp,
						Comment: fmt.Sprintf("redundant load eliminated: %s", instr.Comment),
					}
					function.Instructions[i] = newInstr
					opt.changes++
				} else {
					// Registrar este temporal como el último cargado para esta variable
					lastLoaded[varName] = instr.Dest
				}
			}

		case IR_STORE:
			// Si se almacena en una variable, invalidar su carga cacheada
			if instr.Dest != nil && instr.Dest.IsVariable() {
				delete(lastLoaded, instr.Dest.Name)
			}

		case IR_CALL:
			// Las llamadas a función pueden modificar variables globales
			// Por seguridad, limpiar todas las cargas cacheadas
			lastLoaded = make(map[string]*IROperand)
		}
	}
}

// ==================== ELIMINACIÓN DE ETIQUETAS NO USADAS ====================

func (opt *IROptimizer) removeUnusedLabels(function *IRFunction) {
	// Encontrar todas las etiquetas referenciadas
	usedLabels := make(map[string]bool)

	for _, instr := range function.Instructions {
		switch instr.Op {
		case IR_BRANCH, IR_BRANCH_IF_TRUE, IR_BRANCH_IF_FALSE:
			if instr.Src1 != nil && instr.Src1.Type == IR_OPERAND_LABEL {
				usedLabels[instr.Src1.Name] = true
			}
			if instr.Src2 != nil && instr.Src2.Type == IR_OPERAND_LABEL {
				usedLabels[instr.Src2.Name] = true
			}
		}
	}

	// Eliminar etiquetas no usadas
	newInstructions := make([]*IRInstruction, 0)

	for _, instr := range function.Instructions {
		keep := true

		if instr.Op == IR_LABEL {
			if !usedLabels[instr.Label] {
				keep = false
				opt.changes++
			}
		}

		if keep {
			newInstructions = append(newInstructions, instr)
		}
	}

	function.Instructions = newInstructions
}

// ==================== OPTIMIZACIONES PEEPHOLE ====================

func (opt *IROptimizer) peepholeOptimizations(function *IRFunction) {
	for i := 0; i < len(function.Instructions)-1; i++ {
		current := function.Instructions[i]
		next := function.Instructions[i+1]

		// Patrón: STORE seguido de LOAD de la misma variable
		if opt.isStoreLoadPattern(current, next) {
			// Reemplazar LOAD con copia directa del valor almacenado
			newInstr := &IRInstruction{
				Op:      IR_LOAD_IMMEDIATE,
				Dest:    next.Dest,
				Src1:    current.Src1,
				Comment: "store-load optimization",
			}
			function.Instructions[i+1] = newInstr
			opt.changes++
		}

		// Patrón: Operaciones redundantes (x + 0, x * 1, etc.)
		if opt.isRedundantOperation(current) {
			opt.simplifyRedundantOperation(current)
			opt.changes++
		}

		// Patrón: Saltos a la siguiente instrucción
		if opt.isRedundantJump(current, i, function) {
			current.Op = IR_NOP
			current.Comment = "redundant jump removed"
			opt.changes++
		}
	}
}

func (opt *IROptimizer) isStoreLoadPattern(store, load *IRInstruction) bool {
	if store.Op != IR_STORE || load.Op != IR_LOAD {
		return false
	}

	if store.Dest == nil || load.Src1 == nil {
		return false
	}

	// Verificar que sea la misma variable
	return store.Dest.IsVariable() && load.Src1.IsVariable() &&
		store.Dest.Name == load.Src1.Name
}

func (opt *IROptimizer) isRedundantOperation(instr *IRInstruction) bool {
	if instr.Src1 == nil || instr.Src2 == nil {
		return false
	}

	switch instr.Op {
	case IR_ADD:
		// x + 0 = x
		return opt.isZero(instr.Src2) || opt.isZero(instr.Src1)
	case IR_SUB:
		// x - 0 = x
		return opt.isZero(instr.Src2)
	case IR_MULT:
		// x * 1 = x, x * 0 = 0
		return opt.isOne(instr.Src2) || opt.isOne(instr.Src1) ||
			opt.isZero(instr.Src2) || opt.isZero(instr.Src1)
	case IR_DIV:
		// x / 1 = x
		return opt.isOne(instr.Src2)
	default:
		return false
	}
}

func (opt *IROptimizer) simplifyRedundantOperation(instr *IRInstruction) {
	switch instr.Op {
	case IR_ADD:
		if opt.isZero(instr.Src2) {
			// x + 0 = x
			instr.Op = IR_LOAD_IMMEDIATE
			instr.Src1 = instr.Src1
			instr.Src2 = nil
			instr.Comment += " [simplified: x + 0]"
		} else if opt.isZero(instr.Src1) {
			// 0 + x = x
			instr.Op = IR_LOAD_IMMEDIATE
			instr.Src1 = instr.Src2
			instr.Src2 = nil
			instr.Comment += " [simplified: 0 + x]"
		}

	case IR_SUB:
		if opt.isZero(instr.Src2) {
			// x - 0 = x
			instr.Op = IR_LOAD_IMMEDIATE
			instr.Src1 = instr.Src1
			instr.Src2 = nil
			instr.Comment += " [simplified: x - 0]"
		}

	case IR_MULT:
		if opt.isOne(instr.Src2) {
			// x * 1 = x
			instr.Op = IR_LOAD_IMMEDIATE
			instr.Src1 = instr.Src1
			instr.Src2 = nil
			instr.Comment += " [simplified: x * 1]"
		} else if opt.isOne(instr.Src1) {
			// 1 * x = x
			instr.Op = IR_LOAD_IMMEDIATE
			instr.Src1 = instr.Src2
			instr.Src2 = nil
			instr.Comment += " [simplified: 1 * x]"
		} else if opt.isZero(instr.Src1) || opt.isZero(instr.Src2) {
			// x * 0 = 0 or 0 * x = 0
			instr.Op = IR_LOAD_IMMEDIATE
			instr.Src1 = &IROperand{
				Type:     IR_OPERAND_IMMEDIATE,
				Value:    0,
				DataType: instr.Dest.DataType,
			}
			instr.Src2 = nil
			instr.Comment += " [simplified: x * 0]"
		}

	case IR_DIV:
		if opt.isOne(instr.Src2) {
			// x / 1 = x
			instr.Op = IR_LOAD_IMMEDIATE
			instr.Src1 = instr.Src1
			instr.Src2 = nil
			instr.Comment += " [simplified: x / 1]"
		}
	}
}

func (opt *IROptimizer) isRedundantJump(instr *IRInstruction, index int, function *IRFunction) bool {
	if instr.Op != IR_BRANCH {
		return false
	}

	if instr.Src1 == nil || instr.Src1.Type != IR_OPERAND_LABEL {
		return false
	}

	targetLabel := instr.Src1.Name

	// Verificar si el salto es a la siguiente instrucción
	if index+1 < len(function.Instructions) {
		nextInstr := function.Instructions[index+1]
		if nextInstr.Op == IR_LABEL && nextInstr.Label == targetLabel {
			return true
		}
	}

	return false
}

func (opt *IROptimizer) isZero(operand *IROperand) bool {
	if !operand.IsImmediate() {
		return false
	}

	switch v := operand.Value.(type) {
	case int:
		return v == 0
	case float64:
		return v == 0.0
	default:
		return false
	}
}

func (opt *IROptimizer) isOne(operand *IROperand) bool {
	if !operand.IsImmediate() {
		return false
	}

	switch v := operand.Value.(type) {
	case int:
		return v == 1
	case float64:
		return v == 1.0
	default:
		return false
	}
}

// ==================== ANÁLISIS DE FLUJO DE DATOS ====================

func (opt *IROptimizer) dataFlowAnalysis(function *IRFunction) {
	// Análisis básico de definiciones y usos
	definitions := make(map[string][]*IRInstruction) // Variable -> instrucciones que la definen
	uses := make(map[string][]*IRInstruction)        // Variable -> instrucciones que la usan

	// Recopilar definiciones y usos
	for _, instr := range function.Instructions {
		// Definiciones (dest)
		if instr.Dest != nil && (instr.Dest.IsVariable() || instr.Dest.IsTemp()) {
			varName := instr.Dest.Name
			definitions[varName] = append(definitions[varName], instr)
		}

		// Usos (src1, src2)
		if instr.Src1 != nil && (instr.Src1.IsVariable() || instr.Src1.IsTemp()) {
			varName := instr.Src1.Name
			uses[varName] = append(uses[varName], instr)
		}

		if instr.Src2 != nil && (instr.Src2.IsVariable() || instr.Src2.IsTemp()) {
			varName := instr.Src2.Name
			uses[varName] = append(uses[varName], instr)
		}
	}

	// Eliminar variables que se definen pero nunca se usan
	for varName, defs := range definitions {
		if _, hasUses := uses[varName]; !hasUses {
			// Variable definida pero nunca usada
			for _, defInstr := range defs {
				if opt.hasNoSideEffects(defInstr) {
					defInstr.Op = IR_NOP
					defInstr.Comment = fmt.Sprintf("unused variable %s eliminated", varName)
					opt.changes++
				}
			}
		}
	}
}

// ==================== ESTADÍSTICAS DE OPTIMIZACIÓN ====================

func (opt *IROptimizer) PrintOptimizationStats(program *IRProgram) {
	fmt.Printf("\n📊 Estadísticas de Optimización:\n")

	totalInstructions := 0
	totalFunctions := len(program.Functions)

	for _, function := range program.Functions {
		totalInstructions += len(function.Instructions)
		fmt.Printf("  📋 Función %s: %d instrucciones\n", function.Name, len(function.Instructions))
	}

	fmt.Printf("  📈 Total: %d funciones, %d instrucciones\n", totalFunctions, totalInstructions)
	fmt.Printf("  🎯 Strings en tabla: %d\n", len(program.StringTable))
	fmt.Printf("  🌍 Variables globales: %d\n", len(program.GlobalVars))
}

// ==================== VALIDACIÓN DE IR ====================

func (opt *IROptimizer) ValidateIR(program *IRProgram) []string {
	var errors []string

	for _, function := range program.Functions {
		// Validar que todas las etiquetas referenciadas existen
		labelMap := make(map[string]bool)

		// Recopilar todas las etiquetas definidas
		for _, instr := range function.Instructions {
			if instr.Op == IR_LABEL {
				labelMap[instr.Label] = true
			}
		}

		// Verificar referencias a etiquetas
		for _, instr := range function.Instructions {
			if instr.Src1 != nil && instr.Src1.Type == IR_OPERAND_LABEL {
				if !labelMap[instr.Src1.Name] {
					errors = append(errors, fmt.Sprintf("función %s: etiqueta no definida '%s'",
						function.Name, instr.Src1.Name))
				}
			}

			if instr.Src2 != nil && instr.Src2.Type == IR_OPERAND_LABEL {
				if !labelMap[instr.Src2.Name] {
					errors = append(errors, fmt.Sprintf("función %s: etiqueta no definida '%s'",
						function.Name, instr.Src2.Name))
				}
			}
		}
	}

	return errors
}
