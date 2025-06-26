package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/antlr4-go/antlr/v4"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	// Importa el paquete de pruebas que contiene la lógica de ejecución

	"main.go/ast"
	compiler "main.go/compiler" // NUEVA: nuestro traductor ARM64
	"main.go/cst"
	"main.go/errors"
	interpeter "main.go/grammar"
	"main.go/repl"
)

type executionResult struct {
	Success         bool                  `json:"success"`
	Errors          []repl.Error          `json:"errors"`
	Output          string                `json:"output"`          // Output plano para retrocompatibilidad
	FormattedOutput string                `json:"formattedOutput"` // Output con formato mejorado
	ConsoleMessages []repl.ConsoleMessage `json:"consoleMessages"` // Mensajes estructurados
	CSTSvg          string                `json:"cstSvg"`
	AST             string                `json:"ast"`
	Symbols         []repl.ReportSymbol   `json:"symbols"`
	ScopeTrace      repl.ReportTable      `json:"scopeTrace"`
	ErrorSummary    map[string]int        `json:"errorSummary"`
	ExecutionTime   int64                 `json:"executionTime"`

	// NUEVOS CAMPOS PARA ARM64
	ARM64Code   string   `json:"arm64Code"`   // Código ARM64 generado
	ARM64Errors []string `json:"arm64Errors"` // Errores de traducción
	HasARM64    bool     `json:"hasArm64"`    // Si se generó código ARM64
}

// Agregar esta función para limpiar el HTML del código ARM64
func stripHTMLFromAssembly(htmlCode string) string {
	// Eliminar etiquetas HTML comunes
	htmlCode = regexp.MustCompile(`<span[^>]*>`).ReplaceAllString(htmlCode, "")
	htmlCode = regexp.MustCompile(`</span>`).ReplaceAllString(htmlCode, "")
	htmlCode = regexp.MustCompile(`"asm-[^"]*">`).ReplaceAllString(htmlCode, "")

	// Reemplazar entidades HTML comunes
	htmlCode = strings.ReplaceAll(htmlCode, "&lt;", "<")
	htmlCode = strings.ReplaceAll(htmlCode, "&gt;", ">")
	htmlCode = strings.ReplaceAll(htmlCode, "&amp;", "&")
	htmlCode = strings.ReplaceAll(htmlCode, "&quot;", "\"")

	return htmlCode
}

// Función para traducir a ARM64
func translateToARM64(tree antlr.ParseTree) (string, []string, bool) {
	fmt.Printf("🔹 Iniciando traducción a ARM64...\n")

	// Crear el traductor
	translator := compiler.NewARM64Translator()

	// Traducir el programa
	arm64Code, errors := translator.TranslateProgram(tree)

	// Limpiar cualquier HTML del código generado
	arm64Code = stripHTMLFromAssembly(arm64Code)

	// Asegurarse de que el código esté limpio antes de enviarlo al frontend
	if strings.Contains(arm64Code, "<span") ||
		strings.Contains(arm64Code, "</span") ||
		strings.Contains(arm64Code, "asm-") {
		fmt.Println("⚠️ Se detectaron etiquetas HTML en el código, aplicando limpieza adicional...")
		arm64Code = deepCleanAssemblyCode(arm64Code)
	}

	// Log para debugging
	fmt.Println("🔎 Primeros 100 caracteres del código ARM64 limpio:")
	if len(arm64Code) > 100 {
		fmt.Println(arm64Code[:100] + "...")
	} else {
		fmt.Println(arm64Code)
	}

	if len(errors) > 0 {
		fmt.Printf("❌ Errores en traducción ARM64: %d\n", len(errors))
		for _, err := range errors {
			fmt.Printf("   - %s\n", err)
		}
	} else {
		fmt.Printf("✅ Traducción ARM64 exitosa\n")
	}

	return arm64Code, errors, len(errors) == 0
}

