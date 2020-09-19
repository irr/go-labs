package engines

import (
	"fmt"

	"github.com/gofrs/uuid"
)

// MyEngine ...
type MyEngine struct {
	Name string `json:"name"`
}

// GetMessage ...
func (me MyEngine) GetMessage() (string, error) {
	u4, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	fmt.Printf("MyEngine GetMessage called. [%+v]\n", u4)
	return u4.String(), nil
}

// Validate ...
func (me MyEngine) Validate(msg *string) bool {
	fmt.Printf("MyEngine Validate called. [%+v]\n", *msg)
	return true
}

// Process ...
func (me MyEngine) Process(msg *string) error {
	fmt.Printf("MyEngine Process called. [%+v]\n", me.Name)
	return nil
}

// MyEngine2 ...
type MyEngine2 struct {
	Name  string `json:"name"`
	Other string `json:"other"`
}

// GetMessage ...
func (me MyEngine2) GetMessage() (string, error) {
	u4, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	fmt.Printf("MyEngine2 GetMessage called. [%+v]\n", u4)
	return u4.String(), nil
}

// Validate ...
func (me MyEngine2) Validate(msg *string) bool {
	fmt.Printf("MyEngine2 Validate called. [%+v]\n", *msg)
	return true
}

// Process ...
func (me MyEngine2) Process(msg *string) error {
	fmt.Printf("MyEngine2 Process called. [%+v]\n", me.Name)
	return nil
}
