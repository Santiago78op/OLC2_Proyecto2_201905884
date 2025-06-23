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

// función para traducir a ARM64
func translateToARM64(tree antlr.ParseTree) (string, []string, bool) {
	fmt.Printf("🔹 Iniciando traducción a ARM64...\n")

	// Crear el traductor
	translator := compiler.NewARM64Translator()

	// Traducir el programa
	arm64Code, errors := translator.TranslateProgram(tree)

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

func main() {
	r := mux.NewRouter()

	// API Routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/status", healthCheck).Methods("GET")
	api.HandleFunc("/execute", executeCode).Methods("POST")

	// NUEVA RUTA PARA ARM64
	api.HandleFunc("/arm64", getARM64Code).Methods("POST")

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
