package jsonparser

import "io"

func Parse(r io.Reader) error {
	ps := newParseStream(r)
	tok, err := ps.peek()
	if err != nil {
		return err
	}
	switch tok.k {
	case LeftBracket:
		err := parseObj(ps)
		if err != nil {
			return err
		}
	case LeftBrace:
		err := parseArr(ps)
		if err != nil {
			return err
		}
	default:
		return errUnexpectedTok
	}
	_, err = ps.peek()
	if err == nil {
		return errUnexpectedTok
	}
	return nil
}

func parseExpr(ps *parseStream) error {
	tok, err := ps.peek()
	if err != nil {
		return err
	}
	switch tok.k {
	case LeftBracket:
		return parseObj(ps)
	case LeftBrace:
		return parseArr(ps)
	}
	ps.next()
	return nil
}

func parseObj(ps *parseStream) error {
	if _, err := ps.expect(LeftBracket); err != nil {
		return err
	}
	expectKeyValue := false
	for {
		tok, err := ps.peek()
		if err != nil {
			return err
		}
		if tok.k == RightBracket {
			if expectKeyValue {
				return errUnexpectedTok
			}
			break
		}
		if err := parseKeyValue(ps); err != nil {
			return err
		}
		tok, err = ps.peek()
		if err != nil {
			return err
		}
		switch tok.k {
		case Comma:
			expectKeyValue = true
			ps.next()
		case RightBracket:
			expectKeyValue = false
		default:
			return errUnexpectedTok
		}
	}
	if _, err := ps.expect(RightBracket); err != nil {
		return err
	}
	return nil
}

func parseKeyValue(ps *parseStream) error {
	if _, err := ps.expect(String); err != nil {
		return err
	}
	if _, err := ps.expect(Colon); err != nil {
		return err
	}
	if err := parseExpr(ps); err != nil {
		return err
	}
	return nil
}

func parseArr(ps *parseStream) error {
	if _, err := ps.expect(LeftBrace); err != nil {
		return err
	}
	expectKeyValue := false
	for {
		tok, err := ps.peek()
		if err != nil {
			return err
		}
		if tok.k == RightBrace {
			if expectKeyValue {
				return errUnexpectedTok
			}
			break
		}
		if err := parseExpr(ps); err != nil {
			return err
		}
		tok, err = ps.peek()
		if err != nil {
			return err
		}
		switch tok.k {
		case Comma:
			expectKeyValue = true
			ps.next()
		case RightBrace:
			expectKeyValue = false
		default:
			return errUnexpectedTok
		}
	}
	if _, err := ps.expect(RightBrace); err != nil {
		return err
	}
	return nil
}
