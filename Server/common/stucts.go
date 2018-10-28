package common

type (
	DbCfg struct {
		UsedDatabase string `json:"usedDatabase"`
		NameDatabase string `json:"nameDatabase"`
		User         string `json:"user"`
		Password     string `json:"password"`
		Host         string `json:"host"`
		Port         string `json:"port"`
	}

	ServCfg struct {
		Ip   string `json:"ip"`
		Port string `json:"port"`
	}

	DataStruct map[string]interface{}
	Response struct {
		Command string     `json:"command"`
		Status  bool       `json:"status"`
		Data    DataStruct `json:"data"`
	}

	StaffInfo struct {
		Id     int
		Login  string
		Tables []int
	}

	TableInfo struct {
		Name string
		Visitors []string
		Staff map[int]int
	}

	MenuCategory struct {
		Name 	string 		`json:"nameofcategory"`
		Goods 	[]MenuGood		`json:"goods"`
	}

	MenuGood struct {
		Id 		int 	`json:"id"`
		Name	string 	`json:"name"`
		Price 	int		`json:"price"`
	}
)
