package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/2captcha/2captcha-go"
	"github.com/b9uu/realty/internal/data"
)

func Scrape(realty data.RealtyInterface) error {
	url := "https://rentals.ca/phoenix/api/v1.0.2/listings?obj_path=montreal&details=mid2&suppress-pagination=1&limit=250"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	client := api2captcha.NewClient("087dfb6d3d678490a56ced3912d79874")

	cap := api2captcha.CloudflareTurnstile{
		SiteKey: "0x4AAAAAAADnPIDROrmt1Wwj",
		Url:     "https://rentals.ca",
	}
	d, dd, err := client.Solve(cap.ToRequest())
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(d, dd)
	return nil
	req.Header.Set("cookie", "csrftoken=HodYqpR7Y4brqFYmHxTgH6IYK8iZ4BP4; cf_clearance=wLqyM0MVO_sotLdSJ1b__fkkTlXdPlifU26znZETtsk-1724230202-1.2.1.1-HT_k.HTNFr6WRs897EOMyQaqIQduyU2HZelUTwv3K76PZN6JiEfTSj0ppYeAh1PemqH3LsbfXyYFHyqAAQ9HIxVXnwxEXf_IEFSMbGgEm.y9lwb8WYfNG9p4ZrL_JiPITFcENKfsHnhkadvHZFsu9B1Cn3qHZIhFAuirMLQeHt_bym1TXZ0mCkAX5CL4i6xbKKhnUfvEiAMykGVYbu_KD8r4fEDl_WoHCjTgKYD1iO4W3Z1fLe5pwjvi.Jl_FD_Sc0R7KrbsD94mxnURvaJHKMJNR.4zNcqr5vHW9LwJwkG4CHDM33.VypF8rTlahmKwYirHENNj7lAWjKXaDAVvWm5kBzyn9i3_6UOmG1UUBH0hvjksj7ZLMIetGLTeIHtpfgzgl6Okt4imTqWfD6FVQwDngajzQ2lDZZdO8LZYF5oRZtM3PRr3LIVAdM5WFPZU.eJiouxm1bVBwkco1NtnOd__F5QdJeKaKeWxSc55dJVADA_l17PFvVNnpcAHivO2; __cf_bm=ipG6fYCdeXDZsdiC16DYv7PKVqQVeGqKcIrZQwtpltI-1724230207-1.0.1.1-UFjV.34Ci_3a76RlkLVkKTHj.ZIgm.B3c50JLhFDVbDJ_nAJU1AogTPAk3WU9GUpBE5LwzosOuU1_SlCaQUPhg")
	req.Header.Add("authority", "rentals.ca")
	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("accept-language", "en-US,en;q=0.9")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("referer", "https://rentals.ca/vancouver")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-gpc", "1")
	req.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Add("x-csrftoken", "HodYqpR7Y4brqFYmHxTgH6IYK8iZ4BP4")
	req.Header.Add("x-rentalsapi-apikey", "None")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		fmt.Println(res.StatusCode)
		return fmt.Errorf("Got wrong status code: %d\n", res.StatusCode)
	}

	realties := Resp{}
	err = json.NewDecoder(res.Body).Decode(&realties)
	if err != nil {
		return err
	}
	for _, re := range realties.Data.Listings {
		err = realty.Insert(&re)
		if err != nil {
			fmt.Println(re.ID, re.Name, err)
		}
	}
	return nil
}

type Resp struct {
	Data struct {
		Listings []data.Realty `json:"listings"`
	} `json:"data"`
}
