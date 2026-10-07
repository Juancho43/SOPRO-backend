package bus

// EventHandler es el contrato de acción. Define qué debe hacer un módulo cuando reciba la señal.
type EventHandler func(payload interface{})

// EventBus es la mente maestra (Mastermind) que coordina los módulos sin acoplarlos.
type EventBus interface {
	Publish(topic string, payload interface{})
	Subscribe(topic string, handler EventHandler) // La pieza faltante de tu interfaz
}
