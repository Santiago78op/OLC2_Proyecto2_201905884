// backend/codegen/compiler.go
package codegen

import (
	"fmt"
	"strings"
	"time"

	"github.com/antlr4-go/antlr/v4"
	"main.go/codegen/arm64"
	"main.go/codegen/intermediate"
	"main.go/codegen/output"
	"main.go/errors"
	compiler "main.go/grammar"
	"main.go/repl"
)

// IRCompiler es el compilador principal que integra todo el pipeline
type IRCompiler struct {
	// Componentes del pipeline
	irGenerator  *intermediate.IRGenerator
	irOptimizer  *intermediate.IROptimizer
	assembler    *output.ARM64Assembler
	armOptimizer *arm64.ARM64Optimizer

	// Estado interno
	program  *intermediate.IRProgram
	assembly string
	errors   []repl.Error
	warnings []string

	// Métricas
	stats        CompilerStats
	lastIRString string
}

// CompilerStats contiene estadísticas de compilación
type CompilerStats struct {
	IRGenTime       time.Duration
	IROptimTime     time.Duration
	CodeGenTime     time.Duration
	TotalTime       time.Duration
	IRInstructions  int
	AsmInstructions int
	Optimizations   int
}

// NewIRCompiler crea una nueva instancia del compilador IR
func NewIRCompiler() *IRCompiler {
	return &IRCompiler{
		irGenerator:  intermediate.NewIRGenerator(),
		irOptimizer:  intermediate.NewIROptimizer(),
		assembler:    output.NewARM64Assembler(),
		armOptimizer: arm64.NewARM64Optimizer(),
		errors:       make([]repl.Error, 0),
		warnings:     make([]string, 0),
	}
}

// CompileToIR compila código fuente a IR
func (c *IRCompiler) CompileToIR(sourceCode string) (*intermediate.IRProgram, error) {
	startTime := time.Now()
	c.stats = CompilerStats{} // Reset stats

	fmt.Printf("🚀 Iniciando compilación a IR...\n")

	// ===== FASE 1: ANÁLISIS LÉXICO Y SINTÁCTICO =====
	parseStart := time.Now()

	// Análisis léxico
	lexer := compiler.NewVLangLexer(antlr.NewInputStream(sourceCode))
	lexer.RemoveErrorListeners()

	lexicalErrorListener := errors.NewLexicalErrorListener()
	lexer.AddErrorListener(lexicalErrorListener)

	// Análisis sintáctico
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := compiler.NewVLangGrammar(stream)
	parser.BuildParseTrees = true
	parser.RemoveErrorListeners()

	syntaxErrorListener := errors.NewSyntaxErrorListener(lexicalErrorListener.ErrorTable)
	parser.SetErrorHandler(errors.NewCustomErrorStrategy())
	parser.AddErrorListener(syntaxErrorListener)

	// Generar AST
	tree := parser.Program()

	// Verificar errores de parsing
	if len(syntaxErrorListener.ErrorTable.Errors) > 0 {
		c.errors = syntaxErrorListener.ErrorTable.Errors
		return nil, fmt.Errorf("errores de sintaxis encontrados: %d", len(c.errors))
	}

	fmt.Printf("✅ Análisis sintáctico completado en %v\n", time.Since(parseStart))

	// ===== FASE 2: ANÁLISIS SEMÁNTICO =====
	semanticStart := time.Now()

	// Crear visitor para análisis semántico
	dclVisitor := repl.NewDclVisitor(syntaxErrorListener.ErrorTable)
	dclVisitor.Visit(tree)

	// Verificar errores semánticos
	if len(syntaxErrorListener.ErrorTable.Errors) > 0 {
		c.errors = syntaxErrorListener.ErrorTable.Errors
		return nil, fmt.Errorf("errores semánticos encontrados: %d", len(c.errors))
	}

	fmt.Printf("✅ Análisis semántico completado en %v\n", time.Since(semanticStart))

	// ===== FASE 3: GENERACIÓN DE IR =====
	irGenStart := time.Now()

	c.program = c.irGenerator.GenerateIR(tree, dclVisitor.ScopeTrace)
	if c.program == nil {
		return nil, fmt.Errorf("error generando IR")
	}

	c.stats.IRGenTime = time.Since(irGenStart)
	c.stats.IRInstructions = c.countIRInstructions(c.program)

	fmt.Printf("✅ Generación de IR completada: %d instrucciones en %v\n",
		c.stats.IRInstructions, c.stats.IRGenTime)

	// Actualizar string IR para debugging
	c.lastIRString = c.program.String()

	c.stats.TotalTime = time.Since(startTime)
	return c.program, nil
}

