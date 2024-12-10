package main

import (
	"bytes"
	"fmt"
	"image/color"
	"os"
	"os/exec"
	"path"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Создаем новое приложение
	myApp := app.NewWithID("com.pnd.go-cutter")
	// Создаем окно
	myWindow := myApp.NewWindow("Go-cutter")
	myWindow.Resize(fyne.NewSize(640, 300))

	// Крутящийся кружок (анимированное состояние)
	loadingIndicator := widget.NewProgressBarInfinite()
	loadingIndicator.Hide()

	notificationTitle := canvas.NewText("", color.RGBA{R: 255, G: 0, B: 0, A: 255})
	notificationTitle.TextSize = 14
	notificationTitle.Hide()
	labelSelectSize := widget.NewLabel("Введите размер в мегабайтах:")
	labelSelectedFile := widget.NewLabel("Выбранный файл:")
	resultSelectedFile := widget.NewLabel("")
	optionsHorizontal := []string{"eng", "рус"}

	var fullPath string
	var fileName string
	var fileNameWithoutExt string
	// Создаем горизонтальную группу радио-кнопок
	langRadioGroupHorizontal := widget.NewRadioGroup(optionsHorizontal, func(selected string) {
		println("Выбран ответ:", selected)
	})
	langRadioGroupHorizontal.Horizontal = true
	langRadioGroupHorizontal.Selected = "eng"
	// Создаем список значений для радио-кнопок
	options := []string{"50", "500", "950", "1950", "3950", "4950"}

	// Создаем группу радио-кнопок
	var selectedSize string = "50"
	sizeRadioGroupHorizontal := widget.NewRadioGroup(options, func(selected string) {
		println("Выбрано значение:", selected)
		selectedSize = selected
	})
	sizeRadioGroupHorizontal.Horizontal = true
	sizeRadioGroupHorizontal.Selected = selectedSize

	// Создаем кнопку Старт
	startButton := widget.NewButton("Старт", func() {
		size := fmt.Sprintf("%sm", selectedSize)
		outputFile := fmt.Sprintf("%s.zip", fileNameWithoutExt)

		cmd := exec.Command("zip", "-s", size, outputFile, fullPath)

		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		loadingIndicator.Show()
		err := cmd.Run()
		loadingIndicator.Hide()
		var message string
		if err != nil {
			message = stderr.String()
			myApp.SendNotification(fyne.NewNotification("Ошибка:", message))
			notificationTitle.Text = "Ошибка:" + message + "\n" + err.Error()
			notificationTitle.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
			notificationTitle.Refresh()
			notificationTitle.Show()
			fmt.Printf("Ошибка: %s\n", stderr.String())
		} else {
			message = out.String()
			myApp.SendNotification(fyne.NewNotification("Успех:", message))
			notificationTitle.Text = "Успех:" + message
			notificationTitle.Color = color.RGBA{R: 0, G: 255, B: 0, A: 255}
			notificationTitle.Refresh()
			notificationTitle.Show()
			fmt.Printf("Успех: %s\n", out.String())
		}
	})
	startButton.Disable()

	// Создаем кнопку для открытия диалога выбора файла
	selectFileButton := widget.NewButton("Выбрать файл", func() {
		notificationTitle.Hide()
		// Открываем диалог выбора файла
		dialog.ShowFileOpen(func(uc fyne.URIReadCloser, err error) {
			if err == nil && uc != nil {
				fullPath = uc.URI().Path()

				// Получаем имя файла с расширением
				fileName = path.Base(fullPath)

				// Получаем имя файла без расширения
				ext := path.Ext(fileName)
				fileNameWithoutExt = strings.TrimSuffix(fileName, ext)
				// Выводим результаты
				println("Имя файла без расширения:", fileNameWithoutExt)
				println("Путь файла:", fullPath)
				resultSelectedFile.SetText(fullPath)
				if fullPath != "" && !strings.Contains(fileNameWithoutExt, " ") && !fileExists(fileNameWithoutExt+".zip") {
					startButton.Enable()
					notificationTitle.Hide()
				} else if strings.Contains(fileNameWithoutExt, " ") {
					notificationTitle.Text = "Убедитесь что имя вашего файла не содержит пробелов!"
					notificationTitle.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
					notificationTitle.Refresh()
					notificationTitle.Show()
					startButton.Disable()
				} else if fileExists(fileNameWithoutExt + ".zip") {
					notificationTitle.Text = "Удалите готовый zip файл прежде чем разбивать снова!"
					notificationTitle.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
					notificationTitle.Refresh()
					notificationTitle.Show()
					startButton.Disable()
				}

			}
		}, myWindow)
	})

	// Добавляем кнопки на окно
	myWindow.SetContent(container.NewVBox(
		container.NewHBox(
			widget.NewLabel("Выберите язык"), // Текст для горизонтальных радио кнопок
			langRadioGroupHorizontal,
		),
		container.NewHBox(
			container.NewCenter(labelSelectSize),
			sizeRadioGroupHorizontal,
		),
		container.NewHBox(
			labelSelectedFile,
			resultSelectedFile,
		),
		selectFileButton,
		startButton,
		container.NewHBox(notificationTitle),
		loadingIndicator,
	))

	// Показываем окно
	myWindow.ShowAndRun()
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
