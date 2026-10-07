package bus

import (
	"fmt"
	"sync"
)

// InMemoryEventBus es la implementación de alta velocidad en memoria.
type InMemoryEventBus struct {
	handlers map[string][]EventHandler
	lock     sync.RWMutex // Protege tus datos. El éxito exige orden y rechaza las condiciones de carrera.[cite: 20]
}

var (
	// instance almacena el único centro de mando de tu aplicación.
	instance *InMemoryEventBus

	// once garantiza que la inicialización ocurra exactamente una vez.
	once sync.Once
)

// GetInstance reemplaza a NewInMemoryEventBus y actúa como el único punto de acceso.
func GetInstance() *InMemoryEventBus {
	// once.Do ejecuta la función interna solo la primera vez que se llama.
	once.Do(func() {
		instance = &InMemoryEventBus{
			handlers: make(map[string][]EventHandler), // Inicializamos el mapa exactamente como antes.[cite: 20]
		}
	})

	return instance
}

// Subscribe inscribe a un módulo para que escuche una victoria específica.
func (b *InMemoryEventBus) Subscribe(topic string, handler EventHandler) {
	b.lock.Lock() // Bloqueo de escritura: garantizamos que el registro sea seguro.[cite: 20]
	defer b.lock.Unlock()

	b.handlers[topic] = append(b.handlers[topic], handler)
}

// Publish emite la orden de ejecución.
func (b *InMemoryEventBus) Publish(topic string, payload interface{}) {
	b.lock.RLock() // Bloqueo de lectura: permitimos que múltiples eventos se lean a la vez a máxima velocidad.[cite: 20]
	defer b.lock.RUnlock()

	if handlers, found := b.handlers[topic]; found { //[cite: 20]
		for _, handler := range handlers { //[cite: 20]
			// Al lanzar el handler en una nueva goroutine, desacoplamos el tiempo de ejecución.[cite: 20]
			go handler(payload) //[cite: 20]
			fmt.Print("Topic", topic, "Payload", payload)
		}
	}
}
