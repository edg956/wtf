package mail

import (
	"fmt"
	"github.com/rivo/tview"
	fakemail "github.com/wtfutil/wtf/modules/mail/client/fake-mail"
	"github.com/wtfutil/wtf/modules/mail/pages"
	"github.com/wtfutil/wtf/utils"
	"github.com/wtfutil/wtf/view"
)

// Widget is the container for your module's data
type Widget struct {
	view.ScrollableWidget

	settings *Settings
	pages    *stack
}

// NewWidget creates and returns an instance of Widget
func NewWidget(tviewApp *tview.Application, redrawChan chan bool, pages *tview.Pages, settings *Settings) *Widget {
	widget := Widget{
		ScrollableWidget: view.NewScrollableWidget(tviewApp, redrawChan, pages, settings.common),

		settings: settings,
	}

	widget.init()

	return &widget
}

/* -------------------- Exported Functions -------------------- */

// Refresh updates the onscreen contents of the widget
func (widget *Widget) Refresh() {

	// The last call should always be to the display function
	widget.display()
}

func (widget *Widget) HandleQ() {
	if widget.pages.size() > 1 {
		if _, err := widget.pages.pop(); err != nil {
			widget.initPages()
		}
		widget.SetItemCount(widget.currentPage().NumberOfItems())
		widget.display()
	}
}

func (widget *Widget) Select() {
	if currentPage := widget.currentPage(); currentPage != nil {
		if nextPage := currentPage.Select(widget.Selected); nextPage != nil {
			widget.moveToPage(nextPage)
			widget.display()
		}
	}
}

func (widget *Widget) RowColor(idx int) string {
	theme := widget.CommonSettings().Colors
	var color string
	if idx == widget.Selected {
		color = fmt.Sprintf("%s:%s", theme.RowTheme.HighlightedBackground, theme.RowTheme.EvenForeground)
	} else {
		color = fmt.Sprintf("-:-")
	}

	return color
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) currentPage() pages.IPage {
	current, err := widget.pages.peek()
	if err != nil || current == nil {
		return nil
	}

	return current.(pages.IPage)
}

func (widget *Widget) content() (string, string) {
	currentPage := widget.currentPage()
	if currentPage == nil {
		return "", "Something went wrong."
	}
	return currentPage.Render(widget.renderList)
}

func (widget *Widget) display() {
	title := widget.CommonSettings().Title

	titleExtra, content := widget.content()

	if titleExtra != "" {
		title = title + " - (" + titleExtra + ")"
	}

	widget.Redraw(func() (string, string, bool) {
		return title, content, false
	})
}

func (widget *Widget) init() {
	widget.initializeKeyboardControls()
	widget.initPages()
	widget.SetRenderFunction(widget.display)
}

func (widget *Widget) initPages() {
	client := fakemail.NewClient()
	widget.pages = newStack()
	widget.moveToPage(pages.StartPage(widget.View, client))
}

func (widget *Widget) moveToPage(page pages.IPage) {
	widget.pages.push(page)
	widget.SetItemCount(page.NumberOfItems())
}

func (widget *Widget) renderList(lines []string) string {
	if len(lines) == 0 {
		return " No items to display."
	}
	var result string
	for i, line := range lines {
		result += utils.HighlightableHelper(
			widget.View,
			fmt.Sprintf(
				"[%s] %s",
				widget.RowColor(i),
				line,
			),
			i,
			0,
		)
	}
	return result
}
