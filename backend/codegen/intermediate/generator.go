package intermediate

import (
	"fmt"
	"strconv"

	"github.com/antlr4-go/antlr/v4"
	compiler "main.go/grammar"
	"main.go/repl"
)

// IRGenerator convierte AST a representación intermedia
type IRGenerator struct {
	program      *IRProgram
	currentFunc  *IRFunction
	instructions []*IRInstruction

	// Contadores para generar nombres únicos
	tempCounter  int
	labelCounter int

	// Scope tracking
	scopeTrace  *repl.ScopeTrace
	symbolTable map[string]*IROperand // mapeo de variables a operandos IR

	// Control de flujo
	breakLabels    []string
	continueLabels []string

	// Manejo de strings
	stringTable map[string]int
}

// NewIRGenerator crea un nuevo generador de IR
func NewIRGenerator() *IRGenerator {
	return &IRGenerator{
		program:        NewIRProgram(),
		tempCounter:    0,
		labelCounter:   0,
		symbolTable:    make(map[string]*IROperand),
		breakLabels:    make([]string, 0),
		continueLabels: make([]string, 0),
		stringTable:    make(map[string]int),
	}
}

// GenerateIR convierte un AST a representación intermedia
func (gen *IRGenerator) GenerateIR(tree antlr.ParseTree, scopeTrace *repl.ScopeTrace) *IRProgram {
	gen.scopeTrace = scopeTrace

	// Crear función main implícita para el programa principal
	mainFunc := &IRFunction{
		Name:         "main",
		Parameters:   make([]*IROperand, 0),
		ReturnType:   "int",
		LocalVars:    make([]*IROperand, 0),
		Instructions: make([]*IRInstruction, 0),
		StackSize:    0,
	}

	gen.currentFunc = mainFunc
	gen.instructions = mainFunc.Instructions

	// Procesar el AST
	gen.visit(tree)

	// Agregar return implícito al final de main
	returnVal := gen.newImmediate(0, "int")
	gen.emit(IR_RETURN, nil, returnVal, nil, "implicit return 0")

	// Actualizar instrucciones de la función
	mainFunc.Instructions = gen.instructions
	gen.program.AddFunction(mainFunc)

	return gen.program
}

// visit despacha la visita según el tipo de nodo
func (gen *IRGenerator) visit(tree antlr.ParseTree) *IROperand {
	if tree == nil {
		return nil
	}

	switch node := tree.(type) {
	case *compiler.ProgramContext:
		return gen.visitProgram(node)
	case *compiler.StmtContext:
		return gen.visitStmt(node)
	case *compiler.MutVarDeclContext:
		return gen.visitMutVarDecl(node)
	case *compiler.ValueDeclContext:
		return gen.visitValueDecl(node)
	case *compiler.VarAssDeclContext:
		return gen.visitVarAssDecl(node)
	case *compiler.AssignmentDeclContext:
		return gen.visitAssignmentDecl(node)
	case *compiler.BinaryExprContext:
		return gen.visitBinaryExpr(node)
	case *compiler.UnaryExprContext:
		return gen.visitUnaryExpr(node)
	case *compiler.LiteralExprContext:
		return gen.visitLiteralExpr(node)
	case *compiler.IdPatternExprContext:
		return gen.visitIdPatternExpr(node)
	case *compiler.IntLiteralContext:
		return gen.visitIntLiteral(node)
	case *compiler.FloatLiteralContext:
		return gen.visitFloatLiteral(node)
	case *compiler.StringLiteralContext:
		return gen.visitStringLiteral(node)
	case *compiler.BoolLiteralContext:
		return gen.visitBoolLiteral(node)
	case *compiler.NilLiteralContext:
		return gen.visitNilLiteral(node)
	case *compiler.IfStmtContext:
		return gen.visitIfStmt(node)
	case *compiler.ForStmtCondContext:
		return gen.visitForStmtCond(node)
	case *compiler.FuncCallContext:
		return gen.visitFuncCall(node)
	case *compiler.ReturnStmtContext:
		return gen.visitReturnStmt(node)
	case *compiler.FuncDeclContext:
		return gen.visitFuncDecl(node)
	default:
		// Para nodos que tengan hijos, visitar todos los hijos
		for i := 0; i < tree.GetChildCount(); i++ {
			child := tree.GetChild(i)
			if parseTreeChild, ok := child.(antlr.ParseTree); ok {
				gen.visit(parseTreeChild)
			}
		}
		return nil
	}
}

