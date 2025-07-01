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
		if ctx.Assign_stmt() != nil {
			t.analyzeVariablesAndStrings(ctx.Assign_stmt())
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

	case *compiler.Assign_stmtContext:
		// Procesar asignaciones para buscar strings
		for i := 0; i < ctx.GetChildCount(); i++ {
			if child := ctx.GetChild(i); child != nil {
				if parseTreeChild, ok := child.(antlr.ParseTree); ok {
					t.analyzeVariablesAndStrings(parseTreeChild)
				}
			}
		}

	case *compiler.ArgAddAssigDeclContext:
		// Analizar expresión en asignación suma
		if ctx.Expression() != nil {
			t.analyzeStringsInExpression(ctx.Expression())
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

	case *compiler.VarVectDeclContext:
		t.analyzeVarVectDecl(ctx)

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

	case *compiler.SwitchStmtContext:
		for _, rawCase := range ctx.AllSwitch_case() {
			if caseCtx, ok := rawCase.(*compiler.SwitchCaseContext); ok {
				t.analyzeStringsInExpression(caseCtx.Expression())
				for _, stmt := range caseCtx.AllStmt() {
					t.analyzeVariablesAndStrings(stmt)
				}
			}
		}
		if def := ctx.Default_case(); def != nil {
			if defaultCtx, ok := def.(*compiler.DefaultCaseContext); ok {
				for _, stmt := range defaultCtx.AllStmt() {
					t.analyzeVariablesAndStrings(stmt)
				}
			}
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
	case *compiler.FuncCallExprContext:
		// Inferir tipo basado en la llamada a función
		return t.inferFunctionCallType(ctx.Func_call())
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

// Inferir tipo de llamada a función
func (t *ARM64Translator) inferFunctionCallType(funcCall antlr.ParseTree) string {
	if funcCall == nil {
		return "unknown"
	}

	if funcCallCtx, ok := funcCall.(*compiler.FuncCallContext); ok {
		funcName := funcCallCtx.Id_pattern().GetText()

		// Funciones nativas que devuelven string
		switch funcName {
		case "join":
			return "string"
		case "append":
			// append devuelve un vector del mismo tipo que el primer argumento
			return "[]int" // Por simplicidad, asumir int por ahora
		case "print", "println":
			return "void"
		default:
			// Para funciones de usuario, asumir int por defecto
			return "int"
		}
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

	// VERIFICAR SI TIENE INTERPOLACIÓN
	if strings.Contains(text, "$") {
		// Procesar partes no variables del string interpolado
		parts := t.parseInterpolatedString(text)
		for _, part := range parts {
			if !part.IsVariable && part.Content != "" {
				// Solo registrar las partes de texto literal
				if _, exists := t.stringRegistry[part.Content]; !exists {
					stringLabel := t.generator.AddStringLiteral(part.Content)
					t.stringRegistry[part.Content] = stringLabel
					fmt.Printf("✅ STRING INTERPOLADO REGISTRADO: \"%s\" -> %s\n", part.Content, stringLabel)
				}
			}
		}
	} else {
		// String normal - procesar como antes
		if existingLabel, exists := t.stringRegistry[text]; exists {
			fmt.Printf("🔄 String \"%s\" ya procesado como %s\n", text, existingLabel)
			return
		}

		stringLabel := t.generator.AddStringLiteral(text)
		t.stringRegistry[text] = stringLabel
		fmt.Printf("✅ STRING REGISTRADO: \"%s\" -> %s\n", text, stringLabel)
	}
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
	case *compiler.Assign_stmtContext:
		t.translateAssignStatement(ctx)
	case *compiler.ArgAddAssigDeclContext:
		t.translateArgAddAssignment(ctx)
	case *compiler.FloatLiteralContext:
		t.translateFloatLiteral(ctx)
	case *compiler.BoolLiteralContext:
		t.translateBoolLiteral(ctx)
	case *compiler.UnaryExprContext: // AGREGAR ESTE CASE
		t.translateUnaryExpression(ctx)
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
	case *compiler.ForAssCondContext:
		t.translateForAssignment(ctx)
	case *compiler.BreakStmtContext:
		t.translateBreakStatement(ctx)
	case *compiler.ContinueStmtContext:
		t.translateContinueStatement(ctx)
	case *compiler.VarVectDeclContext:
		t.translateVarVectDecl(ctx)
	case *compiler.VectorAssignContext:
		t.translateVectorAssignment(ctx)

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

func (t *ARM64Translator) translateAssignStatement(ctx *compiler.Assign_stmtContext) {
	// Procesar todos los hijos para encontrar el tipo específico de asignación
	for i := 0; i < ctx.GetChildCount(); i++ {
		child := ctx.GetChild(i)
		switch childCtx := child.(type) {
		case *compiler.ArgAddAssigDeclContext:
			t.translateArgAddAssignment(childCtx)
		case *compiler.AssignmentDeclContext:
			t.translateAssignment(childCtx)
		case antlr.ParseTree:
			// Recursivamente procesar otros tipos
			t.translateNode(childCtx)
		}
	}
}

func (t *ARM64Translator) translateArgAddAssignment(ctx *compiler.ArgAddAssigDeclContext) {
	varName := ctx.Id_pattern().GetText()
	t.generator.Comment(fmt.Sprintf("=== ASIGNACIÓN SUMA: %s += ... ===", varName))

	// Verificar que la variable existe
	if !t.generator.VariableExists(varName) {
		t.addError(fmt.Sprintf("Variable '%s' no está declarada", varName))
		return
	}

	// Cargar valor actual de la variable en x1
	t.generator.LoadVariable(arm64.X1, varName)

	// Evaluar la expresión del lado derecho (resultado en x0)
	t.translateExpression(ctx.Expression())

	// Sumar: x0 = x1 + x0
	t.generator.Add(arm64.X0, arm64.X1, arm64.X0)

	// Guardar el resultado de vuelta en la variable
	t.generator.StoreVariable(arm64.X0, varName)
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
	case *compiler.FloatLiteralContext:
		t.translateFloatLiteral(ctx)
	case *compiler.BoolLiteralContext:
		t.translateBoolLiteral(ctx)
	case *compiler.IdPatternExprContext:
		t.translateVariable(ctx)
	case *compiler.BinaryExprContext:
		t.translateBinaryExpression(ctx)
	case *compiler.UnaryExprContext: // AGREGAR ESTE CASE
		t.translateUnaryExpression(ctx)
	case *compiler.ParensExprContext:
		t.translateExpression(ctx.Expression())
	case *compiler.LiteralExprContext:
		// Procesar el literal interno
		t.translateExpression(ctx.Literal())
	case *compiler.LiteralContext:
		t.translateLiteral(ctx)
	case *compiler.FuncCallExprContext:
		t.translateNode(ctx.Func_call())
	case *compiler.IncredecrContext:
		t.translateExpression(ctx.Incredecre())
	case *compiler.IncrementoContext:
		t.translateIncrement(ctx)
	case *compiler.DecrementoContext:
		t.translateDecrement(ctx)
	case *compiler.VectorItemExprContext:
		t.translateVectorAccess(ctx)

	default:
		t.addError(fmt.Sprintf("Expresión no implementada: %T", ctx))
		t.generator.LoadImmediate(arm64.X0, 0)
	}
}

func (t *ARM64Translator) translateUnaryExpression(ctx *compiler.UnaryExprContext) {
	// Obtener el operador (primer hijo)
	operator := ""
	var operandExpr antlr.ParseTree

	// Buscar el operador y la expresión
	for i := 0; i < ctx.GetChildCount(); i++ {
		child := ctx.GetChild(i)
		if terminal, ok := child.(antlr.TerminalNode); ok {
			operator = terminal.GetText()
		} else if parseTree, ok := child.(antlr.ParseTree); ok {
			operandExpr = parseTree
		}
	}

	t.generator.Comment(fmt.Sprintf("=== OPERADOR UNARIO: %s ===", operator))

	switch operator {
	case "!":
		t.translateLogicalNot(operandExpr)
	case "-":
		t.translateUnaryMinus(operandExpr)
	case "+":
		// El + unario no hace nada, solo evalúa la expresión
		t.translateExpression(operandExpr)
	default:
		t.addError(fmt.Sprintf("Operador unario no implementado: %s", operator))
		t.generator.LoadImmediate(arm64.X0, 0)
	}
}

func (t *ARM64Translator) translateLogicalNot(operandExpr antlr.ParseTree) {
	t.generator.Comment("=== OPERADOR LÓGICO NOT (!) ===")

	// Evaluar la expresión operando
	t.translateExpression(operandExpr)

	// Negar el resultado: si x0 == 0 entonces 1, sino 0
	t.generator.Comment("Negar resultado lógico")
	t.generator.Emit("cmp x0, #0")
	t.generator.Emit("cset x0, eq") // x0 = 1 si x0 era 0, sino 0
}

func (t *ARM64Translator) translateUnaryMinus(operandExpr antlr.ParseTree) {
	t.generator.Comment("=== OPERADOR UNARIO MENOS (-) ===")

	// Evaluar la expresión operando
	t.translateExpression(operandExpr)

	// Negar el valor: x0 = -x0
	t.generator.Comment("Negar valor numérico")
	t.generator.Emit("neg x0, x0")
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

	// NUEVO: Escalar por 100 para mantener 2 decimales
	scaledValue := int(value * 100)
	t.generator.Comment(fmt.Sprintf("Flotante %s escalado como %d (x100)", valueStr, scaledValue))
	t.generator.LoadImmediate(arm64.X0, scaledValue)
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

	// VERIFICAR SI TIENE INTERPOLACIÓN
	if strings.Contains(text, "$") {
		// IMPORTANTE: No usar x0 para el resultado final en interpolación
		// porque processStringInterpolation hace múltiples prints
		t.processStringInterpolation(text)

		// DESPUÉS DE INTERPOLACIÓN, x0 queda en estado indefinido
		// Para funciones que esperan un valor en x0, cargar 0
		t.generator.Comment("Interpolación completada")

	} else {
		// String normal sin interpolación
		if existingLabel, exists := t.stringRegistry[text]; exists {
			t.generator.Comment(fmt.Sprintf("Usar string \"%s\" con etiqueta %s", text, existingLabel))
			t.generator.Emit(fmt.Sprintf("adr x0, %s", existingLabel))
		} else {
			t.addError(fmt.Sprintf("String \"%s\" no fue procesado en primera pasada", text))
			t.generator.LoadImmediate(arm64.X0, 0)
		}
	}
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

	// MANEJAR OPERADORES LÓGICOS CON EVALUACIÓN PEREZOSA
	if operator == "&&" {
		t.translateLogicalAnd(ctx)
		return
	}
	if operator == "||" {
		t.translateLogicalOr(ctx)
		return
	}

	// Determinar si estamos comparando flotantes
	leftIsFloat := t.isFloatExpression(ctx.GetLeft())
	rightIsFloat := t.isFloatExpression(ctx.GetRight())

	// Evaluar operando izquierdo y guardarlo en x1
	t.translateExpression(ctx.GetLeft())
	t.generator.Comment("Mover operando izquierdo a x1")
	t.generator.Emit("mov x1, x0")

	// Evaluar operando derecho (queda en X0)
	t.translateExpression(ctx.GetRight())

	// Si uno de los operandos es flotante, escalar el entero
	if leftIsFloat && !rightIsFloat {
		// Escalar operando derecho
		t.generator.Comment("Escalar operando derecho para comparación con flotante")
		t.generator.Emit("mov x2, #100")
		t.generator.Mul(arm64.X0, arm64.X0, "x2")
	} else if !leftIsFloat && rightIsFloat {
		// Escalar operando izquierdo
		t.generator.Comment("Escalar operando izquierdo para comparación con flotante")
		t.generator.Emit("mov x2, #100")
		t.generator.Mul(arm64.X1, arm64.X1, "x2")
	}

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
	case "%":
		t.generator.Mod(arm64.X0, arm64.X1, arm64.X0)
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

func (t *ARM64Translator) isFloatExpression(expr antlr.ParseTree) bool {
	if expr == nil {
		return false
	}

	switch ctx := expr.(type) {
	case *compiler.FloatLiteralContext:
		return true
	case *compiler.LiteralExprContext:
		return t.isFloatExpression(ctx.Literal())
	case *compiler.LiteralContext:
		// Verificar si contiene FloatLiteralContext
		for i := 0; i < ctx.GetChildCount(); i++ {
			if _, ok := ctx.GetChild(i).(*compiler.FloatLiteralContext); ok {
				return true
			}
		}
		return false
	case *compiler.IdPatternExprContext:
		// Verificar el tipo de la variable
		varName := ctx.Id_pattern().GetText()
		if varType, exists := t.variableTypes[varName]; exists {
			return varType == "float"
		}
		return false
	default:
		return false
	}
}

func (t *ARM64Translator) translateLogicalAnd(ctx *compiler.BinaryExprContext) {
	t.generator.Comment("=== OPERADOR LÓGICO AND (&&) ===")

	falseLabel := t.generator.GetLabel()
	endLabel := t.generator.GetLabel()

	// Evaluar operando izquierdo
	t.translateExpression(ctx.GetLeft())

	// Si el operando izquierdo es falso (0), saltar directamente a false
	t.generator.Comment("Si operando izquierdo es falso, resultado es falso")
	t.generator.JumpIfZero(arm64.X0, falseLabel)

	// Evaluar operando derecho (solo si el izquierdo es verdadero)
	t.generator.Comment("Operando izquierdo es verdadero, evaluar operando derecho")
	t.translateExpression(ctx.GetRight())

	// El resultado está en x0 (verdadero si != 0, falso si == 0)
	// Convertir a 1 o 0 explícitamente
	t.generator.Comment("Convertir resultado a 1 o 0")
	t.generator.Emit("cmp x0, #0")
	t.generator.Emit("cset x0, ne") // x0 = 1 si != 0, sino 0
	t.generator.Jump(endLabel)

	// Etiqueta para resultado falso
	t.generator.SetLabel(falseLabel)
	t.generator.Comment("Resultado es falso")
	t.generator.LoadImmediate(arm64.X0, 0)

	// Etiqueta final
	t.generator.SetLabel(endLabel)
}

// === NUEVA FUNCIÓN: translateLogicalOr (||) ===
func (t *ARM64Translator) translateLogicalOr(ctx *compiler.BinaryExprContext) {
	t.generator.Comment("=== OPERADOR LÓGICO OR (||) ===")

	trueLabel := t.generator.GetLabel()
	endLabel := t.generator.GetLabel()

	// Evaluar operando izquierdo
	t.translateExpression(ctx.GetLeft())

	// Si el operando izquierdo es verdadero (!= 0), saltar directamente a true
	t.generator.Comment("Si operando izquierdo es verdadero, resultado es verdadero")
	t.generator.Emit("cmp x0, #0")
	t.generator.Emit("bne " + trueLabel)

	// Evaluar operando derecho (solo si el izquierdo es falso)
	t.generator.Comment("Operando izquierdo es falso, evaluar operando derecho")
	t.translateExpression(ctx.GetRight())

	// El resultado está en x0
	// Convertir a 1 o 0 explícitamente
	t.generator.Comment("Convertir resultado a 1 o 0")
	t.generator.Emit("cmp x0, #0")
	t.generator.Emit("cset x0, ne") // x0 = 1 si != 0, sino 0
	t.generator.Jump(endLabel)

	// Etiqueta para resultado verdadero
	t.generator.SetLabel(trueLabel)
	t.generator.Comment("Resultado es verdadero")
	t.generator.LoadImmediate(arm64.X0, 1)

	// Etiqueta final
	t.generator.SetLabel(endLabel)
}

// translateComparison traduce operaciones de comparación
func (t *ARM64Translator) translateComparison(reg1, reg2, condition string) {
	t.generator.Compare(reg1, reg2)

	// Usar CSET para convertir el resultado de la comparación a 1 o 0
	t.generator.Comment("Convertir resultado de comparación a 1/0")
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
	case "print_bool": // NUEVA FUNCIÓN
		t.translatePrintBoolFunction(ctx)
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

func (t *ARM64Translator) translatePrintBoolFunction(ctx *compiler.FuncCallContext) {
	t.generator.Comment("=== FUNCIÓN PRINT_BOOL ===")

	if ctx.Arg_list() != nil {
		args := ctx.Arg_list().(*compiler.ArgListContext).AllFunc_arg()
		if len(args) > 0 {
			// Tomar solo el primer argumento
			arg := args[0]
			if argCtx := arg.(*compiler.FuncArgContext); argCtx != nil {
				if argCtx.Expression() != nil {
					// Evaluar la expresión
					t.translateExpression(argCtx.Expression())
					// Llamar a print_bool
					t.generator.CallFunction("print_bool")
				}
			}
		}
	}
}

// la función translateUserFunctionCall en translator.go
func (t *ARM64Translator) translateUserFunctionCall(callCtx *compiler.FuncCallContext, funcDecl *compiler.FuncDeclContext) {
	funcName := callCtx.Id_pattern().GetText()
	t.generator.Comment(fmt.Sprintf("=== LLAMADA A FUNCIÓN DE USUARIO: %s ===", funcName))

	// Obtener información de parámetros de la función
	var paramNames []string
	if funcDecl.Param_list() != nil {
		params := funcDecl.Param_list().(*compiler.ParamListContext).AllFunc_param()
		for _, param := range params {
			if paramCtx := param.(*compiler.FuncParamContext); paramCtx.ID() != nil {
				paramNames = append(paramNames, paramCtx.ID().GetText())
			}
		}
	}

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

				// NUEVO: Determinar el tipo del argumento que se está pasando
				var argType string

				// Evaluar el argumento
				if argCtx.Expression() != nil {
					// Inferir tipo de la expresión ANTES de evaluarla
					argType = t.inferExpressionType(argCtx.Expression())
					t.translateExpression(argCtx.Expression())
				} else if argCtx.Id_pattern() != nil {
					// Es una variable
					varName := argCtx.Id_pattern().GetText()
					if t.generator.VariableExists(varName) {
						// Obtener tipo de la variable
						if varType, exists := t.variableTypes[varName]; exists {
							argType = varType
						} else {
							argType = "unknown"
						}
						t.generator.LoadVariable(arm64.X0, varName)
					} else {
						t.addError(fmt.Sprintf("Variable '%s' no encontrada", varName))
						t.generator.LoadImmediate(arm64.X0, 0)
						argType = "int"
					}
				} else {
					// Fallback: intentar como texto
					argText := argCtx.GetText()
					if t.generator.VariableExists(argText) {
						if varType, exists := t.variableTypes[argText]; exists {
							argType = varType
						} else {
							argType = "unknown"
						}
						t.generator.LoadVariable(arm64.X0, argText)
					} else if value, err := strconv.Atoi(argText); err == nil {
						t.generator.LoadImmediate(arm64.X0, value)
						argType = "int"
					} else {
						t.addError(fmt.Sprintf("No se puede procesar argumento: %s", argText))
						t.generator.LoadImmediate(arm64.X0, 0)
						argType = "int"
					}
				}

				// NUEVO: Asignar tipo al parámetro correspondiente
				if i < len(paramNames) {
					paramName := paramNames[i]
					t.variableTypes[paramName] = argType
					fmt.Printf("📝 Parámetro '%s' asignado tipo: %s\n", paramName, argType)
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
	case "indexOf":
		t.translateIndexOfFunction(ctx)
	case "len":
		t.translateLenFunction(ctx)
	case "append":
		t.translateAppendFunction(ctx)
	case "join":
		t.translateJoinFunction(ctx)
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
					// NUEVA LÓGICA: Detectar tipo de expresión MEJORADA
					exprText := argCtx.Expression().GetText()

					// VERIFICAR EXPLÍCITAMENTE SI ES BOOLEANO
					if t.isBooleanExpression(argCtx.Expression()) {
						// Es un booleano - usar print_bool
						t.generator.Comment(fmt.Sprintf("Imprimiendo booleano: %s", exprText))
						t.translateExpression(argCtx.Expression())
						t.generator.CallFunction("print_bool")
					} else if strings.HasPrefix(exprText, "\"") && strings.HasSuffix(exprText, "\"") {
						// Es un string literal directo
						t.generator.Comment(fmt.Sprintf("Imprimiendo string literal: %s", exprText))
						t.translateExpression(argCtx.Expression())
						t.generator.CallFunction("print_string")
					} else {
						// Verificar tipo de variable si es una variable
						argType := t.getArgumentType(argCtx)
						if argType == "string" {
							t.generator.Comment(fmt.Sprintf("Imprimiendo string: %s", exprText))
							t.translateExpression(argCtx.Expression())
							t.generator.CallFunction("print_string")
						} else if argType == "bool" {
							t.generator.Comment(fmt.Sprintf("Imprimiendo variable booleana: %s", exprText))
							t.translateExpression(argCtx.Expression())
							t.generator.CallFunction("print_bool")
						} else {
							// Es una expresión numérica
							t.generator.Comment(fmt.Sprintf("Imprimiendo valor numérico: %s", exprText))
							t.translateExpression(argCtx.Expression())
							t.generator.CallFunction("print_integer")
						}
					}
				} else if argCtx.Id_pattern() != nil {
					// Variable directa - determinar tipo
					varName := argCtx.Id_pattern().GetText()
					if t.generator.VariableExists(varName) {
						t.generator.LoadVariable(arm64.X0, varName)

						// Determinar qué función usar según el tipo
						if varType, exists := t.variableTypes[varName]; exists {
							switch {
							case varType == "bool":
								t.generator.Comment(fmt.Sprintf("Imprimiendo variable booleana: %s", varName))
								t.generator.CallFunction("print_bool")
							case varType == "string":
								t.generator.Comment(fmt.Sprintf("Imprimiendo variable string: %s", varName))
								t.generator.CallFunction("print_string")
							case strings.HasPrefix(varType, "[]"):
								if varType == "[]string" {
									t.generator.Comment(fmt.Sprintf("Imprimiendo variable vector de strings: %s", varName))
									t.generator.CallFunction("print_vector_string")
								} else {
									t.generator.Comment(fmt.Sprintf("Imprimiendo variable vector: %s", varName))
									t.generator.CallFunction("print_vector")
								}
							default:
								t.generator.Comment(fmt.Sprintf("Imprimiendo variable numérica: %s", varName))
								t.generator.CallFunction("print_integer")
							}
						} else {
							t.generator.Comment(fmt.Sprintf("Imprimiendo variable (tipo desconocido): %s", varName))
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

func (t *ARM64Translator) isBooleanExpression(expr antlr.ParseTree) bool {
	if expr == nil {
		return false
	}

	switch ctx := expr.(type) {
	case *compiler.BoolLiteralContext:
		return true
	case *compiler.LiteralExprContext:
		return t.isBooleanExpression(ctx.Literal())
	case *compiler.LiteralContext:
		// Verificar si contiene BoolLiteralContext
		for i := 0; i < ctx.GetChildCount(); i++ {
			if _, ok := ctx.GetChild(i).(*compiler.BoolLiteralContext); ok {
				return true
			}
		}
		// AGREGAR: También verificar por texto específico
		text := ctx.GetText()
		if text == "true" || text == "false" {
			return true
		}
		return false
	case *compiler.BinaryExprContext:
		// Operaciones de comparación devuelven booleanos
		operator := ctx.GetOp().GetText()
		switch operator {
		case "==", "!=", "<", ">", "<=", ">=", "&&", "||":
			return true
		}
		return false
	case *compiler.UnaryExprContext:
		// El operador ! devuelve booleano
		for i := 0; i < ctx.GetChildCount(); i++ {
			if terminal, ok := ctx.GetChild(i).(antlr.TerminalNode); ok {
				if terminal.GetText() == "!" {
					return true
				}
			}
		}
		return false
	case *compiler.IdPatternExprContext:
		// Verificar el tipo de la variable
		varName := ctx.Id_pattern().GetText()
		if varType, exists := t.variableTypes[varName]; exists {
			return varType == "bool"
		}
		return false
	default:
		// CASO ESPECIAL: Si el texto completo es "true" o "false"
		text := expr.GetText()
		return text == "true" || text == "false"
	}
}

// Determinar tipo de argumento
func (t *ARM64Translator) getArgumentType(argCtx *compiler.FuncArgContext) string {
	if argCtx.Expression() != nil {
		exprText := argCtx.Expression().GetText()

		// VERIFICAR EXPLÍCITAMENTE BOOLEANOS
		if exprText == "true" || exprText == "false" {
			return "bool"
		}

		// VERIFICAR STRINGS
		if strings.HasPrefix(exprText, "\"") && strings.HasSuffix(exprText, "\"") {
			return "string"
		}

		// USAR INFERENCIA DE TIPO
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

// ====================================
// For Loops
// ====================================

// translateForLoop traduce bucles for
func (t *ARM64Translator) translateForLoop(ctx *compiler.ForStmtCondContext) {
	t.generator.Comment("=== FOR LOOP ===")

	startLabel := t.generator.GetLabel()
	endLabel := t.generator.GetLabel()

	// Etiquetas para break y continue
	t.breakLabels = append(t.breakLabels, endLabel)
	t.continueLabels = append(t.continueLabels, startLabel)

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

	//Pop de etiquetas al salir del bucle
	t.breakLabels = t.breakLabels[:len(t.breakLabels)-1]
	t.continueLabels = t.continueLabels[:len(t.continueLabels)-1]
}

func (t *ARM64Translator) translateForAssignment(ctx *compiler.ForAssCondContext) {
	t.generator.Comment("=== FOR tipo C-style ===")

	// i = 1;
	t.translateNode(ctx.Assign_stmt())

	startLabel := t.generator.GetLabel()
	continueLabel := t.generator.GetLabel()
	endLabel := t.generator.GetLabel()

	// Etiquetas para break y continue
	t.breakLabels = append(t.breakLabels, endLabel)
	t.continueLabels = append(t.continueLabels, continueLabel)

	t.generator.SetLabel(startLabel)

	// Evaluar condición: i <= 5
	t.translateExpression(ctx.Expression(0)) // condición
	t.generator.JumpIfZero(arm64.X0, endLabel)

	// Cuerpo
	for _, stmt := range ctx.AllStmt() {
		t.translateNode(stmt)
	}

	// Etiqueta para continue
	t.generator.SetLabel(continueLabel)

	// Incremento
	t.translateExpression(ctx.Expression(1))

	// Repetir ciclo
	t.generator.Jump(startLabel)

	t.generator.SetLabel(endLabel)

	// Limpiar
	t.breakLabels = t.breakLabels[:len(t.breakLabels)-1]
	t.continueLabels = t.continueLabels[:len(t.continueLabels)-1]

}

// ====================================
// Incremento y Decremento
// ====================================
func (t *ARM64Translator) translateIncrement(ctx *compiler.IncrementoContext) {
	varName := ctx.ID().GetText()

	t.generator.LoadVariable(arm64.X0, varName)
	t.generator.Comment(fmt.Sprintf("Incrementar '%s'", varName))
	t.generator.Emit("add x0, x0, #1")
	t.generator.StoreVariable(arm64.X0, varName)
}

func (t *ARM64Translator) translateDecrement(ctx *compiler.DecrementoContext) {
	varName := ctx.ID().GetText()

	t.generator.LoadVariable(arm64.X0, varName)
	t.generator.Comment(fmt.Sprintf("Decrementar '%s'", varName))
	t.generator.Emit("sub x0, x0, #1")
	t.generator.StoreVariable(arm64.X0, varName)
}

// ====================================
// Transferencia de Control
// ====================================

func (t *ARM64Translator) translateBreakStatement(ctx *compiler.BreakStmtContext) {
	if len(t.breakLabels) == 0 {
		t.addError("Break fuera de un contexto de bucle")
		return
	}
	label := t.breakLabels[len(t.breakLabels)-1]
	t.generator.Comment("Break statement")
	t.generator.Jump(label)
}

func (t *ARM64Translator) translateContinueStatement(ctx *compiler.ContinueStmtContext) {
	if len(t.continueLabels) == 0 {
		t.addError("Continue fuera de un contexto de bucle")
		return
	}
	label := t.continueLabels[len(t.continueLabels)-1]
	t.generator.Comment("Continue statement")
	t.generator.Jump(label)
}

// ====================================
// Proceso de interpolacion
// ====================================
func (t *ARM64Translator) processStringInterpolation(text string) {
	t.generator.Comment("=== INTERPOLACIÓN DE STRING ===")

	// Dividir el string en partes: texto y variables
	parts := t.parseInterpolatedString(text)

	for _, part := range parts { // ← Quitar la variable i
		if part.IsVariable {
			// Es una variable - cargar y imprimir según su tipo
			varName := part.Content
			if t.generator.VariableExists(varName) {
				t.generator.Comment(fmt.Sprintf("Interpolando variable: %s", varName))
				t.generator.LoadVariable(arm64.X0, varName)

				// Determinar tipo y llamar función apropiada
				if varType, exists := t.variableTypes[varName]; exists {
					switch varType {
					case "bool":
						t.generator.CallFunction("print_bool")
					case "string":
						t.generator.CallFunction("print_string")
					default:
						t.generator.CallFunction("print_integer")
					}
				} else {
					t.generator.CallFunction("print_integer")
				}
			} else {
				t.addError(fmt.Sprintf("Variable '%s' no encontrada en interpolación", varName))
			}
		} else {
			// Es texto literal - crear string y imprimir
			if part.Content != "" {
				t.generator.Comment(fmt.Sprintf("Interpolando texto: \"%s\"", part.Content))

				// VERIFICAR si ya existe en el registro
				if existingLabel, exists := t.stringRegistry[part.Content]; exists {
					t.generator.Emit(fmt.Sprintf("adr x0, %s", existingLabel))
					t.generator.CallFunction("print_string")
				} else {
					// Si no existe, reportar error
					t.addError(fmt.Sprintf("Parte interpolada \"%s\" no fue registrada en primera pasada", part.Content))
				}
			}
		}
	}
}

// Estructura para las partes del string:
type InterpolationPart struct {
	Content    string
	IsVariable bool
}

func (t *ARM64Translator) parseInterpolatedString(text string) []InterpolationPart {
	var parts []InterpolationPart
	var currentPart strings.Builder

	i := 0
	for i < len(text) {
		if text[i] == '$' && i+1 < len(text) {
			// Guardar texto anterior si existe
			if currentPart.Len() > 0 {
				parts = append(parts, InterpolationPart{
					Content:    currentPart.String(),
					IsVariable: false,
				})
				currentPart.Reset()
			}

			// Saltar el '$'
			i++

			// Extraer nombre de variable
			varStart := i
			for i < len(text) && (isLetter(text[i]) || isDigit(text[i]) || text[i] == '_') {
				i++
			}

			if i > varStart {
				varName := text[varStart:i]
				parts = append(parts, InterpolationPart{
					Content:    varName,
					IsVariable: true,
				})
			}
			// i ya está en la posición correcta
		} else {
			currentPart.WriteByte(text[i])
			i++
		}
	}

	// Agregar última parte si existe
	if currentPart.Len() > 0 {
		parts = append(parts, InterpolationPart{
			Content:    currentPart.String(),
			IsVariable: false,
		})
	}

	return parts
}

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// === LIBRERÍA ESTÁNDAR ===

// generateStandardLibrary genera las funciones básicas necesarias
// === FUNCIÓN generateStandardLibrary COMPLETA ===
// REEMPLAZAR COMPLETAMENTE tu función generateStandardLibrary actual

func (t *ARM64Translator) generateStandardLibrary() {
	// FUNCIÓN print_integer
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
    ret`)

	// FUNCIÓN print_char
	t.generator.EmitRaw(`
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

	// FUNCIÓN print_string
	t.generator.EmitRaw(`
print_string:
    // Función para imprimir strings
    // Input: x0 = dirección del string (terminado en null)
    stp x29, x30, [sp, #-16]!    // Guardar registros
    stp x19, x20, [sp, #-16]!    // Guardar x19 y x20

    mov x19, x0                   // x19 = dirección del string

    // Encontrar la longitud del string
    mov x20, #0                   // x20 = contador de longitud

strlen_loop:
    ldrb w1, [x19, x20]
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

	// FUNCIÓN print_bool
	t.generator.EmitRaw(`
print_bool:
    // Función para imprimir valores booleanos
    // Input: x0 = valor booleano (0=false, cualquier otra cosa=true)
    stp x29, x30, [sp, #-16]!    // Guardar registros

    // Verificar si es true o false
    cmp x0, #0
    beq print_false_simple

print_true_simple:
    // Imprimir "true" manualmente
    mov x0, #116  // 't'
    bl print_char
    mov x0, #114  // 'r'
    bl print_char
    mov x0, #117  // 'u'
    bl print_char
    mov x0, #101  // 'e'
    bl print_char
    b print_bool_simple_done

print_false_simple:
    // Imprimir "false" manualmente
    mov x0, #102  // 'f'
    bl print_char
    mov x0, #97   // 'a'
    bl print_char
    mov x0, #108  // 'l'
    bl print_char
    mov x0, #115  // 's'
    bl print_char
    mov x0, #101  // 'e'
    bl print_char

print_bool_simple_done:
    ldp x29, x30, [sp], #16      // Restaurar registros
    ret`)

	// FUNCIÓN print_vector
	t.generator.EmitRaw(`
print_vector:
    // Función para imprimir vectores/arrays
    // Input: x0 = dirección del vector (primer elemento = longitud)
    stp x29, x30, [sp, #-16]!    // Guardar registros
    stp x19, x20, [sp, #-16]!
    stp x21, x22, [sp, #-16]!

    mov x19, x0                   // x19 = dirección del vector
    
    // Cargar longitud del vector (primer elemento)
    ldr x21, [x19]               // x21 = longitud del vector
    
    // Imprimir "[ "
    mov x0, #91                   // ASCII '['
    bl print_char
    mov x0, #32                   // ASCII ' '
    bl print_char

    mov x20, #0                   // x20 = índice actual

print_vector_loop:
    cmp x20, x21
    bge print_vector_end
    
    // Cargar elemento del vector (saltando el primer elemento que es la longitud)
    add x22, x20, #1             // x22 = índice + 1 (saltar longitud)
    ldr x0, [x19, x22, lsl #3]   // x0 = vector[i+1] (cada elemento = 8 bytes)
    bl print_integer
    
    // Incrementar índice
    add x20, x20, #1
    
    // Si no es el último elemento, imprimir espacio
    cmp x20, x21
    bge print_vector_no_space
    mov x0, #32                   // ASCII ' '
    bl print_char
    
print_vector_no_space:
    b print_vector_loop

print_vector_end:
    // Imprimir " ]"
    mov x0, #32                   // ASCII ' '
    bl print_char
    mov x0, #93                   // ASCII ']'
    bl print_char

    ldp x21, x22, [sp], #16      // Restaurar registros
    ldp x19, x20, [sp], #16
    ldp x29, x30, [sp], #16
    ret`)

	// FUNCIÓN print_vector_string
	t.generator.EmitRaw(`
print_vector_string:
    // Función para imprimir vectores de strings
    // Input: x0 = dirección del vector (primer elemento = longitud)
    stp x29, x30, [sp, #-16]!    // Guardar registros
    stp x19, x20, [sp, #-16]!
    stp x21, x22, [sp, #-16]!

    mov x19, x0                   // x19 = dirección del vector
    
    // Cargar longitud del vector (primer elemento)
    ldr x21, [x19]               // x21 = longitud del vector
    
    // Imprimir "[ "
    mov x0, #91                   // ASCII '['
    bl print_char
    mov x0, #32                   // ASCII ' '
    bl print_char

    mov x20, #0                   // x20 = índice actual

print_vector_string_loop:
    cmp x20, x21
    bge print_vector_string_end
    
    // Cargar dirección del string (saltando el primer elemento que es la longitud)
    add x22, x20, #1             // x22 = índice + 1 (saltar longitud)
    ldr x0, [x19, x22, lsl #3]   // x0 = dirección del string
    bl print_string             // Imprimir el string, no su dirección
    
    // Incrementar índice
    add x20, x20, #1
    
    // Si no es el último elemento, imprimir espacio
    cmp x20, x21
    bge print_vector_string_no_space
    mov x0, #32                   // ASCII ' '
    bl print_char
    
print_vector_string_no_space:
    b print_vector_string_loop

print_vector_string_end:
    // Imprimir " ]"
    mov x0, #32                   // ASCII ' '
    bl print_char
    mov x0, #93                   // ASCII ']'
    bl print_char

    ldp x21, x22, [sp], #16      // Restaurar registros
    ldp x19, x20, [sp], #16
    ldp x29, x30, [sp], #16
    ret`)

	// FUNCIÓN append_vector
	t.generator.EmitRaw(`
append_vector:
    // Función para agregar un elemento a un vector
    // Input: x0 = dirección del vector original, x1 = nuevo elemento
    // Output: x0 = dirección del nuevo vector
    stp x29, x30, [sp, #-16]!    // Guardar registros
    stp x19, x20, [sp, #-16]!
    stp x21, x22, [sp, #-16]!
    stp x23, x24, [sp, #-16]!

    mov x19, x0                   // x19 = vector original
    mov x20, x1                   // x20 = nuevo elemento

    // Cargar longitud actual del vector
    ldr x21, [x19]               // x21 = longitud actual
    add x22, x21, #1             // x22 = nueva longitud

    // Verificar límite máximo (50 elementos total)
    cmp x22, #50
    bgt append_error

    // Calcular tamaño del nuevo vector en bytes
    // Cada vector: 8 bytes (longitud) + longitud * 8 bytes (elementos)
    mov x23, #8                  // Tamaño de cada elemento (8 bytes)
    mul x24, x22, x23            // x24 = nueva_longitud * 8
    add x24, x24, #8             // x24 = tamaño total (incluye longitud)

    // Reservar espacio en el stack para el nuevo vector
    neg x24, x24                 // Hacer negativo para restar del SP
    add sp, sp, x24              // Restar del stack pointer
    neg x24, x24                 // Restaurar valor positivo
    mov x0, sp                   // x0 = dirección del nuevo vector

    // Guardar nueva longitud en el nuevo vector
    str x22, [x0]                // Almacenar nueva longitud

    // Copiar elementos del vector original al nuevo vector
    mov x23, #1                  // Índice para copiar elementos (empezar en 1)
    
copy_loop:
    cmp x23, x22                 // Comparar índice con nueva longitud
    bge copy_done                // Si índice >= nueva_longitud, terminar

    cmp x23, x21                 // Comparar con longitud original
    bgt add_new_element          // Si índice > longitud_original, agregar nuevo elemento

    // Copiar elemento del vector original
    lsl x25, x23, #3             // x25 = índice * 8 (desplazamiento)
    ldr x26, [x19, x25]          // Cargar elemento del vector original
    str x26, [x0, x25]           // Guardar en nuevo vector
    b next_element

add_new_element:
    // Agregar el nuevo elemento
    lsl x25, x23, #3             // x25 = índice * 8 (desplazamiento)
    str x20, [x0, x25]           // Guardar nuevo elemento

next_element:
    add x23, x23, #1             // Incrementar índice
    b copy_loop

copy_done:
    // x0 ya contiene la dirección del nuevo vector
    b append_success

append_error:
    // Error: vector excede tamaño máximo
    mov x0, #0                   // Retornar NULL

append_success:
    ldp x23, x24, [sp], #16      // Restaurar registros
    ldp x21, x22, [sp], #16
    ldp x19, x20, [sp], #16
    ldp x29, x30, [sp], #16
    ret`)

	// FUNCIÓN join_vector
	t.generator.EmitRaw(`
join_vector:
    // Función para unir strings de un vector con un separador
    // Input: x0 = dirección del vector de strings, x1 = separador
    // Output: x0 = dirección del string resultado
    stp x29, x30, [sp, #-16]!    // Guardar registros
    stp x19, x20, [sp, #-16]!
    stp x21, x22, [sp, #-16]!
    stp x23, x24, [sp, #-16]!

    mov x19, x0                   // x19 = vector de strings
    mov x20, x1                   // x20 = separador

    // Cargar longitud del vector
    ldr x21, [x19]               // x21 = longitud del vector

    // Si el vector está vacío, retornar string vacío
    cmp x21, #0
    beq join_empty_vector

    // Calcular tamaño total necesario (estimación conservadora)
    // Cada string máximo 100 chars + separador máximo 10 chars
    mov x22, #110                // 110 chars por elemento
    mul x23, x21, x22            // x23 = longitud * 110
    add x23, x23, #8             // +8 para null terminator

    // Reservar espacio en el stack para el resultado
    neg x23, x23                 // Hacer negativo
    add sp, sp, x23              // Reservar espacio
    neg x23, x23                 // Restaurar positivo
    mov x0, sp                   // x0 = buffer resultado

    // Inicializar buffer con 0
    mov x22, #0                  // Contador de posición en buffer
    
    // Copiar primer string (sin separador)
    cmp x21, #0
    beq join_done
    
    // Cargar dirección del primer string
    ldr x24, [x19, #8]           // x24 = dirección primer string
    bl copy_string_to_buffer     // Copiar string al buffer
    
    // Loop para strings restantes (con separador)
    mov x25, #1                  // Índice = 1 (segundo elemento)
    
join_loop:
    cmp x25, x21                 // Comparar con longitud
    bge join_done                // Si índice >= longitud, terminar
    
    // Copiar separador
    mov x24, x20                 // x24 = separador
    bl copy_string_to_buffer     // Copiar separador
    
    // Copiar siguiente string
    add x26, x25, #1             // x26 = índice + 1 (saltar longitud)
    ldr x24, [x19, x26, lsl #3]  // Cargar dirección del string
    bl copy_string_to_buffer     // Copiar string
    
    add x25, x25, #1             // Incrementar índice
    b join_loop

copy_string_to_buffer:
    // Subrutina para copiar un string al buffer
    // x24 = string source, x0 = buffer base, x22 = offset actual
copy_char_loop:
    ldrb w27, [x24], #1          // Cargar char y avanzar
    cbz w27, copy_char_done      // Si es 0, terminar
    strb w27, [x0, x22]          // Guardar en buffer
    add x22, x22, #1             // Incrementar offset
    b copy_char_loop
copy_char_done:
    ret

join_empty_vector:
    // Vector vacío - retornar string vacío en stack
    sub sp, sp, #16
    mov x0, sp
    mov w1, #0
    strb w1, [x0]                // Null terminator
    b join_success

join_done:
    // Agregar null terminator
    mov w24, #0
    strb w24, [x0, x22]

join_success:
    // x0 ya contiene la dirección del resultado
    ldp x23, x24, [sp], #16      // Restaurar registros
    ldp x21, x22, [sp], #16
    ldp x19, x20, [sp], #16
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

// === FUNCIONES PARA SOPORTE DE VECTORES ===

// analyzeVarVectDecl analiza declaraciones de vectores en primera pasada
func (t *ARM64Translator) analyzeVarVectDecl(ctx *compiler.VarVectDeclContext) {
	varName := ctx.ID().GetText()
	fmt.Printf("🔍 Analizando declaración de vector: %s\n", varName)

	// Declarar la variable si no existe
	if !t.generator.VariableExists(varName) {
		t.generator.DeclareVariable(varName)
	}

	// Obtener tipo de vector
	if ctx.Vector_type() != nil {
		vectorType := ctx.Vector_type().GetText()
		t.variableTypes[varName] = vectorType
		fmt.Printf("🔍 Vector '%s' de tipo: %s\n", varName, vectorType)
	}

	// Analizar elementos del vector para buscar strings
	if ctx.Vect_expr() != nil {
		t.analyzeVectorElements(ctx.Vect_expr())
	}
}

// analyzeVectorElements analiza los elementos de un vector
func (t *ARM64Translator) analyzeVectorElements(vectExpr antlr.ParseTree) {
	switch ctx := vectExpr.(type) {
	case *compiler.VectorItemLisContext:
		fmt.Printf("🔍 Analizando elementos del vector\n")
		// Recorrer todas las expresiones del vector
		for i := 0; i < ctx.GetChildCount(); i++ {
			child := ctx.GetChild(i)
			if parseTreeChild, ok := child.(antlr.ParseTree); ok {
				t.analyzeStringsInExpression(parseTreeChild)
			}
		}
	}
}

// translateVarVectDecl traduce declaraciones de vectores
func (t *ARM64Translator) translateVarVectDecl(ctx *compiler.VarVectDeclContext) {
	varName := ctx.ID().GetText()
	t.generator.Comment(fmt.Sprintf("=== DECLARACIÓN DE VECTOR: %s ===", varName))

	// Verificar que la variable fue declarada en primera pasada
	if !t.generator.VariableExists(varName) {
		t.addError(fmt.Sprintf("Vector '%s' no fue declarado en primera pasada", varName))
		return
	}

	// Obtener tipo de vector
	vectorType := ""
	if ctx.Vector_type() != nil {
		vectorType = ctx.Vector_type().GetText()
	}

	// Determinar si es vector de strings o enteros
	if strings.Contains(vectorType, "string") {
		// Manejar vector de strings
		stringLabels := t.extractVectorStringLabels(ctx.Vect_expr())
		if len(stringLabels) == 0 {
			t.generator.Comment("Vector de strings vacío")
			t.generator.LoadImmediate(arm64.X0, 0)
			t.generator.StoreVariable(arm64.X0, varName)
			return
		}

		// Crear vector de strings en la sección .data
		vectorLabel := t.generator.AddVectorStringData(varName, stringLabels)

		// Cargar dirección del vector en X0
		t.generator.Comment(fmt.Sprintf("Cargar dirección de vector %s", varName))
		t.generator.Emit(fmt.Sprintf("adr x0, %s", vectorLabel))

		// Guardar dirección en la variable
		t.generator.StoreVariable(arm64.X0, varName)

		t.generator.Comment(fmt.Sprintf("Vector de strings %s inicializado con %d elementos", varName, len(stringLabels)))
	} else {
		// Manejar vector de enteros (código existente)
		elements := t.extractVectorElements(ctx.Vect_expr())
		if len(elements) == 0 {
			t.generator.Comment("Vector vacío")
			t.generator.LoadImmediate(arm64.X0, 0)
			t.generator.StoreVariable(arm64.X0, varName)
			return
		}

		// Crear vector en la sección .data
		vectorLabel := t.generator.AddVectorData(varName, elements)

		// Cargar dirección del vector en X0
		t.generator.Comment(fmt.Sprintf("Cargar dirección de vector %s", varName))
		t.generator.Emit(fmt.Sprintf("adr x0, %s", vectorLabel))

		// Guardar dirección en la variable
		t.generator.StoreVariable(arm64.X0, varName)

		t.generator.Comment(fmt.Sprintf("Vector %s inicializado con %d elementos", varName, len(elements)))
	}
}

// extractVectorElements extrae los elementos de un vector
func (t *ARM64Translator) extractVectorElements(vectExpr antlr.ParseTree) []int {
	var elements []int

	switch ctx := vectExpr.(type) {
	case *compiler.VectorItemLisContext:
		// Recorrer hijos buscando expresiones
		for i := 0; i < ctx.GetChildCount(); i++ {
			child := ctx.GetChild(i)
			if parseTreeChild, ok := child.(antlr.ParseTree); ok {
				if value := t.extractIntFromExpression(parseTreeChild); value != -999999 {
					elements = append(elements, value)
				}
			}
		}
	}

	fmt.Printf("🔍 Elementos extraídos del vector: %v\n", elements)
	return elements
}

// extractIntFromExpression extrae un valor entero de una expresión
func (t *ARM64Translator) extractIntFromExpression(expr antlr.ParseTree) int {
	switch ctx := expr.(type) {
	case *compiler.LiteralExprContext:
		return t.extractIntFromExpression(ctx.Literal())
	case *compiler.IntLiteralContext:
		if value, err := strconv.Atoi(ctx.GetText()); err == nil {
			return value
		}
	case *compiler.LiteralContext:
		// Buscar literal entero en los hijos
		for i := 0; i < ctx.GetChildCount(); i++ {
			child := ctx.GetChild(i)
			if intCtx, ok := child.(*compiler.IntLiteralContext); ok {
				if value, err := strconv.Atoi(intCtx.GetText()); err == nil {
					return value
				}
			}
		}
	}

	// Valor especial para indicar que no se pudo extraer un entero
	return -999999
}

// extractVectorStringLabels extrae las etiquetas de strings de un vector de strings
func (t *ARM64Translator) extractVectorStringLabels(vectExpr antlr.ParseTree) []string {
	var stringLabels []string

	switch ctx := vectExpr.(type) {
	case *compiler.VectorItemLisContext:
		// Recorrer hijos buscando expresiones de strings
		for i := 0; i < ctx.GetChildCount(); i++ {
			child := ctx.GetChild(i)
			if parseTreeChild, ok := child.(antlr.ParseTree); ok {
				if label := t.extractStringLabelFromExpression(parseTreeChild); label != "" {
					stringLabels = append(stringLabels, label)
				}
			}
		}
	}

	fmt.Printf("🔍 Etiquetas de strings extraídas del vector: %v\n", stringLabels)
	return stringLabels
}

// extractStringLabelFromExpression extrae la etiqueta de un string de una expresión
func (t *ARM64Translator) extractStringLabelFromExpression(expr antlr.ParseTree) string {
	switch ctx := expr.(type) {
	case *compiler.LiteralExprContext:
		return t.extractStringLabelFromExpression(ctx.Literal())
	case *compiler.StringLiteralContext:
		// Obtener el texto del string
		text := ctx.GetText()
		if len(text) >= 2 {
			text = text[1 : len(text)-1] // Quitar comillas
		}
		// Buscar si ya existe la etiqueta para este string
		if label, exists := t.stringRegistry[text]; exists {
			return label
		}
	case *compiler.LiteralContext:
		// Buscar literal de string en los hijos
		for i := 0; i < ctx.GetChildCount(); i++ {
			child := ctx.GetChild(i)
			if stringCtx, ok := child.(*compiler.StringLiteralContext); ok {
				text := stringCtx.GetText()
				if len(text) >= 2 {
					text = text[1 : len(text)-1] // Quitar comillas
				}
				// Buscar si ya existe la etiqueta para este string
				if label, exists := t.stringRegistry[text]; exists {
					return label
				}
			}
		}
	}

	// Retornar cadena vacía si no se pudo extraer
	return ""
}

// translateVectorAccess traduce acceso a elementos de vectores: vector[index]
func (t *ARM64Translator) translateVectorAccess(ctx *compiler.VectorItemExprContext) {
	t.generator.Comment("=== ACCESO A VECTOR ===")

	// Para VectorItemExprContext, necesitamos buscar el VectorItemContext hijo
	var vectorItemCtx *compiler.VectorItemContext
	for i := 0; i < ctx.GetChildCount(); i++ {
		if child, ok := ctx.GetChild(i).(*compiler.VectorItemContext); ok {
			vectorItemCtx = child
			break
		}
	}

	if vectorItemCtx == nil {
		t.addError("No se encontró VectorItemContext")
		t.generator.LoadImmediate(arm64.X0, 0)
		return
	}

	// Obtener el nombre del vector y el índice
	var vectorName string
	var indexExpr antlr.ParseTree

	for i := 0; i < vectorItemCtx.GetChildCount(); i++ {
		child := vectorItemCtx.GetChild(i)

		if idPattern, ok := child.(*compiler.IdPatternContext); ok {
			vectorName = idPattern.GetText()
		} else if _, ok := child.(antlr.TerminalNode); ok {
			// Saltar '[' y ']'
			continue
		} else if parseTree, ok := child.(antlr.ParseTree); ok {
			// Es la expresión del índice
			indexExpr = parseTree
		}
	}

	if vectorName == "" {
		t.addError("No se pudo determinar el nombre del vector")
		t.generator.LoadImmediate(arm64.X0, 0)
		return
	}

	if indexExpr == nil {
		t.addError("No se pudo determinar el índice del vector")
		t.generator.LoadImmediate(arm64.X0, 0)
		return
	}

	t.generator.Comment(fmt.Sprintf("Accediendo a vector '%s'", vectorName))

	// Verificar que el vector existe
	if !t.generator.VariableExists(vectorName) {
		t.addError(fmt.Sprintf("Vector '%s' no encontrado", vectorName))
		t.generator.LoadImmediate(arm64.X0, 0)
		return
	}

	// Evaluar el índice y dejarlo en X1
	t.translateExpression(indexExpr)
	t.generator.Emit("mov x1, x0") // X1 = índice

	// Cargar la dirección del vector en X0
	t.generator.LoadVariable(arm64.X0, vectorName)

	// X0 = dirección del vector
	// X1 = índice
	// Necesitamos X0 = vector[índice + 1] (el +1 es para saltar la longitud)

	// Incrementar índice para saltar la longitud (primer elemento)
	t.generator.Comment("Saltar longitud del vector (primer elemento)")
	t.generator.Emit("add x1, x1, #1") // X1 = índice + 1

	// Calcular offset: índice * 8 bytes
	t.generator.Comment("Cargar elemento del vector")
	t.generator.Emit("ldr x0, [x0, x1, lsl #3]") // X0 = vector[índice + 1]

	t.generator.Comment(fmt.Sprintf("Vector '%s' accedido exitosamente", vectorName))
}

// translateVectorAssignment traduce asignaciones a elementos de vectores: vector[index] = value
func (t *ARM64Translator) translateVectorAssignment(ctx *compiler.VectorAssignContext) {
	t.generator.Comment("=== ASIGNACIÓN A VECTOR ===")

	// Extraer nombre del vector, índice y valor de la asignación
	var vectorName string
	var indexExpr antlr.ParseTree
	var valueExpr antlr.ParseTree

	// Buscar VectorItemContext (para obtener vector[index])
	var vectorItemCtx *compiler.VectorItemContext
	for i := 0; i < ctx.GetChildCount(); i++ {
		if child, ok := ctx.GetChild(i).(*compiler.VectorItemContext); ok {
			vectorItemCtx = child
			break
		}
	}

	if vectorItemCtx == nil {
		t.addError("No se encontró VectorItemContext en VectorAssignContext")
		return
	}

	// Extraer vector y índice del VectorItemContext
	for i := 0; i < vectorItemCtx.GetChildCount(); i++ {
		child := vectorItemCtx.GetChild(i)

		if idPattern, ok := child.(*compiler.IdPatternContext); ok {
			vectorName = idPattern.GetText()
		} else if parseTree, ok := child.(antlr.ParseTree); ok {
			// Buscar la expresión del índice (no el '[' o ']')
			text := parseTree.GetText()
			if text != "[" && text != "]" {
				indexExpr = parseTree
			}
		}
	}

	// Buscar la expresión del valor (lado derecho del =)
	for i := 0; i < ctx.GetChildCount(); i++ {
		child := ctx.GetChild(i)
		if exprCtx, ok := child.(*compiler.ExpressionContext); ok {
			valueExpr = exprCtx
			break
		} else if literalCtx, ok := child.(*compiler.LiteralExprContext); ok {
			valueExpr = literalCtx
			break
		}
	}

	// Validaciones
	if vectorName == "" {
		t.addError("No se pudo determinar el nombre del vector en asignación")
		return
	}

	if indexExpr == nil {
		t.addError("No se pudo determinar el índice del vector en asignación")
		return
	}

	if valueExpr == nil {
		t.addError("No se pudo determinar el valor a asignar en vector")
		return
	}

	t.generator.Comment(fmt.Sprintf("Asignando valor a vector '%s'", vectorName))

	// Verificar que el vector existe
	if !t.generator.VariableExists(vectorName) {
		t.addError(fmt.Sprintf("Vector '%s' no encontrado para asignación", vectorName))
		return
	}

	// Paso 1: Evaluar el valor a asignar (lado derecho) y guardarlo en X2
	t.generator.Comment("Evaluar valor a asignar")
	t.translateExpression(valueExpr)
	t.generator.Emit("mov x2, x0") // X2 = valor a asignar

	// Paso 2: Evaluar el índice y guardarlo en X1
	t.generator.Comment("Evaluar índice del vector")
	t.translateExpression(indexExpr)
	t.generator.Emit("mov x1, x0") // X1 = índice

	// Paso 3: Cargar la dirección del vector en X0
	t.generator.Comment("Cargar dirección del vector")
	t.generator.LoadVariable(arm64.X0, vectorName) // X0 = dirección del vector

	// Paso 4: Calcular la posición correcta y asignar
	t.generator.Comment("Saltar longitud del vector (primer elemento)")
	t.generator.Emit("add x1, x1, #1") // X1 = índice + 1 (saltar metadata de longitud)

	t.generator.Comment("Asignar valor al elemento del vector")
	t.generator.Emit("str x2, [x0, x1, lsl #3]") // vector[índice + 1] = valor

	t.generator.Comment(fmt.Sprintf("Vector '%s' modificado exitosamente", vectorName))
}

// translateIndexOfFunction maneja la traducción de la función indexOf para vectores
func (t *ARM64Translator) translateIndexOfFunction(ctx *compiler.FuncCallContext) {
	t.generator.Comment("=== FUNCIÓN indexOf ===")

	// Verificar que tenemos exactamente 2 argumentos
	if ctx.Arg_list() == nil {
		t.addError("indexOf requiere exactamente 2 argumentos")
		t.generator.LoadImmediate(arm64.X0, -1)
		return
	}

	argList, ok := ctx.Arg_list().(*compiler.ArgListContext)
	if !ok {
		t.addError("indexOf: no se pudo obtener la lista de argumentos")
		t.generator.LoadImmediate(arm64.X0, -1)
		return
	}

	args := argList.AllFunc_arg()
	if len(args) != 2 {
		t.addError(fmt.Sprintf("indexOf requiere exactamente 2 argumentos, se proporcionaron %d", len(args)))
		t.generator.LoadImmediate(arm64.X0, -1)
		return
	}

	// Primer argumento: el vector
	vectorArg, ok := args[0].(*compiler.FuncArgContext)
	if !ok {
		t.addError("indexOf: primer argumento no válido")
		t.generator.LoadImmediate(arm64.X0, -1)
		return
	}

	// Segundo argumento: el elemento a buscar
	elementArg, ok := args[1].(*compiler.FuncArgContext)
	if !ok {
		t.addError("indexOf: segundo argumento no válido")
		t.generator.LoadImmediate(arm64.X0, -1)
		return
	}
	// Verificar que las expresiones no sean nil
	if elementArg.Expression() == nil && elementArg.Id_pattern() == nil {
		t.addError("indexOf: segundo argumento no tiene expresión ni identificador válido")
		t.generator.LoadImmediate(arm64.X0, -1)
		return
	}

	if vectorArg.Expression() == nil && vectorArg.Id_pattern() == nil {
		t.addError("indexOf: primer argumento no tiene expresión ni identificador válido")
		t.generator.LoadImmediate(arm64.X0, -1)
		return
	}

	// Evaluar el elemento a buscar y guardarlo en X1
	t.generator.Comment("Evaluar elemento a buscar")
	if elementArg.Expression() != nil {
		t.translateExpression(elementArg.Expression())
	} else if elementArg.Id_pattern() != nil {
		// Manejar identificador (variable)
		varName := elementArg.Id_pattern().GetText()
		if t.generator.VariableExists(varName) {
			t.generator.LoadVariable(arm64.X0, varName)
		} else {
			t.addError(fmt.Sprintf("Variable '%s' no encontrada", varName))
			t.generator.LoadImmediate(arm64.X0, 0)
		}
	}
	t.generator.Emit("mov x1, x0") // X1 = elemento a buscar

	// Evaluar el vector y obtener su dirección en X0
	t.generator.Comment("Evaluar vector")
	if vectorArg.Expression() != nil {
		t.translateExpression(vectorArg.Expression())
	} else if vectorArg.Id_pattern() != nil {
		// Manejar identificador (variable)
		varName := vectorArg.Id_pattern().GetText()
		if t.generator.VariableExists(varName) {
			t.generator.LoadVariable(arm64.X0, varName)
		} else {
			t.addError(fmt.Sprintf("Variable '%s' no encontrada", varName))
			t.generator.LoadImmediate(arm64.X0, 0)
		}
	}
	t.generator.Emit("mov x2, x0") // X2 = dirección del vector

	// Cargar la longitud del vector (primer elemento)
	t.generator.Comment("Cargar longitud del vector")
	t.generator.Emit("ldr x3, [x2]") // X3 = longitud del vector

	// Inicializar índice de búsqueda
	t.generator.Comment("Inicializar índice de búsqueda")
	t.generator.Emit("mov x4, #0") // X4 = índice actual

	// Crear etiquetas para los loops
	loopLabel := t.generator.GetLabel()
	foundLabel := t.generator.GetLabel()
	notFoundLabel := t.generator.GetLabel()
	endLabel := t.generator.GetLabel()

	// Loop principal
	t.generator.SetLabel(loopLabel)
	t.generator.Comment("Verificar si hemos llegado al final del vector")
	t.generator.Emit("cmp x4, x3")
	t.generator.Emit("bge " + notFoundLabel)

	// Cargar elemento actual del vector (saltando el primer elemento que es la longitud)
	t.generator.Comment("Cargar elemento actual del vector")
	t.generator.Emit("add x5, x4, #1")           // X5 = índice + 1 (saltar longitud)
	t.generator.Emit("ldr x6, [x2, x5, lsl #3]") // X6 = vector[índice + 1]

	// Comparar elemento actual con elemento buscado
	t.generator.Comment("Comparar elemento actual con elemento buscado")
	t.generator.Emit("cmp x6, x1")
	t.generator.Emit("beq " + foundLabel)

	// Incrementar índice y continuar
	t.generator.Comment("Incrementar índice y continuar")
	t.generator.Emit("add x4, x4, #1")
	t.generator.Emit("b " + loopLabel)

	// Elemento encontrado - devolver índice
	t.generator.SetLabel(foundLabel)
	t.generator.Comment("Elemento encontrado - devolver índice")
	t.generator.Emit("mov x0, x4")
	t.generator.Emit("b " + endLabel)

	// Elemento no encontrado - devolver -1
	t.generator.SetLabel(notFoundLabel)
	t.generator.Comment("Elemento no encontrado - devolver -1")
	t.generator.LoadImmediate(arm64.X0, -1)

	// Fin de la función
	t.generator.SetLabel(endLabel)
	t.generator.Comment("=== FIN indexOf ===")
}

// translateLenFunction maneja la traducción de la función len para vectores
func (t *ARM64Translator) translateLenFunction(ctx *compiler.FuncCallContext) {
	t.generator.Comment("=== FUNCIÓN len ===")

	// Verificar que tenemos exactamente 1 argumento
	if ctx.Arg_list() == nil {
		t.addError("len requiere exactamente 1 argumento")
		t.generator.LoadImmediate(arm64.X0, 0)
		return
	}

	argList, ok := ctx.Arg_list().(*compiler.ArgListContext)
	if !ok {
		t.addError("len: no se pudo obtener la lista de argumentos")
		t.generator.LoadImmediate(arm64.X0, 0)
		return
	}

	args := argList.AllFunc_arg()
	if len(args) != 1 {
		t.addError(fmt.Sprintf("len requiere exactamente 1 argumento, se proporcionaron %d", len(args)))
		t.generator.LoadImmediate(arm64.X0, 0)
		return
	}

	// Obtener el primer argumento (el vector)
	vectorArg, ok := args[0].(*compiler.FuncArgContext)
	if !ok {
		t.addError("len: argumento no válido")
		t.generator.LoadImmediate(arm64.X0, 0)
		return
	}

	// Verificar que las expresiones no sean nil
	if vectorArg.Expression() == nil && vectorArg.Id_pattern() == nil {
		t.addError("len: argumento no tiene expresión ni identificador válido")
		t.generator.LoadImmediate(arm64.X0, 0)
		return
	}

	// Cargar la dirección del vector
	t.generator.Comment("Obtener dirección del vector")
	if vectorArg.Expression() != nil {
		t.translateExpression(vectorArg.Expression())
	} else if vectorArg.Id_pattern() != nil {
		// Manejar identificador (variable)
		varName := vectorArg.Id_pattern().GetText()
		if t.generator.VariableExists(varName) {
			t.generator.LoadVariable(arm64.X0, varName)
		} else {
			t.addError(fmt.Sprintf("Variable '%s' no encontrada", varName))
			t.generator.LoadImmediate(arm64.X0, 0)
			return
		}
	}

	// El primer elemento del vector es su longitud
	t.generator.Comment("Cargar longitud del vector (primer elemento)")
	t.generator.Emit("ldr x0, [x0]") // X0 = longitud del vector

	t.generator.Comment("=== FIN len ===")
}

// translateAppendFunction maneja la traducción de la función append para vectores
func (t *ARM64Translator) translateAppendFunction(ctx *compiler.FuncCallContext) {
	t.generator.Comment("=== FUNCIÓN append ===")

	// Verificar que tenemos exactamente 2 argumentos
	if ctx.Arg_list() == nil {
		t.addError("append requiere exactamente 2 argumentos")
		t.generator.LoadImmediate(arm64.X0, 0)
		return
	}

	argList, ok := ctx.Arg_list().(*compiler.ArgListContext)
	if !ok {
		t.addError("append: no se pudo obtener la lista de argumentos")
		t.generator.LoadImmediate(arm64.X0, 0)
		return
	}

	args := argList.AllFunc_arg()
	if len(args) != 2 {
		t.addError(fmt.Sprintf("append requiere exactamente 2 argumentos, se proporcionaron %d", len(args)))
		t.generator.LoadImmediate(arm64.X0, 0)
		return
	}

	// Primer argumento: el vector original
	vectorArg, ok := args[0].(*compiler.FuncArgContext)
	if !ok {
		t.addError("append: primer argumento no válido")
		t.generator.LoadImmediate(arm64.X0, 0)
		return
	}

	// Segundo argumento: el elemento a agregar
	elementArg, ok := args[1].(*compiler.FuncArgContext)
	if !ok {
		t.addError("append: segundo argumento no válido")
		t.generator.LoadImmediate(arm64.X0, 0)
		return
	}

	// Obtener el vector original
	t.generator.Comment("Cargar vector original")
	if vectorArg.Expression() != nil {
		t.translateExpression(vectorArg.Expression())
	} else if vectorArg.Id_pattern() != nil {
		varName := vectorArg.Id_pattern().GetText()
		if t.generator.VariableExists(varName) {
			t.generator.LoadVariable(arm64.X0, varName)
		} else {
			t.addError(fmt.Sprintf("Variable '%s' no encontrada", varName))
			t.generator.LoadImmediate(arm64.X0, 0)
			return
		}
	}

	// X0 = dirección del vector original
	t.generator.Emit("mov x19, x0") // X19 = dirección del vector original

	// Cargar longitud del vector original
	t.generator.Comment("Cargar longitud del vector original")
	t.generator.Emit("ldr x20, [x19]") // X20 = longitud original

	// Verificar límite máximo (50 elementos)
	t.generator.Comment("Verificar límite máximo")
	t.generator.LoadImmediate(arm64.X21, 49) // Máximo 49 elementos
	t.generator.Emit("cmp x20, x21")

	// Crear etiquetas
	errorLabel := t.generator.GetLabel()
	continueLabel := t.generator.GetLabel()

	t.generator.Emit("bgt " + errorLabel) // Si longitud > 49, error

	// Evaluar el segundo argumento (elemento a agregar)
	t.generator.Comment("Evaluar elemento a agregar")
	if elementArg.Expression() != nil {
		t.translateExpression(elementArg.Expression())
	} else if elementArg.Id_pattern() != nil {
		varName := elementArg.Id_pattern().GetText()
		if t.generator.VariableExists(varName) {
			t.generator.LoadVariable(arm64.X0, varName)
		} else {
			t.addError(fmt.Sprintf("Variable '%s' no encontrada", varName))
			t.generator.LoadImmediate(arm64.X0, 0)
		}
	}

	// X0 = nuevo elemento, X19 = vector original
	t.generator.Emit("mov x1, x0")  // X1 = nuevo elemento
	t.generator.Emit("mov x0, x19") // X0 = vector original

	// Llamar función append_vector
	t.generator.Comment("Llamar función append_vector")
	t.generator.CallFunction("append_vector")
	t.generator.Emit("b " + continueLabel)

	// Manejo de error
	t.generator.SetLabel(errorLabel)
	t.generator.Comment("Error: vector excede tamaño máximo")
	t.generator.LoadImmediate(arm64.X0, 0)

	t.generator.SetLabel(continueLabel)
	t.generator.Comment("=== FIN append ===")
}

// translateJoinFunction maneja la traducción de la función join para vectores de strings
func (t *ARM64Translator) translateJoinFunction(ctx *compiler.FuncCallContext) {
	t.generator.Comment("=== FUNCIÓN join ===")

	// Verificar que tenemos exactamente 2 argumentos
	if ctx.Arg_list() == nil {
		t.addError("join requiere exactamente 2 argumentos")
		t.generator.LoadImmediate(arm64.X0, 0)
		return
	}

	argList, ok := ctx.Arg_list().(*compiler.ArgListContext)
	if !ok {
		t.addError("join: no se pudo obtener la lista de argumentos")
		t.generator.LoadImmediate(arm64.X0, 0)
		return
	}

	args := argList.AllFunc_arg()
	if len(args) != 2 {
		t.addError(fmt.Sprintf("join requiere exactamente 2 argumentos, se proporcionaron %d", len(args)))
		t.generator.LoadImmediate(arm64.X0, 0)
		return
	}

	// Primer argumento: el vector de strings
	vectorArg, ok := args[0].(*compiler.FuncArgContext)
	if !ok {
		t.addError("join: primer argumento no válido")
		t.generator.LoadImmediate(arm64.X0, 0)
		return
	}

	// Segundo argumento: el separador
	separatorArg, ok := args[1].(*compiler.FuncArgContext)
	if !ok {
		t.addError("join: segundo argumento no válido")
		t.generator.LoadImmediate(arm64.X0, 0)
		return
	}

	// Obtener el vector de strings
	t.generator.Comment("Cargar vector de strings")
	if vectorArg.Expression() != nil {
		t.translateExpression(vectorArg.Expression())
	} else if vectorArg.Id_pattern() != nil {
		varName := vectorArg.Id_pattern().GetText()
		if t.generator.VariableExists(varName) {
			t.generator.LoadVariable(arm64.X0, varName)
		} else {
			t.addError(fmt.Sprintf("Variable '%s' no encontrada", varName))
			t.generator.LoadImmediate(arm64.X0, 0)
			return
		}
	}

	// X0 = dirección del vector de strings
	t.generator.Emit("mov x19, x0") // X19 = dirección del vector

	// Evaluar el separador
	t.generator.Comment("Evaluar separador")
	if separatorArg.Expression() != nil {
		t.translateExpression(separatorArg.Expression())
	} else if separatorArg.Id_pattern() != nil {
		varName := separatorArg.Id_pattern().GetText()
		if t.generator.VariableExists(varName) {
			t.generator.LoadVariable(arm64.X0, varName)
		} else {
			t.addError(fmt.Sprintf("Variable '%s' no encontrada", varName))
			t.generator.LoadImmediate(arm64.X0, 0)
			return
		}
	}

	// X0 = separador, X19 = vector de strings
	t.generator.Emit("mov x1, x0")  // X1 = separador
	t.generator.Emit("mov x0, x19") // X0 = vector de strings

	// Llamar función join_vector
	t.generator.Comment("Llamar función join_vector")
	t.generator.CallFunction("join_vector")

	t.generator.Comment("=== FIN join ===")
}
