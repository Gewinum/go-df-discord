package server

import (
    "crypto/rand"
    "fmt"
    "time"
)

type CodeInformation struct {
    Code    string
    XUID    string
    Issued  time.Time
    Expires time.Time
}

type CodeStore interface {
    GetInformation(code string) (*CodeInformation, error)
    GetForXuid(xuid string) (*CodeInformation, error)
    Issue(xuid string) (*CodeInformation, error)
    Revoke(code string) error
}

type defaultCodeStore struct {
    codes map[string]*CodeInformation
}

func newDefaultCodeStore() CodeStore {
    return &defaultCodeStore{
        codes: make(map[string]*CodeInformation),
    }
}

func (s *defaultCodeStore) GetInformation(code string) (*CodeInformation, error) {
    info, exists := s.codes[code]
    if !exists {
        return nil, NewApplicationError(40400, "Code doesn't exist")
    }
    return info, nil
}

func (s *defaultCodeStore) GetForXuid(xuid string) (*CodeInformation, error) {
    for _, info := range s.codes {
        if info.XUID == xuid {
            return info, nil
        }
    }
    return nil, NewApplicationError(40400, "There is no code for this XUID")
}

func (s *defaultCodeStore) Issue(xuid string) (*CodeInformation, error) {
    existing, _ := s.GetForXuid(xuid)
    if existing != nil {
        return nil, NewApplicationError(40000, fmt.Sprintf("Code %s is already issued", existing.Code))
    }
    generatedCode := s.findFreeCode()
    s.codes[generatedCode] = &CodeInformation{
        Code:    generatedCode,
        XUID:    xuid,
        Issued:  time.Now(),
        Expires: time.Now().Add(15 * time.Minute),
    }
    return s.codes[xuid], nil
}

func (s *defaultCodeStore) Revoke(code string) error {
    info, _ := s.GetInformation(code)
    if info == nil {
        return NewApplicationError(40400, "Code doesn't exist")
    }
    delete(s.codes, code)
    return nil
}

func (s *defaultCodeStore) findFreeCode() string {
    for {
        generated, err := generateCode(6)
        if err != nil {
            panic(err)
        }
        info, _ := s.GetInformation(generated)
        if info == nil {
            return generated
        }
    }
}

func generateCode(length int) (string, error) {
    randomChars := "0123456789ABCDEFGHIKLMNOPQRSTVXYZ"
    buffer := make([]byte, length)
    _, err := rand.Read(buffer)
    if err != nil {
        return "", err
    }

    otpCharsLength := len(randomChars)
    for i := 0; i < length; i++ {
        buffer[i] = randomChars[int(buffer[i])%otpCharsLength]
    }

    return string(buffer), nil
}