// ==================== VISITADORES ESPECÍFICOS ====================

func (gen *IRGenerator) visitProgram(ctx *compiler.ProgramContext) *IROperand {
	// Procesar todas las declaraciones globales primero
	for _, stmt := range ctx.AllStmt() {
		if funcDecl, ok := stmt.(*compiler.StmtContext); ok {
			if funcDecl.Func_dcl() != nil {
				gen.visit(funcDecl.Func_dcl())
			}
		}
	}

	// Luego procesar el resto del programa en main
	for _, stmt := range ctx.AllStmt() {
		if stmtCtx, ok := stmt.(*compiler.StmtContext); ok {
			if stmtCtx.Func_dcl() == nil { // Skip function declarations
				gen.visit(stmt)
			}
		}
	}

	return nil
}

func (gen *IRGenerator) visitStmt(ctx *compiler.StmtContext) *IROperand {
	if ctx.Decl_stmt() != nil {
		return gen.visit(ctx.Decl_stmt())
	}
	if ctx.Assign_stmt() != nil {
		return gen.visit(ctx.Assign_stmt())
	}
	if ctx.If_stmt() != nil {
		return gen.visit(ctx.If_stmt())
	}
	if ctx.For_stmt() != nil {
		return gen.visit(ctx.For_stmt())
	}
	if ctx.Func_call() != nil {
		return gen.visit(ctx.Func_call())
	}
	if ctx.Transfer_stmt() != nil {
		return gen.visit(ctx.Transfer_stmt())
	}
	return nil
}

func (gen *IRGenerator) visitMutVarDecl(ctx *compiler.MutVarDeclContext) *IROperand {
	varName := ctx.ID().GetText()
	varType := ctx.Type_().GetText()

	// Crear operando para la variable
	varOperand := gen.newVariable(varName, varType)
	gen.symbolTable[varName] = varOperand

	// Alocar espacio en el stack
	gen.emit(IR_ALLOC_LOCAL, varOperand, nil, nil, fmt.Sprintf("allocate variable %s", varName))

	// Procesar la expresión de inicialización
	if ctx.Expression() != nil {
		initValue := gen.visit(ctx.Expression())
		gen.emit(IR_STORE, varOperand, initValue, nil, fmt.Sprintf("initialize %s", varName))
	}

	return varOperand
}

func (gen *IRGenerator) visitValueDecl(ctx *compiler.ValueDeclContext) *IROperand {
	varName := ctx.ID().GetText()

	// Evaluar la expresión para inferir el tipo
	initValue := gen.visit(ctx.Expression())
	varType := initValue.DataType

	// Crear operando para la variable
	varOperand := gen.newVariable(varName, varType)
	gen.symbolTable[varName] = varOperand

	// Alocar y inicializar
	gen.emit(IR_ALLOC_LOCAL, varOperand, nil, nil, fmt.Sprintf("allocate variable %s", varName))
	gen.emit(IR_STORE, varOperand, initValue, nil, fmt.Sprintf("initialize %s", varName))

	return varOperand
}

func (gen *IRGenerator) visitVarAssDecl(ctx *compiler.VarAssDeclContext) *IROperand {
	varName := ctx.ID().GetText()
	varType := ctx.Type_().GetText()

	// Crear operando para la variable
	varOperand := gen.newVariable(varName, varType)
	gen.symbolTable[varName] = varOperand

	// Alocar espacio
	gen.emit(IR_ALLOC_LOCAL, varOperand, nil, nil, fmt.Sprintf("allocate variable %s", varName))

	// Inicializar si hay expresión
	if ctx.Expression() != nil {
		initValue := gen.visit(ctx.Expression())
		gen.emit(IR_STORE, varOperand, initValue, nil, fmt.Sprintf("initialize %s", varName))
	}

	return varOperand
}

