package logs

type TextFormatter struct {
	IgnoreBasicFields bool
}

func NewTextFormatter(ignoreBasicFields bool) *TextFormatter {
	return &TextFormatter{
		IgnoreBasicFields: ignoreBasicFields,
	}
}

func (t *TextFormatter) Format(entry *Entry) error {
	return nil
}