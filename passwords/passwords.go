package passwords

import (
	"math/rand"
	"time"
	"strings"
	"crypto/md5"
	"io"
	"encoding/base32"
	"bytes"
)

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890 abcdefghijklmnopqrstuvwxyz" +
	"~!@#$%^&*()-_+={}[]\\|<,>.?/\"';:`"

type Password struct {
	Hash string
	Salt string
	Iterations int
}

func toBase32(input []byte) (string) {
	output := &bytes.Buffer{}

	encoder := base32.NewEncoder(base32.StdEncoding, output)
	encoder.Write(input)
	encoder.Close()

	return output.String()[:25]
}

func fromBase32(input string) ([]byte) {
	output := &bytes.Buffer{}

	decoder := base32.NewDecoder(base32.StdEncoding, output)
	decoder.Read([]byte(input))

	return output.Bytes()
}

func seedRandom() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

func ComputeWithSalt(value string, iterations int, salt string) (Password) {
	h := md5.New()

	if iterations < 1 {
		iterations = 1
	}

	for i := 0; i < iterations; i++ {
		io.WriteString(h, value + salt)
		value = string(h.Sum(nil))
	}

	password := Password{
		Hash: toBase32([]byte(value)),
		Salt: salt,
		Iterations: iterations,
	}
	return password;
}

func Compute(value string) (Password) {
	salt := CreateRandomSalt()
  iterations := rand.Intn(10)

	return ComputeWithSalt(value, iterations, salt)
}

func ComputeWithIteration(value string, iterations int) (Password) {
	salt := CreateRandomSalt()
	return ComputeWithSalt(value, iterations, salt)
}

func CreateRandomSalt() (string) {
	seedRandom()
	length := 8

	r := make([]string, length)
	ri := 0
	buf := make([]byte, length)
	known := map[string]bool{}

	for i := 0; i < length; i++ {
	retry:
		l := rand.Intn(length)
		for j := 0; j < l; j++ {
			buf[j] = chars[rand.Intn(len(chars))]
		}
		s := string(buf[0:l])
		if known[s] {
			goto retry
		}
		known[s] = true
		r[ri] = s
		ri++
	}
	return strings.Join(r, "")
}
