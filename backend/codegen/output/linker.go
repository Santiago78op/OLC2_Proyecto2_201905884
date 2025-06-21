// ========================================
// ARM64 LINKER - GENERACIÓN DE EJECUTABLES
// Archivo: backend/codegen/output/linker.go
// ========================================

package output

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ARM64Linker maneja la generación de ejecutables ARM64
type ARM64Linker struct {
	workingDir     string
	outputDir      string
	runtimeLibPath string
	debugMode      bool
	keepTempFiles  bool
}

// LinkingOptions opciones para el proceso de enlazado
type LinkingOptions struct {
	OutputName    string   // Nombre del ejecutable final
	EntryPoint    string   // Punto de entrada (default: main)
	Libraries     []string // Librerías adicionales a enlazar
	StaticLink    bool     // Enlazado estático vs dinámico
	OptimizeSize  bool     // Optimizar para tamaño
	DebugInfo     bool     // Incluir información de debug
	StripSymbols  bool     // Eliminar símbolos para reducir tamaño
	KeepTempFiles bool     // Mantener archivos temporales
}

// LinkingResult resultado del proceso de enlazado
type LinkingResult struct {
	Success        bool          `json:"success"`
	ExecutablePath string        `json:"executablePath"`
	FileSize       int64         `json:"fileSize"`
	LinkTime       time.Duration `json:"linkTime"`
	Warnings       []string      `json:"warnings"`
	TempFiles      []string      `json:"tempFiles"`
}

// NewARM64Linker crea un nuevo linker ARM64
func NewARM64Linker(workingDir string) *ARM64Linker {
	return &ARM64Linker{
		workingDir:     workingDir,
		outputDir:      filepath.Join(workingDir, "build"),
		runtimeLibPath: filepath.Join(workingDir, "runtime"),
		debugMode:      false,
		keepTempFiles:  false,
	}
}

// SetDebugMode configura el modo debug
func (linker *ARM64Linker) SetDebugMode(debug bool) {
	linker.debugMode = debug
}

// SetKeepTempFiles configura si mantener archivos temporales
func (linker *ARM64Linker) SetKeepTempFiles(keep bool) {
	linker.keepTempFiles = keep
}

// LinkExecutable convierte assembly ARM64 en ejecutable
func (linker *ARM64Linker) LinkExecutable(assembly string, options LinkingOptions) (*LinkingResult, error) {
	startTime := time.Now()
	result := &LinkingResult{
		Success:   false,
		Warnings:  make([]string, 0),
		TempFiles: make([]string, 0),
	}

	fmt.Printf("🔗 Iniciando proceso de enlazado ARM64...\n")

	// 1. Crear directorio de trabajo
	if err := linker.ensureOutputDirectory(); err != nil {
		return result, fmt.Errorf("error creando directorio de salida: %v", err)
	}

	// 2. Generar archivo assembly completo
	asmFilePath, err := linker.generateAssemblyFile(assembly, options)
	if err != nil {
		return result, fmt.Errorf("error generando archivo assembly: %v", err)
	}
	result.TempFiles = append(result.TempFiles, asmFilePath)

	// 3. Generar runtime si es necesario
	runtimePath, err := linker.generateRuntime(options)
	if err != nil {
		return result, fmt.Errorf("error generando runtime: %v", err)
	}
	if runtimePath != "" {
		result.TempFiles = append(result.TempFiles, runtimePath)
	}

	// 4. Ensamblar código objeto
	objFilePath, err := linker.assemble(asmFilePath, options)
	if err != nil {
		return result, fmt.Errorf("error ensamblando: %v", err)
	}
	result.TempFiles = append(result.TempFiles, objFilePath)

	// 5. Enlazar ejecutable final
	executablePath, err := linker.link(objFilePath, runtimePath, options)
	if err != nil {
		return result, fmt.Errorf("error enlazando: %v", err)
	}

	// 6. Obtener información del archivo final
	fileInfo, err := os.Stat(executablePath)
	if err != nil {
		return result, fmt.Errorf("error obteniendo info del ejecutable: %v", err)
	}

	// 7. Limpiar archivos temporales si se solicita
	if !linker.keepTempFiles && !options.KeepTempFiles {
		linker.cleanupTempFiles(result.TempFiles)
		result.TempFiles = nil
	}

	// Completar resultado
	result.Success = true
	result.ExecutablePath = executablePath
	result.FileSize = fileInfo.Size()
	result.LinkTime = time.Since(startTime)

	fmt.Printf("✅ Enlazado completado: %s (%d bytes) en %v\n",
		executablePath, result.FileSize, result.LinkTime)

	return result, nil
}

