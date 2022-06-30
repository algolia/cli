package cmdutil

import (
	"sort"

	"github.com/spf13/pflag"
)

type CategoryFlagSet struct {
	Categories map[string]*pflag.FlagSet
	Print      *pflag.FlagSet
	Others     *pflag.FlagSet
}

func NewCategoryFlagSet(flags *pflag.FlagSet) *CategoryFlagSet {
	categories := make(map[string]*pflag.FlagSet)
	others := pflag.NewFlagSet("other", pflag.ContinueOnError)
	print := pflag.NewFlagSet("print", pflag.ContinueOnError)

	flags.VisitAll(func(f *pflag.Flag) {
		if _, ok := f.Annotations["Categories"]; ok {
			mainCategory := f.Annotations["Categories"][0]
			if _, ok := categories[mainCategory]; !ok {
				categories[mainCategory] = pflag.NewFlagSet(mainCategory, pflag.ContinueOnError)
			}
			categories[mainCategory].AddFlag(f)
		} else if _, ok := f.Annotations["IsPrint"]; ok {
			print.AddFlag(f)
		} else {
			others.AddFlag(f)
		}
	})

	return &CategoryFlagSet{
		Print:      print,
		Categories: categories,
		Others:     others,
	}
}

func (c *CategoryFlagSet) SortedCategoryNames() []string {
	var names []string
	for name := range c.Categories {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
