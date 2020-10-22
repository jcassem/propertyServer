package property

// Property model
type Property struct {
	Name string  `json:"name"`
	Rent float32 `json:"rent"`
}

// Example slice
var properties = []Property{
	Property{
		Name: "123 Fake Street",
		Rent: 1200.00,
	},
	Property{
		Name: "2 Main Road",
		Rent: 899.50,
	},
	Property{
		Name: "Flat A 120 Regents Street",
		Rent: 14060.66,
	},
}

// ListProperties Lists all properties
func ListProperties() []Property {
	return properties
}
