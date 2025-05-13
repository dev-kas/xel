package helpers  
  
import (  
	"fmt"  
  
	"github.com/dev-kas/virtlang-go/v2/ast"  
	"github.com/dev-kas/virtlang-go/v2/environment"  
	"github.com/dev-kas/virtlang-go/v2/errors"  
	"github.com/dev-kas/virtlang-go/v2/evaluator"  
	"github.com/dev-kas/virtlang-go/v2/shared"  
	"github.com/dev-kas/virtlang-go/v2/values"  
)  
  
func EvalFnVal(fnValue *shared.RuntimeValue, args []shared.RuntimeValue, env *environment.Environment) (*shared.RuntimeValue, error) {  
	// Verify it's a function  
	if fnValue.Type != shared.Function {  
		return nil, fmt.Errorf("not a function: %s", shared.Stringify(fnValue.Type))  
	}  
  
	// Create argument expressions  
	argExprs := make([]ast.Expr, len(args))  
	  
	// Track which arguments need to be added to the environment  
	envArgs := make(map[int]bool)  
	  
	for i, arg := range args {  
		// Create appropriate expressions based on argument type  
		switch arg.Type {  
		case shared.Number:  
			argExprs[i] = &ast.NumericLiteral{Value: arg.Value.(float64)}  
		case shared.String:  
			argExprs[i] = &ast.StringLiteral{Value: arg.Value.(string)}  
		case shared.Boolean:  
			if arg.Value.(bool) {  
				argExprs[i] = &ast.Identifier{Symbol: "true"}  
			} else {  
				argExprs[i] = &ast.Identifier{Symbol: "false"}  
			}  
		default:  
			// Mark complex types for environment addition  
			envArgs[i] = true  
			argExprs[i] = &ast.Identifier{Symbol: fmt.Sprintf("_arg_%d", i)}  
		}  
	}  
  
	// Create a temporary environment with the function  
	tempEnv := environment.NewEnvironment(env)  
	  
	// Add the function to the environment  
	tempEnv.DeclareVar("_fn", *fnValue, true)  
	  
	// Only add complex arguments to the environment  
	for i, arg := range args {  
		if envArgs[i] {  
			tempEnv.DeclareVar(fmt.Sprintf("_arg_%d", i), arg, true)  
		}  
	}  
  
	// Create a call expression  
	callExpr := &ast.CallExpr{  
		Callee: &ast.Identifier{Symbol: "_fn"},  
		Args:   argExprs,  
	}  
  
	// Evaluate the call expression  
	ret, evalerr := evaluator.Evaluate(callExpr, &tempEnv)  
	if evalerr != nil {  
		if evalerr.InternalCommunicationProtocol != nil &&   
		   evalerr.InternalCommunicationProtocol.Type == errors.ICP_Return {  
			return evalerr.InternalCommunicationProtocol.RValue, nil  
		}  
		return nil, evalerr  
	}  
	return ret, nil  
}
  
func convertArgsToVirtLang(args []interface{}) ([]*shared.RuntimeValue, error) {  
	virtlangArgs := make([]*shared.RuntimeValue, len(args))  
  
	for i, arg := range args {  
		var value shared.RuntimeValue  
  
		switch v := arg.(type) {  
		case int:  
			value = values.MK_NUMBER(float64(v))  
		case float64:  
			value = values.MK_NUMBER(v)  
		case string:  
			value = values.MK_STRING(v)  
		case bool:  
			value = values.MK_BOOL(v)  
		case nil:  
			value = values.MK_NIL()  
		case []interface{}:  
			// Convert slice to VirtLang array  
			arrayElements, err := convertArgsToVirtLang(v)  
			if err != nil {  
				return nil, err  
			}  
  
			elements := make([]shared.RuntimeValue, len(arrayElements))  
			for j, elem := range arrayElements {  
				elements[j] = *elem  
			}  
  
			value = values.MK_ARRAY(elements)  
		case map[string]interface{}:  
			// Convert map to VirtLang object  
			objProperties := make(map[string]*shared.RuntimeValue)  
  
			for k, val := range v {  
				converted, err := convertSingleArgToVirtLang(val)  
				if err != nil {  
					return nil, err  
				}  
				objProperties[k] = converted  
			}  
  
			value = values.MK_OBJECT(objProperties)  
		case *shared.RuntimeValue:  
			// Already a VirtLang value (pointer)  
			virtlangArgs[i] = v  
			continue  
		case shared.RuntimeValue:  
			// Already a VirtLang value (struct)  
			virtlangArgs[i] = &v  
			continue  
		default:  
			return nil, fmt.Errorf("unsupported argument type: %T", arg)  
		}  
  
		virtlangArgs[i] = &value  
	}  
  
	return virtlangArgs, nil  
}  
  
func convertSingleArgToVirtLang(arg interface{}) (*shared.RuntimeValue, error) {  
	args := []interface{}{arg}  
	converted, err := convertArgsToVirtLang(args)  
	if err != nil {  
		return nil, err  
	}  
	return converted[0], nil  
}