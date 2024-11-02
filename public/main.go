package main

import (
	"fmt"
	"syscall/js"
)

type Base struct {
	Value int
	Bonus int
}

type Stat struct {
	Value    int
	Type     string
	Maitrise string
}

// Niveau
var niveau = int(10)
var maitrise = map[string]int{
	"neutre":     0,
	"qualifié":   0,
	"expert":     0,
	"maitre":     0,
	"légendaire": 0,
}

// Base Caractéristiques
var caractéristiques = map[string]Base{
	"intelligence": {Value: 20, Bonus: 5},
	"dexterité":    {Value: 19, Bonus: 4},
	"constitution": {Value: 16, Bonus: 3},
	"sagesse":      {Value: 14, Bonus: 2},
	"charisme":     {Value: 12, Bonus: 1},
	"force":        {Value: 10, Bonus: 0},
	"vitesse":      {Value: 12, Bonus: 0},
}

// Values of Armure
var armure = map[string]int{
	"CA":  1,
	"JDS": 1,
}

// Values of Competence
var statistiques = map[string]Stat{

	// Defense
	"Vitalité": {Value: 0, Type: "neutre", Maitrise: "neutre"},
	"CA":       {Value: 0, Type: "dexterité", Maitrise: "qualifié"},
	"Réflexes": {Value: 0, Type: "dexterité", Maitrise: "expert"},
	"Vigueur":  {Value: 0, Type: "constitution", Maitrise: "expert"},
	"Volonté":  {Value: 0, Type: "sagesse", Maitrise: "expert"},

	// Attaque
	"DD": {Value: 0, Type: "intelligence", Maitrise: "expert"},
	"JS": {Value: 0, Type: "intelligence", Maitrise: "expert"},
	"JA": {Value: 0, Type: "dexterité", Maitrise: "qualifié"},

	// Competences
	"Acrobaties":                         {Value: 0, Type: "dexterité", Maitrise: "qualifié"},
	"Arcanes":                            {Value: 0, Type: "intelligence", Maitrise: "maitre"},
	"Artisanat":                          {Value: 0, Type: "intelligence", Maitrise: "maitre"},
	"Athlétisme":                         {Value: 0, Type: "force", Maitrise: "qualifié"},
	"Connaissances (bombes alchimiques)": {Value: 0, Type: "intelligence", Maitrise: "qualifié"},
	"Connaissances (elfes)":              {Value: 0, Type: "intelligence", Maitrise: "qualifié"},
	"Connaissances (ingénieurie)":        {Value: 0, Type: "intelligence", Maitrise: "qualifié"},
	"Diplomatie":                         {Value: 0, Type: "charisme", Maitrise: "neutre"},
	"Discrétion":                         {Value: 0, Type: "dexterité", Maitrise: "qualifié"},
	"Duperie":                            {Value: 0, Type: "charisme", Maitrise: "neutre"},
	"Intimidation":                       {Value: 0, Type: "charisme", Maitrise: "neutre"},
	"Médecine":                           {Value: 0, Type: "sagesse", Maitrise: "qualifié"},
	"Nature":                             {Value: 0, Type: "sagesse", Maitrise: "qualifié"},
	"Occultisme":                         {Value: 0, Type: "intelligence", Maitrise: "qualifié"},
	"Perception":                         {Value: 0, Type: "sagesse", Maitrise: "qualifié"},
	"Religion":                           {Value: 0, Type: "sagesse", Maitrise: "qualifié"},
	"Représentation":                     {Value: 0, Type: "charisme", Maitrise: "qualifié"},
	"Société":                            {Value: 0, Type: "intelligence", Maitrise: "qualifié"},
	"Survie":                             {Value: 0, Type: "sagesse", Maitrise: "neutre"},
	"Vol":                                {Value: 0, Type: "dexterité", Maitrise: "neutre"},
}

// Calulate all
func calculateAll() {
	maitrise["neutre"] = 0
	maitrise["qualifié"] = niveau + 2
	maitrise["expert"] = niveau + 4
	maitrise["maitre"] = niveau + 6
	maitrise["légendaire"] = niveau + 8

	// Competences
	for key, comp := range statistiques {
		comp.Value = caractéristiques[comp.Type].Bonus + maitrise[comp.Maitrise]
		statistiques[key] = comp
	}

	// Add 10 flat
	statistiques["DD"] = Stat{
		Value:    statistiques["CA"].Value + 10,
		Type:     statistiques["DD"].Type,
		Maitrise: statistiques["DD"].Maitrise,
	}

	// Add 10 flat + bonus armure
	statistiques["CA"] = Stat{
		Value:    statistiques["CA"].Value + 10 + armure["CA"],
		Type:     statistiques["CA"].Type,
		Maitrise: statistiques["CA"].Maitrise,
	}

	// Add bonus armure
	for _, key := range []string{"Réflexes", "Vigueur", "Volonté"} {
		statistiques[key] = Stat{
			Value:    statistiques[key].Value + armure["JDS"],
			Type:     statistiques[key].Type,
			Maitrise: statistiques[key].Maitrise,
		}
	}

	// Value of Vitalité is totally different
	statistiques["Vitalité"] = Stat{
		Value:    6 + 6*niveau + caractéristiques["constitution"].Bonus*niveau + niveau,
		Type:     statistiques["Vitalité"].Type,
		Maitrise: statistiques["Vitalité"].Maitrise,
	}
}

// GetStatistique returns the value of a specific statistiques
func GetStatistique(this js.Value, p []js.Value) interface{} {
	name := p[0].String()
	if def, exists := statistiques[name]; exists {
		return js.ValueOf(def.Value)
	}
	return js.Null()
}

// main function
func main() {
	calculateAll()

	// Create a channel to keep the program running
	c := make(chan struct{}, 0)

	// Set the global variable "caractéristiques" to the caractéristiques map
	js.Global().Set("GetStatistique", js.FuncOf(GetStatistique))

	// Print the caractéristiques map to the console
	fmt.Println("WASM Go initialized")
	fmt.Println("Calculate Caractéristiques: ", caractéristiques)
	fmt.Println("Calculate statistiques: ", statistiques)

	// Wait for the channel to receive a message
	<-c
}
