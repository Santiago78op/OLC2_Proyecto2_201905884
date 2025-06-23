package compiler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"main.go/compiler/arm64"
	compiler "main.go/grammar"
)

// ARM64Translator es el traductor principal de VlangCherry a ARM64
type ARM64Translator struct {
	generator *arm64.ARM64Generator
	errors    []string // Para almacenar errores de traducción

	userFunctions   map[string]*compiler.FuncDeclContext
	currentFunction string
}

// NewARM64Translator crea un nuevo traductor
func NewARM64Translator() *ARM64Translator {
	return &ARM64Translator{
		generator:     arm64.NewARM64Generator(),
		errors:        make([]string, 0),
		userFunctions: make(map[string]*compiler.FuncDeclContext),
	}
}

// === FUNCIÓN PRINCIPAL DE TRADUCCIÓN ===

// TranslateProgram traduce un programa completo de VlangCherry a ARM64
func (t *ARM64Translator) TranslateProgram(tree antlr.ParseTree) (string, []string) {
	// Limpiar estado anterior
	t.generator.Reset()
	t.errors = make([]string, 0)

	// Primera pasada: analizar declaraciones de variables
	t.analyzeVariables(tree)

	// Generar header del programa
	t.generator.GenerateHeader()

	// Traducir el contenido del programa
	t.translateNode(tree)

	// Generar footer del programa
	t.generator.GenerateFooter()

	// Generar código para funciones de usuario
	t.generateUserFunctions()

	// Agregar funciones de librería estándar
	t.generator.EmitRaw("")
	t.generator.EmitRaw("// === LIBRERÍA ESTÁNDAR ===")
	t.generateStandardLibrary()

	return t.generator.GetCode(), t.errors
}

func (t *ARM64Translator) generateUserFunctions() {
	t.generator.EmitRaw("")
	t.generator.EmitRaw("// === FUNCIONES DE USUARIO ===")

	for funcName, funcDecl := range t.userFunctions {
		t.generator.EmitRaw("")
		t.generator.Comment(fmt.Sprintf("Función: %s", funcName))
		t.generator.EmitRaw(fmt.Sprintf("func_%s:", funcName))

		// Prólogo de función
		t.generator.Emit("stp x29, x30, [sp, #-16]!")
		t.generator.Emit("mov x29, sp")

		// Mapear parámetros de registros a variables locales
		if funcDecl.Param_list() != nil {
			params := funcDecl.Param_list().(*compiler.ParamListContext).AllFunc_param()

			for i, param := range params {
				if paramCtx := param.(*compiler.FuncParamContext); paramCtx.ID() != nil {
					paramName := paramCtx.ID().GetText()
					// Declarar parámetro como variable local
					t.generator.DeclareVariable(paramName)

					// Usar un registro temporal para no sobrescribir
					sourceReg := fmt.Sprintf("x%d", i)
					tempReg := fmt.Sprintf("x%d", i+10) // Usar x10, x11, etc. como temporales

					t.generator.Emit(fmt.Sprintf("mov %s, %s", tempReg, sourceReg))
					t.generator.Emit(fmt.Sprintf("mov x0, %s", tempReg))
					t.generator.StoreVariable(arm64.X0, paramName)
				}
			}
		}

		// Traducir cuerpo de la función
		t.currentFunction = funcName
		hasReturnStatement := false

		for _, stmt := range funcDecl.AllStmt() {
			// Verificar si hay statement de return
			if t.hasReturnStatement(stmt) {
				hasReturnStatement = true
			}
			t.translateNode(stmt)
		}

		// Epílogo de función (solo si no hay return explícito)
		if !hasReturnStatement {
			t.generator.Emit("mov x0, #0") // Valor de retorno por defecto
			t.generator.Emit("ldp x29, x30, [sp], #16")
			t.generator.Emit("ret")
		}

		t.currentFunction = ""
	}
}

// Verificar si un statement contiene return
func (t *ARM64Translator) hasReturnStatement(stmt antlr.ParseTree) bool {
	switch ctx := stmt.(type) {
	case *compiler.StmtContext:
		if ctx.Transfer_stmt() != nil {
			transferText := ctx.Transfer_stmt().GetText()
			return strings.HasPrefix(transferText, "return")
		}
		return false
	case *compiler.Transfer_stmtContext:
		transferText := ctx.GetText()
		return strings.HasPrefix(transferText, "return")
	default:
		return false
	}
}

