package controllers

type CompareEvent struct {
	Name         string `json:"name"`
	Id           string `json:"uuid,omitempty"`
	Confidence   string `json:"confidence"`
}
type TestEvent struct {
	Name         string `json:"name"`
	Msg string  `json:"msg"`
}
func NewCompareEvent(Id string, confidence string) CompareEvent {
	return CompareEvent{Name: "compareResult", Id: Id, Confidence: confidence}
}
func NewTestEvent(Id string, confidence string) CompareEvent {
	return CompareEvent{Name: "user.login", Id: Id, Confidence: confidence}
}
