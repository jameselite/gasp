package commands

import (
	"errors"
	"fmt"
	"github.com/jameselite/gasp/helper"
	"github.com/jameselite/gasp/start"
	"os"
	"strings"
)

func MakeController(router string, path string, controllerName string, method string) (string, error) {

	project, tomlErr := start.ParseTOML()

	if tomlErr != nil {
		return "", errors.New(tomlErr.Error())
	}

	switch project.Architecture {
	case "layered":

		findRouter := helper.IsFileExists("internal/routers/" + router + ".go")
		if !findRouter {
			return "", errors.New("router file does not exist")
		}

		controller, makeErr := os.Create("internal/controllers/" + controllerName + ".go")

		if makeErr != nil {
			return "", errors.New(makeErr.Error())
		}

		defer controller.Close()
		
		var controllerContent string

		if project.Framework == "gin" {
			controllerContent = fmt.Sprintf(ControllerTemplateGin, helper.CapitalizeFirstLetter(controllerName))
		}

		if project.Framework == "fiber" {
			controllerContent = fmt.Sprintf(ControllerTemplateFiber, helper.CapitalizeFirstLetter(controllerName))
		}

		_, writeErr := controller.WriteString(controllerContent)

		if writeErr != nil {
			return "", errors.New(writeErr.Error())
		}

		rawRouterContent, readErr := os.ReadFile("internal/routers/" + router +".go")

		if readErr != nil {
			return "", errors.New(readErr.Error())
		}

		splitedContent := strings.Split(string(rawRouterContent), "\n")

		var updatedContent []string

		for _, line := range splitedContent {
			
			if strings.Contains(line, "import (") {
				updatedContent = append(updatedContent, line)
				updatedContent = append(updatedContent, "")
	
				switch project.Architecture {
				case "layered":
					importString := fmt.Sprintf(`"%s/internal/controllers"`, project.Projectname)
					updatedContent = append(updatedContent, "	"+importString)
	
				case "clean":

					importString := fmt.Sprintf(`"%s/internal/controllers"`, project.Projectname)
					updatedContent = append(updatedContent, "	"+importString)
				}

				continue
			}

			if strings.Contains(line, "app.Group") {
				updatedContent = append(updatedContent, line)
				updatedContent = append(updatedContent, "")

				if project.Framework == "gin" {
					method = strings.ToUpper(method)
				}

				if project.Framework == "fiber" {
					method = helper.CapitalizeFirstLetter(method)
				}

				controllerString := fmt.Sprintf(`router.%s("%s", %s)`, method, path, "controllers." + controllerName)

				updatedContent = append(updatedContent, "	" + controllerString)

				continue
			}

			updatedContent = append(updatedContent, line)
		} 

		finalContent := strings.Join(updatedContent, "\n")

		err := os.WriteFile("internal/routers/" + router + ".go", []byte(finalContent), 0644)
		if err != nil {
			return "", errors.New(err.Error())
		}

	case "clean":

		findRouter := helper.IsFileExists("delivery/routers/" + router + ".go")
		if !findRouter {
			return "", errors.New("router file does not exist")
		}

		controller, makeErr := os.Create("delivery/handlers/" + controllerName + ".go")

		if makeErr != nil {
			return "", errors.New(makeErr.Error())
		}

		defer controller.Close()

		var controllerContent string

		if project.Framework == "gin" {
			controllerContent = fmt.Sprintf(ControllerTemplateGin, controllerName)
		}

		if project.Framework == "fiber" {
			controllerContent = fmt.Sprintf(ControllerTemplateFiber, controllerName)
		}
		
		_, writeErr := controller.WriteString(controllerContent)

		if writeErr != nil {
			return "", errors.New(writeErr.Error())
		}

		rawRouterContent, readErr := os.ReadFile("delivery/routers/" + router + ".go")

		if readErr != nil {
			return "", errors.New(readErr.Error())
		}

		splitedContent := strings.Split(string(rawRouterContent), "\n")

		var updatedContent []string

		for _, line := range splitedContent {
			
			if strings.Contains(line, "import (") {
				updatedContent = append(updatedContent, line)
				updatedContent = append(updatedContent, "")
	
				switch project.Architecture {
				case "layered":
					updatedContent = append(updatedContent, "		"+project.Projectname + "/" + "internal" + "/controllers")
	
				case "clean":
					updatedContent = append(updatedContent, "		" + project.Projectname + "/" + "delivery" + "/handlers")
				}
				continue
			}

			if strings.Contains(line, "app.Group") {
				updatedContent = append(updatedContent, line)
				updatedContent = append(updatedContent, "")

				controllerString := fmt.Sprintf(`router.%s("%s", %s)`, method, path, "handlers." + controllerName)

				updatedContent = append(updatedContent, "	" + controllerString)

				continue
			}

			updatedContent = append(updatedContent, line)

		}

		finalContent := strings.Join(updatedContent, "\n")

		err := os.WriteFile("delivery/routers/" + router + ".go", []byte(finalContent), 0644)
		if err != nil {
			return "", errors.New(err.Error())
		}

	default:
		return "", errors.New("sorry, your architecture is not supported")
	}

	return "Controller " + controllerName + "created and added to your router !", nil 
}