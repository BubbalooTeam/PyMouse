package pymouse

import (
	"pymouse/pymouse/client"
	"pymouse/pymouse/modules/afk"
	"pymouse/pymouse/modules/checkers"
	"pymouse/pymouse/modules/miscellaneous"
	"pymouse/pymouse/modules/pm_menu"
	"pymouse/pymouse/modules/sudoers"
	"sort"
	"strings"

	"sync"

	"github.com/sirupsen/logrus"
)

var (
	packageLoadersMutex sync.Mutex
	packageLoaders      = map[string]func(*client.BotStruct){
		"afk":           afk.LoadModule,
		"checkers":      checkers.LoadModule,
		"miscellaneous": miscellaneous.LoadModule,
		"pm_menu":       pm_menu.LoadModule,
		"sudoers":       sudoers.LoadModules,
	}
)

func Register(bS *client.BotStruct) {
	var wg sync.WaitGroup
	done := make(chan struct{}, len(packageLoaders))
	moduleNames := make([]string, 0, len(packageLoaders))

	sortedModules := make([]string, 0, len(packageLoaders))
	for module := range packageLoaders {
		sortedModules = append(sortedModules, module)
	}
	sort.Strings(sortedModules)

	for _, module := range sortedModules {
		wg.Add(1)
		go func(moduleName string, moduleLoader func(*client.BotStruct)) {
			defer wg.Done()
			packageLoadersMutex.Lock()
			defer packageLoadersMutex.Unlock()
			moduleLoader(bS)
			done <- struct{}{}
			moduleNames = append(moduleNames, moduleName)
		}(module, packageLoaders[module])
	}
	go func() {
		wg.Wait()
		close(done)
	}()

	for range done {
	}

	joinedModuleNames := strings.Join(moduleNames, ", ")

	logrus.Infof("Modules Loaded: %s", joinedModuleNames)
}
