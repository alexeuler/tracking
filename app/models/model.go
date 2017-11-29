package models

//An interface for the model, required to decouple controller from model, returned by constructor
//This interface concept is immature and requires further detalization if the program grows
type Model interface {
	Save() bool
}
