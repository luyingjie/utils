package json

func (j *Json) MarshalJSON() ([]byte, error) {
	return j.ToJson()
}

func (j *Json) UnmarshalJSON(b []byte) error {
	r, err := LoadContent(b)
	if r != nil {
		*j = *r
	}
	return err
}

func (j *Json) UnmarshalValue(value interface{}) error {
	if r := New(value); r != nil {
		*j = *r
	}
	return nil
}
