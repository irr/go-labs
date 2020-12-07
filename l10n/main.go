package main

import (
	"fmt"
	"time"

	"golang.org/x/text/currency"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"

	"github.com/goodsign/monday"
)

type entry struct {
	tag, key string
	msg      interface{}
}

var entries = [...]entry{
	{"en", "Hello World", "Hello World"},
	{"el", "Hello World", "Για Σου Κόσμε"},
	{"en", "%d task(s) remaining!", plural.Selectf(1, "%d",
		"=1", "One task remaining!",
		"=2", "Two tasks remaining!",
		"other", "[1]d tasks remaining!",
	)},
	{"el", "%d task(s) remaining!", plural.Selectf(1, "%d",
		"=1", "Μία εργασία έμεινε!",
		"=2", "Μια-δυο εργασίες έμειναν!",
		"other", "[1]d εργασίες έμειναν!",
	)},
}

func init() {
	message.SetString(language.Greek, "%s went to %s.", "%s πήγε στήν %s.")
	message.SetString(language.AmericanEnglish, "%s went to %s.", "%s is in %s.")
	message.SetString(language.Greek, "%s has been stolen.", "%s κλάπηκε.")
	message.SetString(language.AmericanEnglish, "%s has been stolen.", "%s has been stolen.")
	message.SetString(language.Greek, "How are you?", "Πώς είστε?.")

	message.Set(language.Greek, "You have %d. problem",
		plural.Selectf(1, "%d",
			"=1", "Έχεις ένα πρόβλημα",
			"=2", "Έχεις %[1]d πρόβληματα",
			"other", "Έχεις πολλά πρόβληματα",
		))
	message.Set(language.Greek, "You have %d days remaining",
		plural.Selectf(1, "%d",
			"one", "Έχεις μία μέρα ελεύθερη",
			"other", "Έχεις %[1]d μέρες ελεύθερες",
		))

	message.Set(language.Greek, "You are %d minute(s) late.",
		catalog.Var("minutes", plural.Selectf(1, "%d", "one", "λεπτό", "other", "λεπτά")),
		catalog.String("Αργήσατε %[1]d ${minutes}."))

	for _, e := range entries {
		tag := language.MustParse(e.tag)
		switch msg := e.msg.(type) {
		case string:
			message.SetString(tag, e.key, msg)
		case catalog.Message:
			message.Set(tag, e.key, msg)
		case []catalog.Message:
			message.Set(tag, e.key, msg...)
		}
	}
}

func main() {
	p := message.NewPrinter(language.BritishEnglish)
	p.Printf("There are %v flowers in our garden.\n", 1500)

	p = message.NewPrinter(language.Greek)
	p.Printf("There are %v flowers in our garden.\n", 1500)

	p = message.NewPrinter(language.Greek)
	p.Printf("%s went to %s.", "Ο Πέτρος", "Αγγλία")
	fmt.Println()

	p.Printf("%s has been stolen.", "Η πέτρα")
	fmt.Println()

	p = message.NewPrinter(language.AmericanEnglish)
	p.Printf("%s went to %s.", "Peter", "England")
	fmt.Println()

	p.Printf("%s has been stolen.", "The Gem")
	fmt.Println()

	p = message.NewPrinter(language.Greek)
	p.Printf("You have %d. problem", 1)
	fmt.Println()

	p.Printf("You have %d. problem", 2)
	fmt.Println()

	p.Printf("You have %d. problem", 5)
	fmt.Println()

	p.Printf("You have %d days remaining", 1)
	fmt.Println()

	p.Printf("You have %d days remaining", 10)
	fmt.Println()

	p = message.NewPrinter(language.Greek)
	p.Printf("You are %d minute(s) late.", 1) // prints Αργήσατε 1 λεπτό
	fmt.Println()
	p.Printf("You are %d minute(s) late.", 10) // prints Αργήσατε 10 λεπτά
	fmt.Println()

	p = message.NewPrinter(language.BrazilianPortuguese)
	p.Printf("BRA: %d", currency.Symbol(currency.USD.Amount(0.1)))
	fmt.Println()

	p.Printf("BRA: %d", currency.NarrowSymbol(currency.JPY.Amount(1.6)))
	fmt.Println()
	p.Printf("BRA: %d", currency.ISO.Kind(currency.Cash)(currency.EUR.Amount(12.255)))
	fmt.Println()

	p = message.NewPrinter(language.AmericanEnglish)
	p.Printf("USA: %d", currency.Symbol(currency.USD.Amount(0.1)))
	fmt.Println()
	p.Printf("USA: %d", currency.NarrowSymbol(currency.JPY.Amount(1.6)))
	fmt.Println()
	p.Printf("USA: %d", currency.ISO.Kind(currency.Cash)(currency.EUR.Amount(12.255)))
	fmt.Println()

	// https://golang.org/src/time/format.go

	fmt.Printf("BRA: %s", monday.Format(time.Now(), "02-January-2006", "pt_BR"))
	fmt.Println()

	fmt.Printf("USA: %s", monday.Format(time.Now(), "02-January-2006", "en_US"))
	fmt.Println()

	p = message.NewPrinter(language.Greek)
	p.Printf("Hello World")
	p.Println()
	p.Printf("%d task(s) remaining!", 2)
	p.Println()

	p = message.NewPrinter(language.English)
	p.Printf("Hello World")
	p.Println()
	p.Printf("%d task(s) remaining!", 2)
}
