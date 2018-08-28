package nocopy

// NoCopy is the type you should embed as a value (not as a pointer to it)
// in a type you need to make checkable for go vet so it can see that you
// should not copy the type anywhere
type NoCopy struct{}

// Lock is an empty method that shows go vet that developer should not
// copy the type where NoCopy was embedded
func (*NoCopy) Lock() {}
