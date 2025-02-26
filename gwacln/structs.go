package main

// Define the structure for yaml configuration file
type GridRowElement struct {
	Dropdown  Dropdown   `yaml:"dropdown,omitempty"`
	Button    Button     `yaml:"button,omitempty"`
	Form      Form       `yaml:"form,omitempty"`
	Input     Input      `yaml:"input,omitempty"`
	Slider    Slider     `yaml:"slider,omitempty"`
	Textarea  Textarea   `yaml:"textarea,omitempty"`
	Label     Label      `yaml:"label,omitempty"`
	H2        H2         `yaml:"h2,omitempty"`
	Paragraph *Paragraph `yaml:"p,omitempty"`
	Canvas    Canvas     `yaml:"canvas,omitempty"`
	Image     Image      `yaml:"image,omitempty"`
	Date      Date       `yaml:"date,omitempty"`
}

type Date struct {
	Id string `yaml:"id"`
}

type Paragraph struct {
	Id   string `yaml:"id"`
	Text string `yaml:"text"`
}

type Dropdown struct {
	Id         string   `yaml:"id"`
	DefaultInd int      `yaml:"defaultind"`
	Items      []string `yaml:"items"`
}

type Button struct {
	Id   string `yaml:"id"`
	Text string `yaml:"text"`
}

type Form struct {
	Id   string `yaml:"id"`
	Text string `yaml:"text"`
}

type Input struct {
	Id   string `yaml:"id"`
	Text string `yaml:"text"`
}

type Slider struct {
	Id        string `yaml:"id"`
	MinMaxIni []int  `yaml:"minmaxini"`
}

type Textarea struct {
	Id    string `yaml:"id"`
	Text  string `yaml:"text"`
	Lines int    `yaml:"lines"`
}

type Label struct {
	Id      string `yaml:"id"`
	Text    string `yaml:"text"`
	Mutable bool   `yaml:"mutable"`
}

type H2 struct {
	Id      string `yaml:"id"`
	Text    string `yaml:"text"`
	Mutable bool   `yaml:"mutable"`
}

type Canvas struct {
	Id     string `yaml:"id"`
	Width  int    `yaml:"width"`
	Height int    `yaml:"height"`
}

type Image struct {
	Id string `yaml:"id"`
}

// Define the structure for each grid row

type GuiDescr struct {
	Tab Tab `yaml:"tab"`
}

type Tab struct {
	Id   string `yaml:"id"`
	Text string `yaml:"text"`
	Row  []Row  `yaml:"rows"`
}

type Row struct {
	GridRow []GridRowElement `yaml:"gridrow"`
}
