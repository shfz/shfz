package fuzz

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/shfz/shfz/model"
)

var apiData []*model.Api
var fuzzQueue []*model.FuzzText
var totalRequestCount int

func GetApiData() []*model.Api {
	return apiData
}

func GetFuzz(param model.FuzzInfo) (string, error) {
	// fuzzQueueから取得
	if *param.IsGenetic {
		if len(fuzzQueue) != 0 {
			for i := 0; i < len(fuzzQueue); i++ {
				if fuzzQueue[i].Name == param.Name {
					fuzz := fuzzQueue[i].Text
					fuzzQueue = append(fuzzQueue[:i], fuzzQueue[i+1:]...)
					fmt.Println("Genetic", fuzz)
					return fuzz, nil
				}
			}
		} else {
			if err := genGeneticFuzz(param); err != nil {
				return "", err
			}
		}
	}
	// fuzzQueueにfuzzがない場合、ランダム文字列で生成
	text, err := genRandomFuzz(param)
	if err != nil {
		return "", err
	}
	fmt.Println("Random", text)
	return text, nil
}

func genRandomFuzz(param model.FuzzInfo) (string, error) {
	if len(param.Charset) == 0 {
		return "", errors.New("No CharacterSet")
	}
	rand.Seed(time.Now().UnixNano())
	length := rand.Intn(param.MaxLen-param.MinLen) + param.MinLen
	ret := ""
	for i := 0; i < length; i++ {
		j := rand.Intn(len(param.Charset))
		ret += param.Charset[j : j+1]
	}
	return ret, nil
}

func genGeneticFuzz(param model.FuzzInfo) error {
	var used []string
	minLen := 0
	for _, a := range apiData {
		for _, u := range a.UsedFuzzs {
			for _, f := range u.FuzzTexts.FuzzTexts {
				if f.Name == param.Name {
					if u.Framelen >= minLen {
						minLen = u.Framelen
						used = append(used, f.Text)
					}
				}
			}
		}
	}
	fmt.Println(minLen, len(used), used)
	for _, f1 := range used {
		for _, f2 := range used {
			c := model.FuzzText{}
			c.Name = param.Name
			c.Text = geneticAlgo(f1, f2, param)
			fuzzQueue = append(fuzzQueue, &c)
		}
	}
	return nil
}

func geneticAlgo(f1 string, f2 string, param model.FuzzInfo) string {
	rand.Seed(time.Now().UnixNano())
	length := rand.Intn(param.MaxLen-param.MinLen) + param.MinLen
	ret := ""
	for i := 0; i < length; i++ {
		// 突然変異
		k := rand.Intn(100)
		if k < 5 {
			t := rand.Intn(len(param.Charset))
			ret += param.Charset[t : t+1]
		} else if k < 20 { // indexをランダムした交叉
			j := rand.Intn(2)
			if j == 0 {
				s := rand.Intn(len(f1))
				ret += f1[s : s+1]
			} else {
				s := rand.Intn(len(f2))
				ret += f2[s : s+1]
			}
		} else { // indexを維持した交叉
			j := rand.Intn(2)
			if j == 0 {
				if len(f1) > i {
					ret += f1[i : i+1]
				} else {
					ret += f1[i%len(f1) : i%len(f1)+1]
				}
			} else {
				if len(f2) > i {
					ret += f2[i : i+1]
				} else {
					ret += f2[i%len(f2) : i%len(f2)+1]
				}
			}
		}
	}
	return ret
}

func SetFuzzTexts(param model.ApiParam) error {
	c := model.UsedFuzz{}
	c.ID = param.ID
	c.FuzzTexts = param.FuzzTexts
	// すでにAPIがある場合
	for i := range apiData {
		if apiData[i].Name == param.Name {
			apiData[i].UsedFuzzs = append(apiData[i].UsedFuzzs, &c)
			return nil
		}
	}
	// APIがない場合
	t := model.Api{}
	t.Name = param.Name
	t.UsedFuzzs = append(t.UsedFuzzs, &c)
	apiData = append(apiData, &t)
	return nil
}

func SetClientFeedback(id string, param model.ClientFeedback) error {
	for _, i := range apiData {
		for _, j := range i.UsedFuzzs {
			if j.ID == id {
				totalRequestCount += 1
				j.IsClientFeedbacked = true
				j.ClientFeedback = param
				return nil
			}
		}
	}
	return errors.New("SetClientFeedback : No id is found.")
}