// Función de limpieza profunda para asegurar que no hay rastros de HTML
func deepCleanAssemblyCode(code string) string {
	// Remover cualquier etiqueta HTML posiblemente mal formada
	cleanCode := code

	// Eliminar spans y atributos
	cleanCode = regexp.MustCompile(`<span[^>]*>`).ReplaceAllString(cleanCode, "")
	cleanCode = regexp.MustCompile(`</span>`).ReplaceAllString(cleanCode, "")

	// Eliminar fragmentos de etiquetas incompletas como "asm-comment">
	cleanCode = regexp.MustCompile(`"asm-[^"]*">`).ReplaceAllString(cleanCode, "")
	cleanCode = regexp.MustCompile(`class="[^"]*"`).ReplaceAllString(cleanCode, "")

	// Eliminar otros fragmentos HTML comunes
	cleanCode = regexp.MustCompile(`</?[a-z]+[^>]*>`).ReplaceAllString(cleanCode, "")

	// Reemplazar entidades HTML
	cleanCode = strings.ReplaceAll(cleanCode, "&lt;", "<")
	cleanCode = strings.ReplaceAll(cleanCode, "&gt;", ">")
	cleanCode = strings.ReplaceAll(cleanCode, "&amp;", "&")
	cleanCode = strings.ReplaceAll(cleanCode, "&quot;", "\"")

	return cleanCode
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

	fmt.Printf("🔹 Body raw recibido: %s\n", string(bodyBytes))

	if len(bodyBytes) == 0 {
		fmt.Println("❌ Body está vacío")
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	var requestData struct {
		Code string `json:"code"`
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

	codeString := requestData.Code
	for len(codeString) > 0 && (codeString[0] == '\n' || codeString[0] == '\r') {
		codeString = codeString[1:]
	}

	fmt.Printf("✅ Código recibido exitosamente:\n%s\n", codeString)

	// =========== ANÁLISIS Y EJECUCIÓN ===========
	startTime := time.Now()

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
	lexer := interpeter.NewVLangLexer(antlr.NewInputStream(codeString))
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(lexicalErrorListener)

	// 3. Análisis Sintáctico
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := interpeter.NewVLangGrammar(stream)
	parser.BuildParseTrees = true

	syntaxErrorListener := errors.NewSyntaxErrorListener(lexicalErrorListener.ErrorTable)
	parser.RemoveErrorListeners()
	parser.SetErrorHandler(errors.NewCustomErrorStrategy())
	parser.AddErrorListener(syntaxErrorListener)

	// 4. Generar AST
	tree := parser.Program()

	// Verificar si hubo errores críticos
	hasCompilationErrors := len(syntaxErrorListener.ErrorTable.Errors) > 0

	fmt.Printf("🔹 Errores de compilación: %d\n", len(syntaxErrorListener.ErrorTable.Errors))
	if hasCompilationErrors {
		for _, err := range syntaxErrorListener.ErrorTable.Errors {
			fmt.Printf("   - %s (Línea %d, Col %d): %s\n", err.GetDisplayName(), err.Line, err.Column, err.Msg)
		}
	}

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
		fmt.Println("🌳 Generando AST nativo...")

		// Generar AST con timeout para evitar bloqueos
		astChannel := make(chan string, 1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("❌ Error generando AST nativo: %v\n", r)
					astChannel <- generateErrorAST("Error al generar AST")
				}
			}()

			astNode := ast.GenerateNativeAST(tree)
			if astNode != nil {
				astChannel <- ast.GenerateASTSVG(astNode)
			} else {
				astChannel <- generateErrorAST("No se pudo generar el árbol")
			}
		}()

		// Esperar con timeout
		select {
		case finalAST = <-astChannel:
			fmt.Println("✅ AST nativo generado exitosamente")
		case <-time.After(5 * time.Second):
			fmt.Println("⏱️ Timeout generando AST")
			finalAST = generateErrorAST("Timeout al generar AST")
		}
	} else {
		fmt.Println("❌ No se pudo generar el árbol de análisis debido a errores")
		finalAST = generateErrorAST("Error en análisis sintáctico")
	}

	// Si no hay CST report pero sí AST nativo, usar el AST nativo
	if cstReport == "" && finalAST != "" {
		cstReport = finalAST
	}

	reportEndTime := time.Now()

	// =========== GENERAR REPORTES ===========

	// Determinar si la ejecución fue exitosa
	success := !hasCompilationErrors && len(syntaxErrorListener.ErrorTable.Errors) == 0

	// Generar tabla de símbolos
	scopeReport := replVisitor.ScopeTrace.Report()
	symbols := extractSymbolsFromScope(scopeReport)

	// Crear resumen de errores
	errorSummary := syntaxErrorListener.ErrorTable.GetErrorsSummary()

	fmt.Printf("🔹 Resumen de errores: %+v\n", errorSummary)
	fmt.Printf("🔹 Tiempo de interpretación: %v\n", interpretationEndTime.Sub(startTime))
	fmt.Printf("🔹 Tiempo total: %v\n", reportEndTime.Sub(startTime))
	fmt.Printf("🔹 Salida: %s\n", output)

	// =========== TRADUCCIÓN A ARM64 ===========
	var arm64Code string
	var arm64Errors []string
	var hasValidARM64 bool

	// Solo intentar traducir a ARM64 si no hay errores de compilación
	if !hasCompilationErrors {
		fmt.Printf("🔹 Intentando generar código ARM64...\n")
		arm64Code, arm64Errors, hasValidARM64 = translateToARM64(tree)
		fmt.Print("Codigo Arm64 \n", arm64Code)
	} else {
		arm64Code = ""
		arm64Errors = []string{"No se puede generar ARM64 debido a errores de compilación"}
		hasValidARM64 = false
	}

	// Crear resultado con información detallada
	result := executionResult{
		Success:         success,
		Errors:          syntaxErrorListener.ErrorTable.Errors,
		Output:          output,
		FormattedOutput: formattedOutput,
		ConsoleMessages: consoleMessages,
		CSTSvg:          cstReport, // CST del servicio externo
		AST:             finalAST,  // AST nativo generado
		Symbols:         symbols,
		ScopeTrace:      scopeReport,
		ErrorSummary:    errorSummary,
		ExecutionTime:   interpretationEndTime.Sub(startTime).Milliseconds(),

		// NUEVOS CAMPOS ARM64
		ARM64Code:   arm64Code,
		ARM64Errors: arm64Errors,
		HasARM64:    hasValidARM64,
	}

	// Enviar respuesta
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		fmt.Printf("❌ Error encoding response: %v\n", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	fmt.Printf("✅ Respuesta enviada exitosamente\n")
}

