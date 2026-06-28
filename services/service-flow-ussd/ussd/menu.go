package ussd

// Menu is a single screen in a USSD flow.
//
// USSD menus are line-based: each line is one option.
//
// TODO: replace with the real renderer + i18n lookup.
type Menu struct {
	Title   string
	Options []string
}

// Render returns the USSD wire format for the menu.
//
// TODO: implement — typically: "<title>\n1. ...\n2. ...\n0. Back".
func (m Menu) Render() string {
	return m.Title
}

// Parse picks the option the user selected.
//
// TODO: implement — validate that choice is in range, return ErrInvalidChoice
// otherwise.
func (m Menu) Parse(input string) (int, error) {
	_ = input
	return 0, nil
}
