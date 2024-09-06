package main

func YourName(name string, aliases ...string) {

}

func CallYourName() {
	YourName("Neo")
	YourName("Neo", "Amy")
	YourName("Neo", "Amy", "Lydia")
	aliases := []string{"Amy", "Lydia"}
	YourName("Amy", aliases...)
}