func (gen *IRGenerator) visitAssignmentDecl(ctx *compiler.AssignmentDeclContext) *IROperand {
	varName := ctx.Id_pattern().GetText()

	// Buscar la variable en la tabla de símbolos
	varOperand, exists := gen.symbolTable[varName]
	if !exists {
		// Variable no encontrada, crear una nueva (esto podría ser un error en análisis semántico)
		varOperand = gen.newVariable(varName, "unknown")
		gen.symbolTable[varName] = varOperand
		gen.emit(IR_ALLOC_LOCAL, varOperand, nil, nil, fmt.Sprintf("allocate variable %s", varName))
	}

	// Evaluar la expresión del lado derecho
	value := gen.visit(ctx.Expression())

	// Generar instrucción de almacenamiento
	gen.emit(IR_STORE, varOperand, value, nil, fmt.Sprintf("assign to %s", varName))

	return varOperand
}

func (gen *IRGenerator) visitBinaryExpr(ctx *compiler.BinaryExprContext) *IROperand {
	op := ctx.GetOp().GetText()

	// Manejar evaluación con corto circuito para && y ||
	if op == "&&" || op == "||" {
		return gen.visitLogicalExpr(ctx, op)
	}

	left := gen.visit(ctx.GetLeft())
	right := gen.visit(ctx.GetRight())

	// Determinar el tipo de resultado
	resultType := gen.getResultType(left.DataType, right.DataType, op)
	result := gen.newTemp(resultType)

	// Generar la instrucción correspondiente
	var irOp IROpcode
	switch op {
	case "+":
		irOp = IR_ADD
	case "-":
		irOp = IR_SUB
	case "*":
		irOp = IR_MULT
	case "/":
		irOp = IR_DIV
	case "%":
		irOp = IR_MOD
	case "==":
		irOp = IR_CMP_EQ
	case "!=":
		irOp = IR_CMP_NE
	case "<":
		irOp = IR_CMP_LT
	case "<=":
		irOp = IR_CMP_LE
	case ">":
		irOp = IR_CMP_GT
	case ">=":
		irOp = IR_CMP_GE
	default:
		irOp = IR_NOP
	}

	gen.emit(irOp, result, left, right, fmt.Sprintf("%s %s %s", left.String(), op, right.String()))
	return result
}

func (gen *IRGenerator) visitLogicalExpr(ctx *compiler.BinaryExprContext, op string) *IROperand {
	result := gen.newTemp("bool")
	left := gen.visit(ctx.GetLeft())

	if op == "&&" {
		// Para &&: si left es false, resultado es false sin evaluar right
		falseLabel := gen.newLabel("and_false")
		endLabel := gen.newLabel("and_end")

		gen.emit(IR_BRANCH_IF_FALSE, nil, left, gen.newLabelOperand(falseLabel), "short-circuit &&")

		// Left es true, evaluar right
		right := gen.visit(ctx.GetRight())
		gen.emit(IR_STORE, result, right, nil, "store && result")
		gen.emit(IR_BRANCH, nil, gen.newLabelOperand(endLabel), nil, "jump to end")

		// Left es false
		gen.emitLabel(falseLabel)
		gen.emit(IR_STORE, result, gen.newImmediate(false, "bool"), nil, "store false")

		gen.emitLabel(endLabel)

	} else { // op == "||"
		// Para ||: si left es true, resultado es true sin evaluar right
		trueLabel := gen.newLabel("or_true")
		endLabel := gen.newLabel("or_end")

		gen.emit(IR_BRANCH_IF_TRUE, nil, left, gen.newLabelOperand(trueLabel), "short-circuit ||")

		// Left es false, evaluar right
		right := gen.visit(ctx.GetRight())
		gen.emit(IR_STORE, result, right, nil, "store || result")
		gen.emit(IR_BRANCH, nil, gen.newLabelOperand(endLabel), nil, "jump to end")

		// Left es true
		gen.emitLabel(trueLabel)
		gen.emit(IR_STORE, result, gen.newImmediate(true, "bool"), nil, "store true")

		gen.emitLabel(endLabel)
	}

	return result
}

func (gen *IRGenerator) visitUnaryExpr(ctx *compiler.UnaryExprContext) *IROperand {
	op := ctx.GetOp().GetText()
	operand := gen.visit(ctx.Expression())

	result := gen.newTemp(operand.DataType)

	switch op {
	case "-":
		gen.emit(IR_NEG, result, operand, nil, fmt.Sprintf("-%s", operand.String()))
	case "!":
		gen.emit(IR_NOT, result, operand, nil, fmt.Sprintf("!%s", operand.String()))
	default:
		return operand // Operador desconocido, retornar operando original
	}

	return result
}

