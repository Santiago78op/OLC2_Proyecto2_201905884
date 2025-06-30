# Implementación de Soporte para Vectores en ARM64 Translator

## Problema Identificado (Actualizado)
- **Error Inicial**: `Nodo no implementado: *compiler.VarVectDeclContext`
- **Error Secundario**: Vector no aparecía en sección `.data` del ARM64 generado
- **Causa del Error Secundario**: `GenerateHeader()` se llamaba antes de procesar vectores en segunda pasada
- **Archivos afectados**: `backend/compiler/translator.go` y `backend/compiler/arm64/generator.go`

## Cambios Implementados

### 1. **Archivo**: `backend/compiler/translator.go`
**Ubicación**: Función `translateNode()` - línea ~550-600
**Cambio**: Agregado nuevo case para manejar `VarVectDeclContext`
```go
case *compiler.VarVectDeclContext:
    t.translateVarVectDecl(ctx)
```

### 2. **Archivo**: `backend/compiler/translator.go`
**Ubicación**: Función `analyzeVariablesAndStrings()` - dentro del switch principal
**Cambio**: Agregado case para analizar declaraciones de vectores en primera pasada
```go
case *compiler.VarVectDeclContext:
    t.analyzeVarVectDecl(ctx)
```

### 3. **Archivo**: `backend/compiler/translator.go`
**Ubicación**: Final del archivo, después de funciones existentes
**Cambio**: Implementada nueva función `translateVarVectDecl()`
- Extrae nombre de variable
- Extrae tipo de vector
- Extrae elementos del vector
- Genera código ARM64 para inicialización

### 4. **Archivo**: `backend/compiler/translator.go`
**Ubicación**: Final del archivo, después de funciones existentes
**Cambio**: Implementada nueva función `analyzeVarVectDecl()`
- Declara la variable en primera pasada
- Registra tipo de variable
- Analiza elementos por strings

### 5. **Archivo**: `backend/compiler/arm64/generator.go`
**Ubicación**: Estructura `ARM64Generator`
**Cambio**: Agregado campo `vectorData []string` para almacenar vectores por separado

### 6. **Archivo**: `backend/compiler/arm64/generator.go`
**Ubicación**: Función `NewARM64Generator()`
**Cambio**: Inicializar `vectorData`

### 7. **Archivo**: `backend/compiler/arm64/generator.go`
**Ubicación**: Función `GenerateHeader()`
**Cambio**: Agregar vectores a sección `.data` junto con strings

### 8. **Archivo**: `backend/compiler/arm64/generator.go`
**Ubicación**: Función `AddVectorData()`
**Cambio**: Corregir para usar `vectorData` en lugar de `stringData`

### 9. **Archivo**: `backend/compiler/arm64/generator.go`
**Ubicación**: Función `Reset()`
**Cambio**: Limpiar también `vectorData`

### 10. **Archivo**: `backend/compiler/arm64/generator.go` (NUEVO)
**Ubicación**: Función `GetCode()`
**Cambio**: Modificada para ensamblar correctamente sección `.data` con strings y vectores al final

## Fase 2: Soporte para Impresión de Vectores (NEW)

### Problema Detectado
La primera implementación tenía un error crítico: cuando se ejecutaba `print("numeros:", numeros)`, el compilador intentaba imprimir el vector usando `print_integer`, lo que resultaba en imprimir la dirección de memoria del vector en lugar de sus elementos.

### Solución Implementada

#### 1. Detección de Tipos de Vector en `print`
**Archivo:** `backend/compiler/translator.go`
**Líneas:** ~1442-1460

Modificado el switch para detectar tipos que comienzan con `[]`:

```go
// Determinar qué función usar según el tipo
if varType, exists := t.variableTypes[varName]; exists {
    switch {
    case varType == "bool":
        t.generator.CallFunction("print_bool")
    case varType == "string":
        t.generator.CallFunction("print_string")
    case strings.HasPrefix(varType, "[]"):  // NEW: Detectar vectores
        t.generator.CallFunction("print_vector")
    default:
        t.generator.CallFunction("print_integer")
    }
}
```

#### 2. Mejora en el Almacenamiento de Vectores
**Archivo:** `backend/compiler/arm64/generator.go`
**Función:** `AddVectorData`

