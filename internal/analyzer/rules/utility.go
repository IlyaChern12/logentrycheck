package rules

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// logFuncs contains base configs for supported logger packages.
var logFuncs = map[string]map[string]bool{
	"slog": {
		"Debug":        true,
		"Info":         true,
		"Warn":         true,
		"Error":        true,
		"DebugContext": true,
		"InfoContext":  true,
		"WarnContext":  true,
		"ErrorContext": true,
	},
	"zap": {
		"Debug":  true,
		"Info":   true,
		"Warn":   true,
		"Error":  true,
		"Fatal":  true,
		"Panic":  true,
		"Debugf": true,
		"Infof":  true,
		"Warnf":  true,
		"Errorf": true,
		"Fatalf": true,
		"Panicf": true,
	},
	"Sugar": {
		"Infof":  true,
		"Debugf": true,
		"Warnf":  true,
		"Errorf": true,
		"Fatalf": true,
		"Info":   true,
		"Debug":  true,
		"Warn":   true,
		"Error":  true,
		"Fatal":  true,
	},
}

var logMethods = map[string]bool{
	"Debug": true, "Info": true, "Warn": true, "Error": true,
	"Fatal": true, "Panic": true,
	"Debugf": true, "Infof": true, "Warnf": true, "Errorf": true,
	"Fatalf": true, "Panicf": true,
	"DebugContext": true, "InfoContext": true, "WarnContext": true, "ErrorContext": true,
}

// reporter abstracts diagnostic reporting for tests.
type reporter interface {
	Reportf(pos token.Pos, format string, args ...any)
}

// IsContextMethod checks if the method has ctx as the first argument.
func IsContextMethod(name string) bool {
	contextMethods := map[string]bool{
		"DebugContext": true,
		"InfoContext":  true,
		"WarnContext":  true,
		"ErrorContext": true,
	}

	return contextMethods[name]
}

// ExtractLogMessage extracts the log message string from a call expression.
func ExtractLogMessage(pass *analysis.Pass, call *ast.CallExpr) (string, token.Pos, bool) {
	arg := ExtractMessageArg(pass, call)
	if arg == nil {
		return "", 0, false
	}

	lit, ok := arg.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", 0, false
	}

	msg := lit.Value[1 : len(lit.Value)-1]
	return msg, lit.Pos(), true
}

// ExtractMessageArg returns the message argument expression from a log call.
func ExtractMessageArg(pass *analysis.Pass, call *ast.CallExpr) ast.Expr {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}

	if !IsLoggerCall(pass, call) {
		return nil
	}

	if !logMethods[sel.Sel.Name] {
		return nil
	}

	msgIndex := 0
	if IsContextMethod(sel.Sel.Name) {
		msgIndex = 1
	}

	if len(call.Args) <= msgIndex {
		return nil
	}

	return call.Args[msgIndex]
}

// ExtractPackageName extract name of pacakage.
func ExtractPackageName(expr ast.Expr) (string, bool) {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name, true
	case *ast.SelectorExpr:
		return e.Sel.Name, true
	}
	return "", false
}

// IsLoggerCall checks if call is a logger call using type info or package name.
func IsLoggerCall(pass *analysis.Pass, call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	if pass != nil && pass.TypesInfo != nil {
		if isLoggerCallByType(pass, sel) {
			return true
		}
	}

	// fallback to package name matching
	pkg, ok := ExtractPackageName(sel.X)
	if !ok {
		return false
	}

	_, ok = logFuncs[pkg]
	return ok
}

// isLoggerCallByType checks logger call using type information.
func isLoggerCallByType(pass *analysis.Pass, sel *ast.SelectorExpr) bool {
	typ := pass.TypesInfo.TypeOf(sel.X)
	if typ == nil {
		return false
	}

	if ptr, ok := typ.(*types.Pointer); ok {
		typ = ptr.Elem()
	}

	typStr := typ.String()
	return typStr == "go.uber.org/zap.SugaredLogger" ||
		typStr == "go.uber.org/zap.Logger" ||
		typStr == "log/slog.Logger"
}
