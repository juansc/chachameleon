package main

import "math/rand"

type SmallCategory struct {
	entries []string
}

func (s SmallCategory) Entries() []string {
	entries := make([]string, len(s.entries))
	copy(entries, s.entries)
	return entries
}

func (s SmallCategory) GetRandomEntry() string {
	ind := rand.Intn(len(s.entries))
	return s.entries[ind]
}

func CategoryFairyTales() SmallCategory {
	return SmallCategory{
		entries: []string{
			"Cinderella",
			"Goldilocks",
			"Jack and the Beanstalk",
			"Hare and the Tortoise",

			"Snow White",
			"Rapunzel",
			"Aladdin",
			"Princess and the Pea",

			"Peter Pan",
			"Little Red Riding Hood",
			"Pinocchio",
			"Beauty and the Beast",

			"Sleeping Beauty",
			"Hansel and Gretel",
			"Gingerbread Man",
			"Three Little Pigs",
		},
	}
}

func CategoryFood() SmallCategory {
	return SmallCategory{
		entries: []string{
			"Pizza",
			"Potatoes",
			"Fish",
			"Cake",

			"Pasta",
			"Salad",
			"Soup",
			"Bread",

			"Eggs",
			"Cheese",
			"Fruit",
			"Chicken",

			"Sausage",
			"Ice Cream",
			"Chocolate",
			"Beef",
		},
	}
}

func GetSmallCategory() SmallCategory {
	categories := []SmallCategory{
		CategoryFairyTales(),
		CategoryFood(),
	}
	ind := rand.Intn(len(categories))
	return categories[ind]
}
