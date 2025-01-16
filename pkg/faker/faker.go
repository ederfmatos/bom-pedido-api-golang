package faker

import "github.com/go-faker/faker/v4"

func Name() string {
	return faker.Name()
}

func Word() string {
	return faker.Word()
}

func Email() string {
	return faker.Email()
}

func DomainName() string {
	return faker.DomainName()
}

func PhoneNumber() string {
	return faker.Phonenumber()
}

func Jwt() string {
	return faker.Jwt()
}

func URL() string {
	return faker.URL()
}
