// backend/main.go - Integración del sistema IR
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/antlr4-go/antlr/v4"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	// Importaciones existentes
	"main.go/ast"
	"main.go/codegen"
	"main.go/codegen/examples"
	"main.go/cst"
	"main.go/errors"
	compiler "main.go/grammar"
	"main.go/repl"
)

// Estructura de configuración para el servidor
type ServerConfig struct {
	Port        string `json:"port"`
	IREnabled   bool   `json:"irEnabled"`
	DebugMode   bool   `json:"debugMode"`
	OptimizeIR  bool   `json:"optimizeIR"`
	ShowIRStats bool   `json:"showIRStats"`
}

// Configuración global del servidor
var serverConfig = ServerConfig{
	Port:        ":8080",
	IREnabled:   true, // ✨ Habilitar IR por defecto
	DebugMode:   true,
	OptimizeIR:  true,
	ShowIRStats: true,
}

// Resultado de ejecución extendido con soporte IR
type ExtendedExecutionResult struct {
	// Campos existentes
	Success         bool                  `json:"success"`
	Errors          []repl.Error          `json:"errors"`
	Output          string                `json:"output"`
	FormattedOutput string                `json:"formattedOutput"`
	ConsoleMessages []repl.ConsoleMessage `json:"consoleMessages"`
	CSTSvg          string                `json:"cstSvg"`
	AST             string                `json:"ast"`
	Symbols         []repl.ReportSymbol   `json:"symbols"`
	ScopeTrace      repl.ReportTable      `json:"scopeTrace"`
	ErrorSummary    map[string]int        `json:"errorSummary"`
	ExecutionTime   int64                 `json:"executionTime"`

	// ✨ Nuevos campos para IR
	IREnabled       bool     `json:"irEnabled"`
	IRGenerated     bool     `json:"irGenerated"`
	IRString        string   `json:"irString,omitempty"`
	IROptimized     bool     `json:"irOptimized"`
	IRStats         string   `json:"irStats,omitempty"`
	IRErrors        []string `json:"irErrors,omitempty"`
	IRWarnings      []string `json:"irWarnings,omitempty"`
	OptimizationLog []string `json:"optimizationLog,omitempty"`

	// Métricas de rendimiento
	CompilationTime  int64 `json:"compilationTime,omitempty"`
	OptimizationTime int64 `json:"optimizationTime,omitempty"`
	ValidationTime   int64 `json:"validationTime,omitempty"`
	InstructionCount int   `json:"instructionCount,omitempty"`
}

func executeCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Leer y procesar el body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("❌ Error leyendo body: %v\n", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	if len(bodyBytes) == 0 {
		fmt.Println("❌ Body está vacío")
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	var requestData struct {
		Code       string `json:"code"`
		EnableIR   bool   `json:"enableIR,omitempty"`   // ✨ Opción para habilitar IR
		OptimizeIR bool   `json:"optimizeIR,omitempty"` // ✨ Opción para optimizar IR
	}

	if err := json.Unmarshal(bodyBytes, &requestData); err != nil {
		fmt.Printf("❌ Error decodificando JSON: %v\n", err)
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	if requestData.Code == "" {
		fmt.Println("❌ Campo 'code' está vacío")
		http.Error(w, "Code field is required and cannot be empty", http.StatusBadRequest)
		return
	}

	// ✨ Determinar si usar IR (por defecto o por petición)
	useIR := serverConfig.IREnabled || requestData.EnableIR
	optimizeIR := serverConfig.OptimizeIR && requestData.OptimizeIR

	codeString := requestData.Code
	// Limpiar código
	for len(codeString) > 0 && (codeString[0] == '\n' || codeString[0] == '\r') {
		codeString = codeString[1:]
	}

	fmt.Printf("✅ Código recibido (IR: %v, Optimizar: %v):\n%s\n", useIR, optimizeIR, codeString)

	// =========== EJECUCIÓN CON SOPORTE IR ===========
	startTime := time.Now()

	if useIR {
		// 🔧 NUEVA RUTA: Procesamiento con IR
		result := executeWithIR(codeString, optimizeIR, startTime)

		// Enviar respuesta
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(result); err != nil {
			fmt.Printf("❌ Error encoding response: %v\n", err)
			return
		}

		fmt.Printf("✅ Respuesta enviada (IR habilitado)\n")
	} else {
		// 🏃 RUTA EXISTENTE: Procesamiento tradicional (compatibilidad)
		result := executeTraditional(codeString, startTime)

		// Enviar respuesta
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(result); err != nil {
			fmt.Printf("❌ Error encoding response: %v\n", err)
			return
		}

		fmt.Printf("✅ Respuesta enviada (modo tradicional)\n")
	}
}

// ✨ NUEVA FUNCIÓN: Ejecución con soporte IR
func executeWithIR(codeString string, shouldOptimize bool, startTime time.Time) *ExtendedExecutionResult {
	result := &ExtendedExecutionResult{
		IREnabled: true,
	}

	fmt.Printf("🔧 Procesando con sistema IR...\n")

	// ===== FASE 1: COMPILACIÓN A IR =====
	compileStart := time.Now()

	irCompiler := codegen.NewIRCompiler()
	program, err := irCompiler.CompileToIR(codeString)

	compilationTime := time.Since(compileStart)
	result.CompilationTime = compilationTime.Milliseconds()

	if err != nil {
		fmt.Printf("❌ Error compilando a IR: %v\n", err)
		result.Success = false
		result.IRGenerated = false
		result.Errors = irCompiler.GetErrors()
		result.ErrorSummary = getErrorSummary(result.Errors)
		result.ExecutionTime = time.Since(startTime).Milliseconds()
		return result
	}

	result.IRGenerated = true
	result.IRString = irCompiler.GetIRString()

	// Contar instrucciones
	instructionCount := 0
	if program != nil {
		for _, function := range program.Functions {
			instructionCount += len(function.Instructions)
		}
	}
	result.InstructionCount = instructionCount

	fmt.Printf("✅ IR generado: %d instrucciones\n", instructionCount)

	// ===== FASE 2: OPTIMIZACIÓN (OPCIONAL) =====
	if shouldOptimize {
		optimizeStart := time.Now()

		err = irCompiler.OptimizeIR()
		optimizationTime := time.Since(optimizeStart)
		result.OptimizationTime = optimizationTime.Milliseconds()

		if err != nil {
			fmt.Printf("⚠️ Error optimizando IR: %v\n", err)
			result.IRErrors = append(result.IRErrors, err.Error())
		} else {
			result.IROptimized = true
			result.IRString = irCompiler.GetIRString() // IR optimizado
			fmt.Printf("🔧 IR optimizado en %v\n", optimizationTime)
		}
	}

	// ===== FASE 3: VALIDACIÓN =====
	validateStart := time.Now()

	validationErrors := irCompiler.ValidateIR()
	validationTime := time.Since(validateStart)
	result.ValidationTime = validationTime.Milliseconds()

	if len(validationErrors) > 0 {
		result.IRWarnings = validationErrors
		fmt.Printf("⚠️ %d advertencias de validación IR\n", len(validationErrors))
	}

	// ===== FASE 4: GENERAR ESTADÍSTICAS =====
	if serverConfig.ShowIRStats {
		result.IRStats = irCompiler.GetOptimizationStats()
	}

	// ===== FASE 5: FALLBACK A INTERPRETACIÓN TRADICIONAL =====
	// Por ahora, usar el intérprete existente para la ejecución real
	// TODO: En el futuro, aquí iría la generación de código ARM64

	fmt.Printf("📋 Ejecutando con intérprete tradicional (IR como análisis)...\n")

	traditionalResult := executeTraditionalCore(codeString, startTime)

	// Combinar resultados
	result.Success = traditionalResult.Success
	result.Output = traditionalResult.Output
	result.FormattedOutput = traditionalResult.FormattedOutput
	result.ConsoleMessages = traditionalResult.ConsoleMessages
	result.CSTSvg = traditionalResult.CSTSvg
	result.AST = traditionalResult.AST
	result.Symbols = traditionalResult.Symbols
	result.ScopeTrace = traditionalResult.ScopeTrace
	result.ExecutionTime = time.Since(startTime).Milliseconds()

	// Combinar errores (IR + tradicional)
	allErrors := append(result.Errors, traditionalResult.Errors...)
	result.Errors = allErrors
	result.ErrorSummary = getErrorSummary(allErrors)

	fmt.Printf("✅ Procesamiento con IR completado\n")
	return result
}

// 🏃 FUNCIÓN EXISTENTE: Ejecución tradicional (sin cambios)
func executeTraditional(codeString string, startTime time.Time) *ExtendedExecutionResult {
	traditionalResult := executeTraditionalCore(codeString, startTime)

	// Convertir a formato extendido
	return &ExtendedExecutionResult{
		Success:         traditionalResult.Success,
		Errors:          traditionalResult.Errors,
		Output:          traditionalResult.Output,
		FormattedOutput: traditionalResult.FormattedOutput,
		ConsoleMessages: traditionalResult.ConsoleMessages,
		CSTSvg:          traditionalResult.CSTSvg,
		AST:             traditionalResult.AST,
		Symbols:         traditionalResult.Symbols,
		ScopeTrace:      traditionalResult.ScopeTrace,
		ErrorSummary:    traditionalResult.ErrorSummary,
		ExecutionTime:   traditionalResult.ExecutionTime,
		IREnabled:       false,
	}
}

// Núcleo de ejecución tradicional (código existente reorganizado)
func executeTraditionalCore(codeString string, startTime time.Time) *executionResult {
	// 1. Generar CST Report en paralelo
	cstChannel := make(chan string, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Error generando CST Report:", r)
				cstChannel <- ""
			}
		}()
		cstChannel <- cst.CstReport(codeString)
	}()

	// 2. Análisis Léxico
	lexicalErrorListener := errors.NewLexicalErrorListener()
	lexer := compiler.NewVLangLexer(antlr.NewInputStream(codeString))
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(lexicalErrorListener)

	// 3. Análisis Sintáctico
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := compiler.NewVLangGrammar(stream)
	parser.BuildParseTrees = true

	syntaxErrorListener := errors.NewSyntaxErrorListener(lexicalErrorListener.ErrorTable)
	parser.RemoveErrorListeners()
	parser.SetErrorHandler(errors.NewCustomErrorStrategy())
	parser.AddErrorListener(syntaxErrorListener)

	// 4. Generar AST
	tree := parser.Program()

	// Verificar si hubo errores críticos
	hasCompilationErrors := len(syntaxErrorListener.ErrorTable.Errors) > 0

	var replVisitor *repl.ReplVisitor
	var output string = ""
	var formattedOutput string = ""
	var consoleMessages []repl.ConsoleMessage

	// 5. Solo continuar con análisis semántico si no hay errores críticos
	if !hasCompilationErrors {
		// Análisis Semántico y Ejecución
		dclVisitor := repl.NewDclVisitor(syntaxErrorListener.ErrorTable)
		dclVisitor.Visit(tree)

		replVisitor = repl.NewVisitor(dclVisitor)
		replVisitor.Visit(tree)
		output = replVisitor.Console.GetOutput()
		formattedOutput = replVisitor.Console.GetFormattedOutput()
		consoleMessages = replVisitor.Console.GetMessages()
	} else {
		// Si hay errores de compilación, crear visitor básico para reportes
		dclVisitor := repl.NewDclVisitor(syntaxErrorListener.ErrorTable)
		replVisitor = repl.NewVisitor(dclVisitor)
		output = ""
	}

	interpretationEndTime := time.Now()

	// 6. Obtener CST Report
	cstReport := <-cstChannel

	// 7. Generar AST nativo
	var finalAST string
	if tree != nil && !hasCompilationErrors {
		astNode := ast.GenerateNativeAST(tree)
		if astNode != nil {
			finalAST = ast.GenerateASTSVG(astNode)
		}
	}

	if cstReport == "" && finalAST != "" {
		cstReport = finalAST
	}

	// Determinar si la ejecución fue exitosa
	success := !hasCompilationErrors && len(syntaxErrorListener.ErrorTable.Errors) == 0

	// Generar tabla de símbolos
	scopeReport := replVisitor.ScopeTrace.Report()
	symbols := extractSymbolsFromScope(scopeReport)

	// Crear resumen de errores
	errorSummary := syntaxErrorListener.ErrorTable.GetErrorsSummary()

	// Crear resultado
	return &executionResult{
		Success:         success,
		Errors:          syntaxErrorListener.ErrorTable.Errors,
		Output:          output,
		FormattedOutput: formattedOutput,
		ConsoleMessages: consoleMessages,
		CSTSvg:          cstReport,
		AST:             finalAST,
		Symbols:         symbols,
		ScopeTrace:      scopeReport,
		ErrorSummary:    errorSummary,
		ExecutionTime:   interpretationEndTime.Sub(startTime).Milliseconds(),
	}
}

