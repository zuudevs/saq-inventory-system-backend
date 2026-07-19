package dto

type Response[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

func Success[T any](message string, data T) Response[T] {
	return Response[T]{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func Error[T any](message string) Response[T] {
	return Response[T]{
		Success: false,
		Message: message,
	}
}