// Función auxiliar para generar AST de error
func generateErrorAST(errorMsg string) string {
	return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="600" height="200" viewBox="0 0 600 200">
		<rect width="600" height="200" fill="#1e1e1e"/>
		<text x="300" y="90" text-anchor="middle" fill="#ff6b6b" font-family="Arial" font-size="18">
			⚠️ %s
		</text>
		<text x="300" y="120" text-anchor="middle" fill="#cccccc" font-family="Arial" font-size="14">
			Verifica que el código tenga sintaxis válida
		</text>
	</svg>`, errorMsg)
}

// Función auxiliar para extraer símbolos del scope report
func extractSymbolsFromScope(scopeReport repl.ReportTable) []repl.ReportSymbol {
	var allSymbols []repl.ReportSymbol

	// Función recursiva para extraer símbolos de todos los scopes
	var extractFromScope func(scope repl.ReportScope, scopeName string)
	extractFromScope = func(scope repl.ReportScope, scopeName string) {
		// Agregar variables
		for _, symbol := range scope.Vars {
			symbol.Scope = scopeName
			allSymbols = append(allSymbols, symbol)
		}

		// Agregar funciones
		for _, symbol := range scope.Funcs {
			symbol.Scope = scopeName
			allSymbols = append(allSymbols, symbol)
		}

		// Agregar estructuras
		for _, symbol := range scope.Structs {
			symbol.Scope = scopeName
			allSymbols = append(allSymbols, symbol)
		}

		// Procesar scopes hijos recursivamente
		for _, childScope := range scope.ChildScopes {
			childScopeName := scopeName + "." + childScope.Name
			extractFromScope(childScope, childScopeName)
		}
	}

	extractFromScope(scopeReport.GlobalScope, "global")
	return allSymbols
}

// Función auxiliar para generar AST básico si falla el CST report
func generateBasicAST(tree antlr.ParseTree) string {
	if tree == nil {
		return generateErrorAST("No se pudo generar el AST")
	}

	// Generar un SVG básico con información del árbol
	return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="600" height="400" viewBox="0 0 600 400">
		<rect width="600" height="400" fill="#1e1e1e"/>
		<circle cx="300" cy="100" r="40" fill="#007acc" stroke="#ffffff" stroke-width="2"/>
		<text x="300" y="105" text-anchor="middle" fill="#ffffff" font-family="Arial" font-size="12">Program</text>
		<text x="300" y="200" text-anchor="middle" fill="#cccccc" font-family="Arial" font-size="14">
			AST generado exitosamente
		</text>
		<text x="300" y="220" text-anchor="middle" fill="#cccccc" font-family="Arial" font-size="12">
			Texto del árbol: %s
		</text>
	</svg>`, tree.GetText()[:min(50, len(tree.GetText()))])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// Endpoint para obtener solo el código ARM64
