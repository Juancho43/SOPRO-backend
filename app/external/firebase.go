package external

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

var (
	appInstance  *firebase.App
	appErr       error
	firebaseOnce sync.Once
)

func GetFirebaseApp() (*firebase.App, error) {
	firebaseOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var err error
		appInstance, err = firebase.NewApp(ctx, nil)
		if err != nil {
			appErr = err
			log.Printf("Error al inicializar Firebase: %v", err)
		}
		fmt.Println("Firebase ok")
	})

	return appInstance, appErr
}

func GetFirebaseAuthClient() (*auth.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := appInstance.Auth(ctx)
	if err != nil {
		log.Fatalf("Error obteniendo el cliente de Auth: %v\n", err)
		return nil, err
	}
	return client, nil
}