func (gen *IRGenerator) visitLiteralExpr(ctx *compiler.LiteralExprContext) *IROperand {
	return gen.visit(ctx.Literal())
}

func (gen *IRGenerator) visitIdPatternExpr(ctx *compiler.IdPatternExprContext) *IROperand {
	varName := ctx.Id_pattern().GetText()

	// Buscar en la tabla de símbolos
	if varOperand, exists := gen.symbolTable[varName]; exists {
		// Crear un temporal para cargar el valor
		temp := gen.newTemp(varOperand.DataType)
		gen.emit(IR_LOAD, temp, varOperand, nil, fmt.Sprintf("load %s", varName))
		return temp
	}

	// Variable no encontrada, esto debería haber sido detectado en análisis semántico
	return gen.newTemp("unknown")
}

func (gen *IRGenerator) visitIntLiteral(ctx *compiler.IntLiteralContext) *IROperand {
	value, _ := strconv.Atoi(ctx.GetText())
	return gen.newImmediate(value, "int")
}

func (gen *IRGenerator) visitFloatLiteral(ctx *compiler.FloatLiteralContext) *IROperand {
	value, _ := strconv.ParseFloat(ctx.GetText(), 64)
	return gen.newImmediate(value, "float")
}

func (gen *IRGenerator) visitStringLiteral(ctx *compiler.StringLiteralContext) *IROperand {
	// Remover comillas del string literal
	text := ctx.GetText()
	if len(text) >= 2 && text[0] == '"' && text[len(text)-1] == '"' {
		text = text[1 : len(text)-1]
	}

	// Agregar a la tabla de strings
	stringID := gen.program.AddString(text)

	// Crear operando que referencia el string
	return gen.newStringRef(stringID)
}

func (gen *IRGenerator) visitBoolLiteral(ctx *compiler.BoolLiteralContext) *IROperand {
	value := ctx.GetText() == "true"
	return gen.newImmediate(value, "bool")
}

func (gen *IRGenerator) visitNilLiteral(ctx *compiler.NilLiteralContext) *IROperand {
	return gen.newImmediate(nil, "nil")
}

func (gen *IRGenerator) visitIfStmt(ctx *compiler.IfStmtContext) *IROperand {
	elseLabel := gen.newLabel("if_else")
	endLabel := gen.newLabel("if_end")

	// Procesar las cadenas de if
	hasElse := ctx.Else_stmt() != nil
	chainCount := len(ctx.AllIf_chain())

	for i, ifChain := range ctx.AllIf_chain() {
		nextLabel := elseLabel
		if i < chainCount-1 {
			nextLabel = gen.newLabel(fmt.Sprintf("if_chain_%d", i+1))
		}

		// Evaluar condición
		condition := gen.visit(ifChain.(*compiler.IfChainContext).Expression())
		gen.emit(IR_BRANCH_IF_FALSE, nil, condition, gen.newLabelOperand(nextLabel), "if condition")

		// Cuerpo del if
		for _, stmt := range ifChain.(*compiler.IfChainContext).AllStmt() {
			gen.visit(stmt)
		}

		gen.emit(IR_BRANCH, nil, gen.newLabelOperand(endLabel), nil, "jump to end")

		if i < chainCount-1 {
			gen.emitLabel(nextLabel)
		}
	}

	// Else clause
	if hasElse {
		gen.emitLabel(elseLabel)
		if ctx.Else_stmt() != nil {
			elseCtx := ctx.Else_stmt().(*compiler.ElseStmtContext)
			for _, stmt := range elseCtx.AllStmt() {
				gen.visit(stmt)
			}
		}
	} else {
		gen.emitLabel(elseLabel)
	}

	gen.emitLabel(endLabel)
	return nil
}