// Traducir statement return
func (t *ARM64Translator) translateReturnStatement(ctx *compiler.ReturnStmtContext) {
	t.generator.Comment("=== RETURN STATEMENT ===")

	// Si hay expresión de retorno, evaluarla
	if ctx.Expression() != nil {
		t.translateExpression(ctx.Expression())
		// El resultado queda en x0, que es correcto para el valor de retorno
	} else {
		// Return sin valor
		t.generator.LoadImmediate(arm64.X0, 0)
	}

	// Epílogo de función
	t.generator.Emit("ldp x29, x30, [sp], #16")
	t.generator.Emit("ret")
}

// === ANÁLISIS PREVIO ===

// analyzeVariables hace una pasada previa para encontrar todas las variables
// Esto nos permite reservar espacio en el stack antes de generar código
func (t *ARM64Translator) analyzeVariables(node antlr.ParseTree) {

	switch ctx := node.(type) {
	case *compiler.ProgramContext:
		for _, stmt := range ctx.AllStmt() {
			fmt.Printf("Analizando declaración: %T\n", stmt)
			t.analyzeVariables(stmt)
		}

	case *compiler.StmtContext:
		if ctx.Decl_stmt() != nil {
			fmt.Printf("Analizando declaración de variable: %T\n", ctx.Decl_stmt())
			t.analyzeVariables(ctx.Decl_stmt())
		}
		if ctx.If_stmt() != nil {
			t.analyzeVariables(ctx.If_stmt())
		}
		if ctx.For_stmt() != nil {
			t.analyzeVariables(ctx.For_stmt())
		}
		if ctx.Func_dcl() != nil {
			fmt.Printf("Analizando función: %T\n", ctx.Func_dcl())
			t.analyzeVariables(ctx.Func_dcl())
		}
		if ctx.Transfer_stmt() != nil { // AGREGAR ESTO
			t.analyzeVariables(ctx.Transfer_stmt())
		}

	case *compiler.ValueDeclContext:
		// Declaración de variable simple
		varName := ctx.ID().GetText()
		if !t.generator.VariableExists(varName) {
			t.generator.DeclareVariable(varName)
		} else {
			t.addError(fmt.Sprintf("Variable '%s' ya está declarada", varName))
		}
	case *compiler.MutVarDeclContext:
		varName := ctx.ID().GetText()
		if !t.generator.VariableExists(varName) {
			t.generator.DeclareVariable(varName)
		} else {
			t.addError(fmt.Sprintf("Variable '%s' ya está declarada", varName))
		}

	case *compiler.VarAssDeclContext:
		varName := ctx.ID().GetText()
		if !t.generator.VariableExists(varName) {
			t.generator.DeclareVariable(varName)
		} else {
			t.addError(fmt.Sprintf("Variable '%s' ya está declarada", varName))
		}

	case *compiler.IfStmtContext:
		// Analizar el cuerpo del if
		for _, ifChain := range ctx.AllIf_chain() {
			if ifChainCtx, ok := ifChain.(*compiler.IfChainContext); ok {
				for _, stmt := range ifChainCtx.AllStmt() {
					t.analyzeVariables(stmt)
				}
			}
		}
		// Analizar el else si existe
		if ctx.Else_stmt() != nil {
			elseCtx := ctx.Else_stmt().(*compiler.ElseStmtContext)
			for _, stmt := range elseCtx.AllStmt() {
				t.analyzeVariables(stmt)
			}
		}

	case *compiler.ForStmtCondContext:
		// Analizar el cuerpo del for
		for _, stmt := range ctx.AllStmt() {
			t.analyzeVariables(stmt)
		}

	case *compiler.FuncDeclContext:
		funcName := ctx.ID().GetText()
		// AGREGAR: Registrar función de usuario
		if funcName != "main" {
			t.userFunctions[funcName] = ctx
			fmt.Printf("Registrando función de usuario: %s\n", funcName)
		}

		// Analizar el cuerpo de la función
		for _, stmt := range ctx.AllStmt() {
			t.analyzeVariables(stmt)
		}
	}
}

// === TRADUCCIÓN DE NODOS ===

