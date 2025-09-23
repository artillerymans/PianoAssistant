package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/duke-git/lancet/v2/random"
	"github.com/duke-git/lancet/v2/strutil"
)

type MusicNote []Note

const lClef = `
      ██       
   ███ ███ ██  
   ███  ███    
    █   ██ ██  
       ███     
    ████       
   ██          

`
const hClef = `
       ███     
      ████     
      ███      
    ██████     
    ███████    
      ████     
      ████     

`

// 1. 定义 Model
type model struct {
	typeValue int
	count     int
	correct   int //正确数量
	//错误次数
	warningCount        int
	miss                int //没有选择
	listHeightMusicNote MusicNote
	listLowMusicNote    MusicNote
	index               int  //当前的Index
	showTip             bool //是否显示提示
	countDown           Countdown
	flagStyle           string //当前音符的样式
	autoCountDown       bool   //是否自动倒计时
}

const (
	LineTypeUp = iota
	LineTypeMidde
	LineTypeDonw

	/*高音*/
	TypeH
	/*低音*/
	TypeL
)

type Countdown struct {
	total    int
	current  int
	maxCount int
	minCount int
}

type Note struct {
	Tag      int //音标
	UiDes    string
	Selected bool
	lineType int
}

const HorLine = "-"

func drawHeightMusicNote() []Note {
	return []Note{
		{Tag: 4, UiDes: "", Selected: false, lineType: LineTypeUp},
		{Tag: 3, UiDes: HorLine, Selected: false, lineType: LineTypeUp},
		{Tag: 2, UiDes: "", Selected: false, lineType: LineTypeUp},
		{Tag: 1, UiDes: HorLine, Selected: false, lineType: LineTypeUp},
		{Tag: 7, UiDes: "", Selected: false, lineType: LineTypeUp},
		{Tag: 6, UiDes: HorLine, Selected: false, lineType: LineTypeUp},

		{Tag: 5, UiDes: "", Selected: false, lineType: LineTypeMidde},
		{Tag: 4, UiDes: HorLine, Selected: false, lineType: LineTypeMidde},
		{Tag: 3, UiDes: "", Selected: false, lineType: LineTypeMidde},
		{Tag: 2, UiDes: HorLine, Selected: false, lineType: LineTypeMidde},
		{Tag: 1, UiDes: "", Selected: false, lineType: LineTypeMidde},
		{Tag: 7, UiDes: HorLine, Selected: false, lineType: LineTypeMidde},
		{Tag: 6, UiDes: "", Selected: false, lineType: LineTypeMidde},
		{Tag: 5, UiDes: HorLine, Selected: false, lineType: LineTypeMidde},
		{Tag: 4, UiDes: "", Selected: false, lineType: LineTypeMidde},
		{Tag: 3, UiDes: HorLine, Selected: false, lineType: LineTypeMidde},
		{Tag: 2, UiDes: "", Selected: false, lineType: LineTypeMidde},

		{Tag: 1, UiDes: HorLine, Selected: false, lineType: LineTypeDonw},
		{Tag: 7, UiDes: "", Selected: false, lineType: LineTypeDonw},
		{Tag: 6, UiDes: HorLine, Selected: false, lineType: LineTypeDonw},
		{Tag: 5, UiDes: "", Selected: false, lineType: LineTypeDonw},
		{Tag: 4, UiDes: HorLine, Selected: false, lineType: LineTypeDonw},
		{Tag: 3, UiDes: "", Selected: false, lineType: LineTypeDonw},
	}
}

func drawLowMusicNote() []Note {
	return []Note{
		{Tag: 6, UiDes: "", Selected: false, lineType: LineTypeUp},
		{Tag: 5, UiDes: HorLine, Selected: false, lineType: LineTypeUp},
		{Tag: 4, UiDes: "", Selected: false, lineType: LineTypeUp},
		{Tag: 3, UiDes: HorLine, Selected: false, lineType: LineTypeUp},
		{Tag: 2, UiDes: "", Selected: false, lineType: LineTypeUp},
		{Tag: 1, UiDes: HorLine, Selected: false, lineType: LineTypeUp},

		{Tag: 7, UiDes: "", Selected: false, lineType: LineTypeMidde},
		{Tag: 6, UiDes: HorLine, Selected: false, lineType: LineTypeMidde},
		{Tag: 5, UiDes: "", Selected: false, lineType: LineTypeMidde},
		{Tag: 4, UiDes: HorLine, Selected: false, lineType: LineTypeMidde},
		{Tag: 3, UiDes: "", Selected: false, lineType: LineTypeMidde},
		{Tag: 2, UiDes: HorLine, Selected: false, lineType: LineTypeMidde},
		{Tag: 1, UiDes: "", Selected: false, lineType: LineTypeMidde},
		{Tag: 7, UiDes: HorLine, Selected: false, lineType: LineTypeMidde},
		{Tag: 6, UiDes: "", Selected: false, lineType: LineTypeMidde},
		{Tag: 5, UiDes: HorLine, Selected: false, lineType: LineTypeMidde},
		{Tag: 4, UiDes: "", Selected: false, lineType: LineTypeMidde},

		{Tag: 3, UiDes: HorLine, Selected: false, lineType: LineTypeDonw},
		{Tag: 2, UiDes: "", Selected: false, lineType: LineTypeDonw},
		{Tag: 1, UiDes: HorLine, Selected: false, lineType: LineTypeDonw},
		{Tag: 7, UiDes: "", Selected: false, lineType: LineTypeDonw},
		{Tag: 6, UiDes: HorLine, Selected: false, lineType: LineTypeDonw},
		{Tag: 5, UiDes: "", Selected: false, lineType: LineTypeDonw},
	}
}

