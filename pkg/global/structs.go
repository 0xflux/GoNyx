package global

/*
A place to store shared structs. Lots of areas within the project will share code & structs,
so makes sense to keep these in 1 place.
*/

/*
Below structs are to try help unpack the data from the proxy listeners..
*/

type SomeData struct {
	Info string
}

type BrowserConnection struct {
	Data *SomeData
}

type Listener struct {
	ConnectionAddress string
	OtherDetails      ListenerSubData
}

type ListenerSubData struct {
	UninitializedPointer *string
	IntegerValue1        int
	IntegerValue2        int
}