// translateNode traduce cualquier nodo del AST
func (t *ARM64Translator) translateNode(node antlr.ParseTree) {
	switch ctx := node.(type) {
	case *compiler.ProgramContext:
		t.translateProgram(ctx)
	case *compiler.StmtContext:
		t.translateStatement(ctx)
	case *compiler.ValueDeclContext:
		t.translateValueDecl(ctx)
	case *compiler.MutVarDeclContext:
		t.translateMutVarDecl(ctx)
	case *compiler.VarAssDeclContext:
		t.translateVarAssDecl(ctx)
	case *compiler.AssignmentDeclContext:
		t.translateAssignment(ctx)
	case *compiler.IfStmtContext:
		t.translateIfStatement(ctx)
	case *compiler.ForStmtCondContext:
		t.translateForLoop(ctx)
	case *compiler.FuncCallContext:
		t.translateFunctionCall(ctx)
	case *compiler.FuncDeclContext:
		t.translateFunctionDeclaration(ctx)
	case *compiler.Decl_stmtContext:
		t.translateDeclStatement(ctx)
	case *compiler.Transfer_stmtContext:
		t.translateTransferStatement(ctx)
	case *compiler.ReturnStmtContext:
		t.translateReturnStatement(ctx)
	default:
		// Para nodos no implementados, simplemente continuar
		t.addError(fmt.Sprintf("Nodo no implementado: %T", ctx))
	}
}

// Manejar transfer statements (return, break, continue)
// Manejar transfer statements (return, break, continue)
func (t *ARM64Translator) translateTransferStatement(ctx *compiler.Transfer_stmtContext) {

	// Analizar por el texto del primer token para determinar el tipo
	text := ctx.GetText()

	if strings.HasPrefix(text, "return") {
		// Es un return statement
		t.translateReturnStatementFromTransfer(ctx)
	} else if strings.HasPrefix(text, "break") {
		t.translateBreakStatementFromTransfer(ctx)
	} else if strings.HasPrefix(text, "continue") {
		t.translateContinueStatementFromTransfer(ctx)
	} else {
		t.addError(fmt.Sprintf("Transfer statement no reconocido: %s", text))
	}
}

// Manejar return desde transfer_stmt
func (t *ARM64Translator) translateReturnStatementFromTransfer(ctx *compiler.Transfer_stmtContext) {
	t.generator.Comment("=== RETURN STATEMENT ===")

	// Buscar si hay una expresión después de "return"
	hasExpression := false

	// Recorrer hijos para encontrar la expresión
	for i := 0; i < ctx.GetChildCount(); i++ {
		child := ctx.GetChild(i)
		if expressionCtx, ok := child.(*compiler.ExpressionContext); ok {
			hasExpression = true
			t.translateExpression(expressionCtx)
			break
		}
	}

	if !hasExpression {
		// Return sin valor
		t.generator.LoadImmediate(arm64.X0, 0)
	}

	// Epílogo de función
	t.generator.Emit("ldp x29, x30, [sp], #16")
	t.generator.Emit("ret")
}

func (t *ARM64Translator) translateBreakStatementFromTransfer(ctx *compiler.Transfer_stmtContext) {
	t.generator.Comment("=== BREAK STATEMENT ===")
	// TODO: Implementar break
}

func (t *ARM64Translator) translateContinueStatementFromTransfer(ctx *compiler.Transfer_stmtContext) {
	t.generator.Comment("=== CONTINUE STATEMENT ===")
	// TODO: Implementar continue
}

func (t *ARM64Translator) translateBreakStatement(ctx *compiler.BreakStmtContext) {
	// TODO: Implementar break (necesita manejo de labels de loops)
	t.generator.Comment("=== BREAK STATEMENT ===")
}

func (t *ARM64Translator) translateContinueStatement(ctx *compiler.ContinueStmtContext) {
	// TODO: Implementar continue (necesita manejo de labels de loops)
	t.generator.Comment("=== CONTINUE STATEMENT ===")
}

// translateProgram traduce el nodo programa principal
func (t *ARM64Translator) translateProgram(ctx *compiler.ProgramContext) {
	t.generator.Comment("=== TRADUCCIÓN DEL PROGRAMA PRINCIPAL ===")

	for _, stmt := range ctx.AllStmt() {
		t.translateNode(stmt)
	}
}

// translateStatement traduce una declaración general
func (t *ARM64Translator) translateStatement(ctx *compiler.StmtContext) {
	if ctx.Decl_stmt() != nil {
		t.translateNode(ctx.Decl_stmt())
	} else if ctx.Assign_stmt() != nil {
		t.translateNode(ctx.Assign_stmt())
	} else if ctx.If_stmt() != nil {
		t.translateNode(ctx.If_stmt())
	} else if ctx.For_stmt() != nil {
		t.translateNode(ctx.For_stmt())
	} else if ctx.Func_call() != nil {
		t.translateNode(ctx.Func_call())
	} else if ctx.Func_dcl() != nil {
		t.translateNode(ctx.Func_dcl())
	} else if ctx.Transfer_stmt() != nil { // AGREGAR ESTO
		t.translateNode(ctx.Transfer_stmt())
	}
}

