package common

type EndPoint struct {
	Method 		string
	Parameters 	map[string]string
}

func(e EndPoint) GetMethod() string {
	return e.Method
}

func(e EndPoint) GetParameters() map[string]string {
	return e.Parameters
}