// OptimizeIR aplica optimizaciones al IR
func (c *IRCompiler) OptimizeIR() error {
	if c.program == nil {
		return fmt.Errorf("no hay programa IR para optimizar")
	}

	fmt.Printf("🔧 Iniciando optimización de IR...\n")
	optStart := time.Now()

	originalInstructions := c.countIRInstructions(c.program)

	// Aplicar optimizaciones IR
	c.program = c.irOptimizer.Optimize(c.program)

	optimizedInstructions := c.countIRInstructions(c.program)
	c.stats.IROptimTime = time.Since(optStart)
	c.stats.Optimizations = originalInstructions - optimizedInstructions

	fmt.Printf("✅ Optimización IR completada: %d → %d instrucciones en %v\n",
		originalInstructions, optimizedInstructions, c.stats.IROptimTime)

	// Actualizar string IR optimizado
	c.lastIRString = c.program.String()

	return nil
}

// GenerateARM64 genera código ARM64 assembly
func (c *IRCompiler) GenerateARM64() (string, error) {
	if c.program == nil {
		return "", fmt.Errorf("no hay programa IR para generar ARM64")
	}

	fmt.Printf("🔧 Generando código ARM64...\n")
	codeGenStart := time.Now()

	// Generar assembly ARM64
	assembly, err := c.assembler.AssembleProgram(c.program)
	if err != nil {
		return "", fmt.Errorf("error generando ARM64: %v", err)
	}

	c.assembly = assembly
	c.stats.CodeGenTime = time.Since(codeGenStart)

	// Contar instrucciones en el assembly generado
	c.stats.AsmInstructions = c.countAssemblyInstructions(assembly)

	fmt.Printf("✅ Generación ARM64 completada: %d instrucciones en %v\n",
		c.stats.AsmInstructions, c.stats.CodeGenTime)

	return assembly, nil
}

// OptimizeARM64 aplica optimizaciones específicas de ARM64
func (c *IRCompiler) OptimizeARM64() error {
	if c.assembly == "" {
		return fmt.Errorf("no hay código ARM64 para optimizar")
	}

	fmt.Printf("🔧 Optimizando código ARM64...\n")
	optimizeStart := time.Now()

	// Parsear assembly a instrucciones
	instructions := c.parseAssemblyToInstructions(c.assembly)
	originalCount := len(instructions)

	// Aplicar optimizaciones ARM64
	optimizedInstructions := c.armOptimizer.Optimize(instructions)
	optimizedCount := len(optimizedInstructions)

	// Convertir de vuelta a string
	c.assembly = c.instructionsToAssembly(optimizedInstructions)

	c.armOptimizer.PrintOptimizationStats(originalCount, optimizedCount)

	fmt.Printf("✅ Optimización ARM64 completada en %v\n", time.Since(optimizeStart))

	return nil
}

// ValidateIR valida la representación intermedia
func (c *IRCompiler) ValidateIR() []string {
	if c.program == nil {
		return []string{"No hay programa IR para validar"}
	}

	validationErrors := c.irOptimizer.ValidateIR(c.program)
	c.warnings = append(c.warnings, validationErrors...)

	return validationErrors
}

// GetIRString retorna la representación string del IR
func (c *IRCompiler) GetIRString() string {
	return c.lastIRString
}

// GetAssembly retorna el código ARM64 assembly
func (c *IRCompiler) GetAssembly() string {
	return c.assembly
}

// GetErrors retorna los errores de compilación
func (c *IRCompiler) GetErrors() []repl.Error {
	return c.errors
}

// GetWarnings retorna las advertencias de compilación
func (c *IRCompiler) GetWarnings() []string {
	return c.warnings
}

// GetStats retorna las estadísticas de compilación
func (c *IRCompiler) GetStats() CompilerStats {
	return c.stats
}