func getARM64Code(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Leer el código fuente del request
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	var requestData struct {
		Code string `json:"code"`
	}

	if err := json.Unmarshal(bodyBytes, &requestData); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Análisis léxico y sintáctico
	lexer := interpeter.NewVLangLexer(antlr.NewInputStream(requestData.Code))
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := interpeter.NewVLangGrammar(stream)
	tree := parser.Program()

	// Traducir a ARM64
	arm64Code, arm64Errors, success := translateToARM64(tree)

	// Respuesta
	response := map[string]interface{}{
		"success":   success,
		"arm64Code": arm64Code,
		"errors":    arm64Errors,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func executeARM64Code(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Leer el código ARM64 del request
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	var requestData struct {
		ARM64Code string `json:"arm64Code"`
	}

	if err := json.Unmarshal(bodyBytes, &requestData); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	if requestData.ARM64Code == "" {
		http.Error(w, "ARM64 code is required", http.StatusBadRequest)
		return
	}

	// Ejecutar el código ARM64
	output, success, errorMsg := executeARM64Assembly(requestData.ARM64Code)

	// Respuesta
	response := map[string]interface{}{
		"success":   success,
		"output":    output,
		"error":     errorMsg,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Función para ejecutar código ARM64
// Función para ejecutar código ARM64 con pre-procesamiento para corregir errores comunes
func executeARM64Assembly(arm64Code string) (string, bool, string) {
	// PRE-PROCESAR EL CÓDIGO PARA CORREGIR ERRORES CONOCIDOS
	arm64Code = fixARM64Code(arm64Code)

	// Crear archivo temporal
	tmpDir := "/tmp"
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())
	sourceFile := filepath.Join(tmpDir, fmt.Sprintf("temp_program_%s.s", timestamp))
	objectFile := filepath.Join(tmpDir, fmt.Sprintf("temp_program_%s.o", timestamp))
	executableFile := filepath.Join(tmpDir, fmt.Sprintf("temp_program_%s", timestamp))

	// Escribir código ARM64 corregido al archivo
	err := ioutil.WriteFile(sourceFile, []byte(arm64Code), 0644)
	if err != nil {
		return "", false, fmt.Sprintf("Error creating source file: %v", err)
	}
	defer os.Remove(sourceFile)

	// Verificar que las herramientas existen
	if _, err := exec.LookPath("aarch64-linux-gnu-as"); err != nil {
		return "", false, "aarch64-linux-gnu-as not found. Install with: sudo apt-get install gcc-aarch64-linux-gnu"
	}

	// Compilar con aarch64-linux-gnu-as
	asmCmd := exec.Command("aarch64-linux-gnu-as", "-o", objectFile, sourceFile)
	asmOutput, err := asmCmd.CombinedOutput()
	if err != nil {
		return "", false, fmt.Sprintf("Assembly error: %s", string(asmOutput))
	}
	defer os.Remove(objectFile)

	// Enlazar con aarch64-linux-gnu-ld
	ldCmd := exec.Command("aarch64-linux-gnu-ld", "-o", executableFile, objectFile)
	ldOutput, err := ldCmd.CombinedOutput()
	if err != nil {
		return "", false, fmt.Sprintf("Linker error: %s", string(ldOutput))
	}
	defer os.Remove(executableFile)

	// Verificar que qemu existe
	if _, err := exec.LookPath("qemu-aarch64"); err != nil {
		return "", false, "qemu-aarch64 not found. Install with: sudo apt-get install qemu-user"
	}

	// Ejecutar con qemu-aarch64 con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	execCmd := exec.CommandContext(ctx, "qemu-aarch64", "-L", "/usr/aarch64-linux-gnu", executableFile)
	execOutput, err := execCmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		return "", false, "Execution timeout (10 seconds)"
	}

	if err != nil {
		errorMsg := fmt.Sprintf("Execution error: %v", err)
		if len(execOutput) > 0 {
			errorMsg += fmt.Sprintf("\nOutput: %s", string(execOutput))
		}
		return string(execOutput), false, errorMsg
	}

	return string(execOutput), true, ""
}

// Función para corregir errores comunes en el código ARM64 generado
func fixARM64Code(code string) string {
	// 1. Corregir la llamada duplicada a print_string al final de interpolación
	// Buscar el patrón problemático: "// Interpolación completada\n    // Llamar función print_string\n    bl print_string"
	interpolationFix := regexp.MustCompile(`// Interpolación completada\s*\n\s*// Llamar función print_string\s*\n\s*bl print_string`)
	code = interpolationFix.ReplaceAllString(code, "// Interpolación completada")

	// 2. Corregir funciones print_integer que usan registros no preservados
	// Buscar y reemplazar la función print_integer problemática
	printIntegerRegex := regexp.MustCompile(`print_integer:\s*\n(.*?\n)*?\s*ret`)
	code = printIntegerRegex.ReplaceAllString(code, `print_integer:
    // Función mejorada para imprimir enteros
    stp x29, x30, [sp, #-16]!    // Guardar frame pointer y link register
    mov x29, sp                   // Setup frame pointer
    
    // Manejar caso especial: cero
    cmp x0, #0
    bne .L_not_zero
    
    // Imprimir '0'
    mov x0, #48                   // ASCII '0'
    bl print_char
    b .L_print_int_done
    
.L_not_zero:
    // Usar una implementación más simple y robusta
    // Solo manejamos números positivos pequeños por ahora
    cmp x0, #10
    blt .L_single_digit
    
    // Para números >= 10, imprimir recursivamente
    mov x1, x0
    mov x2, #10
    udiv x0, x1, x2               // x0 = x1 / 10
    bl print_integer             // Llamada recursiva
    
    mov x2, #10
    msub x0, x0, x2, x1          // x0 = x1 % 10
    
.L_single_digit:
    add x0, x0, #48              // Convertir a ASCII
    bl print_char
    
.L_print_int_done:
    ldp x29, x30, [sp], #16      // Restaurar registros
    ret`)

	return code
}

func main() {
	r := mux.NewRouter()

	// API Routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/status", healthCheck).Methods("GET")
	api.HandleFunc("/execute", executeCode).Methods("POST")

	// NUEVA RUTA PARA ARM64
	api.HandleFunc("/execute-arm64", executeARM64Code).Methods("POST")

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(r)

	port := ":8080"
	fmt.Printf("🚀 Servidor Go iniciado en http://localhost%s\n", port)
	fmt.Println("📋 API endpoints disponibles:")
	fmt.Println("  - GET    /api/status")
	fmt.Println("  - POST   /api/execute")

	log.Fatal(http.ListenAndServe(port, handler))
}
