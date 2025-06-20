# ¿Qué es "codegen"?

"codegen" es el nombre corto de "Code Generation" (Generación de Código). Es una carpeta nueva que vas a crear dentro de tu carpeta backend/ existente.
¿Por qué una carpeta separada?

Organización: Separar la lógica de generación de código del intérprete existente
Modularidad: Tu intérprete actual seguirá funcionando, el compilador será adicional
Mantenimiento: Código más limpio y fácil de mantener


## Ejemplo de estructura

```
tu-proyecto/
├── backend/                    # Tu carpeta backend existente
│   ├── value/                 # Tus carpetas existentes
│   ├── repl/                  # Tus carpetas existentes  
│   ├── errors/                # Tus carpetas existentes
│   ├── grammar/               # Tus carpetas existentes
│   ├── ast/                   # Tus carpetas existentes
│   ├── cst/                   # Tus carpetas existentes
│   ├── main.go                # Tu archivo main existente
│   └── codegen/               # 🆕 NUEVA CARPETA QUE VAS A CREAR
│       ├── arm64/             # Específico para ARM64
│       │   ├── registers.go   # Gestión de registros
│       │   ├── instructions.go # Definición de instrucciones
│       │   ├── calling_conv.go # Convenciones de llamada
│       │   └── optimizer.go   # Optimizaciones básicas
│       ├── intermediate/      # Representación intermedia
│       │   ├── ir.go         # Definición del IR
│       │   ├── generator.go  # AST → IR
│       │   └── optimizer.go  # Optimizaciones del IR
│       └── output/           # Generación del archivo final
│           ├── assembler.go  # IR → ARM64 Assembly
│           └── linker.go     # Enlazado y archivo ejecutable
├── frontend/                  # Tu carpeta frontend existente
└── README.md                  # Tus archivos existentes
```
