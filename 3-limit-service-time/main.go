//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync"
	"time"
)

// User definisce il modello utente. Usa questo per controllare se un utente è premium o no
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in secondi
}

var userMutex sync.Mutex

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	if u.IsPremium {
		process()
		return true
	}

	userMutex.Lock()
	remainingTime := int64(10) - u.TimeUsed
	if remainingTime <= 0 {
		userMutex.Unlock()
		return false
	}
	userMutex.Unlock()

	done := make(chan bool)
	go func() {
		process()
		done <- true
	}()

	ticker := time.NewTicker(time.Second) // genera eventi ogni secondo e li invia sul suo canale
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return true
		case <-ticker.C: // il canale del ticker riceve un valore ogni secondo
			userMutex.Lock()
			u.TimeUsed++    // quindi ogni secondo viene aggiornato il tempo usato dall'utente
			remainingTime-- // e diminuito il tempo restante
			if remainingTime <= 0 {
				userMutex.Unlock()
				return false
			}
			userMutex.Unlock()
		case <-time.After(10 * time.Second):
			return false
		}
	}
}

func main() {
	RunMockServer()
}