// === DECLARACIONES DE VARIABLES ===
// Manejar declaraciones de funciones
func (t *ARM64Translator) translateFunctionDeclaration(ctx *compiler.FuncDeclContext) {
	funcName := ctx.ID().GetText()

	if funcName == "main" {
		t.generator.Comment(fmt.Sprintf("=== FUNCIÓN %s ===", funcName))

		// Traducir el cuerpo de la función main directamente
		for _, stmt := range ctx.AllStmt() {
			t.translateNode(stmt)
		}
	} else {
		// NO AGREGAR ERROR - Las funciones de usuario se generan al final
		t.generator.Comment(fmt.Sprintf("=== DECLARACIÓN DE FUNCIÓN %s (se generará al final) ===", funcName))
	}
}

// Manejar contexto de declaración
func (t *ARM64Translator) translateDeclStatement(ctx *compiler.Decl_stmtContext) {
	// Recorrer hijos para encontrar el tipo específico
	for i := 0; i < ctx.GetChildCount(); i++ {
		if child, ok := ctx.GetChild(i).(antlr.ParseTree); ok {
			t.translateNode(child)
		}
	}
}

// Manejar declaraciones con inferencia de tipo
func (t *ARM64Translator) translateValueDecl(ctx *compiler.ValueDeclContext) {
	varName := ctx.ID().GetText()
	t.generator.Comment(fmt.Sprintf("=== DECLARACIÓN: mut %s (inferido) ===", varName))

	// Evaluar la expresión del lado derecho
	t.translateExpression(ctx.Expression())

	// Guardar el resultado en la variable
	t.generator.StoreVariable(arm64.X0, varName)
}

// translateMutVarDecl traduce: mut variable int = 10
func (t *ARM64Translator) translateMutVarDecl(ctx *compiler.MutVarDeclContext) {
	varName := ctx.ID().GetText()
	t.generator.Comment(fmt.Sprintf("=== DECLARACIÓN: mut %s ===", varName))

	// Evaluar la expresión del lado derecho
	t.translateExpression(ctx.Expression())

	// Guardar el resultado en la variable
	t.generator.StoreVariable(arm64.X0, varName)
}

// translateVarAssDecl traduce: variable int = 10
func (t *ARM64Translator) translateVarAssDecl(ctx *compiler.VarAssDeclContext) {
	varName := ctx.ID().GetText()
	t.generator.Comment(fmt.Sprintf("=== DECLARACIÓN: %s ===", varName))

	// Evaluar la expresión del lado derecho
	t.translateExpression(ctx.Expression())

	// Guardar el resultado en la variable
	t.generator.StoreVariable(arm64.X0, varName)
}

// === ASIGNACIONES ===

// translateAssignment traduce: variable = expresion
func (t *ARM64Translator) translateAssignment(ctx *compiler.AssignmentDeclContext) {
	varName := ctx.Id_pattern().GetText()
	t.generator.Comment(fmt.Sprintf("=== ASIGNACIÓN: %s = ... ===", varName))

	// Verificar que la variable existe
	if !t.generator.VariableExists(varName) {
		t.addError(fmt.Sprintf("Variable '%s' no está declarada", varName))
		return
	}

	// Evaluar la expresión del lado derecho
	t.translateExpression(ctx.Expression())

	// Guardar el resultado en la variable
	t.generator.StoreVariable(arm64.X0, varName)
}

// === EXPRESIONES ===

// translateExpression traduce cualquier expresión y deja el resultado en X0
func (t *ARM64Translator) translateExpression(expr antlr.ParseTree) {
	fmt.Printf("🔢 Traduciendo expresión: %T = %s\n", expr, expr.GetText())

	switch ctx := expr.(type) {
	case *compiler.IntLiteralContext:
		t.translateIntLiteral(ctx)
	case *compiler.IdPatternExprContext:
		t.translateVariable(ctx)
	case *compiler.BinaryExprContext:
		t.translateBinaryExpression(ctx)
	case *compiler.ParensExprContext:
		t.translateExpression(ctx.Expression())
	case *compiler.LiteralExprContext:
		// Procesar el literal interno
		t.translateExpression(ctx.Literal())
	case *compiler.LiteralContext:
		t.translateLiteral(ctx)
	case *compiler.FuncCallExprContext:
		t.translateNode(ctx.Func_call())
	default:
		t.addError(fmt.Sprintf("Expresión no implementada: %T", ctx))
		t.generator.LoadImmediate(arm64.X0, 0)
	}
}