// GetOptimizationStats retorna estadísticas de optimización como string
func (c *IRCompiler) GetOptimizationStats() string {
	stats := fmt.Sprintf("Estadísticas de Compilación:\n")
	stats += fmt.Sprintf("  ⏱️ Generación IR: %v\n", c.stats.IRGenTime)
	stats += fmt.Sprintf("  ⏱️ Optimización IR: %v\n", c.stats.IROptimTime)
	stats += fmt.Sprintf("  ⏱️ Generación ARM64: %v\n", c.stats.CodeGenTime)
	stats += fmt.Sprintf("  ⏱️ Tiempo total: %v\n", c.stats.TotalTime)
	stats += fmt.Sprintf("  📊 Instrucciones IR: %d\n", c.stats.IRInstructions)
	stats += fmt.Sprintf("  📊 Instrucciones ARM64: %d\n", c.stats.AsmInstructions)
	stats += fmt.Sprintf("  🎯 Optimizaciones aplicadas: %d\n", c.stats.Optimizations)

	if c.stats.IRInstructions > 0 {
		ratio := float64(c.stats.AsmInstructions) / float64(c.stats.IRInstructions)
		stats += fmt.Sprintf("  📈 Ratio ARM64/IR: %.2fx\n", ratio)
	}

	return stats
}

// CompileFullPipeline ejecuta el pipeline completo de compilación
func (c *IRCompiler) CompileFullPipeline(sourceCode string, optimize bool) (string, error) {
	// 1. Compilar a IR
	_, err := c.CompileToIR(sourceCode)
	if err != nil {
		return "", err
	}

	// 2. Optimizar IR si se solicita
	if optimize {
		err = c.OptimizeIR()
		if err != nil {
			return "", err
		}
	}

	// 3. Generar ARM64
	assembly, err := c.GenerateARM64()

	if err != nil {
		return "", err
	}

	// 4. Optimizar ARM64 si se solicita
	if optimize {
		err = c.OptimizeARM64()
		if err != nil {
			return "", err
		}
	}

	// 5. Validar resultado final
	c.ValidateIR()

	return c.assembly, nil
}

// ============ MÉTODOS AUXILIARES ============

// countIRInstructions cuenta el número total de instrucciones IR
func (c *IRCompiler) countIRInstructions(program *intermediate.IRProgram) int {
	count := 0
	for _, function := range program.Functions {
		count += len(function.Instructions)
	}
	return count
}

// countAssemblyInstructions cuenta el número de instrucciones en assembly
func (c *IRCompiler) countAssemblyInstructions(assembly string) int {
	lines := strings.Split(assembly, "\n")
	count := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Contar líneas que no son comentarios, directivas o vacías
		if trimmed != "" && !strings.HasPrefix(trimmed, ".") &&
			!strings.HasPrefix(trimmed, "//") && !strings.HasSuffix(trimmed, ":") {
			count++
		}
	}

	return count
}

// parseAssemblyToInstructions convierte assembly string a slice de instrucciones
func (c *IRCompiler) parseAssemblyToInstructions(assembly string) []*arm64.ARM64Instruction {
	lines := strings.Split(assembly, "\n")
	var instructions []*arm64.ARM64Instruction

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, ".") {
			continue
		}

		// Parsear línea de assembly (implementación simplificada)
		if strings.HasSuffix(trimmed, ":") {
			// Es una etiqueta
			label := strings.TrimSuffix(trimmed, ":")
			instructions = append(instructions, &arm64.ARM64Instruction{
				Label: label,
			})
		} else {
			// Es una instrucción
			parts := strings.Fields(trimmed)
			if len(parts) > 0 {
				opcode := parts[0]
				operands := []string{}

				if len(parts) > 1 {
					// Unir el resto y dividir por comas
					operandStr := strings.Join(parts[1:], " ")
					if strings.Contains(operandStr, "//") {
						operandStr = strings.Split(operandStr, "//")[0]
					}
					operandStr = strings.TrimSpace(operandStr)

					if operandStr != "" {
						operands = strings.Split(operandStr, ",")
						for i := range operands {
							operands[i] = strings.TrimSpace(operands[i])
						}
					}
				}

				instructions = append(instructions, &arm64.ARM64Instruction{
					Opcode:   opcode,
					Operands: operands,
				})
			}
		}
	}

	return instructions
}

// instructionsToAssembly convierte slice de instrucciones a assembly string
func (c *IRCompiler) instructionsToAssembly(instructions []*arm64.ARM64Instruction) string {
	var lines []string

	for _, instr := range instructions {
		lines = append(lines, instr.String())
	}

	return strings.Join(lines, "\n")
}

// Reset reinicia el estado del compilador
func (c *IRCompiler) Reset() {
	c.program = nil
	c.assembly = ""
	c.errors = make([]repl.Error, 0)
	c.warnings = make([]string, 0)
	c.stats = CompilerStats{}
	c.lastIRString = ""
}
