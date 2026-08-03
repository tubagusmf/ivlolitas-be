package helper

func BuildProductSort(sort string) string {
	switch sort {
	case "name":
		return "name ASC"
	case "-name":
		return "name DESC"
	case "created_at":
		return "created_at ASC"
	case "-created_at":
		return "created_at DESC"
	case "price":
		return "price ASC"
	case "-price":
		return "price DESC"
	case "stock":
		return "stock ASC"
	case "-stock":
		return "stock DESC"
	case "sku":
		return "sku ASC"
	case "-sku":
		return "sku DESC"
	default:
		return "created_at DESC"
	}
}