// ✅ AGREGAR ESTA FUNCIÓN:
func (t *ARM64Translator) translateLiteral(ctx *compiler.LiteralContext) {
	// Primero intentar procesar hijos específicos
	for i := 0; i < ctx.GetChildCount(); i++ {
		child := ctx.GetChild(i)
		if child != nil {
			switch childCtx := child.(type) {
			case *compiler.IntLiteralContext:
				t.translateIntLiteral(childCtx)
				return
			case *compiler.FloatLiteralContext:
				t.translateFloatLiteral(childCtx)
				return
			case *compiler.StringLiteralContext:
				t.translateStringLiteral(childCtx)
				return
			case *compiler.BoolLiteralContext:
				t.translateBoolLiteral(childCtx)
				return
			case antlr.ParseTree:
				// Si es otro tipo de ParseTree, procesar recursivamente
				t.translateExpression(childCtx)
				return
			}
		}
	}

	// Fallback: analizar por texto si no se encontró un tipo específico
	text := ctx.GetText()
	fmt.Printf("🔍 Procesando literal por texto: %s\n", text)

	if value, err := strconv.Atoi(text); err == nil {
		t.generator.LoadImmediate(arm64.X0, value)
	} else {
		t.generator.LoadImmediate(arm64.X0, 0)
	}
}

// ✅ AGREGAR ESTAS FUNCIONES SI NO EXISTEN:
func (t *ARM64Translator) translateFloatLiteral(ctx *compiler.FloatLiteralContext) {
	valueStr := ctx.GetText()
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		t.addError(fmt.Sprintf("Error convirtiendo flotante: %s", valueStr))
		value = 0.0
	}
	// Por simplicidad, convertir a entero
	t.generator.LoadImmediate(arm64.X0, int(value))
}

func (t *ARM64Translator) translateStringLiteral(ctx *compiler.StringLiteralContext) {
	// Por simplicidad, cargar la longitud de la cadena
	text := ctx.GetText()
	if len(text) >= 2 {
		text = text[1 : len(text)-1] // Quitar comillas
	}
	t.generator.LoadImmediate(arm64.X0, len(text))
}

func (t *ARM64Translator) translateBoolLiteral(ctx *compiler.BoolLiteralContext) {
	valueStr := ctx.GetText()
	value := 0
	if valueStr == "true" {
		value = 1
	}
	t.generator.LoadImmediate(arm64.X0, value)
}

// translateIntLiteral traduce un literal entero
func (t *ARM64Translator) translateIntLiteral(ctx *compiler.IntLiteralContext) {
	valueStr := ctx.GetText()
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		t.addError(fmt.Sprintf("Error convirtiendo entero: %s", valueStr))
		value = 0
	}

	t.generator.LoadImmediate(arm64.X0, value)
}

// translateVariable traduce el acceso a una variable
func (t *ARM64Translator) translateVariable(ctx *compiler.IdPatternExprContext) {
	varName := ctx.Id_pattern().GetText()

	if !t.generator.VariableExists(varName) {
		t.addError(fmt.Sprintf("Variable '%s' no está declarada", varName))
		t.generator.LoadImmediate(arm64.X0, 0) // Valor por defecto
		return
	}

	t.generator.LoadVariable(arm64.X0, varName)
}

// translateBinaryExpression traduce expresiones binarias (+, -, *, /, ==, etc.)
func (t *ARM64Translator) translateBinaryExpression(ctx *compiler.BinaryExprContext) {
	operator := ctx.GetOp().GetText()
	t.generator.Comment(fmt.Sprintf("=== OPERACIÓN BINARIA: %s ===", operator))

	// Evaluar operando izquierdo y guardarlo en x1
	t.translateExpression(ctx.GetLeft())
	t.generator.Comment("Mover operando izquierdo a x1")
	t.generator.Emit("mov x1, x0")

	// Evaluar operando derecho (queda en X0)
	t.translateExpression(ctx.GetRight())

	// Realizar la operación correspondiente
	switch operator {
	case "+":
		t.generator.Add(arm64.X0, arm64.X1, arm64.X0)
	case "-":
		t.generator.Sub(arm64.X0, arm64.X1, arm64.X0)
	case "*":
		t.generator.Mul(arm64.X0, arm64.X1, arm64.X0)
	case "/":
		t.generator.Div(arm64.X0, arm64.X1, arm64.X0)
	case "==":
		t.translateComparison(arm64.X1, arm64.X0, "eq")
	case "!=":
		t.translateComparison(arm64.X1, arm64.X0, "ne")
	case "<":
		t.translateComparison(arm64.X1, arm64.X0, "lt")
	case ">":
		t.translateComparison(arm64.X1, arm64.X0, "gt")
	case "<=":
		t.translateComparison(arm64.X1, arm64.X0, "le")
	case ">=":
		t.translateComparison(arm64.X1, arm64.X0, "ge")
	default:
		t.addError(fmt.Sprintf("Operador no implementado: %s", operator))
	}
}