// 2. 初始化（程序启动时调用一次）
func (m model) Init() tea.Cmd {
	if m.autoCountDown {
		return taskCount()
	} else {
		return nil
	}
}

func taskCount() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return t
	})
}

// 3. Update（处理消息 & 更新状态）
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg: // 键盘输入
		var key = msg.String()
		if key == "ctrl+c" || key == "q" || key == "esc" {
			return m, tea.Quit // 按 q 退出
		}

		if key == "v" {
			m.showTip = !m.showTip
		} else if key == "up" || key == "down" {
			if key == "down" {
				if m.countDown.total > m.countDown.minCount {
					m.countDown.total -= 1
					m.countDown.current = m.countDown.total
				}
			} else if key == "up" {
				if m.countDown.total < m.countDown.maxCount {
					m.countDown.total += 1
					m.countDown.current = m.countDown.total
				}
			}
		} else if key == "h" {
			if m.typeValue != TypeH {
				for index := range m.listLowMusicNote {
					m.listLowMusicNote[index].Selected = false
				}
				m.typeValue = TypeH
				//随机生成一个
				var index = random.RandInt(0, len(m.listHeightMusicNote))
				m.index = index
				m.listHeightMusicNote[m.index].Selected = true
				m.countDown.current = m.countDown.total
			}
		} else if key == "l" {
			if m.typeValue != TypeL {
				for index := range m.listHeightMusicNote {
					m.listHeightMusicNote[index].Selected = false
				}
				m.typeValue = TypeL
				//随机生成一个
				var index = random.RandInt(0, len(m.listLowMusicNote))
				m.index = index
				m.listLowMusicNote[m.index].Selected = true
				m.countDown.current = m.countDown.total
			}
		} else if key == "c" {
			var enable = !m.autoCountDown
			m.autoCountDown = enable
			var tempCountDown = m.countDown
			tempCountDown.current = tempCountDown.total
			m.countDown = tempCountDown
		} else {
			if m.typeValue == TypeH {
				var tempModel = m.listHeightMusicNote[m.index]
				var tag = strconv.Itoa(tempModel.Tag)
				if strings.EqualFold(tag, key) {
					/*清除掉原先的index选择*/
					//正确的 清空一下原来选中的数据
					m.listHeightMusicNote[m.index].Selected = false
					//随机生成一个
					var index = random.RandInt(0, len(m.listHeightMusicNote))
					m.index = index
					m.listHeightMusicNote[m.index].Selected = true
					m.count += 1
					m.correct += 1
					m.countDown.current = m.countDown.total
				} else {
					m.warningCount += 1
				}
			} else {
				var tempModel = m.listLowMusicNote[m.index]
				var tag = strconv.Itoa(tempModel.Tag)
				if strings.EqualFold(tag, key) {
					/*清除掉原先的index选择*/
					//正确的 清空一下原来选中的数据
					m.listLowMusicNote[m.index].Selected = false
					//随机生成一个
					var index = random.RandInt(0, len(m.listLowMusicNote))
					m.index = index
					m.listLowMusicNote[m.index].Selected = true
					m.count += 1
					m.correct += 1
					m.countDown.current = m.countDown.total
				} else {
					m.warningCount += 1

				}
			}

		}
	case time.Time:
		if m.autoCountDown {
			if m.countDown.current > 0 {
				m.countDown.current -= 1
			} else {
				m.miss += 1
				m.countDown.current = m.countDown.total
				m.count += 1
				if m.typeValue == TypeH {
					m.listHeightMusicNote[m.index].Selected = false
					//随机生成一个
					var index = random.RandInt(0, len(m.listHeightMusicNote))
					m.index = index
					m.listHeightMusicNote[m.index].Selected = true
				} else {
					m.listLowMusicNote[m.index].Selected = false
					//随机生成一个
					var index = random.RandInt(0, len(m.listLowMusicNote))
					m.index = index
					m.listLowMusicNote[m.index].Selected = true
				}
			}
		}
		return m, taskCount()
	}
	return m, nil
}

var titleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#00ADB5"))

var tipStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff2e63"))

var progressStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#eeeeee"))

var noteStyle = lipgloss.NewStyle().
	Align(lipgloss.Right).
	Foreground(lipgloss.Color("#aa96da")).
	Margin(1).
	Padding(1).
	Width(50)

var symbolStyle = lipgloss.NewStyle().
	Align(lipgloss.Right).
	Foreground(lipgloss.Color("#aa96da"))

var inputStyle = lipgloss.NewStyle().Align(lipgloss.Left).Foreground(lipgloss.Color("#ff2e63"))

func (m model) View() string {

	var buildStr = strings.Builder{}

	buildStr.WriteString(titleStyle.Render(fmt.Sprintln(strutil.Pad("速记助手", 140, "*"))))

	buildStr.WriteString(strings.Repeat("\n", 1))

	buildStr.WriteString(tipStyle.Render(fmt.Sprintln("总数: ", m.count, ", 正确:", m.correct, ", 错误:", m.warningCount, ",丢失:", m.miss)))

	buildStr.WriteString(strings.Repeat("\n", 1))

	if m.autoCountDown {
		buildStr.WriteString(progressStyle.Render(fmt.Sprintln(fmt.Sprintf(`计时:%d/%d`, m.countDown.current, m.countDown.total), strings.Repeat("=", m.countDown.current))))
		buildStr.WriteString(strings.Repeat("\n", 1))
	}

	var lBuilStr = m.listLowMusicNote.drawMusicNote(m.flagStyle, m.showTip)
	var hBuilStr = m.listHeightMusicNote.drawMusicNote(m.flagStyle, m.showTip)

	buildStr.WriteString(lipgloss.JoinHorizontal(
		lipgloss.Center,
		symbolStyle.Align(lipgloss.Center).Render(lClef),
		noteStyle.Align(lipgloss.Left).Render(lBuilStr.String()),
		symbolStyle.Align(lipgloss.Center).Render(hClef),
		noteStyle.Align(lipgloss.Right).Render(hBuilStr.String()),
	))

	buildStr.WriteString(strings.Repeat("\n", 1))

	buildStr.WriteString(inputStyle.Render("V 显示隐藏提示"))
	buildStr.WriteString("\n")
	buildStr.WriteString(inputStyle.Padding(1, 0).Render("方向 ⬆️  ⬇️  调节倒计时(默认10s)"))
	buildStr.WriteString("\n")
	buildStr.WriteString(inputStyle.Render("H 调整为高音区(默认) <-> L 调整为低音区"))
	buildStr.WriteString("\n")
	buildStr.WriteString(inputStyle.Render("C 关闭或者开启倒计时(默认开启); Q 退出当前应用"))

	return lipgloss.NewStyle().Padding(1, 1, 1, 1).Render(buildStr.String())
}

func (list MusicNote) drawMusicNote(flagStyle string, showTip bool) strings.Builder {
	var buildStr = strings.Builder{}
	for _, item := range list {
		var source string = ""
		var size int = 40
		if item.Selected {
			source = flagStyle
			size += 2
		} else {
			source = ""
		}
		/*上加 或者 下加线的长度 包括音符在内*/
		var lineTypeMinLength int = 10
		var tipTag string = ""
		if showTip && len(item.UiDes) > 0 {
			tipTag = strconv.Itoa(item.Tag)
		}
		switch item.lineType {
		case LineTypeUp:
			var lastNumber = size - lineTypeMinLength
			var lastStr = fmt.Sprint(strutil.Pad(source, lineTypeMinLength, item.UiDes))
			buildStr.WriteString(fmt.Sprintln(strutil.PadStart("", lastNumber, ""), lastStr, tipTag))
		case LineTypeMidde:
			buildStr.WriteString(fmt.Sprintln(strutil.Pad(source, size, item.UiDes), tipTag))
		case LineTypeDonw:
			var lastNumber = size - lineTypeMinLength
			var lastStr = fmt.Sprint(strutil.Pad(source, lineTypeMinLength, item.UiDes))
			buildStr.WriteString(fmt.Sprintln(lastStr, strutil.PadStart("", lastNumber, ""), tipTag))
		default:

		}
	}
	return buildStr
}

// 5. main 函数
func main() {
	// 创建程序
	var list = drawHeightMusicNote()

	var tempIndex = random.RandInt(0, len(list))
	list[tempIndex].Selected = true

	var lowList = drawLowMusicNote()

	p := tea.NewProgram(
		model{
			typeValue:           TypeH,
			flagStyle:           "🍏🐖", //"⭕️", //"🔴", "🎵"
			listHeightMusicNote: list,
			listLowMusicNote:    lowList,
			index:               tempIndex,
			count:               1,
			countDown: Countdown{
				current:  10,
				total:    10,
				maxCount: 40,
				minCount: 5,
			},
			autoCountDown: true,
		},
	)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
