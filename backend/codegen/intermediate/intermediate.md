# 🎯 ¿Para qué sirve ir.go?

El archivo ir.go es el corazón del sistema de Representación Intermedia (IR). Básicamente, es como crear un "lenguaje intermedio" entre el código 

## 🤔 ¿Por qué necesitamos esto?

Imagina que quieres traducir del español al japonés, pero no sabes japonés directamente. Una estrategia sería:

Español → Inglés (que sí conoces)

Inglés → Japonés (usando un traductor)

En compiladores es igual:

VLan Cherry → IR (representación intermedia)
IR → ARM64 (código de máquina)

### 🏗️ ¿Qué define ir.go?

1. IROpcode - Los "Verbos" del IR

```go
// Go
IR_ADD    // Suma dos valores
IR_SUB    // Resta dos valores
IR_STORE  // Guarda un valor en memoria
IR_LOAD   // Carga un valor de memoria
IR_BRANCH // Salta a otro lugar del código
```

Es como definir las operaciones básicas que cualquier programa puede hacer, sin importar el lenguaje original.

2. IROperand - Los "Sustantivos" del IR

```go
// Tipos de operandos:
IR_OPERAND_TEMP      // %t1, %t2 (temporales)
IR_OPERAND_VAR       // $x, $y (variables)
IR_OPERAND_IMMEDIATE // #5, #10 (valores constantes)
IR_OPERAND_LABEL     // @loop_start (etiquetas para saltos)
```

Son las "cosas" con las que trabajamos: variables, números, direcciones de memoria.

3. IRInstruction - Las "Oraciones" del IR

```go
go// Una instrucción típica:
ADD %t1, $x, #5  // %t1 = $x + 5
```

Cada instrucción dice: "haz esta operación con estos operandos".
🎯 Ejemplo Práctico
Si tienes este código VLan Cherry:

```go
vlancherrymut x int = 10
mut y int = 20
mut sum int = x + y
print(sum)
```

ir.go permite convertirlo a algo como:

```go
irALLOC_LOCAL $x          // Reservar espacio para x
LOAD_IMM %t1, #10       // Cargar 10 en temporal t1
STORE $x, %t1           // Guardar t1 en variable x

ALLOC_LOCAL $y          // Reservar espacio para y  
LOAD_IMM %t2, #20       // Cargar 20 en temporal t2
STORE $y, %t2           // Guardar t2 en variable y

LOAD %t3, $x            // Cargar x en t3
LOAD %t4, $y            // Cargar y en t4
ADD %t5, %t3, %t4       // t5 = t3 + t4
PRINT %t5               // Imprimir t5
```

### 🔄 ¿Por qué es útil?

Simplifica la traducción: Es más fácil convertir IR→ARM64 que VLanCherry→ARM64 directamente
Permite optimizaciones: Puedes optimizar el IR antes de generar código final
Reutilizable: El mismo IR puede generar código para ARM64, x86, RISC-V, etc.
Debugging: Puedes ver exactamente qué está haciendo tu programa paso a paso

### 🎨 Analogía Simple
Piensa en ir.go como definir las piezas de LEGO básicas:

Las formas de las piezas (IROpcode)
Los colores disponibles (IROperandType)
Las instrucciones de construcción (IRInstruction)

Con estas piezas básicas, puedes construir cualquier programa, sin importar qué lenguaje uses originalmente. ¡Y luego convertir esa construcción al código de máquina que necesites!

Es la base fundamental que hace posible todo el resto del compilador. Sin ir.go, no tendrías un "lenguaje común" para representar programas internamente.

## 🔧 ¿Qué es generator.go?

generator.go es el traductor que convierte tu código VLan Cherry (en forma de AST) a la representación intermedia (IR) que definimos en ir.go.
🎯 Su Función Principal
Es como un intérprete de idiomas que toma las estructuras del AST (árbol de sintaxis) y las convierte a instrucciones IR simples y uniformes.

### 🔄 ¿Cómo Funciona?

Entrada: AST (Árbol de Sintaxis)

```go
BinaryExpr
├── Left: Variable("x")
├── Op: "+"
└── Right: Literal(5)
```

Salida: Instrucciones IR

```go
LOAD %t1, $x        // Cargar variable x
LOAD_IMM %t2, #5    // Cargar constante 5
ADD %t3, %t1, %t2   // Sumar t1 + t2 y guardar en t3
```

### 🏗️ Componentes Principales

1. IRGenerator - El "Cerebro"

