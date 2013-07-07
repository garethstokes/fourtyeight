passwords
=========

standard functions needed for user auth

```go
	// pass in the plaintext password with the number
	// of iterations you want the hashing algo to pass through
	password := passwords.ComputeWithIterations("a secret password", 3)

	fmt.Println("hash: ", password.Hash)
	fmt.Println("salt: ", password.Salt)

	result := passwords.ComputeWithSalt("a secret password", 3, password.Salt)
	if result.Hash != password.Hash {
		fmt.Println("ERROR: incorrect password")
	}

  // of course, you could always generate a password without
  // specifying a interation. the library will generate one
  // for you
  password := passwords.Compute("a secret password")
```
