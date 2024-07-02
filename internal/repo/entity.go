package repo

import "time"

const rateDateLayout = `"2006-01-02T15:04:05"`

type RateDate time.Time

func (cd *RateDate) UnmarshalJSON(b []byte) error {
	t, err := time.Parse(rateDateLayout, string(b))
	if err != nil {
		return err
	}
	*cd = RateDate(t)
	return nil
}

func (cd RateDate) MarshalJSON() ([]byte, error) {
	t := time.Time(cd)
	return []byte(t.Format(rateDateLayout)), nil
}

type RateEntity struct {
	ID              *int     `json:"-"`
	CurID           int      `json:"Cur_ID"`
	RateDate        RateDate `json:"Date"`
	CurAbbreviation string   `json:"Cur_Abbreviation"`
	CurScale        int      `json:"Cur_Scale"`
	CurName         string   `json:"Cur_Name"`
	CurOfficialRate float64  `json:"Cur_OfficialRate"`
}