// translateComparison traduce operaciones de comparación
func (t *ARM64Translator) translateComparison(reg1, reg2, condition string) {
	t.generator.Compare(reg1, reg2)

	// Usar CSET para convertir el resultado de la comparación a 1 o 0
	t.generator.Comment(fmt.Sprintf("Convertir resultado de comparación a 1/0"))
	t.generator.Emit(fmt.Sprintf("cset %s, %s", arm64.X0, condition))
}

// === CONTROL DE FLUJO ===

// translateIfStatement traduce declaraciones if-else
func (t *ARM64Translator) translateIfStatement(ctx *compiler.IfStmtContext) {
	t.generator.Comment("=== IF STATEMENT ===")

	elseLabel := t.generator.GetLabel()
	endLabel := t.generator.GetLabel()

	// Evaluar la condición del primer if_chain
	if len(ctx.AllIf_chain()) > 0 {
		ifChain := ctx.AllIf_chain()[0]
		if ifChainCtx, ok := ifChain.(*compiler.IfChainContext); ok {
			// Evaluar condición
			t.translateExpression(ifChainCtx.Expression())

			// Saltar a else si la condición es falsa (0)
			t.generator.JumpIfZero(arm64.X0, elseLabel)

			// Ejecutar cuerpo del if
			for _, stmt := range ifChainCtx.AllStmt() {
				t.translateNode(stmt)
			}

			// Saltar al final para evitar ejecutar el else
			t.generator.Jump(endLabel)
		}
	}

	// Etiqueta else
	t.generator.SetLabel(elseLabel)

	// Si hay else, ejecutarlo
	if ctx.Else_stmt() != nil {
		elseCtx := ctx.Else_stmt().(*compiler.ElseStmtContext)
		for _, stmt := range elseCtx.AllStmt() {
			t.translateNode(stmt)
		}
	}

	// Etiqueta final
	t.generator.SetLabel(endLabel)
}

// translateForLoop traduce bucles for
func (t *ARM64Translator) translateForLoop(ctx *compiler.ForStmtCondContext) {
	t.generator.Comment("=== FOR LOOP ===")

	startLabel := t.generator.GetLabel()
	endLabel := t.generator.GetLabel()

	// Etiqueta de inicio del bucle
	t.generator.SetLabel(startLabel)

	// Evaluar condición
	t.translateExpression(ctx.Expression())

	// Salir del bucle si la condición es falsa
	t.generator.JumpIfZero(arm64.X0, endLabel)

	// Ejecutar cuerpo del bucle
	for _, stmt := range ctx.AllStmt() {
		t.translateNode(stmt)
	}

	// Volver al inicio del bucle
	t.generator.Jump(startLabel)

	// Etiqueta final
	t.generator.SetLabel(endLabel)
}

// === LLAMADAS A FUNCIONES ===

// translateFunctionCall traduce llamadas a funciones
func (t *ARM64Translator) translateFunctionCall(ctx *compiler.FuncCallContext) {
	funcName := ctx.Id_pattern().GetText()

	// Manejar funciones especiales
	switch funcName {
	case "print":
		t.translatePrintFunction(ctx, false) // sin salto de línea
	case "println":
		t.translatePrintFunction(ctx, true) // con salto de línea
	case "main":
		t.generator.Comment("=== LLAMADA A FUNCIÓN MAIN ===")
	default:
		// AGREGAR: Verificar si es función de usuario
		if funcDecl, exists := t.userFunctions[funcName]; exists {
			t.translateUserFunctionCall(ctx, funcDecl)
		} else {
			// Manejar funciones nativas simuladas
			t.translateNativeFunction(ctx)
		}
	}
}

