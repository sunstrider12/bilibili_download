package download

var header =make(map[string]string)

func init() {
	header["user-agent"]="Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36"
	header["sec-ch-ua"]="Not A;Brand\";v=\"99\", \"Chromium\";v=\"90\", \"Google Chrome\";v=\"90\""
	header["sec-ch-ua-mobile"]="?0"
	header["sec-fetch-dest"]="document"
	header["sec-fetch-mode"]="navigate"
	header["sec-fetch-site"]="same-origin"
	header["sec-fetch-user"]="?1"
	header["upgrade-insecure-requests"]=""
	header["accept"]="text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	header["accept-language"]="zh"
	header["cache-control"]="max-age=0"
}