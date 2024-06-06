package main

func main() {
	call(nil)

}

func call(asd func()) {
	asd()
}