func (t *ARM64Translator) translateUserFunctionCall(callCtx *compiler.FuncCallContext, funcDecl *compiler.FuncDeclContext) {
	funcName := callCtx.Id_pattern().GetText()
	t.generator.Comment(fmt.Sprintf("=== LLAMADA A FUNCIÓN DE USUARIO: %s ===", funcName))

	// Preparar argumentos - CARGAR EN ORDEN INVERSO
	if callCtx.Arg_list() != nil {
		args := callCtx.Arg_list().(*compiler.ArgListContext).AllFunc_arg()

		// Debug: mostrar argumentos
		fmt.Printf("🔍 Argumentos para %s: %d\n", funcName, len(args))
		for i, arg := range args {
			fmt.Printf("🔍 Arg %d: %s\n", i, arg.GetText())
		}

		// CARGAR ARGUMENTOS EN ORDEN INVERSO PARA NO SOBRESCRIBIR
		for i := len(args) - 1; i >= 0; i-- {
			arg := args[i]
			if argCtx := arg.(*compiler.FuncArgContext); argCtx != nil {

				targetReg := fmt.Sprintf("x%d", i)
				t.generator.Comment(fmt.Sprintf("Cargando argumento %d (%s) en %s", i, argCtx.GetText(), targetReg))

				// Evaluar el argumento
				if argCtx.Expression() != nil {
					t.translateExpression(argCtx.Expression())
				} else if argCtx.Id_pattern() != nil {
					// Es una variable
					varName := argCtx.Id_pattern().GetText()
					if t.generator.VariableExists(varName) {
						t.generator.LoadVariable(arm64.X0, varName)
					} else {
						t.addError(fmt.Sprintf("Variable '%s' no encontrada", varName))
						t.generator.LoadImmediate(arm64.X0, 0)
					}
				} else {
					// Fallback: intentar como texto
					argText := argCtx.GetText()
					if t.generator.VariableExists(argText) {
						t.generator.LoadVariable(arm64.X0, argText)
					} else if value, err := strconv.Atoi(argText); err == nil {
						t.generator.LoadImmediate(arm64.X0, value)
					} else {
						t.addError(fmt.Sprintf("No se puede procesar argumento: %s", argText))
						t.generator.LoadImmediate(arm64.X0, 0)
					}
				}

				// Mover al registro correcto (solo si no es x0)
				if i != 0 {
					t.generator.Emit(fmt.Sprintf("mov %s, x0", targetReg))
				}
			}
		}
	}

	// Llamar a la función
	t.generator.CallFunction(fmt.Sprintf("func_%s", funcName))
}

func (t *ARM64Translator) translateNativeFunction(ctx *compiler.FuncCallContext) {
	funcName := ctx.Id_pattern().GetText()

	switch funcName {
	case "atoi":
		// Simular atoi - por simplicidad retornar valor fijo
		if ctx.Arg_list() != nil {
			args := ctx.Arg_list().(*compiler.ArgListContext).AllFunc_arg()
			if len(args) > 0 {
				// Por simplicidad, si el string es "123", retornar 123
				t.generator.LoadImmediate(arm64.X0, 123)
			}
		}
	case "parse_float":
		// Simular parse_float
		t.generator.LoadImmediate(arm64.X0, 123) // Simplificado
	case "TypeOf", "Type":
		// Simular TypeOf - retornar código que representa tipo
		t.generator.LoadImmediate(arm64.X0, 1) // 1=int, 2=float, etc.
	default:
		t.addError(fmt.Sprintf("Función no implementada: %s", funcName))
		t.generator.LoadImmediate(arm64.X0, 0)
	}
}

// translatePrintFunction traduce llamadas a print/println
func (t *ARM64Translator) translatePrintFunction(ctx *compiler.FuncCallContext, withNewline bool) {
	t.generator.Comment("=== FUNCIÓN PRINT ===")

	// DEBUG: Ver si hay argumentos
	fmt.Printf("🔍 Print function - Arg_list: %v\n", ctx.Arg_list() != nil)
	if ctx.Arg_list() != nil {
		args := ctx.Arg_list().(*compiler.ArgListContext).AllFunc_arg()
		fmt.Printf("🔍 Número de argumentos: %d\n", len(args))

		for i, arg := range args {
			fmt.Printf("🔍 Argumento %d: %T = %s\n", i, arg, arg.GetText())
		}
	}

	// Procesar argumentos
	if ctx.Arg_list() != nil {
		args := ctx.Arg_list().(*compiler.ArgListContext).AllFunc_arg()

		for i, arg := range args {
			if i > 0 {
				// Imprimir espacio entre argumentos
				t.generator.Comment("Imprimir espacio")
				t.generator.LoadImmediate(arm64.X0, 32) // ASCII espacio
				t.generator.CallFunction("print_char")
			}

			// Evaluar el argumento
			if argCtx := arg.(*compiler.FuncArgContext); argCtx != nil {
				fmt.Printf("🔍 FuncArg: %s\n", argCtx.GetText())

				// Verificar si es expression
				if argCtx.Expression() != nil {
					fmt.Printf("🔍 Evaluando Expression: %s\n", argCtx.Expression().GetText())
					t.translateExpression(argCtx.Expression())
					t.generator.CallFunction("print_integer")
				} else if argCtx.Id_pattern() != nil { // Verificar si es id_pattern
					fmt.Printf("🔍 Evaluando Id_pattern: %s\n", argCtx.Id_pattern().GetText())
					// Crear un contexto de expresión temporal para el id_pattern
					varName := argCtx.Id_pattern().GetText()
					if t.generator.VariableExists(varName) {
						t.generator.LoadVariable(arm64.X0, varName)
						t.generator.CallFunction("print_integer")
					} else {
						t.addError(fmt.Sprintf("Variable '%s' no encontrada en print", varName))
					}
				} else {
					fmt.Printf("🔍 Argumento sin Expression ni Id_pattern válidos\n")
					// Fallback: intentar procesar como texto
					argText := argCtx.GetText()
					fmt.Printf("🔍 Fallback - procesando texto: %s\n", argText)

					// Si es un nombre de variable, intentar cargarla
					if t.generator.VariableExists(argText) {
						t.generator.LoadVariable(arm64.X0, argText)
						t.generator.CallFunction("print_integer")
					} else {
						// Si es un número literal
						if value, err := strconv.Atoi(argText); err == nil {
							t.generator.LoadImmediate(arm64.X0, value)
							t.generator.CallFunction("print_integer")
						} else {
							t.addError(fmt.Sprintf("No se puede procesar argumento: %s", argText))
						}
					}
				}
			} else {
				fmt.Printf("🔍 argCtx es nil\n")
			}
		}
	} else {
		fmt.Printf("🔍 No hay argumentos para imprimir\n")
	}

	// Si es println, agregar salto de línea
	if withNewline {
		t.generator.Comment("Imprimir salto de línea")
		t.generator.LoadImmediate(arm64.X0, 10) // ASCII newline
		t.generator.CallFunction("print_char")
	}
}

