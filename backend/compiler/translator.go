package compiler

import (
	"fmt"
	"strconv"

	"github.com/antlr4-go/antlr/v4"
	"main.go/compiler/arm64"
	compiler "main.go/grammar"
)

// ARM64Translator es el traductor principal de VlangCherry a ARM64
type ARM64Translator struct {
	generator *arm64.ARM64Generator
	errors    []string // Para almacenar errores de traducción
}

// NewARM64Translator crea un nuevo traductor
func NewARM64Translator() *ARM64Translator {
	return &ARM64Translator{
		generator: arm64.NewARM64Generator(),
		errors:    make([]string, 0),
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

	// Agregar funciones de librería estándar
	t.generator.EmitRaw("")
	t.generator.EmitRaw("// === LIBRERÍA ESTÁNDAR ===")
	t.generateStandardLibrary()

	return t.generator.GetCode(), t.errors
}

// === ANÁLISIS PREVIO ===

// analyzeVariables hace una pasada previa para encontrar todas las variables
// Esto nos permite reservar espacio en el stack antes de generar código
func (t *ARM64Translator) analyzeVariables(node antlr.ParseTree) {
	switch ctx := node.(type) {
	case *compiler.ProgramContext:
		for _, stmt := range ctx.AllStmt() {
			fmt.Printf("Analizando declaración: %T\n", stmt.GetText())
			t.analyzeVariables(stmt)
		}

	case *compiler.StmtContext:
		if ctx.Decl_stmt() != nil {
			t.analyzeVariables(ctx.Decl_stmt())
		}
		if ctx.If_stmt() != nil {
			t.analyzeVariables(ctx.If_stmt())
		}
		if ctx.For_stmt() != nil {
			t.analyzeVariables(ctx.For_stmt())
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
	default:
		// Para nodos no implementados, simplemente continuar
		t.addError(fmt.Sprintf("Nodo no implementado: %T", ctx))
	}
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
	}
}

// === DECLARACIONES DE VARIABLES ===

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
	switch ctx := expr.(type) {
	case *compiler.IntLiteralContext:
		t.translateIntLiteral(ctx)
	case *compiler.IdPatternExprContext:
		t.translateVariable(ctx)
	case *compiler.BinaryExprContext:
		t.translateBinaryExpression(ctx)
	case *compiler.ParensExprContext:
		// Expresión entre paréntesis: simplemente traducir la expresión interna
		t.translateExpression(ctx.Expression())
	default:
		t.addError(fmt.Sprintf("Expresión no implementada: %T", ctx))
		// Cargar 0 como valor por defecto
		t.generator.LoadImmediate(arm64.X0, 0)
	}
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

	// Evaluar operando izquierdo y guardarlo en el stack
	t.translateExpression(ctx.GetLeft())
	t.generator.Push(arm64.X0)

	// Evaluar operando derecho (queda en X0)
	t.translateExpression(ctx.GetRight())

	// Recuperar operando izquierdo en X1
	t.generator.Pop(arm64.X1)

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
	default:
		t.addError(fmt.Sprintf("Función no implementada: %s", funcName))
	}
}

// translatePrintFunction traduce llamadas a print/println
func (t *ARM64Translator) translatePrintFunction(ctx *compiler.FuncCallContext, withNewline bool) {
	t.generator.Comment("=== FUNCIÓN PRINT ===")

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
			if argCtx := arg.(*compiler.FuncArgContext); argCtx.Expression() != nil {
				t.translateExpression(argCtx.Expression())

				// Llamar a la función de print
				t.generator.CallFunction("print_integer")
			}
		}
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