func SetServerFeedback(id string, param model.ServerFeedback) error {
	for _, i := range apiData {
		for _, j := range i.UsedFuzzs {
			if j.ID == id {
				j.IsServerFeedbacked = true
				j.ServerFeedback = param
				return nil
			}
		}
	}
	return errors.New("SetServerFeedback : No id is found.")
}

func GenReport(param model.ReportReq) (string, error) {
	fmt.Println(totalRequestCount)
	res := "# shfz summary\n"
	res += fmt.Sprintf("GitHub Actions : [%s #%s](https://github.com/%s/actions/runs/%s)\n", param.Job, param.Number, param.Repo, param.RunID)
	for _, i := range apiData {
		res += fmt.Sprintf("## `%s`\n", i.Name)
		clientErrorCount := 0
		var ce []*model.ReportClientError
		var se []*model.ReportServerError
		for _, j := range i.UsedFuzzs {
			if j.IsClientFeedbacked {
				if *j.IsClientError {
					clientErrorCount += 1
					check := true
					for _, k := range ce {
						if k.ClientError == j.ClientError {
							k.Fuzz = append(k.Fuzz, &j.FuzzTexts)
							check = false
							break
						}
					}
					if check {
						t := model.ReportClientError{}
						t.ClientError = j.ClientError
						t.Fuzz = append(t.Fuzz, &j.FuzzTexts)
						ce = append(ce, &t)
					}
				}
			}
			if j.IsServerFeedbacked {
				if *j.IsServerError {
					check := true
					for _, k := range se {
						if k.ServerError == j.ServerError {
							k.Fuzz = append(k.Fuzz, &j.FuzzTexts)
							check = false
							break
						}
					}
					if check {
						t := model.ReportServerError{}
						t.ServerError = j.ServerError
						t.ServerErrorFile = j.ServerErrorFile
						t.ServerErrorLineNo = j.ServerErrorLineNo
						t.ServerErrorFunc = j.ServerErrorFunc
						t.Fuzz = append(t.Fuzz, &j.FuzzTexts)
						se = append(se, &t)
					}
				}
			}
		}
		per := float64(clientErrorCount) / float64(totalRequestCount) * 100
		res += fmt.Sprintf("\nError rate : %.1f%% (%d/%d)\n\n", per, clientErrorCount, totalRequestCount)
		// client error report
		for _, j := range ce {
			res += fmt.Sprintf("\n:no_entry_sign: `%s`\n", j.ClientError)
			total := len(j.Fuzz)
			count := 0
			for _, k := range j.Fuzz {
				if count < 5 {
					res += "-"
					for _, l := range k.FuzzTexts {
						res += fmt.Sprintf(" `%s` (%s)", l.Text, l.Name)
					}
					res += "\n"
					count += 1
				} else {
					res += fmt.Sprintf("- ...(+%d)\n", total-5)
					break
				}
			}
		}
		// server error report
		for _, j := range se {
			res += fmt.Sprintf("\n:warning: `%s`\n", j.ServerError)
			total := len(j.Fuzz)
			count := 0
			for _, k := range j.Fuzz {
				if count < 5 {
					res += "-"
					for _, l := range k.FuzzTexts {
						res += fmt.Sprintf(" `%s` (%s)", l.Text, l.Name)
					}
					res += "\n"
					count += 1
				} else {
					res += fmt.Sprintf("- ...(+%d)\n", total-5)
					break
				}
			}
			res += fmt.Sprintf("> `%s` is detected at `%s` Line %d in `%s` function.\n", j.ServerError, j.ServerErrorFile, j.ServerErrorLineNo, j.ServerErrorFunc)
			res += fmt.Sprintf("> https://github.com/%s/blob/%s%s%s#L%d\n", param.Repo, param.Hash, param.Path, j.ServerErrorFile, j.ServerErrorLineNo)
		}
	}
	res += "---\n"
	res += fmt.Sprintf("For more information, check [GitHub Actions Log and Artifact](https://github.com/%s/actions/runs/%s).\n", param.Repo, param.RunID)
	res += "\n"
	res += "Generated by [shfz](https://github.com/shfz).\n"
	return res, nil
}
