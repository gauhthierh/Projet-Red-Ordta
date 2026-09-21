package monde

import (
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/window"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func ConfigurerFenetreSimulation(application *app.Application, cameraSimulation *camera.Camera) {
	fenetre, conversionReussie := window.Get().(*window.GlfwWindow)

	if !conversionReussie {
		panic("la fenêtre utilisée n'est pas une fenêtre GLFW")
	}

	fenetre.SetFullscreen(true)

	largeur, hauteur := fenetre.GetFramebufferSize()

	rapportEcran := float32(largeur) / float32(hauteur)
	cameraSimulation.SetAspect(rapportEcran)

	application.Gls().Viewport(0, 0, int32(largeur), int32(hauteur))
}

func VerrouillerSourisSimulation() {
	fenetre, conversionReussie := window.Get().(*window.GlfwWindow)

	if !conversionReussie {
		panic("impossible de récupérer la fenêtre GLFW")
	}

	fenetre.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
}

func LibererSourisSimulation() {
	fenetre, conversionReussie := window.Get().(*window.GlfwWindow)

	if !conversionReussie {
		panic("impossible de récupérer la fenêtre GLFW")
	}

	fenetre.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
}