// ✨ NUEVOS ENDPOINTS PARA IR

// Endpoint para obtener información del sistema IR
func getIRInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	info := map[string]interface{}{
		"irEnabled":   serverConfig.IREnabled,
		"debugMode":   serverConfig.DebugMode,
		"optimizeIR":  serverConfig.OptimizeIR,
		"showIRStats": serverConfig.ShowIRStats,
		"version":     "1.0.0-beta",
		"capabilities": []string{
			"IR Generation",
			"Basic Optimizations",
			"Constant Folding",
			"Dead Code Elimination",
			"Validation",
			"Metrics Collection",
		},
		"supportedOptimizations": []string{
			"Constant Propagation",
			"Constant Folding",
			"Dead Code Elimination",
			"Redundant Load Removal",
			"Peephole Optimizations",
			"Algebraic Simplification",
		},
	}

	json.NewEncoder(w).Encode(info)
}

// Endpoint para configurar el sistema IR
func configureIR(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var config ServerConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Actualizar configuración
	serverConfig.IREnabled = config.IREnabled
	serverConfig.DebugMode = config.DebugMode
	serverConfig.OptimizeIR = config.OptimizeIR
	serverConfig.ShowIRStats = config.ShowIRStats

	fmt.Printf("🔧 Configuración IR actualizada: IR=%v, Debug=%v, Optimize=%v, Stats=%v\n",
		config.IREnabled, config.DebugMode, config.OptimizeIR, config.ShowIRStats)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Configuración actualizada",
		"config":  serverConfig,
	})
}

