package sign

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

func (s *Signer) Verify(r *http.Request, signature string, hf HmacFunc) error {
	signedHeaders := SignedHeaders(r)
	canonicalRequest, err := CanonicalRequest(r, signedHeaders)
	if err != nil {
		return err
	}

	calc, err := hf.Hash(canonicalRequest)
	if err != nil {
		return err
	}

	calcstr := strings.TrimSpace(base64.StdEncoding.EncodeToString(calc))
	calcstr = strings.Replace(calcstr, " ", "+", -1)
	calcstr = url.QueryEscape(calcstr)

	if calcstr != signature {
		return errors.New("signature not match")
	}

	return nil
}
