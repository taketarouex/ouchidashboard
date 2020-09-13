package enum

// Order desc or asc
type Order struct{ value string }

// Asc order
var Asc = Order{"ASC"}

// Desc order
var Desc = Order{"DESC"}