func (gen *IRGenerator) visitForStmtCond(ctx *compiler.ForStmtCondContext) *IROperand {
	loopStart := gen.newLabel("loop_start")
	loopEnd := gen.newLabel("loop_end")

	// Agregar etiquetas para break y continue
	gen.breakLabels = append(gen.breakLabels, loopEnd)
	gen.continueLabels = append(gen.continueLabels, loopStart)

	gen.emitLabel(loopStart)

	// Evaluar condición
	condition := gen.visit(ctx.Expression())
	gen.emit(IR_BRANCH_IF_FALSE, nil, condition, gen.newLabelOperand(loopEnd), "loop condition")

	// Cuerpo del bucle
	for _, stmt := range ctx.AllStmt() {
		gen.visit(stmt)
	}

	gen.emit(IR_BRANCH, nil, gen.newLabelOperand(loopStart), nil, "loop back")
	gen.emitLabel(loopEnd)

	// Remover etiquetas del stack
	gen.breakLabels = gen.breakLabels[:len(gen.breakLabels)-1]
	gen.continueLabels = gen.continueLabels[:len(gen.continueLabels)-1]

	return nil
}

func (gen *IRGenerator) visitFuncCall(ctx *compiler.FuncCallContext) *IROperand {
	funcName := ctx.Id_pattern().GetText()

	// Evaluar argumentos
	var args []*IROperand
	if ctx.Arg_list() != nil {
		if argListCtx, ok := ctx.Arg_list().(*compiler.ArgListContext); ok {
			args = gen.visitArgList(argListCtx)
		}
	}

	// Función especial print
	if funcName == "print" || funcName == "println" {
		return gen.generatePrintCall(funcName, args)
	}

	// Llamada a función general
	result := gen.newTemp("unknown") // El tipo debería determinarse del análisis semántico

	// Push argumentos (en orden inverso para convención de llamada)
	for i := len(args) - 1; i >= 0; i-- {
		gen.emit(IR_PUSH, nil, args[i], nil, fmt.Sprintf("push arg %d", i))
	}

	// Llamada
	funcOperand := gen.newLabelOperand(funcName)
	gen.emit(IR_CALL, result, funcOperand, nil, fmt.Sprintf("call %s", funcName))

	return result
}

func (gen *IRGenerator) visitArgList(ctx *compiler.ArgListContext) []*IROperand {
	var args []*IROperand
	for _, arg := range ctx.AllFunc_arg() {
		// Simplificado: solo procesar la expresión del argumento
		if argCtx, ok := arg.(*compiler.FuncArgContext); ok {
			if argCtx.Expression() != nil {
				argValue := gen.visit(argCtx.Expression())
				args = append(args, argValue)
			}
		}
	}
	return args
}

func (gen *IRGenerator) generatePrintCall(funcName string, args []*IROperand) *IROperand {
	// Generar llamadas específicas para print
	irOp := IR_PRINT
	if funcName == "println" {
		irOp = IR_PRINT_LN
	}

	for _, arg := range args {
		gen.emit(irOp, nil, arg, nil, fmt.Sprintf("%s argument", funcName))
	}

	return gen.newImmediate(nil, "nil")
}

func (gen *IRGenerator) visitReturnStmt(ctx *compiler.ReturnStmtContext) *IROperand {
	var returnValue *IROperand

	if ctx.Expression() != nil {
		returnValue = gen.visit(ctx.Expression())
	} else {
		returnValue = gen.newImmediate(nil, "nil")
	}

	gen.emit(IR_RETURN, nil, returnValue, nil, "return statement")
	return nil
}

func (gen *IRGenerator) visitFuncDecl(ctx *compiler.FuncDeclContext) *IROperand {
	funcName := ctx.ID().GetText()

	// Crear nueva función
	function := &IRFunction{
		Name:         funcName,
		Parameters:   make([]*IROperand, 0),
		ReturnType:   "nil",
		LocalVars:    make([]*IROperand, 0),
		Instructions: make([]*IRInstruction, 0),
		StackSize:    0,
	}

	// Procesar tipo de retorno
	if ctx.Type_() != nil {
		function.ReturnType = ctx.Type_().GetText()
	}

	// Procesar parámetros
	if ctx.Param_list() != nil {
		if paramListCtx, ok := ctx.Param_list().(*compiler.ParamListContext); ok {
			function.Parameters = gen.visitParamList(paramListCtx)
		}
	}

	// Guardar estado actual y cambiar a nueva función
	prevFunc := gen.currentFunc
	prevInstructions := gen.instructions
	prevSymbols := gen.symbolTable

	gen.currentFunc = function
	gen.instructions = make([]*IRInstruction, 0)
	gen.symbolTable = make(map[string]*IROperand)

	// Agregar parámetros a la tabla de símbolos
	for _, param := range function.Parameters {
		gen.symbolTable[param.Name] = param
	}

	// Prólogo de función
	gen.emit(IR_ENTER_FUNCTION, nil, nil, nil, fmt.Sprintf("enter function %s", funcName))

	// Procesar cuerpo de la función
	for _, stmt := range ctx.AllStmt() {
		gen.visit(stmt)
	}

	// Epílogo de función (return implícito si no hay return explícito)
	gen.emit(IR_EXIT_FUNCTION, nil, nil, nil, fmt.Sprintf("exit function %s", funcName))

	// Actualizar instrucciones de la función
	function.Instructions = gen.instructions
	gen.program.AddFunction(function)

	// Restaurar estado anterior
	gen.currentFunc = prevFunc
	gen.instructions = prevInstructions
	gen.symbolTable = prevSymbols

	return nil
}

