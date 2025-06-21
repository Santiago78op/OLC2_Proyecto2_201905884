// backend/codegen/compiler.go
package codegen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

	// ✨ NUEVOS COMPONENTES PARA LINKING
	linker     *output.ARM64Linker
	workingDir string

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
	// Obtener directorio de trabajo
	workingDir, _ := os.Getwd()
	if workingDir == "" {
		workingDir = "/tmp"
	}

	return &IRCompiler{
		irGenerator:  intermediate.NewIRGenerator(),
		irOptimizer:  intermediate.NewIROptimizer(),
		assembler:    output.NewARM64Assembler(),
		armOptimizer: arm64.NewARM64Optimizer(),

		// ✨ NUEVOS INICIALIZADORES
		linker:     output.NewARM64Linker(workingDir),
		workingDir: workingDir,

		errors:   make([]repl.Error, 0),
		warnings: make([]string, 0),
	}
}

// =============== PIPELINE BÁSICO (SIN CAMBIOS) ===============

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
	parser := compiler.NewVLangGrammarParser(stream)
	parser.RemoveErrorListeners()

	syntaxErrorListener := errors.NewSyntaxErrorListener()
	parser.AddErrorListener(syntaxErrorListener)

	// Parse tree
	tree := parser.Program()

	parseTime := time.Since(parseStart)
	fmt.Printf("✅ Análisis completado en %v\n", parseTime)

	// Verificar errores de parsing
	if len(lexicalErrorListener.ErrorTable.Errors) > 0 || len(syntaxErrorListener.ErrorTable.Errors) > 0 {
		// Combinar errores
		c.errors = append(c.errors, lexicalErrorListener.ErrorTable.Errors...)
		c.errors = append(c.errors, syntaxErrorListener.ErrorTable.Errors...)

		return nil, fmt.Errorf("errores de compilación encontrados: %d", len(c.errors))
	}

	// ===== FASE 2: GENERACIÓN DE IR =====
	fmt.Printf("🔧 Generando representación intermedia...\n")
	irGenStart := time.Now()

	// Crear scope trace para la generación IR
	scopeTrace := repl.NewScopeTrace()

	// Generar IR
	c.program = c.irGenerator.GenerateIR(tree, scopeTrace)
	c.stats.IRInstructions = c.countIRInstructions(c.program)
	c.stats.IRGenTime = time.Since(irGenStart)

	fmt.Printf("✅ IR generado: %d instrucciones en %v\n", c.stats.IRInstructions, c.stats.IRGenTime)

	// Guardar IR string
	c.lastIRString = c.program.String()

	// ===== FASE 3: ACTUALIZAR MÉTRICAS =====
	c.stats.TotalTime = time.Since(startTime)

	return c.program, nil
}

// OptimizeIR aplica optimizaciones al IR
func (c *IRCompiler) OptimizeIR() error {
	if c.program == nil {
		return fmt.Errorf("no hay programa IR para optimizar")
	}

	fmt.Printf("🔧 Optimizando IR...\n")
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

// =============== ✨ NUEVAS FUNCIONES PARA LINKING ===============

// CompileToExecutable - NUEVA FUNCIÓN que extiende el pipeline para generar ejecutables
func (c *IRCompiler) CompileToExecutable(sourceCode string, outputName string, optimize bool) (*output.LinkingResult, error) {
	fmt.Printf("🚀 Compilando a ejecutable ARM64: %s\n", outputName)

	// 1. Usar el pipeline existente para generar assembly
	assembly, err := c.CompileFullPipeline(sourceCode, optimize)
	if err != nil {
		return nil, fmt.Errorf("error en pipeline de compilación: %v", err)
	}

	// 2. Configurar opciones de enlazado
	linkOptions := output.LinkingOptions{
		OutputName:    outputName,
		EntryPoint:    "main",
		Libraries:     []string{"c"},
		StaticLink:    false,
		OptimizeSize:  false,
		DebugInfo:     true,
		StripSymbols:  false,
		KeepTempFiles: false,
	}

	// 3. Enlazar ejecutable
	fmt.Printf("🔗 Enlazando ejecutable...\n")
	linkResult, err := c.linker.LinkExecutable(assembly, linkOptions)
	if err != nil {
		return nil, fmt.Errorf("error enlazando ejecutable: %v", err)
	}

	fmt.Printf("✅ Ejecutable creado: %s (%d bytes)\n",
		linkResult.ExecutablePath, linkResult.FileSize)

	return linkResult, nil
}

// CompileAndRun - NUEVA FUNCIÓN que compila y ejecuta
func (c *IRCompiler) CompileAndRun(sourceCode string, outputName string, args []string, optimize bool) (*output.LinkingResult, string, error) {
	// 1. Compilar a ejecutable
	linkResult, err := c.CompileToExecutable(sourceCode, outputName, optimize)
	if err != nil {
		return linkResult, "", err
	}

	// 2. Ejecutar
	fmt.Printf("🎯 Ejecutando programa...\n")
	cmd := exec.Command(linkResult.ExecutablePath, args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return linkResult, string(output), fmt.Errorf("error ejecutando %s: %v", linkResult.ExecutablePath, err)
	}

	fmt.Printf("📤 Programa ejecutado exitosamente\n")
	return linkResult, string(output), nil
}

// ValidateEnvironment - NUEVA FUNCIÓN para validar entorno
func (c *IRCompiler) ValidateEnvironment() error {
	return c.linker.ValidateEnvironment()
}

// SetWorkingDirectory - NUEVA FUNCIÓN para configurar directorio
func (c *IRCompiler) SetWorkingDirectory(dir string) {
	c.workingDir = dir
	c.linker = output.NewARM64Linker(dir)
}

// GetOptimizedAssembly - NUEVA FUNCIÓN que ya necesitas
func (c *IRCompiler) GetOptimizedAssembly() string {
	return c.assembly
}

// CleanBuildDirectory - NUEVA FUNCIÓN para limpiar build
func (c *IRCompiler) CleanBuildDirectory() error {
	buildDir := filepath.Join(c.workingDir, "build")

	if _, err := os.Stat(buildDir); os.IsNotExist(err) {
		return nil // No existe, no hay nada que limpiar
	}

	entries, err := os.ReadDir(buildDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		entryPath := filepath.Join(buildDir, entry.Name())
		if err := os.RemoveAll(entryPath); err != nil {
			return err
		}
	}

	fmt.Printf("🧹 Directorio de build limpiado\n")
	return nil
}

// =============== FUNCIONES EXISTENTES (SIN CAMBIOS) ===============

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

// Reset reinicia el estado del compilador
func (c *IRCompiler) Reset() {
	c.program = nil
	c.assembly = ""
	c.errors = make([]repl.Error, 0)
	c.warnings = make([]string, 0)
	c.stats = CompilerStats{}
	c.lastIRString = ""
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
