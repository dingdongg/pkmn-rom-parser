package tutil

// useful test logging templates 

type templates struct {
	Uint string
	UintHex string
	Int string
	IntHex string
	String string
}

func GetTemplates() templates {
	return templates{
		FatalfUint(false),
		FatalfUint(true),
		FatalfInt(false),
		FatalfInt(true),
		FatalfString(),
	}
}

func FatalfUint(inHex bool) string {
	s := "expected %d but got %d\n"
	if inHex {
		s = "expected 0x%x but got 0x%x\n"
	}

	return s
}

func FatalfInt(inHex bool) string {
	s := "expected %d but got %d\n"
	if inHex {
		s = "expected 0x%x but got 0x%x\n"
	}

	return s
}

func FatalfString() string {
	return "expected '%s' but got '%s'\n"
}

func LogSlice(hex bool) string {
	s := "[ %+v ]\n"
	if hex {
		s = "[ 0x%+x ]\n"
	}

	return s
}