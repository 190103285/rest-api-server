package main

func main() {
	a := App{}
	a.Initialize(
		"postgres",
		"password",
		"coffeeshop",
		"disable")
	//,
	//os.Getenv("disable")

	a.Run(":5433")
}