// ensureOutputDirectory crea el directorio de salida si no existe
func (linker *ARM64Linker) ensureOutputDirectory() error {
	if _, err := os.Stat(linker.outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(linker.outputDir, 0755); err != nil {
			return err
		}
		fmt.Printf("📁 Directorio de salida creado: %s\n", linker.outputDir)
	}
	return nil
}

// generateAssemblyFile genera el archivo .s completo
func (linker *ARM64Linker) generateAssemblyFile(assembly string, options LinkingOptions) (string, error) {
	// Crear assembly completo con runtime headers
	completeAssembly := linker.buildCompleteAssembly(assembly, options)

	// Generar nombre de archivo temporal
	asmFileName := fmt.Sprintf("%s_%d.s", options.OutputName, time.Now().Unix())
	asmFilePath := filepath.Join(linker.outputDir, asmFileName)

	// Escribir archivo
	if err := os.WriteFile(asmFilePath, []byte(completeAssembly), 0644); err != nil {
		return "", err
	}

	fmt.Printf("📝 Archivo assembly generado: %s\n", asmFilePath)
	return asmFilePath, nil
}

// buildCompleteAssembly construye el assembly completo con headers y runtime
func (linker *ARM64Linker) buildCompleteAssembly(userAssembly string, options LinkingOptions) string {
	var builder strings.Builder

	// Header del archivo
	builder.WriteString("// Generated by VLan Cherry Compiler\n")
	builder.WriteString("// Target: ARM64 (AArch64)\n")
	builder.WriteString(fmt.Sprintf("// Build time: %s\n", time.Now().Format(time.RFC3339)))
	builder.WriteString("\n")

	// Directivas de ensamblador
	builder.WriteString(".arch armv8-a\n")
	builder.WriteString(".file \"vlancherry_program.s\"\n")
	builder.WriteString("\n")

	// Declaraciones globales
	builder.WriteString("// Global symbols\n")
	builder.WriteString(".global _start\n")
	if options.EntryPoint != "" && options.EntryPoint != "main" {
		builder.WriteString(fmt.Sprintf(".global %s\n", options.EntryPoint))
	}
	builder.WriteString(".global main\n")
	builder.WriteString("\n")

	// Declaraciones externas del runtime
	builder.WriteString("// Runtime functions\n")
	builder.WriteString(".extern _vlc_println\n")
	builder.WriteString(".extern _vlc_print\n")
	builder.WriteString(".extern _vlc_readln\n")
	builder.WriteString(".extern _vlc_malloc\n")
	builder.WriteString(".extern _vlc_free\n")
	builder.WriteString(".extern _vlc_exit\n")
	builder.WriteString("\n")

	// Sección de texto
	builder.WriteString(".section .text\n")
	builder.WriteString("\n")

	// Punto de entrada del sistema
	builder.WriteString("_start:\n")
	builder.WriteString("    // System entry point\n")
	builder.WriteString("    bl main              // Call user main function\n")
	builder.WriteString("    mov x8, #93          // sys_exit system call\n")
	builder.WriteString("    mov x0, #0           // exit status 0\n")
	builder.WriteString("    svc #0               // invoke system call\n")
	builder.WriteString("\n")

	// Código del usuario
	builder.WriteString("// User program code\n")
	builder.WriteString(userAssembly)
	builder.WriteString("\n")

	// Sección de datos si es necesaria
	if strings.Contains(userAssembly, ".LC") {
		builder.WriteString("// String literals and constants\n")
		builder.WriteString(".section .rodata\n")
	}

	return builder.String()
}

// generateRuntime genera el runtime básico de VLan Cherry
func (linker *ARM64Linker) generateRuntime(options LinkingOptions) (string, error) {
	// Para simplificar, generar un runtime básico en C
	runtimeC := `
// VLan Cherry Runtime Library
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

// Print function
void _vlc_print(long value) {
    printf("%ld", value);
}

// Println function  
void _vlc_println(long value) {
    printf("%ld\n", value);
}

// Read line function
long _vlc_readln(void) {
    char buffer[256];
    if (fgets(buffer, sizeof(buffer), stdin) != NULL) {
        return atol(buffer);
    }
    return 0;
}

// Memory allocation
void* _vlc_malloc(long size) {
    return malloc(size);
}

// Memory deallocation
void _vlc_free(void* ptr) {
    free(ptr);
}

// Exit function
void _vlc_exit(long code) {
    exit(code);
}
`

	// Generar archivo del runtime
	runtimeFileName := fmt.Sprintf("runtime_%d.c", time.Now().Unix())
	runtimePath := filepath.Join(linker.outputDir, runtimeFileName)

	if err := os.WriteFile(runtimePath, []byte(runtimeC), 0644); err != nil {
		return "", err
	}

	fmt.Printf("📚 Runtime generado: %s\n", runtimePath)
	return runtimePath, nil
}

