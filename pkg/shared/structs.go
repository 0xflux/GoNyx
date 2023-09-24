package shared

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