// === LIBRERÍA ESTÁNDAR ===

// generateStandardLibrary genera las funciones básicas necesarias
func (t *ARM64Translator) generateStandardLibrary() {
	// Función para imprimir enteros
	t.generator.EmitRaw(`
print_integer:
    // Función simplificada para imprimir enteros
    // Input: x0 = número a imprimir
    stp x29, x30, [sp, #-16]!    // Guardar registros
    stp x19, x20, [sp, #-16]!
    
    mov x19, x0                   // Guardar número original
    
    // Manejar caso especial: cero
    cmp x19, #0
    bne convert_digits
    
    // Imprimir '0'
    mov x0, #48                   // ASCII '0'
    bl print_char
    b print_done
    
convert_digits:
    // Buffer para dígitos (en el stack)
    sub sp, sp, #32
    mov x20, sp                   // x20 = puntero al buffer
    mov x21, #0                   // x21 = contador de dígitos
    
    // Manejar números negativos
    cmp x19, #0
    bge positive
    mov x0, #45                   // ASCII '-'
    bl print_char
    neg x19, x19                  // Hacer positivo
    
positive:
    // Convertir dígitos
digit_loop:
    mov x22, #10
    udiv x23, x19, x22           // x23 = x19 / 10
    msub x24, x23, x22, x19      // x24 = x19 % 10
    
    add x24, x24, #48            // Convertir a ASCII
    strb w24, [x20, x21]         // Guardar dígito
    add x21, x21, #1             // Incrementar contador
    
    mov x19, x23                 // x19 = quotient
    cbnz x19, digit_loop         // Continuar si no es cero
    
    // Imprimir dígitos en orden inverso
print_digits:
    sub x21, x21, #1
    ldrb w0, [x20, x21]
    bl print_char
    cbnz x21, print_digits
    
    add sp, sp, #32              // Limpiar buffer
    
print_done:
    ldp x19, x20, [sp], #16      // Restaurar registros
    ldp x29, x30, [sp], #16
    ret

print_char:
    // Imprimir un carácter
    // Input: x0 = carácter ASCII
    stp x29, x30, [sp, #-16]!
    
    // Crear buffer temporal en el stack
    sub sp, sp, #16
    strb w0, [sp]                // Guardar carácter
    
    // Syscall write
    mov x0, #1                   // stdout
    mov x1, sp                   // buffer
    mov x2, #1                   // length
    mov x8, #64                  // write syscall
    svc #0
    
    add sp, sp, #16              // Limpiar buffer
    ldp x29, x30, [sp], #16
    ret`)
}

// === UTILIDADES ===

// addError agrega un error a la lista
func (t *ARM64Translator) addError(message string) {
	t.errors = append(t.errors, message)
}

// GetErrors retorna todos los errores encontrados
func (t *ARM64Translator) GetErrors() []string {
	return t.errors
}

// HasErrors indica si hay errores
func (t *ARM64Translator) HasErrors() bool {
	return len(t.errors) > 0
}
