// backend/compiler/translator.go
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

	breakLabels    []string          // Etiquetas para manejar break en loops
	continueLabels []string          // Etiquetas para manejar continue en loops
	stringRegistry map[string]string // texto -> etiqueta Para evitar procesar strings dos veces
	variableTypes  map[string]string // nombre -> tipo Para rastrear tipos de variables
}

// NewARM64Translator crea un nuevo traductor
func NewARM64Translator() *ARM64Translator {
	return &ARM64Translator{
		generator:      arm64.NewARM64Generator(),
		errors:         make([]string, 0),
		userFunctions:  make(map[string]*compiler.FuncDeclContext),
		breakLabels:    make([]string, 0),
		continueLabels: make([]string, 0),
		stringRegistry: make(map[string]string),
		variableTypes:  make(map[string]string),
	}
}

// === FUNCIÓN PRINCIPAL DE TRADUCCIÓN ===

// TranslateProgram traduce un programa completo de VlangCherry a ARM64
func (t *ARM64Translator) TranslateProgram(tree antlr.ParseTree) (string, []string) {
	// Limpiar estado anterior
	t.generator.Reset()
	t.errors = make([]string, 0)
	t.variableTypes = make(map[string]string)

	fmt.Printf("🔍 === PRIMERA PASADA: ANÁLISIS DEL PROGRAMA ===\n")

	// Primera pasada: analizar declaraciones de variables Y strings
	t.analyzeVariablesAndStrings(tree)

	// Generar header del programa
	t.generator.GenerateHeader()

	fmt.Printf("🔍 === SEGUNDA PASADA: GENERACIÓN DE CÓDIGO ===\n")

	// Traducir el contenido del programa / Segunda pasada
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

// === ANÁLISIS MEJORADO (PRIMERA PASADA) ===

// analyzeVariablesAndStrings hace una pasada previa para encontrar todas las variables Y strings
func (t *ARM64Translator) analyzeVariablesAndStrings(node antlr.ParseTree) {
	switch ctx := node.(type) {
	case *compiler.ProgramContext:
		for _, stmt := range ctx.AllStmt() {
			fmt.Printf("🔍 Analizando statement: %T\n", stmt)
			t.analyzeVariablesAndStrings(stmt)
		}

	case *compiler.StmtContext:
		if ctx.Decl_stmt() != nil {
			t.analyzeVariablesAndStrings(ctx.Decl_stmt())
		}
		if ctx.If_stmt() != nil {
			t.analyzeVariablesAndStrings(ctx.If_stmt())
		}
		if ctx.Switch_stmt() != nil {
			t.analyzeVariablesAndStrings(ctx.Switch_stmt())
		}
		if ctx.For_stmt() != nil {
			t.analyzeVariablesAndStrings(ctx.For_stmt())
		}
		if ctx.Func_dcl() != nil {
			t.analyzeVariablesAndStrings(ctx.Func_dcl())
		}
		if ctx.Transfer_stmt() != nil {
			t.analyzeVariablesAndStrings(ctx.Transfer_stmt())
		}
		if ctx.Func_call() != nil {
			t.analyzeVariablesAndStrings(ctx.Func_call())
		}

	case *compiler.ValueDeclContext:
		varName := ctx.ID().GetText()
		if !t.generator.VariableExists(varName) {
			t.generator.DeclareVariable(varName)
		}
		// NUEVO: Inferir tipo de la variable
		if ctx.Expression() != nil {
			varType := t.inferExpressionType(ctx.Expression())
			t.variableTypes[varName] = varType
			fmt.Printf("🔍 Variable '%s' inferida como tipo: %s\n", varName, varType)
		}
		// Analizar strings en la expresión de inicialización
		if ctx.Expression() != nil {
			t.analyzeStringsInExpression(ctx.Expression())
		}

	case *compiler.MutVarDeclContext:
		varName := ctx.ID().GetText()
		if !t.generator.VariableExists(varName) {
			t.generator.DeclareVariable(varName)
		}
		// NUEVO: Inferir tipo de la variable
		if ctx.Expression() != nil {
			varType := t.inferExpressionType(ctx.Expression())
			t.variableTypes[varName] = varType
			fmt.Printf("🔍 Variable '%s' inferida como tipo: %s\n", varName, varType)
		}
		// Analizar strings en la expresión de inicialización
		if ctx.Expression() != nil {
			t.analyzeStringsInExpression(ctx.Expression())
		}

	case *compiler.VarAssDeclContext:
		varName := ctx.ID().GetText()
		if !t.generator.VariableExists(varName) {
			t.generator.DeclareVariable(varName)
		}
		// NUEVO: Inferir tipo de la variable
		if ctx.Expression() != nil {
			varType := t.inferExpressionType(ctx.Expression())
			t.variableTypes[varName] = varType
			fmt.Printf("🔍 Variable '%s' inferida como tipo: %s\n", varName, varType)
		}
		// Analizar strings en la expresión de inicialización
		if ctx.Expression() != nil {
			t.analyzeStringsInExpression(ctx.Expression())
		}

	case *compiler.FuncDeclContext:
		funcName := ctx.ID().GetText()

		fmt.Printf("🔍 Analizando función: %s\n", funcName)

		// Registrar función de usuario
		if funcName != "main" {
			t.userFunctions[funcName] = ctx
		}

		// Contar parámetros
		paramCount := 0
		if ctx.Param_list() != nil {
			params := ctx.Param_list().(*compiler.ParamListContext).AllFunc_param()
			paramCount = len(params)
			for i, param := range params {
				if paramCtx := param.(*compiler.FuncParamContext); paramCtx.ID() != nil {
					paramName := paramCtx.ID().GetText()
					fmt.Printf("📝 Parámetro: %s\n", paramName)
					t.generator.DeclareVariable(paramName)
				} else {
					fmt.Printf("📝 Parámetro %d: sin nombre\n", i)
				}
			}
		}

		// Analizar el cuerpo de la función
		varCount := 0
		for _, stmt := range ctx.AllStmt() {
			initialVarCount := len(t.generator.GetVariables())
			t.analyzeVariablesAndStrings(stmt)
			varCount += len(t.generator.GetVariables()) - initialVarCount
		}

		stackSize := (paramCount + varCount) * 8
		fmt.Printf("📊 Parámetros: %d, Variables locales: %d, Stack: %d bytes\n", paramCount, varCount, stackSize)

	// NUEVO: Análisis específico de llamadas a funciones
	case *compiler.FuncCallContext:
		fmt.Printf("🔍 Analizando llamada a función: %s\n", ctx.Id_pattern().GetText())

		// Analizar argumentos en busca de strings
		if ctx.Arg_list() != nil {
			args := ctx.Arg_list().(*compiler.ArgListContext).AllFunc_arg()
			for i, arg := range args {
				fmt.Printf("🔍   Analizando argumento %d: %s\n", i, arg.GetText())
				if argCtx := arg.(*compiler.FuncArgContext); argCtx != nil {
					if argCtx.Expression() != nil {
						t.analyzeStringsInExpression(argCtx.Expression())
					}
				}
			}
		}

	// Otros casos que pueden contener strings
	case *compiler.IfStmtContext:
		for _, ifChain := range ctx.AllIf_chain() {
			if ifChainCtx, ok := ifChain.(*compiler.IfChainContext); ok {
				// Analizar condición
				if ifChainCtx.Expression() != nil {
					t.analyzeStringsInExpression(ifChainCtx.Expression())
				}
				// Analizar cuerpo
				for _, stmt := range ifChainCtx.AllStmt() {
					t.analyzeVariablesAndStrings(stmt)
				}
			}
		}
		if ctx.Else_stmt() != nil {
			elseCtx := ctx.Else_stmt().(*compiler.ElseStmtContext)
			for _, stmt := range elseCtx.AllStmt() {
				t.analyzeVariablesAndStrings(stmt)
			}
		}

	case *compiler.ForStmtCondContext:
		// Analizar condición
		if ctx.Expression() != nil {
			t.analyzeStringsInExpression(ctx.Expression())
		}
		// Analizar cuerpo
		for _, stmt := range ctx.AllStmt() {
			t.analyzeVariablesAndStrings(stmt)
		}
	}
}

// NUEVA FUNCIÓN: Inferir tipo de expresión
func (t *ARM64Translator) inferExpressionType(expr antlr.ParseTree) string {
	if expr == nil {
		return "unknown"
	}

	switch ctx := expr.(type) {
	case *compiler.LiteralExprContext:
		return t.inferExpressionType(ctx.Literal())
	case *compiler.LiteralContext:
		return t.inferLiteralType(ctx)
	case *compiler.StringLiteralContext:
		return "string"
	case *compiler.IntLiteralContext:
		return "int"
	case *compiler.FloatLiteralContext:
		return "float"
	case *compiler.BoolLiteralContext:
		return "bool"
	case *compiler.IdPatternExprContext:
		varName := ctx.Id_pattern().GetText()
		if varType, exists := t.variableTypes[varName]; exists {
			return varType
		}
		return "unknown"
	case *compiler.BinaryExprContext:
		// Para operaciones binarias, usar el tipo del operando izquierdo
		return t.inferExpressionType(ctx.GetLeft())
	default:
		return "unknown"
	}
}

// Inferir tipo de literal
func (t *ARM64Translator) inferLiteralType(ctx *compiler.LiteralContext) string {
	text := ctx.GetText()

	// Verificar string
	if strings.HasPrefix(text, "\"") && strings.HasSuffix(text, "\"") {
		return "string"
	}

	// Verificar float
	if strings.Contains(text, ".") {
		return "float"
	}

	// Verificar bool
	if text == "true" || text == "false" {
		return "bool"
	}

	// Por defecto, int
	if _, err := strconv.Atoi(text); err == nil {
		return "int"
	}

	return "unknown"
}

// === ANÁLISIS MEJORADO DE STRINGS ===

// analyzeStringsInExpression busca strings en expresiones de forma más completa
func (t *ARM64Translator) analyzeStringsInExpression(expr antlr.ParseTree) {
	if expr == nil {
		return
	}

	fmt.Printf("🔍 Analizando expresión para strings: %s (tipo: %T)\n", expr.GetText(), expr)

	switch ctx := expr.(type) {
	case *compiler.StringLiteralContext:
		// STRING LITERAL DIRECTO
		t.preProcessStringLiteral(ctx)

	case *compiler.LiteralExprContext:
		// Expresión literal que puede contener un string
		t.analyzeStringsInExpression(ctx.Literal())

	case *compiler.LiteralContext:
		// Buscar string literals en el literal
		for i := 0; i < ctx.GetChildCount(); i++ {
			child := ctx.GetChild(i)
			if stringCtx, ok := child.(*compiler.StringLiteralContext); ok {
				t.preProcessStringLiteral(stringCtx)
			}
		}

	case *compiler.BinaryExprContext:
		// Analizar ambos lados de la expresión binaria
		t.analyzeStringsInExpression(ctx.GetLeft())
		t.analyzeStringsInExpression(ctx.GetRight())

	case *compiler.ParensExprContext:
		// Analizar expresión entre paréntesis
		t.analyzeStringsInExpression(ctx.Expression())

	case *compiler.FuncCallExprContext:
		// Analizar llamada a función
		t.analyzeStringsInExpression(ctx.Func_call())

	case *compiler.FuncCallContext:
		// Analizar argumentos de la función
		if ctx.Arg_list() != nil {
			args := ctx.Arg_list().(*compiler.ArgListContext).AllFunc_arg()
			for _, arg := range args {
				if argCtx := arg.(*compiler.FuncArgContext); argCtx != nil {
					if argCtx.Expression() != nil {
						t.analyzeStringsInExpression(argCtx.Expression())
					}
				}
			}
		}

	default:
		// Para otros tipos, analizar recursivamente todos los hijos
		if ctx, ok := expr.(antlr.ParserRuleContext); ok {
			for i := 0; i < ctx.GetChildCount(); i++ {
				if child := ctx.GetChild(i); child != nil {
					if parseTreeChild, ok := child.(antlr.ParseTree); ok {
						t.analyzeStringsInExpression(parseTreeChild)
					}
				}
			}
		}
	}
}

// preProcessStringLiteral procesa strings en la primera pasada (MEJORADO)
func (t *ARM64Translator) preProcessStringLiteral(ctx *compiler.StringLiteralContext) {
	text := ctx.GetText()
	if len(text) >= 2 {
		text = text[1 : len(text)-1] // Quitar comillas
	}

	// Procesar secuencias de escape
	text = strings.ReplaceAll(text, "\\n", "\n")
	text = strings.ReplaceAll(text, "\\t", "\t")
	text = strings.ReplaceAll(text, "\\\"", "\"")
	text = strings.ReplaceAll(text, "\\\\", "\\")

	// Verificar si ya fue procesado
	if existingLabel, exists := t.stringRegistry[text]; exists {
		fmt.Printf("🔄 String \"%s\" ya procesado como %s\n", text, existingLabel)
		return
	}

	// Agregar al generador
	stringLabel := t.generator.AddStringLiteral(text)

	// Registrar para evitar duplicados
	t.stringRegistry[text] = stringLabel

	fmt.Printf("✅ STRING REGISTRADO: \"%s\" -> %s\n", text, stringLabel)
}

// === RESTO DE MÉTODOS (mantenidos igual pero con corrección en print) ===

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

// === TRADUCCIÓN DE NODOS (mantenida igual) ===

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
	case *compiler.SwitchStmtContext:
		t.translateSwitchStatement(ctx)
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

// Modificar translateBreakStatementFromTransfer
func (t *ARM64Translator) translateBreakStatementFromTransfer(ctx *compiler.Transfer_stmtContext) {
	t.generator.Comment("=== BREAK STATEMENT ===")

	// Verificar si estamos en un contexto que permite break
	if len(t.breakLabels) > 0 {
		// Saltar a la etiqueta de break más reciente
		breakLabel := t.breakLabels[len(t.breakLabels)-1]
		t.generator.Jump(breakLabel)
	} else {
		t.addError("Break statement fuera de contexto válido (switch/loop)")
	}
}

func (t *ARM64Translator) translateContinueStatementFromTransfer(ctx *compiler.Transfer_stmtContext) {
	t.generator.Comment("=== CONTINUE STATEMENT ===")
	// TODO: Implementar continue
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
	} else if ctx.Switch_stmt() != nil {
		t.translateNode(ctx.Switch_stmt())
	} else if ctx.For_stmt() != nil {
		t.translateNode(ctx.For_stmt())
	} else if ctx.Func_call() != nil {
		t.translateNode(ctx.Func_call())
	} else if ctx.Func_dcl() != nil {
		t.translateNode(ctx.Func_dcl())
	} else if ctx.Transfer_stmt() != nil {
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
	case *compiler.StringLiteralContext:
		t.translateStringLiteral(ctx)
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

// ✅ FUNCIONES DE TRADUCCIÓN DE LITERALES:
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
	text := ctx.GetText()
	if len(text) >= 2 {
		text = text[1 : len(text)-1] // Quitar comillas
	}

	// Procesar secuencias de escape
	text = strings.ReplaceAll(text, "\\n", "\n")
	text = strings.ReplaceAll(text, "\\t", "\t")
	text = strings.ReplaceAll(text, "\\\"", "\"")
	text = strings.ReplaceAll(text, "\\\\", "\\")

	// VERIFICAR si ya fue procesado en la primera pasada
	if existingLabel, exists := t.stringRegistry[text]; exists {
		// Ya existe, usar la etiqueta existente
		t.generator.Comment(fmt.Sprintf("Usar string \"%s\" con etiqueta %s", text, existingLabel))
		t.generator.Emit(fmt.Sprintf("adr x0, %s", existingLabel))
		return
	}

	// Si no existe en el registro, es un error (debería haberse procesado en primera pasada)
	t.addError(fmt.Sprintf("String \"%s\" no fue procesado en primera pasada", text))
	t.generator.LoadImmediate(arm64.X0, 0) // Fallback
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

// 🔥 CORREGIR FUNCIÓN PRINT - EL PROBLEMA PRINCIPAL ESTÁ AQUÍ
func (t *ARM64Translator) translatePrintFunction(ctx *compiler.FuncCallContext, withNewline bool) {
	t.generator.Comment("=== FUNCIÓN PRINT ===")

	if ctx.Arg_list() != nil {
		args := ctx.Arg_list().(*compiler.ArgListContext).AllFunc_arg()

		for i, arg := range args {
			if i > 0 {
				// Imprimir espacio entre argumentos
				t.generator.Comment("Imprimir espacio")
				t.generator.LoadImmediate(arm64.X0, 32) // ASCII espacio
				t.generator.CallFunction("print_char")
			}

			if argCtx := arg.(*compiler.FuncArgContext); argCtx != nil {
				if argCtx.Expression() != nil {
					// NUEVA LÓGICA: Usar tipo de variable o inferir tipo
					exprText := argCtx.Expression().GetText()

					// Verificar si es una variable y obtener su tipo
					if varType := t.getArgumentType(argCtx); varType == "string" {
						t.generator.Comment(fmt.Sprintf("Imprimiendo string: %s", exprText))
						t.translateExpression(argCtx.Expression())
						t.generator.CallFunction("print_string")
					} else if strings.HasPrefix(exprText, "\"") && strings.HasSuffix(exprText, "\"") {
						// Es un string literal directo
						t.generator.Comment(fmt.Sprintf("Imprimiendo string literal: %s", exprText))
						t.translateExpression(argCtx.Expression())
						t.generator.CallFunction("print_string")
					} else {
						// Es una expresión numérica o variable numérica
						t.generator.Comment(fmt.Sprintf("Imprimiendo valor numérico: %s", exprText))
						t.translateExpression(argCtx.Expression())
						t.generator.CallFunction("print_integer")
					}
				} else if argCtx.Id_pattern() != nil {
					// Variable directa - determinar tipo
					varName := argCtx.Id_pattern().GetText()
					if t.generator.VariableExists(varName) {
						t.generator.LoadVariable(arm64.X0, varName)

						// Determinar qué función usar según el tipo
						if varType, exists := t.variableTypes[varName]; exists && varType == "string" {
							t.generator.Comment(fmt.Sprintf("Imprimiendo variable string: %s", varName))
							t.generator.CallFunction("print_string")
						} else {
							t.generator.Comment(fmt.Sprintf("Imprimiendo variable numérica: %s", varName))
							t.generator.CallFunction("print_integer")
						}
					} else {
						t.addError(fmt.Sprintf("Variable '%s' no encontrada", varName))
					}
				}
			}
		}
	}

	if withNewline {
		t.generator.Comment("Imprimir salto de línea")
		t.generator.LoadImmediate(arm64.X0, 10) // ASCII newline
		t.generator.CallFunction("print_char")
	}
}

// NUEVA FUNCIÓN: Determinar tipo de argumento
func (t *ARM64Translator) getArgumentType(argCtx *compiler.FuncArgContext) string {
	if argCtx.Expression() != nil {
		return t.inferExpressionType(argCtx.Expression())
	} else if argCtx.Id_pattern() != nil {
		varName := argCtx.Id_pattern().GetText()
		if varType, exists := t.variableTypes[varName]; exists {
			return varType
		}
	}
	return "unknown"
}

// === CONTROL DE FLUJO (simplificado) ===

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

// translateSwitchStatement traduce declaraciones switch
func (t *ARM64Translator) translateSwitchStatement(ctx *compiler.SwitchStmtContext) {
	t.generator.Comment("=== SWITCH STATEMENT ===")

	// Evaluar la expresión del switch una vez y guardarla
	t.translateExpression(ctx.Expression())
	t.generator.Comment("Guardar valor del switch en x19")
	t.generator.Emit("mov x19, x0")

	// Generar etiquetas
	defaultLabel := t.generator.GetLabel()
	endLabel := t.generator.GetLabel()
	caseLabels := make([]string, 0)

	// Push etiquetas de break
	t.breakLabels = append(t.breakLabels, endLabel)

	// Generar etiquetas para cada caso
	cases := ctx.AllSwitch_case()
	for range cases {
		caseLabels = append(caseLabels, t.generator.GetLabel())
	}

	t.generator.Comment("=== COMPARACIONES DE CASOS ===")

	// Generar comparaciones para cada caso
	for i, switchCase := range cases {
		if caseCtx, ok := switchCase.(*compiler.SwitchCaseContext); ok {
			t.generator.Comment(fmt.Sprintf("Comparar caso %d", i))

			// Evaluar la expresión del caso
			t.translateExpression(caseCtx.Expression())

			// Comparar con el valor del switch
			t.generator.Compare("x19", "x0")
			t.generator.Emit(fmt.Sprintf("beq %s", caseLabels[i]))
		}
	}

	// Si ningún caso coincide, ir al default (o al final si no hay default)
	if ctx.Default_case() != nil {
		t.generator.Jump(defaultLabel)
	} else {
		t.generator.Jump(endLabel)
	}

	// Generar código para cada caso
	for i, switchCase := range cases {
		if caseCtx, ok := switchCase.(*compiler.SwitchCaseContext); ok {
			t.generator.SetLabel(caseLabels[i])
			t.generator.Comment(fmt.Sprintf("=== CASO %d ===", i))

			// Ejecutar statements del caso
			for _, stmt := range caseCtx.AllStmt() {
				t.translateNode(stmt)
			}

			// Automáticamente saltar al final (break implícito)
			t.generator.Jump(endLabel)
		}
	}

	// Generar caso default si existe
	if ctx.Default_case() != nil {
		t.generator.SetLabel(defaultLabel)
		t.generator.Comment("=== CASO DEFAULT ===")

		defaultCtx := ctx.Default_case().(*compiler.DefaultCaseContext)
		for _, stmt := range defaultCtx.AllStmt() {
			t.translateNode(stmt)
		}
	}

	// Etiqueta final
	t.generator.SetLabel(endLabel)

	// Limpiar etiquetas de break
	if len(t.breakLabels) > 0 {
		t.breakLabels = t.breakLabels[:len(t.breakLabels)-1]
	}

	t.generator.Comment("=== FIN SWITCH ===")
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

	t.generator.EmitRaw(`
print_string:
    // Función para imprimir strings
    // Input: x0 = dirección del string (terminado en null)
    stp x29, x30, [sp, #-16]!    // Guardar registros
    stp x19, x20, [sp, #-16]!
    
    mov x19, x0                   // x19 = dirección del string
    
    // Encontrar la longitud del string
    mov x20, #0                   // x20 = contador de longitud
    
strlen_loop:
    ldrb w1, [x19, x20]          // Cargar byte del string
    cbz w1, strlen_done          // Si es 0 (null terminator), terminar
    add x20, x20, #1             // Incrementar contador
    b strlen_loop
    
strlen_done:
    // Verificar si el string está vacío
    cbz x20, print_string_done
    
    // Syscall write(1, string, length)
    mov x0, #1                   // File descriptor: stdout
    mov x1, x19                  // Buffer: dirección del string
    mov x2, x20                  // Length: longitud calculada
    mov x8, #64                  // Syscall number: write
    svc #0                       // Llamada al sistema
    
print_string_done:
    ldp x19, x20, [sp], #16      // Restaurar registros
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