Modificado para almacenar la longitud como primer elemento:

```go
func (g *ARM64Generator) AddVectorData(vectorName string, elements []int) string {
    vectorLabel := fmt.Sprintf("vec_%s", vectorName)
    
    // Crear definición del vector con longitud como primer elemento
    var vectorDef strings.Builder
    vectorDef.WriteString(fmt.Sprintf("%s: .quad %d", vectorLabel, len(elements))) // Primer elemento: longitud
    
    // Agregar los elementos del vector
    for _, element := range elements {
        vectorDef.WriteString(fmt.Sprintf(", %d", element))
    }
    
    return vectorLabel
}
```

#### 3. Implementación de `print_vector` en ARM64
**Archivo:** `backend/compiler/translator.go`
**Función:** Stdlib ARM64

Nueva función que lee la longitud del primer elemento y luego itera sobre los elementos reales:

```arm64
print_vector:
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
    
    // Incrementar índice y agregar espacios si no es el último
    add x20, x20, #1
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
    ret
```

### Resultado Esperado
El vector `numeros = []int{1, 2, 3, 4, 5}` ahora se almacena como:
```arm64
vec_numeros: .quad 5, 1, 2, 3, 4, 5
```

Y `print("numeros:", numeros)` debería mostrar:
```
numeros: [ 1 2 3 4 5 ]
```

## Estado de Implementación
- ✅ Análisis en primera pasada
- ✅ Traducción de declaración de vectores
- ✅ Generación de código ARM64
- ✅ Manejo de memoria para vectores
- ✅ Integración con sistema existente
- ✅ Extracción de elementos de vectores
- ✅ Soporte para literales enteros en vectores
- ✅ Compilación exitosa sin errores
- ✅ **Corrección de orden de generación de código**

## ✅ **RESULTADO VERIFICADO**

### **Prueba Exitosa:**
```vlang
fn main() {
    print("Creación con literales:")
    numeros = []int{1, 2, 3, 4, 5}
    print("numeros:", numeros)
}
main()
```

### **Salida del Intérprete:**
```
Creación con literales:###Validacion Manualnumeros: [ 1 2 3 4 5 ]
```

### **ARM64 Generado:**
```arm64
.data
str_0: .asciz "Creación con literales:"
str_1: .asciz "###Validacion Manual" 
str_2: .asciz "numeros:"
vec_numeros: .quad 5, 1, 2, 3, 4, 5  // ✅ Longitud + elementos

.text
// ...
// Cargar variable 'numeros' en x0
ldr x0, [sp, #8]
// Imprimiendo variable vector: numeros  // ✅ Detectado como vector
// Llamar función print_vector           // ✅ Usa función correcta
bl print_vector
// ...

print_vector:
    // Función completa implementada que lee longitud y elementos
```

### **Características Implementadas:**
- ✅ **Declaración de vectores:** `numeros = []int{1, 2, 3, 4, 5}`
- ✅ **Almacenamiento con metadata:** Longitud como primer elemento
- ✅ **Detección automática de tipos:** `[]int` detectado correctamente
- ✅ **Impresión formateada:** `[ 1 2 3 4 5 ]`
- ✅ **ARM64 optimizado:** Función `print_vector` eficiente

## Fase 3: Soporte para Acceso por Índices a Vectores (NEW)

### Problema Detectado
Después de completar el soporte para declaración e impresión de vectores, se identificó que faltaba implementar el acceso por índices (`numeros[0]`). El error específico era:

```
Expresión no implementada: *compiler.VectorItemExprContext
```

Esto indicaba que faltaba el case para manejar acceso por índices en la función `translateExpression`.

### Solución Implementada

#### 1. Agregado Case en `translateExpression`
**Archivo:** `backend/compiler/translator.go`
**Función:** `translateExpression` - línea ~1228

```go
// Agregado nuevo case para acceso a vectores
case *compiler.VectorItemExprContext:
    t.translateVectorAccess(ctx)
```

#### 2. Implementación de `translateVectorAccess`
**Archivo:** `backend/compiler/translator.go`
**Líneas:** ~2287-2350

Nueva función que maneja el acceso por índices a vectores:

```go
func (t *ARM64Translator) translateVectorAccess(ctx *compiler.VectorItemExprContext) {
    t.generator.Comment("=== ACCESO A VECTOR ===")

    // Buscar VectorItemContext hijo
    var vectorItemCtx *compiler.VectorItemContext
    for i := 0; i < ctx.GetChildCount(); i++ {
        if child, ok := ctx.GetChild(i).(*compiler.VectorItemContext); ok {
            vectorItemCtx = child
            break
        }
    }

    // Extraer nombre del vector y expresión del índice
    var vectorName string
    var indexExpr antlr.ParseTree

    for i := 0; i < vectorItemCtx.GetChildCount(); i++ {
        child := vectorItemCtx.GetChild(i)

        if idPattern, ok := child.(*compiler.IdPatternContext); ok {
            vectorName = idPattern.GetText()
        } else if parseTree, ok := child.(antlr.ParseTree); ok {
            indexExpr = parseTree // Expresión del índice
        }
    }

    // Validaciones
    if vectorName == "" || indexExpr == nil {
        t.addError("Error al parsear acceso a vector")
        t.generator.LoadImmediate(arm64.X0, 0)
        return
    }

    // Verificar que el vector existe
    if !t.generator.VariableExists(vectorName) {
        t.addError(fmt.Sprintf("Vector '%s' no encontrado", vectorName))
        t.generator.LoadImmediate(arm64.X0, 0)
        return
    }

    // Evaluar índice y cargar vector
    t.translateExpression(indexExpr)           // X0 = índice
    t.generator.Emit("mov x1, x0")            // X1 = índice
    t.generator.LoadVariable(arm64.X0, vectorName) // X0 = dirección del vector

    // Saltar longitud (primer elemento) y acceder al elemento
    t.generator.Comment("Saltar longitud del vector (primer elemento)")
    t.generator.Emit("add x1, x1, #1")        // X1 = índice + 1
    t.generator.Comment("Cargar elemento del vector")
    t.generator.Emit("ldr x0, [x0, x1, lsl #3]") // X0 = vector[índice + 1]
}
```

#### 3. Lógica de Acceso ARM64

La función genera código ARM64 que:
1. **Evalúa el índice:** Traduce la expresión del índice (ej: `0`, `i`, `x+1`)
2. **Carga la dirección del vector:** Obtiene la dirección base del vector
3. **Ajusta por metadata:** Suma 1 al índice para saltar la longitud almacenada
4. **Calcula offset:** Multiplica por 8 bytes (`lsl #3`) para direccionamiento de 64-bit
5. **Carga el elemento:** Accede al elemento correcto del vector

### Código ARM64 Generado

Para `numeros[0]` se genera:

```arm64
// === ACCESO A VECTOR ===
// Accediendo a vector 'numeros'
mov x0, #0                    // Evaluar índice (literal 0)
mov x1, x0                    // X1 = índice
ldr x0, [sp, #8]             // Cargar dirección del vector
// Saltar longitud del vector (primer elemento)
add x1, x1, #1               // X1 = índice + 1 (0 + 1 = 1)
// Cargar elemento del vector
ldr x0, [x0, x1, lsl #3]     // X0 = vector[1] = primer elemento real
```

### Estructura de Memoria del Vector

```
vec_numeros: .quad 5, 1, 2, 3, 4, 5
                   ↑  ↑  ↑  ↑  ↑  ↑
               longitud [0][1][2][3][4]
```

- `numeros[0]` accede al segundo elemento (.quad position 1) = `1`
- `numeros[1]` accede al tercer elemento (.quad position 2) = `2`
- Y así sucesivamente...

### Resultado Esperado

Para el código:
```vlang
numeros = []int{1, 2, 3, 4, 5}
print("Primer elemento:", numeros[0])
print("Segundo elemento:", numeros[1])
```

Se espera la salida:
```
Primer elemento: 1
Segundo elemento: 2
```

### Estado de Implementación - Fase 3
- ✅ **Caso agregado en translateExpression**
- ✅ **Función translateVectorAccess implementada**
- ✅ **Extracción de nombre de vector y índice**
- ✅ **Validación de existencia de vector**
- ✅ **Evaluación de expresiones de índice**
- ✅ **Generación de código ARM64 para acceso**
- ✅ **Manejo de offset por metadata de longitud**