// assemble convierte el archivo .s a archivo objeto .o
func (linker *ARM64Linker) assemble(asmFilePath string, options LinkingOptions) (string, error) {
	objFilePath := strings.Replace(asmFilePath, ".s", ".o", 1)

	// Comando para ensamblar
	args := []string{
		"-c",              // Compile only
		asmFilePath,       // Input file
		"-o", objFilePath, // Output file
	}

	if options.DebugInfo {
		args = append(args, "-g") // Include debug info
	}

	cmd := exec.Command("as", args...)

	if linker.debugMode {
		fmt.Printf("🔧 Ejecutando: as %s\n", strings.Join(args, " "))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error ensamblando: %v\nOutput: %s", err, string(output))
	}

	fmt.Printf("🔧 Archivo objeto generado: %s\n", objFilePath)
	return objFilePath, nil
}

// link enlaza el archivo objeto con el runtime para crear el ejecutable
func (linker *ARM64Linker) link(objFilePath, runtimePath string, options LinkingOptions) (string, error) {
	executablePath := filepath.Join(linker.outputDir, options.OutputName)

	// Compilar runtime a objeto si existe
	var runtimeObjPath string
	if runtimePath != "" {
		runtimeObjPath = strings.Replace(runtimePath, ".c", ".o", 1)

		// Compilar runtime
		runtimeArgs := []string{
			"-c",
			runtimePath,
			"-o", runtimeObjPath,
		}

		if options.DebugInfo {
			runtimeArgs = append(runtimeArgs, "-g")
		}

		cmd := exec.Command("gcc", runtimeArgs...)
		if output, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("error compilando runtime: %v\nOutput: %s", err, string(output))
		}
	}

	// Comando para enlazar
	args := []string{
		objFilePath,          // User object file
		"-o", executablePath, // Output executable
	}

	if runtimeObjPath != "" {
		args = append(args, runtimeObjPath) // Runtime object file
	}

	if options.StaticLink {
		args = append(args, "-static")
	}

	if options.StripSymbols {
		args = append(args, "-s")
	}

	if options.OptimizeSize {
		args = append(args, "-Os")
	}

	// Agregar librerías adicionales
	for _, lib := range options.Libraries {
		args = append(args, "-l"+lib)
	}

	cmd := exec.Command("gcc", args...)

	if linker.debugMode {
		fmt.Printf("🔗 Ejecutando: gcc %s\n", strings.Join(args, " "))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error enlazando: %v\nOutput: %s", err, string(output))
	}

	// Hacer ejecutable
	if err := os.Chmod(executablePath, 0755); err != nil {
		return "", fmt.Errorf("error estableciendo permisos: %v", err)
	}

	fmt.Printf("🎯 Ejecutable creado: %s\n", executablePath)
	return executablePath, nil
}

// cleanupTempFiles elimina archivos temporales
func (linker *ARM64Linker) cleanupTempFiles(tempFiles []string) {
	for _, file := range tempFiles {
		if err := os.Remove(file); err != nil {
			fmt.Printf("⚠️ No se pudo eliminar archivo temporal: %s\n", file)
		}
	}
	fmt.Printf("🧹 Archivos temporales eliminados\n")
}

// ValidateEnvironment verifica que las herramientas necesarias estén disponibles
func (linker *ARM64Linker) ValidateEnvironment() error {
	// Verificar que 'as' (assembler) esté disponible
	if _, err := exec.LookPath("as"); err != nil {
		return fmt.Errorf("assembler 'as' no encontrado: %v", err)
	}

	// Verificar que 'gcc' esté disponible
	if _, err := exec.LookPath("gcc"); err != nil {
		return fmt.Errorf("gcc no encontrado: %v", err)
	}

	// Verificar arquitectura objetivo
	cmd := exec.Command("uname", "-m")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("no se pudo determinar la arquitectura: %v", err)
	}

	arch := strings.TrimSpace(string(output))
	if arch != "aarch64" && arch != "arm64" {
		return fmt.Errorf("arquitectura no compatible: %s (se requiere ARM64/AArch64)", arch)
	}

	fmt.Printf("✅ Entorno ARM64 validado exitosamente\n")
	return nil
}

// GetDefaultLinkingOptions retorna opciones predeterminadas
func GetDefaultLinkingOptions() LinkingOptions {
	return LinkingOptions{
		OutputName:    "vlancherry_program",
		EntryPoint:    "main",
		Libraries:     []string{"c"}, // libc estándar
		StaticLink:    false,
		OptimizeSize:  false,
		DebugInfo:     true,
		StripSymbols:  false,
		KeepTempFiles: false,
	}
}
