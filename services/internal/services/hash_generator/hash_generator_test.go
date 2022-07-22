package hashgenerator

import (
	"testing"
)

const (
	pass string = "test_pass"
	salt string = "test_salt"
)

func TestHashGenerator(t *testing.T) {
	hashGen := HashGenerator{
		salt: salt,
	}
	salted_pass := hashGen.Salt(pass)
	if salted_pass != (pass + salt) {
		t.Errorf("Err: некорректная соль: %v", salted_pass)
	}
	hashed_pass, err := hashGen.Hash(pass)
	if err != nil {
		t.Fatal(err)
	}
	res := hashGen.Check(pass, hashed_pass)
	if !res {
		t.Error("Err: ошибка проверки валидного пароля")
	}
	res = hashGen.Check("incorrect", hashed_pass)
	if res {
		t.Error("Err: ошибка проверки невалидного пароля")
	}
}
