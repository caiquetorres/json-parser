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
		return parseObj(ps)
	case LeftBrace:
		return parseArr(ps)
	}
	return errUnexpectedTok
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
		err = parseKeyValue(ps)
		if err != nil {
			return err
		}
		tok, err = ps.peek()
		if err != nil {
			return err
		}
		if tok.k == Comma {
			expectKeyValue = true
			ps.next()
		} else if tok.k == RightBracket {
			expectKeyValue = false
		} else {
			return errUnexpectedTok
		}
	}
	if _, err := ps.expect(RightBracket); err != nil {
		return err
	}
	return nil
}

func parseKeyValue(ps *parseStream) error {
	_, err := ps.expect(String)
	if err != nil {
		return err
	}
	_, err = ps.expect(Colon)
	if err != nil {
		return err
	}
	err = parseExpr(ps)
	if err != nil {
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
		err = parseExpr(ps)
		if err != nil {
			return err
		}
		tok, err = ps.peek()
		if err != nil {
			return err
		}
		if tok.k == Comma {
			expectKeyValue = true
			ps.next()
		} else if tok.k == RightBrace {
			expectKeyValue = false
		} else {
			return errUnexpectedTok
		}
	}
	if _, err := ps.expect(RightBrace); err != nil {
		return err
	}
	return nil
}