// Endpoint para ejecutar demo del IR
func runIRDemo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Printf("🎯 Ejecutando demo del sistema IR...\n")

	// Capturar salida del demo (en una implementación real usaríamos un logger)
	// Por simplicidad, solo ejecutamos el demo y retornamos estado
	go examples.RunIRDemo()

	response := map[string]interface{}{
		"success":   true,
		"message":   "Demo del sistema IR iniciado",
		"note":      "Revisa los logs del servidor para ver el output del demo",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(response)
}

// Endpoint para análisis de código específico con IR
func analyzeIR(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var requestData struct {
		Code      string `json:"code"`
		Optimize  bool   `json:"optimize"`
		ShowStats bool   `json:"showStats"`
		Validate  bool   `json:"validate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if requestData.Code == "" {
		http.Error(w, "Code field is required", http.StatusBadRequest)
		return
	}

	// Compilar solo a IR (sin ejecución)
	irCompiler := codegen.NewIRCompiler()
	program, err := irCompiler.CompileToIR(requestData.Code)

	result := map[string]interface{}{
		"success": err == nil,
	}

	if err != nil {
		result["error"] = err.Error()
		result["errors"] = irCompiler.GetErrors()
	} else {
		result["irString"] = irCompiler.GetIRString()

		// Optimizar si se solicita
		if requestData.Optimize {
			err = irCompiler.OptimizeIR()
			if err != nil {
				result["optimizationError"] = err.Error()
			} else {
				result["optimized"] = true
				result["irStringOptimized"] = irCompiler.GetIRString()
			}
		}

		// Estadísticas si se solicitan
		if requestData.ShowStats {
			result["stats"] = irCompiler.GetOptimizationStats()
		}

		// Validación si se solicita
		if requestData.Validate {
			validationErrors := irCompiler.ValidateIR()
			result["validationErrors"] = validationErrors
			result["valid"] = len(validationErrors) == 0
		}

		// Información básica del programa
		if program != nil {
			instructionCount := 0
			for _, function := range program.Functions {
				instructionCount += len(function.Instructions)
			}

			result["programInfo"] = map[string]interface{}{
				"functionCount":    len(program.Functions),
				"instructionCount": instructionCount,
				"globalVarCount":   len(program.GlobalVars),
				"stringTableSize":  len(program.StringTable),
			}
		}
	}

	json.NewEncoder(w).Encode(result)
}

// =============== UTILIDADES AUXILIARES ===============

// Función auxiliar para obtener resumen de errores
func getErrorSummary(errors []repl.Error) map[string]int {
	summary := make(map[string]int)
	for _, err := range errors {
		summary[err.Type]++
	}
	return summary
}

// Estructura de resultado original (sin cambios)
type executionResult struct {
	Success         bool                  `json:"success"`
	Errors          []repl.Error          `json:"errors"`
	Output          string                `json:"output"`
	FormattedOutput string                `json:"formattedOutput"`
	ConsoleMessages []repl.ConsoleMessage `json:"consoleMessages"`
	CSTSvg          string                `json:"cstSvg"`
	AST             string                `json:"ast"`
	Symbols         []repl.ReportSymbol   `json:"symbols"`
	ScopeTrace      repl.ReportTable      `json:"scopeTrace"`
	ErrorSummary    map[string]int        `json:"errorSummary"`
	ExecutionTime   int64                 `json:"executionTime"`
}

// Función auxiliar existente
func extractSymbolsFromScope(scopeReport repl.ReportTable) []repl.ReportSymbol {
	var allSymbols []repl.ReportSymbol

	var extractFromScope func(scope repl.ReportScope, scopeName string)
	extractFromScope = func(scope repl.ReportScope, scopeName string) {
		for _, symbol := range scope.Vars {
			symbol.Scope = scopeName
			allSymbols = append(allSymbols, symbol)
		}

		for _, symbol := range scope.Funcs {
			symbol.Scope = scopeName
			allSymbols = append(allSymbols, symbol)
		}

		for _, symbol := range scope.Structs {
			symbol.Scope = scopeName
			allSymbols = append(allSymbols, symbol)
		}

		for _, childScope := range scope.ChildScopes {
			childScopeName := scopeName + "." + childScope.Name
			extractFromScope(childScope, childScopeName)
		}
	}

	extractFromScope(scopeReport.GlobalScope, "global")
	return allSymbols
}

// Endpoint de health check (existente, sin cambios)
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// =============== FUNCIÓN PRINCIPAL ===============

func main() {
	r := mux.NewRouter()

	// API Routes existentes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/status", healthCheck).Methods("GET")
	api.HandleFunc("/execute", executeCode).Methods("POST")

	// ✨ NUEVAS RUTAS PARA IR
	api.HandleFunc("/ir/info", getIRInfo).Methods("GET")
	api.HandleFunc("/ir/config", configureIR).Methods("POST")
	api.HandleFunc("/ir/demo", runIRDemo).Methods("POST")
	api.HandleFunc("/ir/analyze", analyzeIR).Methods("POST")

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(r)

	port := ":8080"
	fmt.Printf("🚀 Servidor VLan Cherry iniciado en http://localhost%s\n", port)
	fmt.Printf("📋 Configuración inicial:\n")
	fmt.Printf("   🔧 IR habilitado: %v\n", serverConfig.IREnabled)
	fmt.Printf("   🐛 Debug mode: %v\n", serverConfig.DebugMode)
	fmt.Printf("   ⚡ Optimizar IR: %v\n", serverConfig.OptimizeIR)
	fmt.Printf("   📊 Mostrar stats: %v\n", serverConfig.ShowIRStats)
	fmt.Println()
	fmt.Println("📋 API endpoints disponibles:")
	fmt.Println("  Existentes:")
	fmt.Println("    - GET    /api/status")
	fmt.Println("    - POST   /api/execute")
	fmt.Println("  ✨ Nuevos (IR):")
	fmt.Println("    - GET    /api/ir/info")
	fmt.Println("    - POST   /api/ir/config")
	fmt.Println("    - POST   /api/ir/demo")
	fmt.Println("    - POST   /api/ir/analyze")
	fmt.Println()

	// Ejecutar demo inicial si está en modo debug
	if serverConfig.DebugMode {
		fmt.Printf("🎯 Ejecutando demo inicial del sistema IR...\n")
		go func() {
			time.Sleep(2 * time.Second) // Esperar a que el servidor inicie
			examples.RunIRDemo()
		}()
	}

	log.Fatal(http.ListenAndServe(port, handler))
}