func (gen *IRGenerator) visitParamList(ctx *compiler.ParamListContext) []*IROperand {
	var params []*IROperand
	for i, param := range ctx.AllFunc_param() {
		if paramCtx, ok := param.(*compiler.FuncParamContext); ok {
			paramName := paramCtx.ID().GetText()
			paramType := paramCtx.Type_().GetText()

			paramOperand := &IROperand{
				Type:     IR_OPERAND_PARAM,
				Name:     paramName,
				DataType: paramType,
				Offset:   i * 8, // Asumiendo 8 bytes por parámetro
			}

			params = append(params, paramOperand)
		}
	}
	return params
}

// ==================== MÉTODOS AUXILIARES ====================

func (gen *IRGenerator) newTemp(dataType string) *IROperand {
	gen.tempCounter++
	return &IROperand{
		Type:     IR_OPERAND_TEMP,
		Name:     fmt.Sprintf("t%d", gen.tempCounter),
		DataType: dataType,
	}
}

func (gen *IRGenerator) newVariable(name, dataType string) *IROperand {
	return &IROperand{
		Type:     IR_OPERAND_VAR,
		Name:     name,
		DataType: dataType,
		Offset:   0, // Se calculará durante la asignación de stack
	}
}

func (gen *IRGenerator) newImmediate(value interface{}, dataType string) *IROperand {
	return &IROperand{
		Type:     IR_OPERAND_IMMEDIATE,
		Value:    value,
		DataType: dataType,
	}
}

func (gen *IRGenerator) newLabel(prefix string) string {
	gen.labelCounter++
	return fmt.Sprintf("%s_%d", prefix, gen.labelCounter)
}

func (gen *IRGenerator) newLabelOperand(label string) *IROperand {
	return &IROperand{
		Type: IR_OPERAND_LABEL,
		Name: label,
	}
}

func (gen *IRGenerator) newStringRef(stringID int) *IROperand {
	return &IROperand{
		Type:     IR_OPERAND_IMMEDIATE,
		Value:    stringID,
		DataType: "string",
		Name:     fmt.Sprintf("str_%d", stringID),
	}
}

func (gen *IRGenerator) emit(op IROpcode, dest, src1, src2 *IROperand, comment string) {
	instr := &IRInstruction{
		Op:      op,
		Dest:    dest,
		Src1:    src1,
		Src2:    src2,
		Comment: comment,
	}

	gen.instructions = append(gen.instructions, instr)
}

func (gen *IRGenerator) emitLabel(label string) {
	instr := &IRInstruction{
		Op:    IR_LABEL,
		Label: label,
	}

	gen.instructions = append(gen.instructions, instr)
}

func (gen *IRGenerator) getResultType(leftType, rightType, op string) string {
	// Reglas de promoción de tipos para operaciones binarias
	switch op {
	case "==", "!=", "<", "<=", ">", ">=", "&&", "||":
		return "bool"
	case "+", "-", "*", "/":
		if leftType == "float" || rightType == "float" {
			return "float"
		}
		return "int"
	case "%":
		return "int"
	default:
		// Por defecto, retornar el tipo del operando izquierdo
		return leftType
	}
}

// GetProgram retorna el programa IR generado
func (gen *IRGenerator) GetProgram() *IRProgram {
	return gen.program
}