```go
type IRGenerator struct {
    program      *IRProgram     // El programa IR que está construyendo
    tempCounter  int            // Contador para generar temporales únicos (%t1, %t2...)
    labelCounter int            // Contador para generar etiquetas únicas (loop_1, if_2...)
    symbolTable  map[string]*IROperand // Mapeo de variables a operandos IR
}
```

2. Visitadores Específicos - Los "Traductores"

Cada tipo de nodo del AST tiene su propio método de traducción:

```go
// Para variables: mut x int = 10
func (gen *IRGenerator) visitMutVarDecl(ctx *MutVarDeclContext) {
    // 1. Crear operando para la variable
    // 2. Generar ALLOC_LOCAL para reservar memoria
    // 3. Evaluar la expresión de inicialización
    // 4. Generar STORE para guardar el valor inicial
}

// Para expresiones: x + y
func (gen *IRGenerator) visitBinaryExpr(ctx *BinaryExprContext) {
    // 1. Traducir lado izquierdo (x)
    // 2. Traducir lado derecho (y)  
    // 3. Generar instrucción de operación (ADD)
    // 4. Retornar temporal con el resultado
}
```

### 🎨 Ejemplo Paso a Paso

Veamos cómo generator.go traduce este código:
vlancherrymut x int = 10

```go
go
if x > 5 {
    print("grande")
}
```

Paso 1: Declaración de variable

```go
govisitMutVarDecl() {
    varName = "x"
    varOperand = newVariable("x", "int")          // Crear $x
    emit(IR_ALLOC_LOCAL, varOperand, ...)         // ALLOC_LOCAL $x
    initValue = visit(expression)                 // Evaluar "10"
    emit(IR_STORE, varOperand, initValue, ...)    // STORE $x, #10
}
```

Genera:

```go
go
irALLOC_LOCAL $x
LOAD_IMM %t1, #10
STORE $x, %t1
```

Paso 2: Condición del if

```go
govisitIfStmt() {
    elseLabel = newLabel("if_else")               // Crear etiqueta
    endLabel = newLabel("if_end")
    
    condition = visit(expression)                 // Evaluar "x > 5"
    emit(IR_BRANCH_IF_FALSE, condition, elseLabel) // Saltar si falso
    
    // Cuerpo del if...
    emit(IR_BRANCH, endLabel)                     // Saltar al final
    emitLabel(elseLabel)
    emitLabel(endLabel)
}
```

Genera:

```go
irLOAD %t2, $x
LOAD_IMM %t3, #5
CMP_GT %t4, %t2, %t3
BRANCH_IF_FALSE %t4, @if_else_1
PRINT_LN "grande"
BRANCH @if_end_1
if_else_1:
if_end_1:
```

🧠 Características Inteligentes

1. Generación de Temporales Únicos

```go
gofunc (gen *IRGenerator) newTemp(dataType string) *IROperand {
    gen.tempCounter++
    return &IROperand{
        Name: fmt.Sprintf("t%d", gen.tempCounter),  // %t1, %t2, %t3...
        Type: IR_OPERAND_TEMP,
        DataType: dataType,
    }
}
```

1. Manejo de Corto Circuito (para && y ||)

```go
go// Para: a && b
visitLogicalExpr() {
    left = visit(a)
    falseLabel = newLabel("and_false")
    
    emit(IR_BRANCH_IF_FALSE, left, falseLabel)  // Si 'a' es falso, saltar
    right = visit(b)                            // Solo evaluar 'b' si 'a' es verdadero
    // ...
}
```

1. Seguimiento de Variables

```go
gosymbolTable["x"] = &IROperand{
    Type: IR_OPERAND_VAR,
    Name: "x",
    DataType: "int",
}
```

### 🎯 ¿Por Qué es Importante?

Simplificación: Convierte estructuras complejas del AST en instrucciones simples
Optimización: Puede aplicar optimizaciones durante la traducción
Abstracción: Separa la lógica del lenguaje de la generación de código final
Debugging: Permite ver exactamente cómo se traduce cada parte del código

🔗 Relación con Otros Archivos

Recibe: AST del parser de ANTLR
Usa: Definiciones de ir.go
Produce: Programa IR para optimizer.go
Alimenta: El optimizador y eventualmente el generador ARM64

📝 Analogía Simple
Si ir.go define el "alfabeto" del IR, entonces generator.go es el "escritor" que:

Lee oraciones en VLan Cherry (AST)
Traduce palabra por palabra al "alfabeto" IR
Escribe un texto completo en IR que significa lo mismo

¡Es el puente crucial entre "lo que escribiste" y "lo que la máquina puede entender"!
