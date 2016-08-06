package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func update(c *cli.Context) error {
	failed := color.New(color.FgRed).SprintFunc()
	success := color.New(color.FgGreen).SprintFunc()
	warn := color.New(color.FgYellow).SprintFunc()
	var addons []string
	if len(c.Args()) == 0 {
		addons = make([]string, len(config.Addons))
		i := 0
		for k := range config.Addons {
			addons[i] = k
			i++
		}
	} else {
		addons = c.Args()
	}

	updated := 0
	wowV := wowVersion()

	for _, name := range addons {
		addon, ok := config.Addons[name]
		if !ok {
			fmt.Printf("%s: isn't installed\n", failed(name))
			continue
		}
		meta, err := downloadURL(name, addon.Source)
		if err != nil {
			fmt.Printf("%s: failed to retrieve metadata: %s\n", failed(name), err.Error())
			continue
		}
		if meta.Version == addon.Version {
			if wowV != 0 && addon.Interface != 0 && addon.Interface < wowV {
				fmt.Printf("%s: (out of date) no update from %s available\n", warn(name), addon.Version)
			} else {
				fmt.Printf("%s: up to date at version %s\n", success(name), addon.Version)
				continue
			}
		}
		err = installAddon(name, addon.Source, "updated")
		if err != nil {
			updated++
		}
	}
	fmt.Printf("%d addons updated\n", updated)
	if !config.KeepCache {
		purgeCache()
	}
	return nil
}

func checkupdate(c *cli.Context) error {
	failed := color.New(color.FgRed).SprintFunc()
	success := color.New(color.FgGreen).SprintFunc()
	updated := 0
	for name, addon := range config.Addons {
		meta, err := downloadURL(name, addon.Source)
		if err != nil {
			fmt.Printf("%s: failed to retrieve metadata: %s\n", failed(name), err.Error())
			continue
		}
		if meta.Version != addon.Version {
			fmt.Printf("%s: can be updated from %s to %s\n", success(name), addon.Version, meta.Version)
			updated++
		}
	}
	if updated > 0 {
		fmt.Printf("%d addons can be updated\n", updated)
	} else {
		fmt.Printf("%s\n", success("You have the latest version of everything"))
	}
	return nil
}